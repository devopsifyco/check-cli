.. DevOpsify Check Tool documentation master file
   Welcome to DevOpsify Check Tool's documentation!

DevOpsify Check Tool
==============

A versatile command-line tool for system checks and dependency management.

Features
--------
- Dependency analysis for multiple package managers
- Version information checks (local and remote)
- SSL certificate validation
- Operating system information
- Network speed testing
- JSON and YAML output formats
- Version history tracking
- Support for multiple dependency file formats

Installation
------------
Download the latest release for your platform from the `releases page <https://github.com/devopsifyco/check-cli/releases/latest>`_.

Manual Download:
^^^^^^^^^^^^^^^^
- Windows: check.exe
- Linux AMD64: check-linux-amd64
- Linux ARM64: check-linux-arm64
- macOS Intel: check-macos-intel
- macOS ARM64: check-macos-arm64

Command Line Download:
^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: bash

   curl -LO https://github.com/devopsifyco/check-cli/releases/download/0.0.19/check-linux-amd64
   chmod +x check-linux-amd64
   ./check-linux-amd64 --help

.. code-block:: powershell

   Invoke-WebRequest -Uri "https://github.com/devopsifyco/check-cli/releases/download/0.0.19/check.exe" -OutFile "check.exe"
   .\check.exe --help

Building from Source
--------------------
You can build the CLI tool locally for all supported platforms using Docker. See the README for platform-specific instructions.

CLI
--------

.. toctree::
   :maxdepth: 2
   :caption: CLI Commands

   code
   version
   ssl
   os
   speed

MCP
--------

.. toctree::
   :maxdepth: 2
   :caption: MCP Servers

   mcp_client
   mcp_server

global_flags
exit_codes

Project Structure
-----------------
- main.go: Entry point
- cmd/: Command definitions
- checks/: System checks (dependencies, version, ssl, os, speed)
- build/: Build scripts
- assets/: Icons/resources
- Dockerfile: Docker build
- .github/: GitHub Actions workflows

Contributing & License
----------------------
See CONTRIBUTING.md and LICENSE for details. 
