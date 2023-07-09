#!/bin/bash
# Script to autoupdate website based on git repo status

# Set some environment variables. We need these for cron
XDG_RUNTIME_DIR=/run/user/$(id -u)
DBUS_SESSION_BUS_ADDRESS=unix:path=${XDG_RUNTIME_DIR}/bus
export DBUS_SESSION_BUS_ADDRESS XDG_RUNTIME_DIR

cd /home/hd-dn/hd-dn
git fetch origin
if git status | grep -q behind; then
  git merge origin/main
  ./build.sh
  systemctl --user restart hd-dn.service
fi