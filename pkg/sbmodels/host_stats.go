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

package sbmodels

// NetworkStats holds info for one NIC only (one instance per 网卡). Dimension: MAC and IPs.
type NetworkStats struct {
	MAC       string   `json:"mac"`
	IPs       []string `json:"ips"`
	IfName    string   `json:"if_name"`
	RxBytes   uint64   `json:"rx_bytes"`
	TxBytes   uint64   `json:"tx_bytes"`
	RxPackets uint64   `json:"rx_packets"`
	TxPackets uint64   `json:"tx_packets"`
	RxErrors  uint64   `json:"rx_errors"`
	TxErrors  uint64   `json:"tx_errors"`
	RxFifo    uint64   `json:"rx_fifo"`
	TxFifo    uint64   `json:"tx_fifo"`
}

// DiskStats holds mountpoint and used percent for one disk.
type DiskStats struct {
	Mountpoint  string  `json:"mountpoint"`
	UsedPercent float64 `json:"used_percent"`
}

// Stats is the stats for host metrics/info.
type Stats struct {
	CPU      float64        `json:"cpu"`
	Memory   float64        `json:"memory"`
	Disk     []DiskStats    `json:"disk"`
	Networks []NetworkStats `json:"networks"`
	Uptime   string         `json:"uptime"`
	Hostname string         `json:"hostname"`
	OS       string         `json:"os"`
	Arch     string         `json:"arch"`
	Kernel   string         `json:"kernel"`
}

// NewHostStats returns a zero-valued Stats.
func NewHostStats() *Stats {
	return &Stats{
		CPU:      0,
		Memory:   0,
		Disk:     []DiskStats{},
		Networks: []NetworkStats{},
		Uptime:   "",
		Hostname: "",
		OS:       "",
		Arch:     "",
		Kernel:   "",
	}
}
