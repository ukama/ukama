/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <ulfius.h>

#include "websocket.h"
#include "http_status.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

int web_socket_cb_ping(const URequest *request,
                       UResponse *response,
                       void *data) {

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

int web_socket_cb_default(const URequest *request,
                          UResponse *response,
                          void *data) {
    
    ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));

    return U_CALLBACK_CONTINUE;
}

int web_socket_cb_post_log(const URequest *request,
                            UResponse *response,
                            void *data) {

    if (ulfius_set_websocket_response(response, NULL, NULL,
                                      &websocket_manager, NULL,
                                      &websocket_incoming_message, data,
                                      &websocket_onclose, NULL) == U_OK) {
		ulfius_add_websocket_deflate_extension(response);
	} else {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
    }

    return U_CALLBACK_CONTINUE;
}
