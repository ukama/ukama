/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef OPS_H_
#define OPS_H_

#include "aisgd.h"

bool aisgd_ops_reconcile(AisgdContext *ctx, JsonObj **response);
bool aisgd_ops_refresh_status(AisgdContext *ctx);
bool aisgd_ops_scan(AisgdContext *ctx, JsonObj **response);
bool aisgd_ops_get_device(AisgdContext *ctx, JsonObj **response);
bool aisgd_ops_get_info(AisgdContext *ctx, JsonObj **response);
bool aisgd_ops_get_alarms(AisgdContext *ctx, JsonObj **response);
bool aisgd_ops_clear_alarms(AisgdContext *ctx, JsonObj **response);
bool aisgd_ops_subscribe_alarms(AisgdContext *ctx, JsonObj **response);
bool aisgd_ops_self_test(AisgdContext *ctx, JsonObj **response);
bool aisgd_ops_configure(AisgdContext *ctx,
                          const char *profile,
                          const char *configPath,
                          JsonObj **response);
bool aisgd_ops_calibrate(AisgdContext *ctx, JsonObj **response);
bool aisgd_ops_get_tilt(AisgdContext *ctx, JsonObj **response);
bool aisgd_ops_set_tilt(AisgdContext *ctx,
                         double targetTiltDeg,
                         JsonObj **response);
bool aisgd_ops_get_device_data(AisgdContext *ctx,
                                int field,
                                JsonObj **response);
bool aisgd_ops_reset(AisgdContext *ctx, JsonObj **response);

#endif /* OPS_H_ */
