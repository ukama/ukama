/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef APP_H
#define APP_H

#include <stdbool.h>
#include <ulfius.h>

#include "config.h"
#include "jobs.h"
#include "snapshot.h"
#include "lanes.h"
#include "web_service.h"
#include "notifier.h"
#include "safety.h"
#include "i2c_bus.h"
#include "gpio_controller.h"

typedef struct {
    Config         cfg;

    Jobs           jobs;
    SnapshotStore  snap;

    I2cBus         busFem1;
    I2cBus         busFem2;
    I2cBus         busCtrl;

    GpioController gpio;

    Notifier       notifier;
    Safety         safety;

    Lanes          lanes;

    UInst          webInst;
    WebCtx         webCtx;
    ServerConfig   serverCfg;

    volatile bool  stop;
} App;

int  app_init(App *app, const char *configPath);
int  app_run(App *app);
void app_request_stop(App *app);
void app_cleanup(App *app);

#endif /* APP_H */
