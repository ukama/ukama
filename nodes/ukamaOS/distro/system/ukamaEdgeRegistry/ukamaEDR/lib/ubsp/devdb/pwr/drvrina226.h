/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DEVDB_PWR_DRVRINA226_H_
#define DEVDB_PWR_DRVRINA226_H_

#include "devdb/pwr/ina226.h"

int drvr_ina226_init ();
int drvr_ina226_registration(Device* p_dev);
int drvr_ina226_read_properties(DevObj* obj, void* prop, uint16_t* count);
int drvr_ina226_configure(void* p_dev, void* prop, void* data );
int drvr_ina226_read(void* p_dev, void* prop, void* data);
int drvr_ina226_write(void* p_dev, void* prop, void* data);
int drvr_ina226_enable(void* p_dev, void* prop, void* data);
int drvr_ina226_disable(void* p_dev, void* prop, void* data);
int drvr_ina226_reg_cb(void* p_dev, SensorCallbackFxn fun);
int drvr_ina226_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int drvr_ina226_enable_irq(void* p_dev, void* prop, void* data);
int drvr_ina226_disable_irq(void* p_dev, void* prop, void* data);

#endif /* DEVDB_PWR_DRVRINA226_H_ */
