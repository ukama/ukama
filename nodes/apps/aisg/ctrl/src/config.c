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
#include "usys_mem.h"

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

const char *config_backend_type_str(BackendType type) {
    switch (type) {
    case BackendTypeRawRs485: return "raw-rs485";
    case BackendTypeStmUart:  return "stm-uart";
    default:                 return "unknown";
    }
}

static BackendType backend_type_from_str(const char *type) {
    if (type != NULL && !strcmp(type, "stm-uart")) return BackendTypeStmUart;
    return BackendTypeRawRs485;
}

void config_set_defaults(Config *config) {
    memset(config, 0, sizeof(Config));
    config->socketPath = strdup(DEF_CTRL_SOCKET);
    config->backendType = BackendTypeRawRs485;
    config->rawDevice = strdup(DEF_RAW_RS485_DEVICE);
    config->rawBaud = DEF_RAW_RS485_BAUD;
    config->stmDevice = strdup(DEF_STM_UART_DEVICE);
    config->stmBaud = DEF_STM_UART_BAUD;
}

bool config_load_from_file(Config *config, const char *path) {
    FILE *fp;
    char errbuf[256];
    toml_table_t *root;
    toml_table_t *service;
    toml_table_t *backend;
    toml_table_t *raw;
    toml_table_t *stm;
    char *type;

    fp = fopen(path, "r");
    if (fp == NULL) return true;

    root = toml_parse_file(fp, errbuf, sizeof(errbuf));
    fclose(fp);
    if (root == NULL) return false;

    service = toml_table_in(root, "service");
    backend = toml_table_in(root, "backend");
    raw = toml_table_in(root, "raw_rs485");
    stm = toml_table_in(root, "stm_uart");

    usys_free(config->socketPath);
    config->socketPath = toml_string_or(service, "socket", DEF_CTRL_SOCKET);

    type = toml_string_or(backend, "type", DEF_BACKEND_TYPE);
    config->backendType = backend_type_from_str(type);
    usys_free(type);

    usys_free(config->rawDevice);
    config->rawDevice = toml_string_or(raw, "device", DEF_RAW_RS485_DEVICE);
    config->rawBaud = toml_int_or(raw, "baud", DEF_RAW_RS485_BAUD);

    usys_free(config->stmDevice);
    config->stmDevice = toml_string_or(stm, "device", DEF_STM_UART_DEVICE);
    config->stmBaud = toml_int_or(stm, "baud", DEF_STM_UART_BAUD);

    toml_free(root);
    return true;
}

void config_free(Config *config) {
    usys_free(config->socketPath);
    usys_free(config->rawDevice);
    usys_free(config->stmDevice);
    memset(config, 0, sizeof(Config));
}
