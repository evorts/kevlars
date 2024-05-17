/**
 * @Author: steven
 * @Description:
 * @File: app_features
 * @Date: 24/12/23 11.21
 */

package scaffold

import (
	"github.com/evorts/kevlars/fflag"
)

type IFeatureFlag interface {
	WithFeatureFlag() IApplication

	FeatureFlag() fflag.Manager
}

func (app *Application) WithFeatureFlag() IApplication {
	// since this feature are tightly dependent with database then need to ensure database are instantiated
	app.WithDatabases()
	app.featureFlag = fflag.New(
		app.DefaultDB(),
		fflag.WithDatabaseRead(app.DefaultDBR()),
		fflag.WithLogger(app.Log()),
	)
	if enabled := app.Config().GetBool("feature_flag.migrations.enabled"); enabled {
		app.featureFlag.AddOptions(
			fflag.WithExecuteMigration(
				enabled,
				app.Config().GetStringSliceOrElse("feature_flag.migrations.dir", []string{})...,
			),
		)
	}
	if v := app.Config().GetBool("feature_flag.lazy_load_data"); v {
		app.featureFlag.AddOptions(fflag.WithLazyLoadData(v))
	}
	app.featureFlag.MustInit()
	return app
}

func (app *Application) FeatureFlag() fflag.Manager {
	return app.featureFlag
}
