/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DRIVERS_TM464_WRAPPER_H_
#define DRIVERS_TM464_WRAPPER_H_

#include "device.h"

int tmp464_wrapper_init ();
int tmp464_wrapper_registration(Device* p_dev);
int tmp464_wrapper_read_properties(DevObj* obj, void* prop, uint16_t* count);
int tmp464_wrapper_configure(void* p_dev, void* prop, void* data );
int tmp464_wrapper_read(void* p_dev, void* prop, void* data);
int tmp464_wrapper_write(void* p_dev, void* prop, void* data);
int tmp464_wrapper_enable(void* p_dev, void* prop, void* data);
int tmp464_wrapper_disable(void* p_dev, void* prop, void* data);
int tmp464_wrapper_reg_cb(void* p_dev, SensorCallbackFxn fun);
int tmp464_wrapper_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int tmp464_wrapper_enable_irq(void* p_dev, void* prop, void* data);
int tmp464_wrapper_disable_irq(void* p_dev, void* prop, void* data);

#endif /* DRIVERS_TM464_WRAPPER_H_ */
