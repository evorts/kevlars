/**
 * @Author: steven
 * @Description:
 * @File: generic
 * @Date: 27/03/24 11.04
 */

package common

type IntegerAll interface {
	Integer | IntegerUnsigned
}

type Integer interface {
	int | int8 | int16 | int32 | int64
}

type IntegerUnsigned interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}
