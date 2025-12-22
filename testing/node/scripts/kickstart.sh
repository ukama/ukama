#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

set -euo pipefail

echo "Starting on-boot group..."
supervisorctl start "on-boot:*"

echo "Waiting for noded to be RUNNING..."
while ! supervisorctl status noded_latest | grep -q 'RUNNING'; do
  sleep 2
done

echo "Starting bootstrap..."
supervisorctl start bootstrap_latest

echo "Waiting for bootstrap to EXIT..."
while ! supervisorctl status bootstrap_latest | grep -q 'EXITED'; do
    sleep 2
done

# Make sure bootstrap exited successfully
if ! supervisorctl status bootstrap_latest | grep -q '(exit status 0;'; then
    echo "ERROR: bootstrap failed:"
    supervisorctl status bootstrap_latest || true
    exit 1
fi

echo "Starting meshd..."
supervisorctl start meshd_latest

echo "Waiting for meshd to be RUNNING..."
while ! supervisorctl status meshd_latest | grep -q 'RUNNING'; do
    sleep 2
done

echo "Starting sys-service group..."
supervisorctl start "sys-service:*" || true

echo "Kickstart complete."
