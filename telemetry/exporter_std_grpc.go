/**
 * @Author: steven
 * @Description:
 * @File: exporter_std_grpc
 * @Date: 18/12/23 00.48
 */

package telemetry

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	"net/url"
)

type otlpExporterGrpc struct {
	path    string
	host    string
	fullUrl string

	exp traceSdk.SpanExporter
}

func (o *otlpExporterGrpc) SpanProcessorType() SpanProcessor {
	return SpanProcessorSimple
}

func (o *otlpExporterGrpc) ExportSpans(ctx context.Context, spans []traceSdk.ReadOnlySpan) error {
	return o.exp.ExportSpans(ctx, spans)
}

func (o *otlpExporterGrpc) Shutdown(ctx context.Context) error {
	return o.exp.Shutdown(ctx)
}

func (o *otlpExporterGrpc) Init() (err error) {
	opts := make([]otlptracegrpc.Option, 0)
	if len(o.fullUrl) > 0 {
		if uri, errUri := url.Parse(o.fullUrl); errUri == nil {
			o.host = uri.Host
			o.path = uri.Path
			if uri.Scheme == "http" {
				opts = append(opts, otlptracegrpc.WithInsecure())
			}
		}
	}
	opts = append(opts, otlptracegrpc.WithCompressor("gzip"))
	if len(o.host) > 0 {
		opts = append(opts, otlptracegrpc.WithEndpoint(o.host))
	}
	o.exp, err = otlptracegrpc.New(context.Background(), opts...)
	return
}

func NewExporterStandardGrpcWithCustomURL(fullUrl string) Exporter {
	return &otlpExporterGrpc{fullUrl: fullUrl}
}
