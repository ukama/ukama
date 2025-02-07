#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# HOW TO USE: ./push_node_event.sh <orgname>

ORGNAME=$1

DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5414/component"
INVENTORY_NODE_IDS=("uk-sa9001-tnode-a1-1234" "uk-sa9001-anode-a1-1234")

NODE_IDS=("uk-sa2450-tnode-v0-4e86" "uk-sa2450-hnode-v0-4e87")
LATITUDES=(-4.3262161 -4.32758)
LONGITUDES=(15.311631 15.3109951)

for i in {0..1}; do
  RESPONSE=$(./bin/msgcli events push --org $ORGNAME --route messaging.mesh.node.online -m "{\"NodeId\":\"${NODE_IDS[$i]}\"}" 2>&1)
  lowercase_nodeid=$(echo "${NODE_IDS[$i]}" | tr '[:upper:]' '[:lower:]')

  SQL_QUERY="UPDATE components SET part_number = '$lowercase_nodeid' WHERE part_number = '${INVENTORY_NODE_IDS[$i]}';"
  psql $DB_URI -t -A -c "$SQL_QUERY"
  echo "Node online event pushed and nodeid updated in inventory for $lowercase_nodeid"

  if [ $i -eq 0 ]; then
    curl_response=$(curl -s -o /dev/stdout -w "\nHTTP Status: %{http_code}\n" -X 'POST' \
      'http://localhost:8036/v1/notify' \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{
      "details": {
        "latitude": -4.3262161,
        "longitude": 15.311631
      },
      "node_id": "uk-sa2450-tnode-v0-4e86",
      "service_name": "health",
      "severity": "low",
      "status": 8100,
      "time": 1733753967,
      "type": "event"
    }')
    echo "Curl response: $curl_response"
  fi
done

