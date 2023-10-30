#!/bin/bash

# Check if a directory argument is provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <directory>"
  exit 1
fi

# Extract the directory from the command-line argument
directory="$1"

# Initialize a counter for executable files
count=0

# Search for executable files (excluding .sh files)
find "$directory" -type f -executable ! -name "*.sh" -print |
  while read -r file; do
    echo "Executable file found: $file"
    ((count++))
  done

if [ "$count" -eq 0 ]; then
  echo "No executable files (excluding .sh files) found in the directory or its subdirectories."
else
  echo "Found $count executable file(s) (excluding .sh files) in the directory and its subdirectories."
fi

