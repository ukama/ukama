/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_ALARMHANDLER_H_
#define INC_ALARMHANDLER_H_

#include <stdbool.h>
#include <stdint.h>
#include <string.h>

#include "registry/alarm.h"
#include "headers/objects/alarm.h"
#include "headers/ubsp/devices.h"

#define ALARM_REPORTED			true
#define ALARM_NOT_REPORTED		false

#define ALARM_RESPONSE_RECVD	true
#define ALARM_RESPONSE_PENDING	false

typedef struct {
	int token;
	DevObj obj;
	uint16_t objid;
	uint16_t instid;
	uint16_t rsrcid;
	bool reported; /* Msg sent to server */
	bool response; /* Response msg received from server */
	uint16_t txcount; /* Send to server */
	AlarmObjInfo data;
} AlarmSchema;

int alarmhandler_proc_node(void* node);
int alramdb_record_data(AlarmSchema* node);
int alarmdb_update(AlarmSchema* aschema);
int alramhandler_prepare_to_tx(AlarmSchema* anode, bool newAlarm);
int alarmhandler_proc_alarm_resp(AlarmSchema* anode, MsgFrame* rmsg, MsgFrame* smsg, int flag);
int alarmhandler_verify_resp(MsgFrame* rmsg, MsgFrame* smsg, void* data, int sflag);
void alarmdb_exit();
void alarmdb_init();
void alarmhandler_exit();
void alarmhandler_init();
void alarmhandler_start();
void alarmhandler_stop(size_t timer);
void alarmdb_pop(AlarmSchema* aschema);
void alarmdb_push(DRDBSchema *node);
void alarmhandler_service(size_t timer_id, void *data);
AlarmSchema* alarmdb_search_node(AlarmSchema* snode);
AlarmObjInfo* alarm_drdbschema_to_objectdb(AlarmData *ndata);
AlarmSchema* alarmdb_create_node(DevObj* obj, uint16_t instid, uint16_t objid,
		uint16_t rsrcid);
MsgFrame* alarmhandler_create_tx_frame(AlarmSchema* anode, bool newtoken);
#endif /* INC_ALARMHANDLER_H_ */
