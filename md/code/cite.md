cite is a static site generator written in ANSI C. It's under active
development, and licensed under the MIT License.

## Get it

You can get source code [here](https://github.com/hannah-scott/cite).

## Use

I build hddn with the discount Markdown toolset, and cite. I chose discount for
three reasons:

1. It fulfils my needs
2. It's smaller and faster than Pandoc
3. It's written in C, so I can read its source without learning a new language

These are glued together with a build script. 

	
	#!/bin/bash
	cd ~/repos/hannahs.ddns.net
	
	# cleanup
	rm *.html 
	rm -r locks
	rm -r books
	
	# make HTML files
	cd ~/repos/repo.hannahs.ddns.net
	rm -r src
	cp -r md src
	
	find src -name "*.md" | cut -d. -f1 | while read i; do
		markdown -o $i.html $i.md
		rm $i.md
	done
	
	# compile and run cite
	make clean cite
	./cite

	
