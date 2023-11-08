/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DEVICEFXN_H_
#define DEVICEFXN_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "device.h"

typedef int (*DevInitFxn)();
typedef int (*DevRegisterFxn)(Device* p_dev);
typedef int (*DevReadPropCountFxn)(Device* obj, uint16_t* count);
typedef int (*DevReadPropFxn)(Device* obj, void* prop);
typedef int (*DevConfigFxn)(void* obj, void* prop, void* data);
typedef int (*DevReadFxn)( void* obj, void* prop, void* data);
typedef int (*DevWriteFxn)(void* obj, void* prop, void* data);
typedef int (*DevEnableFxn)(void* obj, void* prop, void* data);
typedef int (*DevDisableFxn)( void* obj, void* prop, void* data);
typedef int (*DevRegCB)( void* obj, SensorCallbackFxn fun);
typedef int (*DevDregCB)( void* obj, SensorCallbackFxn fun);
typedef int (*DevEnableIRQ)( void* obj, void* prop, void* data);
typedef int (*DevDisableIRQ)( void* obj, void* prop, void* data);
typedef int (*DevGetIRQType)( int pidx, uint8_t* irqtype);
typedef int (*DevConfirmIRQ)( Device *dev, AlertCallBackData** acbdata,
    char* fpath, int* evt);

/* basic read write operation to any devices*/
typedef struct  {
  DevInitFxn init;
  DevRegisterFxn registration;
  DevReadPropCountFxn readPropCount;
  DevReadPropFxn readProp;
  DevConfigFxn configure;
  DevReadFxn read;
  DevWriteFxn write;
  DevEnableFxn enable;
  DevDisableFxn disable;
  DevRegCB	  registerCb;
  DevDregCB	  dregisterCb;
  DevEnableIRQ  enableIrq;
  DevDisableIRQ disableIrq;
  DevConfirmIRQ confirmIrq;
  DevGetIRQType irqType;
} DevOps;

#ifdef __cplusplus
}
#endif

#endif /* DEVICEFXN_H_ */
