/**
 * @Author: steven
 * @Description:
 * @File: map
 * @Date: 17/12/23 23.46
 */

package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GetValueOnMap(dict map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if v, ok := dict[key]; ok {
		return v
	}
	return defaultValue
}

func ValueOnMapByKey[T any](dict map[string]T, key string, defaultValue T) T {
	if v, ok := dict[key]; ok {
		return v
	}
	return defaultValue
}
func ValueOnMap[T string | int | int64 | float64 | interface{}](dict map[string]T, key string, defaultValue T) T {
	if v, ok := dict[key]; ok {
		return v
	}
	return defaultValue
}

func MapToStringArray(value map[string]interface{}, format, separator string) string {
	if value == nil {
		return ""
	}
	var rs []string
	if len(format) < 1 {
		format = "%s=%s"
	}
	if len(separator) < 1 {
		separator = " "
	}
	for k, v := range value {
		rs = append(rs, fmt.Sprintf(format, k, v))
	}
	return strings.Join(rs, separator)
}

func MapInterfaceToMapString(src map[string]interface{}) map[string]string {
	if src == nil {
		return nil
	}
	rs := make(map[string]string, 0)
	for k, v := range src {
		if vs, ok := v.(string); ok {
			rs[k] = vs
		}
	}
	return rs
}

func ToMapInterface(src interface{}) map[string]interface{} {
	rs := map[string]interface{}{}
	if src == nil {
		return rs
	}
	jsonData, err := json.Marshal(src)
	if err != nil {
		return rs
	}
	_ = json.Unmarshal(jsonData, &rs)
	return rs
}

func MapMerge[T string | int8 | int16 | int32 | int | interface{}](src map[string]T, dst map[string]T) map[string]T {
	if src == nil {
		src = make(map[string]T)
	}
	if dst == nil {
		dst = make(map[string]T)
	}
	for k, v := range dst {
		src[k] = v
	}
	return src
}
