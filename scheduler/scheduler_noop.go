/**
 * @Author: steven
 * @Description:
 * @File: scheduler_noop
 * @Date: 14/11/23 11.49
 */

package scheduler

import (
	"github.com/go-co-op/gocron"
	"time"
)

type managerNoop struct {
	zone            *time.Location
	cron            string // cron expression
	activateSeconds bool   // support schedule in seconds timeframe
	schedulers      *gocron.Scheduler
	tasks           []func()
}

func (m *managerNoop) MustInit() Manager {
	return m
}

func (m *managerNoop) Init() error {
	return nil
}

func (m *managerNoop) WithTimeZone(zone *time.Location) Manager {
	return m
}

func (m *managerNoop) WithTasks(tasks ...func()) Manager {
	return m
}

func (m *managerNoop) WithSupportSeconds() Manager {
	return m
}

func (m *managerNoop) StartAsync() {
	// do nothing
}

func (m *managerNoop) StartBlocking() {
	// do nothing
}

func NewNoop(cron string) Manager {
	return &managerNoop{}
}
