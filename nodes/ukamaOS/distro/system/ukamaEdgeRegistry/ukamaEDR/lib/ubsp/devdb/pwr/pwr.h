/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DEVDB_PWR_PWR_H_
#define DEVDB_PWR_PWR_H_

#include "inc/devicefxn.h"
#include "headers/utils/list.h"
#include "headers/ubsp/devices.h"

#define MAX_PWR_SENSOR_TYPE		1

const DevFxnTable* get_dev_pwr_fxn_tbl(char *name);
ListInfo* get_dev_pwr_db();

#endif /* DEVDB_PWR_PWR_H_ */
