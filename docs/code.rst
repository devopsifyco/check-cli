Code Command
============

The `code` command provides code analysis utilities, including dependency analysis and counting lines of code (LOC).

Usage
-----
.. code-block:: bash

   check code [sub_command] [path] [flags]

Subcommands
-----------
- deps: Analyze project dependencies (see below)
- loc: Count lines of code in a directory or file

Global Flags
------------
- --apikey: API key for authentication with remote version service
- -o, --output: Output format (optional: json, yaml)

code deps
---------
Analyzes project dependencies from various package manager files.

.. code-block:: bash

   check code deps [path] [-o format] [--cve]

Options:
- path: Path to dependency file or directory (optional, defaults to current directory)
- -o, --output: Output format (optional: json, yaml)
- --cve: Include CVE information for each dependency (optional)

Examples:
.. code-block:: bash

   check code deps
   check code deps pom.xml -o json
   check code deps ./checks/dependencies/samples/pom.xml --cve -o json

code loc
--------
Counts lines of code in a directory or file using gocloc.

.. code-block:: bash

   check code loc [path] [-o format]

Options:
- path: Path to directory or file (optional, defaults to current directory)
- -o, --output: Output format (optional: json, yaml)

Examples:
.. code-block:: bash

   check code loc
   check code loc ./checks/dependencies/samples
   check code loc ./checks/dependencies/samples -o json

Sample output (table):

.. code-block:: text

   Language        Files    Blank    Comment    Code
   Go              10      100      50         1000
   Python           2       20      10          200
   TOTAL: Files=12 Blank=120 Comment=60 Code=1200

Sample output (JSON):

.. code-block:: json

   {
     "languages": [
       {"name": "Go", "files": 10, "blank": 100, "comment": 50, "code": 1000},
       {"name": "Python", "files": 2, "blank": 20, "comment": 10, "code": 200}
     ],
     "total": {"files": 12, "blank": 120, "comment": 60, "code": 1200}
   } 