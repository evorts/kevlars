/**
 * @Author: steven
 * @Description:
 * @File: config
 * @Date: 29/09/23 10.46
 */

package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
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
	GetMapArray(key string) []map[string]interface{}
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	UnmarshalTo(key string, to interface{}) error
	IsSet(key string) bool
	AllSettings() map[string]interface{}
	Init() error
	MustInit() Manager
}

type Provider interface {
	Init() error
	GetData() map[string]interface{}
}

type configManager struct {
	providers []Provider
	v         *viper.Viper
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
	var arrMap []map[string]interface{}
	if err := c.UnmarshalTo(key, &arrMap); err != nil {
		return make([]map[string]interface{}, 0)
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

func (c *configManager) GetTime(key string) time.Time {
	return c.v.GetTime(key)
}

func (c *configManager) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

func (c *configManager) IsSet(key string) bool {
	return c.v.IsSet(key)
}

func (c *configManager) AllSettings() map[string]interface{} {
	return c.v.AllSettings()
}

func (c *configManager) Init() error {
	// load .env file
	err := godotenv.Load(func() string {
		if v := os.Getenv("ENV_FILE"); len(v) > 0 {
			return v
		}
		return ".env"
	}())
	if err != nil {
		// if failed during load .env file -- there are 2 reasons: invalid format or file not exist
		// thus fallback to default values
		if v := os.Getenv("CONFIG_LOCAL_NAME"); len(v) < 1 {
			_ = os.Setenv("CONFIG_LOCAL_NAME", "config.yaml")
		}
		if v := os.Getenv("SECRET_LOCAL_NAME"); len(v) < 1 {
			_ = os.Setenv("SECRET_LOCAL_NAME", "secrets.yaml")
		}
		if v := os.Getenv("CONFIG_REMOTE_TYPE"); len(v) < 1 {
			_ = os.Setenv("CONFIG_REMOTE_TYPE", TypeYaml.String())
		}
		if v := os.Getenv("SECRET_REMOTE_TYPE"); len(v) < 1 {
			_ = os.Setenv("SECRET_REMOTE_TYPE", TypeYaml.String())
		}
	}
	// container configuration in dev, staging or production should be first class citizen
	envConfigRemoteAddr := os.Getenv("CONFIG_REMOTE_ADDR")
	envConfigRemoteName := os.Getenv("CONFIG_REMOTE_NAME")
	envConfigRemoteType := func() string {
		if v := os.Getenv("CONFIG_REMOTE_TYPE"); len(v) > 0 {
			return v
		}
		return TypeYaml.String()
	}()
	envSecretRemoteAddr := os.Getenv("SECRET_REMOTE_ADDR")
	envSecretRemoteName := os.Getenv("SECRET_REMOTE_NAME")
	envSecretRemoteType := func() string {
		if v := os.Getenv("SECRET_REMOTE_TYPE"); len(v) > 0 {
			return v
		}
		return TypeYaml.String()
	}()
	envRemoteProvider := func() string {
		if v := os.Getenv("REMOTE_PROVIDER"); len(v) > 0 {
			return v
		}
		return RemoteProviderNone.String()
	}()
	if (len(envConfigRemoteAddr) > 0 && len(envSecretRemoteAddr) > 0) || os.Getenv("USE_CONFIG") == "remote" {
		if envRemoteProvider == RemoteProviderGSM.String() {
			c.providers = append(c.providers,
				NewGoogleSecretManager(envConfigRemoteAddr, envConfigRemoteName, envConfigRemoteType),
				NewGoogleSecretManager(envSecretRemoteAddr, envSecretRemoteName, envSecretRemoteType),
			)
		} else {
			c.providers = append(
				c.providers,
				NewRemote(envRemoteProvider, envConfigRemoteAddr, envConfigRemoteType),
				NewRemote(envRemoteProvider, envSecretRemoteAddr, envSecretRemoteType),
			)
		}
	} else {
		c.providers = append(c.providers,
			NewConfigLocal(
				os.Getenv("CONFIG_LOCAL_NAME"),
				func() string {
					if v := os.Getenv("CONFIG_LOCAL_TYPE"); len(v) > 0 {
						return v
					}
					return TypeYaml.String()
				}(),
			),
			NewConfigLocal(
				os.Getenv("SECRET_LOCAL_NAME"),
				func() string {
					if v := os.Getenv("SECRET_LOCAL_TYPE"); len(v) > 0 {
						return v
					}
					return TypeYaml.String()
				}(),
			),
		)
	}
	c.v = viper.New()
	for _, provider := range c.providers {
		if err = provider.Init(); err != nil {
			return err
		}
		if err = c.v.MergeConfigMap(provider.GetData()); err != nil {
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

func New() Manager {
	return &configManager{providers: make([]Provider, 0)}
}
