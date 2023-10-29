/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DEVDB_LED_DRVRLEDTRICOL_H_
#define DEVDB_LED_DRVRLEDTRICOL_H_

#include "devdb/led/ledtricol.h"

int drvr_led_tricol_init ();
int drvr_led_tricol_registration(Device* p_dev);
int drvr_led_tricol_read_properties(DevObj* obj, void* prop, uint16_t* count);
int drvr_led_tricol_configure(void* p_dev, void* prop, void* data );
int drvr_led_tricol_read(void* p_dev, void* prop, void* data);
int drvr_led_tricol_write(void* p_dev, void* prop, void* data);
int drvr_led_tricol_enable(void* p_dev, void* prop, void* data);
int drvr_led_tricol_disable(void* p_dev, void* prop, void* data);

#endif /*DEVDB_LED_DRVRLEDTRICOL_H_*/
