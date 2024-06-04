/**
 * @Author: steven
 * @Description:
 * @File: jwe
 * @Date: 01/06/24 07.42
 */

package jwe

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/evorts/kevlars/common"
	"github.com/go-jose/go-jose/v4"
	"time"
)

type Manager interface {
	Encode(v Claim) (token string, err error)
	Decode(token string) (v Claim, err error)

	Init() error
	MustInit() Manager
}

type jwe struct {
	key               *rsa.PrivateKey
	keyAlgorithm      jose.KeyAlgorithm
	expire            time.Duration //expiration
	signerAlgorithm   jose.SignatureAlgorithm
	signer            jose.Signer
	contentEncryption jose.ContentEncryption
	encrypter         jose.Encrypter
}

func (j *jwe) Init() error {
	var err error
	j.signer, err = jose.NewSigner(jose.SigningKey{Algorithm: j.signerAlgorithm, Key: j.key}, nil)
	if err != nil {
		return err
	}
	j.encrypter, err = jose.NewEncrypter(
		j.contentEncryption,
		jose.Recipient{
			Algorithm: j.keyAlgorithm,
			Key:       j.key.PublicKey,
		}, nil)
	return err
}

func (j *jwe) MustInit() Manager {
	if err := j.Init(); err != nil {
		panic(err)
	}
	return j
}

func NewJWE(key *rsa.PrivateKey, opts ...common.Option[jwe]) Manager {
	m := &jwe{
		key:               key,
		signerAlgorithm:   jose.PS512,
		keyAlgorithm:      jose.RSA_OAEP,
		contentEncryption: jose.A128GCM,
	}
	for _, opt := range opts {
		opt.Apply(m)
	}

	return m
}

func (j *jwe) Encode(v Claim) (token string, err error) {
	// encode to string first
	vb, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	// encrypt it
	enc, errEnc := j.encrypter.Encrypt(vb)
	if errEnc != nil {
		return "", errEnc
	}
	return enc.CompactSerialize()
}

func (j *jwe) Decode(token string) (v Claim, err error) {
	jws, err := jose.ParseEncrypted(token, []jose.KeyAlgorithm{j.keyAlgorithm}, []jose.ContentEncryption{j.contentEncryption})
	if err != nil {
		return Claim{}, err
	}
	rs, errD := jws.Decrypt(j.key)
	if errD != nil {
		return Claim{}, errD
	}
	err = json.Unmarshal(rs, &v)
	return
}
