#!/usr/bin/env bash
# Merge local worker output and reconcile it against the submitted scenario plan.
set -Eeuo pipefail
SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
[[ $# == 1 && -d "$1" ]] || { echo "usage: $0 LOCAL_BATCH_DIRECTORY" >&2; exit 2; }
exec python3 "$SCRIPT_DIR/report.py" "$1"
