/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
