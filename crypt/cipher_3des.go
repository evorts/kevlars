/**
 * @Author: steven
 * @Description:
 * @File: cipher_3des
 * @Date: 07/06/24 22.51
 */

package crypt

import (
	"crypto/cipher"
	"crypto/des"
	"github.com/evorts/kevlars/common"
)

type cipher3DES struct {
	key                common.Bytes
	iv                 common.Bytes
	block              cipher.Block
	encrypterBlockMode cipher.BlockMode
	decrypterBlockMode cipher.BlockMode
}

func (c *cipher3DES) Encrypt(value common.Bytes) (common.Bytes, error) {
	pad := pkcs5Padding(value, c.block.BlockSize())
	enc := make([]byte, len(pad))
	c.encrypterBlockMode.CryptBlocks(enc, pad)
	return enc, nil
}

func (c *cipher3DES) Decrypt(value common.Bytes) (common.Bytes, error) {
	dec := make([]byte, len(value))
	c.decrypterBlockMode.CryptBlocks(dec, value)
	dec = pkcs5UnPadding(dec)
	return dec, nil
}

func newCipher3DES(key common.Bytes) (*cipher3DES, error) {
	c := &cipher3DES{
		key: key,
	}
	v, err := des.NewTripleDESCipher(c.key)
	if err != nil {
		return nil, err
	}
	c.iv = key[:des.BlockSize]
	c.block = v
	c.encrypterBlockMode = cipher.NewCBCEncrypter(c.block, c.key)
	c.decrypterBlockMode = cipher.NewCBCDecrypter(c.block, c.key)
	return c, nil
}
