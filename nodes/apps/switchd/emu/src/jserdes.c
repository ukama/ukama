/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdarg.h>
#include <stdio.h>
#include <string.h>

#include "jserdes.h"
#include "json_types.h"
#include "utils.h"

static int append(char *buf, size_t len, size_t *off, const char *fmt, ...) {
    va_list args;
    int written = 0;

    if (*off >= len) {
        return STATUS_NOK;
    }

    va_start(args, fmt);
    written = vsnprintf(buf + *off, len - *off, fmt, args);
    va_end(args);

    if (written < 0 || (size_t)written >= (len - *off)) {
        return STATUS_NOK;
    }

    *off += (size_t)written;
    return STATUS_OK;
}

int json_serialize_state(const EmuModel *model, char *buf, size_t len) {
    size_t off = 0;

    buf[0] = '\0';

    if (append(buf, len, &off,
               "{\"%s\":\"%s\",\"%s\":%s,\"%s\":\"%s\","
               "\"%s\":%d,\"%s\":%d,\"%s\":%d}",
               JTAG_SCENARIO, model->activeScenario,
               JTAG_REACHABLE, bool_to_json(model->info.reachable),
               JTAG_SOFTWAREVER, model->info.softwareVersion,
               JTAG_SYSTEMTEMPC, model->info.systemTempC,
               JTAG_POEUSEDMW, model->info.poeUsedMw,
               JTAG_POEBUDGETMW, model->info.poeBudgetMw) != STATUS_OK) {
        return STATUS_NOK;
    }

    return STATUS_OK;
}

int json_serialize_ports(const EmuModel *model, char *buf, size_t len) {
    size_t i   = 0;
    size_t off = 0;

    buf[0] = '\0';
    if (append(buf, len, &off, "[") != STATUS_OK) {
        return STATUS_NOK;
    }

    for (i = 0; i < model->portCount; i++) {
        if (append(buf, len, &off,
                   "%s{\"%s\":%u,\"%s\":%s,\"%s\":%s,"
                   "\"%s\":%s,\"%s\":%d}",
                   (i == 0) ? "" : ",",
                   JTAG_ID, model->ports[i].id,
                   JTAG_LINKUP, bool_to_json(model->ports[i].linkUp),
                   JTAG_ADMINUP, bool_to_json(model->ports[i].adminUp),
                   JTAG_POEADMINENABLED,
                   bool_to_json(model->ports[i].poeAdminEnabled),
                   JTAG_POEPOWERMW, model->ports[i].poePowerMw) != STATUS_OK) {
            return STATUS_NOK;
        }
    }

    if (append(buf, len, &off, "]") != STATUS_OK) {
        return STATUS_NOK;
    }

    return STATUS_OK;
}

int json_serialize_firmware(const EmuModel *model, char *buf, size_t len) {
    size_t off = 0;

    buf[0] = '\0';
    if (append(buf, len, &off,
               "{\"%s\":%d,\"%s\":\"%s\",\"%s\":%s,\"%s\":%d}",
               JTAG_STATE, model->firmware.state,
               JTAG_FILENAME, model->firmware.stagedFilename,
               JTAG_APPLYSHOULDFAIL,
               bool_to_json(model->firmware.applyShouldFail),
               JTAG_EXECUTESTATUS, model->firmware.executeStatus) != STATUS_OK) {
        return STATUS_NOK;
    }

    return STATUS_OK;
}

int json_serialize_result_ok(char *buf, size_t len) {
    return snprintf(buf, len, "{\"%s\":\"ok\"}", JTAG_RESULT) > 0 ?
           STATUS_OK : STATUS_NOK;
}

int json_serialize_error(const char *err, char *buf, size_t len) {
    return snprintf(buf, len, "{\"%s\":\"%s\"}", JTAG_ERROR, err) > 0 ?
           STATUS_OK : STATUS_NOK;
}
