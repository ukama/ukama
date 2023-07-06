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

#include "deviced.h"

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig);

int web_service_cb_post_reboot(const URequest *request,
                               UResponse *response,
                               void *epConfig);

int web_service_cb_post_restart(const URequest *request,
                                UResponse *response,
                                void *epConfig);

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig);

#endif /* WEB_SERVICE_H_ */
