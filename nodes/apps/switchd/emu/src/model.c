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

static void model_update_port_counters(EmuPortState *port) {

    uint64_t base;

    if (!port || !port->present || !port->linkUp) {
        return;
    }

    base = (uint64_t)port->id;

    port->rxBytes += 4096U + (base * 257U);
    port->txBytes += 3072U + (base * 193U);
    port->rxPackets += 32U + base;
    port->txPackets += 24U + base;

    if (port->faultLink) {
        port->rxErrors += 1U;
        port->txErrors += 1U;
    }

    if (port->faultPoe) {
        port->txErrors += 1U;
    }
}

void model_recompute(EmuModel *model) {

    size_t i;
    int poeUsed;
    int poeAlarm;
    int linkAlarm;

    i = 0;
    poeUsed = 0;
    poeAlarm = 0;
    linkAlarm = 0;

    for (i = 0; i < model->portCount; i++) {
        EmuPortState *port;

        port = &model->ports[i];

        if (!port->adminUp) {
            port->linkUp = 0;
            port->speedMbps = 0;
        } else if (!port->faultLink &&
                   (port->id <= 3 || port->id == 9)) {
            port->linkUp = 1;
            port->speedMbps = 1000U;
        }

        if (port->poeSupported && port->poeAdminEnabled &&
            !port->faultPoe) {
            port->poeOperStatus = 1;

            switch (port->id) {
            case 1:
                port->poePowerMw = 11500;
                port->poeCurrentMa = 221;
                break;

            case 2:
                port->poePowerMw = 6500;
                port->poeCurrentMa = 125;
                break;

            case 3:
                port->poePowerMw = 8000;
                port->poeCurrentMa = 154;
                break;

            default:
                port->poePowerMw = 0;
                port->poeCurrentMa = 0;
                break;
            }

            port->poeVoltageMv = 52000;
            poeUsed += port->poePowerMw;
        } else {
            port->poeOperStatus = port->faultPoe ? 2 : 0;
            port->poePowerMw = 0;
            port->poeCurrentMa = 0;
            port->poeVoltageMv = port->poeSupported ? 52000 : 0;
        }

        model_update_port_counters(port);

        if (port->faultPoe) {
            poeAlarm = 1;
        }

        if (port->faultLink ||
            (port->adminUp && !port->linkUp && port->id <= 8)) {
            linkAlarm = 1;
        }

        port->updatedAt = time(NULL);
    }

    model->info.poeUsedMw = poeUsed;
    model->info.alarmPoeFailure =
        poeAlarm || (poeUsed > model->info.poeBudgetMw);
    model->info.alarmLinkFailure = linkAlarm;
    model->info.systemPowerMw = 25000 + poeUsed;
    model->info.systemCurrentMa = model->info.systemPowerMw / 48;
    model->info.updatedAt = time(NULL);
}

static void model_init_default_port(EmuPortState *port, uint32_t id) {

    memset(port, 0, sizeof(*port));

    port->id = id;
    snprintf(port->name, sizeof(port->name), "port%u", id);
    snprintf(port->media,
             sizeof(port->media),
             "%s",
             (id <= 8) ? "copper" : "sfp");

    port->present = 1;
    port->adminUp = 1;
    port->linkUp = 0;
    port->speedMbps = 0;
    port->fullDuplex = 1;
    port->poeSupported = (id <= 8);
    port->poeAdminEnabled = 0;
    port->poeClassId = 0;
    port->updatedAt = time(NULL);
}

static void model_apply_ukama_profile(EmuModel *model) {

    if (model->portCount >= 1) {
        snprintf(model->ports[0].name,
                 sizeof(model->ports[0].name),
                 "tnode-poe");
        model->ports[0].linkUp = 1;
        model->ports[0].speedMbps = 1000U;
        model->ports[0].poeSupported = 1;
        model->ports[0].poeAdminEnabled = 1;
        model->ports[0].poeClassId = 4;
    }

    if (model->portCount >= 2) {
        snprintf(model->ports[1].name,
                 sizeof(model->ports[1].name),
                 "cnode-poe");
        model->ports[1].linkUp = 1;
        model->ports[1].speedMbps = 1000U;
        model->ports[1].poeSupported = 1;
        model->ports[1].poeAdminEnabled = 1;
        model->ports[1].poeClassId = 3;
    }

    if (model->portCount >= 3) {
        snprintf(model->ports[2].name,
                 sizeof(model->ports[2].name),
                 "anode-poe");
        model->ports[2].linkUp = 1;
        model->ports[2].speedMbps = 1000U;
        model->ports[2].poeSupported = 1;
        model->ports[2].poeAdminEnabled = 1;
        model->ports[2].poeClassId = 4;
    }

    if (model->portCount >= 4) {
        snprintf(model->ports[3].name,
                 sizeof(model->ports[3].name),
                 "spare-poe");
        model->ports[3].linkUp = 0;
        model->ports[3].speedMbps = 0;
        model->ports[3].poeSupported = 1;
        model->ports[3].poeAdminEnabled = 0;
        model->ports[3].poeClassId = 0;
    }

    if (model->portCount >= 9) {
        snprintf(model->ports[8].name,
                 sizeof(model->ports[8].name),
                 "uplink-sfp");
        snprintf(model->ports[8].media,
                 sizeof(model->ports[8].media),
                 "sfp");
        model->ports[8].linkUp = 1;
        model->ports[8].speedMbps = 1000U;
        model->ports[8].poeSupported = 0;
        model->ports[8].poeAdminEnabled = 0;
        model->ports[8].poeClassId = 0;
    }
}

