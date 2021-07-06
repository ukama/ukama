/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WIMC_PROVIDER_H
#define WIMC_PROVIDER_H

#include <stdio.h>
#include <string.h>

extern int deserialize_provider_response(json_t *resp, AgentCB **agent,
					 int *counter);
extern req_t* create_http_request(char *url, char* ep, char *req_type);

#endif /* WIMC_PROVIDER_H */
