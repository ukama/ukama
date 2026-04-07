/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <getopt.h>
#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "collector.h"
#include "file.h"
#include "network.h"
#include "web_service.h"

#include "usys_error.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_services.h"
#include "usys_string.h"
#include "usys_api.h"

#include "version.h"

#define METRIC_CONFIG "./config/config.toml"
#define DEF_LOG_LEVEL "TRACE"

static struct option longOptions[] = {
    {"config",  required_argument, 0, 'c'},
    {"logs",    required_argument, 0, 'l'},
    {"help",    no_argument,       0, 'h'},
    {"version", no_argument,       0, 'v'},
    {0,         0,                 0,  0 }
};

static void handle_signal(int signum) {
    usys_log_info("received terminate signal");
    collector_exit(signum);
}

static int set_log_level(char *logLevel) {

    int level = USYS_LOG_TRACE;

    if (logLevel == NULL) {
        return USYS_FALSE;
    }

    if (strcmp(logLevel, "TRACE") == 0) {
        level = USYS_LOG_TRACE;
    } else if (strcmp(logLevel, "DEBUG") == 0) {
        level = USYS_LOG_DEBUG;
    } else if (strcmp(logLevel, "INFO") == 0) {
        level = USYS_LOG_INFO;
    } else if (strcmp(logLevel, "ERROR") == 0) {
        level = USYS_LOG_ERROR;
    } else {
        fprintf(stderr, "invalid log level: %s\n", logLevel);
        return USYS_FALSE;
    }

    log_set_level(level);

    return USYS_TRUE;
}

static int verify_file(char *filePath) {

    if (filePath == NULL) {
        usys_log_error("missing config file path");
        return USYS_FALSE;
    }

    if (file_exist(filePath) != USYS_TRUE) {
        usys_log_error("missing config file: %s", filePath);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static void usage(void) {

    usys_puts("usage: metrics.d [options]");
    usys_puts("");
    usys_puts("options:");
    usys_puts("  -c, --config <path>   path to metrics config file");
    usys_puts("  -l, --logs <level>    log level: TRACE DEBUG INFO ERROR");
    usys_puts("  -h, --help            show this help");
    usys_puts("  -v, --version         show version");
}

static int setup_signals(void) {

    struct sigaction action;

    memset(&action, 0, sizeof(action));

    action.sa_handler = handle_signal;
    sigemptyset(&action.sa_mask);

    if (sigaction(SIGINT, &action, NULL) != 0) {
        usys_log_error("failed to install SIGINT handler");
        return USYS_FALSE;
    }

    if (sigaction(SIGTERM, &action, NULL) != 0) {
        usys_log_error("failed to install SIGTERM handler");
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

int main(int argc, char **argv) {

    int ret              = EXIT_SUCCESS;
    int opt              = 0;
    int optIndex         = 0;
    char *configPath     = METRIC_CONFIG;
    char *logLevel       = DEF_LOG_LEVEL;
    UInst adminInstance;

    memset(&adminInstance, 0, sizeof(adminInstance));

    usys_log_set_service(SERVICE_NAME);

    if (set_log_level(logLevel) != USYS_TRUE) {
        return EXIT_FAILURE;
    }

    while (true) {
        opt = getopt_long(argc, argv, "c:l:hv", longOptions, &optIndex);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'c':
            configPath = optarg;
            break;

        case 'l':
            logLevel = optarg;
            if (set_log_level(logLevel) != USYS_TRUE) {
                usage();
                return EXIT_FAILURE;
            }
            break;

        case 'h':
            usage();
            return EXIT_SUCCESS;

        case 'v':
            usys_puts(VERSION);
            return EXIT_SUCCESS;

        default:
            usage();
            return EXIT_FAILURE;
        }
    }

    if (verify_file(configPath) != USYS_TRUE) {
        return EXIT_FAILURE;
    }

    if (setup_signals() != USYS_TRUE) {
        return EXIT_FAILURE;
    }

    usys_log_info("starting metrics collector");

    if (start_admin_web_service(&adminInstance) != USYS_TRUE) {
        usys_log_error("failed to start admin web service");
        return EXIT_FAILURE;
    }

    ret = collector(configPath);

    stop_admin_web_service(&adminInstance);

    usys_log_info("stopping metrics collector");

    return ret;
}
