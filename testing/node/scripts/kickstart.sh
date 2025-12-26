#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

set -euo pipefail

CONF="/etc/supervisor.conf"
SUPERVISORCTL=(supervisorctl -c "$CONF")

# IMPORTANT: noded_latest is part of the group "on-boot"
NODED="on-boot:noded_latest"

ctl() {
  "${SUPERVISORCTL[@]}" "$@"
}

status_line() {
  local name="$1"
  ctl status "$name" 2>/dev/null || true
}

is_state() {
  local name="$1"
  local want="$2"
  ctl status "$name" 2>/dev/null | awk '{print $2}' | grep -qx "$want"
}

wait_state() {
  local name="$1"
  local want="$2"
  local interval="${3:-1}"

  echo "Waiting for ${name} to be ${want}..."
  while ! is_state "$name" "$want"; do
    status_line "$name"
    sleep "$interval"
  done
}

wait_exited_ok() {
  local name="$1"
  local interval="${2:-1}"

  echo "Waiting for ${name} to EXIT..."
  while ! is_state "$name" "EXITED"; do
    status_line "$name"
    sleep "$interval"
  done

  if ! ctl status "$name" 2>/dev/null | grep -q 'exit status 0'; then
    echo "ERROR: ${name} exited but not successfully:"
    status_line "$name"
    exit 1
  fi
}

start_prog() {
  local name="$1"
  echo "Starting ${name}..."
  ctl start "$name"
}

start_group() {
  local group="$1"
  echo "Starting ${group} group..."
  ctl start "${group}:*"
}

echo "Kickstart using supervisorctl config: $CONF"
echo "Supervisor status (sanity):"
ctl status || true

echo "Starting on-boot group..."
start_group "on-boot"

wait_state "$NODED" "RUNNING" 1

echo "Starting bootstrap..."
start_prog "bootstrap_latest"
wait_exited_ok "bootstrap_latest" 1

echo "Starting meshd..."
start_prog "meshd_latest"
wait_state "meshd_latest" "RUNNING" 1

echo "Starting sys-service group..."
ctl start "sys-service:*" || true

echo "Kickstart complete."
