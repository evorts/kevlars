/**
 * @Author: steven
 * @Description:
 * @File: types
 * @Date: 07/06/24 22.17
 */

package common

import "fmt"

type Bytes []byte

func (b Bytes) String() string {
	return string(b)
}

func (b Bytes) ReadableString() string {
	return fmt.Sprintf("%x", b)
}

func (b Bytes) Raw() []byte {
	return b
}
