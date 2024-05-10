/**
 * @Author: steven
 * @Description:
 * @File: eval
 * @Date: 10/05/24 22.09
 */

package eval

import "reflect"

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

func IsEmpty[T any](v T) bool {
	t := reflect.TypeOf(v)
	k := t.Kind()
	switch true {
	case k == reflect.Slice || k == reflect.Array || k == reflect.Map:
		if v == nil || reflect.ValueOf(v).Len() == 0 {
			return true
		}
	case k == reflect.Ptr:
		if v == nil {
			return true
		}
	case k == reflect.Func:
		if v == nil {
			return true
		}
	case k == reflect.String:
		if reflect.ValueOf(v).String() == "" {
			return true
		}
	case k == reflect.Int, k == reflect.Int8, k == reflect.Int16, k == reflect.Int32, k == reflect.Int64:
		if reflect.ValueOf(v).Int() == 0 {
			return true
		}
	case k == reflect.Uint, k == reflect.Uint8, k == reflect.Uint16, k == reflect.Uint32, k == reflect.Uint64:
		if reflect.ValueOf(v).Uint() == 0 {
			return true
		}
	case k == reflect.Float32, k == reflect.Float64:
		if reflect.ValueOf(v).Float() == 0 {
			return true
		}
	}
	return false
}

func IsNil[T any](v T) bool {
	return v == nil
}
