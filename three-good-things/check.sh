#!/bin/bash
# Check that moku has picked up dates
cat atom.atom | grep -E "[0-9]{4}-[0-9]{2}-[0-9]{2}" | head
