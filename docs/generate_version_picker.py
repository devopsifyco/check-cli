import os

HTML_TEMPLATE = """<!DOCTYPE html>
<html>
  <head>
    <title>Check Tool Documentation Versions</title>
    <meta charset=\"utf-8\">
    <style>
      body {{ font-family: sans-serif; margin: 2em; }}
      h1 {{ color: #2c3e50; }}
      ul {{ list-style: none; padding: 0; }}
      li {{ margin: 0.5em 0; }}
      a {{ color: #2980b9; text-decoration: none; font-size: 1.2em; }}
      a:hover {{ text-decoration: underline; }}
    </style>
  </head>
  <body>
    <h1>Check Tool Documentation</h1>
    <p>Select a version:</p>
    <ul>
      {links}
    </ul>
    <hr>
    <small>Built with Sphinx Multiversion</small>
  </body>
</html>
"""

def main():
    html_dir = os.path.join(os.path.dirname(__file__), "_build", "html")
    versions = []
    for name in sorted(os.listdir(html_dir)):
        path = os.path.join(html_dir, name, "index.html")
        if os.path.isdir(os.path.join(html_dir, name)) and os.path.isfile(path):
            versions.append(name)
    links = "\n      ".join(f'<li><a href="{v}/">{v}</a></li>' for v in versions)
    with open(os.path.join(html_dir, "index.html"), "w", encoding="utf-8") as f:
        f.write(HTML_TEMPLATE.format(links=links))

if __name__ == "__main__":
    main() 