/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef CALLBACK_H
#define CALLBACK_H

#include <ulfius.h>

int callback_ping(const URequest *request, UResponse *response,
			       void *user_data);
int callback_get_systems(const URequest *request, UResponse *response,
						 void *user_data);
int callback_not_allowed(const URequest *request, UResponse *response,
			 void *user_data);
int callback_default_webservice(const URequest *request, UResponse *response,
				void *user_data);
#endif /* CALLBACK_H */
