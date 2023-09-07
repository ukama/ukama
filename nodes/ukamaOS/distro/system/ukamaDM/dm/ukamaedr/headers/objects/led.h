/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */


#ifndef OBJ_LED_H_
#define OBJ_LED_H_

#include "headers/objects/objects.h"

#define OBJECT_ID_LED		   3311

typedef struct __attribute__((__packed__)) {
    uint16_t    instanceId;
    bool        onoff;
    int      	dimmer;
    int      	ontime;
    double      cumm_active_pwr;
    double      pwr_factor;
    char      	colour[MAX_LWM2M_OBJ_STR_LEN];
    char        sensor_units[MAX_LWM2M_OBJ_STR_LEN];
    char        application_type[MAX_LWM2M_OBJ_STR_LEN];
} LedObjInfo;

typedef struct _led_info
{
    struct _led_info * next;
    LedObjInfo			   data;
} led_info_t;

// Resource Id's:
#define RES_M_ONOFF_VALUE                       5850
#define RES_M_DIMMER_VALUE                		5851
#define RES_M_ONTIME_VALUE            			5852
#define RES_O_CUMM_ACTIVE_PWR_VALUE             5805
#define RES_O_PWR_FACTOR_VALUE                  5820
#define RES_M_COLOUR_VALUE                      5606
#define RES_M_SENSOR_UNITS                      5701
#define RES_O_APPLICATION_TYPE                  5750

#endif /* OBJ_LED_H_ */
