/**
 * @Author: steven
 * @Description:
 * @File: crypt
 * @Date: 07/06/24 21.59
 */

package crypt

import (
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha512"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/rules"
	"hash"
)

type Manager interface {
	Encrypter
	Decrypter

	common.Init[Manager]
	AddOptions(opts ...common.Option[manager]) Manager
}

type crypter interface {
	Encrypter
	Decrypter
}

type manager struct {
	crypter crypter
	hash    hash.Hash
	cipher  Cipher
	key     common.Bytes
	iv      common.Bytes
	cb      cipher.Block
	pk      *rsa.PrivateKey
}

func (m *manager) Init() error {
	if m.cipher == "" {
		return ErrNoCipherDefined
	}
	if m.key == nil {
		return ErrNoKeyDefined
	}
	var err error
	switch m.cipher {
	case Cipher3DES:
		m.crypter, err = newCipher3DES(m.key)
	case CipherAESCBC, CipherAESGCM:
		if m.iv == nil {
			return ErrNoIVDefined
		}
		m.crypter, err = newCipherAes(rules.Iif(m.cipher == CipherAESGCM, aesGCM, aesCBC), m.key, m.iv)
	case CipherRSA:
		if m.pk == nil {
			return ErrPrivateKeyNotDefined
		}
		if m.hash == nil {
			return ErrHasherNotDefined
		}
		m.crypter, err = newCipherRsa(m.pk, m.hash)
	default:
		err = ErrCipherNotSupported
	}
	if err != nil {
		return err
	}
	return nil
}

func (m *manager) MustInit() Manager {
	if err := m.Init(); err != nil {
		panic(err)
	}
	return m
}

func (m *manager) Encrypt(value common.Bytes) (common.Bytes, error) {
	return m.crypter.Encrypt(value)
}

func (m *manager) Decrypt(value common.Bytes) (common.Bytes, error) {
	return m.crypter.Decrypt(value)
}

func (m *manager) AddOptions(opts ...common.Option[manager]) Manager {
	for _, opt := range opts {
		opt.Apply(m)
	}
	return m
}

func New(opts ...common.Option[manager]) Manager {
	m := &manager{
		hash: sha512.New(),
	}
	for _, opt := range opts {
		opt.Apply(m)
	}
	return m
}
