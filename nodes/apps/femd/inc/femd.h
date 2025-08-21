/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
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

#define SERVICE_NAME              SERVICE_FEM
#define FEM_VERSION               VERSION

#define STATUS_OK                 (0)
#define STATUS_NOK                (-1)
#define STATUS_NOTOK              (-1)

#define DEF_LOG_LEVEL             "INFO"
#define DEF_SERVICE_CLIENT_HOST   "localhost"
#define DEF_NODED_HOST            "localhost"
#define DEF_NOTIFY_HOST           "localhost"
#define DEF_NODED_EP              "/v1/nodeinfo"
#define DEF_NOTIFY_EP             "/notify/v1/event/"
#define DEF_NODE_ID               "ukama-aaa-bbbb-ccc-dddd"
#define DEF_NODE_TYPE             "amplifier"
#define ENV_FEMD_DEBUG_MODE       "FEMD_DEBUG_MODE"

#define EP_BS                     "/"
#define REST_API_VERSION          "v1"
#define URL_PREFIX                EP_BS REST_API_VERSION
#define API_RES_EP(RES)           EP_BS RES

typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

typedef struct {

    Config         *config;
    GpioController *gpioController;
    I2CController  *i2cController;
} ServerConfig;

void handle_sigint(int signum);

int config_init(Config *config);
void config_free(Config *config);
int config_load_from_file(Config *config, const char *filename);
void config_print(const Config *config);

#endif /* FEMD_H */
