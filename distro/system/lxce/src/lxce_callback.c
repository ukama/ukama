/**
 * Copyright (c) 2021-present, Ukama Inc.
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

#include "lxce_callback.h"
#include "log.h"

/*
 * callback_not_allowed -- 
 *
 */
int callback_not_allowed(const URequest *request, UResponse *response,
			 void *user_data) {
  
  ulfius_set_string_body_response(response, 403, "Operation not allowed.");
  return U_CALLBACK_CONTINUE;
}

/*
 * callback_webservice --
 *
 */
int callback_webservice(const URequest *request, UResponse *response,
			void *data) {
  
  ulfius_set_string_body_response(response, 200, "Next time");
  return U_CALLBACK_CONTINUE;
}

/*
 * callback_default -- default callback for no-match
 *
 *
 */
int callback_default(const struct _u_request *request,
		     struct _u_response *response, void *user_data) {
  
  ulfius_set_string_body_response(response, 404, "You are clearly high!");
  return U_CALLBACK_CONTINUE;
}
