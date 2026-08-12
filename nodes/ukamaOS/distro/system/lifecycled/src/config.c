/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "config.h"
#include "lifecycled.h"

#include "usys_log.h"
#include "usys_services.h"

#define DEFAULT_HTTP_PORT                    8097
#define DEFAULT_CHECK_IN_TIMEOUT_SEC           60
#define DEFAULT_CONFIG_TIMEOUT_SEC             60
#define DEFAULT_STARTER_UNAVAILABLE_SEC        15
#define DEFAULT_POLL_INTERVAL_MS             1000
#define DEFAULT_REQUEST_TIMEOUT_SEC             3
#define DEFAULT_STATE_FILE \
    "/ukama/state/lifecycled/state"

static char *cfg_string(const char *name, const char *fallback) {

    const char *value;

    value = getenv(name);
    if (value && *value) return strdup(value);
    return fallback ? strdup(fallback) : NULL;
}

static int cfg_integer(const char *name, int fallback) {

    const char *value;
    char *end;
    long parsed;

    value = getenv(name);
    if (!value || !*value) return fallback;

    errno = 0;
    end = NULL;
    parsed = strtol(value, &end, 10);

    if (errno != 0 || end == value || *end != '\0' ||
        parsed < INT_MIN || parsed > INT_MAX) {
        return fallback;
    }

    return (int)parsed;
}

static int cfg_service_port(const char *envName,
                            const char *service,
                            int fallback) {

    int port;

    port = cfg_integer(envName, 0);
    if (port > 0) return port;

    port = usys_find_service_port(service);
    return port > 0 ? port : fallback;
}

static bool cfg_valid(const Config *config) {

    if (!config) return false;

    if (config->httpPort <= 0 || config->httpPort > 65535 ||
        config->starterPort <= 0 || config->starterPort > 65535 ||
        config->notifyPort <= 0 || config->notifyPort > 65535) {
        usys_log_error("config: invalid service port");
        return false;
    }

    if (!config->httpAddr || !config->starterHost ||
        !config->notifyHost || !config->stateFile) {
        usys_log_error("config: missing required string");
        return false;
    }

    if (config->checkInTimeoutSec <= 0 ||
        config->configTimeoutSec <= 0 ||
        config->starterUnavailableTimeoutSec <= 0 ||
        config->pollIntervalMs <= 0 ||
        config->requestTimeoutSec <= 0) {
        usys_log_error("config: invalid timeout");
        return false;
    }

    return true;
}

bool config_load(Config *config) {

    if (!config) return false;

    memset(config, 0, sizeof(*config));

    config->httpAddr = cfg_string("LIFECYCLED_HTTP_ADDR", "127.0.0.1");
    config->httpPort = cfg_service_port("LIFECYCLED_HTTP_PORT",
                                        LIFECYCLED_SERVICE_NAME,
                                        DEFAULT_HTTP_PORT);

    config->starterHost =
        cfg_string("LIFECYCLED_STARTER_HOST", "127.0.0.1");
    config->starterPort = cfg_service_port("LIFECYCLED_STARTER_PORT",
                                           SERVICE_STARTER,
                                           0);

    config->notifyHost =
        cfg_string("LIFECYCLED_NOTIFY_HOST", "127.0.0.1");
    config->notifyPort = cfg_service_port("LIFECYCLED_NOTIFY_PORT",
                                          SERVICE_NOTIFY,
                                          0);

    config->stateFile =
        cfg_string("LIFECYCLED_STATE_FILE", DEFAULT_STATE_FILE);

    config->checkInTimeoutSec =
        cfg_integer("LIFECYCLED_CHECKIN_TIMEOUT_SEC",
                    DEFAULT_CHECK_IN_TIMEOUT_SEC);
    config->configTimeoutSec =
        cfg_integer("LIFECYCLED_CONFIG_TIMEOUT_SEC",
                    DEFAULT_CONFIG_TIMEOUT_SEC);
    config->starterUnavailableTimeoutSec =
        cfg_integer("LIFECYCLED_STARTER_UNAVAILABLE_TIMEOUT_SEC",
                    DEFAULT_STARTER_UNAVAILABLE_SEC);
    config->pollIntervalMs =
        cfg_integer("LIFECYCLED_POLL_INTERVAL_MS",
                    DEFAULT_POLL_INTERVAL_MS);
    config->requestTimeoutSec =
        cfg_integer("LIFECYCLED_REQUEST_TIMEOUT_SEC",
                    DEFAULT_REQUEST_TIMEOUT_SEC);

    if (!cfg_valid(config)) {
        config_free(config);
        return false;
    }

    return true;
}

void config_free(Config *config) {

    if (!config) return;

    free(config->httpAddr);
    free(config->starterHost);
    free(config->notifyHost);
    free(config->stateFile);
    memset(config, 0, sizeof(*config));
}
