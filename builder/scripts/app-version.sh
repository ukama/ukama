#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
BUILDER_ROOT="$(CDPATH= cd -- "${SCRIPT_DIR}/.." && pwd)"
REPO_ROOT="$(CDPATH= cd -- "${BUILDER_ROOT}/.." && pwd)"
CONFIG_DIR="${BUILDER_ROOT}/configs"
APP_BUILDER="${BUILDER_ROOT}/app_builder"
PACKAGE_DIR="${REPO_ROOT}/build/pkgs"

HUB_URL="${HUB_URL:-https://hub-ukama.udev.ukama.com}"
PREFIX="${UKAMA_TEST_VERSION_PREFIX:-2.0.0}"
TARGET_VERSION=""
COMMAND=""
SELECT_ALL=0
PREFIX_SET=0
TARGET_SET=0

REQUESTED_APPS=()
SELECTED_CONFIGS=()
MISSING_CONFIGS=()
HUB_APPS=()

usage() {
    cat <<'USAGE'
Usage:
  app_version.sh list
  app_version.sh hub [--app APP ...] [--hub-url URL]
  app_version.sh build  (--app APP ... | --all) [--prefix X.Y.Z]
  app_version.sh ensure (--app APP ... | --all) [--prefix X.Y.Z]
                 [--hub-url URL]

Options:
  --app APP          Select by config name or Hub app name. Repeatable.
                     With "hub", omit this to list every Hub app.
  --all              Select all tar.gz apps in builder/configs.
  --prefix X.Y.Z     Version prefix. Default: 2.0.0.
  --target VERSION   Use an exact semantic version.
  --hub-url URL      Hub base URL.
  -h, --help         Show this help.

Default target:
  <prefix>-lab.g<git-hash>

Examples:
  ./builder/scripts/app_version.sh list
  ./builder/scripts/app_version.sh hub
  ./builder/scripts/app_version.sh hub --app gpsd
  ./builder/scripts/app_version.sh ensure --app gpsd
  ./builder/scripts/app_version.sh ensure --app example --app gpsd
  ./builder/scripts/app_version.sh ensure --all --prefix 2.0.0
USAGE
}

log() {
    printf '%-11s %s\n' "$1" "$2"
}

die() {
    log "ERROR" "$*" >&2
    exit 1
}

need() {
    command -v "$1" >/dev/null 2>&1 ||
        die "Required command not found: $1"
}

toml_value() {
    local file="$1"
    local section="$2"
    local key="$3"

    awk -v wanted="[$section]" -v wanted_key="$key" '
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
            if (name != wanted_key) next

            if (value ~ /^".*"$/) {
                value = substr(value, 2, length(value) - 2)
            }
            print value
            exit
        }
    ' "$file"
}

config_key() {
    basename "$1" .toml
}

app_name() {
    toml_value "$1" "capp-exec" "name"
}

app_bin() {
    toml_value "$1" "capp-exec" "bin"
}

app_path() {
    local path

    path="$(toml_value "$1" "capp-exec" "path")"
    [[ -n "$path" ]] || path="/sbin"
    path="/${path#/}"
    path="${path%/}"
    [[ -n "$path" ]] || path="/"
    printf '%s\n' "$path"
}

app_format() {
    toml_value "$1" "capp-output" "format"
}

validate_config() {
    local config="$1"
    local name
    local bin
    local format

    name="$(app_name "$config")"
    bin="$(app_bin "$config")"
    format="$(app_format "$config")"

    [[ -n "$name" ]] || die "Missing capp-exec.name: $config"
    [[ "$name" =~ ^[A-Za-z0-9_-]+$ ]] ||
        die "Invalid app name in $config: $name"
    [[ -n "$bin" ]] || die "Missing capp-exec.bin: $config"
    [[ -n "$format" ]] || die "Missing capp-output.format: $config"
}

add_config() {
    local config="$1"
    local selected

    for selected in "${SELECTED_CONFIGS[@]:-}"; do
        [[ "$selected" == "$config" ]] && return 0
    done
    SELECTED_CONFIGS+=("$config")
}

