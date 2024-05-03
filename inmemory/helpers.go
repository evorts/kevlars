/**
 * @Author: steven
 * @Description:
 * @File: helpers
 * @Date: 03/05/24 07.17
 */

package inmemory

import (
	"context"
	"github.com/evorts/kevlars/utils"
	"go.opentelemetry.io/otel/trace"
)

func injectPrefixWhenDefined(prefix, key string) string {
	if len(prefix) < 1 {
		return key
	}
	return prefix + "_" + key
}

func injectPrefixIntoKeysWhenDefined(prefix string, keys ...string) []string {
	if len(prefix) < 1 {
		return keys
	}
	rs := make([]string, len(keys))
	for _, key := range keys {
		rs = append(rs, prefix+"_"+key)
	}
	return rs
}

func wrapTelemetryTuple1[T any](ctx context.Context, tc trace.Tracer, spanName string, spanAttr []trace.SpanStartOption, task func(ctx context.Context) T) T {
	newCtx, span := tc.Start(ctx, spanName, append(spanAttr, trace.WithSpanKind(trace.SpanKindClient))...)
	defer span.End()
	return task(newCtx)
}

func ValidProvider(v string) bool {
	return utils.InArray([]string{ProviderValKey.String(), ProviderRedis.String()}, v)
}
