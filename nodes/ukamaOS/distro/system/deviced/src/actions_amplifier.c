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
#include "usys_services.h"

int actions_radio_apply(Config *config, ControlState desired) {

    int ret;
    int retCode;

    ret = STATUS_NOK;
    retCode = 0;

    (void)config;

    usys_log_info("radio: %s (via femd:%d)",
                  (desired == CONTROL_STATE_ON) ? "on" : "off",
                  config->femPort);

    ret = wc_put_gpio_to_femd(config, 1, desired, &retCode);
    if (ret != STATUS_OK) {
        usys_log_error("radio: femd gpio apply failed fem=1 http=%d", retCode);
        return ret;
    }

    ret = wc_put_gpio_to_femd(config, 2, desired, &retCode);
    if (ret != STATUS_OK) {
        usys_log_error("radio: femd gpio apply failed fem=2 http=%d", retCode);
        return ret;
    }

    return STATUS_OK;
}
