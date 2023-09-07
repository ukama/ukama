/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef IDB_MFGDATA_H_
#define IDB_MFGDATA_H_

#include "headers/ubsp/ukdblayout.h"

typedef struct __attribute__((__packed__)) {
	char uuid[24];
	uint8_t idx_count;
	UKDB ukdb;
} MFGData;

#endif /* IDB_MFGDATA_H_ */
