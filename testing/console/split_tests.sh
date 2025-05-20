#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

mkdir -p tests/split_tests

grep -n "^test(" tests/patched_test.spec.ts | while read -r line; do
    line_num=$(echo "$line" | cut -d: -f1)
    test_name=$(echo "$line" | sed -E 's/.*test\(["'"'"']([^"'"'"']+)["'"'"'].*/\1/')
    
    test_name=$(echo "$test_name" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')
    
    start_line=$line_num
    brace_count=0
    end_line=$start_line
    
    while IFS= read -r current_line; do
        open_braces=$(echo "$current_line" | tr -cd '{' | wc -c)
        close_braces=$(echo "$current_line" | tr -cd '}' | wc -c)
        
        brace_count=$((brace_count + open_braces - close_braces))
        end_line=$((end_line + 1))
        
        if [ $brace_count -eq 0 ]; then
            break
        fi
    done < <(tail -n +$start_line tests/patched_test.spec.ts)
    
    {
        sed -n '/^import/,/^test(/p' tests/patched_test.spec.ts | sed '$d' | \
        sed 's|from '\''../constants'\''|from '\''../../constants'\''|g' | \
        sed 's|from '\''../utils'\''|from '\''../../utils'\''|g'
        grep "test.setTimeout" tests/patched_test.spec.ts
        echo
        sed -n "${start_line},${end_line}p" tests/patched_test.spec.ts | sed 's/^[[:space:]]*$//' | sed '/^$/d'
    } > "tests/split_tests/${test_name}.spec.ts"
done

echo "Tests have been split into separate files in the tests/split_tests directory"