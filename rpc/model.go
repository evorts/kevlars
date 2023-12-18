/**
 * @Author: steven
 * @Description:
 * @File: model
 * @Version: 1.0.0
 * @Date: 17/08/23 12.13
 */

package rpc

import "context"

type AuthorizationCredential string

func (c AuthorizationCredential) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": string(c),
	}, nil
}
func (c AuthorizationCredential) RequireTransportSecurity() bool {
	return false
}
