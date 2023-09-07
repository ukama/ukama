/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
