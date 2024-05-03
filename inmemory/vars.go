/**
 * @Author: steven
 * @Description:
 * @File: vars
 * @Date: 03/05/24 08.21
 */

package inmemory

type Provider string

const (
	ProviderValKey Provider = "valkey"
	ProviderRedis  Provider = "redis"
)

func (p Provider) String() string { return string(p) }
