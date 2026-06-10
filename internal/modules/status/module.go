package status

import (
	"runtime"
	"time"
)

type HostOverview struct {
	Service       string    `json:"service"`
	Version       string    `json:"version"`
	Runtime       string    `json:"runtime"`
	StartedAt     time.Time `json:"startedAt"`
	UptimeSeconds int64     `json:"uptimeSeconds"`
	OS            string    `json:"os"`
	Arch          string    `json:"arch"`
	GoVersion     string    `json:"goVersion"`
	Goroutines    int       `json:"goroutines"`
	Memory        Memory    `json:"memory"`
}

type Memory struct {
	AllocBytes     uint64 `json:"allocBytes"`
	SysBytes       uint64 `json:"sysBytes"`
	HeapInuseBytes uint64 `json:"heapInuseBytes"`
	HeapIdleBytes  uint64 `json:"heapIdleBytes"`
	NumGC          uint32 `json:"numGC"`
}

func Host(startedAt time.Time) HostOverview {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	now := time.Now()
	return HostOverview{
		Service:       "OpenLucky",
		Version:       "0.1.0-dev",
		Runtime:       "go+hertz",
		StartedAt:     startedAt,
		UptimeSeconds: int64(now.Sub(startedAt).Seconds()),
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		GoVersion:     runtime.Version(),
		Goroutines:    runtime.NumGoroutine(),
		Memory: Memory{
			AllocBytes:     stats.Alloc,
			SysBytes:       stats.Sys,
			HeapInuseBytes: stats.HeapInuse,
			HeapIdleBytes:  stats.HeapIdle,
			NumGC:          stats.NumGC,
		},
	}
}
