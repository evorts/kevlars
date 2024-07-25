/**
 * @Author: steven
 * @Description:
 * @File: config_option
 * @Date: 21/07/24 07.32
 */

package config

import "github.com/evorts/kevlars/common"

func WithStringVar(s StringVars) common.Option[configManager] {
	return common.OptionFunc[configManager](func(c *configManager) {
		c.stringVars = s
	})
}
