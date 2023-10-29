/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef TMP_H_
#define TMP_H_

#include "inc/devicefxn.h"
#include "headers/utils/list.h"
#include "headers/ubsp/devices.h"

#define MAX_TEMP_SENSOR_TYPE 			3
#define TEMP_SNSR_SE98ATP			0x01
#define TEMP_SNSR_ADT				0x02
#define TEMP_SNSR_TMP464			0x03

const DevFxnTable* get_dev_tmp_fxn_tbl(char *name);
ListInfo* get_dev_tmp_db();
#endif /* TMP_H_ */
