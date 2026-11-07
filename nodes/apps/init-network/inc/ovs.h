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
#include <stddef.h>

#include "config.h"
#include "status.h"

bool ovs_setup(Config *config, AppStatus *status);
bool ovs_reconcile(Config *config, AppStatus *status);
bool ovs_check(Config *config, char *reason, size_t reasonSize);

bool ovs_service_set(Config *config, bool enabled);
JsonObj *ovs_service_status(Config *config);
void ovs_service_poll(Config *config);

#endif /* OVS_H_ */
