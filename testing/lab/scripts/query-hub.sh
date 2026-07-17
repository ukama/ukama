#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

HUB_URL="${HUB_URL:-https://hub-ukama.udev.ukama.com}"
APPS=()

usage() {
    cat <<'USAGE'
Usage:
  query-hub.sh [--app APP ...] [--hub-url URL]

Options:
  --app APP       Query one Hub app. Repeatable.
                  Omit to query all Hub apps.
  --hub-url URL   Hub base URL.
  -h, --help      Show this help.

Examples:
  ./query-hub.sh
  ./query-hub.sh --app gpsd
  ./query-hub.sh --app gpsd --app example
  ./query-hub.sh --hub-url https://hub-ukama.udev.ukama.com
USAGE
}

die() {
    printf 'ERROR       %s\n' "$*" >&2
    exit 1
}

need() {
    command -v "$1" >/dev/null 2>&1 ||
        die "Required command not found: $1"
}

hub_get() {
    local path="$1"

    curl --fail --silent --show-error \
        --header 'Accept: application/json' \
        "${HUB_URL%/}${path}"
}

load_all_apps() {
    local response

    response="$(hub_get '/v1/hub/app')" ||
        die "Unable to query Hub app list"

    mapfile -t APPS < <(
        jq -r '(.artifact // .Artifact // [])[]?' <<< "$response" |
            LC_ALL=C sort -u
    )
}

print_app_versions() {
    local app="$1"
    local response

    [[ "$app" =~ ^[A-Za-z0-9._-]+$ ]] ||
        die "Invalid app name: $app"

    response="$(hub_get "/v1/hub/app/${app}")" ||
        die "Unable to query Hub app: $app"

    jq -r --arg app "$app" '
        (.versions // [])
        | sort_by(.version)[]? as $version
        | ($version.FormatInfo // $version.formats // [])[]?
        | select((.Type // .type) == "tar.gz")
        | [
            $app,
            ($version.version // "-"),
            (.createdAt // $version.createdAt // "-")
          ]
        | @tsv
    ' <<< "$response" |
    while IFS=$'\t' read -r name version created; do
        printf '%-20s %-40s %s\n' \
            "$name" "$version" "$created"
    done
}

parse_args() {
    while [[ "$#" -gt 0 ]]; do
        case "$1" in
            --app)
                [[ "$#" -ge 2 ]] || die "--app requires a value"
                APPS+=("$2")
                shift 2
                ;;
            --hub-url)
                [[ "$#" -ge 2 ]] || die "--hub-url requires a value"
                HUB_URL="$2"
                shift 2
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                die "Unknown argument: $1"
                ;;
        esac
    done
}

main() {
    local app

    parse_args "$@"

    need curl
    need jq

    if [[ ${#APPS[@]} -eq 0 ]]; then
        load_all_apps
    fi

    printf '%-20s %-40s %s\n' "APP" "VERSION" "CREATED"

    for app in "${APPS[@]}"; do
        print_app_versions "$app"
    done
}

main "$@"
