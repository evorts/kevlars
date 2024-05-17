/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 17/05/24 08.38
 */

package fflag

import (
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/logger"
)

func WithDatabaseRead(db db.Manager) common.Option[manager] {
	return common.OptionFunc[manager](func(c *manager) {
		c.dbr = db
	})
}

func WithLogger(log logger.Manager) common.Option[manager] {
	return common.OptionFunc[manager](func(c *manager) {
		c.log = log
	})
}

func WithExecuteMigration(enabled bool, dir ...string) common.Option[manager] {
	return common.OptionFunc[manager](func(c *manager) {
		c.migrationDir = dir
		c.migrationEnabled = enabled
	})
}

func WithLazyLoadData(v bool) common.Option[manager] {
	return common.OptionFunc[manager](func(c *manager) {
		c.lazyLoadData = v
	})
}
