#!/bin/bash
for file in *.jpg; do
    mv $file $file.original
    convert $file.original -resize 720 -strip $file
done