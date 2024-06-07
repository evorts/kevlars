/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 13/05/24 16.14
 */

package auth

import (
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/inmemory"
	"github.com/evorts/kevlars/logger"
)

func ClientWithLogger(v logger.Manager) common.Option[clientManager] {
	return common.OptionFunc[clientManager](func(c *clientManager) {
		c.log = v
	})
}

func ClientWithInMemory(v inmemory.Manager) common.Option[clientManager] {
	return common.OptionFunc[clientManager](func(c *clientManager) {
		c.mem = v
	})
}

func ClientWithDatabaseRead(db db.Manager) common.Option[clientManager] {
	return common.OptionFunc[clientManager](func(c *clientManager) {
		c.dbr = db
	})
}

func ClientWithExecuteMigration(enabled bool, dir ...string) common.Option[clientManager] {
	return common.OptionFunc[clientManager](func(c *clientManager) {
		c.migrationDir = dir
		c.migrationEnabled = enabled
	})
}
