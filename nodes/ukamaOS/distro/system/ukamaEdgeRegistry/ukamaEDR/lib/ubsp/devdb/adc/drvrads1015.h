/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DEVDB_ADC_DRVRADS1015_H_
#define DEVDB_ADC_DRVRADS1015_H_

#include "devdb/adc/ads1015.h"

int drvr_ads1015_init ();
int drvr_ads1015_registration(Device* p_dev);
int drvr_ads1015_read_properties(DevObj* obj, void* prop, uint16_t* count);
int drvr_ads1015_configure(void* p_dev, void* prop, void* data );
int drvr_ads1015_read(void* p_dev, void* prop, void* data);
int drvr_ads1015_write(void* p_dev, void* prop, void* data);
int drvr_ads1015_enable(void* p_dev, void* prop, void* data);
int drvr_ads1015_disable(void* p_dev, void* prop, void* data);

#endif /* DEVDB_ADC_DRVRADS1015_H_ */
