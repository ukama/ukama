#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# shutdown.sh 'system' system-path system-name
# shutdown.sh 'node' node-id

if [ "$1" = "system" ]; then

    cd "$2" || exit 1
    docker-compose down -v || { echo "$system shutdown: FAILED"; exit 1; }
    cd -

    echo "$system shutdown: OK"

elif [ "$1" = "node" ]; then

    node_id=$2
    echo '{"execute":"guest-shutdown"}' | nc -U /tmp/qemu-monitor-${node_id}.sock

    # wait for 10 seconds before taking next steps
    sleep 10

    if ps aux | grep -q "qemu-system-x86_64 -hda ${node_id}.img"; then
        echo "Forcing node shutdown. Id: ${node_id}"
        sudo pkill -f "qemu-system-x86_64 -hda ${node_id}.img" \
             || { echo "Unable to kill node with id: ${node_id}"; exit 1; } 
    fi
else
    exit 1
fi

exit 0
