# DevOpsify Check Tool

DevOpsify Check Tool is a fast, cross-platform command-line utility for developers, DevOps engineers, and IT professionals to quickly analyze dependencies, check software versions (with CVE info), validate SSL certificates, inspect OS details, and test network speed. It supports Windows, Linux, and macOS, and outputs results in JSON or YAML for easy integration into CI/CD, security audits, and troubleshooting workflows.

## ✨Features

- 📦 Analyze project dependencies for multiple package managers (Maven, npm, pip, Go, .NET, etc.)
- 🔍 Check and compare software versions locally and remotely, with version history and CVE vulnerability info
- 🔒 Validate SSL certificates for domains
- 🖥️ Display detailed operating system and environment information
- 🚀 Perform network speed tests
- 📝 Output results in JSON and YAML formats for automation and reporting

For full documentation, visit the [official DevOpsify Check Tool documentation](https://devopsifyco.github.io/check-cli).

## Installation

Download the latest release for your platform from the [releases page](https://github.com/devopsifyco/check-cli/releases/latest).

### Manual Download

- [Windows (check.exe)](https://github.com/devopsifyco/check-cli/releases/download/0.0.2/check.exe)
- [Linux AMD64 (check-linux-amd64)](https://github.com/devopsifyco/check-cli/releases/download/0.0.2/check-linux-amd64)
- [Linux ARM64 (check-linux-arm64)](https://github.com/devopsifyco/check-cli/releases/download/0.0.2/check-linux-arm64)
- [macOS Intel (check-macos-intel)](https://github.com/devopsifyco/check-cli/releases/download/0.0.2/check-macos-intel)
- [macOS ARM64 (check-macos-arm64)](https://github.com/devopsifyco/check-cli/releases/download/0.0.2/check-macos-arm64)

### Download via Command Line

**Linux/macOS:**
```sh
curl -LO https://github.com/devopsifyco/check-cli/releases/download/0.0.2/check-linux-amd64
chmod +x check-linux-amd64
./check-linux-amd64 --help
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://github.com/devopsifyco/check-cli/releases/download/0.0.2/check.exe" -OutFile "check.exe"
.\check.exe --help
```

## 🚀 Get started

**Clone the repository:**
   ```sh
   git clone https://github.com/devopsifyco/check-cli.git
   cd check-cli
   ```

### Project Structure

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
