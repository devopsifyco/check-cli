# Updated to support 'Edit on GitHub' links using html_context
import os
import sys
import yaml
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

# sphinx-multiversion configuration
smv_tag_whitelist = r'^[0-9]+\.[0-9]+$'  # Only tags like 1.0, 2.1, etc.
smv_branch_whitelist = r'^(main|master|develop)$'  # Only main branches
smv_remote_whitelist = r'^origin$'
smv_released_pattern = r'^tags/[0-9]+\.[0-9]+$'
smv_outputdir_format = '{ref.name}'

# Version/language picker context
build_all_docs = os.environ.get("build_all_docs")
pages_root = os.environ.get("pages_root", "")

if build_all_docs is not None:
    current_language = os.environ.get("current_language", "en")
    current_version = os.environ.get("current_version", "master")
    html_context = {
        'current_language': current_language,
        'languages': [],
        'current_version': current_version,
        'versions': [],
    }
    with open(os.path.join(os.path.dirname(__file__), "versions.yaml"), "r") as yaml_file:
        docs = yaml.safe_load(yaml_file)
    # Add languages for the current version
    for language in docs[current_version].get('languages', []):
        html_context['languages'].append([language, f"{pages_root}/{current_version}/{language}/"])
    # Add all versions for the current language
    for version, details in docs.items():
        if current_language in details.get('languages', []):
            html_context['versions'].append([version, f"{pages_root}/{version}/{current_language}/"])
