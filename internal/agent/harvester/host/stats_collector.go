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

	"os-artificer/saber/pkg/sbmodels"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	netutil "github.com/shirou/gopsutil/v4/net"
)

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

// CollectDisk returns disk usage for physical devices only (e.g. hard disks, CD-ROM, USB) using gopsutil; virtual/memory partitions (e.g. tmpfs, /dev/shm) are excluded.
func CollectDisk() []sbmodels.DiskStats {
	partitions, err := disk.Partitions(false)
	if err != nil || len(partitions) == 0 {
		return []sbmodels.DiskStats{}
	}

	diskStats := make([]sbmodels.DiskStats, 0)
	for _, p := range partitions {
		u, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		diskStats = append(diskStats, sbmodels.DiskStats{
			Mountpoint:  p.Mountpoint,
			UsedPercent: u.UsedPercent,
		})
	}

	return diskStats
}

// CollectNetwork returns one NetworkStats per physical NIC (one 网卡 per record). Each record has that NIC's MAC, IPs, IfName, and traffic counters.
func CollectNetwork() []sbmodels.NetworkStats {
	ifaces, err := net.Interfaces()
	if err != nil {
		return []sbmodels.NetworkStats{}
	}

	counterMap := make(map[string]netutil.IOCountersStat)
	if counters, err := netutil.IOCounters(true); err == nil {
		for i := range counters {
			counterMap[counters[i].Name] = counters[i]
		}
	}

	out := make([]sbmodels.NetworkStats, 0)
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if !isPhysicalInterface(iface.Name) {
			continue
		}

		mac := ""
		if iface.HardwareAddr != nil {
			mac = iface.HardwareAddr.String()
		}
		ips := make([]string, 0)
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			if ip := ipnet.IP.To4(); ip != nil {
				ips = append(ips, ip.String())
			}
		}

		c := counterMap[iface.Name]
		out = append(out, sbmodels.NetworkStats{
			MAC:       mac,
			IPs:       ips,
			IfName:    iface.Name,
			RxBytes:   c.BytesRecv,
			TxBytes:   c.BytesSent,
			RxPackets: c.PacketsRecv,
			TxPackets: c.PacketsSent,
			RxErrors:  c.Errin,
			TxErrors:  c.Errout,
			RxFifo:    c.Fifoin,
			TxFifo:    c.Fifoout,
		})
	}

	return out
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

// CollectStats fills s with collected host metrics/info.
func CollectStats(s *sbmodels.Stats) error {
	s.CPU = CollectCPU()
	s.Memory = CollectMemory()

	s.Disk = CollectDisk()
	s.Networks = CollectNetwork()

	s.Uptime = CollectUptime()
	s.Hostname = CollectHostname()

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
