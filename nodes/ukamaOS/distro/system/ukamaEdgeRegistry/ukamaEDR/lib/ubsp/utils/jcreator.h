/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef UTILS_JCREATOR_H_
#define UTILS_JCREATOR_H_

#include "headers/ubsp/devices.h"
#include "headers/ubsp/ukdblayout.h"

int jcreator_schema( UnitInfo *unit_info, UnitCfg *unit_cfg,
		ModuleInfo* minfo, char** junit_schema);
#endif /* UTILS_JCREATOR_H_ */
