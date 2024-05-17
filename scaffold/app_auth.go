/**
 * @Author: steven
 * @Description:
 * @File: app_auth
 * @Date: 13/05/24 18.14
 */

package scaffold

import (
	"github.com/evorts/kevlars/auth"
)

type IAuth interface {
	WithClientAuthorization() IApplication

	ClientAuth() auth.ClientManager
}

func (app *Application) WithClientAuthorization() IApplication {
	// since this feature are tightly dependent with database then need to ensure database are instantiated
	app.WithDatabases()
	app.authClient = auth.NewClientManager(
		app.DefaultDB(),
		auth.ClientWithDatabaseRead(app.DefaultDBR()),
		auth.ClientWithLogger(app.Log()),
	)
	if enabled := app.Config().GetBool("auth.client.migrations.enabled"); enabled {
		app.authClient.AddOptions(
			auth.ClientWithExecuteMigration(
				enabled,
				app.Config().GetStringSliceOrElse("auth.client.migrations.dir", []string{})...,
			),
		)
	}
	if v := app.Config().GetBool("auth.client.lazy_load_data"); v {
		app.authClient.AddOptions(auth.ClientWithLazyLoadData(v))
	}
	app.authClient.MustInit()
	return app
}

func (app *Application) ClientAuth() auth.ClientManager {
	return app.authClient
}
