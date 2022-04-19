/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DEVDB_LED_LEDTRICOL_H_
#define DEVDB_LED_LEDTRICOL_H_

#include "devdb/led/led.h"
#include "inc/driverfxn.h"

int led_tricol_init ();
int led_tricol_registration(Device* p_dev);
int led_tricol_read_prop_count(Device* p_dev, uint16_t * count);
int led_tricol_read_properties(Device* p_dev, void* prop);
int led_tricol_configure(void* p_dev, void* prop, void* data );
int led_tricol_read(void* p_dev, void* prop, void* data);
int led_tricol_write(void* p_dev, void* prop, void* data);
int led_tricol_enable(void* p_dev, void* prop, void* data);
int led_tricol_disable(void* p_dev, void* prop, void* data);

#endif /* DEVDB_LED_LEDTRICOL_H_ */
