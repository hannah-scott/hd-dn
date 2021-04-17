#!/bin/bash

cd ~/repos/hannahs.ddns.net

# cleanup
rm *.html 
rm -r locks
rm -r books
rm -r code

# make HTML files
cp -r md src

find src -name "*.md" | cut -d. -f1 | while read i; do
	markdown -o $i.html $i.md
	rm $i.md
done

# compile and run cite
make clean cite
./cite

# cleanup
rm -r src

find md -name "*.html" | while read i; do
	rm $i
done
