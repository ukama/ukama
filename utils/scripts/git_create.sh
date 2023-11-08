#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <name>"
  exit 1
fi

name="$1"
repository_url="https://github.com/ukama/$name.git"

# Create a directory with the provided name and navigate to it
cd "$name"

# Initialize a new Git repository with the 'main' branch
git init -b main

# Add all files to the Git staging area
git add .

# Commit the changes with a message
git commit -m "$name"

# Add a remote named 'origin' with the provided GitHub repository URL
git remote add origin "$repository_url"

# Push the changes to the 'main' branch on the remote repository
git push origin main

# Navigate back to the previous directory
cd ..

echo "Git commands executed successfully for repository: $repository_url"
