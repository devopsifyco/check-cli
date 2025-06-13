package net

import (
	"fmt"
	"runtime"
	"time"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"sort"
	"github.com/devopsifyco/check-cli/checks"
)

// OSResult implements CheckResult interface for OS checks
type OSResult struct {
	OS            string                 `json:"os"`
	Arch          string                 `json:"arch"`
	CPUs          int                    `json:"cpus"`
	GoVersion     string                 `json:"go_version"`
	Uptime        string                 `json:"uptime"`
	MemoryTotal   uint64                 `json:"memory_total"`
	MemoryUsed    uint64                 `json:"memory_used"`
	MemoryPercent float64                `json:"memory_percent"`
	DiskInfo      []DiskInfo             `json:"disk_info"`
	HostInfo      map[string]interface{} `json:"host_info"`
	CPUInfo       []CPUInfo              `json:"cpu_info"`
	NetInfo       []NetInfo              `json:"net_info"`
	ProcessInfo   ProcessInfo            `json:"process_info"`
	SoftwareInfo  []SoftwareInfo         `json:"software_info"`
}

type DiskInfo struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type CPUInfo struct {
	ModelName string    `json:"model_name"`
	Cores     int32     `json:"cores"`
	Mhz       float64   `json:"mhz"`
	Usage     []float64 `json:"usage"`
}

type NetInfo struct {
	Name        string   `json:"name"`
	Addresses   []string `json:"addresses"`
	BytesSent   uint64   `json:"bytes_sent"`
	BytesRecv   uint64   `json:"bytes_recv"`
	PacketsSent uint64   `json:"packets_sent"`
	PacketsRecv uint64   `json:"packets_recv"`
}

type ProcessInfo struct {
	Total     int32         `json:"total"`
	TopCPU    []ProcessStat `json:"top_cpu"`
	TopMemory []ProcessStat `json:"top_memory"`
}

type ProcessStat struct {
	PID      int32   `json:"pid"`
	Name     string  `json:"name"`
	CPUUsage float64 `json:"cpu_usage"`
	MemUsage float64 `json:"mem_usage"`
}

type SoftwareInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Publisher   string `json:"publisher"`
	InstallDate string `json:"install_date"`
}

