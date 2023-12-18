/**
 * @Author: steven
 * @Description:
 * @File: convert
 * @Version: 1.0.0
 * @Date: 14/09/23 15.32
 */

package utils

import "strconv"

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
