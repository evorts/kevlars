/**
 * @Author: steven
 * @Description:
 * @File: app_schedulers
 * @Date: 22/12/23 10.10
 */

package scaffold

import "github.com/evorts/kevlars/scheduler"

type IScheduler interface {
	WithSchedulers() IApplication
	HasScheduler() bool
	Scheduler(key string) scheduler.Manager
	DefaultScheduler() scheduler.Manager
}

func (app *Application) WithSchedulers() IApplication {
	return app
}

func (app *Application) HasScheduler() bool {
	return len(app.schedulers) > 0
}

func (app *Application) Scheduler(key string) scheduler.Manager {
	if v, ok := app.schedulers[key]; ok {
		return v
	}
	panic("scheduler with key " + key + " not found")
}

func (app *Application) DefaultScheduler() scheduler.Manager {
	return app.Scheduler(DefaultKey)
}
