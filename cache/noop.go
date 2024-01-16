/**
 * @Author: steven
 * @Description:
 * @File: cache_noop
 * @Version: 1.0.0
 * @Date: 12/09/23 11.12
 */

package cache

import (
	"context"
	"time"
)

type managerNoop struct {
}

func (m *managerNoop) MustConnect(ctx context.Context) Manager {
	return m
}

func (m *managerNoop) Connect(ctx context.Context) error {
	return nil
}

func (m *managerNoop) Set(ctx context.Context, key string, value interface{}) error {
	return nil
}

func (m *managerNoop) SetWhenNotExist(ctx context.Context, key string, value interface{}) error {
	return nil
}

func (m *managerNoop) SetWithExpire(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return nil
}

func (m *managerNoop) SetWithExpireWhenNotExist(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return nil
}

func (m *managerNoop) SetString(ctx context.Context, key, value string, expire time.Duration) error {
	return nil
}

func (m *managerNoop) Get(ctx context.Context, key string, bindTo interface{}) error {
	return nil
}

func (m *managerNoop) GetString(ctx context.Context, key string) string {
	return ""
}

func (m *managerNoop) HSet(ctx context.Context, key string, value ...interface{}) error {
	return nil
}

func (m *managerNoop) HSetWhenNotExist(ctx context.Context, key, field string, value interface{}) error {
	return nil
}

func (m *managerNoop) HGet(ctx context.Context, key, field string, bindTo interface{}) error {
	return nil
}

func (m *managerNoop) HDel(ctx context.Context, key string, fields ...string) error {
	return nil
}

func (m *managerNoop) HGetAll(ctx context.Context, key string, bindTo interface{}) error {
	return nil
}

func (m *managerNoop) HMSet(ctx context.Context, key string, values ...interface{}) error {
	return nil
}

func (m *managerNoop) HMGet(ctx context.Context, key, field string, bindTo interface{}) error {
	return nil
}

func (m *managerNoop) Ping() error {
	return nil
}

func NewNoop() Manager {
	return &managerNoop{}
}
