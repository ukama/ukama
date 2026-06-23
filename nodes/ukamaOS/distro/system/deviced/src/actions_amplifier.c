/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "actions.h"
#include "deviced.h"
#include "web_client.h"

#include "usys_log.h"
#include "usys_services.h"

static int write_radio_emu(ControlState desired) {

    FILE *fp;

    fp = fopen(DEF_RADIO_EMU_FILE, "w");
    if (!fp) {
        usys_log_error("radio: failed to open emu state file: %s", DEF_RADIO_EMU_FILE);
        return STATUS_NOK;
    }

    fprintf(fp, "%s\n", desired == CONTROL_STATE_ON ? "on" : "off");
    fclose(fp);
    return STATUS_OK;
}

int actions_radio_apply(Config *config, ControlState desired) {

    int retCode;

    if (!config) return STATUS_NOK;

    if (config->clientMode) {
        usys_log_info("radio: %s (client emu)",
                      desired == CONTROL_STATE_ON ? "on" : "off");
        return write_radio_emu(desired);
    }

    if (!config->nodeType) return STATUS_NOK;

    if (strcmp(config->nodeType, UKAMA_TOWER_NODE) == 0) {
        usys_log_info("radio: %s (tower trx client)",
                      desired == CONTROL_STATE_ON ? "on" : "off");
        retCode = -1;
        return wc_send_radio_to_client(config, desired, &retCode) == USYS_OK ?
               STATUS_OK : STATUS_NOK;
    }

    if (strcmp(config->nodeType, UKAMA_AMPLIFIER_NODE) != 0) {
        usys_log_error("radio: unsupported node type: %s", config->nodeType);
        return STATUS_NOK;
    }

    usys_log_info("radio: %s (via femd:%d)",
                  desired == CONTROL_STATE_ON ? "on" : "off",
                  config->femPort);

    if (desired == CONTROL_STATE_OFF) {
        retCode = -1;
        if (wc_put_gpio_to_femd(config, 1, desired, &retCode) != STATUS_OK) {
            usys_log_error("radio: femd gpio apply failed fem=1 http=%d", retCode);
            return STATUS_NOK;
        }

        retCode = -1;
        if (wc_put_gpio_to_femd(config, 2, desired, &retCode) != STATUS_OK) {
            usys_log_error("radio: femd gpio apply failed fem=2 http=%d", retCode);
            return STATUS_NOK;
        }
        return STATUS_OK;
    }

    retCode = -1;
    if (wc_put_gpio_to_femd(config, 1, desired, &retCode) != STATUS_OK) {
        usys_log_error("radio: femd gpio apply failed fem=1 http=%d", retCode);
        return STATUS_NOK;
    }

    retCode = -1;
    if (wc_put_gpio_to_femd(config, 2, desired, &retCode) != STATUS_OK) {
        usys_log_error("radio: femd gpio apply failed fem=2 http=%d", retCode);
        return STATUS_NOK;
    }

    return STATUS_OK;
}
