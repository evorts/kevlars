/**
 * @Author: steven
 * @Description:
 * @File: ifs
 * @Date: 17/12/23 22.12
 */

package utils

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func Iif[T any](expr bool, whenTrue, whenFalse T) T {
	if expr {
		return whenTrue
	}
	return whenFalse
}

func IfNil[T string | int | *time.Time | map[string]interface{}](value interface{}, whenNil T) T {
	if value == nil {
		return whenNil
	}
	// check if its map
	mapValue, isMap := value.(map[string]interface{})
	if isMap && mapValue == nil {
		return whenNil
	}
	if reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil() {
		return whenNil
	}
	if v, ok := value.(T); ok {
		return v
	}
	return whenNil
}

func IfEmpty[T string | int](value T, whenEmpty T) T {
	t := reflect.TypeOf(value)
	switch t.Name() {
	case "string":
		return Iif(fmt.Sprintf("%v", value) == "", whenEmpty, value) //nolint:govet
	case "int":
		v := reflect.ValueOf(value).Int()
		return Iif(v == 0, whenEmpty, value)
	}
	return value
}

func IfNotEmpty[T string | int](value T, whenNotEmpty T) T {
	t := reflect.TypeOf(value)
	switch t.Name() {
	case "string":
		return Iif(fmt.Sprintf("%v", value) != "", whenNotEmpty, value) //nolint:govet
	case "int":
		v := reflect.ValueOf(value).Int()
		return Iif(v != 0, whenNotEmpty, value)
	}
	return value
}

func IfErrorThen(err error, run func()) {
	if err != nil {
		run()
	}
}

func IfTrueThen(expr bool, run func()) {
	if expr {
		run()
	}
}

func IfR[T *any](expr bool, run func() T) T {
	if expr {
		return run()
	}
	return nil
}

func IfRE[T *any](expr bool, run func() (T, error)) (T, error) {
	if expr {
		return run()
	}
	return nil, errors.New("expression are returning false")
}

func IfE(expr bool, then func(), el func()) {
	if expr {
		then()
		return
	}
	el()
}

func IfER[T any](expr bool, then func() T, el func() T) T {
	if expr {
		return then()
	}
	return el()
}

func OR(ors ...bool) bool {
	for _, or := range ors {
		if or {
			return true
		}
	}
	return false
}

func AND(ands ...bool) bool {
	for _, and := range ands {
		if !and {
			return false
		}
	}
	return true
}
