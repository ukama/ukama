/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DRIVERFXN_H_
#define DRIVERFXN_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "device.h"

typedef int (*DrvrInitFxn)(Device* p_dev);
typedef int (*DrvrConfigFxn)(void* obj, void* prop, void* data);
typedef int (*DrvrReadFxn)( void* obj, void* prop, void* data);
typedef int (*DrvrWriteFxn)(void* obj, void* prop, void* data);
typedef int (*DrvrEnableFxn)(void* obj, void* prop, void* data);
typedef int (*DrvrDisableFxn)( void* obj, void* prop, void* data);
typedef int (*DrvrRegCB)( void* obj,SensorCallbackFxn fun);
typedef int (*DrvrDregCB)( void* obj, SensorCallbackFxn fun);
typedef int (*DrvrEnableIRQ)( void* obj, void* prop, void* data);
typedef int (*DrvrDisableIRQ)( void* obj, void* prop, void* data);

/* Driver operation supported */
typedef struct  {
  DrvrInitFxn init;
  DrvrConfigFxn configure;
  DrvrReadFxn read;
  DrvrWriteFxn write;
  DrvrEnableFxn enable;
  DrvrDisableFxn disable;
  DrvrRegCB	  registerCb;
  DrvrDregCB	  dregisterCb;
  DrvrEnableIRQ  enableIrq;
  DrvrDisableIRQ disableIrq;
} DrvrOps;


#ifdef __cplusplus
}
#endif

#endif /* DRIVERFXN_H_ */
