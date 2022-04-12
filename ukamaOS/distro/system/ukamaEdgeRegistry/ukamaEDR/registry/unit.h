/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef REGISTRY_UNIT_H_
#define REGISTRY_UNIT_H_

#include "inc/registry.h"
#include "ifmsg.h"
#include "headers/ubsp/ukdblayout.h"

typedef struct {
	PData UUID;
	PData name;
	PData unit;
	PData asmdate;
	PData oemname;
	PData skew;
	PData mac;
	PData swversion;
	PData pswversion;
	PData modcount;
} UnitData;

typedef struct {
	uint8_t instance;
	char unitUUID[32];
	const void *dbfxntbl;
	UnitData data;
} DRDBUnitSchema;

int drdb_exec_unit_inst_rsrc(DRDBSchema* reg, MsgFrame* rqmsg);
int drdb_read_unit_inst_data_from_dev(DRDBSchema* reg, MsgFrame* temp);
int drdb_write_unit_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);

void drdb_add_unit_to_reg(void* pdev);
void drdb_add_unit_inst_to_reg(UnitInfo* uinfo, uint8_t instance, uint8_t subdev);
void free_unit_data (void* data);
void* copy_unit_data(void *pdata);

#endif /* REGISTRY_UNIT_H_ */
