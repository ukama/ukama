#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.
set -eu
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
action=$1
shift
case "$action" in
  attach)
    [ "$#" -eq 8 ] || exit 2
    "$script_dir/start-ue.sh" "$@"
    "$script_dir/wait-ues-attached.sh" "$3" "$8"
    ;;
  detach)
    [ "$#" -eq 3 ] || exit 2
    ue_id=$1
    run_dir=$2
    imsi=$3
    case "$imsi" in ''|*[!0-9]*) echo 'invalid owned IMSI' >&2; exit 2 ;; esac
    # Stop the retrying agent before closing the EPC session. Otherwise it
    # could immediately reattach while the target population is decreasing.
    state="$run_dir/runtime-ues/$ue_id.env"
    if [ ! -f "$state" ]; then
      # start-ue writes its state after container startup. Recover a partial
      # start by the exact IMSI already journaled before invoking the provider.
      if podman container exists "ue-$imsi"; then
        podman rm -f "ue-$imsi" >/dev/null
        echo "removed owned UE container, but missing state prevents session finalization: $ue_id" >&2
        exit 1
      fi
      exit 0
    fi
    if [ -f "$state" ]; then
      . "$state"
      if [ -n "${UE_CONTAINER:-}" ] && podman container exists "$UE_CONTAINER"; then
        podman stop -t 3 "$UE_CONTAINER" >/dev/null
      fi
    fi
    "$script_dir/detach-ue.sh" "$ue_id" "$run_dir"
    "$script_dir/cleanup-ue.sh" "$ue_id" "$run_dir"
    ;;
  probe)
    exec "$script_dir/verify-ue-session.sh" "$@"
    ;;
  *) echo "unknown workload UE action: $action" >&2; exit 2 ;;
esac
