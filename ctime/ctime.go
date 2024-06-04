/**
 * @Author: steven
 * @Description:
 * @File: time
 * @Version: 1.0.0
 * @Date: 22/08/23 12.09
 */

package ctime

import (
	"github.com/evorts/kevlars/common"
	"sync"
	"time"
)

var (
	loc     *time.Location
	locOnce sync.Once
)

type Manager interface {
	Renew() Manager
	Now() time.Time
	NowPtr() *time.Time
	Get() time.Time
	GetPtr() *time.Time
	GetLocation() *time.Location
}

type manager struct {
	now time.Time
	tz  TimeZone
	loc *time.Location
}

func (m *manager) Now() time.Time {
	m.now = time.Now().In(m.loc)
	return m.now
}

func (m *manager) NowPtr() *time.Time {
	m.Now()
	return &m.now
}

func (m *manager) Renew() Manager {
	m.now = time.Now().In(m.loc)
	return m
}

func (m *manager) Get() time.Time {
	return m.now
}

func (m *manager) GetPtr() *time.Time {
	return &m.now
}

func (m *manager) GetLocation() *time.Location {
	return m.loc
}

func New(opts ...common.Option[manager]) Manager {
	m := &manager{tz: DefaultTimeZone}
	for _, opt := range opts {
		opt.Apply(m)
	}
	var err error
	locOnce.Do(func() {
		loc, err = time.LoadLocation(m.tz.String())
	})
	if err != nil {
		panic(err)
	}
	m.loc = loc
	m.now = time.Now().In(loc)
	return m
}
