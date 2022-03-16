/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef BSP_SE98_H_
#define BSP_SE98_H_

#include "device.h"

int bsp_se98_init (Device* p_dev);
int bsp_se98_registration(Device* p_dev);
int bsp_se98_get_irq_type(int bsp_pidx, uint8_t* alertstate);
int bsp_se98_read_prop_count(Device* p_dev, uint16_t * count);
int bsp_se98_read_properties(Device* p_dev, void* prop);
int bsp_se98_configure(void* p_dev, void* prop, void* data);
int bsp_se98_read(void* p_dev, void* prop, void* data );
int bsp_se98_write(void* p_dev, void* prop, void* data);
int bsp_se98_enable(void* p_dev, void* prop, void* data);
int bsp_se98_disable(void* p_dev, void* prop, void* data);
int bsp_se98_reg_cb(void* p_dev, SensorCallbackFxn fun);
int bsp_se98_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int bsp_se98_enable_irq(void* p_dev, void* prop, void* data);
int bsp_se98_disable_irq(void* p_dev, void* prop, void* data);
int bsp_se98_confirm_irq(Device *dev, AlertCallBackData** acbdata,
    char* fpath, int* evt);

#endif /* BSP_SE98_H_ */
