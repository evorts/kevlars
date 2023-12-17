/**
 * @Author: steven
 * @Description:
 * @File: exporter_zipkin
 * @Date: 18/12/23 00.49
 */

package telemetry

import (
	"context"
	"go.opentelemetry.io/otel/exporters/zipkin"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"os"
)

type zipkinExporter struct {
	url     string
	logName string

	exp traceSdk.SpanExporter
}

func (e *zipkinExporter) ExportSpans(ctx context.Context, spans []traceSdk.ReadOnlySpan) error {
	return e.exp.ExportSpans(ctx, spans)
}

func (e *zipkinExporter) Shutdown(ctx context.Context) error {
	return e.exp.Shutdown(ctx)
}

func (e *zipkinExporter) Init() (err error) {
	var logger = log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile)
	e.exp, err = zipkin.New(e.url, zipkin.WithLogger(logger))
	return err
}

func (e *zipkinExporter) SpanProcessorType() SpanProcessor {
	return SpanProcessorSimple
}

func NewZipkinExporter(url, logName string) Exporter {
	return &zipkinExporter{url: url, logName: logName}
}
