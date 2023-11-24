/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "web_service.h"
#include "http_status.h"

int web_service_cb_ping(const URequest *request, UResponse *response,
                        void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request, UResponse *response,
                           void *epConfig) {
    
    ulfius_set_string_body_response(response, HttpStatus_Forbidden,
                                    HttpStatusStr(HttpStatus_Forbidden));

    return U_CALLBACK_CONTINUE;
}







