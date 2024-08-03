#!/bin/bash

# read the workflow template
WORKFLOW_TEMPLATE=$(cat ukamaos-workflow-template.yaml.templ)

WARNING="# THIS FILE IS GENERATED AUTOMATICALLY BY RUNNING gen-workflow.sh 
# DON'T CHANGE IT MANUALLY TO AVOID YOUR CHANGES BEING OVERWRITTEN
# USE workflow-template.yaml FOR MAKING CHANGES IN WORKFLOWS
"

generate(){
    for APP in "$@"; do
        echo "generating workflow for ${APP}"
        SANITIZED=$(echo "${APP}" | sed "s#-##g")
        SANITIZED_PATH=$(echo "${APP}" | tr / -)
        WORKFLOW_PATH="${SANITIZED_PATH}"
        APP_NAME=$(basename ${APP})

        echo $WORKFLOW_PATH
        # replace template route placeholder with route name
        WORKFLOW=$(echo "${WORKFLOW_TEMPLATE}" | sed "s#{{APP}}#${APP}#g" | sed "s#{{APP_NAME}}#${APP_NAME}#g" | sed "s#{{WORKFLOW_PATH}}#${WORKFLOW_PATH}#g" )

        # save workflow to .github/workflows/{WORKFLOW_PATH}.yaml
        echo "${WARNING}${WORKFLOW}" > ${WORKFLOW_PATH}.yaml
    done
}

generate "nodes/ukamaOS/distro/system/bootstrap" \
         "nodes/ukamaOS/distro/system/configd" \
         "nodes/ukamaOS/distro/system/deviced" \
         "nodes/ukamaOS/distro/system/example" \
         "nodes/ukamaOS/distro/system/init" \
         "nodes/ukamaOS/distro/system/lookoutd" \
         "nodes/ukamaOS/distro/system/meshd" \
         "nodes/ukamaOS/distro/system/metricsd" \
         "nodes/ukamaOS/distro/system/noded" \
         "nodes/ukamaOS/distro/system/notifyd" \
         "nodes/ukamaOS/distro/system/rlog" \
         "nodes/ukamaOS/distro/system/starterd" \
         "nodes/ukamaOS/distro/system/wimcd"

