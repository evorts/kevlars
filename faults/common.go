/**
 * @Author: steven
 * @Description:
 * @File: common
 * @Version: 1.0.0
 * @Date: 25/09/23 00.39
 */

package faults

type IError interface {
	error
	Code() string
	Props() map[string]string
	UnWrap() error
}

type commonError struct {
	code          string
	message       string
	details       map[string]string
	underlyingErr error
}

type IOption[T any] interface {
	apply(o *T)
}

type option[T any] func(o *T) T

func (o option[T]) apply(opt *T) {
	o(opt)
}
