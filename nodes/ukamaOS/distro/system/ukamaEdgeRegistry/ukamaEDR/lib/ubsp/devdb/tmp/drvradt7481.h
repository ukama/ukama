/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DEVDB_TMP_DRVRADT7481_H_
#define DEVDB_TMP_DRVRADT7481_H_

#include "devdb/tmp/adt7481.h"

int drvr_adt7841_init ();
int drvr_adt7841_registration(Device* p_dev);
int drvr_adt7841_read_properties(DevObj* obj, void* prop, uint16_t* count);
int drvr_adt7841_configure(void* p_dev, void* prop, void* data );
int drvr_adt7841_read(void* p_dev, void* prop, void* data);
int drvr_adt7841_write(void* p_dev, void* prop, void* data);
int drvr_adt7841_enable(void* p_dev, void* prop, void* data);
int drvr_adt7841_disable(void* p_dev, void* prop, void* data);
int drvr_adt7841_reg_cb(void* p_dev, SensorCallbackFxn fun);
int drvr_adt7841_dreg_cb(void* p_dev, SensorCallbackFxn fun);
int drvr_adt7841_enable_irq(void* p_dev, void* prop, void* data);
int drvr_adt7841_disable_irq(void* p_dev, void* prop, void* data);

#endif /* DEVDB_TMP_DRVRADT7481_H_ */
