#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPO_ROOT="$(CDPATH= cd -- "${SCRIPT_DIR}/../../.." && pwd)"

if [[ -n "${UKAMA_APP_VERSION:-}" ]]; then
    VERSION="${UKAMA_APP_VERSION}"
else
    LATEST_TAG="$(
        git -C "${REPO_ROOT}" \
            -c safe.directory="${REPO_ROOT}" \
            describe --tags --abbrev=0
    )"
    COMMIT_HASH="$(
        git -C "${REPO_ROOT}" \
            -c safe.directory="${REPO_ROOT}" \
            rev-parse --short HEAD
    )"
    VERSION="${LATEST_TAG}-${COMMIT_HASH}"
fi

if [[ -z "${VERSION}" || "${VERSION}" == *$'\n'* ]]; then
    echo "Invalid application version" >&2
    exit 1
fi

if [[ "${1:-}" == "--print" ]]; then
    printf '%s\n' "${VERSION}"
    exit 0
fi

cat > version.h <<EOF_VERSION
#ifndef VERSION_H_
#define VERSION_H_

#define VERSION "${VERSION}"

#endif /* VERSION_H_ */
EOF_VERSION
