/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef ADC_H_
#define ADC_H_

#include "device.h"
#include "device_ops.h"

#include "usys_list.h"

#define MAX_ADC_SENSOR_TYPE		1

void  clean_dev_adc_prop();
const DevOps* get_adc_dev_ops(char *name);
ListInfo* get_adc_dev_ldgr();

#endif /* ADC_H_ */
