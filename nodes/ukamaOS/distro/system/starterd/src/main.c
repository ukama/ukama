/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>

#include "starterd.h"
#include "config.h"
#include "manifest.h"
#include "state_store.h"
#include "actions.h"
#include "supervisor.h"
#include "network.h"
#include "web_service.h"

#include "usys_log.h"

static volatile sig_atomic_t gTerminate = 0;

static void on_signal(int sig) {

    (void)sig;
    gTerminate = 1;
}

static void setup_signals(void) {

    struct sigaction sa;

    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = on_signal;
    sigaction(SIGTERM, &sa, NULL);
    sigaction(SIGINT, &sa, NULL);
}

static void redirect_logs(const char *path) {

    int fd;

    if (!path || !*path) {
        return;
    }

    fd = open(path, O_CREAT | O_APPEND | O_WRONLY, 0644);
    if (fd < 0) {
        return;
    }

    dup2(fd, STDOUT_FILENO);
    dup2(fd, STDERR_FILENO);

    if (fd > 2) {
        close(fd);
    }
}

static int log_level_from_env(void) {

    const char *v;

    v = getenv("STARTERD_LOG_LEVEL");
    if (!v || !*v) {
        return USYS_LOG_DEBUG;
    }

    if (strcmp(v, "debug") == 0) return USYS_LOG_DEBUG;
    if (strcmp(v, "info") == 0)  return USYS_LOG_INFO;
    if (strcmp(v, "warn") == 0)  return USYS_LOG_WARN;
    if (strcmp(v, "error") == 0) return USYS_LOG_ERROR;

    return USYS_LOG_DEBUG;
}

int main(int argc, char **argv) {

    Config config;
    Space *spaceList;
    ActionQueue queue;
    StarterContext ctx;
    Supervisor *sup;
    Action *a;
    int exitCode;

    (void)argc;
    (void)argv;

    setup_signals();

    usys_log_set_service(STARTERD_SERVICE_NAME);
    usys_log_set_level(log_level_from_env());

    if (!config_load(&config)) {
        usys_log_error("startup: config load failed");
        return 1;
    }

    redirect_logs(config.logPath);

    spaceList = NULL;
    if (!manifest_load(&config, &spaceList)) {
        usys_log_error("startup: manifest load failed");
        config_free(&config);
        return 1;
    }

    state_store_load(&config, spaceList);

    actions_init(&queue);

    memset(&ctx, 0, sizeof(ctx));
    ctx.config             = &config;
    ctx.spaceList          = spaceList;
    ctx.queue              = &queue;
    ctx.supervisor         = NULL;
    ctx.uInstance          = NULL;
    ctx.terminateRequested = 0;
    ctx.switchRequested    = 0;
    ctx.updateInProgress   = 0;
    ctx.exitCode           = 0;

    if (!network_init(&ctx)) {
        usys_log_error("startup: network init failed");
        manifest_free(spaceList);
        config_free(&config);
        return 1;
    }

    sup = supervisor_start(&config, spaceList, &queue, &ctx);
    if (!sup) {
        usys_log_error("startup: supervisor start failed");
        network_shutdown(&ctx);
        manifest_free(spaceList);
        config_free(&config);
        return 1;
    }

    ctx.supervisor = sup;

    if (!web_service_start(&ctx)) {
        usys_log_error("startup: web service start failed");
        supervisor_stop(sup);
        network_shutdown(&ctx);
        manifest_free(spaceList);
        config_free(&config);
        return 1;
    }

    a = action_new(ACTION_RUN_BOOT, NULL, NULL, NULL, NULL);
    actions_enqueue(&queue, a);

    a = action_new(ACTION_RUN_ALL, NULL, NULL, NULL, NULL);
    actions_enqueue(&queue, a);

    supervisor_signal(sup);

    usys_log_info("starterd: running on %s:%d", config.httpAddr, config.httpPort);

    while (!gTerminate && !ctx.switchRequested) {
        sleep(1);
    }

    if (gTerminate) {
        ctx.terminateRequested = 1;
        usys_log_info("starterd: terminating by signal");
    } else if (ctx.switchRequested) {
        usys_log_info("starterd: self-update switch requested");
    }

    web_service_stop(&ctx);
    supervisor_stop(sup);
    network_shutdown(&ctx);
    actions_free(&queue);
    state_store_save(&config, spaceList);
    manifest_free(spaceList);
    config_free(&config);

    exitCode = ctx.switchRequested ? 77 : ctx.exitCode;

    if (exitCode == 77) {
        usys_log_info("starterd: exiting with switch code 77");
    } else {
        usys_log_info("starterd: exiting with code %d", exitCode);
    }

    return exitCode;
}
