/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
