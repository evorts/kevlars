/**
 * @Author: steven
 * @Description:
 * @File: vars
 * @Date: 18/12/23 11.12
 */

package config

type RemoteProvider string

const (
	RemoteProviderGSM    RemoteProvider = "google_secret_manager"
	RemoteProviderConsul RemoteProvider = "consul"
	RemoteProviderNone   RemoteProvider = "none"
)

func (rp RemoteProvider) String() string {
	return string(rp)
}

type Type string

const (
	TypeYaml Type = "yaml"
	TypeJson Type = "json"
)

func (t Type) String() string {
	return string(t)
}
