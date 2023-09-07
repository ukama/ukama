/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef OBJ_DIGITAL_OUTPUT_H_
#define OBJ_DIGITAL_OUTPUT_H_

#include "objects/objects.h"

#define OBJECT_ID_DIGITAL_OUTPUT		   3201

typedef struct __attribute__((__packed__)) {
    uint16_t                instanceId;
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
} DopObjInfo;

typedef struct _digital_output_t
{
    struct _temp_info       *next;
    DopObjInfo					data;
} digital_output_t;

// Resource Id's:
#define RES_M_DIGITAL_OUTPUT_STATE                      5550
#define RES_O_DIGITAL_OUTPUT_POLARITY                   5551
#define RES_O_APPLICATION_TYPE                          5750

#endif
