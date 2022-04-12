/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef REGISTRY_MODULE_H_
#define REGISTRY_MODULE_H_

#include "inc/registry.h"
#include "ifmsg.h"
#include "headers/ubsp/ukdblayout.h"

typedef struct {
	PData UUID;
	PData name;
	PData moduleclass; /* Could be class */
	PData partnumber;
	PData mfgdate;
	PData mfgname;
	PData hwversion;
	PData mac;
	PData swversion;
	PData pswversion;
	PData devcount;
} ModuleData;

typedef struct {
	uint8_t instance;
	char modUUID[32];
	const void *dbfxntbl;
	ModuleData data;
} DRDBModSchema;

int drdb_exec_mod_inst_rsrc(DRDBSchema* reg, MsgFrame* rqmsg);
int drdb_read_mod_inst_data_from_dev(DRDBSchema* reg, MsgFrame* temp);
int drdb_write_mod_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);

void drdb_add_mod_to_reg(void* pdev);
void drdb_add_mod_inst_to_reg(ModuleInfo* minfo, uint8_t instance, uint8_t subdev);
void free_mod_data (void* data);
void* copy_module_data(void *pdata);
#endif /* REGISTRY_MODULE_H_ */
