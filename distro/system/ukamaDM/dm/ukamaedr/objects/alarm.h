/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef OBJ_ALARM_H_
#define OBJ_ALARM_H_

#include "objects/objects.h"

/* Event Type */
#define EVENT_TYPE_DISABLED 			0x00
#define EVENT_TYPE_ALARM_CURR_STATE 	0x01
#define EVENT_TYPE_STATE_CHANGE			0x02
#define	EVENT_TYPE_LOG					0x03

#define OBJECT_ID_DEV_ALARM		  		34570

/* Real TIME */
#define EVENT_REALTIME					true
#define EVENT_ONREQ						false

typedef struct __attribute__((__packed__)) {
	uint16_t		instanceId;
	int    			eventtype;
	bool            realtime;
	uint16_t        state;
	char			disc[MAX_LWM2M_OBJ_STR_LEN];
	double          lowthreshold;
	double          highthreshold;
	double          crithreshold;
	int				eventcount;
	int64_t 		time;
	uint16_t		sobjid;
	uint16_t		sinstid;
	uint16_t		srsrcid;
	double			sensorvalue;
	char            sensorunits[MAX_LWM2M_OBJ_STR_LEN];
	char            applicationtype[MAX_LWM2M_OBJ_STR_LEN];
}AlarmObjInfo;

typedef struct _alarm_info
{
    struct AlarmObjInfo * next;         // matches lwm2m_list_t::next
    AlarmObjInfo data;
} alarm_info_t;

// Resource Id's:
#define RES_M_AL_EVENTTYPE         	6011
#define RES_M_AL_REALTIME         	6012
#define RES_M_AL_STATE              1
#define RES_M_AL_LOW_LIMIT      	2 // Low limit alarm
#define RES_M_AL_HIGH_LIMIT		    3
#define RES_M_AL_CRIT_LIMIT		    4
#define RES_M_AL_LOW_THRESHOLD      5 // Low limit
#define RES_M_AL_HIGH_THRESHOLD     6
#define RES_M_AL_CRIT_THRESHOLD     7
#define RES_M_AL_EVT_COUNT          6018
#define RES_M_AL_RECRD_TIME		    6021
#define RES_M_AL_CLEAR		        6022
#define RES_M_AL_OBJ_ID		        8
#define RES_M_AL_INST_ID            9
#define RES_M_AL_RSRC_ID		    10
#define RES_M_AL_ENABLE		        11 /* TODO: Enable alarms */
#define RES_M_AL_DISABLE	        12 /* TODO: Disable alarms */
#define RES_O_AL_DESCRIPTION        13
#define RES_M_SENSOR_VALUE          5700
#define RES_M_SENSOR_UNITS          5701
#define RES_O_APPLICATION_TYPE      5750


#endif /* OBJ_ALARM_H_ */
