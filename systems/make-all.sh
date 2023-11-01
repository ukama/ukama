#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# To run "chmod +x make-all.sh && ./make-all.sh"
services=(
    "auth/api-gateway" 
    "data-plan/base-rate" 
    "data-plan/package" 
    "data-plan/rate" 
    "data-plan/api-gateway"
    "subscriber/registry"
    "subscriber/sim-manager"
    "subscriber/sim-pool"
    "subscriber/test-agent"
    "subscriber/api-gateway"
    "nucleus/org"
    "nucleus/user"
    "nucleus/api-gateway"
    "registry/node"
    "registry/network"
    "registry/invitation"
    "registry/member"
    "registry/api-gateway"
    "services/msgClient"
    "notification/mailer"
    "notification/node-gateway"
    "notification/notify"
    "notification/api-gateway"
    "init/lookup"
    "init/node-gateway"
    "init/api-gateway"
    "metrics/api-gateway"
    )


# Loop through each path
for path in "${services[@]}"; do
    # Change directory to the current path
    cd "$path" || { echo "Failed to change directory to $path"; exit 1; }

    # Run the "make" command
    go mod tidy && make gen && make

    # Check the exit status of the "make" command
    if [ $? -eq 0 ]; then
        echo "Make completed successfully in $path"
    else
        echo "Make failed in $path"
    fi

    # Return to the original directory (optional)
    cd - >/dev/null
done
