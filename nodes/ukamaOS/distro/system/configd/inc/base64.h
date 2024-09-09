/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef BASE64_H_
#define BASE64_H_

#include "usys_types.h"

int base64_decode(char *decodedData, const char *encodedData);

#endif
