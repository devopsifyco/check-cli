Deps Command
============

Analyzes project dependencies from various package manager files.

Usage
-----
.. code-block:: bash

   check deps [path] [-o format]

Options
-------
- path: Path to dependency file or directory (optional, defaults to current directory)
- -o, --output: Output format (optional: json, yaml)

Supported Package Managers and Files
-----------------------------------

+----------------+----------------------------+---------------------------------------------+
| Package Manager| File(s)                    | Features                                    |
+================+============================+=============================================+
| Maven          | pom.xml                    | Properties, Parent POM, Version Resolution  |
| Node.js        | package.json, package-lock.json | Dev Dependencies, Locked Versions     |
| Python         | requirements.txt, pyproject.toml | Version Specifiers, URL/Optional Deps |
| .NET           | project.json, .csproj, packages.config | Package/Dev References, Framework |
| Go             | go.mod                     | Module/Indirect Dependencies                |
+----------------+----------------------------+---------------------------------------------+

Examples
--------
.. code-block:: bash

   check deps
   check deps pom.xml -o json
   check deps requirements.txt -o json
   check deps pyproject.toml -o yaml
   check deps package.json -o json
   check deps project.csproj -o yaml

Example: Check dependencies and include CVE information
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
.. code-block:: bash

   check deps ./checks/dependencies/samples/pom.xml --cve -o json

Sample output with CVEs:
.. code-block:: json

   {
     "dependencies": [
       {
         "name": "express",
         "version": "4.18.2",
         "manager": "npm",
         "cves": [
           {
             "id": "CVE-2023-12345",
             "severity": "high",
             "description": "Example vulnerability in express"
           }
         ]
       }
     ]
   } 