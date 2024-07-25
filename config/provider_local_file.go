/**
 * @Author: steven
 * @Description:
 * @File: provider_local
 * @Date: 29/09/23 10.49
 */

package config

import "github.com/spf13/viper"

type localFileProvider struct {
	name  string
	cType Type
	paths []string

	v *viper.Viper
}

func (c *localFileProvider) GetData() map[string]interface{} {
	return c.v.AllSettings()
}

func (c *localFileProvider) Init() error {
	c.v = viper.New()
	c.v.SetConfigName(c.name)
	c.v.SetConfigType(c.cType.String())
	if len(c.paths) < 1 {
		c.paths = []string{"."}
	}
	for _, v := range c.paths {
		c.v.AddConfigPath(v)
	}
	return c.v.ReadInConfig()
}

func NewLocalFile(name string, ctype Type, paths ...string) Provider {
	return &localFileProvider{name: name, cType: ctype, paths: paths}
}
