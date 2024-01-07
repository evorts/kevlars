/**
 * @Author: steven
 * @Description:
 * @File: auth
 * @Date: 24/12/23 21.50
 */

package auth

import "context"

type Manager interface {
	GetRequestingPartyPermissions(ctx context.Context, token, realm string, opts RequestingPartyTokenOptions) (RequestingPartyPermissions, error)
}

type manager struct {
}

func (m *manager) GetRequestingPartyPermissions(ctx context.Context, token, realm string, opts RequestingPartyTokenOptions) (RequestingPartyPermissions, error) {
	//TODO implement me
	panic("implement me")
}

func New() Manager {
	return &manager{}
}
