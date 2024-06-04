/**
 * @Author: steven
 * @Description:
 * @File: configNoop
 * @Date: 05/06/24 00.55
 */

package config

import (
	"time"
)

type noop struct{}

func (n *noop) Get(key string) interface{} {
	return nil
}

func (n *noop) GetBool(key string) bool {
	return false
}

func (n *noop) GetBoolOrElse(key string, elseValue bool) bool {
	return elseValue
}

func (n *noop) GetFloat64(key string) float64 {
	return 0
}

func (n *noop) GetFloat64OrElse(key string, elseValue float64) float64 {
	return elseValue
}

func (n *noop) GetInt(key string) int {
	return 0
}

func (n *noop) GetIntOrElse(key string, elseValue int) int {
	return elseValue
}

func (n *noop) GetIntSlice(key string) []int {
	return []int{}
}

func (n *noop) GetIntSliceOrElse(key string, elseValue []int) []int {
	return elseValue
}

func (n *noop) GetString(key string) string {
	return ""
}

func (n *noop) GetStringOrElse(key string, elseValue string) string {
	return elseValue
}

func (n *noop) GetStringMap(key string) map[string]interface{} {
	return map[string]interface{}{}
}

func (n *noop) GetStringMapOrElse(key string, elseValue map[string]interface{}) map[string]interface{} {
	return elseValue
}

func (n *noop) GetStringMapString(key string) map[string]string {
	return map[string]string{}
}

func (n *noop) GetStringMapStringOrElse(key string, elseValue map[string]string) map[string]string {
	return elseValue
}

func (n *noop) GetStringSlice(key string) []string {
	return []string{}
}

func (n *noop) GetStringSliceOrElse(key string, orElse []string) []string {
	return orElse
}

func (n *noop) GetMapArray(key string) []map[string]interface{} {
	return make([]map[string]interface{}, 0)
}

func (n *noop) GetTime(key string) time.Time {
	return time.Time{}
}

func (n *noop) GetDuration(key string) time.Duration {
	return time.Duration(0)
}

func (n *noop) GetDurationOrElse(key string, elseValue time.Duration) time.Duration {
	return elseValue
}

func (n *noop) UnmarshalTo(key string, to interface{}) error {
	return nil
}

func (n *noop) IsSet(key string) bool {
	return false
}

func (n *noop) AllSettings() map[string]interface{} {
	return map[string]interface{}{}
}

func (n *noop) Init() error {
	return nil
}

func (n *noop) MustInit() Manager {
	return n
}

func NewNoop() Manager {
	return &noop{}
}
