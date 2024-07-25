/**
 * @Author: steven
 * @Description:
 * @File: config_loader
 * @Date: 25/07/24 08.04
 */

package config

import (
	"github.com/evorts/kevlars/rules/eval"
	"github.com/joho/godotenv"
	"os"
)

func (c *configManager) loadEnv() {
	// load .env file
	err := godotenv.Load(func() string {
		if v := os.Getenv("ENV_FILE"); len(v) > 0 {
			return v
		}
		return ".env"
	}())
	if err == nil {
		return
	}
	// if failed during load .env file -- there are 2 reasons: invalid format or file not exist
	// thus fallback to default values
	if v := os.Getenv(EnvKeyConfigLocalName.String()); len(v) < 1 {
		_ = os.Setenv(EnvKeyConfigLocalName.String(), defaultConfigName)
	}
	if v := os.Getenv(EnvKeyConfigLocalType.String()); len(v) < 1 {
		_ = os.Setenv(EnvKeyConfigLocalType.String(), TypeYaml.String())
	}
	if v := os.Getenv(EnvKeySecretLocalName.String()); len(v) < 1 {
		_ = os.Setenv(EnvKeySecretLocalName.String(), defaultSecretName)
	}
	if v := os.Getenv(EnvKeySecretLocalType.String()); len(v) < 1 {
		_ = os.Setenv(EnvKeySecretLocalType.String(), TypeYaml.String())
	}

	if v := os.Getenv(EnvKeyConfigRemoteName.String()); len(v) < 1 {
		_ = os.Setenv(EnvKeyConfigRemoteName.String(), defaultConfigName)
	}
	if v := os.Getenv(EnvKeyConfigRemoteType.String()); len(v) < 1 {
		_ = os.Setenv(EnvKeyConfigRemoteType.String(), TypeYaml.String())
	}
	if v := os.Getenv(EnvKeySecretRemoteName.String()); len(v) < 1 {
		_ = os.Setenv(EnvKeySecretRemoteName.String(), defaultSecretName)
	}
	if v := os.Getenv(EnvKeySecretRemoteType.String()); len(v) < 1 {
		_ = os.Setenv(EnvKeySecretRemoteType.String(), TypeYaml.String())
	}
}

