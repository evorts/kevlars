/**
 * @Author: steven
 * @Description:
 * @File: cipher_aes_cbc
 * @Date: 07/06/24 23.13
 */

package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"github.com/evorts/kevlars/common"
)

type cipherAES struct {
	use                cAes
	key                common.Bytes
	iv                 common.Bytes
	block              cipher.Block
	aead               cipher.AEAD
	encrypterBlockMode cipher.BlockMode
	decrypterBlockMode cipher.BlockMode
}

func (c *cipherAES) Encrypt(value common.Bytes) (common.Bytes, error) {
	if c.use == aesGCM {
		return c.aead.Seal(nil, c.iv, value, nil), nil
	}
	value = pkcs7Padding(value, aes.BlockSize)
	if len(value)%aes.BlockSize != 0 {
		return nil, ErrCipherNotMultipleBlockSize
	}
	buf := make([]byte, len(value))
	c.encrypterBlockMode.CryptBlocks(buf, value)
	return buf, nil
}

func (c *cipherAES) Decrypt(value common.Bytes) (common.Bytes, error) {
	if c.use == aesGCM {
		return c.aead.Open(nil, c.iv, value, nil)
	}
	if len(value)%aes.BlockSize != 0 {
		return nil, ErrCipherNotMultipleBlockSize
	}
	c.decrypterBlockMode.CryptBlocks(value, value)
	unPad := pkcs7UnPadding(value)
	return unPad, nil
}

func (c *cipherAES) generateIV() error {
	if len(c.iv) > aes.BlockSize {
		return ErrIVTooLong
	}
	if len(c.iv) == aes.BlockSize {
		return nil
	}
	iv := hex.EncodeToString(c.iv)
	// and blank byte to iv
	for i := len(c.iv); i < aes.BlockSize; i++ {
		iv += "00"
	}
	rs, err := hex.DecodeString(iv)
	if err != nil {
		return err
	}
	c.iv = rs
	return nil
}

func newCipherAes(use cAes, key, iv common.Bytes) (*cipherAES, error) {
	if key == nil {
		return nil, ErrNoKeyDefined
	}
	if iv == nil {
		return nil, ErrNoIVDefined
	}
	c := new(cipherAES)
	c.key = key
	c.iv = iv
	var err error
	c.block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if use == aesGCM {
		if kl := len(key); kl != 16 && kl != 24 && kl != 32 {
			return nil, ErrInvalidKeyDefined
		}
		if ivl := len(iv); ivl != 12 {
			return nil, ErrInvalidIVDefined
		}
		c.aead, err = cipher.NewGCM(c.block)
		if err != nil {
			return nil, err
		}
	} else {
		if err = c.generateIV(); err != nil {
			return nil, err
		}
		c.encrypterBlockMode = cipher.NewCBCEncrypter(c.block, c.iv)
		c.decrypterBlockMode = cipher.NewCBCDecrypter(c.block, c.iv)
	}
	return c, nil
}
