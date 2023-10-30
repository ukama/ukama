/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_EDRDB_H_
#define INC_EDRDB_H_

#include "inc/registry.h"
#include "ifmsg.h"


typedef void (*DBInit)();
typedef void (*DBExit)();
typedef void (*DBAddDevToReg)(void* dev);
typedef int (*DBReadDataFromDev)(DRDBSchema* reg, MsgFrame* rqmsg);
typedef int (*DBWriteDataFromDev)(DRDBSchema* reg, MsgFrame* rqmsg);
typedef int (*DBExecDataFromDev)(DRDBSchema* reg, MsgFrame* rqmsg);
typedef int (*DBSearchInstInReg)(int instance, DeviceType type, DRDBSchema *out);
typedef void (*DBUpdateInstInReg)(double data, PData* min, PData* max, PData* avg, PData* cumm, PData* count);
typedef void (*DBFreeInstDataFromReg)(void* data);

typedef struct {
	DBAddDevToReg db_add_dev_to_reg;
	DBReadDataFromDev db_read_data_from_dev;
	DBWriteDataFromDev db_write_data_from_dev;
	DBExecDataFromDev db_exec;
	DBSearchInstInReg db_search_inst_in_reg;
	DBFreeInstDataFromReg db_free_inst_data_from_reg;
	DBUpdateInstInReg db_update_inst_in_reg;
}DBFxnTable;

#endif /* INC_EDRDB_H_ */
