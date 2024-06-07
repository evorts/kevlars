/**
 * @Author: steven
 * @Description:
 * @File: file
 * @Date: 28/05/24 16.51
 */

package utils

import "os"

func FileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
