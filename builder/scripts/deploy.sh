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
JSON_FILE="../deploy_config.json"
MASTERORGNAME="ukama"
AUTHSYSKEY="auth-services"
IS_INCLUDE_BFF=false
BILLINGSYSKEY="billing"
OWNEREMAIL=$(jq -r '.setup.email' "$JSON_FILE")
PASSWORD=$(jq -r '.setup.password' "$JSON_FILE")
OWNERNAME=$(jq -r '.setup.name' "$JSON_FILE")
ORGNAME=$(jq -r '.setup["org-name"]' "$JSON_FILE")
ORGID=$(jq -r '.setup["org-id"]' "$JSON_FILE")
COUNTRY=$(jq -r '.setup["country"]' "$JSON_FILE")
CURRENCY=$(jq -r '.setup["currency"]' "$JSON_FILE")
SYS=$(jq -r '.systems' "$JSON_FILE")
KEY=$(jq -r '.key' "$JSON_FILE")
METADATA=$(jq -c '.' ../metadata.json)
LAGOAPIKEY=$(jq -r '."lago-api-key"' "$JSON_FILE")
if [[ "$(uname)" == "Darwin" ]]; then
    # For Mac
    LOCAL_HOST_IP=$(ifconfig en0 | grep inet | awk '$1=="inet" {print $2}')
elif [[ "$(uname)" == "Linux" ]]; then
    # For Linux
    LOCAL_HOST_IP=$(ifconfig enp0s25 | grep inet | awk '$1=="inet" {print $2}')
fi
OWNERAUTHID=""
OWNERID=$(uuidgen)

function buildSystems() {
    echo  "$TAG Building systems..."
    ./make-sys-for-mac.sh ../deploy_config.json 2>&1 | tee buildSystems.log
}

while getopts "b" opt; do
    case ${opt} in
        b )
            buildSystems
            ;;
        \? )
            echo "Usage: cmd [-b]"
            exit 1
            ;;
    esac
done

function set_env() {
    export OWNERID=$OWNERID
    export OWNERAUTHID=$OWNERAUTHID
    export OWNERNAME=$OWNERNAME
    export OWNEREMAIL=$OWNEREMAIL
    export ORGID=$ORGID
    export ORGNAME=$ORGNAME
    export KEY=$KEY
    export LAGO_API_KEY=$LAGOAPIKEY
    export MASTERORGNAME=$MASTERORGNAME
    export LOCAL_HOST_IP=$LOCAL_HOST_IP
    export COMPONENT_ENVIRONMENT=test
    export COUNTRY=$COUNTRY
    export CURRENCY=$CURRENCY
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
        "password": "'$PASSWORD'",
        "traits": {
            "email": "'$OWNEREMAIL'",
            "name": "'$OWNERNAME'"
        }
    }')
    sleep 2
    identity=$(echo $response | jq -r '.session.identity')
    sleep 2
    OWNERAUTHID=$(echo $identity | jq -r '.id')
    echo  "$TAG User register with ${GREEN}AUTH ID = $OWNERAUTHID${NC}"
    echo  "$TAG Please verify your email address by visiting ${GREEN}http://localhost:4436${NC}"
}

SYSTEMS=()

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

    for key in "${SYSTEMS_ARRAY[@]}"; do
        case "$key" in
            "services")
                SYSTEMS+=("1 $key")
                ;;
            "auth-services")
                SYSTEMS+=("2 $key")
                ;;
            "init")
                SYSTEMS+=("3 $key")
                ;;
            "nucleus")
                SYSTEMS+=("4 $key")
                ;;
            "inventory")
                SYSTEMS+=("5 $key")
                ;;
            "registry")
                SYSTEMS+=("6 $key")
                ;;
            "dataplan")
                SYSTEMS+=("7 $key")
                ;;
            "node")
                SYSTEMS+=("8 $key")
                ;;
            "billing")
                SYSTEMS+=("9 $key")
                ;;
            "subscriber")
                SYSTEMS+=("10 $key")
                ;;
            *)
                SYSTEMS+=("11 $key")
                ;;
        esac
    done

    SYSTEMS=($(for item in "${SYSTEMS[@]}"; do echo "$item"; done | sort -n -k1,1 | cut -d' ' -f2-))
}

