# Updated to support 'Edit on GitHub' links using html_context
import os
import sys
sys.path.insert(0, os.path.abspath('..'))

project = 'Check Tool'
copyright = 'DevOpsify'
author = 'DevOpsify'
release = ''

extensions = ['sphinx_multiversion']
templates_path = ['_templates']
exclude_patterns = []

html_theme = 'sphinx_rtd_theme'
html_static_path = ['_static']
pygments_style = 'sphinx'

# Update the following html_context values with your actual GitHub username and repository name
html_theme_options = {
    "display_github": True,
    "github_user": "devopsifyco",
    "github_repo": "check-cli",
    "github_version": "master",
    "conf_py_path": "docs",  # Path in the repo to the docs folder
}

html_context = {
    "display_github": True, # Integrate GitHub
    "github_user": "devopsifyco", # Username
    "github_repo": "check-cli", # Repo name
    "github_version": "master", # Branch
    "conf_py_path": "/docs/", # Path in the checkout to the docs root
} 

# sphinx-multiversion configuration (optional, can be removed if not needed)
smv_tag_whitelist = r'^\d+\.\d+$'  # Only tags like 1.0, 2.1, etc.
smv_branch_whitelist = r'^(main|master|develop)$'  # Only main branches
smv_remote_whitelist = r'^origin$'
smv_released_pattern = r'^tags/\d+\.\d+$'
smv_outputdir_format = '{ref.name}'
