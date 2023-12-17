/**
 * @Author: steven
 * @Description:
 * @File: exporter_datadog
 * @Date: 18/12/23 00.48
 */

package telemetry

import (
	"context"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
)

type datadogExporter struct {
	address   string
	agentHost string
	agentPort string
}

func (d *datadogExporter) ExportSpans(ctx context.Context, spans []traceSdk.ReadOnlySpan) error {
	//TODO implement me
	panic("implement me")
}

func (d *datadogExporter) Shutdown(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d *datadogExporter) Init() error {
	//TODO implement me
	panic("implement me")
}

func (d *datadogExporter) SpanProcessorType() SpanProcessor {
	//TODO implement me
	panic("implement me")
}

func NewDatadogExporterUseEndpointUrl(address string) Exporter {
	return &datadogExporter{address: address}
}

func NewDatadogExporterUseAgent(host, port string) Exporter {
	return &datadogExporter{agentHost: host, agentPort: port}
}
