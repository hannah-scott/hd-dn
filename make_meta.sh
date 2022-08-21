#!/bin/bash

# delete old meta page
rm meta/index.html

# get list of every page that's not an index
echo "" > temp
find . -name "*.html" -not \( -name "draft*" -o -name "*index*" \) | \
    # alphabetise by filename
    sed -e "s/\/\([a-z-]*\.html\)/~\1/g" | \
    sort -t~ -k2 | \
    sed -e "s/~/\//g" | \
    # write listitem
    while read l; do
        echo "<li><a href=\"$(echo $l | cut -d. -f2-)\">" >> temp
        echo "$(echo $l | grep -o "[^/]*$")" >> temp
        echo "</a></li>" >> temp
done

# Build meta page
echo """<!DOCTYPE html>
<html lang=\"en\">
    <head>
        <meta charset=\"UTF-8\" />
        <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\" />
        <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />
        <title>HD-DN: Meta</title>
        <link rel=\"stylesheet\" href=\"/css/style.css\" />
    </head>
    <body>
        <nav class=\"menu\">
            <h1><a href=\"/\">HD-DN</a></h1>
            <ul>
                <li><a href=\"/about/\">About</a></li>
                <li><a href=\"/bookmarks/\">Bookmarks</a></li>
                <li><a href=\"/books/\">Books</a></li>
                <li><a href=\"/posts/\">Posts</a></li>
                <li><a href=\"/three-good-things\">Three Good Things</a></li>
                <li><a href=\"/toki-pona/\">Toki Pona</a></li>
                <li><a href=\"/quotes/\">Quotes</a></li>
            </ul>
        </nav>
        <h1 class=\"title\">Meta</h1>
        
        <p>This is a list of every non-index HTML file on the site.</p>

        <ul>
""" > meta/index.html

cat temp >> meta/index.html
rm temp

echo """</ul>
    </body>
</html>""" >> meta/index.html