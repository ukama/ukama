/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INIT_CLIENT_JSERDES_H
#define INIT_CLIENT_JSERDES_H

#include <jansson.h>

#include "initClient.h"

#define JSON_IP          "ip"
#define JSON_PORT        "port"
#define JSON_CERTIFICATE "certificate"

int serialize_request(Request *request, json_t **json);

#endif /* INIT_CLIENT_JSERDES_H */
