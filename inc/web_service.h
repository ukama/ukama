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

#define WEB_SOCKETS                     1
#define WEB_SERVICE                     0

#define EP_BS                           "/"
#define WEBSERVICE                      "noded"
#define REST_API_VERSION                "v1"



/* API End Points */
#define UUID                            "uuid"
#define DEVTYPE                         "type"
#define DEVNAME                         "name"
#define DEVDESC                         "desc"
#define PROPNAME                        "prop"
#define MFGDATA                         "name"

#define API_RES_EP(RES)                  EP_BS WEBSERVICE EP_BS \
    REST_API_VERSION EP_BS RES

#define STATUS_OK                       0
#define STATUS_NOK                      (-1)

// RESPONSE CODE
#define RESP_CODE_SUCCESS               200
#define RESP_CODE_INVALID_REQUEST       400
#define RESP_CODE_RESOURCE_NOT_FOUND    404
#define RESP_CODE_SERVER_FAILURE        500

#define SEND_ADD_SUBSCRIBER_REQ_FAILED  1

#define WEB_SERVICE_PORT                8085

#define METHOD_LENGTH                   7
#define URL_EXT_LENGTH                  64
#define MAX_END_POINTS                  64

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;

typedef int (*HttpCb)(const URequest *request, // Input parameters (set by the framework)
                 UResponse* response,         // Output parameters (set by the user)
                void * user_data);

typedef struct {
   char method[METHOD_LENGTH];
   char endPoint[URL_EXT_LENGTH];
} WebServiceAPI;


int web_service_init();
int web_service_start();

void web_service_add_unit_endpoints();
void web_service_add_module_endpoints();

#ifdef __cplusplus
}
#endif

#endif /* INC_WEB_SERVICE_H_ */
