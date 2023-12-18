/**
 * @Author: steven
 * @Description:
 * @File: provider_gsm_option
 * @Date: 20/11/23 07.11
 */

package config

type GsmOption interface {
	apply(m *googleSecret)
}

type gsmOptionFunc func(m *googleSecret)

func (f gsmOptionFunc) apply(m *googleSecret) {
	f(m)
}

func WithGsmProjectId(projectId string) GsmOption {
	return gsmOptionFunc(func(m *googleSecret) {
		m.projectId = projectId
	})
}

func WithGsmResourceName(resourceName string) GsmOption {
	return gsmOptionFunc(func(m *googleSecret) {
		m.resourceName = resourceName
	})
}

func WithGsmConfigType(configType string) GsmOption {
	return gsmOptionFunc(func(m *googleSecret) {
		m.configType = configType
	})
}

func WithGsmUseJsonCredsFile(jsonCredsFile string) GsmOption {
	return gsmOptionFunc(func(m *googleSecret) {
		m.jsonCredFile = jsonCredsFile
	})
}

func WithGsmUseJsonCreds(jsonCreds []byte) GsmOption {
	return gsmOptionFunc(func(m *googleSecret) {
		m.jsonCred = jsonCreds
	})
}
