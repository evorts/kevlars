/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 07/06/24 22.44
 */

package crypt

import (
	"github.com/evorts/kevlars/common"
	"hash"
)

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