// Print implements CheckResult interface
func (r *OSResult) Print(outputFormat string) {
	switch outputFormat {
	case "json":
		checks.PrintJSON(r)
	case "yaml":
		checks.PrintYAML(r)
	default:
		fmt.Printf("OS: %s\n", r.OS)
		fmt.Printf("Architecture: %s\n", r.Arch)
		fmt.Printf("CPUs: %d\n", r.CPUs)
		fmt.Printf("Go Version: %s\n", r.GoVersion)
		fmt.Printf("Uptime: %s\n", r.Uptime)

		fmt.Printf("\nInstalled Software:\n")
		for _, sw := range r.SoftwareInfo {
			fmt.Printf("Name: %s\n", sw.Name)
			fmt.Printf("  Version: %s\n", sw.Version)
			fmt.Printf("  Publisher: %s\n", sw.Publisher)
			fmt.Printf("  Install Date: %s\n\n", sw.InstallDate)
		}

		fmt.Printf("\nCPU Information:\n")
		for i, cpu := range r.CPUInfo {
			fmt.Printf("CPU %d:\n", i)
			fmt.Printf("  Model: %s\n", cpu.ModelName)
			fmt.Printf("  Cores: %d\n", cpu.Cores)
			fmt.Printf("  Frequency: %.2f MHz\n", cpu.Mhz)
			for j, usage := range cpu.Usage {
				fmt.Printf("  Core %d Usage: %.2f%%\n", j, usage)
			}
		}

		fmt.Printf("\nMemory Information:\n")
		fmt.Printf("Total: %.2f GB\n", float64(r.MemoryTotal)/(1024*1024*1024))
		fmt.Printf("Used: %.2f GB\n", float64(r.MemoryUsed)/(1024*1024*1024))
		fmt.Printf("Usage: %.2f%%\n", r.MemoryPercent)

		fmt.Printf("\nDisk Information:\n")
		for _, disk := range r.DiskInfo {
			fmt.Printf("Path: %s\n", disk.Path)
			fmt.Printf("Total: %.2f GB\n", float64(disk.Total)/(1024*1024*1024))
			fmt.Printf("Used: %.2f GB\n", float64(disk.Used)/(1024*1024*1024))
			fmt.Printf("Usage: %.2f%%\n\n", disk.UsedPercent)
		}

		fmt.Printf("Network Information:\n")
		for _, net := range r.NetInfo {
			fmt.Printf("Interface: %s\n", net.Name)
			fmt.Printf("  Addresses: %v\n", net.Addresses)
			fmt.Printf("  Bytes Sent: %.2f MB\n", float64(net.BytesSent)/(1024*1024))
			fmt.Printf("  Bytes Received: %.2f MB\n", float64(net.BytesRecv)/(1024*1024))
			fmt.Printf("  Packets Sent: %d\n", net.PacketsSent)
			fmt.Printf("  Packets Received: %d\n\n", net.PacketsRecv)
		}

		fmt.Printf("Process Information:\n")
		fmt.Printf("Total Processes: %d\n", r.ProcessInfo.Total)
		fmt.Printf("\nTop CPU Processes:\n")
		for _, p := range r.ProcessInfo.TopCPU {
			fmt.Printf("  PID: %d, Name: %s, CPU: %.2f%%, Memory: %.2f%%\n", p.PID, p.Name, p.CPUUsage, p.MemUsage)
		}
		fmt.Printf("\nTop Memory Processes:\n")
		for _, p := range r.ProcessInfo.TopMemory {
			fmt.Printf("  PID: %d, Name: %s, CPU: %.2f%%, Memory: %.2f%%\n", p.PID, p.Name, p.CPUUsage, p.MemUsage)
		}

		fmt.Printf("\nHost Information:\n")
		for k, v := range r.HostInfo {
			fmt.Printf("%s: %v\n", k, v)
		}
	}
}

// OSCheckCommand implements the CheckCommand interface for OS checks
type OSCheckCommand struct {
	*checks.BaseCheckCommand
}

// NewOSCheckCommand creates a new OS check command
func NewOSCheckCommand() *OSCheckCommand {
	return &OSCheckCommand{
		BaseCheckCommand: checks.NewBaseCheckCommand(
			"os",
			"Check operating system information",
			"os",
			0,
		),
	}
}

// getSoftwareInfo retrieves installed software information based on the platform
func getSoftwareInfo() ([]SoftwareInfo, error) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsSoftwareInfo()
	case "darwin":
		return getDarwinSoftwareInfo()
	case "linux":
		return getLinuxSoftwareInfo()
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// getWindowsSoftwareInfo retrieves installed software information from Windows registry
func getWindowsSoftwareInfo() ([]SoftwareInfo, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("not running on Windows")
	}

	var software []SoftwareInfo
	// TODO: Implement Windows registry scanning
	// This is a placeholder that returns an empty list for now
	return software, nil
}

// getDarwinSoftwareInfo retrieves installed software information from macOS
func getDarwinSoftwareInfo() ([]SoftwareInfo, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("not running on macOS")
	}

	var software []SoftwareInfo
	// TODO: Implement macOS software scanning
	// This is a placeholder that returns an empty list for now
	return software, nil
}

// getLinuxSoftwareInfo retrieves installed software information from Linux package managers
func getLinuxSoftwareInfo() ([]SoftwareInfo, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("not running on Linux")
	}

	var software []SoftwareInfo
	// TODO: Implement Linux package manager scanning
	// This is a placeholder that returns an empty list for now
	return software, nil
}

