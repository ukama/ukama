/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
