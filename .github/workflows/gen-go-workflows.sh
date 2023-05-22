# read the workflow template
WORKFLOW_TEMPLATE=$(cat workflow-template.yaml.templ)
PR_RELEASE_TEMPLATE=$(cat pr-release-template.yaml.templ)

WARNING="# THIS FILE IS GENERATED AUTOMATICALLY BY RUNNING gen-workflow.sh 
# DON'T CHANGE IT MANUALLY TO AVOID YOUR CHANGES BEING OVERWRITTEN
# USE workflow-template.yaml FOR MAKING CHANGES IN WORKFLOWS
"

generate(){
    for SERVICE in $(ls -l ../../${1} | grep ^d | awk '{print $9}'); do
        echo "generating workflow for ${1}/${SERVICE}"
        SANITIZED=$(echo "${SERVICE}" | sed "s#-##g")
        SANITIZED_PATH=$(echo "${1}" | tr / -)
        WORKFLOW_PATH="${SANITIZED_PATH}-${SERVICE}"
        SYSTEM_NAME=${2}

        echo $WORKFLOW_PATH
        # replace template route placeholder with route name
        WORKFLOW=$(echo "${WORKFLOW_TEMPLATE}" | sed "s#{{SERVICE}}#${1}/${SERVICE}#g" | sed "s#{{SERVICE_NAME}}#${SERVICE}#g" | sed "s#{{HELMFILE_PREFIX}}#${2}#g" \
         | sed "s#{{SANITIZED_NAME}}#${SANITIZED}#g" | sed "s#{{WORKFLOW_PATH}}#${WORKFLOW_PATH}#g" | sed "s#{{SYSTEM_NAME}}#${SYSTEM_NAME}#g" )

        PR_RELEASE=$(echo "${PR_RELEASE_TEMPLATE}" | sed "s#{{SERVICE}}#${1}/${SERVICE}#g" | sed "s#{{SERVICE_NAME}}#${SERVICE}#g" | sed "s#{{HELMFILE_PREFIX}}#${2}#g" \
         | sed "s#{{SANITIZED_NAME}}#${SANITIZED}#g" | sed "s#{{WORKFLOW_PATH}}#${WORKFLOW_PATH}#g" | sed "s#{{SYSTEM_NAME}}#${SYSTEM_NAME}#g" )
        
        # save workflow to .github/workflows/{ROUTE}
        echo "${WARNING}${WORKFLOW}" > ${WORKFLOW_PATH}.yaml
        echo "${WARNING}${PR_RELEASE}" > ${WORKFLOW_PATH}-pr-release.yaml
    done
}

generate "services/bootstrap" "bootstrap"
generate "services/cloud" "ukama"
generate "services/hub" "hub"
generate "services/factory" "factory"
generate "testing/services" "testing"
generate "services/metrics" "metrics"
generate "systems/init" "init"
generate "systems/subscriber" "subscriber"
generate "systems/registry" "registry"
generate "systems/data-plan" "data-plan"
generate "systems/services" "services"
generate "systems/auth" "auth"
generate "systems/billing" "billing"