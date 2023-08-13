#!/bin/sh
set -euox
cd static/three-good-things
mv index.lark .index_backup.lark
cat ~/obsidian/Three\ Good\ Things.md > .tgt_1.lark
cat .tgt_*.lark > index.lark
