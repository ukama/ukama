/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef UTILS_IRQHELPER_H_
#define UTILS_IRQHELPER_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "errorcode.h"
#include "device_ops.h"
#include "driver_ops.h"
#include "irqdb.h"
#include "property.h"

/**
 * @fn      int irqhelper_confirm_irq(const DrvrOps*, Device*,
 *              AlertCallBackData**, Property*, int, char*, int*)
 * @brief   Confirm the IRQ raised by the sensor. It compares the actual real
 *          time value of the sensor to the limits. If the alert is real it
 *          reports it back otherwise discard it.
 *
 * @param   drveOps
 * @param   p_dev
 * @param   acbdata
 * @param   prop
 * @param   max_prop
 * @param   fpath
 * @param   evt
 * @return  On success, 0
 *          On failure, -1
 */
int irqhelper_confirm_irq(const DrvrOps *drveOps, Device *p_dev, AlertCallBackData** acbdata, Property* prop,
                        int max_prop, char* fpath, int* evt);

#ifdef __cplusplus
}
#endif

#endif /* UTILS_IRQHELPER_H_ */
