/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DRIVERFXN_H_
#define DRIVERFXN_H_

#include "headers/ubsp/devices.h"

typedef int (*DrvDBInitFxn)(Device* p_dev);
typedef int (*DrvDBConfigFxn)(void* obj, void* prop, void* data);
typedef int (*DrvDBReadFxn)( void* obj, void* prop, void* data);
typedef int (*DrvDBWriteFxn)(void* obj, void* prop, void* data);
typedef int (*DrvDBEnableFxn)(void* obj, void* prop, void* data);
typedef int (*DrvDBDisableFxn)( void* obj, void* prop, void* data);
typedef int (*DrvDBRegCB)( void* obj,SensorCallbackFxn fun);
typedef int (*DrvDBDregCB)( void* obj, SensorCallbackFxn fun);
typedef int (*DrvDBEnableIRQ)( void* obj, void* prop, void* data);
typedef int (*DrvDBDisableIRQ)( void* obj, void* prop, void* data);

typedef struct  {
	DrvDBInitFxn init;
	DrvDBConfigFxn configure;
	DrvDBReadFxn read;
	DrvDBWriteFxn write;
	DrvDBEnableFxn enable;
	DrvDBDisableFxn disable;
	DrvDBRegCB	  register_cb;
	DrvDBDregCB	  dregister_cb;
	DrvDBEnableIRQ  enable_irq;
	DrvDBDisableIRQ disable_irq;
} DrvDBFxnTable;

#endif /* DRIVERFXN_H_ */
