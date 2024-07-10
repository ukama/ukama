#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# Install the required packages
pip install -r requirements.txt

# Loop over each swagger file (one for each app)
for json_file in "specs"/*.json; do
    echo "Processing $json_file..."

    python3 test_server.py "$json_file" &
    server_pid=$!

    # Wait a few seconds to ensure the server is up and running
    sleep 5

    python3 test_client.py "$json_file"

    # Kill the server process
    kill $server_pid

    echo "Finished processing $json_file."
done

echo "All JSON files processed."
