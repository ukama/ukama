/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef OVS_H_
#define OVS_H_

#include <stdbool.h>

#include "config.h"
#include "status.h"

bool ovs_setup(Config *config, AppStatus *status);
bool ovs_reconcile(Config *config, AppStatus *status);

#endif /* OVS_H_ */
