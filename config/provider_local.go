/**
 * @Author: steven
 * @Description:
 * @File: provider_local
 * @Date: 29/09/23 10.49
 */

package config

import "github.com/spf13/viper"

type configLocal struct {
	localConfigName  string
	localConfigType  string
	localConfigPaths []string

	v *viper.Viper
}

func (c *configLocal) GetData() map[string]interface{} {
	return c.v.AllSettings()
}

func (c *configLocal) Init() error {
	c.v = viper.New()
	c.v.SetConfigName(c.localConfigName)
	c.v.SetConfigType(c.localConfigType)
	if len(c.localConfigPaths) < 1 {
		c.localConfigPaths = []string{"."}
	}
	for _, v := range c.localConfigPaths {
		c.v.AddConfigPath(v)
	}
	return c.v.ReadInConfig()
}

func NewConfigLocal(name, ctype string, paths ...string) Provider {
	return &configLocal{localConfigName: name, localConfigType: ctype, localConfigPaths: paths}
}
