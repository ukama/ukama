/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
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
    
    ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));

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

    if (config->clientMode == USYS_FALSE) {
        id = u_map_get(request->map_url, "id");
        if (id == NULL) {
            ulfius_set_string_body_response
                (response, HttpStatus_BadRequest,
                 HttpStatusStr(HttpStatus_BadRequest));
            return U_CALLBACK_CONTINUE;
        } else if (strcmp(id, config->nodeID) != 0) {
            ulfius_set_string_body_response
                (response, HttpStatus_BadRequest,
                 HttpStatusStr(HttpStatus_BadRequest));
            return U_CALLBACK_CONTINUE;
        }
    }

    /* Send alarm to notify.d, wait few sec and reboot linux */
    process_reboot(config);
    ulfius_set_empty_body_response(response, HttpStatus_Accepted);

    return U_CALLBACK_CONTINUE;
}






