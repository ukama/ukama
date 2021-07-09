# read the workflow template
WORKFLOW_TEMPLATE=$(cat workflow-template.yaml.templ)

WARNING="# THIS FILE IS GENERATED AUTOMATICALLY BY RUNNING gen-workflow.sh 
# DON'T CHANGE IT MANUALLY TO AVOID YOUR CHANGES BEING OVERWRITTEN
# USE workflow-template.yaml FOR MAKING CHANGES IN WORKFLOWS
"

# iterate each route in routes directory
for SERVICE in $(ls -l ../../ | grep ^d | awk '{print $9}'); do
    echo "generating workflow for ${SERVICE}"

    # replace template route placeholder with route name
    WORKFLOW=$(echo "${WORKFLOW_TEMPLATE}" | sed "s/{{SERVICE}}/${SERVICE}/g")

    # save workflow to .github/workflows/{ROUTE}
    echo "${WARNING}${WORKFLOW}" > ${SERVICE}.yml
done