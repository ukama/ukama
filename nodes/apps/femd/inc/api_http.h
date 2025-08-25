/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */
#ifndef API_HTTP_H
#define API_HTTP_H

#include <ulfius.h>

#include "femd.h"
#include "gpio_controller.h"

/* Set standard CORS/Allow headers for a response. allow_str examples:
   "GET, OPTIONS" or "GET, PUT, PATCH, OPTIONS" */
void api_set_cors_allow(UResponse *resp, const char *allow_str);
int api_parse_fem_id(const URequest *req, FemUnit *out_unit);
int api_parse_channel_id(const URequest *req, int *out_ch);

#endif /* API_HTTP_H */
