/**
 * @Author: steven
 * @Description:
 * @File: provider_gsm_option
 * @Date: 20/11/23 07.11
 */

package config

import "github.com/evorts/kevlars/common"

func WithGsmProjectId(projectId string) common.Option[gsmProvider] {
	return common.OptionFunc[gsmProvider](func(m *gsmProvider) {
		m.projectId = projectId
	})
}

func WithGsmResourceName(resourceName string) common.Option[gsmProvider] {
	return common.OptionFunc[gsmProvider](func(m *gsmProvider) {
		m.resourceName = resourceName
	})
}

func WithGsmConfigType(configType string) common.Option[gsmProvider] {
	return common.OptionFunc[gsmProvider](func(m *gsmProvider) {
		m.configType = configType
	})
}

func WithGsmUseJsonCredsFile(jsonCredsFile string) common.Option[gsmProvider] {
	return common.OptionFunc[gsmProvider](func(m *gsmProvider) {
		m.jsonCredFile = jsonCredsFile
	})
}

func WithGsmUseJsonCreds(jsonCreds []byte) common.Option[gsmProvider] {
	return common.OptionFunc[gsmProvider](func(m *gsmProvider) {
		m.jsonCred = jsonCreds
	})
}
