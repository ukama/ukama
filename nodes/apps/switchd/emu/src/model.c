/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <time.h>

#include "model.h"

void model_recompute(EmuModel *model) {
    size_t i     = 0;
    int poeUsed  = 0;
    int poeAlarm = 0;
    int linkAlarm = 0;

    for (i = 0; i < model->portCount; i++) {
        EmuPortState *port = &model->ports[i];

        if (port->poeSupported && port->poeAdminEnabled && !port->faultPoe) {
            port->poeOperStatus = 1;
            port->poePowerMw    = 8000 + (int)port->id * 250;
            port->poeCurrentMa  = 180 + (int)port->id * 10;
            port->poeVoltageMv  = 52000;
            poeUsed            += port->poePowerMw;
        } else {
            port->poeOperStatus = port->faultPoe ? 2 : 0;
            port->poePowerMw    = 0;
            port->poeCurrentMa  = 0;
            port->poeVoltageMv  = port->poeSupported ? 52000 : 0;
        }

        if (port->faultPoe) {
            poeAlarm = 1;
        }

        if (port->faultLink || (port->adminUp && !port->linkUp && port->id <= 8)) {
            linkAlarm = 1;
        }
    }

    model->info.poeUsedMw         = poeUsed;
    model->info.alarmPoeFailure   = poeAlarm || (poeUsed > model->info.poeBudgetMw);
    model->info.alarmLinkFailure  = linkAlarm;
    model->info.systemPowerMw     = 25000 + poeUsed;
    model->info.systemCurrentMa   = model->info.systemPowerMw / 48;
    model->info.updatedAt         = time(NULL);
}

void model_init(EmuModel *model, const EmuConfig *cfg) {
    size_t i = 0;

    memset(model, 0, sizeof(*model));
    pthread_mutex_init(&model->lock, NULL);

    model->cfg       = *cfg;
    model->portCount = EMU_MAX_PORTS;
    model->running   = 1;
    model->httpFd    = -1;
    model->snmpFd    = -1;
    model->tftpFd    = -1;

    snprintf(model->activeScenario, sizeof(model->activeScenario), "%s",
             cfg->scenario);
    snprintf(model->info.manufacturer, sizeof(model->info.manufacturer),
             "Tycon Systems");
    snprintf(model->info.serial, sizeof(model->info.serial),
             "EMU-TYCON-0001");
    snprintf(model->info.hardwareVersion, sizeof(model->info.hardwareVersion),
             "HW-EMU-1.0");
    snprintf(model->info.softwareVersion, sizeof(model->info.softwareVersion),
             "SW-EMU-1.0.0");

    model->info.reachable      = 1;
    model->info.systemTempC    = 40;
    model->info.ambientTempC   = 33;
    model->info.inputVoltageMv = 54000;
    model->info.poeBudgetMw    = 240000;

    model->firmware.rebootDelaySec = 5;
    model->firmware.applyDelaySec  = 3;
    model->firmware.state          = FW_IDLE;
    model->firmware.stateSince     = time(NULL);

    for (i = 0; i < model->portCount; i++) {
        EmuPortState *port = &model->ports[i];

        port->id = (uint32_t)(i + 1);
        snprintf(port->name, sizeof(port->name), "port%u", port->id);
        snprintf(port->media, sizeof(port->media), "%s",
                 (port->id <= 8) ? "copper" : "sfp");

        port->present         = 1;
        port->adminUp         = 1;
        port->linkUp          = (port->id <= 4) ? 1 : 0;
        port->speedMbps       = port->linkUp ? 1000U : 0U;
        port->fullDuplex      = 1;
        port->poeSupported    = (port->id <= 8);
        port->poeAdminEnabled = (port->id <= 4);
        port->poeClassId      = (port->id <= 2) ? 3 : 4;
        port->updatedAt       = time(NULL);
    }

    model_recompute(model);
}

void model_set_reachable(EmuModel *model, int reachable) {
    model->info.reachable = reachable ? 1 : 0;
}

int model_set_port_link(EmuModel *model, unsigned int portId, int up) {
    if (portId < 1 || portId > model->portCount) {
        return STATUS_NOK;
    }

    model->ports[portId - 1].linkUp    = up ? 1 : 0;
    model->ports[portId - 1].speedMbps = up ? 1000U : 0U;
    model->ports[portId - 1].updatedAt = time(NULL);

    model_recompute(model);
    return STATUS_OK;
}

int model_set_port_admin(EmuModel *model, unsigned int portId, int up) {
    if (portId < 1 || portId > model->portCount) {
        return STATUS_NOK;
    }

    model->ports[portId - 1].adminUp   = up ? 1 : 0;
    model->ports[portId - 1].updatedAt = time(NULL);
    if (!up) {
        model->ports[portId - 1].linkUp = 0;
    }

    model_recompute(model);
    return STATUS_OK;
}

int model_set_port_poe(EmuModel *model, unsigned int portId, int on) {
    if (portId < 1 || portId > 8) {
        return STATUS_NOK;
    }

    model->ports[portId - 1].poeAdminEnabled = on ? 1 : 0;
    model->ports[portId - 1].updatedAt       = time(NULL);

    model_recompute(model);
    return STATUS_OK;
}

void model_stage_firmware(EmuModel *model, const char *path,
                          const char *filename) {
    snprintf(model->firmware.stagedPath, sizeof(model->firmware.stagedPath),
             "%s", (path != NULL) ? path : "");
    snprintf(model->firmware.stagedFilename,
             sizeof(model->firmware.stagedFilename), "%s",
             (filename != NULL) ? filename : "");
    model->firmware.state      = FW_STAGED;
    model->firmware.stateSince = time(NULL);
}
