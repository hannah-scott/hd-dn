#!/bin/bash
for file in $1/*.jpg; do
    convert $file -strip -resize 720 -define jpeg:extent=100KB $file
done