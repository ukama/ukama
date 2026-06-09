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

#include "aisgd.h"
#include "ops.h"
#include "web_service.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "version.h"

static volatile bool gTerminate = false;

static UsysOption longOptions[] = {
    { "config",  required_argument, 0, 'c' },
    { "logs",    required_argument, 0, 'l' },
    { "help",    no_argument,       0, 'h' },
    { "version", no_argument,       0, 'v' },
    { 0, 0, 0, 0 }
};

static void handle_signal(int signum) {
    (void)signum;
    gTerminate = true;
}

static void set_log_level(char *slevel) {
    int ilevel = USYS_LOG_TRACE;

    if (slevel == NULL) return;
    if (!strcmp(slevel, "DEBUG")) ilevel = USYS_LOG_DEBUG;
    if (!strcmp(slevel, "INFO"))  ilevel = USYS_LOG_INFO;
    usys_log_set_level(ilevel);
}

int main(int argc, char **argv) {
    int opt;
    int optIdx;
    char *debug = DEF_LOG_LEVEL;
    char *configFile = DEF_CONFIG_FILE;
    UInst serviceInst;
    Config config;
    AppStatus status;
    AisgdContext ctx;
    JsonObj *json = NULL;

    usys_log_set_service(AISGD_SERVICE_NAME);

    while (true) {
        opt = usys_getopt_long(argc, argv, "c:l:hv", longOptions, &optIdx);
        if (opt == -1) break;
        if (opt == 'v') {
            printf("%s\n", VERSION);
            return 0;
        } else if (opt == 'c') {
            configFile = optarg;
        } else if (opt == 'l') {
            debug = optarg;
        } else {
            printf("Usage: aisgd [-c config] [-l TRACE|DEBUG|INFO] [-v]\n");
            return opt == 'h' ? 0 : 1;
        }
    }

    set_log_level(debug);
    signal(SIGINT, handle_signal);
    signal(SIGTERM, handle_signal);

    config_set_defaults(&config);
    if (!config_load_from_file(&config, configFile)) return 1;

    status_init(&status);
    memset(&ctx, 0, sizeof(ctx));
    ctx.config = &config;
    ctx.status = &status;
    controller_client_init(&ctx.controller, &config);

    if (start_web_service(&ctx, &serviceInst) != USYS_TRUE) {
        status_destroy(&status);
        config_free(&config);
        return 1;
    }

    if (aisgd_ops_reconcile(&ctx, &json)) json_decref(json);

    while (!gTerminate) pause();

    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);
    status_destroy(&status);
    config_free(&config);

    return 0;
}
