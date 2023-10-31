#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2021-present, Ukama Inc.

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

