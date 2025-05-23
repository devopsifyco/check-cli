const fs = require('fs');
const path = require('path');

const latestTag = process.argv[2];
if (!latestTag) {
  console.error('Usage: node update-readme-download.js <latest_tag>');
  process.exit(1);
}

const readmePath = path.join(__dirname, '..', 'README.md');
let readme = fs.readFileSync(readmePath, 'utf8');

// Find the Installation section
const installSectionRegex = /(## Installation[\s\S]*?)(?=^## |\n## |$)/m;
const match = readme.match(installSectionRegex);

if (match) {
  let installSection = match[1];
  // Replace all {version} with the new tag
  installSection = installSection.replace(/\{version\}/g, latestTag);
  // Replace the section in the README
  readme = readme.replace(installSectionRegex, installSection);
  fs.writeFileSync(readmePath, readme);
  console.log('README.md download links updated to', latestTag);
} else {
  console.error('Installation section not found in README.md');
  process.exit(1);
} 