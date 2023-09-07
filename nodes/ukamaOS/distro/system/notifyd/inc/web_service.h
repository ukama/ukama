/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_WEB_SERVICE_H_
#define INC_WEB_SERVICE_H_

#include "config.h"
#include "ulfius.h"

#define EP_BS                           "/"
#define WEBSERVICE                      "notify"
#define REST_API_VERSION                "v1"

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

#define URL_PREFIX                      EP_BS WEBSERVICE EP_BS REST_API_VERSION
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

int web_service_cb_post_event(const URequest *request,
                              UResponse *response,
                              void *epConfig);

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig);

#endif /* INC_WEB_SERVICE_H_ */
