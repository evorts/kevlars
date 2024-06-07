/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 07/06/24 22.44
 */

package crypt

import (
	"encoding/base64"
	"github.com/evorts/kevlars/common"
	"hash"
)

func WithCipher(v Cipher) common.Option[manager] {
	return common.OptionFunc[manager](func(m *manager) {
		m.cipher = v
	})
}

func WithKey(key common.Bytes) common.Option[manager] {
	return common.OptionFunc[manager](func(m *manager) {
		m.key = key
	})
}

func WithIV(iv common.Bytes) common.Option[manager] {
	return common.OptionFunc[manager](func(m *manager) {
		m.iv = iv
	})
}

func WithHash(hash hash.Hash) common.Option[manager] {
	return common.OptionFunc[manager](func(m *manager) {
		m.hash = hash
	})
}

func WithPrivateKey(b64PEMPrivateKey string) common.Option[manager] {
	return common.OptionFunc[manager](func(m *manager) {
		var dec64 common.Bytes
		_, err := base64.StdEncoding.Decode(dec64, []byte(b64PEMPrivateKey))
		if err != nil {
			panic(err)
		}
		m.pk, err = GenerateRsaPrivateKeyFromPemString(dec64.String())
		if err != nil {
			panic(err)
		}
	})
}
