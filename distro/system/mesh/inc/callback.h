/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef CALLBACK_H
#define CALLBACK_H

#include <ulfius.h>

int callback_websocket(const struct _u_request *request,
		       struct _u_response *response, void *user_data);
int callback_not_allowed(const struct _u_request *request,
			 struct _u_response *response, void *user_data);
int callback_default(const struct _u_request *request,
		     struct _u_response *response, void *user_data);
#endif /* CALLBACK_H */
