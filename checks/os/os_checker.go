package os

// OSChecker defines the interface for OS-specific system information gathering
type OSChecker interface {
	// GatherOSInfo collects basic OS information
	GatherOSInfo() (*OSInfo, error)
	
	// GatherMemoryInfo collects memory usage information
	GatherMemoryInfo() (*MemoryInfo, error)
	
	// GatherDiskInfo collects disk usage information
	GatherDiskInfo() ([]DiskInfo, error)
	
	// GatherCPUInfo collects CPU information and usage
	GatherCPUInfo() ([]CPUInfo, error)
	
	// GatherNetworkInfo collects network interface information
	GatherNetworkInfo() ([]NetworkInfo, error)
	
	// GatherProcessInfo collects process information
	GatherProcessInfo() (*ProcessInfo, error)
}

// NewOSChecker creates a new OS-specific checker based on the current platform
func NewOSChecker() OSChecker {
	switch runtime.GOOS {
	case "windows":
		return newWindowsChecker()
	case "linux":
		return newLinuxChecker()
	case "darwin":
		return newDarwinChecker()
	default:
		return nil
	}
} 