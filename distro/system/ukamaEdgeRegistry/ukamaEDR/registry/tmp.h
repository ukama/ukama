/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef EDRDB_TMP_H_
#define EDRDB_TMP_H_

#include "inc/registry.h"
#include "ifmsg.h"

#define MAX_TMP_RANGE	(125000)
#define MIN_TMP_RANGE	(-40000)


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
	PData applicationtype;
	PData resetcounter;
} TempData;

typedef struct {
	uint8_t instance;
	DevObj obj;
	TempData data;
} DRDBTempSchema;

int drdb_add_adt7481_dev_to_reg(Device* dev, Property* prop);
int drdb_add_se98_dev_to_reg(Device* dev, Property* prop);
int drdb_add_tmp464_dev_to_reg(Device* dev, Property* prop);
int drdb_add_tmp_inst_to_reg(Device *dev, Property *prop, uint8_t instance,
		uint8_t subdev);
int drdb_exec_tmp_inst_rsrc(DRDBSchema* reg, MsgFrame* rqmsg);
int drdb_read_tmp_inst_data_from_dev(DRDBSchema* reg, MsgFrame* temp);
int drdb_write_tmp_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);
void free_tmp_data (void* data);
void drdb_add_tmp_dev_to_reg(void* pdev);
void drdb_tmp_alarm_cb(DevObj *obj, void *data, void *prop_idx,
		void *count);
void drdb_update_tmp_inst_data(double temp, PData* min, PData* max,
		PData* avg, PData* cumm, PData* count);
void* copy_tmp_data(void *pdata);
#endif /* EDRDB_TMP_H_ */
