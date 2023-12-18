/**
 * @Author: steven
 * @Description:
 * @File: tracer
 * @Date: 18/12/23 00.56
 */

package telemetry

import (
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func NewNoopTracer(name string) trace.Tracer {
	return noop.NewTracerProvider().Tracer(name)
}
