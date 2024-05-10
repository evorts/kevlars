/**
 * @Author: steven
 * @Description:
 * @File: app_run
 * @Date: 22/12/23 10.05
 */

package scaffold

import (
	"context"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/midware"
	"github.com/evorts/kevlars/requests"
	"github.com/evorts/kevlars/utils"
	"github.com/evorts/kevlars/validation"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"time"
)

type IRun interface {
	Run(run func(a *Application))
	RunAsDaemon(run func(a *Application))
	RunUseEcho(run func(app *Application, e *echo.Echo))
	RunRestApiUseEcho(run func(app *Application, e *echo.Echo))
	RunGrpcServer(run func(app *Application, rpcServer *grpc.Server))
}

func (app *Application) RunAsDaemon(run func(a *Application)) {
	if app.DefaultScheduler() == nil {
		run(app)
		return
	}
	app.DefaultScheduler().WithTasks(func() {
		run(app)
	}).MustInit().StartBlocking()
}

func (app *Application) Run(run func(a *Application)) {
	run(app)
}

func (app *Application) RunUseEcho(run func(a *Application, e *echo.Echo)) {
	e := echo.New()
	e.HideBanner = true
	e.Debug = app.Config().GetInt("app.log.level") == logger.LogLevelDebug.Id()
	e.Logger.SetLevel(log.Lvl(logger.LogLevel(app.Config().GetInt("app.log.level")).EchoLogLevel()))
	// register predefined routes
	if len(app.routes) > 0 {
		for _, r := range app.routes {
			e.Add(r.method, r.path, r.handlerEcho)
		}
	}
	run(app, e)
	GracefulStopWithContext(
		app.Context(),
		app.PortRest(),
		app.gracefulTimeout,
		e, app.Log(),
	)
}

func (app *Application) RunRestApiUseEcho(run func(a *Application, e *echo.Echo)) {
	e := echo.New()
	e.HideBanner = true
	e.Debug = app.Config().GetInt("app.log.level") == logger.LogLevelDebug.Id()
	e.Validator = validation.NewValidator().MustInit()
	e.Logger.SetLevel(log.Lvl(logger.LogLevel(app.Config().GetInt("app.log.level")).EchoLogLevel()))
	e.Pre(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			reqId := ""
			if app.requestIdGenerator != nil {
				reqId = app.requestIdGenerator()
			}
			if len(reqId) < 1 {
				if uid, erUid := uuid.NewV4(); erUid == nil {
					reqId = uid.String()
				}
			}
			return reqId
		},
		RequestIDHandler: func(c echo.Context, reqId string) {
			// set client id in both echo context and native context
			if app.parseRequestId {
				if uid, err := uuid.FromString(reqId); err == nil {
					reqId = uid.String()
				}
			}
			req := c.Request()
			ctx := req.Context()
			c.Set(requests.ContextId.String(), reqId)
			newCtx := context.WithValue(ctx, requests.ContextId, reqId)
			newReq := req.WithContext(newCtx)
			c.SetRequest(newReq)
		},
		TargetHeader: echo.HeaderXRequestID,
	}))
	skipEndpoints := app.Config().GetStringSlice("telemetry.tracing.skip_endpoints")
	e.Use(
		midware.EchoLogger(app.Log(), skipEndpoints...),
		middleware.Recover(),
	)
	if app.Config().GetBool("app.timeout.enabled") {
		e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout: time.Duration(app.Config().GetIntOrElse("app.timeout.ms", 5000)) * time.Millisecond,
		}))
	}
	e.HTTPErrorHandler = midware.NewEchoHttpError(app.Telemetry(), app.Log()).Handler
	if app.Telemetry().Enabled() {
		echoMwOTel := midware.NewEchoTelemetryMiddleware(
			app.Telemetry(),
			midware.EchoWithTracerName(app.Config().GetString("telemetry.name")),
			midware.EchoWithBodyTrace(app.Config().GetBool("telemetry.body_tracing_enabled")),
			midware.EchoWithMaskedBodyFields(app.MaskedFields()...),
			midware.EchoWithMaskedHeaders(app.MaskedHeaders()...),
			midware.EchoWithSkipper(func(c echo.Context) bool {
				return utils.InArray(app.Config().GetStringSlice("telemetry.tracing.skip_endpoints"), c.Request().URL.Path)
			}),
		)
		e.Pre(echoMwOTel.EchoOTelPre())
		app.midwareTelemetry = echoMwOTel
	}
	if app.Metrics().Enabled() {
		e.Use(midware.EchoMetrics(
			midware.EchoMetricWithMetrics(app.Metrics()),
			midware.EchoMetricWithName(app.Name()+"_"+app.Scope()),
			midware.EchoMetricWithEnv(app.Env()),
		))
	}
	// register predefined routes
	if len(app.routes) > 0 {
		for _, r := range app.routes {
			e.Add(r.method, r.path, r.handlerEcho)
		}
	}
	run(app, e)
	GracefulStopWithContext(
		app.Context(),
		app.PortRest(),
		app.gracefulTimeout,
		e, app.Log(),
	)
}

func (app *Application) RunGrpcServer(run func(app *Application, rpcServer *grpc.Server)) {
	netListen, err := net.Listen("tcp", ":"+utils.IntToString(app.Config().GetIntOrElse("app.grpc_port", 8899)))
	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}
	server := grpc.NewServer()
	run(app, server)
	go func() {
		if err = server.Serve(netListen); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	server.GracefulStop()
}
