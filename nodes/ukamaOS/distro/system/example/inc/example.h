/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef EXAMPLE_H_
#define EXAMPLE_H_

#include <ulfius.h>

#include "usys_log.h"
#include "usys_services.h"
#include "usys_types.h"

#define SERVICE_NAME              "example"

#define DEF_LOG_LEVEL             "TRACE"
#define DEF_SERVICE_PORT          18100

#define ENV_EXAMPLE_SERVICE_PORT  "EXAMPLE_SERVICE_PORT"

#define EP_BS                     "/"
#define REST_API_VERSION          "v1"
#define URL_PREFIX                EP_BS REST_API_VERSION
#define API_RES_EP(RES)           EP_BS RES

typedef struct _u_instance UInst;
typedef struct _u_request  URequest;
typedef struct _u_response UResponse;

typedef struct {

    int servicePort;
} Config;

#endif /* EXAMPLE_H_ */
