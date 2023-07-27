#!/bin/sh
set -euox
mv index.txt .index_backup.txt
cat .tgt_*.txt > index.txt