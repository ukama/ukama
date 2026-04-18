#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -e

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPO_ROOT="$(CDPATH= cd -- "${SCRIPT_DIR}/../../.." && pwd)"

LATEST_TAG="$(git -C "${REPO_ROOT}" -c safe.directory="${REPO_ROOT}" \
    describe --tags --abbrev=0)"
COMMIT_HASH="$(git -C "${REPO_ROOT}" -c safe.directory="${REPO_ROOT}" \
    rev-parse --short HEAD)"
VERSION="${LATEST_TAG}-${COMMIT_HASH}"

if [ "$1" = "--print" ]; then
    echo "${VERSION}"
    exit 0
fi

cat > version.h <<EOF
#ifndef VERSION_H_
#define VERSION_H_
#define VERSION "${VERSION}"
#endif /* VERSION_H_ */
EOF
