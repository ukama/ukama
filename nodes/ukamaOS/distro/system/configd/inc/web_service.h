/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef INC_WEB_SERVICE_H_
#define INC_WEB_SERVICE_H_

#include "config.h"
#include "ulfius.h"

#define EP_BS                           "/"
#define REST_API_VERSION                "v1"
#define URL_PREFIX                      EP_BS REST_API_VERSION
#define API_RES_EP(RES)                 EP_BS RES

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
typedef int (*HttpCb)(const URequest *request, // Input parameters (set by the framework)
                 UResponse* response,         // Output parameters (set by the user)
                void * user_data);

/**
 * @struct WebServiceAPI
 * @brief  Struct used by discover API
 *
 */
typedef struct {
   char method[METHOD_LENGTH];
   char endPoint[URL_EXT_LENGTH];
} WebServiceAPI;

/**
 * @fn      int web_service_start()
 * @brief   Initializes the ulfius framework for REST server.
 *
 * @param   port
 * @return  On success STATUS_OK
 *          On failure STATUS_NOK
 */
int web_service_init(int port);

/**
 * @fn      int web_service_start()
 * @brief   Add API endpoints and start the ulfius HTTP server
 *
 * @return  On success STATUS_OK
 *          On failure STATUS_NOK
 */
int web_service_start();

/**
 * @fn      void web_service_add_node_endpoints()
 * @brief   Add REST API end points to REST framework.
 *
 */
void web_service_add_node_endpoints();

/**
 * @fn      void web_service_exit()
 * @brief   Stop ulfius framework.
 *
 */
void web_service_exit();

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig);

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig);

int web_service_cb_default(const URequest *request,
						   UResponse *response,
                           void *epConfig);

int web_service_cb_post_config(const URequest *request,
                              UResponse *response,
							  void *epConfig);

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data);

#endif /* INC_WEB_SERVICE_H_ */
