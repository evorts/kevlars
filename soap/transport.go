package soap

import (
	"context"
	"errors"
	"fmt"
	"github.com/evorts/kevlars/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelAttr "go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
	otelTrace "go.opentelemetry.io/otel/trace"
	"net"
	"net/http"
	"strings"
	"time"
)

type Transport interface {
	RoundTrip(request *http.Request) (*http.Response, error)
}

type transport struct {
	name       string
	parent     http.RoundTripper
	tm         telemetry.Manager
	meter      metric.Meter
	metrics    metrics
	attributes []attribute.KeyValue
}

func defaultTransportConfig() *transport {
	t := &transport{
		name:   "soap.client",
		parent: http.DefaultTransport,
		meter:  otel.Meter("soap.client"),
	}
	t.metrics = metrics{
		noRequestCounter:        t.mustCounter(t.meter.Int64Counter(fmt.Sprintf("%s.no_request", t.name))),
		errorsCounter:           t.mustCounter(t.meter.Int64Counter(fmt.Sprintf("%s.errors", t.name))),
		successesCounter:        t.mustCounter(t.meter.Int64Counter(fmt.Sprintf("%s.success", t.name))),
		timeoutsCounter:         t.mustCounter(t.meter.Int64Counter(fmt.Sprintf("%s.timeouts", t.name))),
		canceledCounter:         t.mustCounter(t.meter.Int64Counter(fmt.Sprintf("%s.cancelled", t.name))),
		deadlineExceededCounter: t.mustCounter(t.meter.Int64Counter(fmt.Sprintf("%s.deadline_exceeded", t.name))),
		totalDurationCounter:    t.mustHistogram(t.meter.Int64Histogram(fmt.Sprintf("%s.total_duration", t.name))),
		inFlightCounter:         t.mustUpDownCounter(t.meter.Int64UpDownCounter(fmt.Sprintf("%s.in_flight", t.name))),
		attemptsCounter:         t.mustCounter(t.meter.Int64Counter(fmt.Sprintf("%s.attemps", t.name))),
		failureCounter:          t.mustCounter(t.meter.Int64Counter(fmt.Sprintf("%s.failures", t.name))),
		redirectCounter:         t.mustCounter(t.meter.Int64Counter(fmt.Sprintf("%s.redirects", t.name))),
	}
	return t
}

func (t *transport) mustCounter(counter metric.Int64Counter, err error) metric.Int64Counter {
	if err != nil {
		panic(err)
	}
	return counter
}

func (t *transport) mustUpDownCounter(counter metric.Int64UpDownCounter, err error) metric.Int64UpDownCounter {
	if err != nil {
		panic(err)
	}
	return counter
}

func (t *transport) mustHistogram(histogram metric.Int64Histogram, err error) metric.Int64Histogram {
	if err != nil {
		panic(err)
	}
	return histogram
}

func (t *transport) RoundTrip(request *http.Request) (*http.Response, error) {
	if !t.tm.Enabled() {
		return t.parent.RoundTrip(request)
	}
	ctx := t.extractCtx(request)
	attributes := t.requestAttributes(request)
	action := ctx.Value("soap.action")
	if action != nil {
		if av, ok := action.(string); ok {
			attributes = append(attributes, otelAttr.String("req.action", av))
		}
	}

	var span otelTrace.Span
	ctx, span = t.tm.Tracer().Start(
		ctx, "soap.call",
		otelTrace.WithSpanKind(otelTrace.SpanKindClient),
	)
	defer func() {
		if span.IsRecording() {
			span.SetAttributes(attributes...)
		}
		span.End()
	}()

	t.beforeHook(ctx, attributes, request)

	start := time.Now()
	response, err := t.parent.RoundTrip(request)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		t.errorHook(ctx, err, attributes)
		return response, err
	}

	attributes = t.responseAttributes(attributes, response)
	t.afterHook(ctx, duration, attributes)

	if t.isRedirection(response) {
		t.redirectHook(ctx, attributes)
		return response, err
	}

	if t.isFailure(response) {
		t.failureHook(ctx, attributes)
		return response, err
	}

	t.successHook(ctx, attributes)
	return response, err
}

