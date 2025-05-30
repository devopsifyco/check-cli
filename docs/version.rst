DevOpsify Check Tool version information for components, supporting both local and remote version lookups.

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
- --cve: Include CVE information for the specified version (optional)
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

Example output (history, table format):

.. code-block:: text

   Version      Release Date  Active Support  Security Support  EOL Date   
   ------------ ------------ -------------- --------------- ------------
   1.28.0      2024-06-01    2025-06-01     2025-12-01       2026-06-01  
   1.27.0      2023-12-01    2024-12-01     2025-06-01       2025-12-01  
   1.26.0      2023-06-01    2024-06-01     2024-12-01       2025-06-01  
   ...         ...           ...            ...              ...         

Example output (full format):

.. code-block:: json

   {
     "name": "nginx",
     "version": "1.24.0",
     "vendor": "NGINX",
     "release_date": "2023-04-11",
     "active_support_end": "2024-04-11",
     "security_support_end": "2024-10-11",
     "eol_date": "2025-04-11",
     "id": 123,
     "created_at": "2023-04-11"
   }

Example output (with CVEs):

.. code-block:: json

   {
     "version": {
       "name": "postgresql",
       "version": "16.4",
       "vendor": "PostgreSQL",
       "release_date": "2024-01-01",
       "eol_date": "2026-01-01"
     },
     "cves": [
       {
         "cve_id": "CVE-2024-12345",
         "state": "published",
         "published_date": "2024-02-01",
         "score": 8.2,
         "title": "Example vulnerability in postgresql",
         "references": ["https://example.com/cve/CVE-2024-12345"]
       }
     ]
   }

Example output (table, text format with --cve):

.. code-block:: text

   Name: postgresql
   Version: 16.4
   EOL Date: 2026-01-01

   CVEs:
   ------------------  ----------  ------  --------------------------------------------------
   CVE ID              Published   Score   Title
   ------------------  ----------  ------  --------------------------------------------------
   CVE-2025-12345      2024-02-01  8.2     Example vulnerability in postgresql
   ------------------  ----------  ------  -------------------------------------------------- 