//go:build linux

package os

import (
	"bufio"
	"fmt"
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

// GetLinuxSoftwareInfo retrieves installed software information from Linux
func GetLinuxSoftwareInfo() ([]SoftwareInfo, error) {
	// Try dpkg first (Debian/Ubuntu)
	if _, err := exec.LookPath("dpkg"); err == nil {
		cmd := exec.Command("dpkg", "-l")
		output, err := cmd.Output()
		if err == nil {
			return parseDebianPackages(string(output))
		}
	}

	// Try rpm (Red Hat/CentOS/Fedora)
	if _, err := exec.LookPath("rpm"); err == nil {
		cmd := exec.Command("rpm", "-qa", "--queryformat", "%{NAME}\t%{VERSION}\t%{VENDOR}\t%{INSTALLTIME:date}\n")
		output, err := cmd.Output()
		if err == nil {
			return parseRPMPackages(string(output))
		}
	}

	return nil, fmt.Errorf("no supported package manager found")
}

func parseDebianPackages(output string) ([]SoftwareInfo, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var software []SoftwareInfo

	// Skip header lines
	for i := 0; i < 5 && scanner.Scan(); i++ {}

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}

		// Fields are: status, name, version, architecture, description
		if !strings.HasPrefix(fields[0], "ii") {
			continue
		}

		software = append(software, SoftwareInfo{
			Name:      fields[1],
			Version:   fields[2],
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

func parseRPMPackages(output string) ([]SoftwareInfo, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var software []SoftwareInfo

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) < 4 {
			continue
		}

		software = append(software, SoftwareInfo{
			Name:      fields[0],
			Version:   fields[1],
			Publisher: fields[2],
			InstallDate: fields[3],
		})
	}

	// Sort by name
	sort.Slice(software, func(i, j int) bool {
		return strings.ToLower(software[i].Name) < strings.ToLower(software[j].Name)
	})

	return software, nil
}

// GetLinuxDateFormat returns the date format for Linux systems
func GetLinuxDateFormat() string {
	return "2006-01-02" // Use ISO format for Linux systems
} 