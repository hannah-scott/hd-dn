#!/bin/sh
cd ~/hd-dn
go build .
mkdir -p ~/bin
mv hd-dn ~/bin
chmod +x ~/bin/hd-dn