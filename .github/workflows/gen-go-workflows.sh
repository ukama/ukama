# read the workflow template
WORKFLOW_TEMPLATE=$(cat workflow-template.yaml.templ)

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
        
        echo $WORKFLOW_PATH
        # replace template route placeholder with route name
        WORKFLOW=$(echo "${WORKFLOW_TEMPLATE}" | sed "s#{{SERVICE}}#${1}/${SERVICE}#g" | sed "s#{{SERVICE_NAME}}#${SERVICE}#g" | sed "s#{{HELMFILE_PREFIX}}#${2}#g" \
         | sed "s#{{SANITIZED_NAME}}#${SANITIZED}#g" | sed "s#{{WORKFLOW_PATH}}#${WORKFLOW_PATH}#g" )
        
        # save workflow to .github/workflows/{ROUTE}
        echo "${WARNING}${WORKFLOW}" > ${WORKFLOW_PATH}.yaml
    done
}

generate "services/bootstrap" "bootstrap"
generate "services/cloud" "ukama"
generate "services/hub" "hub"
generate "testing/services" "factory"
generate "services/metrics" "metrics"
