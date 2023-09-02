#!/bin/sh
./build_tgt.sh
git add .
git commit -m "TGT $1"
git push
