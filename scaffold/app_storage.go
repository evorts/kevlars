/**
 * @Author: steven
 * @Description:
 * @File: app_storage
 * @Date: 21/12/23 06.37
 */

package scaffold

import (
	"github.com/evorts/kevlars/cache"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/utils"
)

type IStorage interface {
	WithDatabases() IApplication
	DB(key string) db.Manager
	HasDB() bool
	DefaultDB() db.Manager

	WithCaches() IApplication
	Cache(ke string) cache.Manager
	HasCache() bool
	DefaultCache() cache.Manager
}

func (app *Application) WithDatabases() IApplication {
	// get configuration for multi database
	// expected result as follows:
	// {
	//	 "postgres":{"driver":"","dsn":"","telemetry_enabled":bool}
	//	 "mysql":{"driver":"","dsn":"","telemetry_enabled":bool}
	//	}
	dbs := app.config.GetStringMap("dbs")
	if len(dbs) < 1 {
		panic("there's no dbs configuration found")
	}
	for dbk, dbc := range dbs {
		dbcItem, ok := dbc.(map[string]interface{})
		if !ok {
			continue
		}
		driver, dsn, tmEnabled := "", "", true
		maxOpenConnection := 0
		maxIdleConnection := 0
		if v, exist := dbcItem["dsn"]; exist {
			dsn, _ = v.(string)
		}
		if v, exist := dbcItem["driver"]; exist {
			driver, _ = v.(string)
		}
		if v, exist := dbcItem["telemetry_enabled"]; exist {
			tmEnabled, _ = v.(bool)
		}
		if len(driver) < 1 && len(dsn) < 1 {
			continue
		}
		if v, exist := dbcItem["max_open_connection"]; exist {
			maxOpenConnection, _ = v.(int)
		}
		if v, exist := dbcItem["max_idle_connection"]; exist {
			maxIdleConnection, _ = v.(int)
		}
		opts := make([]db.Option, 0)
		if maxOpenConnection > 0 {
			opts = append(opts, db.WithMaxOpenConnection(maxOpenConnection))
		}
		if maxIdleConnection > 0 {
			opts = append(opts, db.WithMaxIdleConnection(maxIdleConnection))
		}
		app.dbs[dbk] = db.New(
			db.SupportedDriver(driver),
			dsn, tmEnabled, opts...,
		)
		utils.IfTrueThen(tmEnabled, func() {
			app.dbs[dbk].SetTelemetry(app.Telemetry())
		})
		app.dbs[dbk].MustConnect(app.startContext)
	}
	return app
}

func (app *Application) DB(key string) db.Manager {
	if v, ok := app.dbs[key]; ok {
		return v
	}
	panic("database with key " + key + " not found")
}

func (app *Application) HasDB() bool {
	return len(app.dbs) > 0
}

func (app *Application) DefaultDB() db.Manager {
	return app.DB(DefaultKey)
}

func (app *Application) WithCaches() IApplication {
	// get configuration for multi database
	// expected result as follows:
	// {
	//	 "default":{"enabled":"[value]","address":"","creds":"","db":[bool]},
	//	 "other":{"enabled":"[value]","address":"","creds":"","db":[bool]}
	//	}
	caches := app.config.GetStringMap("caches")
	if len(caches) < 1 {
		panic("there's no caches configuration found")
	}
	for ck, cv := range caches {
		cItem, ok := cv.(map[string]interface{})
		if !ok {
			continue
		}
		address, pass, dbIdx, enabled, tmEnabled := "", "", 0, true, false
		if v, exist := cItem["address"]; exist {
			address, _ = v.(string)
		}
		if v, exist := cItem["pass"]; exist {
			pass, _ = v.(string)
		}
		if v, exist := cItem["enabled"]; exist {
			enabled, _ = v.(bool)
		}
		if v, exist := cItem["telemetry_enabled"]; exist {
			tmEnabled, _ = v.(bool)
		}
		utils.IfE(enabled, func() {
			opts := make([]cache.Option, 0)
			utils.IfE(tmEnabled, func() {
				opts = append(opts, cache.WithTelemetry(app.Telemetry()))
			}, func() {
				opts = append(opts, cache.WithTelemetry(telemetry.NewNoop()))
			})
			app.caches[ck] = cache.NewRedis(address, pass, dbIdx, opts...)
		}, func() {
			app.caches[ck] = cache.NewNoop()
		})
		app.caches[ck].MustConnect(app.startContext)
	}
	return app
}

func (app *Application) Cache(key string) cache.Manager {
	if v, ok := app.caches[key]; ok {
		return v
	}
	panic("cache manager with key " + key + " not found")
}

func (app *Application) HasCache() bool {
	return len(app.caches) > 0
}

func (app *Application) DefaultCache() cache.Manager {
	return app.Cache(DefaultKey)
}
