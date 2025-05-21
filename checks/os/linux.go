package os

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"runtime"
	"time"
)

type linuxChecker struct{}

func newLinuxChecker() OSChecker {
	return &linuxChecker{}
}

func (l *linuxChecker) GatherOSInfo() (*OSInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %v", err)
	}

	return &OSInfo{
		Name:      runtime.GOOS,
		Arch:      runtime.GOARCH,
		CPUs:      runtime.NumCPU(),
		GoVersion: runtime.Version(),
		Uptime:    time.Duration(hostInfo.Uptime * uint64(time.Second)).String(),
	}, nil
}

func (l *linuxChecker) GatherMemoryInfo() (*MemoryInfo, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %v", err)
	}

	return &MemoryInfo{
		Total:   memInfo.Total,
		Used:    memInfo.Used,
		Percent: memInfo.UsedPercent,
	}, nil
}

func (l *linuxChecker) GatherDiskInfo() ([]DiskInfo, error) {
	var disks []DiskInfo
	partitions, err := disk.Partitions(true) // Include all mount points on Linux
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %v", err)
	}

	for _, partition := range partitions {
		// Skip pseudo filesystems
		if partition.Fstype == "tmpfs" || partition.Fstype == "devtmpfs" {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		disks = append(disks, DiskInfo{
			Path:        partition.Mountpoint,
			Total:       usage.Total,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		})
	}

	return disks, nil
}

func (l *linuxChecker) GatherCPUInfo() ([]CPUInfo, error) {
	var cpus []CPUInfo
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %v", err)
	}

	cpuPercent, err := cpu.Percent(time.Second, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	}

	for i, cpu := range cpuInfo {
		cpus = append(cpus, CPUInfo{
			ModelName: cpu.ModelName,
			Cores:     cpu.Cores,
			Mhz:       cpu.Mhz,
			Usage:     []float64{cpuPercent[i]},
		})
	}

	return cpus, nil
}

func (l *linuxChecker) GatherNetworkInfo() ([]NetworkInfo, error) {
	var networks []NetworkInfo
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %v", err)
	}

	for _, iface := range interfaces {
		// Skip loopback interface on Linux
		if iface.Name == "lo" {
			continue
		}

		ioStats, err := net.IOCounters(true)
		if err != nil {
			continue
		}

		for _, io := range ioStats {
			if io.Name == iface.Name {
				var addresses []string
				for _, addr := range iface.Addrs {
					addresses = append(addresses, addr.Addr)
				}

				networks = append(networks, NetworkInfo{
					Name:        iface.Name,
					Addresses:   addresses,
					BytesSent:   io.BytesSent,
					BytesRecv:   io.BytesRecv,
					PacketsSent: io.PacketsSent,
					PacketsRecv: io.PacketsRecv,
				})
				break
			}
		}
	}

	return networks, nil
}

func (l *linuxChecker) GatherProcessInfo() (*ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %v", err)
	}

	info := &ProcessInfo{
		Total: int32(len(processes)),
	}

	// Get top CPU and memory processes
	for _, p := range processes[:10] { // Get top 10 processes
		name, _ := p.Name()
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryPercent()

		stat := ProcessStat{
			PID:      p.Pid,
			Name:     name,
			CPUUsage: cpu,
			MemUsage: float64(mem),
		}

		info.TopCPU = append(info.TopCPU, stat)
		info.TopMemory = append(info.TopMemory, stat)
	}

	return info, nil
} 