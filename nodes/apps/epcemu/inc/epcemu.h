/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef EPCEMU_H_
#define EPCEMU_H_

#include <stdbool.h>
#include <stdint.h>
#include <pthread.h>

#include "ulfius.h"
#include "jansson.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_services.h"
#include "usys_string.h"
#include "usys_types.h"

typedef struct _u_instance UInst;
typedef struct _u_request URequest;
typedef struct _u_response UResponse;
typedef json_t JsonObj;

#define EPCEMU_SERVICE_NAME       "epcemu"
#define EPCEMU_APP_NAME           "epcemu.d"

#define EPCEMU_PCRF_SERVICE       "pcrf"
#define EPCEMU_INITNET_SERVICE    "init-network"

#define EPCEMU_DEF_APN            "ukama"
#define EPCEMU_DEF_LOG_LEVEL      "TRACE"
#define EPCEMU_MAX_STR            256
#define EPCEMU_MAX_REASON         256

#define EP_BS                     "/"
#define REST_API_VERSION          "v1"
#define URL_PREFIX                EP_BS REST_API_VERSION
#define API_RES_EP(RES)           EP_BS RES

typedef enum {
    EpcemuStateStarting = 0,
    EpcemuStateResolvingServices,
    EpcemuStateCheckingInitNetwork,
    EpcemuStateCheckingPcrf,
    EpcemuStateReady,
    EpcemuStateFailed
} EpcemuState;

typedef struct {
    int servicePort;
    int pcrfPort;
    int initNetworkPort;

    char pcrfUrl[EPCEMU_MAX_STR];
    char initNetworkUrl[EPCEMU_MAX_STR];
    char ueCidr[EPCEMU_MAX_STR];
    char bridge[EPCEMU_MAX_STR];
    char bridgeCidr[EPCEMU_MAX_STR];

    bool pcrfReady;
    bool initNetworkReady;
} EpcemuConfig;

typedef struct {
    EpcemuState state;
    bool ready;
    char reason[EPCEMU_MAX_REASON];
    pthread_mutex_t mutex;
} EpcemuStatus;

typedef struct {
    EpcemuConfig *config;
    EpcemuStatus *status;
} ServiceContext;

void status_init(EpcemuStatus *status);
void status_destroy(EpcemuStatus *status);
void status_set(EpcemuStatus *status,
                EpcemuState state,
                const char *reason);
void status_fail(EpcemuStatus *status, const char *reason);
bool status_is_ready(EpcemuStatus *status);
const char *status_state_str(EpcemuState state);

#endif /* EPCEMU_H_ */
