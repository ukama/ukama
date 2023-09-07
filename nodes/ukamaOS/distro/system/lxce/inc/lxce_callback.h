/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef LXCE_CALLBACK_H
#define LXCE_CALLBACK_H

#include <ulfius.h>

typedef struct _u_instance UInst;
typedef struct _u_request  URequest;
typedef struct _u_response UResponse;
typedef struct _u_map UMap;

int callback_webservice(const URequest *request, UResponse *response,
			void *user_data);
int callback_not_allowed(const URequest *request, UResponse *response,
			 void *user_data);
int callback_default(const URequest *request, UResponse *response,
		     void *user_data);

#endif /* LXCE_CALLBACK_H */
