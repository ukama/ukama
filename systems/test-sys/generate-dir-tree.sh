#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

echo "Readme generated based on you proto file at path:"
cd $1
pwd;find -s . -not -path '*/.*' -not -path '*/*.txt'  -print  2>/dev/null|awk '!/\.$/ {for (i=1;i<NF-1;i++){printf("      │")}print "      ├── "$NF}'  FS='/' &> temp.txt
{ IFS= read -rd '' value <temp.txt;} 2>/dev/null
echo "Generate directory tree"
export DIR_CONTENT=$(printf '%s' "$value")
echo "Modifying readme with directory tree..."
sleep 2
file_contents=$(<./README.md)
echo "${file_contents//#DIR_CONTENT#/$DIR_CONTENT}" > ./README.md
rm -f temp.txt
echo "Completed."
