/**
 * @Author: steven
 * @Description:
 * @File: helpers
 * @Date: 25/12/23 09.21
 */

package queue

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

func wrapTelemetry(ctx context.Context, tc trace.Tracer, spanName string, spanAttr []trace.SpanStartOption, task func(ctx context.Context)) {
	newCtx, span := tc.Start(ctx, spanName, spanAttr...)
	defer span.End()
	task(newCtx)
}

func wrapTelemetryTuple1[T any](ctx context.Context, tc trace.Tracer, spanName string, spanAttr []trace.SpanStartOption, task func(ctx context.Context) T) T {
	newCtx, span := tc.Start(ctx, spanName, spanAttr...)
	defer span.End()
	return task(newCtx)
}

func wrapTelemetryTuple2[T any, T2 any](ctx context.Context, tc trace.Tracer, spanName string, spanAttr []trace.SpanStartOption, task func(ctx context.Context) (T, T2)) (T, T2) {
	newCtx, span := tc.Start(ctx, spanName, spanAttr...)
	defer span.End()
	return task(newCtx)
}
