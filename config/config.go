/**
 * @Author: steven
 * @Description:
 * @File: config
 * @Date: 29/09/23 10.46
 */

package config

import (
	"github.com/evorts/kevlars/common"
	"github.com/spf13/viper"
	"time"
)

type Manager interface {
	Get(key string) interface{}
	GetBool(key string) bool
	GetBoolOrElse(key string, elseValue bool) bool
	GetFloat64(key string) float64
	GetFloat64OrElse(key string, elseValue float64) float64
	GetInt(key string) int
	GetIntOrElse(key string, elseValue int) int
	GetIntSlice(key string) []int
	GetIntSliceOrElse(key string, elseValue []int) []int
	GetString(key string) string
	GetStringOrElse(key string, elseValue string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapOrElse(key string, elseValue map[string]interface{}) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringOrElse(key string, elseValue map[string]string) map[string]string
	GetStringSlice(key string) []string
	GetStringSliceOrElse(key string, orElse []string) []string
	GetMapArray(key string) []map[string]interface{}
	GetMapArrayOrElse(key string, orElse []map[string]interface{}) []map[string]interface{}
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetDurationOrElse(key string, elseValue time.Duration) time.Duration
	UnmarshalTo(key string, to interface{}) error
	IsSet(key string) bool
	AllSettings() map[string]interface{}

	common.Init[Manager]
}

type Provider interface {
	Init() error
	GetData() map[string]interface{}
}

type configManager struct {
	providers  []Provider
	v          *viper.Viper
	stringVars StringVars
}

func (c *configManager) UnmarshalTo(key string, to interface{}) error {
	return c.v.UnmarshalKey(key, to)
}

func (c *configManager) GetBoolOrElse(key string, elseValue bool) bool {
	if v := c.GetBool(key); v {
		return v
	}
	return elseValue
}

func (c *configManager) GetFloat64OrElse(key string, elseValue float64) float64 {
	if v := c.GetFloat64(key); v > 0 {
		return v
	}
	return elseValue
}

func (c *configManager) GetIntOrElse(key string, elseValue int) int {
	if v := c.GetInt(key); v > 0 {
		return v
	}
	return elseValue
}

func (c *configManager) GetIntSliceOrElse(key string, elseValue []int) []int {
	if v := c.GetIntSlice(key); len(v) > 0 {
		return v
	}
	return elseValue
}

func (c *configManager) GetStringOrElse(key string, elseValue string) string {
	if v := c.GetString(key); len(v) > 0 {
		return v
	}
	return elseValue
}

func (c *configManager) GetStringMapOrElse(key string, elseValue map[string]interface{}) map[string]interface{} {
	if v := c.GetStringMap(key); v != nil {
		return v
	}
	return elseValue
}

func (c *configManager) GetStringMapStringOrElse(key string, elseValue map[string]string) map[string]string {
	if v := c.GetStringMapString(key); v != nil {
		return v
	}
	return elseValue
}

func (c *configManager) GetMapArray(key string) []map[string]interface{} {
	return c.GetMapArrayOrElse(key, make([]map[string]interface{}, 0))
}

func (c *configManager) GetMapArrayOrElse(key string, elseValue []map[string]interface{}) []map[string]interface{} {
	var arrMap []map[string]interface{}
	if err := c.UnmarshalTo(key, &arrMap); err != nil {
		return elseValue
	}
	return arrMap
}

func (c *configManager) Get(key string) interface{} {
	return c.v.Get(key)
}

func (c *configManager) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *configManager) GetFloat64(key string) float64 {
	return c.v.GetFloat64(key)
}

func (c *configManager) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *configManager) GetIntSlice(key string) []int {
	return c.v.GetIntSlice(key)
}

func (c *configManager) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *configManager) GetStringMap(key string) map[string]interface{} {
	return c.v.GetStringMap(key)
}

func (c *configManager) GetStringMapString(key string) map[string]string {
	return c.v.GetStringMapString(key)
}

func (c *configManager) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

func (c *configManager) GetStringSliceOrElse(key string, orElse []string) []string {
	if v := c.v.GetStringSlice(key); v != nil {
		return v
	}
	return orElse
}

func (c *configManager) GetTime(key string) time.Time {
	return c.v.GetTime(key)
}

func (c *configManager) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

func (c *configManager) GetDurationOrElse(key string, elseValue time.Duration) time.Duration {
	d := c.v.GetDuration(key)
	if d < 1 {
		return elseValue
	}
	return d
}

func (c *configManager) IsSet(key string) bool {
	return c.v.IsSet(key)
}

func (c *configManager) AllSettings() map[string]interface{} {
	return c.v.AllSettings()
}

func (c *configManager) Init() error {
	c.loadEnv()
	// container configuration in dev, staging or production should be first class citizen
	ev := c.populateEnvVars()
	switch true {
	case (len(ev.remote.configAddress) > 0 && len(ev.remote.secretAddress) > 0) ||
		ev.useConfig.value == UseConfigRemote:
		c.providers = append(
			c.providers, c.loadRemoteProvider(
				RemoteProvider(ev.remote.providerName),
				remoteProviderItem{
					address: ev.remote.configAddress,
					providerItem: &providerItem{
						name:  ev.remote.configName,
						ctype: ev.remote.configType,
					},
				},
				remoteProviderItem{
					address: ev.remote.secretAddress,
					providerItem: &providerItem{
						name:  ev.remote.secretName,
						ctype: ev.remote.secretType,
					},
				})...)
	case ev.useConfig.value == UseConfigDynamic:
		c.providers = append(c.providers, c.loadDynamicProviders(ev)...)

	// local config file and string var
	default:
		c.providers = append(
			c.providers, c.loadLocalProvider(
				LocalProvider(ev.local.providerName),
				localProviderItem{
					providerItem: &providerItem{
						name:  ev.local.configName,
						ctype: ev.local.configType,
					},
					stringVars: c.stringVars,
				},
				localProviderItem{
					providerItem: &providerItem{
						name:  ev.local.secretName,
						ctype: ev.local.secretType,
					},
				})...)
	}
	c.v = viper.New()
	for _, provider := range c.providers {
		if err := provider.Init(); err != nil {
			return err
		}
		if err := c.v.MergeConfigMap(provider.GetData()); err != nil {
			return err
		}
	}
	return nil
}

func (c *configManager) MustInit() Manager {
	if err := c.Init(); err != nil {
		panic(err)
	}
	return c
}

func New(opts ...common.Option[configManager]) Manager {
	c := &configManager{providers: make([]Provider, 0)}
	for _, opt := range opts {
		opt.Apply(c)
	}
	return c
}
