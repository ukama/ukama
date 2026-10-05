#!/bin/sh
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2026-present, Ukama Inc.
# Standalone host build of the same C loader/validators used by ukama-lab.
# No BFF, Podman, browser, ukamaOS tree, or third-party headers are needed.
set -eu
lab_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
lint_build=$(mktemp -d "${TMPDIR:-/tmp}/ukama-lab-lint.XXXXXX")
trap 'rm -rf "$lint_build"' EXIT HUP INT TERM
printf '#define VERSION "scenario-lint"\n' > "$lint_build/version.h"
"${CC:-cc}" -std=gnu11 -D_POSIX_C_SOURCE=200809L \
    -Wall -Wextra -Werror -Wdeclaration-after-statement -O1 \
    -I"$lint_build" -I"$lab_root/inc" \
    "$lab_root/utils/scenario-lint/main.c" \
    "$lab_root/src/scenario.c" "$lab_root/src/validate.c" \
    "$lab_root/src/scenario_webapp.c" "$lab_root/src/scenario_lint.c" \
    "$lab_root/src/util.c" -o "$lint_build/scenario-lint"
"$lint_build/scenario-lint" "$@"
