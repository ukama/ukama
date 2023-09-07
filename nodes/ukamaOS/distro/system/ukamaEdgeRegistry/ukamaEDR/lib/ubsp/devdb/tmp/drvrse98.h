/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DEVDB_TMP_DRVRSE98_H_
#define DEVDB_TMP_DRVRSE98_H_

#include "devdb/tmp/se98.h"

int drvr_se98_init ();
int drvr_se98_registration(Device* p_dev);
int drvr_se98_read_properties(DevObj* obj, void* prop, uint16_t* count);
int drvr_se98_configure(void* p_dev, void* prop, void* data );
int drvr_se98_read(void* p_dev, void* prop, void* data);
int drvr_se98_write(void* p_dev, void* prop, void* data);
int drvr_se98_enable(void* p_dev, void* prop, void* data);
int drvr_se98_disable(void* p_dev, void* prop, void* data);
int drvr_se98_reg_cb(void* p_dev, SensorCallbackFxn fun);
int drvr_se98_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int drvr_se98_enable_irq(void* p_dev, void* prop, void* data);
int drvr_se98_disable_irq(void* p_dev, void* prop, void* data);


#endif /* DEVDB_TMP_DRVRSE98_H_ */
