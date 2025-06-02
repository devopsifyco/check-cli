Check MCP Server
================

A Model Context Protocol server that provides CVE checking capabilities via the Opsify API. This server enables LLMs to search for and retrieve CVE (Common Vulnerabilities and Exposures) information using a flexible set of filters.

.. warning::
   This server can access external APIs and may represent a security risk if misconfigured. Ensure your API key and usage comply with your organization's security policies.

Available Tools
---------------

- ``search_cve`` - Search CVEs with various filters via the Opsify API.

  - ``cve_id`` (string, optional): CVE ID to search for
  - ``title`` (string, optional): Title to search in CVE description
  - ``state`` (string, optional): State to filter by
  - ``priority`` (string, optional): Priority level to filter by
  - ``severity`` (string, optional): Severity level to filter by
  - ``score`` (float, optional): CVSS score to filter by
  - ``product_name`` (string, optional): Product name to filter affected products
  - ``product_version`` (string, optional): Product version to filter affected products
  - ``vendor`` (string, optional): Vendor name to filter affected products
  - ``from_date`` (string, optional): Start date for filtering (YYYY-MM-DD or ISO 8601)
  - ``to_date`` (string, optional): End date for filtering (YYYY-MM-DD or ISO 8601)
  - ``skip`` (int, optional): Number of records to skip (pagination)
  - ``limit`` (int, optional): Maximum number of records to return (pagination)

- ``search_release`` - Search releases with optional filters for vendor, product name, and date range. Supports pagination.

  - ``vendor`` (string, optional): Vendor name to filter by (case-insensitive)
  - ``product_name`` (string, optional): Product name to filter by (case-insensitive)
  - ``from_date`` (string, optional): Start date (inclusive) for filtering (YYYY-MM-DD or ISO datetime)
  - ``to_date`` (string, optional): End date (inclusive) for filtering (YYYY-MM-DD or ISO datetime)
  - ``date_field`` (string, optional): Which date field to filter on (e.g., "release_date")
  - ``page`` (int, optional): Page number (starting from 1)
  - ``page_size`` (int, optional): Number of items per page

- ``get_version_cves`` - Get CVEs for a specific version of a product. Optionally filter by vendor. Uses caching (TTL: 1 day).

  - ``product_name`` (string, required): Product name (e.g., 'nginx')
  - ``version`` (string, required): Specific version to retrieve CVEs for (e.g., '1.0.0')
  - ``vendor`` (string, optional): Vendor name to filter by (case-insensitive)

Usage
-----
.. code-block:: bash

   python -m check_mcp

Or as a script/entry point:

.. code-block:: bash

   check-mcp

Example Tool Call
-----------------
.. code-block:: json

   {
     "tool": "search_cve",
     "arguments": {
       "product_name": "nginx",
       "severity": "high",
       "from_date": "2023-01-01",
       "limit": 5
     }
   }

Installation
------------

Using uv (recommended)
^^^^^^^^^^^^^^^^^^^^^^
When using `uv <https://docs.astral.sh/uv/>`_ no specific installation is needed. We will use `uvx <https://docs.astral.sh/uv/guides/tools/>`_ to directly run *check-mcp*.

Using pip
^^^^^^^^^
Alternatively you can install ``check-mcp`` via pip:

.. code-block:: bash

   pip install check-mcp

After installation, you can run it as a script using:

.. code-block:: bash

   python -m check_mcp

Configuration
-------------

Configure for Claude.app
^^^^^^^^^^^^^^^^^^^^^^^^
Add to your Claude settings:

.. code-block:: json
   :caption: Using uvx

   "mcpServers": {
     "check": {
       "command": "uvx",
       "args": ["check-mcp"]
     }
   }

.. code-block:: json
   :caption: Using docker

   "mcpServers": {
     "check": {
       "command": "docker",
       "args": ["run", "-i", "--rm", "mcp/check"]
     }
   }

.. code-block:: json
   :caption: Using pip installation

   "mcpServers": {
     "check": {
       "command": "python",
       "args": ["-m", "check_mcp"]
     }
   }

Configure for VS Code
^^^^^^^^^^^^^^^^^^^^^
For quick installation, use one of the one-click install buttons below:

.. raw:: html

   <a href="https://insiders.vscode.dev/redirect/mcp/install?name=check&config=%7B%22command%22%3A%22uvx%22%2C%22args%22%3A%5B%22check-mcp%22%5D%7D"><img src="https://img.shields.io/badge/VS_Code-UV-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white"/></a>
   <a href="https://insiders.vscode.dev/redirect/mcp/install?name=check&config=%7B%22command%22%3A%22uvx%22%2C%22args%22%3A%5B%22check-mcp%22%5D%7D&quality=insiders"><img src="https://img.shields.io/badge/VS_Code_Insiders-UV-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white"/></a>
   <a href="https://insiders.vscode.dev/redirect/mcp/install?name=check&config=%7B%22command%22%3A%22docker%22%2C%22args%22%3A%5B%22run%22%2C%22-i%22%2C%22--rm%22%2C%22mcp%2Fcheck%22%5D%7D"><img src="https://img.shields.io/badge/VS_Code-Docker-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white"/></a>
   <a href="https://insiders.vscode.dev/redirect/mcp/install?name=check&config=%7B%22command%22%3A%22docker%22%2C%22args%22%3A%5B%22run%22%2C%22-i%22%2C%22--rm%22%2C%22mcp%2Fcheck%22%5D%7D&quality=insiders"><img src="https://img.shields.io/badge/VS_Code_Insiders-Docker-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white"/></a>

For manual installation, add the following JSON block to your User Settings (JSON) file in VS Code. You can do this by pressing ``Ctrl + Shift + P`` and typing ``Preferences: Open User Settings (JSON)``.

Optionally, you can add it to a file called ``.vscode/mcp.json`` in your workspace. This will allow you to share the configuration with others.

.. note::
   The ``mcp`` key is needed when using the ``mcp.json`` file.

.. code-block:: json
   :caption: Using uvx

   {
     "mcp": {
       "servers": {
         "check": {
           "command": "uvx",
           "args": ["check-mcp"]
         }
       }
     }
   }

.. code-block:: json
   :caption: Using Docker

   {
     "mcp": {
       "servers": {
         "check": {
           "command": "docker",
           "args": ["run", "-i", "--rm", "mcp/check"]
         }
       }
     }
   }

Debugging
---------

You can use the MCP inspector to debug the server. For uvx installations:

.. code-block:: bash

   npx @modelcontextprotocol/inspector uvx check-mcp

Or if you've installed the package in a specific directory or are developing on it:

.. code-block:: bash

   cd path/to/servers/src/check
   npx @modelcontextprotocol/inspector uv run check-mcp