/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef FEMD_H
#define FEMD_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <unistd.h>
#include <errno.h>

#include <getopt.h>
#include <stdint.h>
#include <stdbool.h>

#include "config.h"
#include "version.h"
#include "gpio_controller.h"
#include "i2c_controller.h"
#include "web_api.h"
#include "yaml_config.h"
#include "safety_monitor.h"

// Use simple printf on macOS for development
#ifdef __APPLE__
#define usys_log_info(fmt, ...) printf("[INFO] " fmt "\n", ##__VA_ARGS__)
#define usys_log_error(fmt, ...) printf("[ERROR] " fmt "\n", ##__VA_ARGS__)
#define usys_log_warn(fmt, ...) printf("[WARN] " fmt "\n", ##__VA_ARGS__)
#define usys_log_debug(fmt, ...) printf("[DEBUG] " fmt "\n", ##__VA_ARGS__)
#define usys_log_set_service(name) ((void)0)
#define usys_log_remote_init(name) ((void)0)
#else
#include "usys_log.h"
#endif

#define SERVICE_NAME              "femd"
#define FEM_VERSION               VERSION

#define STATUS_OK                 0
#define STATUS_NOK               -1

#define DEF_LOG_LEVEL            "INFO"

// Global variables
extern volatile sig_atomic_t g_running;

// Function declarations
void handle_sigint(int signum);
void print_usage(const char *program);
void print_version(void);

#endif /* FEMD_H */