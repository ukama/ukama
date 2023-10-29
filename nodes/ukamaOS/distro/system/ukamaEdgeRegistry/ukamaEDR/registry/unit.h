/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
