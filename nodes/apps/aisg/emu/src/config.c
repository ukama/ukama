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

#include <toml.h>

#include "config.h"
#include "usys_log.h"
#include "usys_services.h"

int usys_find_service_port(const char *service);

static void copy_str(char *dst, size_t size, const char *src) {
    snprintf(dst, size, "%s", src ? src : "");
}

void emu_config_init(EmuConfig *config) {
    if (config == NULL) {
        return;
    }

    memset(config, 0, sizeof(EmuConfig));
    copy_str(config->socketPath,
             sizeof(config->socketPath),
             AISG_EMU_SOCKET_PATH);
    copy_str(config->scenario, sizeof(config->scenario), "normal");
    config->servicePort = usys_find_service_port(AISG_EMU_SERVICE_NAME);
    if (config->servicePort <= 0) {
        config->servicePort = AISG_EMU_DEFAULT_PORT;
    }
}

bool emu_config_load(EmuConfig *config, const char *file) {
    FILE *fp;
    char errbuf[256];
    toml_table_t *root;
    toml_table_t *service;
    toml_table_t *scenario;
    toml_datum_t value;

    if (config == NULL || file == NULL) {
        return false;
    }

    fp = fopen(file, "r");
    if (fp == NULL) {
        return false;
    }

    root = toml_parse_file(fp, errbuf, sizeof(errbuf));
    fclose(fp);
    if (root == NULL) {
        usys_log_error("failed to parse %s: %s", file, errbuf);
        return false;
    }

    service = toml_table_in(root, "service");
    if (service) {
        value = toml_string_in(service, "socket");
        if (value.ok) {
            copy_str(config->socketPath, sizeof(config->socketPath), value.u.s);
            free(value.u.s);
        }
    }

    scenario = toml_table_in(root, "scenario");
    if (scenario) {
        value = toml_string_in(scenario, "name");
        if (value.ok) {
            copy_str(config->scenario, sizeof(config->scenario), value.u.s);
            free(value.u.s);
        }
    }

    toml_free(root);
    return true;
}

void emu_config_free(EmuConfig *config) {
    (void)config;
}
