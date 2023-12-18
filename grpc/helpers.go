/**
 * @Author: steven
 * @Description:
 * @File: helpers
 * @Version: 1.0.0
 * @Date: 28/08/23 17.23
 */

package rpc

import (
	"context"
	"encoding/json"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/utils"
	"github.com/evorts/kevlars/vars"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"time"
)

type GrpcTimeoutCallOption struct {
	grpc.EmptyCallOption

	forcedTimeout time.Duration
}

func GrpcWithForcedTimeout(forceTimeout time.Duration) GrpcTimeoutCallOption {
	return GrpcTimeoutCallOption{forcedTimeout: forceTimeout}
}

func grpcGetForcedTimeout(callOptions []grpc.CallOption) (time.Duration, bool) {
	for _, opt := range callOptions {
		if co, ok := opt.(GrpcTimeoutCallOption); ok {
			return co.forcedTimeout, true
		}
	}
	return 0, false
}

func GrpcTimeoutInterceptor(t time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		timeout := t
		if v, ok := grpcGetForcedTimeout(opts); ok {
			timeout = v
		}
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func GrpcTraceInterceptor(tc trace.Tracer) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) error {
		newCtx, span := tc.Start(ctx, method, trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

func GrpcMetricInterceptor(metric telemetry.MetricsManager) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		iRes := invoker(ctx, method, req, reply, cc, opts...)
		metric.StartDefault(method).Push("grpc:" + utils.Iif(iRes == nil, "success", "error"))
		return iRes
	}
}

func GrpcLogRequestPayloadInterceptor(inJson bool, logWithProps func(props map[string]interface{}, messages ...interface{})) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		v := utils.IfER(inJson, func() interface{} {
			vv, _ := json.Marshal(req)
			return string(vv)
		}, func() interface{} {
			return req
		})
		logWithProps(map[string]interface{}{
			"method": method,
			"req_id": ctx.Value(vars.RequestContextId),
		}, v)
		iRes := invoker(ctx, method, req, reply, cc, opts...)
		return iRes
	}
}
