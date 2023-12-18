package scheduler

import (
	"errors"
	"github.com/evorts/kevlars/logger"
	"github.com/go-co-op/gocron"
	"time"
)

type Manager interface {
	MustInit() Manager
	Init() error
	WithTimeZone(zone *time.Location) Manager
	WithTasks(tasks ...func()) Manager
	WithSupportSeconds() Manager
	StartAsync()
	StartBlocking()
}

type manager struct {
	zone            *time.Location
	cron            string // cron expression
	activateSeconds bool   // support schedule in seconds timeframe
	schedulers      *gocron.Scheduler
	tasks           []func()
	name            string
	log             logger.Manager
}

var (
	defaultTimeLocation = time.FixedZone("UTC+7", 7*3600)
)

func (m *manager) WithTimeZone(zone *time.Location) Manager {
	m.zone = zone
	return m
}

func (m *manager) WithTasks(tasks ...func()) Manager {
	m.tasks = tasks
	return m
}

func (m *manager) WithSupportSeconds() Manager {
	m.activateSeconds = true
	return m
}
func (m *manager) MustInit() Manager {
	if err := m.Init(); err != nil {
		panic(err)
	}
	return m
}

func (m *manager) Init() error {
	if len(m.tasks) < 1 {
		return errors.New("no tasks to run")
	}
	m.schedulers = gocron.NewScheduler(m.zone)
	if m.activateSeconds {
		m.schedulers.CronWithSeconds(m.cron)
	} else {
		m.schedulers.Cron(m.cron)
	}
	for _, task := range m.tasks {
		_, err := m.schedulers.Do(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) instantiate() {
	if m.schedulers == nil {
		m.MustInit()
	}
}

func (m *manager) StartAsync() {
	m.instantiate()
	m.log.WarnWithProps(map[string]interface{}{
		"name": m.name,
		"cron": m.cron,
	}, "schedule started as async")
	m.schedulers.StartAsync()
}

func (m *manager) StartBlocking() {
	m.instantiate()
	m.log.WarnWithProps(map[string]interface{}{
		"name": m.name,
		"cron": m.cron,
	}, "schedule started as blocking")
	m.schedulers.StartBlocking()
}

func New(cron string, opts ...Option) Manager {
	m := &manager{cron: cron, zone: defaultTimeLocation, log: logger.NewNoop()}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}
