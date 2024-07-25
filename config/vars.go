/**
 * @Author: steven
 * @Description:
 * @File: vars
 * @Date: 18/12/23 11.12
 */

package config

import (
	"github.com/evorts/kevlars/utils"
	"sort"
	"strings"
)

type CProvider string

func (c CProvider) String() string {
	return string(c)
}

// Level the larger the significance
func (c CProvider) Level() int {
	switch c {
	case CProvider(RemoteProviderDB):
		return 20
	case CProvider(RemoteProviderConsul):
		return 21
	case CProvider(RemoteProviderGSM):
		return 21
	case CProvider(LocalProviderFile):
		return 11
	case CProvider(LocalProviderStringVar):
		return 10
	default:
		return 0
	}
}

type LocalProvider CProvider

const (
	LocalProviderFile      LocalProvider = "file"
	LocalProviderStringVar LocalProvider = "string"
	LocalProviderNone      LocalProvider = "none"
)

func (lp LocalProvider) String() string {
	return CProvider(lp).String()
}

type RemoteProvider string

const (
	RemoteProviderGSM    RemoteProvider = "gsm" // Google Secret Manager
	RemoteProviderConsul RemoteProvider = "consul"
	RemoteProviderDB     RemoteProvider = "db" // database
	RemoteProviderNone   RemoteProvider = "none"
)

func (rp RemoteProvider) String() string {
	return CProvider(rp).String()
}

type Type string

const (
	TypeYaml Type = "yaml"
	TypeJson Type = "json"
)

func (t Type) String() string {
	return string(t)
}

const (
	defaultTableName = "app_config"
)

type UseConfig string

const (
	UseConfigLocal   UseConfig = "local"
	UseConfigRemote  UseConfig = "remote"
	UseConfigDynamic UseConfig = "dynamic"
)

func (u UseConfig) String() string {
	return string(u)
}

type UseConfigDynamicValueProviderItem struct {
	ContextScope UseConfig
	Provider     CProvider
}

type UseConfigDynamicValueProviderItems []UseConfigDynamicValueProviderItem

func (i UseConfigDynamicValueProviderItems) Sort() UseConfigDynamicValueProviderItems {
	sort.Slice(i, func(a, b int) bool {
		return i[a].Provider.Level() > i[b].Provider.Level()
	})
	return i
}

type UseConfigDynamicValueEnv string

func (c UseConfigDynamicValueEnv) String() string {
	return string(c)
}
func (c UseConfigDynamicValueEnv) Parse() UseConfigDynamicValueProviderItems {
	ca := strings.Split(c.String(), ",")
	rs := make(UseConfigDynamicValueProviderItems, 0)
	for _, s := range ca {
		v := strings.Split(strings.TrimSpace(s), ".")
		// has to produce two segment or else skip
		if len(v) != 2 {
			continue
		}
		useConfigSegment := UseConfig(v[0])
		providerSegment := CProvider(v[1])
		// when use_config segment not registered then skip
		if !utils.InArray([]UseConfig{
			UseConfigRemote,
			UseConfigLocal,
		}, useConfigSegment) {
			continue
		}
		// when provider segment not registered then skip
		if useConfigSegment == UseConfigRemote && !utils.InArray([]CProvider{
			CProvider(RemoteProviderConsul),
			CProvider(RemoteProviderGSM),
			CProvider(RemoteProviderDB),
		}, providerSegment) {
			continue
		}
		if useConfigSegment == UseConfigLocal && !utils.InArray([]CProvider{
			CProvider(LocalProviderFile),
			CProvider(LocalProviderStringVar),
		}, providerSegment) {
			continue
		}
		rs = append(rs, UseConfigDynamicValueProviderItem{
			ContextScope: useConfigSegment,
			Provider:     providerSegment,
		})
	}
	return rs
}

const (
	EnvContextName    = "NAME"
	EnvContextAddress = "ADDR"
	EnvContextType    = "TYPE"
	EnvContextPrefix  = "PREFIX"

	defaultConfigName = "config.yaml"
	defaultSecretName = "secret.yaml"
)

type EnvKey string

const (
	EnvKeyLocalProvider   EnvKey = "LOCAL_PROVIDER"
	EnvKeyConfigLocalName EnvKey = "CONFIG_LOCAL_" + EnvContextName
	EnvKeyConfigLocalType EnvKey = "CONFIG_LOCAL_" + EnvContextType
	EnvKeySecretLocalName EnvKey = "SECRET_LOCAL_" + EnvContextName
	EnvKeySecretLocalType EnvKey = "SECRET_LOCAL_" + EnvContextType

	EnvKeyRemoteProvider      EnvKey = "REMOTE_PROVIDER"
	EnvKeyConfigRemoteName    EnvKey = "CONFIG_REMOTE_" + EnvContextName
	EnvKeyConfigRemoteType    EnvKey = "CONFIG_REMOTE_" + EnvContextType
	EnvKeyConfigRemoteAddress EnvKey = "CONFIG_REMOTE_" + EnvContextAddress
	EnvKeyConfigRemotePrefix  EnvKey = "CONFIG_REMOTE_" + EnvContextPrefix
	EnvKeySecretRemoteName    EnvKey = "SECRET_REMOTE_" + EnvContextName
	EnvKeySecretRemoteType    EnvKey = "SECRET_REMOTE_" + EnvContextType
	EnvKeySecretRemoteAddress EnvKey = "SECRET_REMOTE_" + EnvContextAddress

	EnvKeyUseConfig              EnvKey = "USE_CONFIG"
	EnvKeyUseConfigDynamicValues EnvKey = "USE_CONFIG_DYN_VALUES"
)

func (e EnvKey) String() string {
	return string(e)
}

type EnvKeyPattern string

const (
	EnvKeyPatternConfigLocal  EnvKeyPattern = "CONFIG_LOCAL_{Provider}_{Context}"
	EnvKeyPatternSecretLocal  EnvKeyPattern = "SECRET_LOCAL_{Provider}_{Context}"
	EnvKeyPatternConfigRemote EnvKeyPattern = "CONFIG_REMOTE_{Provider}_{Context}"
	EnvKeyPatternSecretRemote EnvKeyPattern = "SECRET_REMOTE_{Provider}_{Context}"
)

func (k EnvKeyPattern) String(provider, context string) string {
	v := strings.ReplaceAll(string(k), "{Provider}", provider)
	return strings.ReplaceAll(v, "{Context}", context)
}
