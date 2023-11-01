/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DRIVERS_LEDTRICOL_WRAPPER_H_
#define DRIVERS_LEDTRICOL_WRAPPER_H_

#include "device.h"

int led_tricol_wrapper_init ();
int led_tricol_wrapper_registration(Device* p_dev);
int led_tricol_wrapper_read_properties(DevObj* obj, void* prop, uint16_t* count);
int led_tricol_wrapper_configure(void* p_dev, void* prop, void* data );
int led_tricol_wrapper_read(void* p_dev, void* prop, void* data);
int led_tricol_wrapper_write(void* p_dev, void* prop, void* data);
int led_tricol_wrapper_enable(void* p_dev, void* prop, void* data);
int led_tricol_wrapper_disable(void* p_dev, void* prop, void* data);

#endif /*DRIVERS_LEDTRICOL_WRAPPER_H_*/