func (c *configManager) populateEnvVars() envVars {
	ev := envVars{
		useConfig: useConfigVarItem{
			value: func() UseConfig {
				if v := os.Getenv(EnvKeyUseConfig.String()); len(v) > 0 {
					return UseConfig(v)
				}
				return UseConfigLocal
			}(),
			combination: func() UseConfigDynamicValueEnv {
				if v := os.Getenv(EnvKeyUseConfigDynamicValues.String()); len(v) > 0 {
					return UseConfigDynamicValueEnv(v)
				}
				return ""
			}(),
		},
		remote: remoteVarItem{
			providerName: func() CProvider {
				if v := os.Getenv(EnvKeyRemoteProvider.String()); len(v) > 0 {
					return CProvider(v)
				}
				return CProvider(RemoteProviderNone)
			}(),
			configAddress: os.Getenv(EnvKeyConfigRemoteAddress.String()),
			configName:    os.Getenv(EnvKeyConfigRemoteName.String()),
			configType: func() Type {
				if v := os.Getenv(EnvKeyConfigRemoteType.String()); len(v) > 0 {
					return Type(v)
				}
				return TypeYaml
			}(),
			secretAddress: os.Getenv(EnvKeySecretRemoteAddress.String()),
			secretName:    os.Getenv(EnvKeySecretRemoteName.String()),
			secretType: func() Type {
				if v := os.Getenv(EnvKeySecretRemoteType.String()); len(v) > 0 {
					return Type(v)
				}
				return TypeYaml
			}(),
		},
		local: localVarItem{
			providerName: func() CProvider {
				if v := os.Getenv(EnvKeyLocalProvider.String()); len(v) > 0 {
					return CProvider(v)
				}
				return CProvider(LocalProviderFile)
			}(),
			configName: os.Getenv(EnvKeyConfigLocalName.String()),
			configType: func() Type {
				if v := os.Getenv(EnvKeyConfigLocalType.String()); len(v) > 0 {
					return Type(v)
				}
				return TypeYaml
			}(),
			secretName: os.Getenv(EnvKeySecretLocalName.String()),
			secretType: func() Type {
				if v := os.Getenv(EnvKeySecretLocalType.String()); len(v) > 0 {
					return Type(v)
				}
				return TypeYaml
			}(),
		},
		dynamic: dynamicVarItem{
			localStringVars: localVarItem{
				providerName: CProvider(LocalProviderStringVar),
			},
			localFile: localVarItem{
				providerName: CProvider(LocalProviderFile),
			},
			remoteDB: remoteVarItem{
				providerName: CProvider(RemoteProviderDB),
			},
			remoteGSM: remoteVarItem{
				providerName: CProvider(RemoteProviderGSM),
			},
			remoteConsul: remoteVarItem{
				providerName: CProvider(RemoteProviderConsul),
			},
		},
	}
	// fill out dynamic value
	ev.dynamic.localStringVars = localVarItem{
		providerName: ev.dynamic.localStringVars.providerName,
		configName:   EnvKeyPatternConfigLocal.String(ev.dynamic.localStringVars.providerName.String(), EnvContextName),
		configType:   Type(EnvKeyPatternConfigLocal.String(ev.dynamic.localStringVars.providerName.String(), EnvContextType)),
		secretName:   EnvKeyPatternSecretLocal.String(ev.dynamic.localStringVars.providerName.String(), EnvContextName),
		secretType:   Type(EnvKeyPatternSecretLocal.String(ev.dynamic.localStringVars.providerName.String(), EnvContextType)),
	}
	ev.dynamic.localFile = localVarItem{
		providerName: ev.dynamic.localFile.providerName,
		configName:   EnvKeyPatternConfigLocal.String(ev.dynamic.localFile.providerName.String(), EnvContextName),
		configType:   Type(EnvKeyPatternConfigLocal.String(ev.dynamic.localFile.providerName.String(), EnvContextType)),
		secretName:   EnvKeyPatternSecretLocal.String(ev.dynamic.localFile.providerName.String(), EnvContextName),
		secretType:   Type(EnvKeyPatternSecretLocal.String(ev.dynamic.localFile.providerName.String(), EnvContextType)),
	}
	ev.dynamic.remoteDB = remoteVarItem{
		providerName:  ev.dynamic.remoteDB.providerName,
		configAddress: EnvKeyPatternConfigRemote.String(ev.dynamic.remoteDB.providerName.String(), EnvContextAddress),
		configName:    EnvKeyPatternConfigRemote.String(ev.dynamic.remoteDB.providerName.String(), EnvContextName),
		configType:    Type(EnvKeyPatternConfigRemote.String(ev.dynamic.remoteDB.providerName.String(), EnvContextType)),
		secretAddress: EnvKeyPatternSecretRemote.String(ev.dynamic.remoteDB.providerName.String(), EnvContextAddress),
		secretName:    EnvKeyPatternSecretRemote.String(ev.dynamic.remoteDB.providerName.String(), EnvContextName),
		secretType:    Type(EnvKeyPatternSecretRemote.String(ev.dynamic.remoteDB.providerName.String(), EnvContextType)),
	}
	ev.dynamic.remoteGSM = remoteVarItem{
		providerName:  ev.dynamic.remoteGSM.providerName,
		configAddress: EnvKeyPatternConfigRemote.String(ev.dynamic.remoteGSM.providerName.String(), EnvContextAddress),
		configName:    EnvKeyPatternConfigRemote.String(ev.dynamic.remoteGSM.providerName.String(), EnvContextName),
		configType:    Type(EnvKeyPatternConfigRemote.String(ev.dynamic.remoteGSM.providerName.String(), EnvContextType)),
		secretAddress: EnvKeyPatternSecretRemote.String(ev.dynamic.remoteGSM.providerName.String(), EnvContextAddress),
		secretName:    EnvKeyPatternSecretRemote.String(ev.dynamic.remoteGSM.providerName.String(), EnvContextName),
		secretType:    Type(EnvKeyPatternSecretRemote.String(ev.dynamic.remoteGSM.providerName.String(), EnvContextType)),
	}
	ev.dynamic.remoteConsul = remoteVarItem{
		providerName:  ev.dynamic.remoteConsul.providerName,
		configAddress: EnvKeyPatternConfigRemote.String(ev.dynamic.remoteConsul.providerName.String(), EnvContextAddress),
		configName:    EnvKeyPatternConfigRemote.String(ev.dynamic.remoteConsul.providerName.String(), EnvContextName),
		configType:    Type(EnvKeyPatternConfigRemote.String(ev.dynamic.remoteConsul.providerName.String(), EnvContextType)),
		secretAddress: EnvKeyPatternSecretRemote.String(ev.dynamic.remoteConsul.providerName.String(), EnvContextAddress),
		secretName:    EnvKeyPatternSecretRemote.String(ev.dynamic.remoteConsul.providerName.String(), EnvContextName),
		secretType:    Type(EnvKeyPatternSecretRemote.String(ev.dynamic.remoteConsul.providerName.String(), EnvContextType)),
	}
	return ev
}

