#!/bin/bash

# Define the coverage packages
cover_packages="github.com/ukama/ukama/systems/nucleus/user/pb/gen,github.com/ukama/ukama/systems/registry/network/pb/gen"

# Run tests and generate a coverage profile
go test -coverpkg="$cover_packages" ./... -coverprofile=coverage.out -covermode=count && go tool cover -func ./coverage.out

# Input coverage file
coverage_file="coverage.out"

# Use awk to process the file in-place
awk 'NR==1 {print; next} /ukama/ {print}' "$coverage_file" > "$coverage_file.temp" && mv "$coverage_file.temp" "$coverage_file"