// Execute implements the CheckCommand interface
func (c *OSCheckCommand) Execute(args []string) (checks.CheckResult, error) {
	// Get memory information
	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %v", err)
	}

	// Get disk information
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %v", err)
	}

	diskInfo := make([]DiskInfo, 0)
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		diskInfo = append(diskInfo, DiskInfo{
			Path:        partition.Mountpoint,
			Total:       usage.Total,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		})
	}

	// Get CPU information
	cpuInfo := make([]CPUInfo, 0)
	cpus, err := cpu.Info()
	if err == nil {
		percentages, err := cpu.Percent(time.Second, true)
		if err == nil {
			for _, cpu := range cpus {
				cpuInfo = append(cpuInfo, CPUInfo{
					ModelName: cpu.ModelName,
					Cores:     cpu.Cores,
					Mhz:       cpu.Mhz,
					Usage:     percentages,
				})
			}
		}
	}

	// Get network information
	netInfo := make([]NetInfo, 0)
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			addrs := make([]string, 0)
			for _, addr := range iface.Addrs {
				addrs = append(addrs, addr.Addr)
			}

			ioStats, err := net.IOCounters(true)
			if err == nil {
				for _, io := range ioStats {
					if io.Name == iface.Name {
						netInfo = append(netInfo, NetInfo{
							Name:        iface.Name,
							Addresses:   addrs,
							BytesSent:   io.BytesSent,
							BytesRecv:   io.BytesRecv,
							PacketsSent: io.PacketsSent,
							PacketsRecv: io.PacketsRecv,
						})
						break
					}
				}
			}
		}
	}

	// Get process information
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get process info: %v", err)
	}

	processStats := make([]ProcessStat, 0)
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		cpu, err := p.CPUPercent()
		if err != nil {
			continue
		}

		mem, err := p.MemoryPercent()
		if err != nil {
			continue
		}

		processStats = append(processStats, ProcessStat{
			PID:      p.Pid,
			Name:     name,
			CPUUsage: cpu,
			MemUsage: float64(mem),
		})
	}

	// Sort processes by CPU and memory usage
	sort.Slice(processStats, func(i, j int) bool {
		return processStats[i].CPUUsage > processStats[j].CPUUsage
	})
	topCPU := processStats
	if len(topCPU) > 5 {
		topCPU = topCPU[:5]
	}

	sort.Slice(processStats, func(i, j int) bool {
		return processStats[i].MemUsage > processStats[j].MemUsage
	})
	topMemory := processStats
	if len(topMemory) > 5 {
		topMemory = topMemory[:5]
	}

	// Get host information
	hostInfo, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %v", err)
	}

	hostInfoMap := map[string]interface{}{
		"Hostname":        hostInfo.Hostname,
		"Platform":        hostInfo.Platform,
		"PlatformVersion": hostInfo.PlatformVersion,
		"KernelVersion":   hostInfo.KernelVersion,
		"KernelArch":      hostInfo.KernelArch,
	}

	// Get software information
	softwareInfo, err := getSoftwareInfo()
	if err != nil {
		// Log the error but don't fail the entire command
		fmt.Printf("Warning: Failed to get software info: %v\n", err)
	}

	return &OSResult{
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		CPUs:          runtime.NumCPU(),
		GoVersion:     runtime.Version(),
		Uptime:        time.Duration(hostInfo.Uptime * uint64(time.Second)).String(),
		MemoryTotal:   memory.Total,
		MemoryUsed:    memory.Used,
		MemoryPercent: memory.UsedPercent,
		DiskInfo:      diskInfo,
		HostInfo:      hostInfoMap,
		CPUInfo:       cpuInfo,
		NetInfo:       netInfo,
		ProcessInfo: ProcessInfo{
			Total:     int32(len(processes)),
			TopCPU:    topCPU,
			TopMemory: topMemory,
		},
		SoftwareInfo:  softwareInfo,
	}, nil
} 