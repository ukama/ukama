/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
