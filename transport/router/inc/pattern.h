/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef PATTERN_H
#define PATTERN_H

#include "router.h"

#define ASTERIK_ONLY "*"

void free_service(Service *service);
int find_matching_service(Router *router, Pattern *requestPattern,
			  Forward **forward);

#endif /* PATTERN_H */
