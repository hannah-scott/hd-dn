#!/bin/sh
set -euox
cd static/three-good-things
mv index.lark .index_backup.lark
cat .tgt_*.lark > index.lark