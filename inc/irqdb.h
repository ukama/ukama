/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef UTILS_IRQDB_H_
#define UTILS_IRQDB_H_

#include "device.h"

#include "usys_thread.h"
#include "usys_types.h"

#define MAX_IRQ 	100

typedef enum {
  IRQ_GPIO = 0,
  IRQ_SYSFS = 1,
}IRQSrcType;

typedef union {
  char sysFsName[64];
  DevI2cCfg devGpio;
}IRQSrc;

typedef struct  __attribute__((__packed__)) {
  IRQSrcType type;
  IRQSrc src;
  DevObj obj;
}IRQSrcInfo;

typedef void (*ThreadedIRQCallback)(void *context);
typedef void (*ThreadedIRQEnable)(void *context);

typedef struct IRQCfg {
  pthread_t pthread;
  SensorCallbackFxn cb; /*Callback to run when interrupt occurs*/
    DevObj obj;           /*Pointer to pass to cb function*/
    char fName[64];  /* Added if we are looking notifications from file change.*/
} IRQCfg;

void irqdb_init();
void irqdb_exit();
IRQCfg* irqdb_search_for_device_irq(Device* dev);
int irqdb_register_for_device_irq(IRQSrcInfo* rsrc, SensorCallbackFxn cb,
    ThreadedIRQEnable IRQ_enable);
int irqdb_deregister_for_device_irq(IRQSrcInfo* rsrc, ThreadedIRQEnable IRQ_disable);
void irqdb_print_list();
#endif /* UTILS_IRQDB_H_ */
