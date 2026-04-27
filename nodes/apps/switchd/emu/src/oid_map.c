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

/* Standard IF-MIB OID prefixes */
#define OID_IF_NUMBER          "1.3.6.1.2.1.2.1.0"
#define OID_IF_DESCR_PFX       "1.3.6.1.2.1.2.2.1.2."
#define OID_IF_SPEED_PFX       "1.3.6.1.2.1.2.2.1.5."
#define OID_IF_ADMIN_PFX       "1.3.6.1.2.1.2.2.1.7."
#define OID_IF_OPER_PFX        "1.3.6.1.2.1.2.2.1.8."
#define OID_IF_IN_OCT_PFX      "1.3.6.1.2.1.2.2.1.10."
#define OID_IF_IN_UCAST_PFX    "1.3.6.1.2.1.2.2.1.11."
#define OID_IF_IN_DISC_PFX     "1.3.6.1.2.1.2.2.1.13."
#define OID_IF_IN_ERR_PFX      "1.3.6.1.2.1.2.2.1.14."
#define OID_IF_OUT_OCT_PFX     "1.3.6.1.2.1.2.2.1.16."
#define OID_IF_OUT_UCAST_PFX   "1.3.6.1.2.1.2.2.1.17."
#define OID_IF_OUT_DISC_PFX    "1.3.6.1.2.1.2.2.1.19."
#define OID_IF_OUT_ERR_PFX     "1.3.6.1.2.1.2.2.1.20."
#define OID_IF_NAME_PFX        "1.3.6.1.2.1.31.1.1.1.1."
#define OID_IF_HC_IN_PFX       "1.3.6.1.2.1.31.1.1.1.6."
#define OID_IF_HC_OUT_PFX      "1.3.6.1.2.1.31.1.1.1.10."

/* Private Tycon MIB OID prefixes */
#define OID_POE_EXIST_PFX  "1.3.6.1.4.1.12284.5.2.1.1.2."
#define OID_POE_ADMIN_PFX  "1.3.6.1.4.1.12284.5.2.1.1.3."
#define OID_POE_OPER_PFX   "1.3.6.1.4.1.12284.5.2.1.1.4."
#define OID_POE_POWER_PFX  "1.3.6.1.4.1.12284.5.2.1.1.5."
#define OID_POE_CURR_PFX   "1.3.6.1.4.1.12284.5.2.1.1.6."
#define OID_POE_VOLT_PFX   "1.3.6.1.4.1.12284.5.2.1.1.7."
#define OID_POE_CLASS_PFX  "1.3.6.1.4.1.12284.5.2.1.1.8."

static int port_index(const char *oid, const char *prefix) {
    size_t len = strlen(prefix);

    if (strncmp(oid, prefix, len) != 0) {
        return -1;
    }

    return atoi(oid + len);
}

