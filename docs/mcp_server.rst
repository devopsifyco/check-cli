SSE Protocol
====================

Check MCP Server SSE provides a Model Context Protocol (MCP) server implementation with Server-Sent Events (SSE) support, enabling real-time streaming of results and updates to clients. This is particularly useful for LLMs and other clients that require live feedback or incremental results from long-running operations.

Overview
--------

The SSE mode allows the MCP server to push updates to connected clients over a single HTTP connection, using the Server-Sent Events protocol. This is ideal for scenarios where clients need to receive data as it becomes available, such as CVE search results, version checks, or other DevOps intelligence tasks.

Available Tools
---------------

- ``get_specific_version`` - Retrieve information for a specific version of a product.

  - ``product_name`` (string, required): Product name (e.g., 'nginx')
  - ``version`` (string, required): Specific version to retrieve (e.g., '1.0.0')
  - ``vendor`` (string, optional): Vendor name to filter by (case-insensitive)

- ``get_latest_version`` - Retrieve the latest version information for a product.

  - ``product_name`` (string, required): Product name (e.g., 'nginx')
  - ``vendor`` (string, optional): Vendor name to filter by (case-insensitive)

- ``search_releases`` - Search releases with optional filters for vendor, product name, and date range. Supports pagination.

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

Features
--------
- Real-time streaming of results using SSE
- Compatible with LLMs and automation tools that support event streams
- Secure, configurable, and extensible
- Shares the same toolset as the standard Check MCP Server (see :doc:`mcp_client`)

Usage
-----

To start the Check MCP Server in SSE mode:

.. code-block:: bash

   python -m main

Or, if you have installed the package as an entry point:

.. code-block:: bash

   check-mcp-server

The server will automatically use SSE transport if the configuration specifies ``transport = 'sse'``. By default, it listens on the host and port defined in your configuration (see below).

SSE Configuration
-----------------

To connect to the server using SSE, use the following configuration in your MCP client:

.. code-block:: json
   :caption: Example MCP client configuration

   {
     "mcpServers": {
       "devopsify": {
         "transport": "sse",
         "url": "http://localhost:8050/sse"
       }
     }
   }

.. note::
   For Windsurf users, use ``serverUrl`` instead of ``url``.

   For n8n users, use ``host.docker.internal`` instead of ``localhost`` if connecting from another container:

   ::

      http://host.docker.internal:8050/sse

Configuration
-------------

The server transport mode is controlled by the ``Config.server.transport`` setting. To enable SSE, set this value to ``sse`` in your configuration file or environment variables.

.. code-block:: ini
   :caption: Example .env or config.ini

   SERVER_TRANSPORT=sse
   SERVER_HOST=0.0.0.0
   SERVER_PORT=8050

You can also configure other parameters such as session timeout and cleanup interval in your configuration.

Starting with Docker Compose
---------------------------

You can start the MCP Server SSE service using Docker Compose. Make sure you have Docker and Docker Compose installed, and a valid ``.env`` file in the project root.

.. code-block:: bash

   docker compose up --build

This will build and start the service as defined in ``docker-compose.yml``:

.. code-block:: yaml

   services:
     mcp:
       build:
         context: .
         dockerfile: Dockerfile
       container_name: check-mcp
       ports:
         - "${PORT:-8050}:8050"
       env_file:
         - .env
       environment:
         - TRANSPORT=${TRANSPORT:-sse}
         - HOST=${HOST:-0.0.0.0}
         - PORT=${PORT:-8050}
         - OPSIFY_API_KEY=${OPSIFY_API_KEY}
         - OPSIFY_API_BASE_URL=${OPSIFY_API_BASE_URL}
         - SESSION_TIMEOUT=3600
         - SESSION_CLEANUP_INTERVAL=300
         - LOG_LEVEL=DEBUG
         - PYTHONUNBUFFERED=1
       restart: unless-stopped
       networks:
         - mcp_network

   networks:
     mcp_network:
       driver: bridge

How It Works
------------

When running in SSE mode, the server uses the ``run_sse_async()`` method of the ``FastMCP`` class. This method starts the server with SSE endpoints, allowing clients to subscribe to event streams and receive updates as they are generated.

The main entry point for the server is ``src/main.py``, which detects the transport mode and starts the server accordingly:

.. code-block:: python

   async def main():
       if Config.server.transport == 'sse':
           await mcp.run_sse_async()
       else:
           await mcp.run_async()

   if __name__ == "__main__":
       asyncio.run(main())

Clients can connect to the SSE endpoints to receive real-time updates for supported tools and operations.

See Also
--------
- :doc:`mcp_client` for available tools and usage examples
- `Server-Sent Events (MDN) <https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events>`_

License
-------

Check MCP Server SSE is licensed under the MIT License. See the LICENSE file for details. 