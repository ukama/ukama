/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef INC_WEB_SERVICE_H_
#define INC_WEB_SERVICE_H_

#include "ulfius.h"

#define EP_BS                           "/"
#define REST_API_VERSION                "v1"
#define URL_PREFIX                      EP_BS REST_API_VERSION
#define API_RES_EP(RES)                 EP_BS RES

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
    
int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *data);
int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *data);
int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *data);
int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *data);

#endif /* INC_WEB_SERVICE_H_ */
