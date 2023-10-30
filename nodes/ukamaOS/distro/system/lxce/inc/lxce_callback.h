/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
