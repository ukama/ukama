#!/bin/sh
set -e
TARGET_EXEC=$1
rm -rf ./libs 
mkdir ./libs

# Logs
ldd $TARGET_EXEC | awk 'NF == 4 { system("echo cp " $3 " ./libs") }'

echo "Copying libs for $TARGET_EXEC"

#Copying dependencies 
ldd $TARGET_EXEC | awk 'NF == 4 { system("cp " $3 " ./libs") }'


if [ -d "/home/runner/work/ukama/ukama/nodes/ukamaOS/distro/vendor/build/lib" ]; then
    echo "Workaround for microhttpd."
    cp /home/runner/work/ukama/ukama/nodes/ukamaOS/distro/vendor/build/lib/libmicrohttpd.* ./libs 
else 
    echo "Nothing required"
fi

sleep 5;

echo "Copied files"
ls -ltr ./libs

