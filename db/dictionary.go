/**
 * @Author: steven
 * @Description:
 * @File: dictionary
 * @Date: 18/05/24 15.26
 */

package db

var (
	valueByDriver = map[SupportedDriver]map[ValueKey]string{
		DriverPostgreSQL: {
			ValueKeyNow: "current_timestamp",
		},
		DriverMySQL: {
			ValueKeyNow: "now()",
		},
	}
)

type ValueKey string

const (
	ValueKeyNow ValueKey = "now"
)

func (k ValueKey) String(driver SupportedDriver) string {
	v, ok := valueByDriver[driver]
	if !ok {
		return ""
	}
	if vv, ok2 := v[k]; ok2 {
		return vv
	}
	return ""
}
