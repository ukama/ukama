/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Callback functions for various endpoints and REST methods.
 */

#include <ulfius.h>
#include <string.h>

#include "httpStatus.h"
#include "initClient.h"
#include "log.h"

/*
 * decode a u_map into a string
 */
static char *print_map(const struct _u_map * map) {

	char * line, * to_return = NULL;
	const char **keys, * value;

	int len, i;

	if (map != NULL) {
		keys = u_map_enum_keys(map);
		for (i=0; keys[i] != NULL; i++) {
			value = u_map_get(map, keys[i]);
			len = snprintf(NULL, 0, "key is %s, value is %s", keys[i], value);
			line = o_malloc((len+1)*sizeof(char));
			snprintf(line, (len+1), "key is %s, value is %s", keys[i], value);
			if (to_return != NULL) {
				len = o_strlen(to_return) + o_strlen(line) + 1;
				to_return = o_realloc(to_return, (len+1)*sizeof(char));
				if (o_strlen(to_return) > 0) {
					strcat(to_return, "\n");
				}
			} else {
				to_return = o_malloc((o_strlen(line) + 1)*sizeof(char));
				to_return[0] = 0;
			}
			strcat(to_return, line);
			o_free(line);
		}
		return to_return;
	} else {
		return NULL;
	}
}

/*
 * log_request -- log various parameters for the incoming request.
 *
 */
static void log_request(const struct _u_request *request) {

	log_debug("Recevied: %s %s %s", request->http_protocol, request->http_verb,
			  request->http_url);
}

/*
 * callback_default -- default callback 404
 *
 */
int callback_default_webservice(const URequest *request, UResponse *response,
								void *data) {

	ulfius_set_string_body_response(response, 404, "");
	return U_CALLBACK_CONTINUE;
}

/*
 * callback_ping --
 *
 */
int callback_ping(const URequest *request, UResponse *response,
						void *data) {

	ulfius_set_string_body_response(response, 200, "ok");
	return U_CALLBACK_CONTINUE;
}

/*
 * callback_get_systems --
 *
 */
int callback_get_systems(const URequest *request, UResponse *response,
						 void *data) {

	/* GET /systems/?name=system_name */
	int statusCode = 200;
	char *systemName = NULL, *responseStr = NULL;

	log_request(request);

	systemName = (char *)u_map_get(request->map_url, INIT_CLIENT_NAME_STR);
	if (!systemName) {
		log_error("Invalid system name in the GET request for EP: %s.",
				  EP_SYSTEMS);
		statusCode = 400;
		responseStr = msprintf("%s", INIT_CLIENT_ERROR_INVALID_KEY_STR);
		ulfius_set_string_body_response(response, HttpStatus_BadRequest,
										responseStr);

		return U_CALLBACK_CONTINUE;
	}

	/* REST call to init and get requested info about the system */
	if (get_system_info((Config *)data, systemName, &responseStr) == QUERY_OK) {
		ulfius_set_string_body_response(response, HttpStatus_OK, responseStr);
	} else {
		ulfius_set_string_body_response(response, HttpStatus_BadRequest,
										INIT_CLIENT_ERROR_INVALID_SYSTEM_NAME);
	}

	if (responseStr) free(responseStr);

	return U_CALLBACK_CONTINUE;
}
