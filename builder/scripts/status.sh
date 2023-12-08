#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[38;5;39m'
NC='\033[0m'
IFS=',' read -r -a SYSTEMS <<< "$1"

#echo "$TAG Checking Docker container status..."
while true; do
    for PROJECT_NAME in "${SYSTEMS[@]}"; do

        printf "System: $PROJECT_NAME \n"

        DOCKER_PS_OUTPUT=$(docker ps -a --no-trunc --filter "label=com.docker.compose.project="$PROJECT_NAME"" --filter "network=services_dev-net" --format "{{json . }}" | jq -c 'if (.Networks | test(":")) then (.Networks |= fromjson) else . end | {ID: .ID, Names: .Names, Networks: .Networks, Status: .Status}')
        if [ -z "$DOCKER_PS_OUTPUT" ]; then
            continue
        fi

        DOCKER_PS_OUTPUT_NEW=()
        while IFS= read -r line; do
            ID=$(echo "$line" | jq -r '.ID')
            if [ -z "$ID" ]; then
                continue
            fi
            LIVE_STATUS=$(docker inspect -f '{{.State.Status}}' "$ID")
            line=$(echo "$line" | jq --arg LIVE_STATUS "$LIVE_STATUS" '.Status = $LIVE_STATUS')
            DOCKER_PS_OUTPUT_NEW+=("$line")
        done <<< "$DOCKER_PS_OUTPUT"

        NEW_OBJECT=$(echo "${DOCKER_PS_OUTPUT_NEW[@]}" | jq -s 'reduce .[] as $item ({}; .ID += [$item.ID] | .Names += [$item.Names] | .Networks += [$item.Networks] | .Status += [$item.Status])')
        NAMES=($(jq -r '.Names[]' <<< "$NEW_OBJECT"))
        STATUSES=($(jq -r '.Status[]' <<< "$NEW_OBJECT"))

        INDEX=0
        for i in ${NAMES[@]}; do
            case "${STATUSES[INDEX]}" in
                "running") COLOR=$GREEN;;
                "exited") COLOR=$RED;;
                *) COLOR=$YELLOW;;
            esac
            printf "%-30s %b%s%b\n" "    ${NAMES[INDEX]}" "$COLOR" "${STATUSES[INDEX]}" "$NC"
            INDEX=$((INDEX+1))
        done
    done

    sleep $2
    clear
done
