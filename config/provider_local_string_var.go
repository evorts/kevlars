/**
 * @Author: steven
 * @Description:
 * @File: provider_local
 * @Date: 29/09/23 10.49
 */

package config

import (
	"bytes"
	"github.com/spf13/viper"
)

type StringVar struct {
	Name      string
	Value     string
	ValueType Type
}

type StringVars []StringVar

type localStringVarProvider struct {
	name      string
	valueType Type
	value     string

	v *viper.Viper
}

func (c *localStringVarProvider) GetData() map[string]interface{} {
	return c.v.AllSettings()
}

func (c *localStringVarProvider) Init() error {
	c.v = viper.New()
	c.v.SetConfigName(c.name)
	c.v.SetConfigType(c.valueType.String())
	return c.v.ReadConfig(bytes.NewBuffer([]byte(c.value)))
}

func NewLocalStringVar(name string, valueType Type, value string) Provider {
	return &localStringVarProvider{name: name, valueType: valueType, value: value}
}
