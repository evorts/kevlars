/**
 * @Author: steven
 * @Description:
 * @File: init
 * @Date: 15/01/24 21.02
 */

package common

type Init[T any] interface {
	Init() error
	MustInit() T
}
