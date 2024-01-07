/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 25/12/23 08.22
 */

package queue

type Option[T any] interface {
	apply(t *T)
}

type optionFunc[T any] func(t *T)

func (o optionFunc[T]) apply(t *T) {
	o(t)
}

type PublishOption[T any] interface {
	apply(t *T)
}

type publishOptionFunc[T any] func(t *T)

func (o publishOptionFunc[T]) apply(t *T) {
	o(t)
}

type ActOption[T any] interface {
	apply(t *T)
}

type actOptionFunc[T any] func(t *T)

func (o actOptionFunc[T]) apply(t *T) {
	o(t)
}
