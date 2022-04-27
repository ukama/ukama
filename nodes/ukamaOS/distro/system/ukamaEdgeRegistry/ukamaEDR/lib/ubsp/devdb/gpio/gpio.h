/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DEVDB_GPIO_GPIO_H_
#define DEVDB_GPIO_GPIO_H_

#include "inc/driverfxn.h"
#include "devdb/gpio/gpiowrapper.h"

/* EDGE */
#define GPIO_EDGE_NONE				0x00
#define GPIO_EDGE_RISING			0x01
#define GPIO_EDGE_FALLING			0x02
#define GPIO_EDGE_BOTH				0x03

/* Direction */
#define GPIO_DIRECTION_INPUT		0x00
#define GPIO_DIRECTION_OUTPUT		0x01

/*Active Low */
#define GPIO_NORMAL					0x00
#define GPIO_INVERT					0x01

int gpio_init ();
int gpio_registration(Device* p_dev);
int gpio_read_prop_count(Device* p_dev, uint16_t * count);
int gpio_read_properties(Device* p_dev, void* prop);
int gpio_configure(void* p_dev, void* prop, void* data );
int gpio_read(void* p_dev, void* prop, void* data);
int gpio_write(void* p_dev, void* prop, void* data);
int gpio_enable(void* p_dev, void* prop, void* data);
int gpio_disable(void* p_dev, void* prop, void* data);

#endif /* DEVDB_GPIO_GPIO_H_ */
