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
#include <pthread.h>

#include "config.h"
#include "model.h"
#include "server.h"
#include "web_service.h"
#include "usys_log.h"
#include "version.h"

static volatile sig_atomic_t gRunning = 1;

typedef struct {
    EmuConfig *config;
    EmuModel *model;
} ServerThreadArg;

static void handle_signal(int sig)
{
    (void)sig;
    gRunning = 0;
}

static void *server_thread(void *arg)
{
    ServerThreadArg *threadArg;

    threadArg = arg;
    emu_server_run(threadArg->config, threadArg->model, &gRunning);

    return NULL;
}

static bool install_signal_handlers(void)
{
    struct sigaction sa;

    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = handle_signal;
    sigemptyset(&sa.sa_mask);

    if (sigaction(SIGINT, &sa, NULL) < 0) {
        return false;
    }

    if (sigaction(SIGTERM, &sa, NULL) < 0) {
        return false;
    }

    return true;
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

static void usage(const char *prog)
{
    fprintf(stderr, "Usage: %s [-c config] [-l log-level]\n", prog);
}

int main(int argc, char **argv)
{
    EmuConfig config;
    EmuModel model;
    UInst web;
    ServerThreadArg threadArg;
    pthread_t tid;
    const char *configFile;
    int opt;

    configFile = AISG_EMU_CONFIG_FILE;

    if (!emu_config_init(&config)) {
        return 1;
    }

    emu_model_init(&model);
    usys_log_set_level(LOG_INFO);
    usys_log_set_service(AISG_EMU_SERVICE_NAME);

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
            emu_model_free(&model);
            emu_config_free(&config);
            return 0;
        case 'h':
        default:
            usage(argv[0]);
            emu_model_free(&model);
            emu_config_free(&config);
            return 1;
        }
    }

    if (!emu_config_load(&config, configFile)) {
        usys_log_warn("using default aisg-emu configuration");
    }

    emu_model_load_scenario(&model, config.scenario);

    if (!install_signal_handlers()) {
        usys_log_error("failed to install signal handlers");
        emu_model_free(&model);
        emu_config_free(&config);
        return 1;
    }

    threadArg.config = &config;
    threadArg.model  = &model;

    if (pthread_create(&tid, NULL, server_thread, &threadArg) != 0) {
        usys_log_error("failed to create emulator socket thread");
        emu_model_free(&model);
        emu_config_free(&config);
        return 1;
    }

    start_web_service(&web, &config, &model);

    while (gRunning) {
        pause();
    }

    usys_log_info("stopping aisg-emu");

    stop_web_service(&web);
    pthread_join(tid, NULL);

    emu_model_free(&model);
    emu_config_free(&config);

    return 0;
}
