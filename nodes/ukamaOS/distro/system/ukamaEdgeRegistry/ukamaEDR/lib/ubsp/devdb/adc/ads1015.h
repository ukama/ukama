/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DEVDB_ADC_ADS1015_H_
#define DEVDB_ADC_ADS1015_H_

#include "devdb/adc/adc.h"
#include "inc/driverfxn.h"

int ads1015_init ();
int ads1015_registration(Device* p_dev);
int ads1015_read_prop_count(Device* p_dev, uint16_t * count);
int ads1015_read_properties(Device* p_dev, void* prop);
int ads1015_configure(void* p_dev, void* prop, void* data );
int ads1015_read(void* p_dev, void* prop, void* data);
int ads1015_write(void* p_dev, void* prop, void* data);
int ads1015_enable(void* p_dev, void* prop, void* data);
int ads1015_disable(void* p_dev, void* prop, void* data);

#endif /* DEVDB_ADC_ADS1015_H_ */
