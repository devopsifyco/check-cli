# Updated to support 'Edit on GitHub' links using html_context
import os
import sys
sys.path.insert(0, os.path.abspath('..'))

project = 'Check Tool'
copyright = 'DevOpsify'
author = 'DevOpsify'
release = ''

extensions = []
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
