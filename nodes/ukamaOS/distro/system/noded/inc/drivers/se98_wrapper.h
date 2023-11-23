/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DRIVERS_SE98_WRAPPER_H_
#define DRIVERS_SE98_WRAPPER_H_

#include "device.h"

int se98_wrapper_init ();
int se98_wrapper_registration(Device* p_dev);
int se98_wrapper_read_properties(DevObj* obj, void* prop, uint16_t* count);
int se98_wrapper_configure(void* p_dev, void* prop, void* data );
int se98_wrapper_read(void* p_dev, void* prop, void* data);
int se98_wrapper_write(void* p_dev, void* prop, void* data);
int se98_wrapper_enable(void* p_dev, void* prop, void* data);
int se98_wrapper_disable(void* p_dev, void* prop, void* data);
int se98_wrapper_reg_cb(void* p_dev, SensorCallbackFxn fun);
int se98_wrapper_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int se98_wrapper_enable_irq(void* p_dev, void* prop, void* data);
int se98_wrapper_disable_irq(void* p_dev, void* prop, void* data);

#endif /* DRIVERS_SE98_WRAPPER_H_ */
