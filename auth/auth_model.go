/**
 * @Author: steven
 * @Description:
 * @File: auth_model
 * @Date: 14/05/24 11.51
 */

package auth

type OAuthPermission struct {
	ResourceName string   `json:"resource_name"`
	Scopes       []string `json:"scopes"`
}

type OAuthPermissions []*OAuthPermission
