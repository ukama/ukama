/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "web_service.h"
#include <string.h>

#include "configd.h"
#include "web_client.h"
#include "http_status.h"
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

	ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
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

int web_service_cb_ready(const URequest *request,
                         UResponse *response,
                         void *epConfig) {
    JsonObj *json;

    (void)request;

    if (!epConfig) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        "unavailable");
        return U_CALLBACK_CONTINUE;
    }

    json = json_pack("{s:b,s:s}", "ready", 1, "reason", "ready");
    if (!json) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        "allocation_failed");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_json_body_response(response,
                                  HttpStatus_OK,
                                  json);
    json_decref(json);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_config_status(const URequest *request,
                                  UResponse *response,
                                  void *epConfig) {

    Config *config = epConfig;
    ConfigRecord record;
    JsonObj *json;
    const char *mode;
    const char *phase;

    (void)request;
    if (!config || !config->stateStore) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        "state_unavailable");
        return U_CALLBACK_CONTINUE;
    }

    config_store_snapshot(config->stateStore, &record);

    mode = (record.mode == CONFIG_MODE_NOCONFIG ? "NOCONFIG" :
            record.mode == CONFIG_MODE_CONFIG ? "CONFIG" : "NONE");
    phase = (record.phase == CONFIG_PHASE_COMPLETED ? "completed" :
             record.phase == CONFIG_PHASE_PENDING ? "pending" :
             record.phase == CONFIG_PHASE_FAILED ? "failed" : "awaiting");
    json = json_pack("{s:i,s:s,s:s,s:s,s:I,s:i,s:s}",
                     "schemaVersion", 1,
                     "mode", mode,
                     "phase", phase,
                     "requestId", record.requestId,
                     "generation", (json_int_t)record.generation,
                     "revision", record.revision,
                     "error", record.error);
    if (!json) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        "allocation_failed");
        return U_CALLBACK_CONTINUE;
    }

    /* A failed transaction is a valid status response, not a transport error. */
    ulfius_set_json_body_response(response,
                                  HttpStatus_OK,
                                  json);
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

	JsonObj *json      = NULL;
    JsonObj *mode      = NULL;
    JsonObj *requestId = NULL;

    int result;
    
	json = ulfius_get_json_body_request(request, NULL);
	if (json == NULL) {
		ulfius_set_string_body_response(response,
                                        HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
		return U_CALLBACK_CONTINUE;
	}

    if (!json_is_object(json)) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_BadRequest,
                                        "expected_object");
        json_decref(json);
        return U_CALLBACK_CONTINUE;
    }

    if (json_object_get(json, "mode")) {
        mode      = json_object_get(json, "mode");
        requestId = json_object_get(json, "requestId");

        if (!json_is_string(mode) ||
            json_string_length(mode) != strlen(json_string_value(mode)) ||
            strcmp(json_string_value(mode), "NOCONFIG") != 0 ||
            !json_is_string(requestId) ||
            json_string_length(requestId) != strlen(json_string_value(requestId)) ||
            json_object_size(json) != 2) {
            ulfius_set_string_body_response(response,
                                            HttpStatus_BadRequest,
                                            "expected_NOCONFIG_and_requestId");
            json_decref(json);
            return U_CALLBACK_CONTINUE;
        }

        result = process_noconfig(json_string_value(requestId), epConfig);
        json_decref(json);
        if (result == HttpStatus_OK) {
            return web_service_cb_config_status(request, response, epConfig);
        }
        ulfius_set_string_body_response(response, result,
            result == HttpStatus_BadRequest ? "invalid_requestId" :
            result == HttpStatus_Conflict   ? "configuration_conflict" : "state_unavailable");
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

int web_service_cb_delete_config(const URequest *request,
                                 UResponse *response,
                                 void *epConfig) {

    int result;

    if (request->binary_body_length != 0) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        "DELETE_does_not_accept_a_body");
        return U_CALLBACK_CONTINUE;
    }

    result = process_delete_config(epConfig);
    if (result == HttpStatus_OK) {
        return web_service_cb_config_status(request, response, epConfig);
    }

    ulfius_set_string_body_response(response, result, "state_unavailable");
    return U_CALLBACK_CONTINUE;
}
