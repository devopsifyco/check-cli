//go:build darwin

package os

import (
	"encoding/json"
	"os/exec"
	"strings"
	"sort"
)

// SoftwareInfo represents installed software information
type SoftwareInfo struct {
	Name        string
	Version     string
	Publisher   string
	InstallDate string
}

// GetDarwinSoftwareInfo retrieves installed software information from macOS
func GetDarwinSoftwareInfo() ([]SoftwareInfo, error) {
	// Get list of installed applications using system_profiler
	cmd := exec.Command("system_profiler", "SPApplicationsDataType", "-json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result struct {
		Apps []struct {
			Name    string `json:"_name"`
			Version string `json:"version"`
			Path    string `json:"path"`
		} `json:"SPApplicationsDataType"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	software := make([]SoftwareInfo, 0, len(result.Apps))
	for _, app := range result.Apps {
		software = append(software, SoftwareInfo{
			Name:      app.Name,
			Version:   app.Version,
			Publisher: "",
			InstallDate: "",
		})
	}

	// Sort by name
	sort.Slice(software, func(i, j int) bool {
		return strings.ToLower(software[i].Name) < strings.ToLower(software[j].Name)
	})

	return software, nil
}

// GetDarwinDateFormat returns the date format for macOS systems
func GetDarwinDateFormat() string {
	return "2006-01-02" // Use ISO format for macOS systems
} 