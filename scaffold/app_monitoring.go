/**
 * @Author: steven
 * @Description:
 * @File: app_monitoring
 * @Date: 20/12/23 23.48
 */

package scaffold

import (
	"github.com/evorts/kevlars/health"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/ts"
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
	logLevel := app.Config().GetInt("app.log.level")
	if logLevel < 1 || logLevel > 7 {
		logLevel = logger.LogLevelError.Id()
	}
	svc := []string{app.name}
	if len(app.env) > 0 {
		svc = append(svc, app.env)
	}
	if len(app.version) > 0 {
		svc = append(svc, app.version)
	}
	if len(app.scope) > 0 {
		svc = append(svc, app.scope)
	}
	opts := []logger.Option{
		logger.WithServiceName(strings.Join(svc, ".")),
	}
	if app.Config().GetBool("app.log.use_local_tz") {
		opts = append(opts, logger.WithTZTimeFormatter(ts.DefaultTimeZone))
	}
	app.log = logger.NewLogger(logger.LogLevel(logLevel), os.Stdout, opts...)
	return app
}

func (app *Application) withHealthcheck() IApplication {
	if app.Config().GetBool("app.healthcheck.health") {
		app.routes = append(app.routes,
			route{
				method:      http.MethodGet,
				path:        "/health",
				handlerEcho: app.healthEchoHandler,
			},
		)
	}
	if app.Config().GetBool("app.healthcheck.metrics") {
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
	if app.Config().GetBool("app.healthcheck.dependencies") {
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
		"db":    make([]deps, 0),
		"cache": make([]deps, 0),
	}
	utils.IfTrueThen(app.HasDB(), func() {
		for dbk, dbm := range app.dbs {
			result["db"] = append(result["db"], deps{
				name:   dbk,
				status: utils.IfER(dbm.Ping() == nil, func() string { return health.OK }, func() string { return health.NOK }),
			})
		}
	})
	utils.IfTrueThen(app.HasCache(), func() {
		for ck, cm := range app.caches {
			result["cache"] = append(result["cache"], deps{
				name:   ck,
				status: utils.IfER(cm.Ping() == nil, func() string { return health.OK }, func() string { return health.NOK }),
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
	return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}
