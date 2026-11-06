/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SERVICE_GROUP_H_
#define SERVICE_GROUP_H_

#include "config.h"
#include "deviced.h"

bool service_group_start(Config *config);
void service_group_stop(Config *config);
JsonObj *service_group_state(Config *config);
int service_group_command(const URequest *request, UResponse *response,
                          void *data);
int service_group_watchdog(const URequest *request, UResponse *response,
                           void *data);
int service_group_maintenance(const URequest *request, UResponse *response,
                              void *data);

#endif /* SERVICE_GROUP_H_ */
