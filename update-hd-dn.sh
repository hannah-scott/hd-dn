#!/bin/bash
# Script to autoupdate website based on git repo status
cd /home/hd-dn/hd-dn
git fetch origin
if git status | grep -q behind; then
  git merge origin/main
  ./build.sh
  systemctl --user restart hd-dn.service
fi