#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[38;5;39m'
NC='\033[0m' # No Color
TAG="${BLUE}Ukama>${NC}"
# Arg one: email
OWNEREMAIL=$1 
# Arg two: name
OWNERNAME=$2
# Arg three: ORG NAME
ORGNAME=$3
# Arg four: ORG ID
ORGID=$4
DOCKER_PROJECT_NAMES=("ukama-auth" "services" "init" "nucleus")

# Check if EMAIL or NAME is empty

# Check if EMAIL or NAME is empty
if [ -z "$OWNEREMAIL" ] || [ -z "$OWNERNAME" ]; then
    echo -e "Error: Both EMAIL and NAME must be provided"
    exit 1
fi

function run_docker_compose() {
    echo -e "$TAG Running $2 docker compose..."
    cd $1
    docker compose up --build -d > /dev/null 2>&1
    echo -e "$TAG $2 docker container is up"
}

function register_user() {
    echo -e "$TAG Signing up Owner user"

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
    OWNER_AUTH_ID=$(echo $identity | jq -r '.id')
    echo -e "$TAG User register with ${GREEN}AUTH ID = $OWNER_AUTH_ID${NC}"
    echo -e "$TAG Please verify your email address by visiting ${GREEN}http://localhost:4436${NC}"
}

# Main
if [[ "$(uname)" == "Darwin" ]]; then
    # For Mac
    export LOCAL_HOST_IP=$(ifconfig en0 | grep inet | awk '$1=="inet" {print $2}')
elif [[ "$(uname)" == "Linux" ]]; then
    # For Linux
    export LOCAL_HOST_IP=$(ifconfig enp0s25 | grep inet | awk '$1=="inet" {print $2}')
fi

# Run docker compose for ukama-auth
export COMPOSE_PROJECT_NAME="ukama-auth"
run_docker_compose "../ukama-auth" "Auth"
register_user

## Add env variables
export OWNERID=$(uuidgen)
export OWNERAUTHID=$OWNER_AUTH_ID
export OWNERNAME=$OWNERNAME
export OWNEREMAIL=$OWNEREMAIL
export ORGID=$ORGID
export ORGNAME=$ORGNAME

sleep 3

## Run docker compose for ukama/system/services
export COMPOSE_PROJECT_NAME="services"
run_docker_compose "../ukama/systems/services" "Ukama Services"

sleep 3

## Run docker compose for ukama/system/init
export COMPOSE_PROJECT_NAME="init"
run_docker_compose "../init" "Ukama Init system"

sleep 3

## Connect to init-lookup db and add org in orgs table
echo "Add org in lookup..."
DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5401/lookup"
QUERY="INSERT INTO \"public\".\"orgs\" (\"created_at\", \"updated_at\", \"name\", \"org_id\", \"certificate\") VALUES (NOW(), NOW(), '$ORGNAME', '$ORGID', 'ukama-cert')"
psql $DB_URI -c "$QUERY"

sleep 3

## Run docker compose for ukama/system/nucleus
export COMPOSE_PROJECT_NAME="nucleus"
run_docker_compose "../nucleus" "Ukama Nucleus system"

# Check Docker container status
echo -e "$TAG Checking Docker container status..."
while true; do
    for PROJECT_NAME in "${DOCKER_PROJECT_NAMES[@]}"; do
        DOCKER_PS_OUTPUT=$(docker ps -a --no-trunc --filter "label=com.docker.compose.project=$PROJECT_NAME" --format "{{json . }}" | jq -c 'if (.Networks | test(":")) then (.Networks |= fromjson) else . end | {ID: .ID, Names: .Names, Networks: .Networks, Status: .Status}')

        echo ""
        echo "System: $PROJECT_NAME"
        DOCKER_PS_OUTPUT_NEW=()
        while IFS= read -r line; do
            ID=$(echo "$line" | jq -r '.ID')
            LIVE_STATUS=$(docker inspect -f '{{.State.Status}}' "$ID")
            line=$(echo "$line" | jq --arg LIVE_STATUS "$LIVE_STATUS" '.Status = $LIVE_STATUS')
            DOCKER_PS_OUTPUT_NEW+=("$line")
        done <<< "$DOCKER_PS_OUTPUT"
        NEW_OBJECT=$(echo "${DOCKER_PS_OUTPUT_NEW[@]}" | jq -s 'reduce .[] as $item ({}; .ID += [$item.ID] | .Names += [$item.Names] | .Networks += [$item.Networks] | .Status += [$item.Status])')
        echo -e "${BLUE}Names\t\t\t\tStatus${NC}"
        NAMES=($(jq -r '.Names[]' <<< "$NEW_OBJECT"))
        STATUSES=($(jq -r '.Status[]' <<< "$NEW_OBJECT"))
        INDEX=-1
        for i in ${NAMES[@]}; do
            INDEX=${INDEX}+1
            case "${STATUSES[INDEX]}" in
                "running") COLOR=$GREEN;;
                "exited") COLOR=$RED;;
                *) COLOR=$YELLOW;;
            esac
            printf "%-30s %b%s%b\n" "${NAMES[INDEX]}" "$COLOR" "${STATUSES[INDEX]}" "$NC"
        done
    done
    sleep 5
    clear
done
