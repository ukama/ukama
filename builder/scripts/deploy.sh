# How to run example
# ./deploy.sh ../deploy_config.json

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
OWNEREMAIL=$(jq -r '.setup.email' "$1")
OWNERNAME=$(jq -r '.setup.name' "$1")
ORGNAME=$(jq -r '.setup["org-name"]' "$1")
ORGID=$(jq -r '.setup["org-id"]' "$1")
SYSTEMS=$(jq -r '.systems' "$1")
KEY=$(jq -r '.key' "$1")
METADATA=$(jq -c '.' ../metadata.json)
MAILERHOST=$(jq -r '.mailer.host' "$1")
MAILERPORT=$(jq -r '.mailer.port' "$1")
MAILERUSERNAME=$(jq -r '.mailer.username' "$1")
MAILERPASSWORD=$(jq -r '.mailer.password' "$1")
LAGOAPIKEY=$(jq -r '."lago-api-key"' "$1")
OWNERAUTHID=$(uuidgen)

function set_env() {
    export OWNERID=$(uuidgen)
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
    if [[ "$(uname)" == "Darwin" ]]; then
        # For Mac
        export LOCAL_HOST_IP=$(ifconfig en0 | grep inet | awk '$1=="inet" {print $2}')
    elif [[ "$(uname)" == "Linux" ]]; then
        # For Linux
        export LOCAL_HOST_IP=$(ifconfig enp0s25 | grep inet | awk '$1=="inet" {print $2}')
    fi
}

function run_docker_compose() {
    set_env
    echo  "$TAG Running $2 docker compose..."
    cd $1
    # docker-compose down  > /dev/null 2>&1
    # docker-compose build > /dev/null 2>&1
    # docker-compose up -d   > /dev/null 2>&1
    docker compose up --build -d > /dev/null 2>&1
    echo  "$TAG $2 docker container is up"
}

function register_user() {
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

if [ $# q 0 ]; then
    echo "No arguments provided. Please provide the path to the deploy_config JSON file."
    exit 1
fi

IFS=', ' read -r -a SYSTEMS_ARRAY <<< "$SYSTEMS"


if [[ $SYSTEMS != *"auth-services"* ]]; then
    echo "Please add auth-services in the systems array in the deploy_config JSON file."
else
    auth=$(echo "$METADATA" | jq -c --arg AUTHSYSKEY "$AUTHSYSKEY" '.[$AUTHSYSKEY]')
    # Run docker compose for ukama-auth
    export COMPOSE_PROJECT_NAME="$(echo "$auth" | jq -r '.key')"
    run_docker_compose "$(echo "$auth" | jq -r '.path')" "$(echo "$auth" | jq -r '.name')"
    cd ../ukama/builder/scripts
    register_user
fi

#Navigate to Ukama repo
cd ../../systems

# Loop through the SYSTEMS array
for SYSTEM in "${SYSTEMS_ARRAY[@]}"; do
    if [ "$SYSTEM" == $AUTHSYSKEY ]; then
        continue
    fi
    SYSTEM_OBJECT=$(echo "$METADATA" | jq -c --arg SYSTEM "$SYSTEM" '.[$SYSTEM]')
    export COMPOSE_PROJECT_NAME=$(echo "$SYSTEM_OBJECT" | jq -r '.key')
    run_docker_compose "$(echo "$SYSTEM_OBJECT" | jq -r '.path')" "$(echo "$SYSTEM_OBJECT" | jq -r '.name')"
    case $SYSTEM in
    "init")
        sleep 2
        ## Connect to init-lookup db and add org in orgs table
        echo  "$TAG Add org in lookup..."
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5401/lookup"
        QUERY="INSERT INTO \"public\".\"orgs\" (\"created_at\", \"updated_at\", \"name\", \"org_id\", \"certificate\") VALUES (NOW(), NOW(), '$ORGNAME', '$ORGID', 'ukama-cert')"
        psql $DB_URI -c "$QUERY"
        ;;
    esac
    cd ../
done