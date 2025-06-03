Check SSL
===========

Validate SSL certificates for domains.

Usage
-----

.. code-block:: bash

   check ssl example.com [-o format]

Options
-------
- -o, --output: Output format (optional: json, yaml)

Arguments
---------
- example.com: The domain to check (required)

Example Output (text)
---------------------

.. code-block:: text

   Domain: example.com
   Valid: true
   Issuer: Let's Encrypt
   Subject: CN=example.com
   Valid From: 2024-01-01
   Valid Until: 2024-04-01
   Expires In: 90 days
   Serial Number: 1234567890
   Subject Alternative Names:
     - example.com
     - www.example.com

Example Output (JSON)
---------------------

.. code-block:: json

   {
     "domain": "example.com",
     "valid": true,
     "issuer": "Let's Encrypt",
     "subject": "CN=example.com",
     "not_before": "2024-01-01",
     "not_after": "2024-04-01",
     "expires_in_days": 90,
     "serial_number": "1234567890",
     "sans": ["example.com", "www.example.com"],
     "error": null
   } 