/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DEVDB_GPIO_DRVRGPIO_H_
#define DEVDB_GPIO_DRVRGPIO_H_

#include "devdb/gpio/gpio.h"

int drvr_gpio_init ();
int drvr_gpio_registration(Device* p_dev);
int drvr_gpio_read_properties(DevObj* obj, void* prop, uint16_t* count);
int drvr_gpio_configure(void* p_dev, void* prop, void* data );
int drvr_gpio_read(void* p_dev, void* prop, void* data);
int drvr_gpio_write(void* p_dev, void* prop, void* data);
int drvr_gpio_enable(void* p_dev, void* prop, void* data);
int drvr_gpio_disable(void* p_dev, void* prop, void* data);


#endif /* DEVDB_GPIO_DRVRGPIO_H_ */
