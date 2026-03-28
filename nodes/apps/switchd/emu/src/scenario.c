/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "scenario.h"
#include "model.h"

int scenario_apply(EmuModel *model, const char *name) {
    size_t i = 0;

    snprintf(model->activeScenario, sizeof(model->activeScenario), "%s",
             (name != NULL) ? name : DEF_SCENARIO);

    model->faults.snmpDelayMs = 0;
    model->faults.unreachable = 0;
    model->faults.flapPortId  = 0;
    model->faults.flapPeriodSec = 0;
    model->faults.tftpFail    = 0;
    model->faults.snmpSetFail = 0;
    model->firmware.applyShouldFail = 0;
    model->info.systemTempC = 40;
    model_set_reachable(model, 1);

    for (i = 0; i < model->portCount; i++) {
        model->ports[i].faultPoe  = 0;
        model->ports[i].faultLink = 0;
    }

    if (name == NULL || strcmp(name, "normal") == 0) {
        ;
    } else if (strcmp(name, "high_temp") == 0) {
        model->info.systemTempC = 78;
    } else if (strcmp(name, "poe_fault_port4") == 0) {
        model->ports[3].faultPoe = 1;
    } else if (strcmp(name, "firmware_fail") == 0) {
        model->firmware.applyShouldFail = 1;
    } else if (strcmp(name, "switch_down") == 0) {
        model->faults.unreachable = 1;
        model_set_reachable(model, 0);
    } else if (strcmp(name, "slow_snmp") == 0) {
        model->faults.snmpDelayMs = 1200;
    } else if (strcmp(name, "link_flap_port3") == 0) {
        model->faults.flapPortId    = 3;
        model->faults.flapPeriodSec = 4;
    }

    model_recompute(model);
    return STATUS_OK;
}
