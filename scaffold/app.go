/**
 * @Author: steven
 * @Description:
 * @File: app
 * @Date: 20/12/23 23.28
 */

package scaffold

import (
	"context"
	"github.com/evorts/kevlars/algo"
	"github.com/evorts/kevlars/audit"
	"github.com/evorts/kevlars/cache"
	"github.com/evorts/kevlars/config"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/fflag"
	"github.com/evorts/kevlars/health"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/rest"
	"github.com/evorts/kevlars/rpc"
	"github.com/evorts/kevlars/scheduler"
	"github.com/evorts/kevlars/soap"
	"github.com/evorts/kevlars/telemetry"
	"os"
	"time"
)

type IApplication interface {
	Name() string
	Version() string
	Env() string
	Config() config.Manager

	IMonitoring
	IRemotes
	IRun
	IStorage
}

type Application struct {
	name            string
	version         string
	env             string
	scope           string
	config          config.Manager
	startContext    context.Context
	gracefulTimeout time.Duration // in seconds

	// monitoring
	log           logger.Manager
	tm            telemetry.Manager
	metrics       telemetry.MetricsManager
	audit         audit.Manager
	health        *health.Service
	maskedFields  []string
	maskedHeaders []string

	// remotes
	restClients map[string]rest.Manager
	grpcClients map[string]rpc.ClientManager
	soapClients map[string]soap.Manager

	// scheduler
	schedulers map[string]scheduler.Manager

	// storage
	dbs    map[string]db.Manager    // multi database
	caches map[string]cache.Manager // multi cache

	// algorithm
	similarity algo.SimilarityManager // provider => similarity

	// feature flag
	featureFlag fflag.Manager

	routes           []route // routes for rest endpoint
	midwareTelemetry interface{}

	requestIdGenerator func() string
	parseRequestId     bool // parseRequestId into uuidV4 in case it got cleanup during framework process
}

func (app *Application) Name() string {
	return app.name
}

func (app *Application) Scope() string {
	return app.scope
}

func (app *Application) Version() string {
	return app.version
}

func (app *Application) Env() string {
	return app.env
}

func (app *Application) Config() config.Manager {
	return app.config
}

func (app *Application) MaskedFields() []string {
	return app.maskedFields
}

func (app *Application) MaskedHeaders() []string {
	return app.maskedHeaders
}

func NewApp(opts ...Option) IApplication {
	app := &Application{
		caches:      make(map[string]cache.Manager),
		dbs:         make(map[string]db.Manager),
		restClients: make(map[string]rest.Manager),
		grpcClients: make(map[string]rpc.ClientManager),
		soapClients: make(map[string]soap.Manager),
	}
	app.config = config.New().MustInit()
	app.name = app.Config().GetString("app.name")
	app.version = app.Config().GetString("app.version")
	app.env = func() string {
		if v := os.Getenv("APP_ENV"); len(v) > 0 {
			return v
		}
		return app.Config().GetString("app.env")
	}()
	app.startContext = context.Background()
	app.gracefulTimeout = func() time.Duration {
		if v := app.Config().GetInt("app.graceful_timeout"); v > 0 {
			return time.Duration(v) * time.Second
		}
		return 5 * time.Second
	}()
	app.tm = telemetry.NewNoop()
	for _, opt := range opts {
		opt.apply(app)
	}
	app.withLogger()
	app.withHealthcheck()
	return app
}
