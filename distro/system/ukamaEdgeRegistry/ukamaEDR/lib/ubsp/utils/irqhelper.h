/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef UTILS_IRQHELPER_H_
#define UTILS_IRQHELPER_H_

#include "headers/errorcode.h"
#include "inc/driverfxn.h"
#include "utils/irqdb.h"
#include "headers/ubsp/property.h"

int irqhelper_confirm_irq(const DrvDBFxnTable *drvr_db_fx_tbl, Device *p_dev, AlertCallBackData** acbdata, Property* prop,
                        int max_prop, char* fpath, int* evt);

#endif /* UTILS_IRQHELPER_H_ */
