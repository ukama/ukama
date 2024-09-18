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
ORG_COMMUNITY="community"
# Parse the JSON file and initialize the variables
IS_INIT_SYS="init"
BILLINGSYSKEY="billing"
AUTHSYSKEY="auth-services"
INVENTORY_SYS_KEY="inventory"
IS_INVENTORY_SYS=false
METADATA=$(jq -c '.' ../metadata.json)
JSON_FILE="../deploy_orgs_config.json"
if [ "$1" == "-d" ]; then
  ISDEBUGMODE=true
else
  ISDEBUGMODE=false
fi

if [[ "$(uname)" == "Darwin" ]]; then
    # For Mac
    LOCAL_HOST_IP=$(ifconfig en0 | grep inet | awk '$1=="inet" {print $2}')
elif [[ "$(uname)" == "Linux" ]]; then
    # For Linux
    LOCAL_HOST_IP=$(ifconfig enp0s25 | grep inet | awk '$1=="inet" {print $2}')
fi

MASTERORGNAME=$(jq -r '.["master-org-name"]' "$JSON_FILE")

jq -c '.orgs[]' "$JSON_FILE" | while read -r ORG; do
    ORGNAME=$(echo "$ORG" | jq -r '.["org-name"]')
    echo "${TAG} Processing organization: ${YELLOW}${ORGNAME}${NC}"

    # Initialize the variables
    SUBNET=$(echo "$ORG" | jq -r '.subnet')
    ORG_TYPE=$(echo "$ORG" | jq -r '.type')
    OWNEREMAIL=$(echo "$ORG" | jq -r '.email')
    PASSWORD=$(echo "$ORG" | jq -r '.password')
    OWNERNAME=$(echo "$ORG" | jq -r '.name')
    ORGID=$(echo "$ORG" | jq -r '.["org-id"]')
    KEY=$(echo "$ORG" | jq -r '.key')
    MAILERHOST=$(echo "$ORG" | jq -r '.mailer.host')
    MAILERPORT=$(echo "$ORG" | jq -r '.mailer.port')
    MAILERUSERNAME=$(echo "$ORG" | jq -r '.mailer.username')
    MAILERPASSWORD=$(echo "$ORG" | jq -r '.mailer.password')
    LAGOAPIKEY=$(echo "$ORG" | jq -r '.["lago-api-key"]')
    WITHSUBAUTH=$(echo "$ORG" | jq -r '.["with-subscriber-auth"]')
    SYS=$(echo "$ORG" | jq -r '.systems')
    OWNERAUTHID=""
    OWNERID=$(uuidgen)
    SYSTEMS=()

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
        export COMPONENT_ENVIRONMENT=test
    }

    function run_docker_compose() {
        set_env
        echo  "$TAG Running $2 docker compose..."
        cd $1
        CONTAINER_NAME=$3
        while true; do
            if [ "$ISDEBUGMODE" = true ]; then
                docker compose -p $ORGNAME up --build -d
            else
                docker compose -p $ORGNAME up --build -d > /dev/null 2>&1
            fi
            sleep 5
            break
            # Need to figure out a way to verify the container status
            # if docker ps | grep -q $CONTAINER_NAME; then
            #     echo  "$TAG $2 docker container is up"
            # else
            #     echo "Container $CONTAINER_NAME is not running. Retrying..."
            # fi
        done
    }

    function register_user() {
        echo  "$TAG Signing up Owner for org ${GREEN}$ORGNAME${NC}"
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

    register_org_to_init(){
        echo  "$TAG Add ${ORGNAME} org in lookup..."
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5401/lookup"
        QUERY="INSERT INTO \"public\".\"orgs\" (\"created_at\", \"updated_at\", \"name\", \"org_id\", \"certificate\",  \"country\") VALUES (NOW(), NOW(), '$ORGNAME', '$ORGID', 'ukama-cert', 'CD')"
        psql $DB_URI -c "$QUERY"
    }

    get_user_id() {
        echo  "$TAG Fetching user ID..."
        SQL_QUERY="SELECT PUBLIC.users.id FROM PUBLIC.users WHERE users.auth_id = '$OWNERAUTHID' AND users.deleted_at IS NULL ORDER BY users.id LIMIT 1;"
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5411/users"
        OWNERID=$(psql $DB_URI -t -A -c "$SQL_QUERY")
        echo "$TAG User ID: ${GREEN} $OWNERID ${NC}"
    }

    create_org() {
       echo  "$TAG Register org in nucleus..."
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5411/org"
        QUERY="INSERT INTO \"public\".\"orgs\" (\"created_at\", \"updated_at\", \"name\", \"owner\", \"certificate\", \"id\", \"deactivated\") VALUES (NOW(), NOW(), '$ORGNAME', '$OWNERID', 'ukama-cert', '$ORGID', FALSE)"
        psql $DB_URI -c "$QUERY"
    }

    sort_systems_by_dependency() {
        IFS=', ' read -r -a SYSTEMS_ARRAY <<< "$SYS"
        
        if [[ "$ORG_TYPE" == "$ORG_COMMUNITY" ]] && 
        ([[ ! " ${SYSTEMS_ARRAY[@]} " =~ " auth-services " ]] || 
            [[ ! " ${SYSTEMS_ARRAY[@]} " =~ " init " ]] || 
            [[ ! " ${SYSTEMS_ARRAY[@]} " =~ " nucleus " ]]); then
            echo "Error: ${ORGNAME} Required systems are missing for this type of org."
            exit 1
        fi
        
        if [[ ! "$ORG_TYPE" == "$ORG_COMMUNITY" ]]; then
            invalid_system=false
            for system in "auth-services" "init" "nucleus" "inventory" "bff" "console"; do
                if [[ " ${SYSTEMS_ARRAY[@]} " =~ " ${system} " ]]; then
                    invalid_system=true
                    break
                fi
            done

            if [[ "${invalid_system}" == true ]]; then
                echo "Error: Invalid system for empowerment type org."
                exit 1
            fi
        fi

        if [[ ! " ${SYSTEMS_ARRAY[@]} " =~ " services " ]]; then
            echo "Error: 'services' required in the systems array. Please modify systems in multi_org_config JSON file."
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
                "subscriber")
                    SYSTEMS+=("8 $key")
                    ;;
                *)
                    SYSTEMS+=("9 $key")
                    ;;
            esac
        done

        SYSTEMS=($(for item in "${SYSTEMS[@]}"; do echo "$item"; done | sort -n -k1,1 | cut -d' ' -f2-))
    }

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
            sed -i '' "s/services_/${ORGNAME}_/g" docker-compose.yml
            sed -i '' "s/10.1.0.0/${SUBNET}/g" docker-compose.yml
            if [[ "$(uname)" == "Darwin" ]]; then
                sed -i '' "s/build: \.\.\/services\/initClient/image: main-init/g" docker-compose.yml
            fi
            
            if [[ ! "$ORG_TYPE" == "$ORG_COMMUNITY" ]]; then
                sed -i '' '/ports:/d' docker-compose.yml
                sed -i '' '/- 8090:80/d' docker-compose.yml
                sed -i '' '/- 5672:5672/d' docker-compose.yml
                sed -i '' '/- 15672:15672/d' docker-compose.yml
                sed -i '' '/- 8075:8080/d' docker-compose.yml
                sed -i '' '/- 8036:8080/d' docker-compose.yml
                sed -i '' '/- 8097:8080/d' docker-compose.yml
                sed -i '' '/- 8058:8080/d' docker-compose.yml
                sed -i '' '/- 8078:8080/d' docker-compose.yml
                sed -i '' '/- 8074:8080/d' docker-compose.yml
                sed -i '' '/- 5405:5432/d' docker-compose.yml # REGISTRY SYS
                sed -i '' '/- 5489:5432/d' docker-compose.yml # NODE SYS
                sed -i '' '/- 5632:5432/d' docker-compose.yml # NOTIFICATION SYS
                sed -i '' '/- 5412:5432/d' docker-compose.yml # SUBSCRIBER SYS
                sed -i '' '/- 5404:5432/d' docker-compose.yml # DATAPLAN SYS
            fi
            if [[ $WITHSUBAUTH == false ]]; then
                sed -i '' '/- 4446:4446/d' docker-compose.yml # SUBSCRIBER MAILSERVER
                sed -i '' '/- 4447:4447/d' docker-compose.yml # SUBSCRIBER MAILSERVER
                sed -i '' '/- 4423:4423/d' docker-compose.yml # SUBSCRIBER AUTH
                sed -i '' '/- 4424:4424/d' docker-compose.yml # SUBSCRIBER AUTH
            fi
        done
        cd $root_dir
    }

    pre_deploy_config_for_other_org(){
        if [[ ! "$ORG_TYPE" == "$ORG_COMMUNITY" ]]; then
            echo "$TAG Preparing deploy config for empowerment org..."
            register_org_to_init
            sleep 3
            register_user
            sleep 5
            get_user_id
            sleep 5
            create_org
        fi
    }

    sort_systems_by_dependency
    setup_docker_compose_files
    pre_deploy_config_for_other_org

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
        if [ "$SYSTEM" != $AUTHSYSKEY ]; then
            cd ../../systems
        fi
        if [ "$SYSTEM" == $BILLINGSYSKEY ]; then
            cd ./billing/provider
            chmod +x start_provider.sh
            ./start_provider.sh
            cd ../..
        fi
        if [[ " ${SYSTEM} " == " bff " ]]; then
           cp .env.dev.example .env
               echo ".env file created and content copied from .env.dev.example for bff"
        fi
        if [[ " ${SYSTEM} " == " console " ]]; then
           cp .env.dev.example .env.local
           echo ".env.local file created and content copied from .env.dev.example for console"
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
            register_org_to_init
            ;;

        "dataplan")
            # TODO: NEED TO BE FIXED
            # sleep 2
            # echo  "$TAG Add default baserate in dataplan..."
            # DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5404/baserate"
            # QUERY1="INSERT INTO "base_rates" ("created_at","updated_at","deleted_at","uuid","country","provider","vpmn","imsi","sms_mo","sms_mt","data","x2g","x3g","x5g","lte","lte_m","apn","effective_at","end_at","sim_type","currency") VALUES ('2024-05-22 17:53:30.57','2024-05-22 17:53:30.57',NULL,'dd153d7f-d4aa-45e0-9e6a-0cc6407015ca','CD','OWS Tel','TTC',1,0,0,0,true,true,false,true,false,'Manual entry required','2024-06-10 00:00:00','2026-02-10 00:00:00',2,'Dollar')"
            # psql $DB_URI -c "$QUERY1"

            # echo  "$TAG Set default markup..."
            # DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5404/rate"
            # QUERY2="INSERT INTO "default_markups" ("created_at","updated_at","deleted_at","markup") VALUES ('2024-05-22 17:51:33.322','2024-05-22 17:51:33.322',NULL,1)"
            # psql $DB_URI -c "$QUERY2"
        esac
        cd ../
    done

    cleanup

    if [[ "${ORG_TYPE}" =~ "${ORG_COMMUNITY}" ]]; then
        sleep 3
        if ($IS_INVENTORY_SYS); then
            echo "$TAG Syncing up org inventory..."
            components=$(curl --location --silent --request PUT "http://${LOCAL_HOST_IP}:8077/v1/components/sync")
            echo "$TAG Org inventory synced up."
        fi

        sleep 5

        SYS_QUERY_1="UPDATE PUBLIC.systems SET url = 'http://api-gateway-registry:8080' WHERE systems."name" = 'registry'";
        SYS_QUERY_2="UPDATE PUBLIC.systems SET url = 'http://api-gateway-notification:8080' WHERE systems."name" = 'notification'";
        SYS_QUERY_3="UPDATE PUBLIC.systems SET url = 'http://api-gateway-nucleus:8080' WHERE systems."name" = 'nucleus'";
        SYS_QUERY_4="UPDATE PUBLIC.systems SET url = 'http://api-gateway-subscriber:8080' WHERE systems."name" = 'subscriber'";
        SYS_QUERY_5="UPDATE PUBLIC.systems SET url = 'http://api-gateway-dataplan:8080' WHERE systems."name" = 'dataplan'";
        SYS_QUERY_6="UPDATE PUBLIC.systems SET url = 'http://api-gateway-inventory:8080' WHERE systems."name" = 'inventory'";
        SYS_QUERY_7="UPDATE PUBLIC.systems SET url = 'http://subscriber-auth:4423' WHERE systems."name" = 'subscriber-auth'";
        SYS_QUERY_8="UPDATE PUBLIC.systems SET url = 'http://api-gateway-node:8080' WHERE systems."name" = 'node'";

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
    fi

    if [[ ! "$ORG_TYPE" == "$ORG_COMMUNITY" ]]; then
        echo "${TAG}Connect global services with containers...${NC}"
        if [[ " ${SYSTEMS[@]} " =~ " registry " ]]; then
            docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-member-1
            docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-network-1
            docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-invitation-1
            docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-site-1
        fi
        if [[ " ${SYSTEMS[@]} " =~ " subscriber " ]]; then
            docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-registry-1
            docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-simmanager-1
        fi
        if [[ " ${SYSTEMS[@]} " =~ " notification " ]]; then
            docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-eventnotify-1
            docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-distributor-1
        fi
        if [[ " ${SYSTEMS[@]} " =~ " node " ]]; then
            docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-controller-1
            docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-configurator-1
        fi
        docker network connect ${ORGNAME}_ukama-net ${MASTERORGNAME}-bff-1

        echo  "$TAG Updateing inventory..."
        SQL_QUERY="UPDATE PUBLIC.components SET user_id = '$OWNERID';"
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5414/component"
        psql $DB_URI -t -A -c "$SQL_QUERY"
    fi
done

echo "${TAG} Task Done${NC}"