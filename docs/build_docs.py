import os
import subprocess
import yaml

def build_doc(version, language, tag):
    os.environ["current_version"] = version
    os.environ["current_language"] = language
    # Save Makefile, conf.py, versions.yaml, and index.rst
    subprocess.run("cp Makefile /tmp/Makefile", shell=True, check=True)
    subprocess.run("cp conf.py /tmp/conf.py", shell=True, check=True)
    subprocess.run("cp versions.yaml /tmp/versions.yaml", shell=True, check=True)
    subprocess.run("cp index.rst /tmp/index.rst", shell=True, check=True)
    subprocess.run(f"git checkout {tag}", shell=True, check=True)
    # Restore Makefile, conf.py, versions.yaml, and index.rst
    subprocess.run("cp /tmp/Makefile Makefile", shell=True, check=True)
    subprocess.run("cp /tmp/conf.py conf.py", shell=True, check=True)
    subprocess.run("cp /tmp/versions.yaml versions.yaml", shell=True, check=True)
    subprocess.run("cp /tmp/index.rst index.rst", shell=True, check=True)
    os.environ['SPHINXOPTS'] = f"-D language='{language}'"
    subprocess.run("make html", shell=True, check=True)

def move_dir(src, dst):
    os.makedirs(dst, exist_ok=True)
    subprocess.run(f"cp -r {src}* {dst}", shell=True, check=True)

os.environ["build_all_docs"] = str(True)
os.environ["pages_root"] = "https://devopsifyco.github.io/check-cli"

with open("versions.yaml", "r") as yaml_file:
    docs = yaml.safe_load(yaml_file)

for version, details in docs.items():
    tag = details.get('tag', '')
    for language in details.get('languages', []):
        build_doc(version, language, tag)
        move_dir("./_build/html/", f"../pages/{version}/{language}/") 