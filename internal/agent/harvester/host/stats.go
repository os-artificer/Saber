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

package host

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	netutil "github.com/shirou/gopsutil/v4/net"
)

// Stats is the stats for host metrics/info.
type Stats struct {
	CPU      float64 `json:"cpu"`
	Memory   float64 `json:"memory"`
	Disk     string  `json:"disk"`
	Network  string  `json:"network"`
	Uptime   string  `json:"uptime"`
	Hostname string  `json:"hostname"`
	IP       string  `json:"ip"`
	MAC      string  `json:"mac"`
	OS       string  `json:"os"`
	Arch     string  `json:"arch"`
	Kernel   string  `json:"kernel"`
}

// NewStats creates a new Stats.
func NewStats() *Stats {
	return &Stats{
		CPU:      0,
		Memory:   0,
		Disk:     "",
		Network:  "",
		Uptime:   "",
		Hostname: "",
		IP:       "",
		MAC:      "",
		OS:       "",
		Arch:     "",
		Kernel:   "",
	}
}

// CollectCPU returns CPU usage percentage (0-100) using gopsutil.
func CollectCPU() float64 {
	percent, err := cpu.Percent(100*time.Millisecond, false)
	if err != nil || len(percent) == 0 {
		return 0
	}
	return percent[0]
}

// CollectMemory returns memory usage percentage (0-100) using gopsutil.
func CollectMemory() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return v.UsedPercent
}

// CollectDisk returns root filesystem usage string (e.g. "45.2%") using gopsutil.
func CollectDisk() string {
	u, err := disk.Usage("/")
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f%%", u.UsedPercent)
}

// CollectNetwork returns "rx: N tx: N" (bytes) using gopsutil.
func CollectNetwork() string {
	counters, err := netutil.IOCounters(false)
	if err != nil || len(counters) == 0 {
		return ""
	}
	// pernic=false returns one entry with name "all"
	for i := range counters {
		if counters[i].Name == "all" {
			return fmt.Sprintf("rx: %d tx: %d", counters[i].BytesRecv, counters[i].BytesSent)
		}
	}
	var rx, tx uint64
	for i := range counters {
		if strings.HasPrefix(counters[i].Name, "lo") {
			continue
		}
		rx += counters[i].BytesRecv
		tx += counters[i].BytesSent
	}
	return fmt.Sprintf("rx: %d tx: %d", rx, tx)
}

// CollectUptime returns uptime string (e.g. "3d12h") using gopsutil.
func CollectUptime() string {
	sec, err := host.Uptime()
	if err != nil {
		return ""
	}
	d := time.Duration(sec) * time.Second
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dd%dh%dm", days, hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}

// CollectHostname returns the hostname using gopsutil.
func CollectHostname() string {
	info, err := host.Info()
	if err != nil {
		return ""
	}
	return info.Hostname
}

// CollectIP returns the first IPv4 address of a physical network interface.
func CollectIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if !isPhysicalInterface(iface.Name) {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			if ip := ipnet.IP.To4(); ip != nil {
				return ip.String()
			}
		}
	}
	return ""
}

// CollectMAC returns the first non-loopback interface MAC address.
func CollectMAC() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.HardwareAddr == nil {
			continue
		}
		if strings.HasPrefix(iface.Name, "lo") {
			continue
		}
		return iface.HardwareAddr.String()
	}
	return ""
}

// CollectOS returns OS/platform info using gopsutil host.
func CollectOS() string {
	info, err := host.Info()
	if err != nil {
		return ""
	}
	if info.Platform != "" && info.PlatformVersion != "" {
		return info.Platform + " " + info.PlatformVersion
	}
	if info.Platform != "" {
		return info.Platform
	}
	return info.OS
}

// CollectArch returns the kernel/machine architecture using gopsutil.
func CollectArch() string {
	arch, err := host.KernelArch()
	if err != nil {
		return ""
	}
	return arch
}

// CollectKernel returns kernel version using gopsutil.
func CollectKernel() string {
	version, err := host.KernelVersion()
	if err != nil {
		return ""
	}
	return version
}

// CollectStats collects the stats for host metrics/info.
func (s *Stats) CollectStats() error {
	s.CPU = CollectCPU()
	s.Memory = CollectMemory()
	s.Disk = CollectDisk()
	s.Network = CollectNetwork()
	s.Uptime = CollectUptime()
	s.Hostname = CollectHostname()
	s.IP = CollectIP()
	s.MAC = CollectMAC()
	s.OS = CollectOS()
	s.Arch = CollectArch()
	s.Kernel = CollectKernel()
	return nil
}

// isPhysicalInterface returns true if the interface name looks like a real physical NIC.
// Excludes loopback, bridges, veth, docker, virbr, tun/tap and other virtual interfaces.
func isPhysicalInterface(name string) bool {
	if name == "lo" {
		return false
	}
	lower := strings.ToLower(name)
	virtualPrefixes := []string{
		"veth", "docker", "br-", "virbr", "vb-", "tun", "tap",
		"cali", "flannel", "cni", "kube", "ovs", "vlan",
	}
	for _, p := range virtualPrefixes {
		if strings.HasPrefix(lower, p) {
			return false
		}
	}
	return true
}
