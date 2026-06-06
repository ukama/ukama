#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

REPO="ukama/ukama"

check_auth() {
    gh auth status >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "${RED}GitHub auth failed.${NC}"
        echo "Run:"
        echo "  gh auth login -h github.com -p https -w"
        exit 1
    fi
}

get_workflows() {
    gh workflow list --repo "$REPO" --limit 300 --all | grep "^build-" | awk '{print $1}'
}

get_last_run_id() {
    local workflow_name=$1
    local branch=$2

    gh run list \
        --repo "$REPO" \
        --workflow "$workflow_name" \
        --branch "$branch" \
        --limit 1 \
        --json databaseId \
        --jq '.[0].databaseId // empty'
}

get_run_status() {
    local run_id=$1

    gh run view "$run_id" \
        --repo "$REPO" \
        --json status,conclusion \
        --jq 'if .conclusion == null then .status else .conclusion end'
}

get_run_jobs() {
    local run_id=$1

    gh run view "$run_id" \
        --repo "$REPO" \
        --json jobs \
        --jq '
            .jobs
            | map(
                .name + "=" +
                (
                    if .conclusion == null then
                        .status
                    else
                        .conclusion
                    end
                )
              )
            | join(", ")
        '
}

print_status_line() {
    local workflow_name=$1
    local status=$2
    local jobs=$3

    if [ "$status" = "success" ]; then
        printf "%-45s ${GREEN}%-15s${NC} %s\n" "$workflow_name" "$status" "$jobs"
    elif [ "$status" = "failure" ]; then
        printf "%-45s ${RED}%-15s${NC} %s\n" "$workflow_name" "$status" "$jobs"
    else
        printf "%-45s ${YELLOW}%-15s${NC} %s\n" "$workflow_name" "$status" "$jobs"
    fi
}

show_status() {
    local branch=${1:-main}

    check_auth

    workflows=$(get_workflows)

    if [ -z "$workflows" ]; then
        echo "No build-* workflows found."
        exit 1
    fi

    printf "\n%-45s %-15s %s\n" "Workflow Name" "Status" "Jobs"
    printf "%-45s %-15s %s\n" "-------------" "------" "----"

    for workflow_name in $workflows; do
        run_id=$(get_last_run_id "$workflow_name" "$branch")

        if [ -z "$run_id" ]; then
            print_status_line "$workflow_name" "no-runs" "-"
            continue
        fi

        status=$(get_run_status "$run_id")
        jobs=$(get_run_jobs "$run_id")

        print_status_line "$workflow_name" "$status" "$jobs"
    done
}

run_workflows() {
    local branch=${1:-main}

    check_auth

    workflows=$(get_workflows)

    if [ -z "$workflows" ]; then
        echo "No build-* workflows found."
        exit 1
    fi

    declare -A workflow_runs

    echo "Triggering build-* workflows on branch: $branch"
    echo

    for workflow_name in $workflows; do
        echo "Triggering: $workflow_name"

        gh workflow run "$workflow_name" --repo "$REPO" --ref "$branch"
        if [ $? -ne 0 ]; then
            echo -e "${YELLOW}Skipping: $workflow_name cannot be manually triggered${NC}"
            echo
            continue
        fi

        sleep 3

        run_id=$(get_last_run_id "$workflow_name" "$branch")

        if [ -z "$run_id" ]; then
            echo -e "${YELLOW}Could not find run id for: $workflow_name${NC}"
            echo
            continue
        fi

        workflow_runs["$workflow_name"]=$run_id

        echo "Run ID: $run_id"
        echo
    done

    if [ ${#workflow_runs[@]} -eq 0 ]; then
        echo "No workflows were triggered."
        exit 1
    fi

    echo "Waiting for triggered workflows to complete..."
    echo

    all_done=false
    while [ "$all_done" = false ]; do
        all_done=true

        for workflow_name in "${!workflow_runs[@]}"; do
            run_id=${workflow_runs[$workflow_name]}
            status=$(get_run_status "$run_id")

            if [ "$status" = "queued" ] || [ "$status" = "in_progress" ] || [ "$status" = "waiting" ] || [ "$status" = "pending" ]; then
                all_done=false
                echo "$workflow_name is still $status..."
            fi
        done

        if [ "$all_done" = false ]; then
            sleep 15
        fi
    done

    echo
    echo "Final status:"
    echo

    printf "%-45s %-15s %s\n" "Workflow Name" "Status" "Jobs"
    printf "%-45s %-15s %s\n" "-------------" "------" "----"

    for workflow_name in "${!workflow_runs[@]}"; do
        run_id=${workflow_runs[$workflow_name]}
        status=$(get_run_status "$run_id")
        jobs=$(get_run_jobs "$run_id")

        print_status_line "$workflow_name" "$status" "$jobs"
    done
}

if [ "$1" = "run" ]; then
    run_workflows "$2"

elif [ "$1" = "status" ]; then
    show_status "$2"

else
    echo "Usage:"
    echo "  $0 run [branch]"
    echo "  $0 status [branch]"
    echo
    echo "Examples:"
    echo "  $0 run main"
    echo "  $0 status main"
    exit 1
fi
