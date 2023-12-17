/**
 * @Author: steven
 * @Description:
 * @File: helpers
 * @Version: 1.0.0
 * @Date: 29/08/23 23.23
 */

package ts

import (
	"sync"
	"time"
)

var (
	defaultZoneLocationOnce sync.Once
)

var defaultLoc *time.Location

func Now() time.Time {
	defaultZoneLocationOnce.Do(func() {
		defaultLoc, _ = time.LoadLocation(DefaultTimeZone.String())
	})
	return time.Now().In(defaultLoc)
}

func NowPtr() *time.Time {
	now := Now()
	return &now
}

func NowWithTZ(zone TimeZone) time.Time {
	location, _ := time.LoadLocation(zone.String())
	return time.Now().In(location)
}

func NowWithTZPtr(zone TimeZone) *time.Time {
	nowTZ := NowWithTZ(zone)
	return &nowTZ
}
