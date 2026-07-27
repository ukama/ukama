#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

# Ukama package and starter manifest utilities.

_manifest_builder_root() {
    cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd
}

_manifest_toml_value() {
    local file="$1"
    local section="$2"
    local key="$3"

    awk -v wanted="[$section]" -v key="$key" '
        /^[[:space:]]*\[/ {
            in_section = ($0 == wanted)
            next
        }
        in_section {
            line = $0
            sub(/[[:space:]]*#.*/, "", line)
            split(line, parts, "=")
            name = parts[1]
            gsub(/^[[:space:]]+|[[:space:]]+$/, "", name)
            if (name != key) next

            value = substr(line, index(line, "=") + 1)
            gsub(/^[[:space:]]+|[[:space:]]+$/, "", value)
            if (value ~ /^".*"$/) {
                value = substr(value, 2, length(value) - 2)
            }
            print value
            exit
        }
    ' "$file"
}

_manifest_json_escape() {
    local value="$1"

    value=${value//\\/\\\\}
    value=${value//\"/\\\"}
    value=${value//$'\n'/\\n}
    printf '%s' "$value"
}

_manifest_has_app() {
    local wanted="$1"
    shift
    local app

    for app in "$@"; do
        if [[ "$app" == "$wanted" ]]; then
            return 0
        fi
    done

    return 1
}

_manifest_emit_argv() {
    local binary="$1"
    local args="$2"
    local token
    local first=1
    local tokens=()

    tokens+=("$binary")
    if [[ -n "$args" ]]; then
        read -r -a extra <<< "$args"
        tokens+=("${extra[@]}")
    fi

    printf '['
    for token in "${tokens[@]}"; do
        if [[ $first -eq 0 ]]; then
            printf ', '
        fi
        printf '"%s"' "$(_manifest_json_escape "$token")"
        first=0
    done
    printf ']'
}

_manifest_env_entries() {
    local file="$1"

    awk -v wanted="[capp-exec.env]" '
        /^[[:space:]]*\[/ {
            in_section = ($0 == wanted)
            next
        }
        in_section {
            line = $0
            sub(/[[:space:]]*#.*/, "", line)
            if (line !~ /=/) next

            name = substr(line, 1, index(line, "=") - 1)
            value = substr(line, index(line, "=") + 1)
            gsub(/^[[:space:]]+|[[:space:]]+$/, "", name)
            gsub(/^[[:space:]]+|[[:space:]]+$/, "", value)
            if (name == "" || value == "") next

            if (value ~ /^".*"$/) {
                value = substr(value, 2, length(value) - 2)
            }
            printf "%s\t%s\n", name, value
        }
    ' "$file"
}

_manifest_emit_env() {
    local config_file="$1"
    local entry
    local key
    local value
    local first=1
    local entries=()

    mapfile -t entries < <(_manifest_env_entries "$config_file")
    if [[ ${#entries[@]} -eq 0 ]]; then
        return 1
    fi

    printf '          "env": {\n'
    for entry in "${entries[@]}"; do
        key="${entry%%$'\t'*}"
        value="${entry#*$'\t'}"
        if [[ $first -eq 0 ]]; then
            printf ',\n'
        fi
        printf '            "%s": "%s"' \
            "$(_manifest_json_escape "$key")" \
            "$(_manifest_json_escape "$value")"
        first=0
    done
    printf '\n          }'
    return 0
}

_manifest_emit_app() {
    local manifest_file="$1"
    local config_key="$2"
    local comma="$3"
    local builder_root
    local config_file
    local name
    local version
    local binary
    local path
    local args
    local cmd
    local readiness

    builder_root="$(_manifest_builder_root)"
    config_file="${builder_root}/configs/${config_key}.toml"

    if [[ ! -f "$config_file" ]]; then
        log "ERROR" "Missing app config: ${config_file}"
        return 1
    fi

    name="$(_manifest_toml_value "$config_file" "capp-exec" "name")"
    version="$(_manifest_toml_value "$config_file" "capp-exec" "version")"
    binary="$(_manifest_toml_value "$config_file" "capp-exec" "bin")"
    path="$(_manifest_toml_value "$config_file" "capp-exec" "path")"
    args="$(_manifest_toml_value "$config_file" "capp-exec" "args")"
    readiness="$(_manifest_toml_value \
        "$config_file" "capp-exec" "readiness")"

    [[ -n "$name" ]] || name="$config_key"
    [[ -n "$version" ]] || version="latest"
    [[ -n "$binary" ]] || {
        log "ERROR" "Missing capp-exec.bin in ${config_file}"
        return 1
    }
    [[ -n "$path" ]] || path="/sbin"

    path="${path%/}"
    cmd="${path#/}/${binary}"

    if [[ "$comma" == "yes" ]]; then
        printf ',\n' >> "$manifest_file"
    fi

    {
        printf '        {\n'
        printf '          "name": "%s",\n' \
            "$(_manifest_json_escape "$name")"
        printf '          "tag": "%s",\n' \
            "$(_manifest_json_escape "$version")"
        printf '          "cmd": "%s",\n' \
            "$(_manifest_json_escape "$cmd")"
        printf '          "argv": '
        _manifest_emit_argv "$binary" "$args"

        if [[ "$readiness" == "true" ]]; then
            printf ',\n'
            printf '          "readiness": "yes"'
        fi

        if _manifest_env_entries "$config_file" | grep -q .; then
            printf ',\n'
            _manifest_emit_env "$config_file"
            printf '\n'
        else
            printf '\n'
        fi
        printf '        }'
    } >> "$manifest_file"
}

create_starter_manifest() {
    local manifest_file="$1"
    shift
    local app_names=("$@")
    local boot_order=("init-network" "noded" "bootstrap" "meshd")
    local service_order=()
    local app
    local first

    log "INFO" "Creating starter manifest at ${manifest_file}"
    mkdir -p "$(dirname "$manifest_file")"

    if _manifest_has_app "rlog" "${app_names[@]}"; then
        service_order+=("rlog")
    fi

    for app in "${app_names[@]}"; do
        case "$app" in
            starterd|init-network|noded|bootstrap|meshd|rlog)
                continue
                ;;
        esac
        service_order+=("$app")
    done

    cat > "$manifest_file" <<'JSON'
{
  "version": "0.1",
  "target": "ukama-node",
  "spaces": [
    {
      "name": "boot",
      "apps": [
JSON

    first=yes
    for app in "${boot_order[@]}"; do
        if ! _manifest_has_app "$app" "${app_names[@]}"; then
            continue
        fi
        if [[ "$first" == "yes" ]]; then
            _manifest_emit_app "$manifest_file" "$app" "no"
            first=no
        else
            _manifest_emit_app "$manifest_file" "$app" "yes"
        fi
    done

    cat >> "$manifest_file" <<'JSON'

      ]
    },
    {
      "name": "services",
      "apps": [
JSON

    first=yes
    for app in "${service_order[@]}"; do
        if [[ "$first" == "yes" ]]; then
            _manifest_emit_app "$manifest_file" "$app" "no"
            first=no
        else
            _manifest_emit_app "$manifest_file" "$app" "yes"
        fi
    done

    cat >> "$manifest_file" <<'JSON'

      ]
    }
  ]
}
JSON
}

# Backward-compatible function name used by the image scripts.
create_manifest_file() {
    create_starter_manifest "$@"
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
}

get_enabled_apps() {
    local common_config="$1"
    local board_config="$2"
    declare -A app_map
    local line key val

    while IFS='=' read -r key val; do
        [[ -n "$key" && "$key" != \#* ]] && app_map["$key"]="$val"
    done < "$common_config"

    if [[ -n "$board_config" && -f "$board_config" ]]; then
        while IFS='=' read -r key val; do
            [[ -n "$key" && "$key" != \#* ]] && app_map["$key"]="$val"
        done < "$board_config"
    fi

    APPS=()
    while IFS= read -r key; do
        if [[ "${app_map[$key]}" == "yes" ]]; then
            APPS+=("$key")
        fi
    done < <(printf '%s\n' "${!app_map[@]}" | sort)

    export APPS
}
