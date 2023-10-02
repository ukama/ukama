/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WEB_SERVICE_H_
#define WEB_SERVICE_H_

#include "starter.h"

int web_service_cb_ping(const URequest *request,
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
