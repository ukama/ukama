#!/bin/bash -x
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

# Ukama pkg related utilities

function create_manifest_file() {
    local manifest_file="$1"
    shift
    local app_names=("$@")

    log "INFO" "Creating manifest file at ${manifest_file}"

    cat <<EOF > "${manifest_file}"
{
    "version": "0.1",
    "spaces" : [
        { "name" : "boot" },
        { "name" : "services" },
        { "name" : "reboot" }
    ],
    "capps": [
        {
            "name"   : "noded",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "boot"
        },
        {
            "name"   : "bootstrap",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "boot",
            "depends_on" : [
                {
                    "capp"  : "noded",
                    "state" : "active"
                }
            ]
        },
        {
            "name"   : "meshd",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "boot",
            "depends_on" : [
                {
                    "capp"  : "bootstrap",
                    "state" : "done"
                }
            ]
        }
    ],
    "services": [
EOF

    local entries=()
    for app in "${app_names[@]}"; do
        case "$app" in
            "wimcd"|"configd"|"metricsd"|"lookoutd"|"deviced"|"notifyd")
                entries+=(
"        {
            \"name\"   : \"$app\",
            \"tag\"    : \"latest\",
            \"restart\" : \"yes\",
            \"space\"  : \"services\"
         }"
                )
                ;;
        esac
    done

    (IFS=,
     printf "%s\n" "${entries[*]}" >> "${manifest_file}")

    echo '    ]' >> "${manifest_file}"
    echo '}' >> "${manifest_file}"
}

function copy_all_apps() {
    local repo_pkg="$1"
    local dest_pkg="$2"

    log "INFO" "Copying selected apps from ${repo_pkg} to ${dest_pkg}"

    mkdir -p "$dest_pkg"

    for app in "${APPS[@]}"; do
        app_file="${repo_pkg}/${app}_latest.tar.gz"
        if [[ -f "$app_file" ]]; then
            log "INFO" "Copying $app_file"
            cp "$app_file" "$dest_pkg/"
        else
            log "WARN" "App package not found: $app_file"
        fi
    done
}

function copy_required_libs() {
    local lib_pkg="$1"
    local dest="$2"
    local tmp_dir

    log "INFO" "Installing required libs from ${lib_pkg}"
    tmp_dir=$(mktemp -d)
    tar -zxf "${lib_pkg}/vendor_libs.tgz" -C "$tmp_dir"
    cp -f "${tmp_dir}"/* "${dest}/"

#    rm -rf "$tmp_dir"
}

get_enabled_apps() {
    local common_config="$1"
    local board_config="$2"
    declare -A app_map
    local line key val

    # Read common config
    while IFS='=' read -r key val; do
        [[ -n "$key" && "$key" != \#* ]] && app_map["$key"]="$val"
    done < "$common_config"

    # Read board-specific config if provided
    if [[ -n "$board_config" && -f "$board_config" ]]; then
        while IFS='=' read -r key val; do
            [[ -n "$key" && "$key" != \#* ]] && app_map["$key"]="$val"
        done < "$board_config"
    fi

    # Build global APPS array
    APPS=()
    for key in "${!app_map[@]}"; do
        if [[ "${app_map[$key]}" == "yes" ]]; then
            APPS+=("$key")
        fi
    done

    export APPS
}
