/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "web_service.h"

#include "configd.h"
#include "web_client.h"
#include "httpStatus.h"
#include "jserdes.h"
#include "service.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

#include "version.h"

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig) {

	ulfius_set_string_body_response(response, HttpStatus_OK,
			HttpStatusStr(HttpStatus_OK));

	return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

	ulfius_set_string_body_response(response, HttpStatus_NotFound,
			HttpStatusStr(HttpStatus_NotFound));

	return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_MethodNotAllowed,
                                    HttpStatusStr(HttpStatus_MethodNotAllowed));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_config(const URequest *request,
                               UResponse *response,
                               void *epConfig) {

	JsonObj *json = NULL;

	json = ulfius_get_json_body_request(request, NULL);
	if (json == NULL) {
		ulfius_set_string_body_response(response,
                                        HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
		return U_CALLBACK_CONTINUE;
	}

	if (process_received_config(json, (Config *)epConfig)) {
		ulfius_set_string_body_response(response,
                                        HttpStatus_Created,
                                        HttpStatusStr(HttpStatus_Created));
	} else {
		ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
	}

    json_decref(json);
	return U_CALLBACK_CONTINUE;
}





