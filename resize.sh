#!/bin/bash
cd ./static/img
find . -name "*.jpg" | while read file; do
    convert $file -resize 720x -strip -define jpeg:extent=64KB $file
done
