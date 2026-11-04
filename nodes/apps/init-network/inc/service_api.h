/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SERVICE_API_H_
#define SERVICE_API_H_

#include "http_status.h"

#define HTTP_DELETE "DELETE"
#define HTTP_GET    "GET"
#define HTTP_POST   "POST"
#define HTTP_PUT    "PUT"

#define API_EP_PING      API_RES_EP("ping")
#define API_EP_READY     API_RES_EP("ready")
#define API_EP_RECONCILE API_RES_EP("reconcile")
#define API_EP_SERVICE   API_RES_EP("service")
#define API_EP_STATUS    API_RES_EP("status")
#define API_EP_VERSION   API_RES_EP("version")

#endif /* SERVICE_API_H_ */
