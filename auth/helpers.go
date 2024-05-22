/**
 * @Author: steven
 * @Description:
 * @File: helpers
 * @Date: 22/05/24 17.59
 */

package auth

import (
	"github.com/evorts/kevlars/db"
	"github.com/huandu/go-sqlbuilder"
)

func getFlavorByDriver(driver db.SupportedDriver) sqlbuilder.Flavor {
	if driver == db.DriverPostgreSQL {
		return sqlbuilder.PostgreSQL
	}
	panic("unsupported flavor")
}
