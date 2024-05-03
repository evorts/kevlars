/**
 * @Author: steven
 * @Description:
 * @File: app_storage
 * @Date: 21/12/23 06.37
 */

package scaffold

import (
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/inmemory"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/utils"
)

type IStorage interface {
	WithDatabases() IApplication
	DB(key string) db.Manager
	HasDB() bool
	DefaultDB() db.Manager

	WithInMemories() IApplication
	InMemory(key string) inmemory.Manager
	HasInMemory() bool
	DefaultInMemory() inmemory.Manager
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
		utils.IfTrueThen(tmEnabled, func() {
			opts = append(opts, db.WithTelemetry(app.Telemetry()), db.WithTelemetryEnabled(tmEnabled))
		})
		app.dbs[dbk] = db.New(db.SupportedDriver(driver), dsn, opts...)
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

func (app *Application) WithInMemories() IApplication {
	// get configuration for multi database
	// expected result as follows:
	// {
	//	 "default":{"provider":"valkey","enabled":"[value]","address":"","creds":"","db":[bool]},
	//	 "other":{"provider":"redis","enabled":"[value]","address":"","creds":"","db":[bool]}
	//	}
	inMemories := app.config.GetStringMap("in_memory")
	if len(inMemories) < 1 {
		panic("there's no in memory configuration found")
	}
	for ck, cv := range inMemories {
		cItem, ok := cv.(map[string]interface{})
		if !ok {
			continue
		}
		provider, address, pass, dbIdx, enabled, tmEnabled := "", "", "", 0, true, false
		if v, exist := cItem["provider"]; exist {
			provider, _ = v.(string)
		}
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
		utils.IfTrueThen(!inmemory.ValidProvider(provider), func() {
			provider = inmemory.ProviderValKey.String()
		})
		utils.IfE(enabled, func() {
			utils.IfE(provider == inmemory.ProviderRedis.String(), func() {
				opts := append(
					inmemory.NewRedisOptions(),
					inmemory.RedisWithPassword(pass),
					inmemory.RedisWithDB(dbIdx),
				)
				utils.IfE(tmEnabled, func() {
					opts = append(opts, inmemory.RedisWithTelemetry(app.Telemetry()))
				}, func() {
					opts = append(opts, inmemory.RedisWithTelemetry(telemetry.NewNoop()))
				})
				app.in_memories[ck] = inmemory.NewRedis(
					address,
					opts...,
				)
			}, func() {
				opts := append(
					inmemory.NewValKeyOptions(),
					inmemory.ValKeyWithPassword(pass),
					inmemory.ValKeyWithDB(dbIdx),
				)
				utils.IfE(tmEnabled, func() {
					opts = append(opts, inmemory.ValKeyWithTelemetry(app.Telemetry()))
				}, func() {
					opts = append(opts, inmemory.ValKeyWithTelemetry(telemetry.NewNoop()))
				})
				app.in_memories[ck] = inmemory.NewValKey(
					address,
					opts...,
				)
			})
		}, func() {
			app.in_memories[ck] = inmemory.NewNoop()
		})
		app.in_memories[ck].MustConnect(app.startContext)
	}
	return app
}

func (app *Application) InMemory(key string) inmemory.Manager {
	if v, ok := app.in_memories[key]; ok {
		return v
	}
	panic("cache manager with key " + key + " not found")
}

func (app *Application) HasInMemory() bool {
	return len(app.in_memories) > 0
}

func (app *Application) DefaultInMemory() inmemory.Manager {
	return app.InMemory(DefaultKey)
}
