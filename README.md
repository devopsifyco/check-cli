# Check CLI Tool

A versatile command-line tool for system checks and dependency management.

## Features

- Dependency analysis for multiple package managers
- Version information checks with support for local and remote version lookups
- SSL certificate validation
- Operating system information
- Network speed testing
- Support for JSON and YAML output formats
- Comprehensive version history tracking
- Support for multiple package managers and dependency file formats

## Installation

Download the latest release for your platform from the [releases page](https://github.com/devopsifyco/check-cli/releases/latest).

### Manual Download

- [Windows (check.exe)](https://github.com/devopsifyco/check-cli/releases/download/0.0.5/check.exe)
- [Linux AMD64 (check-linux-amd64)](https://github.com/devopsifyco/check-cli/releases/download/0.0.5/check-linux-amd64)
- [Linux ARM64 (check-linux-arm64)](https://github.com/devopsifyco/check-cli/releases/download/0.0.5/check-linux-arm64)
- [macOS Intel (check-macos-intel)](https://github.com/devopsifyco/check-cli/releases/download/0.0.5/check-macos-intel)
- [macOS ARM64 (check-macos-arm64)](https://github.com/devopsifyco/check-cli/releases/download/0.0.5/check-macos-arm64)

### Download via Command Line

**Linux/macOS:**
```sh
curl -LO https://github.com/devopsifyco/check-cli/releases/download//check-linux-amd64
chmod +x check-linux-amd64
./check-linux-amd64 --help
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://github.com/devopsifyco/check-cli/releases/download//check.exe" -OutFile "check.exe"
.\check.exe --help
```
```

## Building from Source

You can build the CLI tool locally for all supported platforms using Docker. This will generate executables for Windows, Linux (amd64/arm64), and macOS (Intel/ARM).

### On Windows

Run:
```bat
build\windows.bat
```

### On Linux or macOS

Run:
```bash
bash build/linux.sh
```

or, after making it executable:

```bash
chmod +x build/linux.sh
./build/linux.sh
```

The built executables will be placed in the `dist` directory:
- check.exe (Windows)
- check-linux-amd64 (Linux AMD64)
- check-linux-arm64 (Linux ARM64)
- check-macos-intel (macOS Intel)
- check-macos-arm64 (macOS ARM64)

## Commands

### `deps` - Dependency Check

Analyzes project dependencies from various package manager files.

#### Usage
```bash
check deps [path] [-o format]
```

#### Options
- `path`: Path to dependency file or directory (optional, defaults to current directory)
- `-o, --output`: Output format (optional: json, yaml)

#### Supported Package Managers and Files

| Package Manager | File(s)                    | Features                    |
|----------------|----------------------------|----------------------------|
| Maven          | `pom.xml`                  | Properties, Parent POM, Version Resolution |
| Node.js        | `package.json`, `package-lock.json` | Dev Dependencies, Locked Versions |
| Python         | `requirements.txt`         | Version Specifiers, URL Dependencies |
|                | `pyproject.toml`           | Optional Dependencies     |
| .NET           | `project.json`, `.csproj`, `packages.config`  | Package References, Target Framework, Development Dependencies |
| Go             | `go.mod`                   | Module Dependencies, Indirect Dependencies |

#### Examples

Check dependencies in current directory:
```bash
check deps
```

Check specific files with JSON/YAML output:
```bash
# Maven dependencies
check deps pom.xml -o json
check deps pom.xml -o yaml

# Python dependencies (both formats supported)
check deps requirements.txt -o json
check deps pyproject.toml -o yaml

# Node.js dependencies
check deps package.json -o json

# .NET dependencies
check deps project.csproj -o yaml
```

# Example: Check dependencies and include CVE information
check deps ./checks/dependencies/samples/pom.xml --cve -o json
```

Sample output with CVEs:
```json
{
  "dependencies": [
    {
      "name": "express",
      "version": "4.18.2",
      "manager": "npm",
      "cves": [
        {
          "id": "CVE-2023-12345",
          "severity": "high",
          "description": "Example vulnerability in express"
        }
      ]
    }
  ]
}
```

Output formats:
```json
{
  "dependencies": [
    {
      "name": "org.springframework.boot:spring-boot-starter-web",
      "version": "3.2.0",
      "manager": "maven"
    },
    {
      "name": "pytest",
      "version": "7.4.3",
      "manager": "pip",
      "tags": ["optional:test"]
    }
  ]
}
```

```yaml
dependencies:
  - name: express
    version: 4.18.2
    manager: npm
  - name: Django
    version: 4.2.7
    manager: pip
```

### `version` - Version Check

Check version information for components, supporting both local and remote version lookups.

```bash
check version [component] [version] [flags]
```

#### Arguments
- `component`: The component to check (e.g., nginx, istio)
- `version`: (Optional) Specific version to check. If omitted, returns the latest version.

#### Flags
- `--apikey`: API key for authentication with remote version service
- `--full`: Show full version information including release dates and support timeline
- `--history`: Show version history
- `--client`: Check local client version instead of remote API
- `-o, --output`: Output format (json, yaml)

#### Supported Components for Local Version Check
- CLI tool itself
- Docker
- Kubectl
- Helm

#### Examples

Check latest version of a component:
```bash
check version nginx
check version nginx --full
check version nginx --history
```

Check specific version of a component:
```bash
check version istio 1.23
check version nginx 1.24.0 --full
```

Check cves of specific version of product:
```bash
check version postgresql 16.4 --cve
```

Check local component version:
```bash
check version docker --client
check version kubectl --client
```

Example output (full format):
```
Name: nginx
Version: 1.24.0
Vendor: NGINX
Release Date: 2023-04-11
Active Support End: 2024-04-11
Security Support End: 2024-10-11
EOL Date: 2025-04-11
ID: 123
Created At: 2023-04-11
```

Example output (specific version):
```
Name: istio
Version: 1.23
EOL Date: 2025-04-16
```

### `ssl` - SSL Certificate Check

Validate SSL certificates for domains.

#### Usage
```bash
check ssl example.com
```

#### Example Output
```
Domain: example.com
Issuer: Let's Encrypt
Valid From: 2024-01-01
Valid To: 2024-04-01
Status: valid
```

### `os` - Operating System Check

Display system information.

#### Usage
```bash
check os
```

#### Example Output
```
OS: Windows 10
Architecture: amd64
Kernel Version: 10.0.22631
```

### `speed` - Network Speed Test

Perform network speed test.

#### Usage
```bash
check speed
```

#### Example Output
```
Download: 100 Mbps
Upload: 20 Mbps
Latency: 15 ms
```

## Global Flags

- `-o, --output`: Output format (json, yaml)
- `--apikey`: API key for authentication

## Exit Codes

- 0: Success
- 1: Error
- 2: Invalid arguments

## Project Structure

The main directories and files in this project are:

- `main.go` - Entry point for the CLI tool
- `cmd/` - Command definitions and CLI logic
- `checks/` - Implementation of system checks (dependencies, version, ssl, os, speed, etc.)
  - `dependencies/` - Dependency checkers for various package managers
  - `version/` - Version check logic
  - `os/` - OS information logic
- `build/` - Build scripts for different platforms
- `assets/` - Icons and other resources
- `Dockerfile` - For building the CLI tool as a Docker image
- `.github/` - GitHub Actions workflows and templates

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

## License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file or visit [https://devopsify.co/license](https://devopsify.co/license) for details.

## Support

For support, please open an issue in the GitHub repository or contact the maintainers.

## Acknowledgments

- [Ookla Speedtest](https://www.speedtest.net/) for the speed test functionality
- [Go](https://golang.org/) for the programming language
- [Docker](https://www.docker.com/) for containerization support 
