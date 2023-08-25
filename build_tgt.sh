#!/bin/sh
set -uox
cd static/three-good-things
# mv index.html .index_backup.html
cat ~/obsidian/Three\ Good\ Things.md > .tgt_1.md
cat .tgt_*.md > index.md

## Add article tags
cat _index.md > index.md
sed -E "s/^## /<\/article>\n<article>\n## /g" .tgt_1.md .tgt_2.md | tail -n +2 >> index.md
echo "</article>" >> index.md

# Convert to HTML
pandoc index.md -f markdown-smart -o index.html
rm index.md
tidy -imq index.html 2>/dev/null
./moku.py