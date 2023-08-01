/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef AGENT_CALLBACK_H
#define AGENT_CALLBACK_H

int agent_web_service_cb_default(const URequest *request,
                                 UResponse *response,
                                 void *data);
int agent_web_service_cb_post_capp(const URequest *request,
                                   UResponse *response,
                                   void *data);
int agent_web_service_cb_ping(const URequest *request,
                              UResponse *response,
                              void *data);
int agent_web_service_cb_default(const URequest *request,
                                 UResponse *response,
                                 void *data);

#endif /* AGENT_CALLBACK_H */
