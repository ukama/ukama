#!/bin/sh
set -e
TARGET_EXEC=$1
echo "Copying libs for $TARGET_EXEC"
rm -rf ./libs 
mkdir ./libs

#Copying dependencies 
ldd $TARGET_EXEC | awk 'NF == 4 { system("cp -v " $3 " ./libs") }'
