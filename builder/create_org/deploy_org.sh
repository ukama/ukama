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
# Parse the JSON file and initialize the variables
ORG_NAME=$(jq -r '."org-name"' "$1")
ORG_ID=$(jq -r '."org-id"' "$1")
OWNER_AUTH_ID=$(jq -r '."owner-auth-id"' "$1")
OWNER_USER_ID=$(jq -r '."owner-user-id"' "$1")
OWNER_USER_EMAIL=$(jq -r '."owner-user-email"' "$1")
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
root_dir=$(pwd)
echo $root_dir

function set_env() {
    export KEY=$KEY
    export ORG_ID=$ORG_ID
    export ORG_NAME=$ORG_NAME
    export MAILERHOST=$MAILERHOST
    export MAILERPORT=$MAILERPORT
    export LAGOAPIKEY=$LAGOAPIKEY   
    export OWNER_AUTH_ID=$OWNER_AUTH_ID
    export OWNER_USER_ID=$OWNER_USER_ID
    export MAILERUSERNAME=$MAILERUSERNAME
    export MAILERPASSWORD=$MAILERPASSWORD
    export OWNER_USER_EMAIL=$OWNER_USER_EMAIL
}

function run_docker_compose() {
    set_env
    echo  "$TAG Running $2 docker compose..."
    cd $1

    CONTAINER_NAME="$ORG_NAME-$3"
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
    cd ..
}

sort_systems_by_dependency() {
    IFS=', ' read -r -a SYSTEMS_ARRAY <<< "$SYS"

    SYSTEMS=($(for key in "${SYSTEMS_ARRAY[@]}"; do
      if [ "$key" == "services" ]; then
          echo "1 $key"
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

echo  "$TAG Add org in lookup..."
DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5401/lookup"
QUERY1="INSERT INTO \"public\".\"orgs\" (\"created_at\", \"updated_at\", \"name\", \"org_id\", \"certificate\") VALUES (NOW(), NOW(), '$ORG_NAME', '$ORG_ID', 'ukama-cert')"
psql $DB_URI -c "$QUERY1"
DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5411/org"
# QUERY2="INSERT INTO "users" ("uuid","deactivated","deleted_at","id") VALUES ('$OWNER_USER_ID',false,NULL,1) ON CONFLICT DO NOTHING"
# psql $DB_URI -c "$QUERY2"
QUERY4="INSERT INTO "orgs" ("id","name","owner","certificate","deactivated","created_at","updated_at","deleted_at") VALUES ('$ORG_ID','$ORG_NAME','$OWNER_USER_ID','',false,NOW(),NOW(),NULL)"
psql $DB_URI -c "$QUERY4"
# psql $DB_URI -c "$QUERY3"
QUERY3="INSERT INTO "org_users" ("org_id","user_id") VALUES ('$ORG_ID',2) ON CONFLICT DO NOTHING"
psql $DB_URI -c "$QUERY3"
cd ~
cd $root_dir

run_docker_compose "./services" "Services" "services"
run_docker_compose "./registry" "Registry" "registry"
run_docker_compose "./dataplan" "DataPlan" "dataplan"

 sleep 2
echo  "$TAG Add default baserate in dataplan..."
DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5404/baserate"
QUERY6="INSERT INTO "base_rates" ("created_at","updated_at","deleted_at","uuid","country","provider","vpmn","imsi","sms_mo","sms_mt","data","x2g","x3g","x5g","lte","lte_m","apn","effective_at","end_at","sim_type","currency") VALUES ('2024-05-22 17:53:30.57','2024-05-22 17:53:30.57',NULL,'dd153d7f-d4aa-45e0-9e6a-0cc6407015ca','CD','OWS Tel','TTC',1,0,0,0,true,true,false,true,false,'Manual entry required','2024-06-10 00:00:00','2026-02-10 00:00:00',2,'Dollar')"
psql $DB_URI -c "$QUERY6"

echo  "$TAG Set default markup..."
DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5404/rate"
QUERY7="INSERT INTO "default_markups" ("created_at","updated_at","deleted_at","markup") VALUES ('2024-05-22 17:51:33.322','2024-05-22 17:51:33.322',NULL,1)"
psql $DB_URI -c "$QUERY7"

run_docker_compose "./subscriber" "Subscriber" "subscriber"
run_docker_compose "./notification" "Notification" "notification"
