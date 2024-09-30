#!/bin/bash
# To run "chmod +x make-sys-for-mac.sh && ./make-sys-for-mac.sh ../deploy_config.json"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[38;5;39m'
NC='\033[0m'
TAG="${BLUE}Ukama>${NC}"
root_dir=$(pwd)
METADATA=$(jq -c '.' ../metadata.json)
SYS=$(jq -r '.systems' "$1")

declare -a PATHS
filter_make_sys() {
    IFS=', ' read -r -a SYSTEMS_ARRAY <<< "$SYS"

    for key in "${SYSTEMS_ARRAY[@]}"; do
        case "$key" in
            "nucleus")
                PATHS+=("nucleus/org" "nucleus/user" "nucleus/api-gateway")
                ;;
            "registry")
                PATHS+=("registry/network" "registry/site" "registry/member" "registry/invitation" "registry/node" "registry/api-gateway")
                ;;
            "subscriber")
                PATHS+=("subscriber/registry" "subscriber/sim-manager" "subscriber/sim-pool" "subscriber/api-gateway")
                ;;
            "dataplan")
                PATHS+=("data-plan/base-rate" "data-plan/package" "data-plan/rate" "data-plan/api-gateway")
                ;;
            "notification")
                PATHS+=("notification/distributor" "notification/event-notify" "notification/mailer" "notification/api-gateway")
                ;;
            "node")
                PATHS+=("node/configurator" "node/controller" "node/health" "node/node-gateway" "node/software" "node/api-gateway" "node/notify")
                ;;
            "init")
                PATHS+=("init/lookup" "init/api-gateway" "init/node-gateway")
                ;;
            "inventory")
                PATHS+=("inventory/accounting" "inventory/component" "inventory/api-gateway")
                ;;
            "messaging")
                PATHS+=("messaging/nns" "messaging/node-feeder")
                ;;
            "services")
                PATHS+=("services/msgClient")
                ;;
        esac
    done
}

filter_make_sys

cd ../../systems
root_dir=$(pwd)

for path in "${PATHS[@]}"; do
    cd "$root_dir/$path" || { echo "Failed to change directory to $path"; exit 1; }
   
    go mod tidy && make lint && make

    if [ $? -eq 0 ]; then
        echo "Make completed successfully in $path"
    else
        echo "Make failed in $path"
    fi

    cd - >/dev/null
done
