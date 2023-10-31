#!/bin/sh
find . -name "*.pcss" | while read l; do
  name=$(echo $l |cut -d. -f2|cut -d/ -f2-)
  echo $name
  cat "$l" > "$name.css"
  echo "" >> "$name.css"
  cat style.css >> "$name.css"
done