sort_systems_by_dependency

IS_INVENTORY_SYS=false
INVENTORY_SYS_KEY="inventory"
IS_INIT_SYS="init"

cleanup() {
    echo  "$TAG Cleaning up..."
    cd $root_dir
    cd ../../systems
    for SYSTEM in "${SYSTEMS[@]}"; do
        cd ~
        cd $root_dir
        if [ "$SYSTEM" == $AUTHSYSKEY ]; then
            cd ../../../ukama-auth/kratos
            sed -i '' "s/\${LOCAL_HOST_IP}/$LOCAL_HOST_IP/g" "kratos.yml"
            cd ../../ukama/builder/scripts
        fi
        if [ "$SYSTEM" != $AUTHSYSKEY ]; then
            cd ../../systems
        fi
        SYSTEM_OBJECT=$(echo "$METADATA" | jq -c --arg SYSTEM "$SYSTEM" '.[$SYSTEM]')
        cd "$(echo "$SYSTEM_OBJECT" | jq -r '.path')"
        if [ -d ".temp" ]; then
            rm -rf docker-compose.yml
            mv ".temp/docker-compose.yml" .
            rm -rf ".temp"
        fi
    done
    cd $root_dir
}

setup_docker_compose_files(){
    for SYSTEM in "${SYSTEMS[@]}"; do
        cd ~
        cd $root_dir
        if [ "$SYSTEM" == $AUTHSYSKEY ]; then
            cd ../../../ukama-auth/kratos
            sed -i '' "s/\${LOCAL_HOST_IP}/$LOCAL_HOST_IP/g" "kratos.yml"
            cd ../../ukama/builder/scripts
        fi
        if [ "$SYSTEM" != $AUTHSYSKEY ]; then
            cd ../../systems
        fi
        SYSTEM_OBJECT=$(echo "$METADATA" | jq -c --arg SYSTEM "$SYSTEM" '.[$SYSTEM]')
        cd "$(echo "$SYSTEM_OBJECT" | jq -r '.path')"
        mkdir -p ".temp"
        cp docker-compose.yml ".temp"
        if [[ "$(uname)" == "Darwin" ]]; then
            sed -i '' "s/build: \.\.\/services\/initClient/image: main-init/g" docker-compose.yml
        fi
    done
    cd $root_dir
}

setup_docker_compose_files

