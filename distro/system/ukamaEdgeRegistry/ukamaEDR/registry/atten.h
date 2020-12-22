/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */


#ifndef EDRDB_ATTEN_H_
#define EDRDB_ATTEN_H_

#include "inc/registry.h"
#include "headers/edr/ifmsg.h"

typedef struct {
	PData attvalue;
	PData minrange;
	PData maxrange;
	PData latchenable; /* TODO: check if required or not */
	PData units;
	PData applicationtype;
}AttData;

typedef struct {
	uint8_t instance;
	DevObj obj;
	AttData data;
} DRDBAttSchema;

int drdb_read_att_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);
int drdb_write_att_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);

void drdb_add_att_dev_to_reg(void* pdev);
void drdb_add_att_inst_to_reg(Device* dev, uint8_t inst, uint8_t subdev);
void free_att_data (void* data);
Property *get_att_dev_property(DevObj* obj, int pid);
void* copy_att_data(void *pdata);
#endif /* EDRDB_ATTEN_H_ */
