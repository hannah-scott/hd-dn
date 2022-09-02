#!/bin/bash
cp backups/*.jpg .
for file in *.jpg; do
    convert $file -strip -resize 720 -define jpeg:extent=50KB $file
done