# Loop through the SYSTEMS array
for SYSTEM in "${SYSTEMS[@]}"; do
    cd ~
    cd $root_dir
   if [ "$SYSTEM" == $AUTHSYSKEY ]; then
        cd ../../../ukama-auth/kratos
        sed -i '' "s/\${LOCAL_HOST_IP}/$LOCAL_HOST_IP/g" "kratos.yml"
        cd ../app
        cp .env.example .env.local
        cd ../../ukama/builder/scripts
    fi
    if [ "$SYSTEM" == $INVENTORY_SYS_KEY ]; then
        IS_INVENTORY_SYS=true
    fi
    if [ "$SYSTEM" == 'bff' ]; then
        IS_INCLUDE_BFF=true
    fi
    if [ "$SYSTEM" != $AUTHSYSKEY ]; then
        cd ../../systems
    fi
    if [ "$SYSTEM" == $BILLINGSYSKEY ]; then
        echo "$TAG Current directory before billing setup: $(pwd)"
        
        # First, run the billing provider script
        cd ./billing/provider
        echo "$TAG Running billing provider in: $(pwd)"
        export COMPOSE_PROJECT_NAME="billing-provider"
        chmod +x start_provider.sh
        ./start_provider.sh
        cd ../..
        
        # Navigate to payments system and run deployment
        echo "$TAG Setting up payments system..."
        cd ../../payments/builder-script
        export COMPOSE_PROJECT_NAME="payments"
        chmod +x deploy.sh
        ./deploy.sh services $ORGNAME

        # Navigate to webhooks system and run deployment
        echo "$TAG Setting up webhooks system..."
        cd ../../webhooks/builder-script
        export COMPOSE_PROJECT_NAME="webhooks"
        chmod +x deploy.sh
        ./deploy.sh services $ORGNAME
        
        # Now enter the testing and hooks setup
        echo "$TAG Setting up hooks..."
        cd ../../ukama/testing/services/hooks
        export COMPOSE_PROJECT_NAME="hooks"
        docker-compose up -d
        cd ../../../../ukama/systems

        export COMPOSE_PROJECT_NAME="billing"
        echo "$TAG Returned to systems directory: $(pwd)"
    fi
    
    SYSTEM_OBJECT=$(echo "$METADATA" | jq -c --arg SYSTEM "$SYSTEM" '.[$SYSTEM]')
    export COMPOSE_PROJECT_NAME=$(echo "$SYSTEM_OBJECT" | jq -r '.key')
    run_docker_compose "$(echo "$SYSTEM_OBJECT" | jq -r '.path')" "$(echo "$SYSTEM_OBJECT" | jq -r '.name')" "$(echo "$SYSTEM_OBJECT" | jq -r '.key')"
    case $SYSTEM in
    "auth-services")
        sleep 2
        register_user
        ;;
    "console")
        cp .env.dev.example .env.local
        echo ".env.local file created and content copied from .env.dev.example for console"
        ;;
     "bff")
        cp .env.dev.example .env
        echo ".env file created and content copied from .env.dev.example for bff"
        ;;
    "init")
        sleep 2
        echo  "$TAG Add org in lookup..."
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5401/lookup"
        QUERY="INSERT INTO \"public\".\"orgs\" (\"created_at\", \"updated_at\", \"name\", \"org_id\", \"certificate\") VALUES (NOW(), NOW(), '$ORGNAME', '$ORGID', 'ukama-cert')"
        psql $DB_URI -c "$QUERY"
        ;;
    "dataplan")
        sleep 2
        echo  "$TAG Add default baserate in dataplan..."
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5404/baserate"
        QUERY1="INSERT INTO "base_rates" ("created_at","updated_at","deleted_at","uuid","country","provider","vpmn","imsi","sms_mo","sms_mt","data","x2g","x3g","x5g","lte","lte_m","apn","effective_at","end_at","sim_type","currency") VALUES ('2024-05-22 17:53:30.57','2024-05-22 17:53:30.57',NULL,'dd153d7f-d4aa-45e0-9e6a-0cc6407015ca','CD','OWS Tel','TTC',1,0,0,0,true,true,false,true,false,'Manual entry required','2024-06-10 00:00:00','2026-02-10 00:00:00',2,'Dollar')"
        psql $DB_URI -c "$QUERY1"

        echo  "$TAG Set default markup..."
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5404/rate"
        QUERY2="INSERT INTO "default_markups" ("created_at","updated_at","deleted_at","markup") VALUES ('2024-05-22 17:51:33.322','2024-05-22 17:51:33.322',NULL,1)"
        psql $DB_URI -c "$QUERY2"
    esac
    cd ../
done

sleep 3

if [ "$IS_INVENTORY_SYS" = true ]; then
    echo "$TAG Syncing up org inventory..."
    components=$(curl --location --silent --request PUT "http://${LOCAL_HOST_IP}:8077/v1/components/sync")
    echo "$TAG Org inventory synced up."
fi

# Update system url in lookup db
sleep 5

