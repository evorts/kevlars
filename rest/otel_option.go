package rest

import (
	"github.com/evorts/kevlars/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"net/http"
	"strings"
)

type TransportOption interface {
	apply(t *otelConfig)
}

type transportOptionFunc func(t *otelConfig)

func (fn transportOptionFunc) apply(t *otelConfig) {
	fn(t)
}

// TransportWithParent sets the underlying http.RoundTripper which is wrapped by this round tripper.
// If the provided http.RoundTripper is nil, http.DefaultTransport will be used as the base http.RoundTripper
func TransportWithParent(parent http.RoundTripper) TransportOption {
	return transportOptionFunc(func(t *otelConfig) {
		if parent != nil {
			t.parent = parent
		}
	})
}

// TransportWithName sets the prefix for the metrics emitted by this round tripper.
// by default, the "http.client" name is used.
func TransportWithName(name string) TransportOption {
	return transportOptionFunc(func(t *otelConfig) {
		if strings.TrimSpace(name) != "" {
			t.name = strings.TrimSpace(name)
		}
	})
}

// TransportWithMeter sets the underlying  metric.Meter that is used to create metric instruments
// By default the no-op meter is used.
func TransportWithMeter(meter metric.Meter) TransportOption {
	return transportOptionFunc(func(t *otelConfig) {
		t.meter = meter
	})
}

// TransportWithTelemetry trace the request with assigned telemetry manager
func TransportWithTelemetry(tm telemetry.Manager) TransportOption {
	return transportOptionFunc(func(t *otelConfig) {
		t.tm = tm
	})
}

// TransportWithAttributes sets a list of attribute.KeyValue labels for all metrics associated with this round tripper
func TransportWithAttributes(attributes ...attribute.KeyValue) TransportOption {
	return transportOptionFunc(func(t *otelConfig) {
		t.attributes = attributes
	})
}

func TransportReqIdHeaders(headers ...string) TransportOption {
	return transportOptionFunc(func(t *otelConfig) {
		t.reqIdHeaders = headers
	})
}

func TransportAppendReqIdHeaders(headers ...string) TransportOption {
	return transportOptionFunc(func(t *otelConfig) {
		for _, header := range headers {
			t.reqIdHeaders = append(t.reqIdHeaders, header)
		}
	})
}
