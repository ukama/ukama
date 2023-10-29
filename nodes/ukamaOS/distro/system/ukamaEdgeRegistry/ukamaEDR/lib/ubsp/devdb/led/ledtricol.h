/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
