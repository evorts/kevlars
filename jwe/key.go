/**
 * @Author: steven
 * @Description:
 * @File: key
 * @Date: 01/06/24 07.51
 */

package jwe

import (
	"crypto/rsa"
	"github.com/evorts/kevlars/crypt"
)

type (
	PrivateKey string
	KeyStorage interface {
		GetPrivate() *rsa.PrivateKey
	}
)

func (k PrivateKey) String() string {
	return string(k)
}
func (k PrivateKey) GetKey() (*rsa.PrivateKey, error) {
	return crypt.GenerateRsaPrivateKeyFromPemString(k.String())
}
