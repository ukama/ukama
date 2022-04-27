/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
