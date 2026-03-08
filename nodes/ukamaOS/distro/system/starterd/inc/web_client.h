/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>

#include "ulfius.h"

#include "config.h"
#include "app.h"

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

bool wc_app_ping(Config *config, App *app);
bool wc_app_version_matches(Config *config,
                            App *app,
                            const char *tag);
bool wc_fetch_package(Config *config,
                      const char *appName,
                      const char *tag,
                      const char *dstPath);
