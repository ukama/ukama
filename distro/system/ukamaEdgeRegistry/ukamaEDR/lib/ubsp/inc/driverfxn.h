/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
