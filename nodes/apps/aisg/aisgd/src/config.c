/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "toml.h"
#include "config.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_services.h"

int usys_find_service_port(const char *service);

static char *toml_string_or(toml_table_t *table,
                            const char *key,
                            const char *defValue) {
    toml_datum_t datum;

    if (table != NULL) {
        datum = toml_string_in(table, key);
        if (datum.ok) return datum.u.s;
    }

    return defValue ? strdup(defValue) : NULL;
}

static int toml_int_or(toml_table_t *table, const char *key, int defValue) {
    toml_datum_t datum;

    if (table == NULL) return defValue;

    datum = toml_int_in(table, key);
    return datum.ok ? (int)datum.u.i : defValue;
}

static bool toml_bool_or(toml_table_t *table,
                         const char *key,
                         bool defValue) {
    toml_datum_t datum;

    if (table == NULL) return defValue;

    datum = toml_bool_in(table, key);
    return datum.ok ? datum.u.b : defValue;
}

static void resolve_service_port(Config *config, int fallbackPort) {
    int port;

    if (config == NULL || config->serviceName == NULL) return;

    port = usys_find_service_port(config->serviceName);
    config->servicePort = port > 0 ? port : fallbackPort;
}

void config_set_defaults(Config *config) {
    if (config == NULL) return;

    memset(config, 0, sizeof(Config));
    config->serviceName = strdup(AISGD_SERVICE_NAME);
    config->servicePort = DEF_SERVICE_PORT;
    config->controllerPath = strdup(DEF_CTRL_SOCKET);
    config->controllerTimeoutMs = DEF_CTRL_TIMEOUT_MS;
    config->requireConfigBeforeCalibrate = true;
    config->requireCalibrateBeforeSetTilt = true;
    config->stateFile = strdup(DEF_STATE_FILE);
}

bool config_load_from_file(Config *config, const char *path) {
    FILE *fp;
    char errbuf[256];
    toml_table_t *root;
    toml_table_t *controller;
    toml_table_t *policy;
    toml_table_t *files;
    int fallbackPort;

    if (config == NULL || path == NULL) return false;

    fallbackPort = config->servicePort;
    fp = fopen(path, "r");
    if (fp == NULL) {
        resolve_service_port(config, fallbackPort);
        return true;
    }

    root = toml_parse_file(fp, errbuf, sizeof(errbuf));
    fclose(fp);
    if (root == NULL) {
        usys_log_error("failed to parse config %s: %s", path, errbuf);
        return false;
    }

    controller = toml_table_in(root, "controller");
    policy     = toml_table_in(root, "policy");
    files      = toml_table_in(root, "files");

    usys_free(config->controllerPath);
    config->controllerPath = toml_string_or(controller,
                                             "path",
                                             DEF_CTRL_SOCKET);
    config->controllerTimeoutMs = toml_int_or(controller,
                                             "timeout_ms",
                                             DEF_CTRL_TIMEOUT_MS);

    config->requireConfigBeforeCalibrate =
        toml_bool_or(policy, "require_config_before_calibrate", true);
    config->requireCalibrateBeforeSetTilt =
        toml_bool_or(policy, "require_calibrate_before_set_tilt", true);

    usys_free(config->stateFile);
    config->stateFile = toml_string_or(files, "state", DEF_STATE_FILE);

    resolve_service_port(config, fallbackPort);
    toml_free(root);

    return true;
}

void config_free(Config *config) {
    if (config == NULL) return;

    usys_free(config->serviceName);
    usys_free(config->controllerPath);
    usys_free(config->stateFile);
    memset(config, 0, sizeof(Config));
}
