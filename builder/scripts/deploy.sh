# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[38;5;39m'
NC='\033[0m'
TAG="${BLUE}Ukama>${NC}"
root_dir=$(pwd)
# Parse the JSON file and initialize the variables
MASTERORGNAME="ukama"
AUTHSYSKEY="auth-services"
BILLINGSYSKEY="billing"
OWNEREMAIL=$(jq -r '.setup.email' "$1")
OWNERNAME=$(jq -r '.setup.name' "$1")
ORGNAME=$(jq -r '.setup["org-name"]' "$1")
ORGID=$(jq -r '.setup["org-id"]' "$1")
SYS=$(jq -r '.systems' "$1")
KEY=$(jq -r '.key' "$1")
METADATA=$(jq -c '.' ../metadata.json)
MAILERHOST=$(jq -r '.mailer.host' "$1")
MAILERPORT=$(jq -r '.mailer.port' "$1")
MAILERUSERNAME=$(jq -r '.mailer.username' "$1")
MAILERPASSWORD=$(jq -r '.mailer.password' "$1")
LAGOAPIKEY=$(jq -r '."lago-api-key"' "$1")
if [[ "$(uname)" == "Darwin" ]]; then
    # For Mac
    LOCAL_HOST_IP=$(ifconfig en0 | grep inet | awk '$1=="inet" {print $2}')
elif [[ "$(uname)" == "Linux" ]]; then
    # For Linux
    LOCAL_HOST_IP=$(ifconfig enp0s25 | grep inet | awk '$1=="inet" {print $2}')
fi
OWNERAUTHID=""
OWNERID=$(uuidgen)

function set_env() {
    export OWNERID=$OWNERID
    export OWNERAUTHID=$OWNERAUTHID
    export OWNERNAME=$OWNERNAME
    export OWNEREMAIL=$OWNEREMAIL
    export ORGID=$ORGID
    export ORGNAME=$ORGNAME
    export KEY=$KEY
    export MAILER_PORT=$MAILERPORT
    export MAILER_HOST=$MAILERHOST
    export MAILER_PASSWORD=$MAILERPASSWORD
    export MAILER_USERNAME=$MAILERUSERNAME
    export MAILER_FROM=$OWNEREMAIL
    export TEMPLATESPATH=member-invite
    export LAGO_API_KEY=$LAGOAPIKEY
    export MASTERORGNAME=$MASTERORGNAME
    export LOCAL_HOST_IP=$LOCAL_HOST_IP
    
}

function run_docker_compose() {
    set_env
    echo  "$TAG Running $2 docker compose..."
    cd $1

    CONTAINER_NAME=$3
    while true; do
        docker compose up --build -d > /dev/null 2>&1
        # docker-compose down  > /dev/null 2>&1
        # docker-compose build > /dev/null 2>&1
        # docker compose up --build -d > /dev/null 2>&1
        sleep 5
        if docker ps | grep -q $CONTAINER_NAME; then
            echo  "$TAG $2 docker container is up"
            break
        else
            echo "Container $CONTAINER_NAME is not running. Retrying..."
        fi
    done
}

function register_user() {
    echo  "$TAG Signing up Owner user"
    flow=$(curl --location --silent "http://${LOCAL_HOST_IP}:4434/self-service/registration/api")
    action=$(echo $flow | jq -r '.ui.action')
    if [[ ! $action =~ ^http(s)?:// ]]; then
        echo "Invalid URL: $flow"
        exit 1
    fi
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
    sleep 2
    identity=$(echo $response | jq -r '.session.identity')
    sleep 2
    OWNERAUTHID=$(echo $identity | jq -r '.id')
    echo  "$TAG User register with ${GREEN}AUTH ID = $OWNERAUTHID${NC}"
    echo  "$TAG Please verify your email address by visiting ${GREEN}http://localhost:4436${NC}"
}

sort_systems_by_dependency() {
    IFS=', ' read -r -a SYSTEMS_ARRAY <<< "$SYS"

    if [[ ! " ${SYSTEMS_ARRAY[@]} " =~ " services " ]] || [[ ! " ${SYSTEMS_ARRAY[@]} " =~ " auth-services " ]]; then
        echo "Error: 'services' and 'auth-services' are required in the systems array in the deploy_config JSON file."
        exit 1
    fi

    if [[ " ${SYSTEMS_ARRAY[@]} " =~ " billing " ]] && ([[ ! " ${SYSTEMS_ARRAY[@]} " =~ " dataplan " ]] || 
        [[ ! " ${SYSTEMS_ARRAY[@]} " =~ " subscriber " ]] || [[ ! " ${SYSTEMS_ARRAY[@]} " =~ " notification " ]]); then
        echo "Error: 'billing' depend on dataplan, subscriber and notification, please make sure these systems are added in the deploy_config JSON file."
        exit 1
    fi

    SYSTEMS=($(for key in "${SYSTEMS_ARRAY[@]}"; do
      if [ "$key" == "services" ]; then
          echo "1 $key"
      elif [ "$key" == "auth-services" ]; then
          echo "2 $key"
      elif [ "$key" == "init" ]; then
          echo "3 $key"
      elif [ "$key" == "nucleus" ]; then
          echo "4 $key"
      elif [ "$key" == "registry" ]; then
          echo "5 $key"
      elif [ "$key" == "dataplan" ]; then
          echo "6 $key"
      elif [ "$key" == "subscriber" ]; then
          echo "7 $key"
      elif [ "$key" == "notification" ]; then
          echo "8 $key"
      else
          echo "9 $key"
      fi
    done | sort -n -k1,1 | cut -d' ' -f2-))
}

sort_systems_by_dependency

# Loop through the SYSTEMS array
for SYSTEM in "${SYSTEMS[@]}"; do
    cd ~
    cd $root_dir
    if [ "$SYSTEM" != $AUTHSYSKEY ]; then
        cd ../../systems
    fi
    if [ "$SYSTEM" == $BILLINGSYSKEY ]; then
        cd ./billing/provider
        chmod +x start_provider.sh
        ./start_provider.sh
        cd ../..
    fi
    SYSTEM_OBJECT=$(echo "$METADATA" | jq -c --arg SYSTEM "$SYSTEM" '.[$SYSTEM]')
    export COMPOSE_PROJECT_NAME=$(echo "$SYSTEM_OBJECT" | jq -r '.key')
    run_docker_compose "$(echo "$SYSTEM_OBJECT" | jq -r '.path')" "$(echo "$SYSTEM_OBJECT" | jq -r '.name')" "$(echo "$SYSTEM_OBJECT" | jq -r '.key')"
    case $SYSTEM in
    "auth-services")
        sleep 2
        register_user
        ;;
    "init")
        sleep 2
        echo  "$TAG Add org in lookup..."
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5401/lookup"
        QUERY="INSERT INTO \"public\".\"orgs\" (\"created_at\", \"updated_at\", \"name\", \"org_id\", \"certificate\") VALUES (NOW(), NOW(), '$ORGNAME', '$ORGID', 'ukama-cert')"
        psql $DB_URI -c "$QUERY"
        ;;
    esac
    cd ../
done