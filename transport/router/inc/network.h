/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef NETWORK_H
#define NETWORK_H

#include "router.h"

int init_frameworks(Config *config, struct _u_instance *webInst);
void setup_endpoints(Router *router, struct _u_instance *instance);
int start_framework(struct _u_instance *instance);
int start_web_service(Router *router, struct _u_instance *webInst);

#endif /* WIMC_NETWORK_H */
