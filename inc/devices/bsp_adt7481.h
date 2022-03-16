/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef BSP_ADT7841_H_
#define BSP_ADT7841_H_

#include "device.h"

int bsp_adt7481_init (Device* p_dev);
int bsp_adt7481_registration(Device* p_dev);
int bsp_adt7481_get_irq_type(int bsp_pidx, uint8_t* alertstate);
int bsp_adt7481_read_prop_count(Device* p_dev, uint16_t * count);
int bsp_adt7481_read_properties(Device* p_dev, void* prop);
int bsp_adt7481_configure(void* p_dev, void* prop, void* data);
int bsp_adt7481_read(void* p_dev, void* prop, void* data );
int bsp_adt7481_write(void* p_dev, void* prop, void* data);
int bsp_adt7481_enable(void* p_dev, void* prop, void* data);
int bsp_adt7481_disable(void* p_dev, void* prop, void* data);
int bsp_adt7481_reg_cb(void* p_dev, SensorCallbackFxn fun);
int bsp_adt7481_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int bsp_adt7481_enable_irq(void* p_dev, void* prop, void* data);
int bsp_adt7481_disable_irq(void* p_dev, void* prop, void* data);
int bsp_adt7481_confirm_irq(Device *p_dev, AlertCallBackData** acbdata, char* fpath, int* count);

#endif /* BSP_ADT7841_H_ */
