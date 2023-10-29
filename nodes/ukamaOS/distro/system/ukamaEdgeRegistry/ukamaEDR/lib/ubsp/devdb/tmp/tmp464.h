/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef TMP464_H_
#define TMP464_H_

#include "inc/driverfxn.h"
#include "devdb/tmp/tmp.h"

int tmp464_init ();
int tmp464_registration(Device* p_dev);
int tmp464_get_irq_type(int pidx, uint8_t* alertstate);
int tmp464_read_prop_count(Device* p_dev, uint16_t * count);
int tmp464_read_properties(Device* p_dev, void* prop);
int tmp464_configure(void* p_dev, void* prop, void* data );
int tmp464_read(void* p_dev, void* prop, void* data);
int tmp464_write(void* p_dev, void* prop, void* data);
int tmp464_enable(void* p_dev, void* prop, void* data);
int tmp464_disable(void* p_dev, void* prop, void* data);
int tmp464_reg_cb(void* p_dev, SensorCallbackFxn fun);
int tmp464_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int tmp464_enable_irq(void* p_dev, void* prop, void* data);
int tmp464_disable_irq(void* p_dev, void* prop, void* data);
int tmp464_confirm_irq(Device *dev, AlertCallBackData** acbdata,
		char* fpath, int* evt);

#endif /* TMP464_H_ */
