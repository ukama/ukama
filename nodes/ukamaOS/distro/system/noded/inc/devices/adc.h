/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef ADC_H_
#define ADC_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "device.h"
#include "device_ops.h"

#include "usys_list.h"

#define MAX_ADC_SENSOR_TYPE		1

/**
 * @fn      const DevOps get_adc_dev_ops*(char*)
 * @brief   Get list of operations that can be performed on ADC sensors.
 *
 * @param   name
 * @return  On success, return ADC sensor operations.
 *          On failure, NULL
 */
const DevOps* get_adc_dev_ops(char *name);

/**
 * @fn      ListInfo get_adc_dev_ldgr*()
 * @brief   Return list of the ADC sensors registered
 *
 * @return  On success List of sensors registered
 *          On failure NULL
 */
ListInfo* get_adc_dev_ldgr();

#ifdef __cplusplus
}
#endif

#endif /* ADC_H_ */
