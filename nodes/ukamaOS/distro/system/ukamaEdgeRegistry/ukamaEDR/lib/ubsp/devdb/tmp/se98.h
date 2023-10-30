/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
