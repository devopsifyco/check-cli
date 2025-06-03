Check Speed
=============

Perform network speed test.

Usage
-----

.. code-block:: bash

   check speed [-o format]

Options
-------
- -o, --output: Output format (optional: json, yaml)

Example Output (text)
---------------------

.. code-block:: text

   Download: 100.00 Mbps
   Upload: 20.00 Mbps
   Ping: 15.00 ms
   Server: ExampleServer (US) - ExampleSponsor

Example Output (JSON)
---------------------

.. code-block:: json

   {
     "download": {"bandwidth": 12500000},
     "upload": {"bandwidth": 2500000},
     "ping": {"latency": 15.0},
     "server": {"name": "ExampleServer", "country": "US", "sponsor": "ExampleSponsor"}
   } 