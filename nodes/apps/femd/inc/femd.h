/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#ifndef FEMD_H
#define FEMD_H

#include <stdint.h>

#include "ulfius.h"

#include "usys_types.h"
#include "usys_services.h"

#ifndef STATUS_OK
#define STATUS_OK  0
#endif

#ifndef STATUS_NOK
#define STATUS_NOK (-1)
#endif

#define SERVICE_NAME SERVICE_FEM
#define MODULE_FEM   "fem"
#define ALARM_NODE   "node"
#define EMPTY_STRING ""

#define ALARM_HIGH "HIGH"
#define ALARM_LOW  "LOW"

#define ALARM_PA_AUTO_OFF        "pa_auto_off"
#define ALARM_PA_AUTO_ON         "pa_auto_on"
#define ALARM_PA_AUTO_OFF_DESCRP "PA disabled automatically due to safety"
#define ALARM_PA_AUTO_ON_DESCRP  "PA restored automatically after safety cleared"

#define ALARM_TYPE_PA_OFF 1
#define ALARM_TYPE_PA_ON  2

#define DEF_NODED_HOST         "localhost"
#define DEF_NOTIFY_HOST        "localhost"
#define DEF_NODED_EP           "/v1/nodeinfo"
#define DEF_NOTIFY_EP          "/notify/v1/event/"
#define DEF_NODE_ID            "ukama-aaa-bbbb-ccc-dddd"
#define DEF_NODE_TYPE          "amplifier"
#define ENV_FEMDD_DEBUG_MODE   "FEMD_DEBUG_MODE"
#define DEF_GPIO_BASE_PATH     "/sys/devices/platform"
#define DEF_SAFETY_CONFIG_PATH "/ukama/configs/femd/saftey_config.json"

typedef enum {
    FEM_UNIT_NONE = 0,
    FEM_UNIT_1    = 1,
    FEM_UNIT_2    = 2
} FemUnit;

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

#endif /* FEMD_H */
