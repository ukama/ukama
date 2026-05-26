/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef HTTP_CLIENT_H_
#define HTTP_CLIENT_H_

#include "jansson.h"

int http_get_json(const char *url,
                  JsonObj **outJson,
                  long *httpCode);

int http_send_json(const char *method,
                   const char *url,
                   JsonObj *body,
                   JsonObj **outJson,
                   long *httpCode);

int http_send_json_timeout(const char *method,
                           const char *url,
                           JsonObj *body,
                           JsonObj **outJson,
                           long *httpCode,
                           int timeoutSec);

#endif /* HTTP_CLIENT_H_ */
