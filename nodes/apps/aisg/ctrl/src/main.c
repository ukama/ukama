/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <signal.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "backend.h"
#include "backend_raw_rs485.h"
#include "backend_stm_uart.h"
#include "config.h"
#include "ctrl.h"
#include "server.h"
#include "usys_log.h"
#include "version.h"

static volatile bool gRunning = true;

static void handle_signal(int sig) {
    (void)sig;
    gRunning = false;
}


static void set_log_level(const char *level)
{
    int logLevel = USYS_LOG_INFO;

    if (level == NULL) {
        return;
    }

    if (!strcmp(level, "TRACE")) {
        logLevel = USYS_LOG_TRACE;
    } else if (!strcmp(level, "DEBUG")) {
        logLevel = USYS_LOG_DEBUG;
    } else if (!strcmp(level, "INFO")) {
        logLevel = USYS_LOG_INFO;
    }

    usys_log_set_level(logLevel);
}

static void usage(const char *prog) {
    fprintf(stderr, "Usage: %s [-c config] [-l log-level]\n", prog);
}

int main(int argc, char **argv) {
    Config config;
    Backend backend;
    const char *configFile;
    int opt;
    bool ok;

    configFile = DEF_CONFIG_FILE;
    config_set_defaults(&config);
    usys_log_set_level(LOG_INFO);

    while ((opt = getopt(argc, argv, "c:l:hv")) != -1) {
        switch (opt) {
        case 'c':
            configFile = optarg;
            break;
        case 'l':
            set_log_level(optarg);
            break;
        case 'v':
            printf("%s\n", VERSION);
            return 0;
        case 'h':
        default:
            usage(argv[0]);
            return 1;
        }
    }

    if (!config_load_from_file(&config, configFile)) {
        usys_log_warn("using default aisg-ctrl configuration");
    }

    signal(SIGINT, handle_signal);
    signal(SIGTERM, handle_signal);

    if (!backend_init(&backend, &config)) {
        usys_log_error("failed to initialize backend");
        config_free(&config);
        return 1;
    }

    if (!backend_open(&backend)) {
        usys_log_error("failed to open backend");
        backend_close(&backend);
        config_free(&config);
        return 1;
    }

    ok = ctrl_server_run(&config, &backend, &gRunning);
    backend_close(&backend);
    config_free(&config);

    return ok ? 0 : 1;
}
