/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DRIVERS_DAT31R5A_WRAPPER_H_
#define DRIVERS_DAT31R5A_WRAPPER_H_

#include "device.h"

int dat31r5a_wrapper_init ();
int dat31r5a_wrapper_registration(Device* p_dev);
int dat31r5a_wrapper_read_properties(DevObj* obj, void* prop, uint16_t* count);
int dat31r5a_wrapper_configure(void* p_dev, void* prop, void* data );
int dat31r5a_wrapper_read(void* p_dev, void* prop, void* data);
int dat31r5a_wrapper_write(void* p_dev, void* prop, void* data);
int dat31r5a_wrapper_enable(void* p_dev, void* prop, void* data);
int dat31r5a_wrapper_disable(void* p_dev, void* prop, void* data);

#endif /* DRIVERS_DAT31R5A_WRAPPER_H_ */
