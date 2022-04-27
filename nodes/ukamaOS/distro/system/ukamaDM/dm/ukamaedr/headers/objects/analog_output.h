/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#ifndef OBJ_ANALOG_OUTPUT_H_
#define OBJ_ANALOG_OUTPUT_H_

#include "headers/objects/objects.h"

#define MINRANGE_ADCVALUE             	 0
#define MAXRANGE_ADCVALUE             	 12000

#define OBJECT_ID_ANALOG_OUTPUT		      3203

typedef struct __attribute__((__packed__)) {
	uint16_t    			instanceId;
	double                	outputcurr;             // matches lwm2m_list_t::id
	double                  minrange;
	double                  maxrange;
	char                    application_type[MAX_LWM2M_OBJ_STR_LEN];
}AdcObjInfo;

typedef struct _analog_output_info
{
    struct _analog_output_info * next;
    AdcObjInfo data;
} analog_output_info_t;

// Resource Id's:
#define RES_M_OUT_CURR_VALUE                    5650
#define RES_O_MIN_RANGE_VALUE                   5603
#define RES_O_MAX_RANGE_VALUE                   5604
#define RES_O_APPLICATION_TYPE                  5750

#endif /* OBJ_ANALOG_OUTPUT_H_ */
