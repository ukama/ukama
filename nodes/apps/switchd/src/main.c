/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <signal.h>
#include <string.h>
#include <unistd.h>

#include "config.h"
#include "switchd.h"
#include "web_service.h"

#include "usys_api.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"

#include "version.h"

static volatile sig_atomic_t gTerminate = 0;

static UsysOption gLongOptions[] = {
    { "logs",    required_argument, 0, 'l' },
    { "help",    no_argument,       0, 'h' },
    { "version", no_argument,       0, 'v' },
    { 0, 0, 0, 0 }
};

static void handle_signal(int signum) {
    (void)signum;
    gTerminate = 1;
}

static void set_log_level(const char *level) {
    int logLevel = USYS_LOG_INFO;

    if (level == NULL) {
        return;
    }

    if (strcmp(level, "TRACE") == 0) {
        logLevel = USYS_LOG_TRACE;
    } else if (strcmp(level, "DEBUG") == 0) {
        logLevel = USYS_LOG_DEBUG;
    } else if (strcmp(level, "INFO") == 0) {
        logLevel = USYS_LOG_INFO;
    } else if (strcmp(level, "WARNING") == 0) {
        logLevel = USYS_LOG_WARN;
    } else if (strcmp(level, "ERROR") == 0) {
        logLevel = USYS_LOG_ERROR;
    }

    usys_log_set_level(logLevel);
}

int main(int argc, char **argv) {
    int opt;
    int optIdx;
    int exitCode;
    char *logLevel;
    UInst serviceInst;

    exitCode = 0;
    logLevel = DEF_LOG_LEVEL;
    memset(&serviceInst, 0, sizeof(serviceInst));

    usys_log_set_service(SERVICE_NAME);
    set_log_level(logLevel);

    while (true) {
        opt = 0;
        optIdx = 0;
        opt = usys_getopt_long(argc,
                               argv,
                               "hvl:",
                               gLongOptions,
                               &optIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            config_usage();
            usys_exit(0);
            break;

        case 'v':
            usys_puts(VERSION);
            usys_exit(0);
            break;

        case 'l':
            logLevel = optarg;
            set_log_level(logLevel);
            break;

        default:
            config_usage();
            usys_exit(1);
        }
    }

    signal(SIGINT, handle_signal);
    signal(SIGTERM, handle_signal);

    if (switchd_init(&gSwitchd) != SWITCHD_OK) {
        usys_log_error("Failed to initialize %s", SERVICE_NAME);
        return 1;
    }

    if (switchd_start(&gSwitchd) != SWITCHD_OK) {
        usys_log_error("Failed to start %s workers", SERVICE_NAME);
        switchd_cleanup(&gSwitchd);
        return 1;
    }

    if (web_service_start(&gSwitchd, &serviceInst) != STATUS_OK) {
        usys_log_error("Failed to start %s web service", SERVICE_NAME);
        switchd_stop(&gSwitchd);
        switchd_cleanup(&gSwitchd);
        return 1;
    }

    usys_log_info("%s started", SERVICE_NAME);

    while (!gTerminate) {
        sleep(1);
    }

    usys_log_info("%s stopping", SERVICE_NAME);

    switchd_request_terminate(&gSwitchd);
    web_service_stop(&serviceInst);
    switchd_stop(&gSwitchd);
    switchd_cleanup(&gSwitchd);

    return exitCode;
}
