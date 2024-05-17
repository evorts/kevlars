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
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/telemetry"
)

type IStorage interface {
	WithDatabases() IApplication
	DB(key string) db.Manager
	HasDBS() bool
	HasDB(key string) bool
	HasDefaultDB() bool
	DefaultDB() db.Manager
	DefaultDBR() db.Manager

	WithInMemories() IApplication
	InMemory(key string) inmemory.Manager
	HasInMemory() bool
	DefaultInMemory() inmemory.Manager
}

func (app *Application) WithDatabases() IApplication {
	// no need to reinitiate when its already initiate
	if app.HasDBS() {
		return app
	}
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
		rules.WhenTrue(tmEnabled, func() {
			opts = append(opts, db.WithTelemetry(app.Telemetry()), db.WithTelemetryEnabled(tmEnabled))
		})
		app.dbs[dbk] = db.New(db.SupportedDriver(driver), dsn, opts...)
		app.dbs[dbk].MustConnect(app.startContext)
	}
	if !app.HasDB(DefaultKey) {
		panic("please define default database")
	}
	return app
}

func (app *Application) DB(key string) db.Manager {
	if v, ok := app.dbs[key]; ok {
		return v
	}
	panic("database with key " + key + " not found")
}

func (app *Application) HasDBS() bool {
	return len(app.dbs) > 0
}

func (app *Application) HasDB(key string) bool {
	if _, ok := app.dbs[key]; ok {
		return true
	}
	return false
}

func (app *Application) HasDefaultDB() bool {
	return app.HasDB(DefaultKey)
}

func (app *Application) DefaultDB() db.Manager {
	return app.DB(DefaultKey)
}

func (app *Application) DefaultDBR() db.Manager {
	if app.HasDB(DefaultKey + "_read") {
		return app.DB(DefaultKey + "_read")
	}
	return app.DefaultDB()
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
		rules.WhenTrue(!inmemory.ValidProvider(provider), func() {
			provider = inmemory.ProviderValKey.String()
		})
		rules.WhenTrueE(enabled, func() {
			rules.WhenTrueE(provider == inmemory.ProviderRedis.String(), func() {
				opts := append(
					inmemory.NewRedisOptions(),
					inmemory.RedisWithPassword(pass),
					inmemory.RedisWithDB(dbIdx),
				)
				rules.WhenTrueE(tmEnabled, func() {
					opts = append(opts, inmemory.RedisWithTelemetry(app.Telemetry()))
				}, func() {
					opts = append(opts, inmemory.RedisWithTelemetry(telemetry.NewNoop()))
				})
				app.inMemories[ck] = inmemory.NewRedis(
					address,
					opts...,
				)
			}, func() {
				opts := append(
					inmemory.NewValKeyOptions(),
					inmemory.ValKeyWithPassword(pass),
					inmemory.ValKeyWithDB(dbIdx),
				)
				rules.WhenTrueE(tmEnabled, func() {
					opts = append(opts, inmemory.ValKeyWithTelemetry(app.Telemetry()))
				}, func() {
					opts = append(opts, inmemory.ValKeyWithTelemetry(telemetry.NewNoop()))
				})
				app.inMemories[ck] = inmemory.NewValKey(
					address,
					opts...,
				)
			})
		}, func() {
			app.inMemories[ck] = inmemory.NewNoop()
		})
		app.inMemories[ck].MustConnect(app.startContext)
	}
	return app
}

func (app *Application) InMemory(key string) inmemory.Manager {
	if v, ok := app.inMemories[key]; ok {
		return v
	}
	panic("cache manager with key " + key + " not found")
}

func (app *Application) HasInMemory() bool {
	return len(app.inMemories) > 0
}

func (app *Application) DefaultInMemory() inmemory.Manager {
	return app.InMemory(DefaultKey)
}
