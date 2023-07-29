#!/bin/sh
set -euox
mv index.lark .index_backup.lark
cat .tgt_*.lark > index.lark