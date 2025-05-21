package os

// OSInfo represents the core operating system information
type OSInfo struct {
	Name        string `json:"name"`
	Arch        string `json:"arch"`
	CPUs        int    `json:"cpus"`
	GoVersion   string `json:"go_version"`
	Uptime      string `json:"uptime"`
}

// MemoryInfo represents memory usage information
type MemoryInfo struct {
	Total   uint64  `json:"total"`
	Used    uint64  `json:"used"`
	Percent float64 `json:"percent"`
}

// DiskInfo represents disk usage information
type DiskInfo struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

// CPUInfo represents CPU information
type CPUInfo struct {
	ModelName string    `json:"model_name"`
	Cores     int32     `json:"cores"`
	Mhz       float64   `json:"mhz"`
	Usage     []float64 `json:"usage"`
}

// NetworkInfo represents network interface information
type NetworkInfo struct {
	Name        string   `json:"name"`
	Addresses   []string `json:"addresses"`
	BytesSent   uint64   `json:"bytes_sent"`
	BytesRecv   uint64   `json:"bytes_recv"`
	PacketsSent uint64   `json:"packets_sent"`
	PacketsRecv uint64   `json:"packets_recv"`
}

// ProcessInfo represents process information
type ProcessInfo struct {
	Total     int32         `json:"total"`
	TopCPU    []ProcessStat `json:"top_cpu"`
	TopMemory []ProcessStat `json:"top_memory"`
}

// ProcessStat represents individual process statistics
type ProcessStat struct {
	PID      int32   `json:"pid"`
	Name     string  `json:"name"`
	CPUUsage float64 `json:"cpu_usage"`
	MemUsage float64 `json:"mem_usage"`
}

// SoftwareInfo represents installed software information
type SoftwareInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Publisher   string `json:"publisher"`
	InstallDate string `json:"install_date"`
}

// SystemInfo represents the complete system information
type SystemInfo struct {
	OS        OSInfo                 `json:"os"`
	Memory    MemoryInfo            `json:"memory"`
	Disks     []DiskInfo           `json:"disks"`
	Host      map[string]interface{} `json:"host"`
	CPUs      []CPUInfo            `json:"cpus"`
	Network   []NetworkInfo        `json:"network"`
	Processes ProcessInfo          `json:"processes"`
	Software  []SoftwareInfo       `json:"software"`
} 