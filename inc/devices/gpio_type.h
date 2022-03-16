/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef GPIOTYPE_H_
#define GPIOTYPE_H_

#include "usys_list.h"

#define MAX_GPIO_SENSOR_TYPE		1

const DevOps* get_gpio_type_dev_ops(char *name);
ListInfo* get_gpio_type_dev_ldgr();

#endif /* GPIOWRAPPER_H_ */
