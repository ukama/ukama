/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef FORWARD_H
#define FORWARD_H

#include "router.h"

req_t *create_forward_request(Forward *forward, Pattern *reqPattern,
			      const req_t *request);
int valid_forward_route(char *host, char *port);

#endif /* FORWARD_H */
