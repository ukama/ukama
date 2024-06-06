/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef WEB_SERVICE_H_
#define WEB_SERVICE_H_

#include "starter.h"

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig);

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig);

int web_service_cb_get_status(const URequest *request,
                              UResponse *response,
                              void *epConfig);

int web_service_cb_post_terminate(const URequest *request,
                                  UResponse *response,
                                  void *epConfig);

int web_service_cb_post_update(const URequest *request,
                               UResponse *response,
                               void *epConfig);

int web_service_cb_post_exec(const URequest *request,
                             UResponse *response,
                             void *epConfig);

int web_service_cb_get_all_capps_status(const URequest *request,
                                        UResponse *response,
                                        void *epConfig);

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig);

#endif /* WEB_SERVICE_H_ */
