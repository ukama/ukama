/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef OBJ_ATTEN_H_
#define OBJ_ATTEN_H_

#include "objects/objects.h"

#define MINRANGE_ATTVALUE              0
#define MAXRANGE_ATTVALUE              126 /* 2*63 Each unit correspond to 0.5 dB change */

#define OBJECT_ID_ATTEN_OUTPUT		   34569

typedef struct __attribute__((__packed__)) {
	uint16_t    			instanceId;
	int                		attvalue;             // matches lwm2m_list_t::id
	bool                  	latchenable;
	int                  	minrange;
	int                  	maxrange;
	char                    sensor_units[MAX_LWM2M_OBJ_STR_LEN];
	char                    application_type[MAX_LWM2M_OBJ_STR_LEN];
}AttObjInfo;

typedef struct _atten_info
{
    struct _atten_info * next;         // matches lwm2m_list_t::next
    AttObjInfo data;
} atten_info_t;

// Resource Id's:
#define RES_M_ATTVALUE                      0
#define RES_M_MINRANGE          		    1
#define RES_M_MAXRANGE           	        2
#define RES_M_LATCH              			3
#define RES_M_SENSOR_UNITS                  5701
#define RES_O_APPLICATION_TYPE              5750

#endif
