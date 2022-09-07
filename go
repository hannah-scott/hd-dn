#!/bin/bash
cd three-good-things
./moku.py
cd ..
cd img
./resize.sh
cd ..
./make_meta.sh
git add .
git commit
git push