/**
 * @Author: steven
 * @Description:
 * @File: encrypt
 * @Date: 07/06/24 22.00
 */

package crypt

import (
	"github.com/evorts/kevlars/common"
	"golang.org/x/crypto/sha3"
	"hash"
)

type Encrypter interface {
	Encrypt(value common.Bytes) (common.Bytes, error)
}

type Hasher interface {
	SHA3(use SHA3, bytes common.Bytes) (common.Bytes, error)
}

type hasher struct{}

func (h *hasher) SHA3(use SHA3, bytes common.Bytes) (common.Bytes, error) {
	var hashing hash.Hash
	switch use {
	case Shake128:
		hashing = sha3.NewShake128()
	case Shake256:
		hashing = sha3.NewShake256()
	case Bit224:
		hashing = sha3.New224()
	case Bit256:
		hashing = sha3.New256()
	case Bit384:
		hashing = sha3.New384()
	case Bit512:
		hashing = sha3.New512()
	default:
		return nil, ErrInvalidSHA3Instance
	}
	hashing.Write(bytes)
	return hashing.Sum(nil), nil
}

func NewHasher() Hasher {
	return &hasher{}
}
