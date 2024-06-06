/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef WEB_SERVICE_H_
#define WEB_SERVICE_H_

#include <ulfius.h>

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *data);

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *data);

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *data);

int web_service_cb_get_level(const URequest *request,
                             UResponse *response,
                             void *data);

int web_service_cb_get_output(const URequest *request,
                              UResponse *response,
                              void *data);

int web_service_cb_post_level(const URequest *request,
                              UResponse *response,
                              void *data);

int web_service_cb_post_output(const URequest *request,
                               UResponse *response,
                               void *data);

#endif /* WEB_SERVICE_H_ */
