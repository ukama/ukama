/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "nodes.h"

bool node_is_tower(Config *config) {

    if (!config || !config->nodeType) return false;

    return strcmp(config->nodeType, UKAMA_TOWER_NODE) == 0;
}

bool node_is_amplifier(Config *config) {

    if (!config || !config->nodeType) return false;

    return strcmp(config->nodeType, UKAMA_AMPLIFIER_NODE) == 0;
}

bool node_is_controller(Config *config) {

    if (!config || !config->nodeType) return false;

    return strcmp(config->nodeType, UKAMA_CONTROLLER_NODE) == 0;
}
