/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "actions.h"
#include "deviced.h"
#include "web_client.h"

#include "usys_log.h"

int actions_service_apply(Config *config, ControlState desired) {

    int retCode = -1;

    if (!config) {
        return STATUS_NOK;
    }

    if (config->clientMode || !config->nodeType ||
        strcmp(config->nodeType, UKAMA_TOWER_NODE) != 0) {
        usys_log_error("service: unsupported on this device.d instance");
        return STATUS_NOK;
    }

    usys_log_info("service: %s", desired == CONTROL_STATE_ON ? "on" : "off");

    return wc_post_service_to_pcrf(config, desired, &retCode);
}
