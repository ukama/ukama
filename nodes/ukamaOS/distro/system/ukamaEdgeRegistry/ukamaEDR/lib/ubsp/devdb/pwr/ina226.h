/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DEVDB_PWR_INA226_H_
#define DEVDB_PWR_INA226_H_

#include "devdb/pwr/pwr.h"
#include "inc/driverfxn.h"

int ina226_init ();
int ina226_registration(Device* p_dev);
int ina226_get_irq_type(int pidx, uint8_t* alertstate);
int ina226_read_prop_count(Device* p_dev, uint16_t * count);
int ina226_read_properties(Device* p_dev, void* prop);
int ina226_configure(void* p_dev, void* prop, void* data );
int ina226_read(void* p_dev, void* prop, void* data);
int ina226_write(void* p_dev, void* prop, void* data);
int ina226_enable(void* p_dev, void* prop, void* data);
int ina226_disable(void* p_dev, void* prop, void* data);
int ina226_reg_cb(void* p_dev, SensorCallbackFxn fun);
int ina226_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int ina226_enable_irq(void* p_dev, void* prop, void* data);
int ina226_disable_irq(void* p_dev, void* prop, void* data);
int ina226_confirm_irq(Device *dev, AlertCallBackData** acbdata,
		char* fpath, int* evt);

#endif /* DEVDB_PWR_INA226_H_ */
