# DevOpsify Check Tool

DevOpsify Check Tool is a fast, cross-platform command-line utility for developers, DevOps engineers, and IT professionals to quickly analyze dependencies, check software versions (with CVE info), validate SSL certificates, inspect OS details, test network speed, and count lines of code. It supports Windows, Linux, and macOS, and outputs results in JSON or YAML for easy integration into CI/CD, security audits, and troubleshooting workflows.

## ✨Features

- 📦 Analyze project dependencies for multiple package managers (Maven, npm, pip, Go, .NET, etc.)
- 🔍 Check and compare software versions locally and remotely, with version history and **CVE vulnerability info**
- 🔒 Validate SSL certificates for domains
- 🖥️ Display detailed operating system and environment information
- 🚀 Perform network speed tests
- 📊 Count lines of code in a directory or file (LOC)
- 📝 Output results in JSON and YAML formats for automation and reporting
- 📜 Version history tracking
- 🗂️ Support for multiple dependency file formats

For full documentation, visit the [official DevOpsify Check Tool documentation](https://devopsifyco.github.io/check-cli).

## 🚀 Quick Usage Example

```sh
# Check the latest version of nginx
./check version nginx

# Check a specific version and show full details
./check version nginx 1.24.0 --full

# Check a specific version and include CVE information
./check version postgresql 16.4 --cve

# Check dependencies in a project and include CVE information
./check code deps --cve -o json

# Check dependencies for a specific file
./check code deps pom.xml -o json

# Count lines of code in the current directory
./check code loc

# Count lines of code in a specific directory or file
./check code loc ./checks/dependencies/samples
./check code loc ./checks/dependencies/samples -o json

# Validate SSL certificate for a domain
./check ssl example.com

# Show operating system details
./check os

# Run a network speed test
./check speed
```

## Installation

Download the latest release for your platform from the [releases page](https://github.com/devopsifyco/check-cli/releases/latest).

### Manual Download

- [Windows (check.exe)](https://github.com/devopsifyco/check-cli/releases/download/0.0.14/check.exe)
- [Linux AMD64 (check-linux-amd64)](https://github.com/devopsifyco/check-cli/releases/download/0.0.14/check-linux-amd64)
- [Linux ARM64 (check-linux-arm64)](https://github.com/devopsifyco/check-cli/releases/download/0.0.14/check-linux-arm64)
- [macOS Intel (check-macos-intel)](https://github.com/devopsifyco/check-cli/releases/download/0.0.14/check-macos-intel)
- [macOS ARM64 (check-macos-arm64)](https://github.com/devopsifyco/check-cli/releases/download/0.0.14/check-macos-arm64)

### Download via Command Line

**Linux/macOS:**
```sh
curl -LO https://github.com/devopsifyco/check-cli/releases/download/0.0.14/check-linux-amd64
chmod +x check-linux-amd64
./check-linux-amd64 --help
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://github.com/devopsifyco/check-cli/releases/download/0.0.14/check.exe" -OutFile "check.exe"
.\check.exe --help
```

## 🏗️ Project Structure

- `main.go` - Entry point for the CLI tool
- `cmd/` - Command definitions and CLI logic
- `checks/` - Implementation of system checks (dependencies, version, ssl, os, speed, etc.)
  - `dependencies/` - Dependency checkers for various package managers (Maven, Node.js, Python, .NET, Go)
  - `version/` - Version check logic (local and remote)
  - `os/` - OS information logic (Windows, Linux, macOS)
- `build/` - Build scripts for different platforms
- `assets/` - Icons and other resources
- `docs/` - Documentation sources (reStructuredText, Sphinx config, etc.)
- `Dockerfile` - For building the CLI tool as a Docker image
- `.github/` - GitHub Actions workflows and templates
- `CONTRIBUTING.md` - Contribution guidelines
- `LICENSE` - License file

**Build using Docker (recommended for all platforms):**

On **Windows**:

     ```bat
     build\windows.bat
     ```
On **Linux/macOS**:

     ```sh
     bash build/linux.sh
     # or, after making it executable:
     chmod +x build/linux.sh
     ./build/linux.sh
     ```

   The built executables will be placed in the `dist` directory:
   - `check.exe` (Windows)
   - `check-linux-amd64` (Linux AMD64)
   - `check-linux-arm64` (Linux ARM64)
   - `check-macos-intel` (macOS Intel)
   - `check-macos-arm64` (macOS ARM64)

**Run the CLI:**
   - On your platform, use the appropriate executable from the `dist` directory. For example:
     ```sh
     ./dist/check-linux-amd64 --help
     ./dist/check-macos-intel --help
     ./dist/check.exe --help
     ```

For more details, see the [official documentation](https://devopsifyco.github.io/check-cli).


## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

## 🎫 License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file or visit [https://devopsify.co/license](https://devopsify.co/license) for details.

## Support

For support, please open an issue in the GitHub repository or contact the maintainers.

## Acknowledgments

- [Ookla Speedtest](https://www.speedtest.net/) for the speed test functionality
- [Go](https://golang.org/) for the programming language
- [Docker](https://www.docker.com/) for containerization support
