/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef LIFECYCLE_GUARD_H_
#define LIFECYCLE_GUARD_H_

#include "app.h"

bool lifecycle_guard_begin(App *app);
void lifecycle_guard_end(App *app);

#endif /* LIFECYCLE_GUARD_H_ */
