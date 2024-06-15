/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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

/**
 * @fn      int web_service_cb_ping(const URequest*, UResponse*, void*)
 * @brief   reports ping response to client
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return
 */
int web_service_cb_ping(const URequest *request, UResponse *response,
		void *epConfig) {

	ulfius_set_string_body_response(response, HttpStatus_OK,
			HttpStatusStr(HttpStatus_OK));

	return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int web_service_cb_default(const URequest*, UResponse*, void*)
 * @brief   default callback used by REST framework if valid endpoint is not
 *          requested.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
int web_service_cb_default(const URequest *request, UResponse *response,
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

/**
 * @fn      int web_service_cb_post_config(const URequest*, UResponse*, void*)
 * @brief   Receive a new config from service.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
int web_service_cb_post_config(const URequest *request,
		UResponse *response, void *epConfig) {

	int ret = STATUS_NOK;
	char *service=NULL;
	JsonObj *json=NULL;

	json = ulfius_get_json_body_request(request, NULL);
	if (json == NULL) {
		ulfius_set_string_body_response(response,
				HttpStatus_BadRequest,
				HttpStatusStr(HttpStatus_BadRequest));
		return U_CALLBACK_CONTINUE;
	}
	usys_log_trace("config.d:: Received POST for an config from %s.", service);

	ret = configd_process_incoming_config(service,
			json,
			(Config *)epConfig);
	if (ret == STATUS_OK) {
		ulfius_set_empty_body_response(response, HttpStatus_Created);
		usys_log_trace("config.d:: Received POST for an config from %s is responsed with %d.", service, HttpStatus_Created);
	} else {
		ulfius_set_empty_body_response(response,
				HttpStatus_InternalServerError);
		usys_log_trace("config.d:: Received POST for an config from %s is responsed with %d.", service, HttpStatus_InternalServerError);
	}

	json_free(&json);
	return U_CALLBACK_CONTINUE;
}





