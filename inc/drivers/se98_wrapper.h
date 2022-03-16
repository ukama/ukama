/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DRIVERS_SE98_WRAPPER_H_
#define DRIVERS_SE98_WRAPPER_H_

#include "device.h"

int se98_wrapperinit ();
int se98_wrapperregistration(Device* p_dev);
int se98_wrapperread_properties(DevObj* obj, void* prop, uint16_t* count);
int se98_wrapperconfigure(void* p_dev, void* prop, void* data );
int se98_wrapperread(void* p_dev, void* prop, void* data);
int se98_wrapperwrite(void* p_dev, void* prop, void* data);
int se98_wrapperenable(void* p_dev, void* prop, void* data);
int se98_wrapperdisable(void* p_dev, void* prop, void* data);
int se98_wrapperreg_cb(void* p_dev, SensorCallbackFxn fun);
int se98_wrapperdreg_cb(void* p_dev, SensorCallbackFxn fun);
int se98_wrapperenable_irq(void* p_dev, void* prop, void* data);
int se98_wrapperdisable_irq(void* p_dev, void* prop, void* data);


#endif /* DRIVERS_SE98_WRAPPER_H_ */
