#!/bin/bash
today=$(date +%F)

echo """---
title: $today
date: $today
author: Hannah Scott
---

""" > $today.md

vim $today.md
