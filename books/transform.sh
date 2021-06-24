#!/bin/bash
file=$(echo $1 | cut -d. -f1)
sed -e "s/\#\#\#/\#\#/g" -i $file.md 
markdown -o $file.html $file.md
sed -e "s/<p>/<p class=\"my\-1\">/g" -i $file.html
sed -e "s/<h2>/<\/div><div class\=\"book\-summary my\-1\"><h2>/g" -i $file.html 