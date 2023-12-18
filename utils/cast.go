package utils

import (
	"errors"
	"strconv"
)

func CastToMapStringND(v interface{}) map[string]string {
	return CastToMapString(v, make(map[string]string))
}

func CastToMapString(v interface{}, defaultValue map[string]string) map[string]string {
	if vv, okv := v.(map[string]string); okv {
		return vv
	}
	return defaultValue
}

func CastToMapInterface(v interface{}, defaultValue map[string]interface{}) map[string]interface{} {
	if vv, okv := v.(map[string]interface{}); okv {
		return vv
	}
	return defaultValue
}

func CastToStringND(v interface{}) string {
	return CastToString(v, "")
}

func CastToString(v interface{}, defaultValue string) string {
	if vv, okv := v.(string); okv {
		return vv
	}
	return defaultValue
}

func CastToIntND(v interface{}) int {
	return CastToInt(v, 0)
}

func CastToInt(v interface{}, defaultValue int) int {
	if vv, okv := v.(int); okv {
		return vv
	}
	return defaultValue
}

func CastToInt64ND(v interface{}) int64 {
	return CastToInt64(v, 0)
}

func CastToInt64(v interface{}, defaultValue int64) int64 {
	if vv, okv := v.(int64); okv {
		return vv
	}
	return defaultValue
}

func CastToBoolND(v interface{}) bool {
	return CastToBool(v, false)
}

func CastToBool(v interface{}, defaultValue bool) bool {
	if vv, okv := v.(bool); okv {
		return vv
	}
	return defaultValue
}

func CastToFloat64ND(v interface{}) float64 {
	return CastToFloat64(v, 0.0)
}

func CastToFloat64(v interface{}, defaultValue float64) float64 {
	if vv, okv := v.(float64); okv {
		return vv
	}
	return defaultValue
}

func Float64ToString(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func CastToStruct[T any](v interface{}, defaultValue T) T {
	if vv, ok := v.(T); ok {
		return vv
	}
	return defaultValue
}

func CastToStructND[T any](v interface{}, result func(T)) error {
	if vv, ok := v.(T); ok {
		result(vv)
		return nil
	}
	return errors.New("failed to cast into struct")
}
