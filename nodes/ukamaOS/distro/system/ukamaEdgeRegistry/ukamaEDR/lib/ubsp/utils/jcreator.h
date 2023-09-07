/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef UTILS_JCREATOR_H_
#define UTILS_JCREATOR_H_

#include "headers/ubsp/devices.h"
#include "headers/ubsp/ukdblayout.h"

int jcreator_schema( UnitInfo *unit_info, UnitCfg *unit_cfg,
		ModuleInfo* minfo, char** junit_schema);
#endif /* UTILS_JCREATOR_H_ */
