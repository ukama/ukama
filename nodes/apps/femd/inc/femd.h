/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef FEMD_H
#define FEMD_H

#include <pthread.h>
#include <signal.h>
#include <getopt.h>
#include <stdbool.h>

#include "config.h"
#include "version.h"
#include "gpio_controller.h"
#include "i2c_controller.h"
#include "yaml_config.h"
#include "safety_monitor.h"

#include "ulfius.h"
#include "usys_types.h"
#include "usys_services.h"
#include "usys_log.h"
#include "jansson.h"

#define SERVICE_NAME              SERVICE_FEMD
#define FEM_VERSION               VERSION

#define STATUS_OK                 (0)
#define STATUS_NOK                (-1)
#define STATUS_NOTOK              (-1)

#define DEF_LOG_LEVEL             "INFO"
#define DEF_CONFIG_FILE           "./config/femd.conf"

#define ERR_FEMD_JSON_CREATION_ERR (-1)

#define EP_BS                     "/"
#define REST_API_VERSION          "v1"
#define URL_PREFIX                EP_BS REST_API_VERSION
#define API_RES_EP(RES)           EP_BS RES

typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

extern volatile sig_atomic_t g_running;

void handle_sigint(int signum);

int config_init(Config *config);
void config_free(Config *config);
int config_load_from_file(Config *config, const char *filename);
void config_print(const Config *config);

#endif /* FEMD_H */
