/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
