/**
 * @Author: steven
 * @Description:
 * @File: manager
 * @Date: 30/04/24 18.13
 */

package inmemory

import (
	"context"
	"time"
)

type Manager interface {
	MustConnect(ctx context.Context) Manager
	Connect(ctx context.Context) error

	Set(ctx context.Context, key string, value interface{}) error
	SetWhenNotExist(ctx context.Context, key string, value interface{}) error
	SetWithExpire(ctx context.Context, key string, value interface{}, expire time.Duration) error
	SetWithExpireWhenNotExist(ctx context.Context, key string, value interface{}, expire time.Duration) error
	SetString(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string, bindTo interface{}) error
	GetString(ctx context.Context, key string) string
	Del(ctx context.Context, keys ...string) error

	HSet(ctx context.Context, key string, value ...interface{}) error
	HSetWhenNotExist(ctx context.Context, key, field string, value interface{}) error
	HGet(ctx context.Context, key, field string, bindTo interface{}) error
	HDel(ctx context.Context, key string, fields ...string) error
	HGetAll(ctx context.Context, key string, bindTo interface{}) error

	HMSet(ctx context.Context, key string, values ...interface{}) error
	HMGet(ctx context.Context, key, field string, bindTo interface{}) error

	Ping() error
}
