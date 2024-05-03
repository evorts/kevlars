/**
 * @Author: steven
 * @Description:
 * @File: convert
 * @Version: 1.0.0
 * @Date: 14/09/23 15.32
 */

package utils

import (
	"encoding/json"
	"strconv"
	"unsafe"
)

func StringToInt(v string) int {
	return StringToIntDV(v, 0)

}

func StringToIntDV(v string, dv int) int {
	if rs, err := strconv.Atoi(v); err == nil {
		return rs
	}
	return dv
}

func IntToString(v int) string {
	return strconv.Itoa(v)
}

func ToStruct[T any](src interface{}, dst *T) *T {
	b, err := json.Marshal(src)
	if err != nil {
		return nil
	}
	if err = json.Unmarshal(b, dst); err != nil {
		return nil
	}
	return dst
}

func ToPtr[T any](v T) *T {
	return &v
}

func BytesToString(b []byte) string {
	return string(b)
}

func StringToBytes(s string) []byte {
	return []byte(s)
}

// BytesToStringUnsafe converts byte slice to string.
func BytesToStringUnsafe(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytesUnsafe converts string to byte slice.
func StringToBytesUnsafe(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func Atoi(b []byte) (int, error) {
	return strconv.Atoi(BytesToString(b))
}

func ParseInt(b []byte, base int, bitSize int) (int64, error) {
	return strconv.ParseInt(BytesToString(b), base, bitSize)
}

func ParseUint(b []byte, base int, bitSize int) (uint64, error) {
	return strconv.ParseUint(BytesToString(b), base, bitSize)
}

func ParseFloat(b []byte, bitSize int) (float64, error) {
	return strconv.ParseFloat(BytesToString(b), bitSize)
}
