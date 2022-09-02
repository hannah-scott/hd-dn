#!/bin/bash
cd three-good-things
./moku.py
cd ..
cd img
./resize.sh
cd ..
git add .
git commit
git push