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

    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    VERSION);

    return U_CALLBACK_CONTINUE;
}

static const char *config_state_reason(ConfigApplyState state) {

    switch (state) {
    case CONFIG_APPLY_IN_PROGRESS:
        return CONFIG_REASON_IN_PROGRESS;
    case CONFIG_APPLY_APPLIED:
        return CONFIG_REASON_APPLIED;
    case CONFIG_APPLY_FAILED:
        return CONFIG_REASON_FAILED;
    default:
        return CONFIG_REASON_AWAITING;
    }
}

int web_service_cb_ready(const URequest *request,
                         UResponse *response,
                         void *epConfig) {

    Config *config;
    ConfigApplyState state;
    char requestId[CONFIG_REQUEST_ID_LEN];
    JsonObj *json;
    int status;

    (void)request;

    config = (Config *)epConfig;
    if (!config) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(
                                            HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    config_status_snapshot(config,
                           &state,
                           requestId,
                           sizeof(requestId));

    status = HttpStatus_OK;
    if (state == CONFIG_APPLY_IN_PROGRESS) {
        status = HttpStatus_Accepted;
    } else if (state == CONFIG_APPLY_FAILED) {
        status = HttpStatus_ServiceUnavailable;
    }

    json = json_object();
    if (!json) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(
                                            HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    json_object_set_new(json,
                        "ready",
                        json_boolean(state == CONFIG_APPLY_AWAITING ||
                                     state == CONFIG_APPLY_APPLIED));
    json_object_set_new(json,
                        "reason",
                        json_string(config_state_reason(state)));
    if (requestId[0]) {
        json_object_set_new(json,
                            "requestId",
                            json_string(requestId));
    }

    ulfius_set_json_body_response(response, status, json);
    json_decref(json);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

	ulfius_set_string_body_response(response,
                                    HttpStatus_NotFound,
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
