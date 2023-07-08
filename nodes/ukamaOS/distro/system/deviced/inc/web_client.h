/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
int wc_send_alarm_to_notifyd(Config *config);

#endif /* WEB_CLIENT_H_ */
