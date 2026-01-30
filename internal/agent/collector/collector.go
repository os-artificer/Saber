/**
 * Copyright 2025 Saber authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/

package collector

import (
	"encoding/json"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// HostMetrics host performance metrics for transfer
type HostMetrics struct {
	Hostname   string  `json:"hostname"`
	Timestamp  int64   `json:"timestamp"`
	CPUPercent float64 `json:"cpu_percent"`
	MemPercent float64 `json:"mem_percent"`
	MemUsedMB  uint64  `json:"mem_used_mb"`
	MemTotalMB uint64  `json:"mem_total_mb"`
	Uptime     uint64  `json:"uptime_sec"`
	BootTime   uint64  `json:"boot_time"`
}

// Collect gathers host metrics and returns JSON payload
func Collect() ([]byte, error) {
	hostname, _ := os.Hostname()
	now := time.Now().Unix()

	cpuPercents, err := cpu.Percent(0, false)
	cpuPercent := 0.0
	if err == nil && len(cpuPercents) > 0 {
		cpuPercent = cpuPercents[0]
	}

	vm, err := mem.VirtualMemory()
	memPercent := 0.0
	memUsedMB := uint64(0)
	memTotalMB := uint64(0)
	if err == nil && vm != nil {
		memPercent = vm.UsedPercent
		memUsedMB = vm.Used / (1024 * 1024)
		memTotalMB = vm.Total / (1024 * 1024)
	}

	hi, err := host.Info()
	uptime := uint64(0)
	bootTime := uint64(0)
	if err == nil && hi != nil {
		uptime = hi.Uptime
		bootTime = hi.BootTime
	}

	m := HostMetrics{
		Hostname:   hostname,
		Timestamp:  now,
		CPUPercent: cpuPercent,
		MemPercent: memPercent,
		MemUsedMB:  memUsedMB,
		MemTotalMB: memTotalMB,
		Uptime:     uptime,
		BootTime:   bootTime,
	}
	return json.Marshal(m)
}
