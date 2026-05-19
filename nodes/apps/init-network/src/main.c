/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "config.h"
#include "init_network.h"
#include "ovs.h"
#include "status.h"

#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

#include "version.h"

static volatile bool gTerminate = false;

static UsysOption longOptions[] = {
    { "config",  required_argument, 0, 'c' },
    { "logs",    required_argument, 0, 'l' },
    { "help",    no_argument,       0, 'h' },
    { "version", no_argument,       0, 'v' },
    { 0, 0, 0, 0 }
};

static void handle_sigint(int signum) {

    (void)signum;

    usys_log_debug("Terminate signal");
    gTerminate = true;
}

static void set_log_level(char *slevel) {

    int ilevel = USYS_LOG_TRACE;

    if (slevel == NULL) return;

    if (!strcmp(slevel, "TRACE")) {
        ilevel = USYS_LOG_TRACE;
    } else if (!strcmp(slevel, "DEBUG")) {
        ilevel = USYS_LOG_DEBUG;
    } else if (!strcmp(slevel, "INFO")) {
        ilevel = USYS_LOG_INFO;
    }

    usys_log_set_level(ilevel);
}

static void usage(void) {

    printf("Usage: init-network.d [options]\n");
    printf("Options:\n");
    printf("-h, --help                    Help menu\n");
    printf("-c, --config <file>           Config file\n");
    printf("-l, --logs <TRACE|DEBUG|INFO> Log level for the process\n");
    printf("-v, --version                 Software version\n");
}

int main(int argc, char **argv) {

    int opt;
    int optIdx;
    char *debug;
    char *configFile;
    UInst serviceInst;
    Config config;
    AppStatus status;
    ServiceContext ctx;
    bool ready;

    debug = DEF_LOG_LEVEL;
    configFile = DEF_CONFIG_FILE;
    ready = false;

    usys_log_set_service(INIT_NETWORK_SERVICE_NAME);

    while (true) {
        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "c:l:hv", longOptions, &optIdx);
        if (opt == -1) break;

        switch (opt) {
        case 'h':
            usage();
            exit(0);
            break;

        case 'v':
            printf("%s\n", VERSION);
            exit(0);
            break;

        case 'c':
            configFile = optarg;
            break;

        case 'l':
            debug = optarg;
            break;

        default:
            usage();
            exit(1);
        }
    }

    set_log_level(debug);

    config_set_defaults(&config);
    if (!config_load_from_file(&config, configFile)) {
        config_free(&config);
        exit(1);
    }

    status_init(&status);

    memset(&ctx, 0, sizeof(ServiceContext));
    ctx.config = &config;
    ctx.status = &status;

    signal(SIGINT, handle_sigint);
    signal(SIGTERM, handle_sigint);

    usys_log_debug("Starting %s", INIT_NETWORK_SERVICE_NAME);

    if (start_web_service(&ctx, &serviceInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to start");
        status_destroy(&status);
        config_free(&config);
        exit(1);
    }

    if (!ovs_setup(&config, &status)) {
        usys_log_error("Network initialization failed");
    }

    while (!gTerminate) {
        pause();
    }

    ready = status_is_ready(&status);

    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);
    status_destroy(&status);
    config_free(&config);

    usys_log_debug("Exiting %s", INIT_NETWORK_SERVICE_NAME);

    return ready ? 0 : 1;
}
