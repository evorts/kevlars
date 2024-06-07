/**
 * @Author: steven
 * @Description:
 * @File: decrypt
 * @Date: 07/06/24 22.00
 */

package crypt

import "github.com/evorts/kevlars/common"

type Decrypter interface {
	Decrypt(value common.Bytes) (common.Bytes, error)
}
