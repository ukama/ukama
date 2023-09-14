#!/bin/bash

# Input coverage file
coverage_file="coverage.out"

# Use awk to process the file in-place
awk 'NR==1 {print; next} /ukama/ {print}' "$coverage_file" > "$coverage_file.temp" && mv "$coverage_file.temp" "$coverage_file"

