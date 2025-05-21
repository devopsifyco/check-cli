//go:build windows

package version

import (
	"fmt"
	"os/exec"
	"strings"
	"golang.org/x/sys/windows/registry"
)

// GetWindowsVersion gets the version of any component on Windows
func GetWindowsVersion(component string) (string, error) {
	// First try using WMIC to find the product
	cmd := exec.Command("wmic", "product", "where", fmt.Sprintf("name like '%%%s%%'", component), "get", "version")
	output, err := cmd.Output()
	if err == nil {
		// Parse WMIC output
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && line != "Version" {
				return line, nil
			}
		}
	}

	// If WMIC fails, try to find the executable and get its version
	cmd = exec.Command("where", component)
	output, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not find %s executable", component)
	}

	path := strings.TrimSpace(string(output))
	if path == "" {
		return "", fmt.Errorf("could not find %s executable", component)
	}

	// Try to get version from the executable
	return GetVersionFromExecutable(path)
}

// GetWindowsDateFormat reads the date format from Windows Registry
func GetWindowsDateFormat() string {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\International`, registry.QUERY_VALUE)
	if err != nil {
		return "1/2/2006" // Default US format
	}
	defer k.Close()

	sShortDate, _, err := k.GetStringValue("sShortDate")
	if err != nil {
		return "1/2/2006" // Default US format
	}

	// Convert Windows date format to Go date format
	format := strings.NewReplacer(
		"dddd", "Monday",
		"ddd", "Mon",
		"dd", "02",
		"d", "2",
		"MMMM", "January",
		"MMM", "Jan",
		"MM", "01",
		"M", "1",
		"yyyy", "2006",
		"yy", "06",
	).Replace(sShortDate)

	// Handle separators
	format = strings.ReplaceAll(format, "/", "/")
	format = strings.ReplaceAll(format, "-", "-")
	format = strings.ReplaceAll(format, ".", ".")

	return format
} 