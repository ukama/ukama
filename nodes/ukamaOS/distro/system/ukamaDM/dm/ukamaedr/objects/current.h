/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef OBJ_CURRENT_H_
#define OBJ_CURRENT_H_

#include "objects/objects.h"

#define OBJECT_ID_CURR		   3317

typedef struct __attribute__((__packed__)) {
	uint16_t    instanceId;             // matches lwm2m_list_t::id
	char        sensor_units[MAX_LWM2M_OBJ_STR_LEN];
	double      sensor_value;
	double      min_measured_value;
	double      max_measured_value;
	double      avg_measured_value;
	double      min_range_value;
	double      max_range_value;
	double      calibration_value;
	char        application_type[MAX_LWM2M_OBJ_STR_LEN];
} CurrObjInfo;

typedef struct _curr_info
{
    struct _curr_info * next;
    CurrObjInfo data;
} curr_info_t;

// Resource Id's:
#define RES_M_SENSOR_VALUE                      5700
#define RES_O_MIN_MEASURED_VALUE                5601
#define RES_O_MAX_MEASURED_VALUE                5602
#define RES_O_MIN_RANGE_VALUE                   5603
#define RES_O_MAX_RANGE_VALUE                   5604
#define RES_M_SENSOR_UNITS                      5701
#define RES_O_RESET_MIN_AND_MAX_MEASURED_VALUE  5605
#define RES_O_CURR_CALIBRATION_VALUE			5821
#define RES_O_APPLICATION_TYPE                  5750

#endif /* OBJ_CURRENT_H_ */
