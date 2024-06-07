/**
 * @Author: steven
 * @Description:
 * @File: cipher_rsa
 * @Date: 07/06/24 23.51
 */

package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/evorts/kevlars/common"
	"hash"
)

type cipherRsa struct {
	pubKey     *rsa.PublicKey
	privateKey *rsa.PrivateKey
	hash       hash.Hash
}

func (c *cipherRsa) Encrypt(value common.Bytes) (common.Bytes, error) {
	return rsa.EncryptOAEP(c.hash, rand.Reader, c.pubKey, value, nil)
}

func (c *cipherRsa) Decrypt(value common.Bytes) (common.Bytes, error) {
	return rsa.DecryptOAEP(c.hash, rand.Reader, c.privateKey, value, nil)
}

func newCipherRsa(privateKey *rsa.PrivateKey, hash hash.Hash) (*cipherRsa, error) {
	return &cipherRsa{
		privateKey: privateKey,
		pubKey:     &privateKey.PublicKey,
		hash:       hash,
	}, nil
}
