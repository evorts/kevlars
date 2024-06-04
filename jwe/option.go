/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 01/06/24 14.19
 */

package jwe

import (
	"crypto/rsa"
	"github.com/evorts/kevlars/common"
	"github.com/go-jose/go-jose/v4"
	"time"
)

func WithPrivateKey(v *rsa.PrivateKey) common.Option[jwe] {
	return common.OptionFunc[jwe](func(j *jwe) {
		j.key = v
	})
}

func WithExpiration(v time.Duration) common.Option[jwe] {
	return common.OptionFunc[jwe](func(j *jwe) {
		j.expire = v
	})
}

func WithKeyAlgorithm(v jose.KeyAlgorithm) common.Option[jwe] {
	return common.OptionFunc[jwe](func(j *jwe) {
		j.keyAlgorithm = v
	})
}

func WithSignerAlgorithm(v jose.SignatureAlgorithm) common.Option[jwe] {
	return common.OptionFunc[jwe](func(j *jwe) {
		j.signerAlgorithm = v
	})
}

func WithContentEncryption(v jose.ContentEncryption) common.Option[jwe] {
	return common.OptionFunc[jwe](func(j *jwe) {
		j.contentEncryption = v
	})
}
