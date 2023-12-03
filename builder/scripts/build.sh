#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

if [ "$1" = "systems" ]; then

    cd "$2" || exit 1

    find . -type f -name 'go.mod' | while read -r modfile; do
        dir=$(dirname "$modfile")
        cd "$dir" || exit 1
        go mod tidy
        make clean
        make
        cd - || exit 1
    done

    # cleanup and build with using cached layers
    docker-compose down -v || exit 1
    docker image prune -f || exit 1
    docker-compose build --no-cache || exit 1

elif [ "$1" = "node" ]; then
    echo "Node build not required."
else
    echo "Invalid argument: $1. Use 'systems' or 'node'."
    exit 1
fi

exit 0
