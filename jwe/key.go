/**
 * @Author: steven
 * @Description:
 * @File: key
 * @Date: 01/06/24 07.51
 */

package jwe

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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
	return generateRsaPrivateKeyFromPemString(k.String())
}

func generateRsaPrivateKeyFromPemString(privatePem string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privatePem))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pri, nil
}
