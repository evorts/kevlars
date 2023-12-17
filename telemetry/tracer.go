/**
 * @Author: steven
 * @Description:
 * @File: tracer
 * @Date: 18/12/23 00.56
 */

package telemetry

import "go.opentelemetry.io/otel/trace"

func NewNoopTracer(name string) trace.Tracer {
	return trace.NewNoopTracerProvider().Tracer(name)
}
