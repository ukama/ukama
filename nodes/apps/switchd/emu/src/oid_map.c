/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "model.h"
#include "oid_map.h"

#define OID_SERIAL       "1.3.6.1.4.1.12284.5.1.7.0"
#define OID_MANUFACTURER "1.3.6.1.4.1.12284.5.1.8.0"
#define OID_HWVER        "1.3.6.1.4.1.12284.5.1.9.0"
#define OID_SWVER        "1.3.6.1.4.1.12284.5.1.10.0"
#define OID_EXECUTE      "1.3.6.1.4.1.12284.5.1.14.0"
#define OID_EXEC_STATUS  "1.3.6.1.4.1.12284.5.1.15.0"
#define OID_POE_USED     "1.3.6.1.4.1.12284.5.2.2.0"
#define OID_POE_BUDGET   "1.3.6.1.4.1.12284.5.2.3.0"
#define OID_SYS_TEMP     "1.3.6.1.4.1.12284.5.6.3.0"
#define OID_AMB_TEMP     "1.3.6.1.4.1.12284.5.6.6.0"
#define OID_IN_VOLT      "1.3.6.1.4.1.12284.5.6.13.0"
#define OID_SYS_CURR     "1.3.6.1.4.1.12284.5.6.14.0"
#define OID_SYS_POWER    "1.3.6.1.4.1.12284.5.6.17.0"
#define OID_ALARM_LINK   "1.3.6.1.4.1.12284.5.6.28.0"
#define OID_ALARM_POE    "1.3.6.1.4.1.12284.5.6.29.0"

static int port_index(const char *oid, const char *prefix) {
    size_t len = strlen(prefix);

    if (strncmp(oid, prefix, len) != 0) {
        return -1;
    }

    return atoi(oid + len);
}

int oid_get_string(EmuModel *model, const char *oid, char *buf, size_t len) {
    if (strcmp(oid, OID_SERIAL) == 0) {
        snprintf(buf, len, "%s", model->info.serial);
    } else if (strcmp(oid, OID_MANUFACTURER) == 0) {
        snprintf(buf, len, "%s", model->info.manufacturer);
    } else if (strcmp(oid, OID_HWVER) == 0) {
        snprintf(buf, len, "%s", model->info.hardwareVersion);
    } else if (strcmp(oid, OID_SWVER) == 0) {
        snprintf(buf, len, "%s", model->info.softwareVersion);
    } else {
        return STATUS_NOK;
    }

    return STATUS_OK;
}

int oid_get_int(EmuModel *model, const char *oid, int *value) {
    int idx = 0;

    if (strcmp(oid, OID_EXEC_STATUS) == 0) {
        *value = model->firmware.executeStatus;
    } else if (strcmp(oid, OID_POE_USED) == 0) {
        *value = model->info.poeUsedMw;
    } else if (strcmp(oid, OID_POE_BUDGET) == 0) {
        *value = model->info.poeBudgetMw;
    } else if (strcmp(oid, OID_SYS_TEMP) == 0) {
        *value = model->info.systemTempC;
    } else if (strcmp(oid, OID_AMB_TEMP) == 0) {
        *value = model->info.ambientTempC;
    } else if (strcmp(oid, OID_IN_VOLT) == 0) {
        *value = model->info.inputVoltageMv;
    } else if (strcmp(oid, OID_SYS_CURR) == 0) {
        *value = model->info.systemCurrentMa;
    } else if (strcmp(oid, OID_SYS_POWER) == 0) {
        *value = model->info.systemPowerMw;
    } else if (strcmp(oid, OID_ALARM_LINK) == 0) {
        *value = model->info.alarmLinkFailure;
    } else if (strcmp(oid, OID_ALARM_POE) == 0) {
        *value = model->info.alarmPoeFailure;
    } else if ((idx = port_index(oid, "1.3.6.1.4.1.12284.5.2.1.1.2.")) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].poeSupported;
    } else if ((idx = port_index(oid, "1.3.6.1.4.1.12284.5.2.1.1.3.")) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].poeAdminEnabled;
    } else if ((idx = port_index(oid, "1.3.6.1.4.1.12284.5.2.1.1.4.")) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].poeOperStatus;
    } else if ((idx = port_index(oid, "1.3.6.1.4.1.12284.5.2.1.1.5.")) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].poePowerMw;
    } else if ((idx = port_index(oid, "1.3.6.1.4.1.12284.5.2.1.1.6.")) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].poeCurrentMa;
    } else if ((idx = port_index(oid, "1.3.6.1.4.1.12284.5.2.1.1.7.")) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].poeVoltageMv;
    } else if ((idx = port_index(oid, "1.3.6.1.4.1.12284.5.2.1.1.8.")) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].poeClassId;
    } else {
        return STATUS_NOK;
    }

    return STATUS_OK;
}

int oid_set_int(EmuModel *model, const char *oid, int value) {
    int idx = 0;

    if (model->faults.snmpSetFail) {
        return STATUS_NOK;
    }

    if (strcmp(oid, OID_EXECUTE) == 0) {
        if (model->firmware.state == FW_STAGED) {
            model->firmware.state        = FW_APPLYING;
            model->firmware.executeStatus = 0;
            model->firmware.stateSince   = time(NULL);
            return STATUS_OK;
        }
        return STATUS_NOK;
    }

    idx = port_index(oid, "1.3.6.1.4.1.12284.5.2.1.1.3.");
    if (idx > 0 && idx <= 8) {
        model->ports[idx - 1].poeAdminEnabled = value ? 1 : 0;
        model_recompute(model);
        return STATUS_OK;
    }

    return STATUS_NOK;
}
