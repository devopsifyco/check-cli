Check OS
==========

Display system information.

Usage
-----

.. code-block:: bash

   check os [-f format]

Options
-------
- -f, --format: Output format (optional, default: json)

Example Output (text)
---------------------

.. code-block:: text

   OS: Windows 10
   Architecture: amd64
   CPUs: 8
   Go Version: go1.21.0
   Uptime: 2h30m

   Installed Software:
   Name: Google Chrome
     Version: 114.0.5735.199
     Publisher: Google LLC
     Install Date: 2023-01-01

   CPU Information:
   CPU 0:
     Model: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
     Cores: 8
     Frequency: 3600.00 MHz
     Core 0 Usage: 10.00%
     ...

   Memory Information:
   Total: 32.00 GB
   Used: 12.00 GB
   Usage: 37.50%

   Disk Information:
   Path: C:\
   Total: 512.00 GB
   Used: 200.00 GB
   Usage: 39.06%

   Network Information:
   ...

Example Output (JSON)
---------------------

.. code-block:: json

   {
     "os": "Windows 10",
     "arch": "amd64",
     "cpus": 8,
     "go_version": "go1.21.0",
     "uptime": "2h30m",
     "software_info": [
       {
         "name": "Google Chrome",
         "version": "114.0.5735.199",
         "publisher": "Google LLC",
         "install_date": "2023-01-01"
       }
     ],
     "cpu_info": [
       {
         "model_name": "Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz",
         "cores": 8,
         "mhz": 3600.00,
         "usage": [10.0, 12.0, 8.0, 9.0, 11.0, 7.0, 13.0, 10.0]
       }
     ],
     "memory_total": 34359738368,
     "memory_used": 12884901888,
     "memory_percent": 37.5,
     "disk_info": [
       {
         "path": "C:\\",
         "total": 549755813888,
         "used": 214748364800,
         "used_percent": 39.06
       }
     ],
     "network_info": [
       {
         "name": "Ethernet",
         "mac": "00:1A:2B:3C:4D:5E",
         "ip": "192.168.1.100"
       }
     ]
   } 