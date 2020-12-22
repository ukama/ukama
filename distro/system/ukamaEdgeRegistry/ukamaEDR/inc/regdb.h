/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
