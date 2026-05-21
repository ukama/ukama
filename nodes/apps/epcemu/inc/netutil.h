/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef NETUTIL_H_
#define NETUTIL_H_

#include <stdint.h>
#include <stdbool.h>

#include "jansson.h"

typedef json_t JsonObj;

int imsi_valid(const char *imsi);
int ip_to_uint32(const char *ip, uint32_t *out);
int ip_in_cidr(const char *ip, const char *cidr);
JsonObj *imsi_to_json_array(const char *imsi);

#endif /* NETUTIL_H_ */

