/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <time.h>
#include <unistd.h>

#include "alarms.h"
#include "engine.h"
#include "model.h"

static void *engine_main(void *arg) {
    EmuModel *model = (EmuModel *)arg;
    int flapTick = 0;

    while (model->running) {
        pthread_mutex_lock(&model->lock);

        if (model->faults.flapPortId > 0 &&
            model->faults.flapPortId <= (int)model->portCount &&
            model->faults.flapPeriodSec > 0) {
            flapTick++;
            if (flapTick >= model->faults.flapPeriodSec) {
                EmuPortState *port = &model->ports[model->faults.flapPortId - 1];
                port->linkUp       = !port->linkUp;
                port->speedMbps    = port->linkUp ? 1000U : 0U;
                flapTick           = 0;
            }
        }

        if (model->info.reachable) {
            size_t i = 0;
            for (i = 0; i < model->portCount; i++) {
                if (model->ports[i].linkUp) {
                    model->ports[i].rxBytes += 2048;
                    model->ports[i].txBytes += 1536;
                    model->ports[i].rxPackets += 8;
                    model->ports[i].txPackets += 6;
                }
            }
        }

        if (model->firmware.state == FW_APPLYING &&
            (time(NULL) - model->firmware.stateSince) >=
            model->firmware.applyDelaySec) {
            if (model->firmware.applyShouldFail) {
                model->firmware.state = FW_FAILED;
                model->firmware.executeStatus = -1;
            } else {
                model->firmware.state = FW_REBOOTING;
                model->firmware.stateSince = time(NULL);
                model->info.rebooting = 1;
                model->info.reachable = 0;
            }
        } else if (model->firmware.state == FW_REBOOTING &&
                   (time(NULL) - model->firmware.stateSince) >=
                   model->firmware.rebootDelaySec) {
            model->firmware.state = FW_DONE;
            model->firmware.executeStatus = 1;
            model->info.rebooting = 0;
            model->info.reachable = 1;
        }

        model_recompute(model);
        alarms_refresh(model);
        pthread_mutex_unlock(&model->lock);
        sleep(1);
    }

    return NULL;
}

int engine_start(EmuModel *model) {
    return pthread_create(&model->engineThread, NULL, engine_main, model);
}

void engine_stop(EmuModel *model) {
    if (model->engineThread != 0U) {
        pthread_join(model->engineThread, NULL);
        model->engineThread = 0;
    }
}
