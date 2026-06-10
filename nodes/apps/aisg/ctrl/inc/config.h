/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

#include <stdbool.h>

#include "usys_services.h"

#define AISG_CTRL_SERVICE_NAME  SERVICE_AISG_CTRL
#define DEF_CONFIG_FILE         "/ukama/configs/aisg-ctrl/config.toml"
#define DEF_LOG_LEVEL           "TRACE"
#define DEF_CTRL_SOCKET         "/var/run/aisg-ctrl.sock"
#define DEF_BACKEND_TYPE        "raw-rs485"
#define DEF_RAW_RS485_DEVICE    "/dev/ttyUSB0"
#define DEF_RAW_RS485_BAUD      9600
#define DEF_STM_UART_DEVICE     "/dev/ttyAMA1"
#define DEF_STM_UART_BAUD       115200
#define CONFIG_MAX_STR          256

typedef enum {
    BackendTypeRawRs485 = 0,
    BackendTypeStmUart
} BackendType;

typedef struct {
    char *socketPath;
    BackendType backendType;
    char *rawDevice;
    int  rawBaud;
    char *stmDevice;
    int  stmBaud;
} Config;

void config_set_defaults(Config *config);
bool config_load_from_file(Config *config, const char *path);
void config_free(Config *config);
const char *config_backend_type_str(BackendType type);

#endif /* CONFIG_H_ */
