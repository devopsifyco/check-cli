Version Command
===============

Check version information for components, supporting both local and remote version lookups.

Usage
-----
.. code-block:: bash

   check version [component] [version] [flags]

Arguments
---------
- component: The component to check (e.g., nginx, istio)
- version: (Optional) Specific version to check. If omitted, returns the latest version.

Flags
-----
- --apikey: API key for authentication with remote version service
- --full: Show full version information including release dates and support timeline
- --history: Show version history
- --client: Check local client version instead of remote API
- -o, --output: Output format (json, yaml)

Supported Components for Local Version Check
-------------------------------------------
- CLI tool itself
- Docker
- Kubectl
- Helm

Examples
--------
.. code-block:: bash

   check version nginx
   check version nginx --full
   check version nginx --history
   check version istio 1.23
   check version nginx 1.24.0 --full
   check version postgresql 16.4 --cve
   check version docker --client
   check version kubectl --client

Example output (full format):
.. code-block:: text

   Name: nginx
   Version: 1.24.0
   Vendor: NGINX
   Release Date: 2023-04-11
   Active Support End: 2024-04-11
   Security Support End: 2024-10-11
   EOL Date: 2025-04-11
   ID: 123
   Created At: 2023-04-11

Example output (specific version):
.. code-block:: text

   Name: istio
   Version: 1.23
   EOL Date: 2025-04-16 