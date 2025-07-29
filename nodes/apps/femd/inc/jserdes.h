/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef JSERDES_H
#define JSERDES_H

#include <jansson.h>

#define JSON_OK                        STATUS_OK
#define JSON_FAILURE                   STATUS_NOTOK
#define JSON_ENCODING_OK               JSON_OK
#define JSON_DECODING_OK               JSON_OK

void json_log(json_t *json);
void json_free(json_t** json);

#endif /* JSERDES_H_ */

