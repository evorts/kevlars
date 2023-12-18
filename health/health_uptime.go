package health

import (
	"github.com/shirou/gopsutil/v3/host"
	"net/http"
	"time"
)

type processUptime struct {
	start time.Time
}

func (u *processUptime) HealthChecks() map[string][]*Checks {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return map[string][]*Checks{
		"uptime": {
			{
				ComponentType: "process",
				ObservedValue: time.Now().UTC().Sub(u.start).Seconds(),
				ObservedUnit:  "s",
				Status:        Pass,
				Time:          now,
			},
		},
	}
}

func (*processUptime) AuthorizeHealth(*http.Request) bool {
	return true
}

// ProcessUptime returns a ChecksProvider for health checks about the process uptime.
// Note that it does not really return the process uptime, but the time since calling this function.
func ProcessUptime() ChecksProvider {
	return &processUptime{start: time.Now().UTC()}
}

type systemUptime struct {
}

func (u *systemUptime) HealthChecks() map[string][]*Checks {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	var (
		err  error
		info *host.InfoStat
	)
	var uptime = &Checks{
		ComponentType: "system",
		Status:        Fail,
		Time:          now,
	}
	if info, err = host.Info(); err != nil {
		uptime.Output = err.Error()
	} else {
		uptime.Status = Pass
		uptime.ObservedValue = info.Uptime
		uptime.ObservedUnit = "s"
	}
	return map[string][]*Checks{
		"uptime": {
			uptime,
		},
	}
}

func (*systemUptime) AuthorizeHealth(*http.Request) bool {
	return true
}

// SystemUptime returns a ChecksProvider for health checks about the system uptime.
func SystemUptime() ChecksProvider {
	return &systemUptime{}
}
