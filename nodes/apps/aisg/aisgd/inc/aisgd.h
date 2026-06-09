/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISGD_H_
#define AISGD_H_

#include "ulfius.h"
#include "jansson.h"
#include "config.h"
#include "controller_client.h"
#include "status.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_services.h"
#include "usys_types.h"

typedef struct _u_instance UInst;
typedef struct _u_request URequest;
typedef struct _u_response UResponse;
typedef json_t JsonObj;

typedef struct {
    Config *config;
    AppStatus *status;
    ControllerClient controller;
} AisgdContext;

#endif /* AISGD_H_ */
