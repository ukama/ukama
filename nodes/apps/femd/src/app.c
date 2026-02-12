/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <unistd.h>

#include "app.h"

#include "usys_log.h"

int app_init(App *app, const char *configPath) {

    if (!app) return STATUS_NOK;

    memset(app, 0, sizeof(*app));

    (void)config_set_defaults(&app->cfg);
    if (configPath) (void)config_load(&app->cfg, configPath);

    if (config_validate(&app->cfg) != STATUS_OK) {
        return STATUS_NOK;
    }

    config_print(&app->cfg);

    if (snapshot_init(&app->snap) != STATUS_OK)                               return STATUS_NOK;
    if (jobs_init(&app->jobs) != STATUS_OK)                                   return STATUS_NOK;
    if (i2c_bus_init(&app->busFem1, app->cfg.i2cBusFem1) != STATUS_OK)        return STATUS_NOK;
    if (i2c_bus_init(&app->busFem2, app->cfg.i2cBusFem2) != STATUS_OK)        return STATUS_NOK;
    if (i2c_bus_init(&app->busCtrl, app->cfg.i2cBusCtrl) != STATUS_OK)        return STATUS_NOK;
    if (gpio_controller_init(&app->gpio, app->cfg.gpioBasePath) != STATUS_OK) return STATUS_NOK;
    if (notifier_init(&app->notifier, &app->cfg) != STATUS_OK)                return STATUS_NOK;

    if (safety_init(&app->safety, &app->jobs, &app->snap, app->cfg.safetyConfigPath) != STATUS_OK) {
        return STATUS_NOK;
    }

    if (lanes_init(&app->lanes,
                   &app->jobs,
                   &app->snap,
                   &app->safety,
                   &app->busFem1,
                   &app->busFem2,
                   &app->busCtrl,
                   &app->gpio,
                   app->cfg.samplePeriodMs,
                   app->cfg.safetyPeriodMs) != STATUS_OK) {
        return STATUS_NOK;
    }

    app->webCtx.config = &app->cfg;
    app->webCtx.jobs   = &app->jobs;
    app->webCtx.snap   = &app->snap;

    app->serverCfg.config = &app->cfg;

    if (app->cfg.enableWeb) {
        if (start_web_service(&app->serverCfg, &app->webInst, &app->webCtx) != USYS_TRUE) {
            return STATUS_NOK;
        }
    }

    return STATUS_OK;
}

int app_run(App *app) {

    if (!app) return STATUS_NOK;

    if (lanes_start(&app->lanes) != STATUS_OK) {
        return STATUS_NOK;
    }

    while (!app->stop) {
        usleep(100000);
    }

    return STATUS_OK;
}

void app_request_stop(App *app) {
    if (app) app->stop = true;
}

void app_cleanup(App *app) {

    if (!app) return;

    lanes_stop(&app->lanes);

    if (app->cfg.enableWeb) {
        ulfius_stop_framework(&app->webInst);
        ulfius_clean_instance(&app->webInst);
    }

    gpio_controller_cleanup(&app->gpio);

    i2c_bus_cleanup(&app->busCtrl);
    i2c_bus_cleanup(&app->busFem2);
    i2c_bus_cleanup(&app->busFem1);

    jobs_cleanup(&app->jobs);
    snapshot_cleanup(&app->snap);

    memset(app, 0, sizeof(*app));
}
