#!/bin/bash
in=$1
out=$2

echo "<article>
    <h2>$(date +%F)</h2>" > $out

markdown $in >> $out

echo "</article>" >> $out