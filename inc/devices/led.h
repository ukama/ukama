/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DEVDB_LED_LED_H_
#define DEVDB__LED_LED_H_

#include "usys_list.h"

#define MAX_LED_SENSOR_TYPE		1

const DevOps* get_led_dev_ops(char *name);
ListInfo* get_led_dev_ledgr();

#endif /* DEVDB_LED_LED_H_ */
