package checks

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"gopkg.in/yaml.v3"
)

// SpeedResult implements CheckResult interface for speed checks
// (Struct fields unchanged from SpeedtestResult)
type SpeedResult struct {
	Download struct {
		Bandwidth float64 `json:"bandwidth" yaml:"bandwidth"`
		Bytes    float64 `json:"bytes" yaml:"bytes"`
		Elapsed  float64 `json:"elapsed" yaml:"elapsed"`
		Latency  struct {
			IQM    float64 `json:"iqm" yaml:"iqm"`
			Low    float64 `json:"low" yaml:"low"`
			High   float64 `json:"high" yaml:"high"`
			Jitter float64 `json:"jitter" yaml:"jitter"`
		} `json:"latency" yaml:"latency"`
	} `json:"download" yaml:"download"`
	Upload struct {
		Bandwidth float64 `json:"bandwidth" yaml:"bandwidth"`
		Bytes    float64 `json:"bytes" yaml:"bytes"`
		Elapsed  float64 `json:"elapsed" yaml:"elapsed"`
		Latency  struct {
			IQM    float64 `json:"iqm" yaml:"iqm"`
			Low    float64 `json:"low" yaml:"low"`
			High   float64 `json:"high" yaml:"high"`
			Jitter float64 `json:"jitter" yaml:"jitter"`
		} `json:"latency" yaml:"latency"`
	} `json:"upload" yaml:"upload"`
	Ping struct {
		Jitter  float64 `json:"jitter" yaml:"jitter"`
		Latency float64 `json:"latency" yaml:"latency"`
		Low     float64 `json:"low" yaml:"low"`
		High    float64 `json:"high" yaml:"high"`
	} `json:"ping" yaml:"ping"`
	Server struct {
		Name     string `json:"name" yaml:"name"`
		Country  string `json:"country" yaml:"country"`
		Sponsor  string `json:"sponsor" yaml:"sponsor"`
		ID       int    `json:"id" yaml:"id"`
		Host     string `json:"host" yaml:"host"`
		Port     int    `json:"port" yaml:"port"`
		Location string `json:"location" yaml:"location"`
		IP       string `json:"ip" yaml:"ip"`
	} `json:"server" yaml:"server"`
	PacketLoss float64 `json:"packetLoss" yaml:"packetLoss"`
	ISP        string  `json:"isp" yaml:"isp"`
}

