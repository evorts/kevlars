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
	"github.com/evorts/kevlars/logger"
	"io/fs"
)

func ClientWithLogger(v logger.Manager) common.Option[clientManager] {
	return common.OptionFunc[clientManager](func(c *clientManager) {
		c.log = v
	})
}

func ClientWithDatabaseRead(db db.Manager) common.Option[clientManager] {
	return common.OptionFunc[clientManager](func(c *clientManager) {
		c.dbr = db
	})
}

func ClientWithLazyLoadData(v bool) common.Option[clientManager] {
	return common.OptionFunc[clientManager](func(c *clientManager) {
		c.lazyLoad = v
	})
}

func ClientWithExecuteMigration(fs fs.FS, migrationTerms func() bool) common.Option[clientManager] {
	return common.OptionFunc[clientManager](func(c *clientManager) {
		c.fs = fs
		c.migrationTerms = migrationTerms
	})
}
