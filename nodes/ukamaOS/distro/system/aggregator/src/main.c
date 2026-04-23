/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <getopt.h>
#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "aggregator.h"
#include "network.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_file.h"
#include "usys_services.h"

#include "version.h"

typedef struct {
    Config   config;
    AppState state;
    UInst    metricsInst;
    UInst    adminInst;
    int      stateInit;
    int      stateStarted;
    int      metricsStarted;
    int      adminStarted;
} App;

static volatile sig_atomic_t gStopRequested = 0;

static struct option longOptions[] = {
    {"config",  required_argument, 0, 'c'},
    {"logs",    required_argument, 0, 'l'},
    {"help",    no_argument,       0, 'h'},
    {"version", no_argument,       0, 'v'},
    {0, 0, 0, 0}
};

static int set_log_level(char *logLevel) {

    int level = USYS_LOG_TRACE;

    if (logLevel == NULL) {
        return RETURN_NOTOK;
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
        return RETURN_NOTOK;
    }

    log_set_level(level);
    return RETURN_OK;
}

static void usage(void) {

    usys_puts("usage: aggregator [options]");
    usys_puts("");
    usys_puts("options:");
    usys_puts("  -c, --config <path>   path to aggregator config file");
    usys_puts("  -l, --logs <level>    log level: TRACE DEBUG INFO ERROR");
    usys_puts("  -h, --help            show this help");
    usys_puts("  -v, --version         show version");
}

static void app_init(App *app) {

    memset(app, 0, sizeof(*app));
}

static void app_cleanup(App *app) {

    if (app == NULL) {
        return;
    }

    if (app->adminStarted) {
        stop_web_service(&app->adminInst);
        app->adminStarted = USYS_FALSE;
    }

    if (app->metricsStarted) {
        stop_web_service(&app->metricsInst);
        app->metricsStarted = USYS_FALSE;
    }

    if (app->stateStarted) {
        app_state_stop(&app->state);
        app->stateStarted = USYS_FALSE;
    }

    if (app->stateInit) {
        app_state_cleanup(&app->state);
        app->stateInit = USYS_FALSE;
    }

    config_free(&app->config);
}

static void handle_signal(int signum) {

    (void)signum;
    gStopRequested = 1;
}

static int setup_signals(void) {

    struct sigaction action;

    memset(&action, 0, sizeof(action));
    action.sa_handler = handle_signal;
    sigemptyset(&action.sa_mask);

    if (sigaction(SIGINT, &action, NULL) != 0) {
        return RETURN_NOTOK;
    }

    if (sigaction(SIGTERM, &action, NULL) != 0) {
        return RETURN_NOTOK;
    }

    return RETURN_OK;
}

int main(int argc, char **argv) {

    int opt = 0;
    int optIndex = 0;

    int rc           = EXIT_FAILURE;
    char *configPath = DEFAULT_CONFIG_PATH;
    char *logLevel   = DEFAULT_LOG_LEVEL;
    App  app;

    app_init(&app);
    usys_log_set_service(SERVICE_NAME);

    if (set_log_level(logLevel) != RETURN_OK) {
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
            if (set_log_level(logLevel) != RETURN_OK) {
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

    if (setup_signals() != RETURN_OK) {
        return EXIT_FAILURE;
    }

    if (config_load(configPath, &app.config) != RETURN_OK) {
        goto done;
    }

    if (app_state_init(&app.state, &app.config) != RETURN_OK) {
        goto done;
    }
    app.stateInit = USYS_TRUE;

    if (app_state_start(&app.state) != RETURN_OK) {
        goto done;
    }
    app.stateStarted = USYS_TRUE;

    if (start_metrics_web_service(&app.metricsInst, &app.state) != RETURN_OK) {
        goto done;
    }
    app.metricsStarted = USYS_TRUE;

    if (start_admin_web_service(&app.adminInst, &app.state) != RETURN_OK) {
        goto done;
    }
    app.adminStarted = USYS_TRUE;

    usys_log_info("aggregator started ...");

    while (!gStopRequested) {
        sleep(1);
    }

    usys_log_info("... aggregator stopping");
    rc = EXIT_SUCCESS;

done:
    app_cleanup(&app);
    return rc;
}
