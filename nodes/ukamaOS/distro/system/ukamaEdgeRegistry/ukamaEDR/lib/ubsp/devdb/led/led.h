/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "inc/devicefxn.h"
#include "headers/utils/list.h"
#include "headers/ubsp/devices.h"

#ifndef DEVDB_LED_LED_H_
#define DEVDB__LED_LED_H_

#define MAX_LED_SENSOR_TYPE		1

const DevFxnTable* get_dev_led_fxn_tbl(char *name);
ListInfo* get_dev_led_db();

#endif /* DEVDB_LED_LED_H_ */
