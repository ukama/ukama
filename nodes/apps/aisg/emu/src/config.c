/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <ctype.h>
#include <math.h>
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

bool emu_config_set_mode(EmuConfig *config, const char *mode)
{
    if (config == NULL || mode == NULL) {
        return false;
    }

    if (!strcmp(mode, "contract")) {
        config->mode = EmuModeContract;
        copy_str(config->modeName, sizeof(config->modeName), "contract");
        return true;
    }

    if (!strcmp(mode, "ret")) {
        config->mode = EmuModeRet;
        copy_str(config->modeName, sizeof(config->modeName), "ret");
        return true;
    }

    return false;
}

bool emu_config_set_bool(int *dst, const char *value)
{
    if (dst == NULL || value == NULL) {
        return false;
    }

    if (!strcmp(value, "1") || !strcmp(value, "true") ||
        !strcmp(value, "yes") || !strcmp(value, "on")) {
        *dst = 1;
        return true;
    }

    if (!strcmp(value, "0") || !strcmp(value, "false") ||
        !strcmp(value, "no") || !strcmp(value, "off")) {
        *dst = 0;
        return true;
    }

    return false;
}

int emu_config_tilt_arg_to_tenths(const char *value, int fallback)
{
    char *end = NULL;
    double tilt;

    if (value == NULL) {
        return fallback;
    }

    tilt = strtod(value, &end);
    if (end == value) {
        return fallback;
    }

    return (int)((tilt * 10.0) + ((tilt >= 0.0) ? 0.5 : -0.5));
}

bool emu_config_init(EmuConfig *config) {

    if (config == NULL) {
        return false;
    }

    memset(config, 0, sizeof(EmuConfig));

    emu_config_set_mode(config, "contract");
    copy_str(config->socketPath,
             sizeof(config->socketPath),
             AISG_EMU_SOCKET_PATH);
    copy_str(config->scenario, sizeof(config->scenario), "normal");

    copy_str(config->retPtyPath,
             sizeof(config->retPtyPath),
             AISG_EMU_RET_PTY_PATH);
    copy_str(config->retVendorCode, sizeof(config->retVendorCode), "UK");
    copy_str(config->retSerial,
             sizeof(config->retSerial),
             "UKAMA00000000001");
    config->retRequiresConfig = 1;
    config->retInitialTiltTenths = 30;
    config->retMinTiltTenths = 0;
    config->retMaxTiltTenths = 100;
    config->retCalibrateDelayMs = 0;
    config->retMoveDelayMs = 0;

    config->servicePort = usys_find_service_port(AISG_EMU_SERVICE_NAME);
    if (config->servicePort <= 0) {
        usys_log_error("Unable to find service port");
        return false;
    }

    return true;
}

static void load_string(toml_table_t *table,
                        const char *key,
                        char *dst,
                        size_t size)
{
    toml_datum_t value;

    if (table == NULL || key == NULL || dst == NULL) {
        return;
    }

    value = toml_string_in(table, key);
    if (value.ok) {
        copy_str(dst, size, value.u.s);
        free(value.u.s);
    }
}

static void load_int(toml_table_t *table, const char *key, int *dst)
{
    toml_datum_t value;

    if (table == NULL || key == NULL || dst == NULL) {
        return;
    }

    value = toml_int_in(table, key);
    if (value.ok) {
        *dst = (int)value.u.i;
    }
}

static void load_bool(toml_table_t *table, const char *key, int *dst)
{
    toml_datum_t value;

    if (table == NULL || key == NULL || dst == NULL) {
        return;
    }

    value = toml_bool_in(table, key);
    if (value.ok) {
        *dst = value.u.b ? 1 : 0;
    }
}

bool emu_config_load(EmuConfig *config, const char *file) {

    FILE *fp;
    char errbuf[256];
    toml_table_t *root;
    toml_table_t *service;
    toml_table_t *scenario;
    toml_table_t *ret;
    toml_datum_t value;

    if (config == NULL || file == NULL) {
        return false;
    }

    fp = fopen(file, "r");
    if (fp == NULL) {
        usys_log_error("Unable to open config file: %s", file);
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
        load_string(service, "socket", config->socketPath, sizeof(config->socketPath));
        value = toml_string_in(service, "mode");
        if (value.ok) {
            if (!emu_config_set_mode(config, value.u.s)) {
                usys_log_warn("unknown aisg-emu mode in config: %s", value.u.s);
            }
            free(value.u.s);
        }
    }

    scenario = toml_table_in(root, "scenario");
    if (scenario) {
        load_string(scenario, "name", config->scenario, sizeof(config->scenario));
    }

    ret = toml_table_in(root, "ret");
    if (ret) {
        load_string(ret, "pty", config->retPtyPath, sizeof(config->retPtyPath));
        load_string(ret, "vendor", config->retVendorCode, sizeof(config->retVendorCode));
        load_string(ret, "serial", config->retSerial, sizeof(config->retSerial));
        load_bool(ret, "requires_config", &config->retRequiresConfig);
        load_int(ret, "initial_tilt_tenths", &config->retInitialTiltTenths);
        load_int(ret, "min_tilt_tenths", &config->retMinTiltTenths);
        load_int(ret, "max_tilt_tenths", &config->retMaxTiltTenths);
        load_int(ret, "calibrate_delay_ms", &config->retCalibrateDelayMs);
        load_int(ret, "move_delay_ms", &config->retMoveDelayMs);
    }

    toml_free(root);
    return true;
}

void emu_config_free(EmuConfig *config) {
    (void)config;
}
