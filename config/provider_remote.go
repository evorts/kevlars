/**
 * @Author: steven
 * @Description:
 * @File: provider_remote
 * @Date: 29/09/23 10.49
 */

package config

import "github.com/spf13/viper"

type remoteProvider struct {
	provider string
	address  string
	path     string
	cType    Type

	v *viper.Viper
}

func (c *remoteProvider) Init() error {
	c.v = viper.New()
	if err := c.v.AddRemoteProvider(c.provider, c.address, c.path); err != nil {
		return err
	}
	c.v.SetConfigType(c.cType.String())
	return c.v.ReadRemoteConfig()
}

func (c *remoteProvider) GetData() map[string]interface{} {
	return c.v.AllSettings()
}

func NewRemote(provider, address, path string, configType Type) Provider {
	return &remoteProvider{provider: provider, address: address, path: path, cType: configType}
}
