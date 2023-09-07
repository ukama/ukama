/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_DRDB_H_
#define INC_DRDB_H_

#include "headers/devices.h"
#include "headers/ifmsg.h"
#include "headers/property.h"

#include <stdbool.h>
#include <stdint.h>
#include <string.h>

typedef int (*ExecFunc)(void* data);

typedef union  {
	bool boolval;
	int intval;
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


int reg_add_dev(Device* dev);
int reg_exec_dev(MsgFrame *req);
int reg_read_dev(MsgFrame *req);
int reg_read_max_instance(DeviceType type);
int reg_read_inst_data_from_dev(DevObj* obj, PData* sd);
int reg_write_dev(MsgFrame *req);
int reg_write_inst_data_to_dev(DevObj* obj, PData* sd);
int reg_register_devices();
int reg_register_misc();
int reg_search_inst(int instance, uint16_t misc, DeviceType type, DRDBSchema *out);
int reg_update_dev(DRDBSchema *reg);

void free_data(DRDBSchema* reg);
void free_sdata(Property *prop);
void reg_exit();
void reg_exit_type(DeviceType type);
void reg_init();
void reg_init_type(DeviceType type);

void reg_mod_init();
void reg_mod_exit();
void reg_unit_init();
void reg_unit_exit();
void* reg_data_value(PData* sd);

DRDBSchema *reg_read_instance(int instance, DeviceType type);
ListInfo* reg_getdb(uint16_t misc, DeviceType type);
ListInfo *reg_getdevdb(DeviceType type);
ListInfo *reg_getmoddb();
ListInfo *reg_getunitdb();
#endif /* INC_DRDB_H_ */
