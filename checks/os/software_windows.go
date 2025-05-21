//go:build windows

package os

import (
	"golang.org/x/sys/windows/registry"
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

// GetWindowsSoftwareInfo retrieves installed software information from the Windows Registry
func GetWindowsSoftwareInfo() ([]SoftwareInfo, error) {
	// Registry paths for installed software
	regPaths := []string{
		"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall",          // 64-bit applications
		"SOFTWARE\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall", // 32-bit applications
	}

	softwareMap := make(map[string]SoftwareInfo) // Use map to avoid duplicates

	for _, regPath := range regPaths {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.READ)
		if err != nil {
			continue
		}
		defer k.Close()

		// Get list of subkeys (these are the application GUIDs)
		subKeys, err := k.ReadSubKeyNames(-1) // -1 means read all subkeys
		if err != nil {
			continue
		}

		for _, subKey := range subKeys {
			subK, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath+"\\"+subKey, registry.READ)
			if err != nil {
				continue
			}
			defer subK.Close()

			// Get application name
			displayName, _, err := subK.GetStringValue("DisplayName")
			if err != nil || displayName == "" {
				continue
			}

			// Get version
			displayVersion, _, _ := subK.GetStringValue("DisplayVersion")
			publisher, _, _ := subK.GetStringValue("Publisher")
			installDate, _, _ := subK.GetStringValue("InstallDate")

			// Use display name as key to avoid duplicates
			softwareMap[displayName] = SoftwareInfo{
				Name:        displayName,
				Version:     displayVersion,
				Publisher:   publisher,
				InstallDate: installDate,
			}
		}
	}

	// Convert map to slice
	software := make([]SoftwareInfo, 0, len(softwareMap))
	for _, sw := range softwareMap {
		software = append(software, sw)
	}

	// Sort by name
	sort.Slice(software, func(i, j int) bool {
		return strings.ToLower(software[i].Name) < strings.ToLower(software[j].Name)
	})

	return software, nil
}

// GetWindowsDateFormat returns the date format for Windows systems
func GetWindowsDateFormat() string {
	return "01/02/2006" // Use US format for Windows systems
} 