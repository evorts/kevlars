/**
 * @Author: steven
 * @Description:
 * @File: array
 * @Date: 29/09/23 10.46
 */

package utils

import "fmt"

func InArray[T string | int | int64 | uint | uint64 | float32 | float64](arr []T, v T) bool {
	for _, av := range arr {
		if av == v {
			return true
		}
	}
	return false
}

func ToArrayOfInterface[T string | int | int64 | float32 | float64](arrays []T) []interface{} {
	rs := make([]interface{}, 0)
	for _, item := range arrays {
		rs = append(rs, item)
	}
	return rs
}

func ArrayToMapInt8(arrOfString []string) map[string]int8 {
	rs := make(map[string]int8, 0)
	if len(arrOfString) < 1 {
		return rs
	}
	for _, str := range arrOfString {
		if _, ok := rs[str]; !ok {
			rs[str] = 1
			continue
		}
		rs[str]++
	}
	return rs
}

func ArrayInterfaceToString(arr []interface{}) []string {
	rs := make([]string, 0)
	if arr == nil {
		return rs
	}
	for _, item := range arr {
		rs = append(rs, fmt.Sprintf("%v", item))
	}
	return rs
}

func GetItemFromMapArray(mapArray []map[string]interface{}, filter func(item map[string]interface{}) bool) map[string]interface{} {
	for _, item := range mapArray {
		if filter(item) {
			return item
		}
	}
	return nil
}

func FindByInMapArray[T string | int | bool](field string, value T, collection []map[string]interface{}) map[string]interface{} {
	for _, item := range collection {
		if _, ok := item[field]; !ok {
			continue
		}
		if item[field] == value {
			return item
		}
	}
	return nil
}

func SliceItems[T any](items []T, page, limit int) []T {
	page = Iif(page < 1, 1, page)
	limit = Iif(limit < 1, 10, limit)
	offset := (page - 1) * limit
	if offset < len(items) {
		items = items[offset:]
	}
	if len(items) > limit {
		items = items[:limit]
	}
	return items
}
