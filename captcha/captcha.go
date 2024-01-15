/**
 * @Author: steven
 * @Description:
 * @File: captcha
 * @Date: 15/01/24 21.16
 */

package captcha

import "github.com/evorts/kevlars/common"

type Manager[T any] interface {
	Generate() (id, answer string, result T)
	Verify(id, value string, clear bool) bool

	common.Init[Manager[T]]
}
