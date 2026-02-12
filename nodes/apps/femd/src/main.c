/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include <signal.h>
#include <string.h>
#include <unistd.h>

#include "app.h"
#include "usys_log.h"

static App gApp;

static void on_signal(int sig) {
    (void)sig;
    app_request_stop(&gApp);
}

int main(int argc, char **argv) {

    const char *cfgPath = NULL;

    if (argc >= 2) {
        cfgPath = argv[1];
    }

    signal(SIGINT,  on_signal);
    signal(SIGTERM, on_signal);

    if (app_init(&gApp, cfgPath) != STATUS_OK) {
        usys_log_error("app init failed");
        app_cleanup(&gApp);
        return 1;
    }

    (void)app_run(&gApp);
    app_cleanup(&gApp);

    return 0;
}
