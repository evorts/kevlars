package health

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"net/http"
	"os"
	"time"
)

type sysInfo struct {
}

func (u *sysInfo) HealthChecks() map[string][]*Checks {
	rs := map[string][]*Checks{
		"uptime":             make([]*Checks, 0),
		"hostname":           make([]*Checks, 0),
		"cpu:utilization":    make([]*Checks, 0),
		"memory:utilization": make([]*Checks, 0),
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	var (
		err  error
		info *host.InfoStat
		hn   string // hostname
		avg  *load.AvgStat
		vm   *mem.VirtualMemoryStat
	)
	var uptime = &Checks{
		ComponentType: "system",
		Status:        Fail,
		Time:          now,
	}
	var processes = uptime

	if info, err = host.Info(); err != nil {
		uptime.Output = err.Error()
		processes = uptime
	} else {
		uptime.Status = Pass
		uptime.ObservedValue = info.Uptime
		uptime.ObservedUnit = "s"
		processes.ComponentID = "Processes"
		processes.Status = Pass
		processes.ObservedValue = info.Procs
	}
	rs["uptime"] = append(rs["uptime"], uptime)
	rs["cpu:utilization"] = append(rs["cpu:utilization"], processes)

	var hostname = &Checks{
		ComponentID:   "hostname",
		ComponentType: "system",
		Status:        Fail,
		Time:          now,
	}
	if hn, err = os.Hostname(); err != nil {
		hostname.Output = err.Error()
	} else {
		hostname.Status = Pass
		hostname.ObservedValue = hn
	}
	rs["hostname"] = append(rs["hostname"], hostname)

	var cpuUtil = func(componentId string) *Checks {
		return &Checks{
			ComponentType: "system",
			ComponentID:   componentId,
			Status:        Fail,
			Time:          now,
		}
	}
	var cpuUtils = []*Checks{
		cpuUtil("1 minute"),
		cpuUtil("5 minute"),
		cpuUtil("15 minute"),
	}
	if avg, err = load.Avg(); err != nil {
		for _, util := range cpuUtils {
			util.Output = err.Error()
		}
	} else {
		for _, util := range cpuUtils {
			util.ObservedUnit = "%"
			util.Status = Pass
		}
		cpuUtils[0].ObservedValue = avg.Load1 / 65536.0
		cpuUtils[1].ObservedValue = avg.Load5 / 65536.0
		cpuUtils[2].ObservedValue = avg.Load15 / 65536.0
	}
	rs["cpu:utilization"] = append(cpuUtils, rs["cpu:utilization"]...)

	var memUtil = func(componentId string) *Checks {
		return &Checks{
			ComponentType: "system",
			ComponentID:   componentId,
			Status:        Fail,
			Time:          now,
		}
	}
	var memUtils = []*Checks{
		memUtil("Total Ram"),
		memUtil("Free Ram"),
		memUtil("Shared Ram"),
		memUtil("Buffer Ram"),
		memUtil("Total Swap"),
		memUtil("Free Swap"),
		memUtil("Total High"),
		memUtil("Free High"),
	}
	if vm, err = mem.VirtualMemory(); err != nil {
		for _, util := range memUtils {
			util.Output = err.Error()
		}
	} else {
		memUnit := fmt.Sprintf("%d bytes", vm.Total)
		for _, util := range memUtils {
			util.ObservedUnit = memUnit
			util.Status = Pass
		}
		memUtils[0].ObservedValue = vm.Total
		memUtils[1].ObservedValue = vm.Free
		memUtils[2].ObservedValue = vm.Shared
		memUtils[3].ObservedValue = vm.Buffers
		memUtils[4].ObservedValue = vm.SwapTotal
		memUtils[5].ObservedValue = vm.SwapFree
		memUtils[6].ObservedValue = vm.HighTotal
		memUtils[7].ObservedValue = vm.HighFree
	}
	rs["memory:utilization"] = append(rs["memory:utilization"], memUtils...)
	return rs
}

func (*sysInfo) AuthorizeHealth(*http.Request) bool {
	return true
}

// SysInfoHealth returns a ChecksProvider that provides sysinfo statistics.
func SysInfoHealth() ChecksProvider {
	return &sysInfo{}
}
