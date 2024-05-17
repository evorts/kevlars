/**
 * @Author: steven
 * @Description:
 * @File: app_features
 * @Date: 24/12/23 11.21
 */

package scaffold

import "github.com/evorts/kevlars/fflag"

type IFeatureFlag interface {
	WithFeatureFlag() IApplication

	FeatureFlag() fflag.Manager
}

func (app *Application) WithFeatureFlag() IApplication {
	// since this feature are tightly dependent with database then need to ensure database are instantiated
	app.WithDatabases()
	app.featureFlag = fflag.New(app.DefaultDB())
	_ = app.featureFlag.Init()
	return app
}

func (app *Application) FeatureFlag() fflag.Manager {
	return app.featureFlag
}
