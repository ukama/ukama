/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <signal.h>
#include <stdlib.h>
#include <unistd.h>

#include "config.h"
#include "engine.h"
#include "model.h"
#include "network.h"
#include "persistence.h"
#include "scenario.h"
#include "snmp_agent.h"
#include "switchemu.h"
#include "tftp_server.h"

#ifdef __has_include
#  if __has_include("usys_log.h")
#    include "usys_log.h"
#    define EMU_LOG_SET_SERVICE(name) usys_log_set_service(name)
#    define EMU_LOG_SET_LEVEL(level)  usys_log_set_level(level)
#    define EMU_LOG_DEBUG(...)        usys_log_debug(__VA_ARGS__)
#    define EMU_LOG_ERROR(...)        usys_log_error(__VA_ARGS__)
#  else
#    include <stdio.h>
#    define EMU_LOG_SET_SERVICE(name) ((void)(name))
#    define EMU_LOG_SET_LEVEL(level)  ((void)(level))
#    define EMU_LOG_DEBUG(...)        fprintf(stderr, __VA_ARGS__)
#    define EMU_LOG_ERROR(...)        fprintf(stderr, __VA_ARGS__)
#  endif
#else
#  include <stdio.h>
#  define EMU_LOG_SET_SERVICE(name) ((void)(name))
#  define EMU_LOG_SET_LEVEL(level)  ((void)(level))
#  define EMU_LOG_DEBUG(...)        fprintf(stderr, __VA_ARGS__)
#  define EMU_LOG_ERROR(...)        fprintf(stderr, __VA_ARGS__)
#endif

static volatile sig_atomic_t gTerminate = 0;

static void handle_sigterm(int signum) {
    (void)signum;
    gTerminate = 1;
}

static void setup_signal_handlers(void) {
    struct sigaction sa;

    sigemptyset(&sa.sa_mask);
    sa.sa_flags   = 0;
    sa.sa_handler = handle_sigterm;

    sigaction(SIGINT, &sa, NULL);
    sigaction(SIGTERM, &sa, NULL);
}

int main(int argc, char **argv) {
    EmuConfig cfg;
    EmuModel model;
    int exitCode = EXIT_SUCCESS;

    EMU_LOG_SET_SERVICE(SERVICE_NAME);

    if (config_load(&cfg, argc, argv) != STATUS_OK) {
        return EXIT_FAILURE;
    }

    EMU_LOG_SET_LEVEL(cfg.logLevel);
    setup_signal_handlers();

    model_init(&model, &cfg);
    pthread_mutex_lock(&model.lock);
    scenario_apply(&model, cfg.scenario);
    pthread_mutex_unlock(&model.lock);

    if (cfg.stateFile[0] != '\0') {
        persistence_load(&model, cfg.stateFile);
    }

    if (network_start(&model) != STATUS_OK ||
        snmp_agent_start(&model) != STATUS_OK ||
        tftp_server_start(&model) != STATUS_OK ||
        engine_start(&model) != STATUS_OK) {
        EMU_LOG_ERROR("Failed to start switchemu services\n");
        exitCode = EXIT_FAILURE;
        goto done;
    }

    EMU_LOG_DEBUG("Started %s\n", SERVICE_NAME);

    while (!gTerminate) {
        sleep(1);
    }

    model.running = 0;

done:
    network_stop(&model);
    snmp_agent_stop(&model);
    tftp_server_stop(&model);
    engine_stop(&model);
    persistence_save(&model, cfg.stateFile);
    pthread_mutex_destroy(&model.lock);

    return exitCode;
}