int oid_get_string(EmuModel *model, const char *oid, char *buf, size_t len) {

    int idx;

    if (!model || !oid || !buf || len == 0) {
        return STATUS_NOK;
    }

    model_recompute(model);

    if (strcmp(oid, OID_SERIAL) == 0) {
        snprintf(buf, len, "%s", model->info.serial);

    } else if (strcmp(oid, OID_MANUFACTURER) == 0) {
        snprintf(buf, len, "%s", model->info.manufacturer);

    } else if (strcmp(oid, OID_HWVER) == 0) {
        snprintf(buf, len, "%s", model->info.hardwareVersion);

    } else if (strcmp(oid, OID_SWVER) == 0) {
        snprintf(buf, len, "%s", model->info.softwareVersion);

    } else if (strcmp(oid, OID_SYS_TEMP) == 0) {
        snprintf(buf, len, "%d C", model->info.systemTempC);

    } else if (strcmp(oid, OID_AMB_TEMP) == 0) {
        snprintf(buf, len, "%d C", model->info.ambientTempC);

    } else if (strcmp(oid, OID_IN_VOLT) == 0) {
        snprintf(buf, len, "%.2f V",
                 ((double)model->info.inputVoltageMv) / 1000.0);

    } else if (strcmp(oid, OID_SYS_CURR) == 0) {
        snprintf(buf, len, "%.2f A",
                 ((double)model->info.systemCurrentMa) / 1000.0);

    } else if (strcmp(oid, OID_SYS_POWER) == 0) {
        snprintf(buf, len, "%.2f W",
                 ((double)model->info.systemPowerMw) / 1000.0);

    } else if ((idx = port_index(oid, OID_IF_NAME_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        snprintf(buf, len, "%s", model->ports[idx - 1].name);

    } else if ((idx = port_index(oid, OID_IF_DESCR_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        snprintf(buf, len, "%s", model->ports[idx - 1].name);

    } else {
        return STATUS_NOK;
    }

    return STATUS_OK;
}

int oid_get_int(EmuModel *model, const char *oid, int *value) {

    int idx;

    idx = 0;

    if (!model || !oid || !value) {
        return STATUS_NOK;
    }

    model_recompute(model);

    if (strcmp(oid, OID_IF_NUMBER) == 0) {
        *value = (int)model->portCount;

    } else if (strcmp(oid, OID_EXEC_STATUS) == 0) {
        *value = model->firmware.executeStatus;

    } else if (strcmp(oid, OID_POE_USED) == 0) {
        /*
         * switch.d expects this KPI in watts.
         * Emulator keeps internal value in milliwatts.
         */
        *value = (model->info.poeUsedMw + 500) / 1000;

    } else if (strcmp(oid, OID_POE_BUDGET) == 0) {
        /*
         * switch.d expects this KPI in watts.
         * Emulator keeps internal value in milliwatts.
         */
        *value = (model->info.poeBudgetMw + 500) / 1000;

    } else if (strcmp(oid, OID_ALARM_LINK) == 0) {
        *value = model->info.alarmLinkFailure;

    } else if (strcmp(oid, OID_ALARM_POE) == 0) {
        *value = model->info.alarmPoeFailure;

    } else if ((idx = port_index(oid, OID_IF_SPEED_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = (int)(model->ports[idx - 1].speedMbps * 1000000U);

    } else if ((idx = port_index(oid, OID_IF_ADMIN_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].adminUp ? 1 : 2;

    } else if ((idx = port_index(oid, OID_IF_OPER_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].linkUp ? 1 : 2;

    } else if ((idx = port_index(oid, OID_IF_IN_OCT_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = (int)(model->ports[idx - 1].rxBytes & 0xFFFFFFFF);

    } else if ((idx = port_index(oid, OID_IF_OUT_OCT_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = (int)(model->ports[idx - 1].txBytes & 0xFFFFFFFF);

    } else if ((idx = port_index(oid, OID_IF_IN_UCAST_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = (int)(model->ports[idx - 1].rxPackets & 0xFFFFFFFF);

    } else if ((idx = port_index(oid, OID_IF_OUT_UCAST_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = (int)(model->ports[idx - 1].txPackets & 0xFFFFFFFF);

    } else if ((idx = port_index(oid, OID_IF_IN_ERR_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = (int)(model->ports[idx - 1].rxErrors & 0xFFFFFFFF);

    } else if ((idx = port_index(oid, OID_IF_OUT_ERR_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = (int)(model->ports[idx - 1].txErrors & 0xFFFFFFFF);

    } else if ((idx = port_index(oid, OID_IF_IN_DISC_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = 0;

    } else if ((idx = port_index(oid, OID_IF_OUT_DISC_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = 0;

    } else if ((idx = port_index(oid, OID_IF_HC_IN_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = (int)(model->ports[idx - 1].rxBytes & 0xFFFFFFFF);

    } else if ((idx = port_index(oid, OID_IF_HC_OUT_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = (int)(model->ports[idx - 1].txBytes & 0xFFFFFFFF);

    } else if ((idx = port_index(oid, OID_POE_EXIST_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].poeSupported;

    } else if ((idx = port_index(oid, OID_POE_ADMIN_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].poeAdminEnabled;

    } else if ((idx = port_index(oid, OID_POE_OPER_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        *value = model->ports[idx - 1].poeOperStatus;

    } else if ((idx = port_index(oid, OID_POE_POWER_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        /*
         * switch.d expects port PoE power in watts.
         * Emulator keeps it in milliwatts.
         */
        *value = (model->ports[idx - 1].poePowerMw + 500) / 1000;

    } else if ((idx = port_index(oid, OID_POE_CURR_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        /*
         * switch.d expects PoE current in mA and divides by 1000.
         */
        *value = model->ports[idx - 1].poeCurrentMa;

    } else if ((idx = port_index(oid, OID_POE_VOLT_PFX)) > 0 &&
               idx <= (int)model->portCount) {
        /*
         * switch.d expects PoE voltage in decivolts and divides by 10.
         * Emulator keeps it in millivolts.
         */
        *value = model->ports[idx - 1].poeVoltageMv / 100;

    } else if ((idx = port_index(oid, OID_POE_CLASS_PFX)) > 0 &&
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

    idx = port_index(oid, OID_POE_ADMIN_PFX);
    if (idx > 0 && idx <= 8) {
        model->ports[idx - 1].poeAdminEnabled = value ? 1 : 0;
        model_recompute(model);
        return STATUS_OK;
    }

    return STATUS_NOK;
}

int oid_get_next(EmuModel *model, const char *oid, char *nextOid, size_t nextLen) {
    static const char *portPrefixes[] = {
        OID_IF_DESCR_PFX,
        OID_IF_SPEED_PFX,
        OID_IF_ADMIN_PFX,
        OID_IF_OPER_PFX,
        OID_IF_IN_OCT_PFX,
        OID_IF_IN_UCAST_PFX,
        OID_IF_IN_DISC_PFX,
        OID_IF_IN_ERR_PFX,
        OID_IF_OUT_OCT_PFX,
        OID_IF_OUT_UCAST_PFX,
        OID_IF_OUT_DISC_PFX,
        OID_IF_OUT_ERR_PFX,
        OID_IF_NAME_PFX,
        OID_IF_HC_IN_PFX,
        OID_IF_HC_OUT_PFX,
        OID_POE_EXIST_PFX,
        OID_POE_ADMIN_PFX,
        OID_POE_OPER_PFX,
        OID_POE_POWER_PFX,
        OID_POE_CURR_PFX,
        OID_POE_VOLT_PFX,
        OID_POE_CLASS_PFX,
        NULL
    };

    static const char *scalarOids[] = {
        OID_IF_NUMBER,
        OID_SERIAL,
        OID_MANUFACTURER,
        OID_HWVER,
        OID_SWVER,
        OID_EXECUTE,
        OID_EXEC_STATUS,
        OID_POE_USED,
        OID_POE_BUDGET,
        OID_SYS_TEMP,
        OID_AMB_TEMP,
        OID_IN_VOLT,
        OID_SYS_CURR,
        OID_SYS_POWER,
        OID_ALARM_LINK,
        OID_ALARM_POE,
        NULL
    };

    char candidates[1024][128];
    int count = 0;

    for (int i = 0; scalarOids[i] != NULL && count < 1024; i++) {
        snprintf(candidates[count++], 128, "%s", scalarOids[i]);
    }

    for (int i = 0; portPrefixes[i] != NULL; i++) {
        for (int p = 1; p <= (int)model->portCount && count < 1024; p++) {
            snprintf(candidates[count++], 128, "%s%d", portPrefixes[i], p);
        }
    }

    const char *best = NULL;

    for (int i = 0; i < count; i++) {
        const char *a = oid;
        const char *b = candidates[i];
        const char *pa = a, *pb = b;
        int cmp = 0;

        while (*pa && *pb && cmp == 0) {
            long va = strtol(pa, (char **)&pa, 10);
            long vb = strtol(pb, (char **)&pb, 10);
            if (va < vb) cmp = -1;
            else if (va > vb) cmp = 1;
            if (*pa == '.') pa++;
            if (*pb == '.') pb++;
        }
        if (cmp == 0) {
            if (*pb) cmp = -1;
            else if (*pa) cmp = 1;
        }

        if (cmp < 0) {
            if (best == NULL) {
                best = candidates[i];
            } else {
                const char *x = candidates[i], *y = best;
                const char *px = x, *py = y;
                int c2 = 0;
                while (*px && *py && c2 == 0) {
                    long vx = strtol(px, (char **)&px, 10);
                    long vy = strtol(py, (char **)&py, 10);
                    if (vx < vy) c2 = -1;
                    else if (vx > vy) c2 = 1;
                    if (*px == '.') px++;
                    if (*py == '.') py++;
                }
                if (c2 == 0) {
                    if (*py) c2 = -1;
                    else if (*px) c2 = 1;
                }
                if (c2 < 0) best = candidates[i];
            }
        }
    }

    if (best == NULL) return STATUS_NOK;
    snprintf(nextOid, nextLen, "%s", best);
    return STATUS_OK;
}
