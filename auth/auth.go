/**
 * @Author: steven
 * @Description:
 * @File: auth
 * @Date: 24/12/23 21.50
 */

package auth

import (
	"context"
	"github.com/evorts/kevlars/common"
)

type OAuthManager interface {
	GetPermissions(ctx context.Context, token, realm string) (OAuthPermissions, error)
}

type oauthManager struct {
}

func (o *oauthManager) GetPermissions(ctx context.Context, token, realm string) (OAuthPermissions, error) {
	//TODO implement me
	panic("implement me")
}

func NewOAuthManager(opts ...common.Option[oauthManager]) OAuthManager {
	m := &oauthManager{}
	for _, opt := range opts {
		opt.Apply(m)
	}
	return m
}
