/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_DEVHELPER_H_
#define INC_DEVHELPER_H_
#include "inc/driverfxn.h"

int dhelper_registration(const DrvDBFxnTable *drvr, Device *p_dev);
void dhelper_irq_callback(DevObj *obj, void *prop, void *data);
int dhelper_validate_property(Property* prop, int pidx);
int dhelper_init_property_from_parser(Device *p_dev, Property** prop,
		void* count);
int dhelper_init_driver(const DrvDBFxnTable *drvr, Device *dev);
int dhelper_configure(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
		int pidx, void *data);
int dhelper_read(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
		int pidx, void *data) ;
int dhelper_write(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
		int pidx, void *data);
int dhelper_enable(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
		int pidx, void *data);
int dhelper_disable(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
		int pidx, void *data);
int dhelper_enable_irq(const DrvDBFxnTable *drvr, SensorCallbackFxn sensor_cb,
		Device *dev, Property *prop, int pidx, void *data);
int dhelper_disable_irq(const DrvDBFxnTable *drvr, Device *dev, Property *prop,
		int pidx, void *data);
int dhelper_confirm_irq(const DrvDBFxnTable *drvr,Device *dev, Property *prop, AlertCallBackData **acbdata, char *fpath,
                        int maxpcount, int *evt);
#endif /* INC_DEVHELPER_H_ */
