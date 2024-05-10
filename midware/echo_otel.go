package midware

import (
	"fmt"
	"github.com/evorts/kevlars/requests"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/utils"
	"github.com/evorts/kevlars/vars"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	otelAttr "go.opentelemetry.io/otel/attribute"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
	otelTrace "go.opentelemetry.io/otel/trace"
)

type EchoTelemetry interface {
	EchoOTelPre() echo.MiddlewareFunc
	EchoOTel() echo.MiddlewareFunc
}

type EchoOpt interface {
	apply(cfg *echoCfg)
}

type echoOpt func(config *echoCfg)

func (o echoOpt) apply(config *echoCfg) {
	o(config)
}

type echoCfg struct {
	name             string
	otm              telemetry.Manager
	skipper          middleware.Skipper
	withBodyTrace    bool
	maskedBodyFields map[string]int8
	maskedHeaders    map[string]int8
}

func trace(ec echo.Context, eCfg *echoCfg, body []byte) {
	req := ec.Request()
	ctx := req.Context()
	reqId := requests.IdEcho(ec)
	ecClientId := requests.ClientIdEcho(ec)
	clientId := rules.Iif(len(ecClientId) < 1, vars.ClientIdUnknown, ecClientId)
	//ctx := ec.propagators.Extract(savedCtx, propagation.HeaderCarrier(request.Header))
	opts := []otelTrace.SpanStartOption{
		//otelTrace.WithAttributes(semConv.NetAttributesFromHTTPRequest("tcp", request)...),
		//otelTrace.WithAttributes(semConv.EndUserAttributesFromHTTPRequest(request)...),
		//otelTrace.WithAttributes(semConv.HTTPClientAttributesFromHTTPRequest(request)...),
		//otelTrace.WithAttributes(semConv.HTTPServerMetricAttributesFromHTTPRequest(service, request)...),
		//otelTrace.WithAttributes(semConv.HTTPServerAttributesFromHTTPRequest(service, c.Path(), request)...),
		otelTrace.WithAttributes(otelAttr.String("header.content_type", req.Header.Get("Content-Type"))),
		otelTrace.WithAttributes(otelAttr.String("header.accept", req.Header.Get("Accept"))),
		otelTrace.WithAttributes(otelAttr.String("req.id", reqId)),
		otelTrace.WithAttributes(otelAttr.String("req.client_id", clientId)),
		otelTrace.WithSpanKind(otelTrace.SpanKindServer),
	}
	if eCfg.withBodyTrace {
		opts = append(
			opts,
			otelTrace.WithAttributes(
				telemetry.UseBytesMapAttributes(
					"req.body", body,
					eCfg.maskedBodyFields,
				).Exec()...,
			),
		)
	}
	spanName := ec.Path()
	if spanName == "" {
		spanName = fmt.Sprintf("HTTP %s route not found", req.Method)
	}
	newCtx, span := eCfg.otm.Tracer().Start(ctx, spanName, opts...)
	defer span.End()
	req = req.WithContext(newCtx)
	ec.SetRequest(req)

	attrs := semConv.HTTPAttributesFromHTTPStatusCode(ec.Response().Status)
	spanStatus, spanMessage := semConv.SpanStatusFromHTTPStatusCodeAndSpanKind(ec.Response().Status, otelTrace.SpanKindServer)
	span.SetAttributes(attrs...)
	span.SetStatus(spanStatus, spanMessage)
}

func (cfg *echoCfg) EchoOTelPre() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.skipper(c) {
				return next(c)
			}
			req := c.Request()
			ctx := req.Context()
			path := req.URL.Path
			// start parent span on incoming request
			newCtx, span := cfg.otm.Tracer().Start(
				ctx, path,
				otelTrace.WithSpanKind(otelTrace.SpanKindServer),
			)
			defer func() {
				span.End()
			}()
			req = req.WithContext(newCtx)
			c.SetRequest(req)
			return next(c)
		}
	}
}

func (cfg *echoCfg) EchoOTel() echo.MiddlewareFunc {
	if cfg.withBodyTrace {
		return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
			Skipper: cfg.skipper,
			Handler: func(c echo.Context, body []byte, response []byte) {
				trace(c, cfg, body)
			},
		})
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.skipper(c) {
				return next(c)
			}
			trace(c, cfg, nil)
			return next(c)
		}
	}
}

func EchoWithTracerName(name string) EchoOpt {
	return echoOpt(func(config *echoCfg) {
		config.name = name
	})
}

func EchoWithBodyTrace(bt bool) EchoOpt {
	return echoOpt(func(config *echoCfg) {
		config.withBodyTrace = bt
	})
}

func EchoWithMaskedBodyFields(mbf ...string) EchoOpt {
	return echoOpt(func(config *echoCfg) {
		config.maskedBodyFields = utils.ArrayToMapInt8(mbf)
	})
}

func EchoWithMaskedHeaders(mh ...string) EchoOpt {
	return echoOpt(func(config *echoCfg) {
		config.maskedHeaders = utils.ArrayToMapInt8(mh)
	})
}

func EchoWithSkipper(skipper middleware.Skipper) EchoOpt {
	return echoOpt(func(config *echoCfg) {
		config.skipper = skipper
	})
}

func EchoAppendMiddleware(echoTelemetry EchoTelemetry, mf ...echo.MiddlewareFunc) []echo.MiddlewareFunc {
	if echoTelemetry != nil {
		return append(mf, echoTelemetry.EchoOTel())
	}
	return mf
}

func NewEchoTelemetryMiddleware(otm telemetry.Manager, opts ...EchoOpt) EchoTelemetry {
	ec := &echoCfg{
		otm:              otm,
		maskedBodyFields: make(map[string]int8, 0),
		maskedHeaders:    make(map[string]int8, 0),
		withBodyTrace:    false,
	}
	for _, opt := range opts {
		opt.apply(ec)
	}
	return ec
}
