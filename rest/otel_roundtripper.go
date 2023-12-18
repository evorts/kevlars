package rest

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
	otelTrace "go.opentelemetry.io/otel/trace"
	"net"
	"net/http"
	"strings"
	"time"
)

type otelRoundTripper struct {
	*transport
}

// RoundTrip executes a single HTTP transaction, returning a Response for the provided Request.
func (ort *otelRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if !ort.tm.Enabled() {
		return ort.parent.RoundTrip(request)
	}
	ctx := ort.extractCtx(request)
	attributes := ort.requestAttributes(request)
	var span otelTrace.Span
	ctx, span = ort.tm.Tracer().Start(ctx, "rest.call", otelTrace.WithSpanKind(otelTrace.SpanKindClient))
	defer func() {
		if span.IsRecording() {
			span.SetAttributes(attributes...)
		}
		span.End()
	}()

	ort.beforeHook(ctx, attributes, request)

	start := time.Now()
	response, err := ort.parent.RoundTrip(request)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		ort.errorHook(ctx, err, attributes)
		return response, err
	}

	attributes = ort.responseAttributes(attributes, response)
	ort.afterHook(ctx, duration, attributes)

	if ort.isRedirection(response) {
		ort.redirectHook(ctx, attributes)
		return response, err
	}

	if ort.isFailure(response) {
		ort.failureHook(ctx, attributes)
		return response, err
	}

	ort.successHook(ctx, attributes)
	return response, err
}

func (ort *otelRoundTripper) isFailure(response *http.Response) bool {
	if response == nil {
		return false
	}
	return response.StatusCode >= http.StatusBadRequest
}

func (ort *otelRoundTripper) isRedirection(response *http.Response) bool {
	if response == nil {
		return false
	}
	return response.StatusCode >= http.StatusMultipleChoices && response.StatusCode < http.StatusBadRequest
}

func (ort *otelRoundTripper) failureHook(
	ctx context.Context,
	attributes []attribute.KeyValue,
) {
	ort.metrics.inFlightCounter.Add(ctx, -1, metric.WithAttributes(attributes...))
	ort.metrics.failureCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
}

func (ort *otelRoundTripper) redirectHook(
	ctx context.Context,
	attributes []attribute.KeyValue,
) {
	ort.metrics.inFlightCounter.Add(ctx, -1, metric.WithAttributes(attributes...))
	ort.metrics.redirectCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
}

func (ort *otelRoundTripper) successHook(
	ctx context.Context,
	attributes []attribute.KeyValue,
) {
	ort.metrics.inFlightCounter.Add(ctx, -1, metric.WithAttributes(attributes...))
	ort.metrics.successesCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
}

func (ort *otelRoundTripper) beforeHook(
	ctx context.Context,
	attributes []attribute.KeyValue,
	request *http.Request,
) {
	ort.metrics.inFlightCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
	ort.metrics.attemptsCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
	if request == nil {
		ort.metrics.noRequestCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
	}
}

func (ort *otelRoundTripper) afterHook(
	ctx context.Context,
	duration int64,
	attributes []attribute.KeyValue,
) {
	ort.metrics.totalDurationCounter.Record(ctx, duration, metric.WithAttributes(attributes...))
}

func (ort *otelRoundTripper) responseAttributes(
	attributes []attribute.KeyValue,
	response *http.Response,
) []attribute.KeyValue {
	return append(
		append([]attribute.KeyValue{}, attributes...),
		ort.extractResponseAttributes(response)...,
	)
}

func (ort *otelRoundTripper) requestAttributes(request *http.Request) []attribute.KeyValue {
	return append(
		append(
			[]attribute.KeyValue{},
			ort.attributes...,
		),
		ort.extractRequestAttributes(request)...,
	)
}

func (ort *otelRoundTripper) errorHook(ctx context.Context, err error, attributes []attribute.KeyValue) {
	ort.metrics.inFlightCounter.Add(ctx, -1, metric.WithAttributes(attributes...))
	ort.metrics.errorsCounter.Add(ctx, 1, metric.WithAttributes(attributes...))

	var timeoutErr net.Error
	if errors.As(err, &timeoutErr) && timeoutErr.Timeout() {
		ort.metrics.timeoutsCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
	}

	if strings.HasSuffix(err.Error(), context.Canceled.Error()) {
		ort.metrics.canceledCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
	}
}

func (ort *otelRoundTripper) extractResponseAttributes(response *http.Response) []attribute.KeyValue {
	if response != nil {
		return []attribute.KeyValue{
			semConv.HTTPStatusCodeKey.Int(response.StatusCode),
			semConv.HTTPFlavorKey.String(response.Proto),
		}
	}
	return nil
}

func (ort *otelRoundTripper) extractRequestAttributes(request *http.Request) []attribute.KeyValue {
	if request != nil {
		reqIdHeader, reqIdValue := ort.getReqIdFromHeaders(request)
		return []attribute.KeyValue{
			semConv.HTTPURLKey.String(request.URL.String()),
			semConv.HTTPMethodKey.String(request.Method),
			semConv.HTTPUserAgentKey.String(request.Header.Get("User-Agent")),
			attribute.String("http.content_type", request.Header.Get("Content-Type")),
			attribute.String("http.accept", request.Header.Get("Accept")),
			attribute.String("http.request_id", reqIdValue),
			attribute.String("http.request_id_header", reqIdHeader),
		}
	}
	return nil
}

func (ort *otelRoundTripper) extractCtx(request *http.Request) context.Context {
	if request != nil && request.Context() != nil {
		return request.Context()
	}
	return context.Background()
}
