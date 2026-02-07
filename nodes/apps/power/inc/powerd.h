/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef POWERD_H
#define POWERD_H

#include <stddef.h>
#include <time.h>

#include "usys_api.h"

#define SERVICE_NAME      "power"
#define DEF_LOG_LEVEL     "INFO"

#define URL_PREFIX        "/v1"
#define API_RES_EP(x)     "/" x

/* Common OK/NOK */
#ifndef STATUS_OK
#define STATUS_OK   0
#endif
#ifndef STATUS_NOK
#define STATUS_NOK  1
#endif

#endif
