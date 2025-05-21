package os

import (
	"encoding/json"
	"fmt"
	"runtime"
)

// OSCheckCommand represents the command for checking OS information
type OSCheckCommand struct {
	Format string
	checker OSChecker
}

// NewOSCheckCommand creates a new OS check command
func NewOSCheckCommand() *OSCheckCommand {
	return &OSCheckCommand{
		Format: "json", // default format
		checker: NewOSChecker(),
	}
}

// Execute runs the OS check command
func (c *OSCheckCommand) Execute() error {
	if c.checker == nil {
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	info := &SystemInfo{}
	var err error

	// Gather OS info
	info.OS, err = c.checker.GatherOSInfo()
	if err != nil {
		return fmt.Errorf("failed to gather OS info: %v", err)
	}

	// Gather memory info
	info.Memory, err = c.checker.GatherMemoryInfo()
	if err != nil {
		return fmt.Errorf("failed to gather memory info: %v", err)
	}

	// Gather disk info
	info.Disks, err = c.checker.GatherDiskInfo()
	if err != nil {
		return fmt.Errorf("failed to gather disk info: %v", err)
	}

	// Gather CPU info
	info.CPUs, err = c.checker.GatherCPUInfo()
	if err != nil {
		return fmt.Errorf("failed to gather CPU info: %v", err)
	}

	// Gather network info
	info.Network, err = c.checker.GatherNetworkInfo()
	if err != nil {
		return fmt.Errorf("failed to gather network info: %v", err)
	}

	// Gather process info
	info.Processes, err = c.checker.GatherProcessInfo()
	if err != nil {
		return fmt.Errorf("failed to gather process info: %v", err)
	}

	return c.printResults(info)
}

// printResults outputs the system information in the specified format
func (c *OSCheckCommand) printResults(info *SystemInfo) error {
	switch c.Format {
	case "json":
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		fmt.Println(string(data))
	default:
		return fmt.Errorf("unsupported format: %s", c.Format)
	}
	return nil
} 