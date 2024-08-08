/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef INC_WEB_SERVICE_H_
#define INC_WEB_SERVICE_H_

#include "config.h"
#include "ulfius.h"

/* API URL parameters*/
#define UUID                            "uuid"
#define DEVTYPE                         "type"
#define DEVNAME                         "name"
#define DEVDESC                         "desc"
#define PROPNAME                        "prop"
#define MFGDATA                         "name"
#define ALERTSTATE                      "state"
#define ENABLE                          "enable"
#define DISABLE                         "disable"

#define EP_BS                           "/"
#define REST_API_VERSION                "v1"
#define URL_PREFIX                      EP_BS REST_API_VERSION
#define API_RES_EP(RES)                 EP_BS RES

#define METHOD_LENGTH                   7
#define URL_EXT_LENGTH                  64
#define MAX_END_POINTS                  64
#define MAX_URL_LENGTH                  128

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
    
/* Callback function used by ulfius when APU is request */
typedef int (*HttpCb)(const URequest *request,
                      UResponse* response,
                      void * user_data);

typedef struct {
   char method[METHOD_LENGTH];
   char endPoint[URL_EXT_LENGTH];
} WebServiceAPI;

int web_service_init(int port);
int web_service_start();
void web_service_add_node_endpoints();
void web_service_exit();

int web_service_cb_ping(const URequest *request, UResponse *response, void *data);
int web_service_cb_version(const URequest *request, UResponse *response, void *data);
int web_service_cb_post_event(const URequest *request, UResponse *response, void *data);
int web_service_cb_default(const URequest *request, UResponse *response, void *data);
int web_service_cb_not_allowed(const URequest *request, UResponse *response, void *data);
int web_service_cb_get_output(const URequest *request, UResponse *response, void *data);
int web_service_cb_post_output(const URequest *request, UResponse *response, void *data);
int web_service_cb_get_count(const URequest *request, UResponse *response, void *data);

#endif /* INC_WEB_SERVICE_H_ */
