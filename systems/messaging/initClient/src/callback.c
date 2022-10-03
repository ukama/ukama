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

#include "initClient.h"

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
