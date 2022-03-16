/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DRIVERS_GPIO_WRAPPER_H_
#define DRIVERS_GPIO_WRAPPER_H_

#include "device.h"

int gpio_wrapper_init ();
int gpio_wrapper_registration(Device* p_dev);
int gpio_wrapper_read_properties(DevObj* obj, void* prop, uint16_t* count);
int gpio_wrapper_configure(void* p_dev, void* prop, void* data );
int gpio_wrapper_read(void* p_dev, void* prop, void* data);
int gpio_wrapper_write(void* p_dev, void* prop, void* data);
int gpio_wrapper_enable(void* p_dev, void* prop, void* data);
int gpio_wrapper_disable(void* p_dev, void* prop, void* data);


#endif /* DRIVERS_GPIO_WRAPPER_H_ */
