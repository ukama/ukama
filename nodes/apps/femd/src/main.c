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
#include "usys_getopt.h"
#include "usys_api.h"

#include "version.h"

static App gApp;

static void on_signal(int sig) {
    (void)sig;
    app_request_stop(&gApp);
}

static UsysOption longOptions[] = {
    { "config",  required_argument, 0, 'c' },
    { "logs",    required_argument, 0, 'l' },
    { "help",    no_argument,       0, 'h' },
    { "version", no_argument,       0, 'v' },
    { 0, 0, 0, 0 }
};

static void set_log_level(char *slevel) {

    int ilevel = USYS_LOG_TRACE;

    if (!strcmp(slevel, "TRACE"))      ilevel = USYS_LOG_TRACE;
    else if (!strcmp(slevel, "DEBUG")) ilevel = USYS_LOG_DEBUG;
    else if (!strcmp(slevel, "INFO"))  ilevel = USYS_LOG_INFO;

    usys_log_set_level(ilevel);
}

static void usage(void) {

    usys_puts("Usage: backhaul.d [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-c, --config                  Saftey config (json)");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-v, --version                 Software version");
}

int main(int argc, char **argv) {

    int opt, optIdx;
    const char *cfgPath = NULL;

    usys_log_set_service(SERVICE_NAME);

    signal(SIGINT,  on_signal);
    signal(SIGTERM, on_signal);

    while (1) {
        opt = 0; optIdx = 0;
        opt = usys_getopt_long(argc, argv, "vh:l:c:", longOptions, &optIdx);
        if (opt == -1) break;

        switch (opt) {
        case 'h':
            usage();
            usys_exit(0);
            break;
        case 'v':
            usys_puts(VERSION);
            usys_exit(0);
            break;
        case 'l':
            set_log_level(optarg);
            break;
        case 'c':
            cfgPath = optarg;
            break;
        default:
            usage();
            usys_exit(0);
        }
    }
    
    if (cfgPath == NULL) {
        usys_log_error("Missing arg");
        usage();
        return 1;
    }

    if (app_init(&gApp, cfgPath) != STATUS_OK) {
        usys_log_error("app init failed");
        app_cleanup(&gApp);
        return 1;
    }

    (void)app_run(&gApp);
    app_cleanup(&gApp);

    return 0;
}
