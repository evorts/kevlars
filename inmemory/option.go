/**
 * @Author: steven
 * @Description:
 * @File: Option
 * @Version: 1.0.0
 * @Date: 20/06/23 11.44
 */

package inmemory

type Option[T any] interface {
	apply(*T)
}

type option[T any] func(*T)

func (o option[T]) apply(m *T) {
	o(m)
}
