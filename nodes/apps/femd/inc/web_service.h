/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#ifndef WEB_SERVICE_H
#define WEB_SERVICE_H

#include <ulfius.h>

#include "config.h"
#include "jobs.h"
#include "snapshot.h"
#include "gpio_controller.h"
#include "safety.h"
#include "notifier.h"

#ifndef URL_PREFIX
#define URL_PREFIX "/v1"
#endif

#ifndef API_RES_EP
#define API_RES_EP(x) "/" x
#endif

typedef struct {
    Jobs          *jobs;
    SnapshotStore *snap;
    GpioController*gpio;
    Safety        *safety;
    Notifier      *notifier;
    Config        *cfg;
} WebCtx;

int start_web_service(ServerConfig *serverConfig, UInst *serviceInst, WebCtx *ctx);

#endif /* WEB_SERVICE_H */
