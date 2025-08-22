/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#ifndef WEB_CLIENT_H_
#define WEB_CLIENT_H_

#include "ulfius.h"
#include "config.h"

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;

int get_nodeid_and_type_from_noded(Config *config);
int wc_send_alarm_to_notifyd(Config *config, int *retCode, int type);

#endif /* WEB_CLIENT_H_ */
