/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef REGISTRY_ALARM_H_
#define REGISTRY_ALARM_H_

#include "inc/registry.h"
#include "headers/edr/ifmsg.h"

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
