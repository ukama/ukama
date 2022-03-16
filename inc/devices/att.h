/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef ATT_H_
#define ATT_H_

#include "usys_list.h"

#define MAX_ATT_SENSOR_TYPE		1

const DevOps* get_att_dev_ops(char *name);
ListInfo* get_att_dev_ldgr();

#endif /* ATT_H_ */
