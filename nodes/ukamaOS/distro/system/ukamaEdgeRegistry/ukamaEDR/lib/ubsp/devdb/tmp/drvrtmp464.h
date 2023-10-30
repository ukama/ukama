/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef TM464DRIVER_H_
#define TM464DRIVER_H_

#include "devdb/tmp/tmp464.h"

int drvr_tmp464_init ();
int drvr_tmp464_registration(Device* p_dev);
int drvr_tmp464_read_properties(DevObj* obj, void* prop, uint16_t* count);
int drvr_tmp464_configure(void* p_dev, void* prop, void* data );
int drvr_tmp464_read(void* p_dev, void* prop, void* data);
int drvr_tmp464_write(void* p_dev, void* prop, void* data);
int drvr_tmp464_enable(void* p_dev, void* prop, void* data);
int drvr_tmp464_disable(void* p_dev, void* prop, void* data);
int drvr_tmp464_reg_cb(void* p_dev, SensorCallbackFxn fun);
int drvr_tmp464_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int drvr_tmp464_enable_irq(void* p_dev, void* prop, void* data);
int drvr_tmp464_disable_irq(void* p_dev, void* prop, void* data);

#endif /* TM464DRIVER_H_ */
