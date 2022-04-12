/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef SE98_H_
#define SE98_H_

#include "inc/driverfxn.h"
#include "devdb/tmp/tmp.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int se98_init (Device* p_dev);
int se98_registration(Device* p_dev);
int se98_get_irq_type(int pidx, uint8_t* alertstate);
int se98_read_prop_count(Device* p_dev, uint16_t * count);
int se98_read_properties(Device* p_dev, void* prop);
int se98_configure(void* p_dev, void* prop, void* data);
int se98_read(void* p_dev, void* prop, void* data );
int se98_write(void* p_dev, void* prop, void* data);
int se98_enable(void* p_dev, void* prop, void* data);
int se98_disable(void* p_dev, void* prop, void* data);
int se98_reg_cb(void* p_dev, SensorCallbackFxn fun);
int se98_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int se98_enable_irq(void* p_dev, void* prop, void* data);
int se98_disable_irq(void* p_dev, void* prop, void* data);
int se98_confirm_irq(Device *dev, AlertCallBackData** acbdata,
		char* fpath, int* evt);

#endif /* SE98_H_ */
