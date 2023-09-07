/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
