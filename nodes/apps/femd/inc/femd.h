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

#ifndef STATUS_OK
#define STATUS_OK  0
#endif

#ifndef STATUS_NOK
#define STATUS_NOK (-1)
#endif

#ifndef USYS_TRUE
#define USYS_TRUE  1
#endif

#ifndef USYS_FALSE
#define USYS_FALSE 0
#endif

#define MODULE_FEM "fem"
#define ALARM_NODE "node"
#define EMPTY_STRING ""

#define ALARM_HIGH "HIGH"
#define ALARM_LOW  "LOW"

#define ALARM_PA_AUTO_OFF        "pa_auto_off"
#define ALARM_PA_AUTO_ON         "pa_auto_on"
#define ALARM_PA_AUTO_OFF_DESCRP "PA disabled automatically due to safety"
#define ALARM_PA_AUTO_ON_DESCRP  "PA restored automatically after safety cleared"

#define ALARM_TYPE_PA_OFF 1
#define ALARM_TYPE_PA_ON  2

#endif /* FEMD_H */
