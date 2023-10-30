/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DRVRSYSFS_H_
#define DRVRSYSFS_H_

#include "inc/driverfxn.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define IF_SYSFS_SUPPORT(file) 		((!strcmp(file, "") && !strcmp(file, " "))?0:1)

//TODO: void* hwattr check if this is still required.
const DrvDBFxnTable* drvr_sysfs_get_fxn_tbl();
int drvr_sysfs_init ();
int drvr_sysfs_registration(Device* dev);
int drvr_sysfs_configure(void* hwattr, void* prop , void* data);
int drvr_sysfs_read(void* hwattr, void* prop , void* data );
int drvr_sysfs_write(void* hwattr, void* prop , void* data);
int drvr_sysfs_enable( void* hwattr, void* prop , void* data);
int drvr_sysfs_disable( void* hwattr, void* prop , void* data);
int drvr_sysfs_reg_cb( void* hwattr, SensorCallbackFxn fun );
int drvr_sysfs_dreg_cb( void* hwattr, SensorCallbackFxn fun );
int drvr_sysfs_enable_irq( void* hwattr, void* prop , void* data);
int drvr_sysfs_disable_irq( void* hwattr, void* prop , void* data);

#endif /* DRVRSYSFS_H_ */
