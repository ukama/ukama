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
ISDEBUGMODE=false
MASTERORGNAME=$(jq -r '.["master-org-name"]' "$JSON_FILE")


if [[ "$(uname)" == "Darwin" ]]; then
    # For Mac
    LOCAL_HOST_IP=$(ifconfig en0 | grep inet | awk '$1=="inet" {print $2}')
elif [[ "$(uname)" == "Linux" ]]; then
    # For Linux
    LOCAL_HOST_IP=$(ifconfig enp0s25 | grep inet | awk '$1=="inet" {print $2}')
fi

function buildSystems() {
    echo  "$TAG Building systems..."
    ./make-sys-for-mac.sh ../deploy_config.json 2>&1 | tee buildSystems.log
}

while getopts "bd" opt; do
    case ${opt} in
        b )
            buildSystems
            ;;
        d )
            ISDEBUGMODE=true
            ;;
        \? )
            echo "Usage: cmd [-b, -d]"
            exit 1
            ;;
    esac
done

jq -c '.orgs[]' "$JSON_FILE" | while read -r ORG; do
    ORGNAME=$(echo "$ORG" | jq -r '.["org-name"]')
    echo "${TAG} Processing organization: ${YELLOW}${ORGNAME}${NC}"
    docker network rm ${ORGNAME}_ukama-net > /dev/null 2>&1
    # Initialize the variables
    SUBNET=$(echo "$ORG" | jq -r '.subnet')
    ORG_TYPE=$(echo "$ORG" | jq -r '.type')
    OWNEREMAIL=$(echo "$ORG" | jq -r '.email')
    PASSWORD=$(echo "$ORG" | jq -r '.password')
    OWNERNAME=$(echo "$ORG" | jq -r '.name')
    ORGID=$(echo "$ORG" | jq -r '.["org-id"]')
    KEY=$(echo "$ORG" | jq -r '.key')
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
        export LAGO_API_KEY=$LAGOAPIKEY
        export MASTERORGNAME=$MASTERORGNAME
        export MASTER_ORG_NAME=$MASTERORGNAME
        export LOCAL_HOST_IP=$LOCAL_HOST_IP
        export COMPONENT_ENVIRONMENT=test
        export PROMETHEUS_HTTP_URL="http://${LOCAL_HOST_IP}:${PROMETHEUS}/v1/prometheus"
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
        QUERY="INSERT INTO \"public\".\"orgs\" (\"created_at\", \"updated_at\", \"name\", \"org_id\", \"certificate\") VALUES (NOW(), NOW(), '$ORGNAME', '$ORGID', 'ukama-cert')"
        psql $DB_URI -c "$QUERY"
    }

    get_user_id() {
        echo  "$TAG Fetching user ID..."
        response=$(curl --location --silent "http://localhost:8060/v1/users/auth/${OWNERAUTHID}" -H 'accept: application/json')
        user_id=$(echo "$response" | jq -r '.user.id')
        OWNERID="$user_id"
        echo "$TAG User ID: ${GREEN} $OWNERID ${NC}"
    }

    create_org() {
       echo  "$TAG Register org in nucleus..."
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5406/org"
        QUERY="INSERT INTO \"public\".\"orgs\" (\"created_at\", \"updated_at\", \"name\", \"owner\", \"certificate\", \"id\", \"deactivated\", \"country\", \"currency\") VALUES (NOW(), NOW(), '$ORGNAME', '$OWNERID', 'ukama-cert', '$ORGID', FALSE, 'cod', 'cdf')"
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

    getSysPorts(){
        MET_PORT=$(echo "$ORG" | jq -r '.["sys-db-ports"]["metrics"]')
        DP_PORT=$(echo "$ORG" | jq -r '.["sys-db-ports"]["dataplan"]')
        SUB_PORT=$(echo "$ORG" | jq -r '.["sys-db-ports"]["subscriber"]')
        REG_PORT=$(echo "$ORG" | jq -r '.["sys-db-ports"]["registry"]')
        NOT_PORT=$(echo "$ORG" | jq -r '.["sys-db-ports"]["notification"]')
        NODE_PORT=$(echo "$ORG" | jq -r '.["sys-db-ports"]["node"]')

        PGA_PORT=$(echo "$ORG" | jq -r '.["sys-ports"]["pg-admin"]')
        RABBITMQ_P1=$(echo "$ORG" | jq -r '.["sys-ports"]["rbitmq-1"]')
        RABBITMQ_P2=$(echo "$ORG" | jq -r '.["sys-ports"]["rbitmq-2"]')
        PROMETHEUS=$(echo "$ORG" | jq -r '.["sys-ports"]["prometheus"]')
        REGAPI_PORT=$(echo "$ORG" | jq -r '.["sys-ports"]["registry"]')
        NODEAPI_PORT=$(echo "$ORG" | jq -r '.["sys-ports"]["node"]')
        NODEGW_PORT=$(echo "$ORG" | jq -r '.["sys-ports"]["nodegw"]')
        NOTAPI_PORT=$(echo "$ORG" | jq -r '.["sys-ports"]["notification"]')
        SUBAPI_PORT=$(echo "$ORG" | jq -r '.["sys-ports"]["subscriber"]')
        DPAPI_PORT=$(echo "$ORG" | jq -r '.["sys-ports"]["dataplan"]')
        METRICS_PORT=$(echo "$ORG" | jq -r '.["sys-ports"]["metrics"]')
        MSGCLIENT_PORT=$(echo "$ORG" | jq -r '.["sys-ports"]["msgclient"]')
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
                getSysPorts
                INITCLIENT_HOST=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${MASTERORGNAME}-api-gateway-init-1)
                NUCLEUSCLIENT_HOST=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${MASTERORGNAME}-api-gateway-nucleus-1)
                INVENTORYCLIENT_HOST=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${MASTERORGNAME}-api-gateway-inventory-1)

               if [ "$SYSTEM" == "services" ]; then
                    sed -i '' "s/- 8090:80/- ${PGA_PORT}:80/g" docker-compose.yml # PG ADMIN
                    sed -i '' "s/- 5672:5672/- ${RABBITMQ_P1}:5672/g" docker-compose.yml # RABBIT MQ P1
                    sed -i '' "s/- 15672:15672/- ${RABBITMQ_P2}:15672/g" docker-compose.yml # RABBIT MQ P2
                    continue
                fi

                sed -i '' "s/- 8075:8080/- ${REGAPI_PORT}:8080/g" docker-compose.yml # REGISTRY SYS APIGW
                sed -i '' "s/- 8036:8080/- ${NODEAPI_PORT}:8080/g" docker-compose.yml # NODE SYS APIGW
                sed -i '' "s/- 8097:8080/- ${NODEGW_PORT}:8080/g" docker-compose.yml # NODE SYS NODEGW
                sed -i '' "s/- 8058:8080/- ${NOTAPI_PORT}:8080/g" docker-compose.yml # NOTIFICATION SYS APIGW
                sed -i '' "s/- 8097:8080/- ${SUBAPI_PORT}:8080/g" docker-compose.yml # SUBSCRIBER SYS APIGW
                sed -i '' "s/- 8074:8080/- ${DPAPI_PORT}:8080/g" docker-compose.yml # DATAPLAN SYS APIGW
                sed -i '' "s/- 8067:8080/- ${METRICS_PORT}:8080/g" docker-compose.yml # METRICS SYS APIGW

                sed -i '' "s/- 5405:5432/- ${REG_PORT}:5432/g" docker-compose.yml # REGISTRY SYS PG
                sed -i '' "s/- 5489:5432/- ${NODE_PORT}:5432/g" docker-compose.yml # NODE SYS PG
                sed -i '' "s/- 5632:5432/- ${NOT_PORT}:5432/g" docker-compose.yml # NOTIFICATION SYS PG
                sed -i '' "s/- 5412:5432/- ${SUB_PORT}:5432/g" docker-compose.yml # SUBSCRIBER SYS PG
                sed -i '' "s/- 5404:5432/- ${DP_PORT}:5432/g" docker-compose.yml # DATAPLAN SYS PG
                sed -i '' "s/- 5407:5432/- ${MET_PORT}:5432/g" docker-compose.yml # METRICS SYS PG
                
                sed -i '' "s/api-gateway-init:8080/${INITCLIENT_HOST}:8080/g" docker-compose.yml
                sed -i '' "s/api-gateway-nucleus:8080/${NUCLEUSCLIENT_HOST}:8080/g" docker-compose.yml
                sed -i '' "s/api-gateway-inventory:8080/${INVENTORYCLIENT_HOST}:8080/g" docker-compose.yml

                sed -i '' "s/9095/${MSGCLIENT_PORT}/g" docker-compose.yml
                sed -i '' "s/5672/${RABBITMQ_P1}/g" docker-compose.yml
                sed -i '' "s/9079/${PROMETHEUS}/g" docker-compose.yml
            fi

            if [[ $WITHSUBAUTH == false ]]; then
                sed -i '' '/- 4446:4436/d' docker-compose.yml # SUBSCRIBER MAILSERVER
                sed -i '' '/- 4447:4437/d' docker-compose.yml # SUBSCRIBER MAILSERVER
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
            cd console-bff
            cp .env.dev.example .env
            echo ".env file created and content copied from .env.dev.example for bff"
            cd ..
        fi
        if [[ " ${SYSTEM} " == " console " ]]; then
            cd ../interfaces/console
            cp .env.dev.example .env.local
            echo ".env.local file created and content copied from .env.dev.example for console"
            cd ../../systems
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
            sleep 2
            echo  "$TAG Add default baserate in dataplan..."
            DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:${DP_PORT}/baserate"
            QUERY1="INSERT INTO "base_rates" ("created_at","updated_at","deleted_at","uuid","country","provider","vpmn","imsi","sms_mo","sms_mt","data","x2g","x3g","x5g","lte","lte_m","apn","effective_at","end_at","sim_type","currency") VALUES ('2024-05-22 17:53:30.57','2024-05-22 17:53:30.57',NULL,'dd153d7f-d4aa-45e0-9e6a-0cc6407015ca','CD','OWS Tel','TTC',1,0,0,0,true,true,false,true,false,'Manual entry required','2024-06-10 00:00:00','2026-02-10 00:00:00',2,'Dollar')"
            psql $DB_URI -c "$QUERY1"

            echo  "$TAG Set default markup..."
            DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:${DP_PORT}/rate"
            QUERY2="INSERT INTO "default_markups" ("created_at","updated_at","deleted_at","markup") VALUES ('2024-05-22 17:51:33.322','2024-05-22 17:51:33.322',NULL,1)"
            psql $DB_URI -c "$QUERY2"
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
        docker network connect ${ORGNAME}_ukama-net ${MASTERORGNAME}-org-1
        docker network connect ${MASTERORGNAME}_ukama-net ${ORGNAME}-subscriber-auth-1
        
        echo  "$TAG Updateing inventory for ${ORGNAME}..."
        SQL_QUERY="UPDATE PUBLIC.components SET user_id = '$OWNERID';"
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5414/component"
        psql $DB_URI -t -A -c "$SQL_QUERY"
        
        sleep 3
        SYS_QUERY_1="UPDATE PUBLIC.systems SET url = 'http://salman-org-subscriber-auth-1:4423' WHERE systems."name" = 'subscriber-auth'";
        DB_URI="postgresql://postgres:Pass2020!@127.0.0.1:5401/lookup"
        psql $DB_URI -c "$SYS_QUERY_1"
    fi
done

echo "${TAG} Task Done${NC}"