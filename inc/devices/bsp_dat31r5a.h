/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef BSP_DAT31R5A_H_
#define BSP_DAT31R5A_H_

#include "device.h"

int bsp_dat31r5a_init ();
int bsp_dat31r5a_registration(Device* p_dev);
int bsp_dat31r5a_read_prop_count(Device* p_dev, uint16_t * count);
int bsp_dat31r5a_read_properties(Device* p_dev, void* prop);
int bsp_dat31r5a_configure(void* p_dev, void* prop, void* data );
int bsp_dat31r5a_read(void* p_dev, void* prop, void* data);
int bsp_dat31r5a_write(void* p_dev, void* prop, void* data);
int bsp_dat31r5a_enable(void* p_dev, void* prop, void* data);
int bsp_dat31r5a_disable(void* p_dev, void* prop, void* data);

#endif /* BSP_DAT31R5A_H_ */
