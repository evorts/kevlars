/**
 * @Author: steven
 * @Description:
 * @File: vars
 * @Date: 07/06/24 22.09
 */

package crypt

import "errors"

type SHA3 string

const (
	Shake128 SHA3 = "sha3-shake-128"
	Shake256 SHA3 = "sha3-shake-256"
	Bit224   SHA3 = "sha3-bits-224"
	Bit256   SHA3 = "sha3-bits-256"
	Bit384   SHA3 = "sha3-bits-384"
	Bit512   SHA3 = "sha3-bits-512"
)

type Cipher string

const (
	Cipher3DES   Cipher = "3DES"
	CipherAESCBC Cipher = "AES-CBC"
	CipherAESGCM Cipher = "AES-GCM"
	CipherRSA    Cipher = "RSA"
)

type cAes string

const (
	aesCBC cAes = cAes(CipherAESCBC)
	aesGCM cAes = cAes(CipherAESGCM)
)

func (c Cipher) String() string { return string(c) }

var (
	ErrNoKeyDefined         = errors.New("no key defined")
	ErrInvalidKeyDefined    = errors.New("invalid key defined")
	ErrNoIVDefined          = errors.New("no iv defined")
	ErrInvalidIVDefined     = errors.New("invalid iv defined")
	ErrNoCipherDefined      = errors.New("no cipher defined")
	ErrNoEncrypterDefined   = errors.New("no encrypter defined")
	ErrNoDecrypterDefined   = errors.New("no decrypter defined")
	ErrInvalidSHA3Instance  = errors.New("invalid SHA3 instance")
	ErrPrivateKeyNotDefined = errors.New("private key not defined")
	ErrHasherNotDefined     = errors.New("hasher not defined")
	ErrCipherNotSupported   = errors.New("cipher not supported")

	ErrIVTooLong                  = errors.New("iv too long")
	ErrCipherNotMultipleBlockSize = errors.New("cipherText is not a multiple of the block size")
)
