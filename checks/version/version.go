package version

import (
	"fmt"
	"os/exec"
	"strings"
	"github.com/devopsifyco/check-cli/checks/utilities/output"
)

var (
	Version   = "0.0.16"
	Revision  = "e4d08a7"
	BuildDate = "2025-06-10"
)

// Result implements CheckResult interface for version checks
type Result struct {
	Version string `json:"version" yaml:"version"`
	Full    bool   `json:"full,omitempty" yaml:"full,omitempty"`
}

// Print implements CheckResult interface
func (r *Result) Print(outputFormat string) {
	switch outputFormat {
	case "json":
		output.PrintJSON(r)
	case "yaml":
		output.PrintYAML(r)
	default:
		fmt.Println(r.Version)
	}
}

// CheckLocal checks the version of a component using local client
func CheckLocal(component string, outputFormat string, fullOutput bool) (*Result, error) {
	var version string
	var err error

	switch component {
	case "cli":
		version, err = GetCLIVersion()
	case "docker":
		version, err = GetDockerVersion()
	case "kubectl":
		version, err = GetKubectlVersion()
	case "helm":
		version, err = GetHelmVersion()
	default:
		return nil, fmt.Errorf("unsupported component: %s", component)
	}

	if err != nil {
		return nil, err
	}

	return &Result{
		Version: version,
		Full:    fullOutput,
	}, nil
}

// GetCLIVersion gets the version of the CLI tool itself
func GetCLIVersion() (string, error) {
	return Version, nil
}

// GetCLIMetadata returns version, revision, and build date
func GetCLIMetadata() map[string]string {
	return map[string]string{
		"version": Version,
		"revision": Revision,
		"buildDate": BuildDate,
	}
}

// PrintCLIVersionInfo prints version info in the requested format
func PrintCLIVersionInfo(outputFormat string) {
	info := GetCLIMetadata()
	switch outputFormat {
	case "json":
		output.PrintJSON(info)
	case "yaml":
		output.PrintYAML(info)
	default:
		fmt.Printf("check version %s (rev: %s, date: %s)\n", Version, Revision, BuildDate)
	}
}

// GetDockerVersion gets the version of Docker
func GetDockerVersion() (string, error) {
	cmd := exec.Command("docker", "version", "--format", "{{.Client.Version}}")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Docker version: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetKubectlVersion gets the version of kubectl
func GetKubectlVersion() (string, error) {
	cmd := exec.Command("kubectl", "version", "--client", "--short")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get kubectl version: %w", err)
	}
	// Extract version from output like: "Client Version: v1.25.0"
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Client Version:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", fmt.Errorf("failed to parse kubectl version")
}

// GetHelmVersion gets the version of Helm
func GetHelmVersion() (string, error) {
	cmd := exec.Command("helm", "version", "--short")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Helm version: %w", err)
	}
	// Extract version from output like: "v3.10.0+gce66412"
	version := strings.TrimSpace(string(output))
	if strings.HasPrefix(version, "v") {
		// Remove the +gce66412 part if it exists
		if idx := strings.Index(version, "+"); idx != -1 {
			version = version[:idx]
		}
		return version, nil
	}
	return "", fmt.Errorf("failed to parse Helm version")
}

// GetVersionFromExecutable gets version information from an executable file
func GetVersionFromExecutable(path string) (string, error) {
	// Try common version flags
	versionFlags := []string{
		"--version",
		"-v",
		"-V",
		"version",
	}

	for _, flag := range versionFlags {
		cmd := exec.Command(path, flag)
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			// Try to extract version from output
			version := strings.TrimSpace(string(output))
			// Remove any "version" prefix if present
			version = strings.TrimPrefix(strings.ToLower(version), "version")
			version = strings.TrimSpace(version)
			if version != "" {
				return version, nil
			}
		}
	}

	return "", fmt.Errorf("could not determine version")
} 