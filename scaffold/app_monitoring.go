/**
 * @Author: steven
 * @Description:
 * @File: app_monitoring
 * @Date: 20/12/23 23.48
 */

package scaffold

import (
	"github.com/evorts/kevlars/ctime"
	"github.com/evorts/kevlars/health"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

type IMonitoring interface {
	withLogger() IApplication
	withHealthcheck() IApplication

	Log() logger.Manager
	Telemetry() telemetry.Manager
	Metrics() telemetry.MetricsManager
	MiddlewareTelemetry() interface{}
}

func (app *Application) Log() logger.Manager {
	return app.log
}

func (app *Application) Telemetry() telemetry.Manager {
	return app.tm
}

func (app *Application) MiddlewareTelemetry() interface{} {
	return app.midwareTelemetry
}

func (app *Application) Metrics() telemetry.MetricsManager {
	return app.metrics
}

func (app *Application) withLogger() IApplication {
	// 0 = PANIC, 1 = FATAL, 2 = ERROR, 3 = WARN, 4 = INFO, 5 = DEBUG, 6 = TRACE
	logLevel := app.Config().GetInt(AppLogLevel.String())
	if !utils.NumberInRange(logLevel, logger.LogLevelPanic.Id(), logger.LogLevelOff.Id()) {
		logLevel = logger.LogLevelError.Id()
	}
	svc := []string{app.name}
	rules.WhenTrue(len(app.Env()) > 0, func() {
		svc = append(svc, app.env)
	})
	rules.WhenTrue(len(app.Version()) > 0, func() {
		svc = append(svc, app.version)
	})
	rules.WhenTrue(len(app.Scope()) > 0, func() {
		svc = append(svc, app.scope)
	})
	opts := []logger.Option{
		logger.WithServiceName(strings.Join(svc, ".")),
	}
	rules.WhenTrue(app.Config().GetBool(AppLogUseCustomTimezone.String()), func() {
		if v := app.Config().GetString(AppLogTimezone.String()); len(v) > 0 {
			opts = append(opts, logger.WithTZTimeFormatter(ctime.TimeZone(v)))
		}
	})
	app.log = logger.NewLogger(logger.LogLevel(logLevel), os.Stdout, opts...)
	return app
}

func (app *Application) withHealthcheck() IApplication {
	if app.Config().GetBool(AppHealth.String()) {
		app.routes = append(app.routes,
			route{
				method:      http.MethodGet,
				path:        "/health",
				handlerEcho: app.healthEchoHandler,
			},
		)
	}
	if app.Config().GetBool(AppMetric.String()) {
		app.health = health.New(
			health.Health{
				Version: app.version,
			},
			health.SystemUptime(),
			health.ProcessUptime(),
			health.SysInfoHealth(),
		)
		app.routes = append(app.routes,
			route{
				method:      http.MethodGet,
				path:        "/health/metrics",
				handlerEcho: app.healthInfoEchoHandler,
			},
		)
	}
	if app.Config().GetBool(AppDependencies.String()) {
		app.routes = append(app.routes, route{
			method:      http.MethodGet,
			path:        "/health/dependencies",
			handlerEcho: app.healthDependenciesEchoHandler,
		})
	}
	return app
}

// healthInfoEchoHandler godoc
// @Summary      System healthcheck
// @Description  Check the health status of the system where the app reside
// @Tags         monitoring
// @Accept       json
// @Produce      json
// @Success      200  {object}  health.Health
// @Router       /health/metrics [get]
func (app *Application) healthInfoEchoHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, app.health.Build())
}

// healthDependenciesEchoHandler godoc
// @Summary      System healthcheck
// @Description  Check the health status of the system where the app reside
// @Tags         monitoring
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health/dependencies [get]
func (app *Application) healthDependenciesEchoHandler(c echo.Context) error {
	type deps struct {
		name   string
		status string
	}
	result := map[string][]deps{
		DbKey.String():       make([]deps, 0),
		InMemoryKey.String(): make([]deps, 0),
	}
	rules.WhenTrue(app.HasDBS(), func() {
		for dbk, dbm := range app.dbs {
			result[DbKey.String()] = append(result[DbKey.String()], deps{
				name:   dbk,
				status: rules.WhenTrueRE1(dbm.Ping() == nil, func() string { return health.OK }, func() string { return health.NOK }),
			})
		}
	})
	rules.WhenTrue(app.HasInMemories(), func() {
		for ck, cm := range app.inMemories {
			result[InMemoryKey.String()] = append(result[InMemoryKey.String()], deps{
				name:   ck,
				status: rules.WhenTrueRE1(cm.Ping() == nil, func() string { return health.OK }, func() string { return health.NOK }),
			})
		}
	})
	return c.JSON(http.StatusOK, result)
}

// healthEchoHandler godoc
// @Summary      System ping
// @Description  Check reachable of the app by ping
// @Tags         monitoring
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]string
// @Router       /health [get]
func (app *Application) healthEchoHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{StatusKey.String(): "OK"})
}