void model_init(EmuModel *model, const EmuConfig *cfg) {

    size_t i;

    memset(model, 0, sizeof(*model));
    pthread_mutex_init(&model->lock, NULL);

    model->cfg = *cfg;
    model->portCount = EMU_MAX_PORTS;
    model->running = 1;
    model->httpFd = -1;
    model->snmpFd = -1;
    model->tftpFd = -1;

    snprintf(model->activeScenario,
             sizeof(model->activeScenario),
             "%s",
             cfg->scenario);
    snprintf(model->info.manufacturer,
             sizeof(model->info.manufacturer),
             "Tycon Systems");
    snprintf(model->info.serial,
             sizeof(model->info.serial),
             "EMU-TYCON-0001");
    snprintf(model->info.hardwareVersion,
             sizeof(model->info.hardwareVersion),
             "HW-EMU-1.0");
    snprintf(model->info.softwareVersion,
             sizeof(model->info.softwareVersion),
             "SW-EMU-1.0.0");

    model->info.reachable = 1;
    model->info.systemTempC = 40;
    model->info.ambientTempC = 33;
    model->info.inputVoltageMv = 54000;
    model->info.poeBudgetMw = 240000;

    model->firmware.rebootDelaySec = 5;
    model->firmware.applyDelaySec = 3;
    model->firmware.state = FW_IDLE;
    model->firmware.stateSince = time(NULL);

    for (i = 0; i < model->portCount; i++) {
        model_init_default_port(&model->ports[i], (uint32_t)(i + 1));
    }

    model_apply_ukama_profile(model);
    model_recompute(model);
}

void model_set_reachable(EmuModel *model, int reachable) {

    model->info.reachable = reachable ? 1 : 0;
}

int model_set_port_link(EmuModel *model, unsigned int portId, int up) {

    if (portId < 1 || portId > model->portCount) {
        return STATUS_NOK;
    }

    model->ports[portId - 1].faultLink = up ? 0 : 1;
    model->ports[portId - 1].linkUp = up ? 1 : 0;
    model->ports[portId - 1].speedMbps = up ? 1000U : 0U;
    model->ports[portId - 1].updatedAt = time(NULL);

    model_recompute(model);

    return STATUS_OK;
}

int model_set_port_admin(EmuModel *model, unsigned int portId, int up) {

    if (portId < 1 || portId > model->portCount) {
        return STATUS_NOK;
    }

    model->ports[portId - 1].adminUp = up ? 1 : 0;
    model->ports[portId - 1].updatedAt = time(NULL);

    if (!up) {
        model->ports[portId - 1].linkUp = 0;
        model->ports[portId - 1].speedMbps = 0;
    }

    model_recompute(model);

    return STATUS_OK;
}

int model_set_port_poe(EmuModel *model, unsigned int portId, int on) {

    if (portId < 1 || portId > 8) {
        return STATUS_NOK;
    }

    model->ports[portId - 1].poeAdminEnabled = on ? 1 : 0;
    model->ports[portId - 1].updatedAt = time(NULL);

    model_recompute(model);

    return STATUS_OK;
}

void model_stage_firmware(EmuModel *model,
                          const char *path,
                          const char *filename) {

    snprintf(model->firmware.stagedPath,
             sizeof(model->firmware.stagedPath),
             "%s",
             (path != NULL) ? path : "");
    snprintf(model->firmware.stagedFilename,
             sizeof(model->firmware.stagedFilename),
             "%s",
             (filename != NULL) ? filename : "");

    model->firmware.state = FW_STAGED;
    model->firmware.stateSince = time(NULL);
}
