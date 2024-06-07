/**
 * @Author: steven
 * @Description:
 * @File: helpers
 * @Date: 07/06/24 22.59
 */

package crypt

import "bytes"

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	bufLen := len(ciphertext)
	padLen := blockSize - bufLen%blockSize
	padded := make([]byte, bufLen+padLen)
	copy(padded, ciphertext)
	for i := 0; i < padLen; i++ {
		padded[bufLen+i] = byte(padLen)
	}
	return padded
}

func pkcs7UnPadding(ciphertext []byte) []byte {
	padding := len(ciphertext) - int(ciphertext[len(ciphertext)-1])
	buf := make([]byte, padding)
	copy(buf, ciphertext[:padding])
	return buf
}
