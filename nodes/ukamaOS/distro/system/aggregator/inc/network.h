/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AGG_NETWORK_H_
#define AGG_NETWORK_H_

#include "aggregator.h"
#include "ulfius.h"

typedef struct _u_instance UInst;
typedef struct _u_request URequest;
typedef struct _u_response UResponse;

int start_metrics_web_service(UInst *metricsInst, AppState *state);
int start_admin_web_service(UInst *adminInst, AppState *state);
void stop_web_service(UInst *inst);

#endif /* AGG_NETWORK_H_ */
