#!/bin/bash
tidy -m index.html
./moku.py
./check.sh
