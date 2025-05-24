#!/bin/bash
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

    for app in "${app_names[@]}"; do
        case "$app" in
            "wimcd"|"configd"|"metricsd"|"lookoutd"|"deviced"|"notifyd")
                cat <<EOF >> "${manifest_file}"
        {
            "name"   : "$app",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "services"
        },
EOF
                ;;
        esac
    done

    # Remove last comma
    sed -i '$ s/,$//' "${manifest_file}"
    echo '    ]' >> "${manifest_file}"
    echo '}' >> "${manifest_file}"
}

function copy_all_apps() {
    local repo_pkg="$1"
    local dest_pkg="$2"

    log "INFO" "Copying apps from ${repo_pkg} to ${dest_pkg}"
    cp -rvf "${repo_pkg}" "${dest_pkg}"
}

function copy_required_libs() {
    local lib_pkg="$1"
    local rootfs="$2"

    log "INFO" "Installing required libs from ${lib_pkg}"
    pushd "${lib_pkg}" > /dev/null
    tar zxvf vendor_libs.tgz 
    cp -vrf ./build/* "${rootfs}/usr/"
    popd > /dev/null
}

