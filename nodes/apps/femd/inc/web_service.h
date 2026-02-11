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

#include "femd.h"
#include "config.h"
#include "jobs.h"
#include "snapshot.h"

#define URL_PREFIX "/v1"

typedef struct {
    Config        *config;
    Jobs          *jobs;
    SnapshotStore *snap;
} WebCtx;

int start_web_service(ServerConfig *serverConfig, UInst *serviceInst, WebCtx *ctx);

#endif /* WEB_SERVICE_H */
