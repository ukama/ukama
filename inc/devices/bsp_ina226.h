/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef BSP_INA226_H_
#define BSP_INA226_H_

#include "device.h"

int bsp_ina226_init ();
int bsp_ina226_registration(Device* p_dev);
int bsp_ina226_get_irq_type(int bsp_pidx, uint8_t* alertstate);
int bsp_ina226_read_prop_count(Device* p_dev, uint16_t * count);
int bsp_ina226_read_properties(Device* p_dev, void* prop);
int bsp_ina226_configure(void* p_dev, void* prop, void* data );
int bsp_ina226_read(void* p_dev, void* prop, void* data);
int bsp_ina226_write(void* p_dev, void* prop, void* data);
int bsp_ina226_enable(void* p_dev, void* prop, void* data);
int bsp_ina226_disable(void* p_dev, void* prop, void* data);
int bsp_ina226_reg_cb(void* p_dev, SensorCallbackFxn fun);
int bsp_ina226_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int bsp_ina226_enable_irq(void* p_dev, void* prop, void* data);
int bsp_ina226_disable_irq(void* p_dev, void* prop, void* data);
int bsp_ina226_confirm_irq(Device *dev, AlertCallBackData** acbdata,
    char* fpath, int* evt);

#endif /* BSP_INA226_H_ */