resolve_selection() {
    local config
    local requested
    local match
    local found

    shopt -s nullglob

    if [[ "$SELECT_ALL" -eq 1 ]]; then
        for config in "${CONFIG_DIR}"/*.toml; do
            validate_config "$config"
            [[ "$(app_format "$config")" == "tar.gz" ]] || continue
            add_config "$config"
        done
    fi

    for requested in "${REQUESTED_APPS[@]:-}"; do
        found=""
        for config in "${CONFIG_DIR}"/*.toml; do
            validate_config "$config"
            match="$(app_name "$config")"
            if [[ "$requested" == "$(config_key "$config")" ||
                  "$requested" == "$match" ]]; then
                [[ -z "$found" ]] ||
                    die "App name is ambiguous: $requested"
                found="$config"
            fi
        done

        [[ -n "$found" ]] || die "Unknown app: $requested"
        [[ "$(app_format "$found")" == "tar.gz" ]] ||
            die "Only tar.gz apps are supported: $requested"
        add_config "$found"
    done

    shopt -u nullglob
}

set_target_version() {
    local hash

    if [[ -n "$TARGET_VERSION" ]]; then
        TARGET_VERSION="${TARGET_VERSION#v}"
    else
        PREFIX="${PREFIX#v}"
        [[ "$PREFIX" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] ||
            die "Invalid version prefix: $PREFIX"

        hash="$(
            git -C "$REPO_ROOT" \
                -c safe.directory="$REPO_ROOT" \
                rev-parse --short HEAD
        )"
        TARGET_VERSION="${PREFIX}-lab.g${hash}"
    fi

    [[ "$TARGET_VERSION" =~ \
       ^[0-9]+\.[0-9]+\.[0-9]+([+-][0-9A-Za-z.-]+)?$ ]] ||
        die "Invalid semantic version: $TARGET_VERSION"
}

hub_response() {
    local app="$1"
    local output="$2"
    local status

    status="$(
        curl --silent --show-error \
            --output "$output" \
            --write-out '%{http_code}' \
            --header 'Accept: application/json' \
            "${HUB_URL%/}/v1/hub/app/${app}"
    )" || return 2

    case "$status" in
        200) return 0 ;;
        404) return 1 ;;
        *)
            log "HUB" "GET ${app} returned HTTP ${status}" >&2
            cat "$output" >&2 || true
            return 2
            ;;
    esac
}

hub_list_response() {
    local output="$1"
    local status

    status="$(
        curl --silent --show-error \
            --output "$output" \
            --write-out '%{http_code}' \
            --header 'Accept: application/json' \
            "${HUB_URL%/}/v1/hub/app"
    )" || return 1

    if [[ "$status" != "200" ]]; then
        log "HUB" "GET app list returned HTTP ${status}" >&2
        cat "$output" >&2 || true
        return 1
    fi
}

add_hub_app() {
    local app="$1"
    local selected

    for selected in "${HUB_APPS[@]:-}"; do
        [[ "$selected" == "$app" ]] && return 0
    done
    HUB_APPS+=("$app")
}

resolve_hub_app_name() {
    local requested="$1"
    local config
    local name

    shopt -s nullglob
    for config in "${CONFIG_DIR}"/*.toml; do
        name="$(app_name "$config")"
        if [[ "$requested" == "$(config_key "$config")" ||
              "$requested" == "$name" ]]; then
            shopt -u nullglob
            printf '%s\n' "$name"
            return
        fi
    done
    shopt -u nullglob

    printf '%s\n' "$requested"
}

load_hub_apps() {
    local requested
    local response
    local app

    if [[ ${#REQUESTED_APPS[@]} -gt 0 ]]; then
        for requested in "${REQUESTED_APPS[@]}"; do
            app="$(resolve_hub_app_name "$requested")"
            [[ -n "$app" ]] || die "Invalid app: $requested"
            add_hub_app "$app"
        done
        return
    fi

    response="$(mktemp)"
    hub_list_response "$response" || {
        rm -f "$response"
        die "Unable to query Hub app list"
    }

    if ! jq empty "$response" >/dev/null 2>&1; then
        rm -f "$response"
        die "Invalid JSON response from Hub app list"
    fi

    while IFS= read -r app; do
        [[ -n "$app" ]] && add_hub_app "$app"
    done < <(
        jq -r '(.artifact // .Artifact // [])[]?' "$response" |
            LC_ALL=C sort
    )

    rm -f "$response"
}

print_hub_versions() {
    local app
    local response
    local rc

    printf '%-20s %-40s %-12s %s\n' \
        "APP" "VERSION" "SIZE" "CREATED"

    for app in "${HUB_APPS[@]}"; do
        response="$(mktemp)"
        if hub_response "$app" "$response"; then
            :
        else
            rc=$?
            rm -f "$response"
            if [[ "$rc" -eq 1 ]]; then
                die "Hub app not found: $app"
            fi
            die "Unable to query Hub for $app"
        fi

        if ! jq empty "$response" >/dev/null 2>&1; then
            rm -f "$response"
            die "Invalid JSON response for $app"
        fi

        jq -r --arg app "$app" '
            (.versions // [])
            | sort_by(.version)[]? as $version
            | ($version.FormatInfo // $version.formats // [])[]?
            | select((.Type // .type) == "tar.gz")
            | [
                $app,
                $version.version,
                ((.Size // .size // "0") | tostring),
                (.createdAt // "")
              ]
            | @tsv
        ' "$response" |
        while IFS=$'\t' read -r name version size created; do
            printf '%-20s %-40s %-12s %s\n' \
                "$name" "$version" "$size" "$created"
        done

        rm -f "$response"
    done
}

run_hub() {
    need curl
    need jq

    load_hub_apps
    print_hub_versions
}

hub_has_version() {
    local app="$1"
    local response
    local rc

    response="$(mktemp)"
    if hub_response "$app" "$response"; then
        rc=0
    else
        rc=$?
        rm -f "$response"
        [[ "$rc" -eq 1 ]] && return 1
        return 2
    fi

    if ! jq empty "$response" >/dev/null 2>&1; then
        rm -f "$response"
        log "HUB" "Invalid JSON response for ${app}" >&2
        return 2
    fi

    if jq -e --arg version "$TARGET_VERSION" '
        .versions[]?
        | select(.version == $version)
        | (.FormatInfo // .formats // [])[]?
        | select(
            .Type == "tar.gz" and
            ((.Size | tonumber?) // 0) > 0 and
            (.Url | type == "string") and
            (.Url | endswith(".tar.gz"))
          )
    ' "$response" >/dev/null; then
        rm -f "$response"
        return 0
    fi

    rm -f "$response"
    return 1
}

ensure_app_builder() {
    log "BUILDER" "building app_builder"
    make -C "$BUILDER_ROOT" app_builder
    [[ -x "$APP_BUILDER" ]] || die "app_builder was not created"
}

package_path() {
    local config="$1"

    printf '%s/%s_%s.tar.gz\n' \
        "$PACKAGE_DIR" "$(app_name "$config")" "$TARGET_VERSION"
}

validate_package() {
    local config="$1"
    local name
    local bin
    local path
    local package
    local root
    local expected_bin
    local embedded_version
    local listing

    name="$(app_name "$config")"
    bin="$(app_bin "$config")"
    path="$(app_path "$config")"
    package="$(package_path "$config")"
    root="${name}_${TARGET_VERSION}"
    expected_bin="${root}${path%/}/${bin}"
    listing="$(mktemp)"

    [[ -s "$package" ]] || die "Package missing or empty: $package"
    gzip -t "$package" || die "Invalid gzip package: $package"
    tar -tzf "$package" > "$listing" ||
        die "Invalid tar archive: $package"

    if awk '
        /^\// { bad = 1 }
        /(^|\/)\.\.($|\/)/ { bad = 1 }
        END { exit bad ? 0 : 1 }
    ' "$listing"; then
        rm -f "$listing"
        die "Unsafe path found in package: $package"
    fi

    if ! awk -v root="${root}/" '
        $0 != substr(root, 1, length(root) - 1) &&
        index($0, root) != 1 { bad = 1 }
        END { exit bad ? 1 : 0 }
    ' "$listing"; then
        rm -f "$listing"
        die "Package contains files outside ${root}: $package"
    fi

    grep -Fxq "${root}/VERSION" "$listing" || {
        rm -f "$listing"
        die "VERSION missing from package: $package"
    }

    grep -Fxq "$expected_bin" "$listing" || {
        rm -f "$listing"
        die "Expected binary missing from package: $expected_bin"
    }

    embedded_version="$(tar -xOf "$package" "${root}/VERSION")"
    [[ "$embedded_version" == "$TARGET_VERSION" ]] || {
        rm -f "$listing"
        die "VERSION mismatch in package: $package"
    }

    rm -f "$listing"
    log "VERIFY" "${name}:${TARGET_VERSION}"
}

build_package() {
    local config="$1"
    local name

    name="$(app_name "$config")"
    log "BUILD" "${name} using $(config_key "$config").toml"

    (
        cd "$BUILDER_ROOT"
        UKAMA_ROOT="$REPO_ROOT" \
        UKAMA_APP_VERSION="$TARGET_VERSION" \
        UKAMA_SKIP_LATEST_LINK=1 \
            "$APP_BUILDER" --create --config "$config"
    )

    validate_package "$config"
}

upload_package() {
    local config="$1"
    local name
    local package
    local response
    local status

    name="$(app_name "$config")"
    package="$(package_path "$config")"
    response="$(mktemp)"

    log "UPLOAD" "${name}:${TARGET_VERSION}"
    status="$(
        curl --silent --show-error \
            --output "$response" \
            --write-out '%{http_code}' \
            --request PUT \
            --header 'Accept: application/json' \
            --header 'Content-Type: application/gzip' \
            --data-binary "@${package}" \
            "${HUB_URL%/}/v1/hub/app/${name}/${TARGET_VERSION}"
    )" || {
        rm -f "$response"
        die "Hub upload failed for ${name}:${TARGET_VERSION}"
    }

    if [[ "$status" != "201" ]]; then
        cat "$response" >&2 || true
        rm -f "$response"
        die "Hub upload returned HTTP ${status} for ${name}"
    fi

    rm -f "$response"
}

list_apps() {
    local config
    local path

    printf '%-22s %-20s %-8s %s\n' \
        "CONFIG" "HUB APP" "FORMAT" "BINARY"

    shopt -s nullglob
    for config in "${CONFIG_DIR}"/*.toml; do
        validate_config "$config"
        path="$(app_path "$config")"
        printf '%-22s %-20s %-8s %s%s%s\n' \
            "$(config_key "$config")" \
            "$(app_name "$config")" \
            "$(app_format "$config")" \
            "${path%/}" \
            "$([[ "$path" == "/" ]] && printf '' || printf '/')" \
            "$(app_bin "$config")"
    done
    shopt -u nullglob
}

run_build() {
    local config

    [[ ${#SELECTED_CONFIGS[@]} -gt 0 ]] ||
        die "Select at least one app with --app or --all"

    need git
    need gzip
    need tar
    need make

    set_target_version
    log "TARGET" "$TARGET_VERSION"
    ensure_app_builder

    for config in "${SELECTED_CONFIGS[@]}"; do
        build_package "$config"
    done

    log "RESULT" "all packages built"
}

run_ensure() {
    local config
    local name
    local rc

    [[ ${#SELECTED_CONFIGS[@]} -gt 0 ]] ||
        die "Select at least one app with --app or --all"

    need git
    need curl
    need jq
    need gzip
    need tar
    need make

    set_target_version
    log "TARGET" "$TARGET_VERSION"

    for config in "${SELECTED_CONFIGS[@]}"; do
        name="$(app_name "$config")"
        if hub_has_version "$name"; then
            log "HUB" "${name}:${TARGET_VERSION} available"
        else
            rc=$?
            [[ "$rc" -eq 1 ]] || die "Unable to query Hub for ${name}"
            log "HUB" "${name}:${TARGET_VERSION} missing"
            MISSING_CONFIGS+=("$config")
        fi
    done

    if [[ ${#MISSING_CONFIGS[@]} -eq 0 ]]; then
        log "RESULT" "all packages already available"
        return
    fi

    ensure_app_builder

    for config in "${MISSING_CONFIGS[@]}"; do
        build_package "$config"
    done

    for config in "${MISSING_CONFIGS[@]}"; do
        upload_package "$config"
    done

    for config in "${MISSING_CONFIGS[@]}"; do
        name="$(app_name "$config")"
        hub_has_version "$name" ||
            die "Hub verification failed for ${name}:${TARGET_VERSION}"
        log "HUB" "${name}:${TARGET_VERSION} published"
    done

    log "RESULT" "all packages available"
}

parse_args() {
    [[ "$#" -gt 0 ]] || {
        usage
        exit 1
    }

    COMMAND="$1"
    shift

    while [[ "$#" -gt 0 ]]; do
        case "$1" in
            --app)
                [[ "$#" -ge 2 ]] || die "--app requires a value"
                REQUESTED_APPS+=("$2")
                shift 2
                ;;
            --all)
                SELECT_ALL=1
                shift
                ;;
            --prefix)
                [[ "$#" -ge 2 ]] || die "--prefix requires a value"
                PREFIX="$2"
                PREFIX_SET=1
                shift 2
                ;;
            --target)
                [[ "$#" -ge 2 ]] || die "--target requires a value"
                TARGET_VERSION="$2"
                TARGET_SET=1
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

    [[ "$SELECT_ALL" -eq 0 || ${#REQUESTED_APPS[@]} -eq 0 ]] ||
        die "Use either --all or --app, not both"
    [[ "$TARGET_SET" -eq 0 || "$PREFIX_SET" -eq 0 ]] ||
        die "Use either --target or --prefix, not both"
}

main() {
    parse_args "$@"

    case "$COMMAND" in
        list)
            [[ "$SELECT_ALL" -eq 0 &&
               ${#REQUESTED_APPS[@]} -eq 0 ]] ||
                die "list does not accept app selection"
            list_apps
            ;;
        hub)
            [[ "$SELECT_ALL" -eq 0 ]] ||
                die "hub lists all apps when --app is omitted"
            [[ "$PREFIX_SET" -eq 0 && "$TARGET_SET" -eq 0 ]] ||
                die "hub does not accept --prefix or --target"
            run_hub
            ;;
        build)
            resolve_selection
            run_build
            ;;
        ensure)
            resolve_selection
            run_ensure
            ;;
        help|-h|--help)
            usage
            ;;
        *)
            die "Unknown command: $COMMAND"
            ;;
    esac
}

main "$@"
