#!/usr/bin/env bash
set -euo pipefail

# podman-nuke.sh
#
# !!! DANGEROUS !!!
# This script DESTROYS ALL podman state for the current user:
# - containers (running + stopped)
# - pods
# - images
# - volumes (DATA LOSS)
# - networks
# - build cache
#
# SAFETY:
#   You MUST pass --all or the script will exit.
#
# Usage:
#   ./podman-nuke.sh --all
#   ./podman-nuke.sh --all --dry-run
#   ./podman-nuke.sh --all --no-reset

CONFIRM_ALL=0
DRY_RUN=0
NO_RESET=0

for arg in "${@:-}"; do
  case "$arg" in
    --all) CONFIRM_ALL=1 ;;
    --dry-run) DRY_RUN=1 ;;
    --no-reset) NO_RESET=1 ;;
    -h|--help)
      cat <<'EOF'
podman-nuke.sh — DESTROY ALL podman resources (DANGEROUS)

REQUIRED:
  --all        Required safety flag. Script will NOT run without it.

OPTIONAL:
  --dry-run    Show what would be deleted, do nothing
  --no-reset   Do not use "podman system reset -f"
               (still removes containers, images, volumes, cache)

Examples:
  ./podman-nuke.sh --all
  ./podman-nuke.sh --all --dry-run
  ./podman-nuke.sh --all --no-reset
EOF
      exit 0
      ;;
    *)
      echo "Unknown option: $arg" >&2
      exit 2
      ;;
  esac
done

if [[ "$CONFIRM_ALL" -ne 1 ]]; then
  echo "ERROR: Refusing to run."
  echo "This script DESTROYS ALL podman data."
  echo
  echo "You must explicitly pass: --all"
  echo "Example:"
  echo "  ./podman-nuke.sh --all"
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

section "Podman sanity check"
command -v podman >/dev/null || { echo "podman not found"; exit 1; }
podman --version

section "State BEFORE"
echo "== Containers =="
podman ps -a || true
echo
echo "== Pods =="
podman pod ls -a || true
echo
echo "== Images =="
podman images -a || true
echo
echo "== Volumes =="
podman volume ls || true
echo
echo "== Networks =="
podman network ls 2>/dev/null || true

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
podman ps -a || true
podman pod ls -a || true
podman images -a || true
podman volume ls || true
podman network ls 2>/dev/null || true

echo
echo "✅ Podman nuke complete."
