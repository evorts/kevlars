/**
 * @Author: steven
 * @Description:
 * @File: zone
 * @Version: 1.0.0
 * @Date: 22/08/23 12.09
 */

package ts

type TimeZone string

const (
	TZAsiaJakarta TimeZone = "Asia/Jakarta"
)

func (t TimeZone) String() string {
	return string(t)
}

var (
	DefaultTimeZone = TZAsiaJakarta
)