func (t *transport) isFailure(response *http.Response) bool {
	if response == nil {
		return false
	}
	return response.StatusCode >= http.StatusBadRequest
}

func (t *transport) isRedirection(response *http.Response) bool {
	if response == nil {
		return false
	}
	return response.StatusCode >= http.StatusMultipleChoices && response.StatusCode < http.StatusBadRequest
}

func (t *transport) failureHook(
	ctx context.Context,
	attributes []attribute.KeyValue,
) {
	t.metrics.inFlightCounter.Add(ctx, -1, metric.WithAttributes(attributes...))
	t.metrics.failureCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
}

func (t *transport) redirectHook(
	ctx context.Context,
	attributes []attribute.KeyValue,
) {
	t.metrics.inFlightCounter.Add(ctx, -1, metric.WithAttributes(attributes...))
	t.metrics.redirectCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
}

func (t *transport) successHook(
	ctx context.Context,
	attributes []attribute.KeyValue,
) {
	t.metrics.inFlightCounter.Add(ctx, -1, metric.WithAttributes(attributes...))
	t.metrics.successesCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
}

func (t *transport) beforeHook(
	ctx context.Context,
	attributes []attribute.KeyValue,
	request *http.Request,
) {
	t.metrics.inFlightCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
	t.metrics.attemptsCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
	if request == nil {
		t.metrics.noRequestCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
	}
}

func (t *transport) afterHook(
	ctx context.Context,
	duration int64,
	attributes []attribute.KeyValue,
) {
	t.metrics.totalDurationCounter.Record(ctx, duration, metric.WithAttributes(attributes...))
}

func (t *transport) responseAttributes(
	attributes []attribute.KeyValue,
	response *http.Response,
) []attribute.KeyValue {
	return append(
		append([]attribute.KeyValue{}, attributes...),
		t.extractResponseAttributes(response)...,
	)
}

func (t *transport) requestAttributes(request *http.Request) []attribute.KeyValue {
	return append(
		append(
			[]attribute.KeyValue{},
			t.attributes...,
		),
		t.extractRequestAttributes(request)...,
	)
}

func (t *transport) errorHook(ctx context.Context, err error, attributes []attribute.KeyValue) {
	t.metrics.inFlightCounter.Add(ctx, -1, metric.WithAttributes(attributes...))
	t.metrics.errorsCounter.Add(ctx, 1, metric.WithAttributes(attributes...))

	var timeoutErr net.Error
	if errors.As(err, &timeoutErr) && timeoutErr.Timeout() {
		t.metrics.timeoutsCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
	}

	if strings.HasSuffix(err.Error(), context.Canceled.Error()) {
		t.metrics.canceledCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
	}
}

func (t *transport) extractResponseAttributes(response *http.Response) []attribute.KeyValue {
	if response != nil {
		return []attribute.KeyValue{
			semConv.HTTPStatusCodeKey.Int(response.StatusCode),
			semConv.HTTPFlavorKey.String(response.Proto),
		}
	}
	return nil
}

func (t *transport) extractRequestAttributes(request *http.Request) []attribute.KeyValue {
	if request != nil {
		return []attribute.KeyValue{
			semConv.HTTPURLKey.String(request.URL.String()),
			semConv.HTTPMethodKey.String(request.Method),
		}
	}
	return nil
}

func (t *transport) extractCtx(request *http.Request) context.Context {
	if request != nil && request.Context() != nil {
		return request.Context()
	}
	return context.Background()
}

func NewTransport(opts ...TransportOption) Transport {
	t := defaultTransportConfig()
	for _, opt := range opts {
		opt.apply(t)
	}
	return t
}
