#!/bin/bash

# Get the latest Git tag
LATEST_TAG=$(git describe --tags --abbrev=0)
# Get the current commit hash
COMMIT_HASH=$(git rev-parse --short HEAD)

# Combine to form the full version
VERSION="${LATEST_TAG}-${COMMIT_HASH}"

# Generate version.h
echo "#ifndef VERSION_H_" > version.h
echo "#define VERSION_H_" >> version.h
echo "#define VERSION \"${VERSION}\"" >> version.h
echo "#endif /* VERSION_H_ */" >> version.h
