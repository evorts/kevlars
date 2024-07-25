/**
 * @Author: steven
 * @Description:
 * @File: config_model
 * @Date: 21/07/24 21.38
 */

package config

type providerItem struct {
	name  string
	ctype Type
}
type remoteProviderItem struct {
	address string
	*providerItem
}

type localProviderItem struct {
	stringVars StringVars
	*providerItem
}

type envVarItem struct {
	value string
	*providerItem
}

type useConfigVarItem struct {
	value       UseConfig
	combination UseConfigDynamicValueEnv
}

type remoteVarItem struct {
	providerName  CProvider
	configAddress string
	configName    string
	configType    Type
	secretAddress string
	secretName    string
	secretType    Type
}

type localVarItem struct {
	providerName CProvider
	configName   string
	configType   Type
	secretName   string
	secretType   Type
}

type dynamicVarItem struct {
	localStringVars localVarItem
	localFile       localVarItem
	remoteDB        remoteVarItem
	remoteGSM       remoteVarItem
	remoteConsul    remoteVarItem
}

type envVars struct {
	useConfig useConfigVarItem
	remote    remoteVarItem
	local     localVarItem
	dynamic   dynamicVarItem
}
