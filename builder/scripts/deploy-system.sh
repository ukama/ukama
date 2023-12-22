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
TAG="${BLUE}Ukama>${NC}"

set_env() {

    json_file=$1

    export OWNEREMAIL=$(jq -r '.deploy.env.owneremail' "$json_file")
    export OWNERNAME=$(jq -r '.deploy.env.ownername' "$json_file")
    export ORGNAME=$(jq -r '.deploy.env.orgname' "$json_file")
    export ORGID=$(jq -r '.deploy.env.orgid' "$json_file")
    export KEY=$(jq -r '.deploy.env.key' "$json_file")
    export MAILERHOST=$(jq -r '.deploy.env.mailer_host' "$json_file")
    export MAILERPORT=$(jq -r '.deploy.env.mailer_port' "$json_file")
    export MAILERUSERNAME=$(jq -r '.deploy.env.mailer_username' "$json_file")
    export MAILERPASSWORD=$(jq -r '.deploy.env.mailer_password' "$json_file")
    export LAGOAPIKEY=$(jq -r '.deploy.env.lago-api-key' "$json_file")
    export LOCAL_HOST_IP=$(jq -r '.deploy.env.local_host_ip' "$json_file")
}

register_user() {
    echo  "$TAG Signing up Owner user"
    flow=$(curl --location --silent 'http://localhost:4434/self-service/registration/api')
    action=$(echo $flow | jq -r '.ui.action')
    response=$(curl --location --request POST "$action" \
     --header 'Content-Type: application/json' \
     --data-raw '{
        "method": "password",
        "password": "@Pass2021",
        "traits": {
            "email": "'$OWNEREMAIL'",
            "name": "'$OWNERNAME'",
            "firstVisit": true
        }
    }')

    identity=$(echo $response | jq -r '.session.identity')
    OWNERAUTHID=$(echo $identity | jq -r '.id')

    echo  "$TAG User register with ${GREEN}AUTH ID = $OWNERAUTHID${NC}"
    echo  "$TAG Please verify your email address by visiting ${GREEN}http://localhost:4436${NC}"
}

if [ "$1" = "system" ]; then

    system=$2
    path=$3
    json_file=$4
    cwd=`pwd`

    set_env $json_file
    cd "$path" || exit 1

    echo  "$TAG Running $system ..."
    docker-compose up -d > /dev/null 2>&1
    echo  "$TAG $system is up ..."

    case $system in
        "ukama-auth")
            register_user
            ;;
        "init")
            sleep 2
            ## Connect to init-lookup db and add org in orgs table
            echo  "$TAG Add org in lookup..."
            DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5401/lookup"
            QUERY="INSERT INTO \"public\".\"orgs\" (\"created_at\", \"updated_at\", \"name\", \"org_id\", \"certificate\") VALUES (NOW(), NOW(), '$ORGNAME', '$ORGID', 'ukama-cert')"
            psql $DB_URI -c "$QUERY" || exit 1
            ;;
    esac

    cd $cwd

elif [ "$1" = "node" ]; then
    image_file=$2.img
    sudo qemu-system-x86_64 -hda ${image_file} -m 1024 -kernel ./vmlinuz-5.4.0-26-generic \
         -initrd ./initrd.img-5.4.0-26-generic -append "root=/dev/sda1" || exit 1
fi

exit 0
