#!/bin/bash

source data_seeder_test.sh

DB_CONTAINER_NAME="postgresd"
DB_USERNAME="postgres"
DB_PASSWORD="ThisIsUkamaXPass"
PROJECT_NAME="ukama-systems"

SIMS_CN="users"
REGISTRY_CN="registry"
SIMMANAGER_CN="simmanager"
TESTAGENT_CN="testagent"
BASERATE_CN="baserate"
PACKAGE_CN="package"
RATE_CN="rate"
NETWORK_CN="network"
ORG_CN="org"
NODE_CN="node"
USERS_CN="users"
DATA_PLAN_CN="data-plan"
SUBSCRIBER_CN="subscriber"
AUTH_CN="auth"
MSG_CLIENT_SUBSCRIBER_CN="msg-client-subscriber"
MSG_CLIENT_DATA_PLAN_CN="msg-client-dataplan"
MSG_CLIENT_REGISTRY_CN="msg-client-registry"

declare -a CONTAINER_NAMES=(
    "users"
    "registry"
    "simmanager"
    "testagent"
    "baserate"
    "package"
    "rate"
    "network"
    "org"
    "node"
    "users"
    "data-plan"
    "subscriber"
    "auth"
    "msg-client-subscriber"
    "msg-client-dataplan"
    "msg-client-registry"
)

declare -a TABLE_NAMES=(
    "users.users"
    "org.org_users"
    "org.orgs"
    "org.users"
    "network.sites"
    "network.networks"
    "network.orgs"
    "network.sites"
    "node.attached_nodes"
    "node.nodes"
    "package.package_details"
    "package.package_markups"
    "package.package_rates"
    "package.packages"
    "baserate.base_rates"
    "rate.default_markups"
    "rate.markups"
    "simmanager.packages"
    "simmanager.sims"
    "sim.sims"
    "msgclient.service_routes"
    "msgclient.services"
    "msgclient.routes"
)

if docker ps -a --format "{{.Names}}" | grep -q "^${ORG_CN}\$"; then
    if docker ps -a --format "{{.Names}}" | grep -q "^${DB_CONTAINER_NAME}\$"; then
        echo "Container '${DB_CONTAINER_NAME}' found."
        echo "Cleaning DB..."

        for table_name in "${TABLE_NAMES[@]}"
        do
            IFS='.' read -ra parts <<< "$table_name"
            db=${parts[0]}
            table=${parts[1]}
            docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $db -c "DELETE FROM $table;"
        done
    else
        echo "Container '${DB_CONTAINER_NAME}' does not exist."
    fi
fi

docker compose up --build -d

echo "Project $PROJECT_NAME is up and running"
echo "Connecting with Postgres container..."
sleep 10

echo "Inserting data in db's..."

echo "Inserting data into Users DB..."
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $USERS_DB -c "$USER_QUERY"

echo "Inserting data into Org DB..."
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $ORG_DB -c "$ORG_QUERY"
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $ORG_DB -c "$USERS_IN_ORG_QUERY"
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $ORG_DB -c "$ORG_USERS_QUERY"

echo "Inserting data into Networks DB..."
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $NETWORK_DB -c "$NETWORK_ORGS_QUERY"
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $NETWORK_DB -c "$NETWORKS_QUERY"
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $NETWORK_DB -c "$SITES_QUERY"

echo "Inserting data into Baserate DB..."
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $BASERATE_DB -c "$BASERATE_QUERY"

echo "Inserting data into Markup DB..."
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $RATE_DB -c "$MARKUP_DEFAULT_QUERY"
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $RATE_DB -c "$MARKUPS_QUERY"

echo "Inserting data into Package DB..."
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $PACKAGE_DB -c "$PACKAGES_QUERY"
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $PACKAGE_DB -c "$PACKAGE_DETAILS_QUERY"
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $PACKAGE_DB -c "$PACKAGE_MARKUPS_QUERY"
docker exec -it $DB_CONTAINER_NAME psql -U $DB_USERNAME -d $PACKAGE_DB -c "$PACKAGE_RATES_QUERY"