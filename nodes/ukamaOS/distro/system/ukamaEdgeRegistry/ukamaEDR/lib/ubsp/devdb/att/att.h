/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DEVDB_ATT_ATT_H_
#define DEVDB_ATT_ATT_H_

#include "inc/devicefxn.h"
#include "headers/utils/list.h"
#include "headers/ubsp/devices.h"

#define MAX_ATT_SENSOR_TYPE		1

const DevFxnTable* get_dev_att_fxn_tbl(char *name);
ListInfo* get_dev_att_db();

#endif /* DEVDB_ATT_ATT_H_ */
