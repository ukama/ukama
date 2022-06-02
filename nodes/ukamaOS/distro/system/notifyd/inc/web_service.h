/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_WEB_SERVICE_H_
#define INC_WEB_SERVICE_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "ulfius.h"

#define WEB_SERVICE_PORT                8095

#define WEB_SOCKETS                     1
#define WEB_SERVICE                     0

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

#define API_RES_EP(RES)                  EP_BS WEBSERVICE EP_BS \
    REST_API_VERSION EP_BS RES

/* RESPONSE CODE */
#define RESP_CODE_SUCCESS               200
#define RESP_CODE_INVALID_REQUEST       400
#define RESP_CODE_RESOURCE_NOT_FOUND    404
#define RESP_CODE_SERVER_FAILURE        500

#define METHOD_LENGTH                   7
#define URL_EXT_LENGTH                  64
#define MAX_END_POINTS                  64

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


#ifdef __cplusplus
}
#endif

#endif /* INC_WEB_SERVICE_H_ */
