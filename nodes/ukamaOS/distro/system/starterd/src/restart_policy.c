/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "restart_policy.h"

#include <stdlib.h>

int restart_policy_next_delay(Config *config, App *app, time_t now) {

    int delay;

    if (!config || !app) return 0;

    delay = app->nextBackoffSec;
    if (delay <= 0) delay = 1;

    if (delay > config->restartMaxBackoffSec) delay = config->restartMaxBackoffSec;
    return delay;
}

void restart_policy_on_start(Config *config, App *app, time_t now) {

    if (!config || !app) return;

    app->lastStartTime = now;

    if (app->lastExitTime > 0) {
        if ((now - app->lastExitTime) >= config->restartStableResetSec) {
            app->restartCount       = 0;
            app->nextBackoffSec     = 1;
            app->restartWindowStart = now;
        }
    }
}

void restart_policy_on_exit(Config *config, App *app, time_t now) {

    if (!config || !app) return;

    app->lastExitTime = now;

    if (app->restartWindowStart == 0) {
        app->restartWindowStart = now;
    }

    app->restartCount += 1;

    if (app->nextBackoffSec <= 0) app->nextBackoffSec = 1;
    app->nextBackoffSec *= 2;
    if (app->nextBackoffSec > config->restartMaxBackoffSec) {
        app->nextBackoffSec = config->restartMaxBackoffSec;
    }
}
