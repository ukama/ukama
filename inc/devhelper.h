/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_DEVHELPER_H_
#define INC_DEVHELPER_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "driver_ops.h"

int dhelper_registration(const DrvrOps *drvr, Device *p_dev);
void dhelper_irq_callback(DevObj *obj, void *prop, void *data);
int dhelper_validate_property(Property* prop, int pidx);
int dhelper_init_property_from_parser(Device *p_dev, Property** prop,
    void* count);
int dhelper_init_driver(const DrvrOps *drvr, Device *dev);
int dhelper_configure(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data);
int dhelper_read(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data) ;
int dhelper_write(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data);
int dhelper_enable(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data);
int dhelper_disable(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data);
int dhelper_enable_irq(const DrvrOps *drvr, SensorCallbackFxn sensor_cb,
    Device *dev, Property *prop, int pidx, void *data);
int dhelper_disable_irq(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data);
int dhelper_confirm_irq(const DrvrOps *drvr,Device *dev, Property *prop, AlertCallBackData **acbdata, char *fpath,
                        int maxpcount, int *evt);

#ifdef __cplusplus
}
#endif

#endif /* INC_DEVHELPER_H_ */
