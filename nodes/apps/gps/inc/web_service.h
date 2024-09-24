/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef WEB_SERVICE_H_
#define WEB_SERVICE_H_

#include "gpsd.h"

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *config);

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *config);

int web_service_cb_lock(const URequest *request,
                        UResponse *response,
                        void *config);

int web_service_cb_coordinates(const URequest *request,
                               UResponse *response,
                               void *config);

int web_service_cb_time(const URequest *request,
                        UResponse *response,
                        void *config);

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *config);

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *config);

#endif /* WEB_SERVICE_H_ */
