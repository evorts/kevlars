/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 15/01/24 21.09
 */

package common

type Option[T any] interface {
	Apply(*T)
}

type OptionFunc[T any] func(*T)

func (o OptionFunc[T]) Apply(t *T) {
	o(t)
}
