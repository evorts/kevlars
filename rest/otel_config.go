package rest

import (
	"github.com/evorts/kevlars/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"net/http"
)

type otelConfig struct {
	name         string
	parent       http.RoundTripper
	meter        metric.Meter
	attributes   []attribute.KeyValue
	tm           telemetry.Manager
	reqIdHeaders []string
}

var reqIdHeaders = []string{"X-Request-ID", "Idempotency-Key"}

func defaultOtelConfig() *otelConfig {
	return &otelConfig{
		name:         "http.client",
		parent:       http.DefaultTransport,
		meter:        otel.Meter("http.client"),
		reqIdHeaders: reqIdHeaders,
	}
}

func (t *otelConfig) getReqIdFromHeaders(req *http.Request) (string, string) {
	for _, header := range t.reqIdHeaders {
		if reqId := req.Header.Get(header); len(reqId) > 0 {
			return header, reqId
		}
	}
	return "", ""
}
