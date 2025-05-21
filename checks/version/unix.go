//go:build !windows

package version

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GetUnixVersion gets the version of any component on Unix-like systems
func GetUnixVersion(component string) (string, error) {
	// First try using package managers
	switch runtime.GOOS {
	case "darwin":
		// Try brew first
		cmd := exec.Command("brew", "list", "--versions", component)
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					return fields[1], nil
				}
			}
		}

		// Then try pkgutil
		cmd = exec.Command("pkgutil", "--pkg-info", component)
		output, err = cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "version:") {
					return strings.TrimSpace(strings.TrimPrefix(line, "version:")), nil
				}
			}
		}
	default:
		// Try dpkg first
		cmd := exec.Command("dpkg", "-l", component)
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) >= 3 && strings.HasPrefix(fields[0], "ii") {
					return fields[2], nil
				}
			}
		}

		// Then try rpm
		cmd = exec.Command("rpm", "-q", "--queryformat", "%{VERSION}", component)
		output, err = cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output)), nil
		}
	}

	// If package managers fail, try to find the executable
	cmd := exec.Command("which", component)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	path := strings.TrimSpace(string(output))
	if path == "" {
		return "", err
	}

	return GetVersionFromExecutable(path)
}

// GetUnixDateFormat reads the date format from system locale
func GetUnixDateFormat() string {
	// Try to get LC_TIME from environment
	lcTime := os.Getenv("LC_TIME")
	if lcTime == "" {
		lcTime = os.Getenv("LANG")
	}
	if lcTime == "" {
		return "2006-01-02" // Default ISO format
	}

	var dateFormat string

	switch runtime.GOOS {
	case "darwin":
		// On macOS, try to get date format using defaults command
		cmd := exec.Command("defaults", "read", "NSGlobalDomain", "AppleICUDateFormatStrings")
		output, err := cmd.Output()
		if err == nil {
			// Parse the output to get the short date format
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "ShortFormat") {
					dateFormat = strings.TrimSpace(strings.Split(line, "=")[1])
					break
				}
			}
		}
	case "linux":
		// Try to get date format using locale command
		cmd := exec.Command("locale", "d_fmt")
		cmd.Env = append(os.Environ(), "LC_TIME="+lcTime)
		output, err := cmd.Output()
		if err == nil {
			dateFormat = strings.TrimSpace(string(output))
		}

		// If locale command fails, try date command
		if dateFormat == "" {
			cmd = exec.Command("date", "+%x")
			output, err = cmd.Output()
			if err == nil {
				dateFormat = strings.TrimSpace(string(output))
			}
		}
	}

	// If we couldn't get the system format, use ISO format
	if dateFormat == "" {
		return "2006-01-02"
	}

	// Convert locale format to Go date format
	format := strings.NewReplacer(
		"%A", "Monday",
		"%a", "Mon",
		"%B", "January",
		"%b", "Jan",
		"%d", "02",
		"%e", "2",
		"%m", "01",
		"%Y", "2006",
		"%y", "06",
		"%F", "2006-01-02",
		"%D", "01/02/06",
	).Replace(dateFormat)

	// Clean up any remaining % sequences
	format = strings.ReplaceAll(format, "%", "")

	// Handle common separators
	format = strings.ReplaceAll(format, "/", "/")
	format = strings.ReplaceAll(format, "-", "-")
	format = strings.ReplaceAll(format, ".", ".")

	return format
} 