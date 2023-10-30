/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_DRDB_H_
#define INC_DRDB_H_

#include "headers/objects/objects.h"
#include "headers/ubsp/devices.h"
#include "headers/ubsp/property.h"

#include <stdbool.h>
#include <stdint.h>
#include <string.h>

typedef int (*ExecFunc)(void* data);

typedef union  {
	bool boolval;
	uint16_t sintval;
	int intval;
	int64_t lintval;
	char stringval[32];
	double doubleval;
} PValue;

/* It would have been really nice if we had it as union but no way to figure out
 * which data is present. So for now keeping it as struct. */
typedef struct {
	/* If property exist than read item from property otherwise value
	 *  contains the required item */
	int resourceId;
	Property* prop;
	union {
		PValue value;
		ExecFunc execFunc;
	};
} PData;

typedef struct {
	uint16_t instance;
	ObjectType type;
	char UUID[32]; /* Only Valid for unit and Module registry */
	DevObj obj;
	const void *dbfxntbl;
	void *data;
} DRDBSchema;

typedef struct {
	uint16_t instance;
	char UUID[32];
	const void *dbfxntbl;
	void *data;
} ModDBSchema;

typedef struct {
	uint16_t instance;
	char UUID[32];
	const void *dbfxntbl;
	void *data;
} UnitDBSchema;

#endif /* INC_DRDB_H_ */
