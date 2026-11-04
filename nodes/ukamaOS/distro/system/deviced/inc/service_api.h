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

#define HTTP_GET  "GET"
#define HTTP_POST "POST"

#define API_EP_RADIO               API_RES_EP("radio")
#define API_EP_REBOOT              API_RES_EP("reboot")
#define API_EP_SERVICE             API_RES_EP("service")
#define API_EP_SERVICE_MAINTENANCE API_RES_EP("service/maintenance")
#define API_EP_SERVICE_WATCHDOG    API_RES_EP("service/watchdog")
#define API_EP_STATE               API_RES_EP("state")

#endif /* SERVICE_API_H_ */
