/**
 * @Author: steven
 * @Description:
 * @File: provider_remote_db_option
 * @Date: 21/07/24 06.02
 */

package config

import "github.com/evorts/kevlars/common"

func WithDBTableName(tableName string) common.Option[dbProvider] {
	return common.OptionFunc[dbProvider](func(m *dbProvider) {
		m.tableName = tableName
	})
}
