/**
 * @Author: steven
 * @Description:
 * @File: numbers
 * @Date: 27/03/24 11.01
 */

package utils

import "github.com/evorts/kevlars/common"

func NumberInRange[T common.IntegerAll | common.Float](v T, min T, max T) bool {
	return v >= min && v <= max
}

func NumberInRangeEx[T common.IntegerAll | common.Float](v T, min T, max T) bool {
	return v > min && v < max
}
