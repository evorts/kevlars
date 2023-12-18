package rest

import (
	"fmt"
	"go.opentelemetry.io/otel/metric"
	"net/http"
)

type Transport interface {
	BuildRoundTripper() http.RoundTripper
}

type transport struct {
	*otelConfig
	metrics httpMetric
}

func (t *transport) BuildRoundTripper() http.RoundTripper {
	return &otelRoundTripper{transport: &transport{
		otelConfig: t.otelConfig,
		metrics: httpMetric{
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
		},
	}}
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

func NewTransport(opts ...TransportOption) Transport {
	t := &transport{
		otelConfig: defaultOtelConfig(),
	}
	for _, opt := range opts {
		opt.apply(t.otelConfig)
	}
	return t
}
