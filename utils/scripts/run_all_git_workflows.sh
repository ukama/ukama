#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# ANSI color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

declare -A workflow_runs

# Function to check the status of a workflow run
check_run_status() {
    local run_id=$1
    gh run view "$run_id" --json conclusion --jq '.conclusion'
}

# Fetch workflows
workflows=$(gh workflow list --limit 300 --all | grep "systems-" | awk '{print $1}')

# Run each workflow and store the run ID
for workflow_name in $workflows; do
    gh workflow run "$workflow_name" --ref main
    sleep 2
    run_id=$(gh run list --workflow="$workflow_name" --limit=1 --json databaseId --jq '.[0].databaseId')
    workflow_runs["$workflow_name"]=$run_id
    echo "Triggered workflow: $workflow_name, Run ID: $run_id"
done

echo "Waiting for workflows to complete..."

all_done=false
while [ "$all_done" = false ]; do
  all_done=true
  for workflow_name in "${!workflow_runs[@]}"; do
    run_id=${workflow_runs[$workflow_name]}
    status=$(check_run_status "$run_id")
    if [ "$status" = "" ]; then
      all_done=false
      echo "Workflow: $workflow_name, Run ID: $run_id is still running..."
      break
    else
      echo "Workflow: $workflow_name, Run ID: $run_id has completed with status: $status"
    fi
  done
  if [ "$all_done" = false ]; then
    sleep 10  # Wait before checking again
  fi
done

echo "All workflows have completed."

# Print final status of all workflows
printf "\n%-40s %-15s\n" "Workflow-name" "Status"
printf "%-40s %-15s\n" "-----------------" "-----------"
for workflow_name in "${!workflow_runs[@]}"; do
  run_id=${workflow_runs[$workflow_name]}
  status=$(check_run_status "$run_id")
  if [ "$status" = "success" ]; then
    printf "%-40s ${GREEN}%-15s${NC}\n" "$workflow_name" "$status"
  else
    printf "%-40s ${RED}%-15s${NC}\n" "$workflow_name" "$status"
  fi
done
