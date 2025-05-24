# Updated to support 'Edit on GitHub' links using html_context
import os
import sys
sys.path.insert(0, os.path.abspath('..'))

project = 'Check Tool'
copyright = ''
author = ''
release = ''

extensions = []
templates_path = ['_templates']
exclude_patterns = []

html_theme = 'sphinx_rtd_theme'
html_static_path = ['_static']

# Update the following html_context values with your actual GitHub username and repository name
html_context = {
    "display_github": True,  # Integrate GitHub
    "github_user": "devopsifyco",  # GitHub username or org
    "github_repo": "check-cli",  # Repository name
    "github_version": "master",  # Branch or tag
    "doc_path": "docs",  # Path to your docs root in the repo
} 