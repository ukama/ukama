/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
