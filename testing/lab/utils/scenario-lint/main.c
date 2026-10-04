/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */
#include "scenario.h"

int main(int argc, char **argv) {
    return scenario_lint_main(argc - 1, argv + 1);
}