func (c *configManager) loadRemoteProvider(provider RemoteProvider, items ...remoteProviderItem) []Provider {
	rs := make([]Provider, 0)
	for _, item := range items {
		switch true {
		case provider == RemoteProviderGSM:
			rs = append(rs, NewRemoteGSM(item.address, item.name, item.ctype))
		case provider == RemoteProviderDB:
			rs = append(rs, NewRemoteDB([]string{item.name}))
		default:
			rs = append(rs, NewRemote(provider.String(), item.address, item.name, item.ctype))
		}
	}
	return rs
}

func (c *configManager) loadLocalProvider(provider LocalProvider, items ...localProviderItem) []Provider {
	rs := make([]Provider, 0)
	for _, item := range items {
		switch true {
		case provider == LocalProviderStringVar:
			for _, sv := range item.stringVars {
				rs = append(rs, NewLocalStringVar(sv.Name, item.ctype, sv.Value))
			}
		default:
			rs = append(rs, NewLocalFile(item.name, item.ctype))
		}
	}
	return rs
}

func (c *configManager) loadDynamicProviders(ev envVars) []Provider {
	items := ev.useConfig.combination.Parse().Sort()
	secrets := make([]Provider, 0)
	configs := make([]Provider, 0)
	for _, item := range items {
		switch item.Provider {
		case CProvider(LocalProviderStringVar):
			// evaluate config from string vars
			if len(c.stringVars) < 1 {
				continue
			}
			for _, sv := range c.stringVars {
				configs = append(configs, NewLocalStringVar(sv.Name, sv.ValueType, sv.Value))
			}
		case CProvider(LocalProviderFile):
			if eval.AND(
				len(ev.dynamic.localFile.configName) > 0,
				len(ev.dynamic.localFile.configType) > 0,
			) {
				configs = append(configs, NewLocalFile(ev.dynamic.localFile.configName, ev.dynamic.localFile.configType))
			}
			if eval.AND(
				len(ev.dynamic.localFile.secretName) > 0,
				len(ev.dynamic.localFile.secretType) > 0,
			) {
				secrets = append(secrets, NewLocalFile(ev.dynamic.localFile.secretName, ev.dynamic.localFile.secretType))
			}
		case CProvider(RemoteProviderDB):
			// evaluate config from db
			if eval.AND(
				len(ev.dynamic.remoteDB.configName) > 0,
				len(ev.dynamic.remoteDB.configAddress) > 0,
			) {
				prefix := func() string {
					if v := os.Getenv(EnvKeyConfigRemotePrefix.String()); len(v) > 0 {
						return v
					}
					return ""
				}()
				configs = append(configs, NewRemoteDB([]string{prefix}))
			}
			// should never store secret in DB, thus will not be evaluated
		case CProvider(RemoteProviderGSM):
			if eval.AND(
				len(ev.dynamic.remoteGSM.configName) > 0,
				len(ev.dynamic.remoteGSM.configAddress) > 0,
			) {
				configs = append(configs, NewRemoteGSM(
					ev.dynamic.remoteGSM.configAddress,
					ev.dynamic.remoteGSM.configName,
					ev.dynamic.remoteGSM.configType,
				))
			}
			if eval.AND(
				len(ev.dynamic.remoteGSM.secretName) > 0,
				len(ev.dynamic.remoteGSM.secretAddress) > 0,
			) {
				secrets = append(secrets, NewRemoteGSM(
					ev.dynamic.remoteGSM.secretAddress,
					ev.dynamic.remoteGSM.secretName,
					ev.dynamic.remoteGSM.secretType,
				))
			}
		case CProvider(RemoteProviderConsul):
			if eval.AND(
				len(ev.dynamic.remoteConsul.configName) > 0,
				len(ev.dynamic.remoteConsul.configAddress) > 0,
			) {
				configs = append(configs, NewRemote(
					RemoteProviderConsul.String(),
					ev.dynamic.remoteConsul.configAddress,
					ev.dynamic.remoteConsul.configName,
					ev.dynamic.remoteConsul.configType,
				))
			}
			if eval.AND(
				len(ev.dynamic.remoteConsul.secretName) > 0,
				len(ev.dynamic.remoteConsul.secretAddress) > 0,
			) {
				secrets = append(secrets, NewRemote(
					RemoteProviderConsul.String(),
					ev.dynamic.remoteConsul.secretAddress,
					ev.dynamic.remoteConsul.secretName,
					ev.dynamic.remoteConsul.secretType,
				))
			}
		}
	}
	return append(configs, secrets...)
}
