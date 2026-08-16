/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "config.h"
#include "lifecycled.h"
#include "web_service.h"

#include "usys_log.h"
#include "usys_services.h"

static volatile sig_atomic_t gTerminate = 0;

static void on_signal(int signalNumber) {

    (void)signalNumber;
    gTerminate = 1;
}

static void setup_signals(void) {

    struct sigaction action;

    memset(&action, 0, sizeof(action));
    action.sa_handler = on_signal;
    sigaction(SIGINT, &action, NULL);
    sigaction(SIGTERM, &action, NULL);
    signal(SIGPIPE, SIG_IGN);
}

static int log_level(void) {

    const char *value;

    value = getenv("LIFECYCLED_LOG_LEVEL");
    if (!value || !*value) return USYS_LOG_DEBUG;

    if (strcmp(value, "debug") == 0) return USYS_LOG_DEBUG;
    if (strcmp(value, "info") == 0)  return USYS_LOG_INFO;
    if (strcmp(value, "warn") == 0)  return USYS_LOG_WARN;
    if (strcmp(value, "error") == 0) return USYS_LOG_ERROR;

    return USYS_LOG_DEBUG;
}

int main(int argc, char **argv) {

    LifecycleContext ctx;
    Config config;
    int result;

    (void)argc;
    (void)argv;

    result = 1;
    setup_signals();
    usys_log_set_service(SERVICE_LIFECYCLE);
    usys_log_set_level(log_level());

    if (!config_load(&config)) {
        usys_log_error("startup: configuration failed");
        return 1;
    }

    if (curl_global_init(CURL_GLOBAL_DEFAULT) != CURLE_OK) {
        usys_log_error("startup: curl initialization failed");
        config_free(&config);
        return 1;
    }

    if (!lifecycle_context_init(&ctx, &config)) {
        usys_log_error("startup: lifecycle context failed");
        goto cleanup_curl;
    }

    if (!lifecycle_context_start(&ctx)) {
        usys_log_error("startup: lifecycle worker failed");
        goto cleanup_context;
    }

    if (!web_service_start(&ctx)) {
        usys_log_error("startup: web service failed");
        lifecycle_context_stop(&ctx);
        goto cleanup_context;
    }

    usys_log_info("lifecycle.d running");
    while (!gTerminate) sleep(1);

    web_service_stop(&ctx);
    lifecycle_context_stop(&ctx);
    result = 0;

cleanup_context:
    lifecycle_context_free(&ctx);

cleanup_curl:
    curl_global_cleanup();
    config_free(&config);
    return result;
}
