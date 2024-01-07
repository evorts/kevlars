/**
 * @Author: steven
 * @Description:
 * @File: model
 * @Date: 24/12/23 21.51
 */

package auth

type RequestingPartyTokenOptions struct {
	GrantType string
	Audience  string
}

type RequestingPartyPermissions []RequestingPartyPermission

type RequestingPartyPermission struct {
	ResourceName string
	Scopes       []string
}
