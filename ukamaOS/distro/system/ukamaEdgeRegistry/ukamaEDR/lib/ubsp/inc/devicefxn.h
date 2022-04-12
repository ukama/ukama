/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DEVICEFXN_H_
#define DEVICEFXN_H_

#include "headers/ubsp/devices.h"

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

/*basic read write operation to Temperature devices*/
typedef struct  {
	DevInitFxn init;
	DevRegisterFxn registration;
	DevReadPropCountFxn read_prop_count;
	DevReadPropFxn read_prop;
	DevConfigFxn configure;
	DevReadFxn read;
	DevWriteFxn write;
	DevEnableFxn enable;
	DevDisableFxn disable;
	DevRegCB	  register_cb;
	DevDregCB	  dregister_cb;
	DevEnableIRQ  enable_irq;
	DevDisableIRQ disable_irq;
	DevConfirmIRQ confirm_irq;
	DevGetIRQType irq_type;
} DevFxnTable;

#endif /* DEVICEFXN_H_ */
