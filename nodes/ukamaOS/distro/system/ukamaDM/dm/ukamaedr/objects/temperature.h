/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef OBJ_TEMP_H_
#define OBJ_TEMP_H_

#include "objects/objects.h"

#define OBJECT_ID_TMP				3303

typedef struct __attribute__((__packed__)) {
	uint16_t    instanceId;             // matches lwm2m_list_t::id
	char        sensor_units[MAX_LWM2M_OBJ_STR_LEN];
	double      sensor_value;
	double      min_measured_value;
	double      max_measured_value;
	double      avg_measured_value;
	double      min_range_value;
	double      max_range_value;
	char        application_type[MAX_LWM2M_OBJ_STR_LEN];
}TempObjInfo;

typedef struct _temp_info
{
    struct _temp_info * next;
    TempObjInfo	data;
} temp_info_t;

// Resource Id's:
#define RES_M_SENSOR_VALUE                      5700
#define RES_O_MIN_MEASURED_VALUE                5601
#define RES_O_MAX_MEASURED_VALUE                5602
#define RES_O_MIN_RANGE_VALUE                   5603
#define RES_O_MAX_RANGE_VALUE                   5604
#define RES_M_SENSOR_UNITS                      5701
#define RES_O_RESET_MIN_AND_MAX_MEASURED_VALUE  5605
#define RES_O_APPLICATION_TYPE                  5750
#endif
