/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef BSP_ADS1015_H_
#define BSP_ADS1015_H_

#include "device.h"

int bsp_ads1015_init();
int bsp_ads1015_registration(Device* p_dev);
int bsp_ads1015_read_prop_count(Device* p_dev, uint16_t * count);
int bsp_ads1015_read_properties(Device* p_dev, void* prop);
int bsp_ads1015_configure(void* p_dev, void* prop, void* data );
int bsp_ads1015_read(void* p_dev, void* prop, void* data);
int bsp_ads1015_write(void* p_dev, void* prop, void* data);
int bsp_ads1015_enable(void* p_dev, void* prop, void* data);
int bsp_ads1015_disable(void* p_dev, void* prop, void* data);

#endif /* DEVICES_BSP_ADS1015_H_ */
