Deps Command
============

Analyzes project dependencies from various package manager files.

Usage
-----
.. code-block:: bash

   check deps [path] [-o format] [--cve]

Options
-------
- path: Path to dependency file or directory (optional, defaults to current directory)
- -o, --output: Output format (optional: json, yaml)
- --cve: Include CVE information for each dependency (optional)

Supported Package Managers and Files
-----------------------------------

+----------------+-----------------------------------------------+---------------------------------------------+
| Package Manager| File(s)                                       | Features                                    |
+================+===============================================+=============================================+
| Maven          | pom.xml                                       | Properties, Parent POM, Version Resolution  |
| Node.js        | package.json, package-lock.json                | Dev Dependencies, Locked Versions           |
| Python         | requirements.txt, pyproject.toml               | Version Specifiers, URL/Optional Deps       |
| .NET           | project.json, .csproj, packages.config         | Package/Dev References, Framework           |
| Go             | go.mod                                        | Module/Indirect Dependencies                |
+----------------+-----------------------------------------------+---------------------------------------------+

Examples
--------
.. code-block:: bash

   check deps
   check deps pom.xml -o json
   check deps requirements.txt -o json
   check deps pyproject.toml -o yaml
   check deps package.json -o json
   check deps project.csproj -o yaml
   check deps ./checks/dependencies/samples/pom.xml --cve -o json

Sample output with CVEs (JSON):

.. code-block:: json

   [
     {
       "name": "express",
       "version": "4.18.2",
       "manager": "npm",
       "cves": [
         {
           "cve_id": "CVE-2023-12345",
           "state": "published",
           "published_date": "2023-01-01",
           "score": 7.5,
           "title": "Example vulnerability in express",
           "references": ["https://example.com/cve/CVE-2023-12345"]
         }
       ]
     },
     {
       "name": "lodash",
       "version": "4.17.21",
       "manager": "npm"
     }
   ]

Sample output (table, text format with --cve):

.. code-block:: text

   - express (4.18.2) [npm]
     ------------------  ----------  ------  ----------------------------------------------
     CVE ID              Published   Score   Title
     ------------------  ----------  ------  ----------------------------------------------
     CVE-2025-12345      2025-01-01  7.5     Example vulnerability in express
     ------------------  ----------  ------  ----------------------------------------------
   - lodash (4.17.21) [npm]
     No CVE 