/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
