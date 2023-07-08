/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "web_service.h"
#include "web_client.h"
#include "http_status.h"
#include "config.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

extern void process_reboot(Config *config);

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {
    
    ulfius_set_string_body_response(response, HttpStatus_Unauthorized,
                                    HttpStatusStr(HttpStatus_Unauthorized));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_restart(const URequest *request,
                                UResponse *response,
                                void *epConfig) {

    int ret = STATUS_NOK;
    char *id=NULL;

    id = u_map_get(request->map_url, "id");
    ulfius_set_empty_body_response(response, HttpStatus_Accepted);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_reboot(const URequest *request,
                                UResponse *response,
                                void *epConfig) {

    char *id=NULL;
    Config *config=NULL;

    config = (Config *)epConfig;

    id = u_map_get(request->map_url, "id");
    if (id == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
    } else if (strcmp(id, config->nodeID) != 0) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
    } else {

        /* Send alarm to notify.d, wait few sec and reboot linux */
        process_reboot(config);

        ulfius_set_empty_body_response(response, HttpStatus_Accepted);
    }

    return U_CALLBACK_CONTINUE;
}






