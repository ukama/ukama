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

# Search for Makefiles, C files, and header files
find "$directory" -type f \( -name "Makefile" -o -name "*.c" -o -name "*.h" \) -print0 |
  while IFS= read -r -d '' file; do
    # Check if the file contains the specified copyright line
    if grep -qE "(^#|^\s*/\*|^\*)?\s*Copyright \(c\) 20(21|22|23)-present, Ukama Inc\." "$file"; then
      echo "File with specified copyright: $file"
    fi
  done

echo "Search completed."

