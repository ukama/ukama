#!/bin/bash
# To run "chmod +x make-sys-for-mac.sh && ./make-sys-for-mac.sh ../deploy_config.json"

TAG="\033[38;5;39mUkama>\033[0m"
YELLOW='\033[1;33m'
NC='\033[0m' # No Color
RED='\033[0;31m'
GREEN='\033[0;32m'

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
                PATHS+=("node/configurator" "node/controller" "node/health" "node/node-gateway" "node/software" "node/api-gateway" "node/notify" "node/state")
                ;;
            "billing")
                PATHS+=("node/configurator" "billing/report" "billing/api-gateway" "billing/collector" )
                ;;
            "init")
                PATHS+=("init/lookup" "init/api-gateway" "init/node-gateway")
                ;;
            "inventory")
                PATHS+=("inventory/accounting" "inventory/component" "inventory/api-gateway")
                ;;
            "metrics")
                PATHS+=("metrics/exporter" "metrics/api-gateway" "metrics/sanitizer")
                ;;
            "messaging")
                PATHS+=("messaging/mesh" "messaging/nns" "messaging/node-feeder" "messaging/api-gateway")
                ;;
            "services")
                PATHS+=("services/msgClient")
                ;;
            "dummy")
                PATHS+=("dummy/dnode" "dummy/dsubscriber" "dummy/api-gateway" "dummy/dcontroller" "dummy/dsimfactory")
                ;;
            "ukamaagent")
                PATHS+=("ukama-agent/api-gateway" "ukama-agent/cdr" "ukama-agent/asr" "ukama-agent/node-gateway")
                ;;
        esac
    done
}

filter_make_sys

cd ../../systems
root_dir=$(pwd)

for path in "${PATHS[@]}"; do
    cd "$root_dir"
    if [[ "$path" == dummy* ]]; then
        cd ../testing/services
        cd "$path" || { echo "Failed to change directory to $path"; exit 1; }
    else
        cd "$path" || { echo "Failed to change directory to $path"; exit 1; }
    fi
    
    go mod tidy && make lint && make

    IFS='/' read -r -a path_array <<< "$path"
    system=${path_array[0]}
    service=${path_array[1]}
   
    if [ $? -eq 0 ]; then
        echo -e "${TAG} Successfully build system: ${GREEN}${system}${NC} - service: ${GREEN}${service}${NC}"
    else
        echo -e "${TAG} Failed to build system: ${RED}${system}${NC} - service: ${RED}${service}${NC}"
    fi

    cd - >/dev/null || exit
done