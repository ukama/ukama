/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef OBJ_DIGITAL_INPUT_H_
#define OBJ_DIGITAL_INPUT_H_

#include "headers/objects/objects.h"

#define OBJECT_ID_DIGITAL_INPUT		   3200

typedef struct __attribute__((__packed__)) {
	uint16_t    			instanceId;             // matches lwm2m_list_t::id
	bool                    digital_state;
	int                     direction;
	int                     digital_counter;
	bool                    digital_polarity;
	int                     digital_debounce;
	int                     digitial_edge_selection;
	int 					ontime;
	int 					offtime;
	char                    application_type[MAX_LWM2M_OBJ_STR_LEN];
	char                    sensor_type[MAX_LWM2M_OBJ_STR_LEN];
} DipObjInfo;

typedef DipObjInfo DigObjInfo;

typedef struct _digital_input_t
{
    struct _temp_info       *next;
    DipObjInfo                data;
} digital_input_t;

// Resource Id's:
#define RES_M_DIGITAL_INPUT_STATE                      5500
#define RES_O_DIGITAL_INPUT_COUNTER                    5501
#define RES_O_DIGITAL_INPUT_POLARITY                   5502
#define RES_O_DIGITAL_INPUT_DEBOUNCE                   5503
#define RES_O_DIGITIAL_INPUT_EDGE_SELECTION            5504
#define RES_O_APPLICATION_TYPE                         5750
#define RES_O_SENSOR_TYPE                              5751
#define RES_O_DIGITAL_INPUT_COUNTER_RESET              5505

#endif