if [ "$IS_INCLUDE_BFF" = true ]; then
    SYS_QUERY_1="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-registry:8080' WHERE systems."name" = 'registry'";
    SYS_QUERY_2="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-notification:8080' WHERE systems."name" = 'notification'";
    SYS_QUERY_3="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-nucleus:8080' WHERE systems."name" = 'nucleus'";
    SYS_QUERY_4="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-subscriber:8080' WHERE systems."name" = 'subscriber'";
    SYS_QUERY_5="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-dataplan:8080' WHERE systems."name" = 'dataplan'";
    SYS_QUERY_6="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-inventory:8080' WHERE systems."name" = 'inventory'";
    SYS_QUERY_7="UPDATE PUBLIC.systems SET api_gw_url = 'http://subscriber-auth:4423' WHERE systems."name" = 'subscriber-auth'";
    SYS_QUERY_8="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-node:8080' WHERE systems."name" = 'node'";
    SYS_QUERY_9="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-metrics:8080' WHERE systems."name" = 'metrics'";
    SYS_QUERY_10="UPDATE PUBLIC.systems SET api_gw_url = 'http://report-api-gateway:8080' WHERE systems."name" = 'report';"
    SYS_QUERY_11="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-ukama-agent:8080' WHERE systems."name" = 'ukamaagent'";
    SYS_QUERY_12="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-dummy:8080' WHERE systems."name" = 'dummy'";
    SYS_QUERY_13="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-messaging:8080' WHERE systems."name" = 'messaging'";
    SYS_QUERY_14="UPDATE PUBLIC.systems SET api_gw_url = 'http://api-gateway-hub:8080' WHERE systems."name" = 'hub'";
fi
if [ "$IS_INCLUDE_BFF" = false ]; then
    SYS_QUERY_1="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8075' WHERE systems."name" = 'registry'";
    SYS_QUERY_2="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8058' WHERE systems."name" = 'notification'";
    SYS_QUERY_3="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8060' WHERE systems."name" = 'nucleus'";
    SYS_QUERY_4="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8078' WHERE systems."name" = 'subscriber'";
    SYS_QUERY_5="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8074' WHERE systems."name" = 'dataplan'";
    SYS_QUERY_6="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8077' WHERE systems."name" = 'inventory'";
    SYS_QUERY_7="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:4423' WHERE systems."name" = 'subscriber-auth'";
    SYS_QUERY_8="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8097' WHERE systems."name" = 'node'";
    SYS_QUERY_9="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8067' WHERE systems."name" = 'metrics'";
    SYS_QUERY_10="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8079' WHERE systems."name" = 'report'";
    SYS_QUERY_11="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8073' WHERE systems."name" = 'ukamaagent'";
    SYS_QUERY_12="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8086' WHERE systems."name" = 'dummy'";
    SYS_QUERY_13="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8079' WHERE systems."name" = 'messaging'";
    SYS_QUERY_14="UPDATE PUBLIC.systems SET api_gw_url = 'http://localhost:8000' WHERE systems."name" = 'hub'";
fi

echo "$TAG Registering systems URL in lookup db..."
DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5401/lookup"
psql $DB_URI -c "$SYS_QUERY_1"
psql $DB_URI -c "$SYS_QUERY_2"
psql $DB_URI -c "$SYS_QUERY_3"
psql $DB_URI -c "$SYS_QUERY_4"
psql $DB_URI -c "$SYS_QUERY_5"
psql $DB_URI -c "$SYS_QUERY_6"
psql $DB_URI -c "$SYS_QUERY_7"
psql $DB_URI -c "$SYS_QUERY_8"
psql $DB_URI -c "$SYS_QUERY_9"
psql $DB_URI -c "$SYS_QUERY_10"
psql $DB_URI -c "$SYS_QUERY_11"
psql $DB_URI -c "$SYS_QUERY_12"
psql $DB_URI -c "$SYS_QUERY_13"
psql $DB_URI -c "$SYS_QUERY_14"

cleanup

echo "$TAG Task done."