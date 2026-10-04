#!/bin/sh
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2026-present, Ukama Inc.
set -eu
lab_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
if ! command -v node >/dev/null 2>&1; then
    echo "WEBAPP requires Node.js 22 or newer" >&2
    exit 1
fi
if [ ! -f "$lab_root/adapters/webapp/dist/cli.js" ]; then
    echo "WEBAPP build missing. From the lab directory, run:" >&2
    echo "  npm --prefix adapters/webapp ci" >&2
    echo "  npm --prefix adapters/webapp run build" >&2
    echo "  npm --prefix adapters/webapp run install:browser" >&2
    exit 1
fi
# Keep auth/profile/artifact paths consistent with future C-runner execution.
cd "$lab_root"
exec node "$lab_root/adapters/webapp/dist/cli.js" "$@"
