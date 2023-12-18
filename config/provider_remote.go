/**
 * @Author: steven
 * @Description:
 * @File: provider_remote
 * @Date: 29/09/23 10.49
 */

package config

import "github.com/spf13/viper"

type configRemote struct {
	provider   string
	address    string
	configPath string
	configType string

	v *viper.Viper
}

func (c *configRemote) Init() error {
	c.v = viper.New()
	if err := c.v.AddRemoteProvider(c.provider, c.address, c.configPath); err != nil {
		return err
	}
	c.v.SetConfigType(c.configType)
	return c.v.ReadRemoteConfig()
}

func (c *configRemote) GetData() map[string]interface{} {
	return c.v.AllSettings()
}

func NewRemote(provider, address, path string) Provider {
	return &configRemote{provider: provider, address: address, configPath: path}
}
