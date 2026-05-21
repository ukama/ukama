/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WEB_SERVICE_H_
#define WEB_SERVICE_H_

#include "epcemu.h"

int start_web_service(ServiceContext *ctx, UInst *serviceInst);

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *data);

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *data);

int web_service_cb_status(const URequest *request,
                          UResponse *response,
                          void *data);

int web_service_cb_attach(const URequest *request,
                          UResponse *response,
                          void *data);

int web_service_cb_detach(const URequest *request,
                          UResponse *response,
                          void *data);

int web_service_cb_get_ue(const URequest *request,
                          UResponse *response,
                          void *data);

int web_service_cb_list_ues(const URequest *request,
                            UResponse *response,
                            void *data);

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *data);

#endif /* WEB_SERVICE_H_ */
