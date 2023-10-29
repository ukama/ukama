/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef REGISTRY_ALARM_H_
#define REGISTRY_ALARM_H_

#include "ifmsg.h"
#include "inc/registry.h"

#define PROPERY_IDX_NOT_APPLICABLE		(-1)

typedef struct {
	uint16_t		sobjid;
	uint16_t		sinstid;
	uint16_t		srsrcid;
} DevIDT;

/* Storing property index for threshold.
 * Used for initializing AlarmPropertydata
 * if value is negative property idx is not available
 */
typedef struct {
	int plowthresholdidx;
	int phighthresholdidx;
	int pcrithresholdidx;
	int psensorvalueidx;
	int plowlimitalarmidx;
	int phighlimitalarmidx;
	int pcrilimitalarmidx;
} AlarmSensorData;

/* Hold the device properties for thresholds */
typedef struct {
	Property		*plowthreshold;
	Property		*phighthreshold;
	Property		*pcritthreshold;
	Property		*plowlimitalarm;
	Property		*phighlimitalarm;
	Property		*pcrilimitalarm;
	Property		*psensorvalue;
} AlarmPropertyData;

typedef struct {
	PData eventtype;
	PData realtime;
	PData state;
	PData disc;
	PData lowthreshold;
	PData highthreshold;
	PData crithreshold;
	PData eventcount;
	PData time;
	PData clear;
	PData enbalarm;
	PData disalarm;
	PData sobjid;
	PData sinstid;
	PData srsrcid;
	PData sensorvalue;
	PData sensorunits;
	PData applicationtype;
} AlarmData;

typedef struct  {
	uint8_t instance;
	DevObj obj;
	AlarmData data;
} DRDBAlarmSchema;

int drdb_exec_alarm_inst_rsrc(DRDBSchema *reg, MsgFrame *rqmsg);
int drdb_read_alarm_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg);
int drdb_write_alarm_inst_data_to_dev(DRDBSchema *reg, MsgFrame *rqmsg);
void drdb_alarm_cb(DevObj *obj, AlertCallBackData** acbdata, int* count);
void drdb_add_alarm_inst_to_reg(DevObj* obj, DevIDT* idt, AlarmPropertyData* pdata);
void free_alarm_data(void *pdata);
void* copy_alarm_data(void *pdata);
#endif /* REGISTRY_ALARM_H_ */
