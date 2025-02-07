#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# Set the username
USERNAME="$1"

# Prompt the user for the password securely
read -s -p "Enter password for $USERNAME: " PASSWORD
echo

# Set the path to the file containing the links
LINK_FILE="srtm90m_test.txt"

# Read each line of the file and process it as a URL
while read -r URL; do
  echo "Processing URL: $URL"

  # Use wget to download the file with authentication
  wget --user="$USERNAME" --password="$PASSWORD" "$URL"

done < "$LINK_FILE"


# Loop through each zip file and extract it
for ZIP_FILE in *.zip; do
  echo "Extracting $ZIP_FILE..."

  # Extract the zip file using unzip
  unzip "$ZIP_FILE"

  # Delete the zip file after extraction
  rm "$ZIP_FILE"
done

shopt -s nullglob
for f in ./*.hgt; do
    echo "Converting $f"
    srtm2sdf $f
    rm $f
done
