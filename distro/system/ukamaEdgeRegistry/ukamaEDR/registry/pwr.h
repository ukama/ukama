/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef EDRDB_PWR_H_
#define EDRDB_PWR_H_

#include "inc/registry.h"
#include "headers/edr/ifmsg.h"
#include "headers/utils/list.h"

typedef struct {
	PData value;
	PData min;
	PData max;
	PData avg;
	PData cumm;
	PData counter;
	PData minrange;
	PData maxrange;
	PData units;
	PData calibration;
	PData applicationtype;
	PData resetcounter;
} GenPwrData;

typedef GenPwrData VoltData;
typedef GenPwrData CurrData;
typedef GenPwrData PwrData;

typedef struct {
	uint16_t instance;
	DevObj obj;
	GenPwrData data;
} DRDBGenPwrSchema;

int drdb_exec_pwr_inst_rsrc(DRDBSchema* reg, MsgFrame* rqmsg);
int drdb_read_pwr_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);
int drdb_write_pwr_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);

void free_pwr_data (void* data);
void drdb_add_pwr_dev_to_reg(void* pdev);
void drdb_add_curr_inst_to_reg(Device *dev, Property* prop, uint8_t inst, int pidx);
void drdb_add_pwr_inst_to_reg(Device *dev, Property* prop, uint8_t inst, int pidx);
void drdb_add_volt_inst_to_reg(Device *dev, Property* prop, uint8_t inst, int pidx);
void drdb_update_pwr_inst_data(double val, PData* min, PData* max, PData* avg, PData* cumm, PData* count);

void* copy_pwr_data(void *pdata);

#endif /* EDRDB_PWR_H_ */
