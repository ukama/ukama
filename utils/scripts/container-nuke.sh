#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

# container-nuke.sh
#
# !!! DANGEROUS !!!
# DESTROYS ALL container engine state (for selected engine/context):
# - containers (running + stopped)
# - pods (podman only)
# - images
# - volumes (DATA LOSS)
# - networks (where possible)
# - build cache
#
# SAFETY:
#   You MUST pass --all or the script will exit.
#
# Usage:
#   ./container-nuke.sh --all --engine podman
#   ./container-nuke.sh --all --engine docker
#   ./container-nuke.sh --all              # auto-detect (podman)
#   ./container-nuke.sh --all --dry-run
#   ./container-nuke.sh --all --no-reset

CONFIRM_ALL=0
DRY_RUN=0
NO_RESET=0
ENGINE="auto"

usage() {
    cat <<'EOF'
container-nuke.sh — DESTROY ALL docker/podman resources (DANGEROUS)

REQUIRED:
  --all               Required safety flag. Script will NOT run without it.

OPTIONAL:
  --engine <auto|podman|docker>  Select engine (default: auto)
  --dry-run           Show what would be deleted, do nothing
  --no-reset          Do not use engine "reset-like" operation (where applicable)
                      (still removes containers, images, volumes, cache)

Examples:
  ./container-nuke.sh --all --engine podman
  ./container-nuke.sh --all --engine docker
  ./container-nuke.sh --all --dry-run
  ./container-nuke.sh --all --no-reset
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --all) CONFIRM_ALL=1; shift ;;
        --dry-run) DRY_RUN=1; shift ;;
        --no-reset) NO_RESET=1; shift ;;
        --engine)
            ENGINE="${2:-}"
            [[ -n "$ENGINE" ]] || { echo "ERROR: --engine requires value" >&2; exit 2; }
            shift 2
            ;;
        -h|--help) usage; exit 0 ;;
        *)
            echo "Unknown option: $1" >&2
            usage
            exit 2
            ;;
    esac
done

if [[ "$CONFIRM_ALL" -ne 1 ]]; then
    echo "ERROR: Refusing to run."
    echo "This script DESTROYS ALL container engine data."
    echo
    echo "You must explicitly pass: --all"
    echo "Example:"
    echo "  ./container-nuke.sh --all --engine podman"
    exit 1
fi

run() {
    if [[ "$DRY_RUN" -eq 1 ]]; then
        echo "[dry-run] $*"
        return 0
    fi
    echo "+ $*"
    "$@"
}

section() {
    echo
    echo "============================================================"
    echo "$1"
    echo "============================================================"
}

detect_engine() {
    if [[ "$ENGINE" == "podman" || "$ENGINE" == "docker" ]]; then
        return 0
    fi
    if command -v podman >/dev/null 2>&1; then
        ENGINE="podman"
        return 0
    fi
    if command -v docker >/dev/null 2>&1; then
        ENGINE="docker"
        return 0
    fi
    echo "ERROR: Could not auto-detect podman or docker." >&2

    exit 1
}

podman_show_state() {
    echo "== Containers =="; podman ps -a || true
    echo; echo "== Pods =="; podman pod ls -a || true
    echo; echo "== Images =="; podman images -a || true
    echo; echo "== Volumes =="; podman volume ls || true
    echo; echo "== Networks =="; podman network ls 2>/dev/null || true
}

podman_nuke() {
    section "Podman sanity check"
    command -v podman >/dev/null || { echo "podman not found"; exit 1; }
    podman --version

    section "State BEFORE"
    podman_show_state

    if [[ "$DRY_RUN" -eq 1 ]]; then
        section "DRY RUN — nothing will be deleted"
    fi

    section "Stopping containers"
    run podman stop -a || true

    section "Removing containers"
    run podman rm -a -f || true

    section "Removing pods"
    run podman pod rm -a -f || true

    if [[ "$NO_RESET" -eq 0 ]]; then
        section "SYSTEM RESET (hard nuke)"
        run podman system reset -f
    else
        section "Removing images"
        run podman rmi -a -f || true

        section "Removing volumes (DATA LOSS)"
        run podman volume rm -a -f || true

        section "Pruning everything else"
        run podman system prune -a --volumes -f || true
        run podman builder prune -a -f || true
        run podman network prune -f || true
    fi

    section "State AFTER"
    podman_show_state

    echo
    echo "✅ Podman nuke complete."
}

docker_show_state() {
    echo "== Containers =="; docker ps -a || true
    echo; echo "== Images =="; docker images -a || true
    echo; echo "== Volumes =="; docker volume ls || true
    echo; echo "== Networks =="; docker network ls || true
}

docker_nuke() {
    section "Docker sanity check"
    command -v docker >/dev/null || { echo "docker not found"; exit 1; }
    docker --version

    section "State BEFORE"
    docker_show_state

    if [[ "$DRY_RUN" -eq 1 ]]; then
        section "DRY RUN — nothing will be deleted"
    fi

    # Stop/remove containers
    section "Stopping all containers"
    # docker stop wants IDs; "docker stop $(docker ps -q)" fails on empty, so guard it.
    ids="$(docker ps -q || true)"
    if [[ -n "${ids}" ]]; then
        run docker stop ${ids} || true
    else
        echo "(none running)"
    fi

    section "Removing all containers"
    ids_all="$(docker ps -aq || true)"
    if [[ -n "${ids_all}" ]]; then
        run docker rm -f ${ids_all} || true
    else
        echo "(none to remove)"
    fi

    # Reset-like behavior
    if [[ "$NO_RESET" -eq 0 ]]; then
        section "RESET-LIKE NUKE (aggressive prune)"
        # This is the closest to podman system reset:
        run docker system prune -a --volumes -f || true
        run docker builder prune -a -f || true
    else
        section "Removing images"
        img_ids="$(docker images -aq || true)"
        if [[ -n "${img_ids}" ]]; then
            run docker rmi -f ${img_ids} || true
        else
            echo "(none to remove)"
        fi

        section "Removing volumes (DATA LOSS)"
        vol_names="$(docker volume ls -q || true)"
        if [[ -n "${vol_names}" ]]; then
            # shellcheck disable=SC2086
            run docker volume rm -f ${vol_names} || true
        else
            echo "(none to remove)"
        fi

        section "Pruning everything else"
        run docker system prune -a --volumes -f || true
        run docker builder prune -a -f || true
    fi

    # Networks: we can prune and attempt remove user networks; default networks cannot be removed.
    section "Pruning networks"
    run docker network prune -f || true

    section "State AFTER"
    docker_show_state

    echo
    echo "✅ Docker nuke complete."
}

# main
detect_engine

case "$ENGINE" in
    podman) podman_nuke ;;
    docker) docker_nuke ;;
    *)
        echo "ERROR: Unsupported engine: $ENGINE" >&2
        exit 2
        ;;
esac

exit 0