// Print implements CheckResult interface
func (r *SpeedResult) Print(outputFormat string) {
	switch outputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(r)
		if err != nil {
			fmt.Printf("Error formatting YAML: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
	default:
		// Convert bandwidth from bytes per second to megabits per second
		downloadMbps := r.Download.Bandwidth * 8 / 1000000
		uploadMbps := r.Upload.Bandwidth * 8 / 1000000

		fmt.Printf("Download: %.2f Mbps\n", downloadMbps)
		fmt.Printf("Upload: %.2f Mbps\n", uploadMbps)
		fmt.Printf("Ping: %.2f ms\n", r.Ping.Latency)
		fmt.Printf("Server: %s (%s) - %s\n", r.Server.Name, r.Server.Country, r.Server.Sponsor)
	}
}

// SpeedCheckCommand implements the CheckCommand interface for speed checks
type SpeedCheckCommand struct {
	*BaseCheckCommand
	showServer bool
}

// NewSpeedCheckCommand creates a new speed check command
func NewSpeedCheckCommand() *SpeedCheckCommand {
	return &SpeedCheckCommand{
		BaseCheckCommand: NewBaseCheckCommand(
			"speed",
			"Run a network speed test",
			"speed",
			0,
		),
	}
}

// Execute implements the CheckCommand interface
func (c *SpeedCheckCommand) Execute(args []string) (CheckResult, error) {
	var result SpeedResult
	var err error

	if runtime.GOOS == "windows" {
		result, err = runSpeedWindows()
	} else {
		result, err = runSpeedUnix()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to run speed: %v", err)
	}

	return &result, nil
}

// CheckSpeed runs a speed test and returns the results
func CheckSpeed(jsonOutput bool) {
	var result SpeedResult
	var err error

	if runtime.GOOS == "windows" {
		result, err = runSpeedWindows()
	} else {
		result, err = runSpeedUnix()
	}

	if err != nil {
		fmt.Printf("Error running speed: %v\n", err)
		os.Exit(1)
	}

	if jsonOutput {
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		// Convert bandwidth from bytes per second to megabits per second
		downloadMbps := result.Download.Bandwidth * 8 / 1000000
		uploadMbps := result.Upload.Bandwidth * 8 / 1000000

		fmt.Printf("Speed Results:\n")
		fmt.Printf("Download: %.2f Mbps\n", downloadMbps)
		fmt.Printf("Upload: %.2f Mbps\n", uploadMbps)
		fmt.Printf("Ping: %.2f ms\n", result.Ping.Latency)
		fmt.Printf("Server: %s (%s) - %s\n", result.Server.Name, result.Server.Country, result.Server.Sponsor)
	}
}

// runSpeedWindows runs speed test on Windows
func runSpeedWindows() (SpeedResult, error) {
	// Get the directory of the current executable
	checkExePath, err := os.Executable()
	if err != nil {
		return SpeedResult{}, fmt.Errorf("failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(checkExePath)

	// First try to find speedtest.exe in the speedtest directory
	speedtestPath := filepath.Join(exeDir, "speedtest", "speedtest.exe")
	if _, err := os.Stat(speedtestPath); os.IsNotExist(err) {
		// If not found, try to download and install
		fmt.Println("Speedtest CLI not found. Attempting to download...")
		err = downloadSpeedtestWindows()
		if err != nil {
			return SpeedResult{}, fmt.Errorf("failed to download speedtest: %v", err)
		}
	}

	// Run speedtest with JSON output
	cmd := exec.Command(speedtestPath, "--accept-license=true", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return SpeedResult{}, fmt.Errorf("failed to run speedtest: %v\nstderr: %s", err, string(exitErr.Stderr))
		}
		return SpeedResult{}, fmt.Errorf("failed to run speedtest: %v", err)
	}

	// Parse JSON output
	var result SpeedResult
	if err := json.Unmarshal(output, &result); err != nil {
		return SpeedResult{}, fmt.Errorf("failed to parse speedtest output: %v\nraw output: %s", err, string(output))
	}

	return result, nil
}

// runSpeedUnix runs speed test on Unix-like systems
func runSpeedUnix() (SpeedResult, error) {
	// First try to find speedtest in system PATH
	speedtestPath, err := exec.LookPath("speedtest")
	if err != nil {
		// If not found in PATH, try to download and install
		fmt.Println("Speedtest CLI not found in PATH. Attempting to download...")
		// Get the directory of the current executable
		checkExePath, err := os.Executable()
		if err != nil {
			return SpeedResult{}, fmt.Errorf("failed to get executable path: %v", err)
		}
		exeDir := filepath.Dir(checkExePath)
		speedtestPath = filepath.Join(exeDir, "speedtest", "speedtest")
		// Check if we already have a downloaded version
		if _, err := os.Stat(speedtestPath); os.IsNotExist(err) {
			err = downloadSpeedtestUnix()
			if err != nil {
				return SpeedResult{}, fmt.Errorf("failed to download speedtest: %v", err)
			}
		}
	}

	// Try to run speedtest command
	cmd := exec.Command(speedtestPath, "--accept-license=true", "--format=json")
	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return SpeedResult{}, fmt.Errorf("failed to run speedtest: %v\nstderr: %s", err, stderr.String())
	}

	output := stdout.Bytes()
	if len(output) == 0 {
		return SpeedResult{}, fmt.Errorf("speedtest produced no output\nstderr: %s", stderr.String())
	}

	// Parse JSON output
	var result SpeedResult
	if err := json.Unmarshal(output, &result); err != nil {
		return SpeedResult{}, fmt.Errorf("failed to parse speedtest output: %v\noutput: %s", err, string(output))
	}

	return result, nil
}

// downloadSpeedtestWindows downloads and installs Speedtest CLI for Windows
func downloadSpeedtestWindows() error {
	// Get the directory of the current executable
	checkExePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(checkExePath)

	// Create speedtest directory in the same folder as the executable
	speedtestDir := filepath.Join(exeDir, "speedtest")
	if err := os.MkdirAll(speedtestDir, 0755); err != nil {
		return fmt.Errorf("failed to create speedtest directory: %v", err)
	}

	// Download the installer
	zipPath := filepath.Join(speedtestDir, "speedtest.zip")
	// Use the official Ookla Speedtest CLI download link
	resp, err := http.Get("https://install.speedtest.net/app/cli/ookla-speedtest-1.2.0-win64.zip")
	if err != nil {
		return fmt.Errorf("failed to download speedtest: %v", err)
	}
	defer resp.Body.Close()

	// Check if the download was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download speedtest: HTTP status %d", resp.StatusCode)
	}

	// Create the zip file
	out, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write zip file: %v", err)
	}

	// Close the file before trying to read it
	out.Close()

	// Open the zip file
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer zipReader.Close()

	// Extract the executable from the zip
	extractDir := filepath.Join(speedtestDir, "extract")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("failed to create extract directory: %v", err)
	}

	// Find and extract the speedtest.exe file
	var foundExePath string
	for _, file := range zipReader.File {
		if strings.HasSuffix(strings.ToLower(file.Name), "speedtest.exe") {
			// Create the file
			outputPath := filepath.Join(extractDir, filepath.Base(file.Name))
			outputFile, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("failed to create output file: %v", err)
			}
			defer outputFile.Close()

			// Open the file in the zip
			zipFile, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to open file in zip: %v", err)
			}
			defer zipFile.Close()

			// Copy the contents
			_, err = io.Copy(outputFile, zipFile)
			if err != nil {
				return fmt.Errorf("failed to extract file: %v", err)
			}

			foundExePath = outputPath
			break
		}
	}

	if foundExePath == "" {
		// List files in the zip for debugging
		var fileList []string
		for _, file := range zipReader.File {
			fileList = append(fileList, file.Name)
		}
		return fmt.Errorf("executable not found in zip file. Found files: %v", fileList)
	}

	// Clean up the zip file
	os.Remove(zipPath)

	// Make the file executable
	if err := os.Chmod(foundExePath, 0755); err != nil {
		return fmt.Errorf("failed to make file executable: %v", err)
	}

	// Copy the executable to the speedtest directory
	finalPath := filepath.Join(speedtestDir, "speedtest.exe")
	if err := copyFile(foundExePath, finalPath); err != nil {
		return fmt.Errorf("failed to copy executable: %v", err)
	}

	// Clean up the extract directory
	os.RemoveAll(extractDir)

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// downloadSpeedtestUnix downloads and installs Speedtest CLI for Unix-like systems
func downloadSpeedtestUnix() error {
	// First try to install without sudo
	var cmd *exec.Cmd
	var err error

	// Try to install speedtest-cli based on the package manager
	switch {
	case CommandExists("apt-get"):
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "speedtest-cli")
	case CommandExists("yum"):
		cmd = exec.Command("sudo", "yum", "install", "-y", "speedtest-cli")
	case CommandExists("brew"):
		cmd = exec.Command("brew", "install", "speedtest-cli")
	default:
		return fmt.Errorf("no supported package manager found")
	}

	if err = cmd.Run(); err == nil {
		return nil
	}

	// If package manager installation failed, download the official binary
	fmt.Println("Package manager installation failed. Downloading official Speedtest CLI...")
	// Create a temporary directory for the speedtest binary
	tempDir, err := os.MkdirTemp("", "speedtest")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Determine architecture
	var arch string
	switch runtime.GOARCH {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		arch = "aarch64"
	default:
		return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	// Download the official binary
	downloadURL := fmt.Sprintf("https://install.speedtest.net/app/cli/ookla-speedtest-1.2.0-linux-%s.tgz", arch)
	tarPath := filepath.Join(tempDir, "speedtest.tgz")
	// Download the file
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download speedtest: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download speedtest: HTTP status %d", resp.StatusCode)
	}

	// Create the tar.gz file
	out, err := os.Create(tarPath)
	if err != nil {
		return fmt.Errorf("failed to create tar file: %v", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write tar file: %v", err)
	}
	out.Close()

	// Extract the tar.gz file
	cmd = exec.Command("tar", "-xzf", tarPath, "-C", tempDir)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract tar file: %v", err)
	}

	// Find the extracted binary
	var binaryPath string
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temporary directory: %v", err)
	}
	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Name()), "speedtest") {
			binaryPath = filepath.Join(tempDir, file.Name())
			break
		}
	}

	if binaryPath == "" {
		return fmt.Errorf("failed to find speedtest binary in downloaded package")
	}

	// Get the directory of the current executable
	checkExePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(checkExePath)

	// Create speedtest directory in the same folder as the executable
	speedtestDir := filepath.Join(exeDir, "speedtest")
	if err := os.MkdirAll(speedtestDir, 0755); err != nil {
		return fmt.Errorf("failed to create speedtest directory: %v", err)
	}

	// Copy the binary to the speedtest directory
	finalPath := filepath.Join(speedtestDir, "speedtest")
	if err := copyFile(binaryPath, finalPath); err != nil {
		return fmt.Errorf("failed to copy binary: %v", err)
	}

	// Make the binary executable
	if err := os.Chmod(finalPath, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %v", err)
	}

	return nil
} 