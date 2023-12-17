/**
 * @Author: steven
 * @Description:
 * @File: exporter_std
 * @Date: 18/12/23 00.48
 */

package telemetry

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	"net/url"
)

type otlpExporter struct {
	path    string
	host    string
	fullUrl string

	exp traceSdk.SpanExporter
}

func (o *otlpExporter) SpanProcessorType() SpanProcessor {
	return SpanProcessorSimple
}

func (o *otlpExporter) ExportSpans(ctx context.Context, spans []traceSdk.ReadOnlySpan) error {
	return o.exp.ExportSpans(ctx, spans)
}

func (o *otlpExporter) Shutdown(ctx context.Context) error {
	return o.exp.Shutdown(ctx)
}

func (o *otlpExporter) Init() (err error) {
	opts := make([]otlptracehttp.Option, 0)
	if len(o.fullUrl) > 0 {
		if uri, errUri := url.Parse(o.fullUrl); errUri == nil {
			o.host = uri.Host
			o.path = uri.Path
			if uri.Scheme == "http" {
				opts = append(opts, otlptracehttp.WithInsecure())
			}
		}
	}
	opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.NoCompression))
	if len(o.host) > 0 {
		opts = append(opts, otlptracehttp.WithEndpoint(o.host))
	}
	if len(o.path) > 0 {
		opts = append(opts, otlptracehttp.WithURLPath(o.path))
	}
	client := otlptracehttp.NewClient(opts...)
	o.exp, err = otlptrace.New(context.Background(), client)
	return
}

func NewExporterStandardWithCustomHost(host string) Exporter {
	return &otlpExporter{
		host: host,
	}
}

func NewExporterStandardWithCustomPath(path string) Exporter {
	return &otlpExporter{path: path}
}

func NewExporterStandardWithCustomURL(fullUrl string) Exporter {
	return &otlpExporter{fullUrl: fullUrl}
}
