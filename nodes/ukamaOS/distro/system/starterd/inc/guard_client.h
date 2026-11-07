/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef GUARD_CLIENT_H_
#define GUARD_CLIENT_H_

#include <jansson.h>

#define HTTP_TRANSPORT_ERROR 0

int guard_http(const char *method, const char *url, json_t *body,
               json_t **reply);

#endif /* GUARD_CLIENT_H_ */
