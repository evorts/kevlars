/**
 * @Author: steven
 * @Description:
 * @File: builder
 * @Date: 25/05/24 23.30
 */

package utils

func RepeatInSlice[T any](v T, count int) []T {
	rs := make([]T, count)
	for i := 0; i < count; i++ {
		rs[i] = v
	}
	return rs
}
