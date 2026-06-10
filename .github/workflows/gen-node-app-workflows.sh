#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_FILE="${SCRIPT_DIR}/node-app-workflow-template.yaml.templ"

if [[ ! -f "${TEMPLATE_FILE}" ]]; then
    echo "template not found: ${TEMPLATE_FILE}" >&2
    exit 1
fi

WORKFLOW_TEMPLATE="$(cat "${TEMPLATE_FILE}")"

WARNING="# THIS FILE IS GENERATED AUTOMATICALLY BY RUNNING gen-node-app-workflows.sh
# DON'T CHANGE IT MANUALLY TO AVOID YOUR CHANGES BEING OVERWRITTEN
# USE node-app-workflow-template.yaml.templ FOR MAKING CHANGES IN WORKFLOWS

"

generate() {
    for APP in "$@"; do
        APP_PATH="${APP%%:*}"
        APP_NAME="${APP##*:}"

        # If no explicit display/build name is provided, use the app directory name.
        if [[ "${APP_PATH}" == "${APP_NAME}" ]]; then
            APP_NAME="$(basename "${APP_PATH}")"
        fi

        WORKFLOW_PATH="$(echo "${APP_PATH}" | tr / -)"

        echo "generating workflow for ${APP_PATH} -> ${WORKFLOW_PATH}.yaml"

        WORKFLOW="$(printf '%s' "${WORKFLOW_TEMPLATE}" \
            | sed "s#{{APP}}#${APP_PATH}#g" \
            | sed "s#{{APP_NAME}}#${APP_NAME}#g" \
            | sed "s#{{WORKFLOW_PATH}}#${WORKFLOW_PATH}#g")"

        printf '%s%s\n' "${WARNING}" "${WORKFLOW}" > "${WORKFLOW_PATH}.yaml"
    done
}

# Format:
#   "repo/path/to/app"              -> workflow name uses basename
#   "repo/path/to/app:display-name" -> workflow name uses explicit display-name
#
# GPS keeps the existing sample workflow name: build-node-apps-gpsd.
generate \
    "nodes/apps/backhaul" \
    "nodes/apps/controller" \
    "nodes/apps/epcemu" \
    "nodes/apps/femd" \
    "nodes/apps/init-network" \
    "nodes/apps/pcrf" \
    "nodes/apps/power" \
    "nodes/apps/switchd" \
    "nodes/apps/aisg"
