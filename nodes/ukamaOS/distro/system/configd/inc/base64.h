/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef BASE64_H_
#define BASE64_H_

#include "usys_types.h"

int base64_decode(char *decodedData, const char *encodedData);

#endif
