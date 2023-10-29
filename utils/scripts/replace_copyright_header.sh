#!/bin/bash
SOURCE_DIR=`dirname $0`
OK_COLOR="\e[32m"
SKIPPED_COLOR="\e[31m"
RESET_COLOR="\e[0m"

if [ $# != 1 ]; then
    echo "Usage: replace_header.sh <file>"
    exit 1
fi

input_file="$1"

# Check if "Ukama Inc" exists in the input file
if ! grep -q "Ukama Inc" "$input_file"; then
    printf "%-40s copyright change ... [${SKIPPED_COLOR}Skipped${RESET_COLOR}]\n" "$input_file"
    exit 1
fi

# Extract the year from the input file
year=$(grep -oE '\b(2021|2022|2023)\b' "$input_file")

if [ -z "$year" ]; then
    printf "%-40s copyright change ... [${SKIPPED_COLOR}Skipped${RESET_COLOR}]\n" "$input_file"
    exit 1
fi

# Determine the corresponding header template
template_file="$SOURCE_DIR/header.template.$year"

if [ ! -f "$template_file" ]; then
    echo "Header template for the year $year not found."
    exit 1
fi

cat $template_file > $SOURCE_DIR/tmp
cat $1 | awk -f $SOURCE_DIR/remove_header.awk >> $SOURCE_DIR/tmp
mv $SOURCE_DIR/tmp $1

printf "%-40s copyright change ... [${OK_COLOR}OK${RESET_COLOR}]\n" "$input_file"

