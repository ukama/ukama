/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DRIVERS_SYSFS_H_
#define DRIVERS_SYSFS_H_

#include "driver_ops.h"


#define IF_SYSFS_SUPPORT(file) 		((!strcmp(file, "") && !strcmp(file, " "))?0:1)

//TODO: void* hwAttrs check if this is still required.
const DrvrOps* sysfs_wrapper_get_ops();
int sysfs_wrapper_init ();
int sysfs_wrapper_registration(Device* dev);
int sysfs_wrapper_configure(void* hwAttrs, void* prop , void* data);
int sysfs_wrapper_read(void* hwAttrs, void* prop , void* data );
int sysfs_wrapper_write(void* hwAttrs, void* prop , void* data);
int sysfs_wrapper_enable( void* hwAttrs, void* prop , void* data);
int sysfs_wrapper_disable( void* hwAttrs, void* prop , void* data);
int sysfs_wrapper_reg_cb( void* hwAttrs, SensorCallbackFxn fun );
int sysfs_wrapper_dreg_cb( void* hwAttrs, SensorCallbackFxn fun );
int sysfs_wrapper_enable_irq( void* hwAttrs, void* prop , void* data);
int sysfs_wrapper_disable_irq( void* hwAttrs, void* prop , void* data);

#endif /* DRIVERS_SYSFS_H_ */
