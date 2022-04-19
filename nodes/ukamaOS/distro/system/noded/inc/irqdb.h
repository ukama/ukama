/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef UTILS_IRQDB_H_
#define UTILS_IRQDB_H_

#ifdef __cplusplus
extern "C" {
#endif

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
    char fName[64];
    /* Added if we are looking notifications from file change.*/
} IRQCfg;

/**
 * @fn       int irqdb_deregister_for_device_irq(IRQSrcInfo*, ThreadedIRQEnable)
 * @brief    De-register interrupts for the sensor device by canceling the
 *           thread which monitors for the changes in sysfs file and removes
 *           the thread id from list
 *
 * @param    rsrc
 * @param    IRQ_disable
 * @return   On success, 0
 *           On failure, non zero value
 */
int irqdb_deregister_for_device_irq(IRQSrcInfo* rsrc,
                ThreadedIRQEnable IRQ_disable);

/**
 * @fn      int irqdb_register_for_device_irq(IRQSrcInfo*, SensorCallbackFxn,
 *          ThreadedIRQEnable)
 * @brief   Register for alerts by sensor device. Creates a thread which
 *          monitor for changes in sysfs files and adds thread id to list.
 *
 * @param   rsrc
 * @param   cb
 * @param   IRQ_enable
 * @return  On success, 0
 *          On failure, non zero value
 */
int irqdb_register_for_device_irq(IRQSrcInfo* rsrc, SensorCallbackFxn cb,
    ThreadedIRQEnable IRQ_enable);

/**
 * @fn      void irqdb_exit()
 * @brief   Remove all the entries from IRQdb list and destroys it.
 *
 */
void irqdb_exit();

/**
 * @fn      void irqdb_init()
 * @brief   Create a IRQDV list for storing registered alerts.
 *
 */
void irqdb_init();

/**
 * @fn      void irqdb_print_list()
 * @brief   print the information associated with registered alerts.
 *
 */
void irqdb_print_list();

/**
 * @fn      IRQCfg irqdb_search_for_device_irq*(Device*)
 * @brief   Search for a IRQ's registerd for the device in IRQDB list.
 *
 * @param   dev
 * @return  On success, 0
 *          On failure, non zero value
 */
IRQCfg* irqdb_search_for_device_irq(Device* dev);

#ifdef __cplusplus
}
#endif

#endif /* UTILS_IRQDB_H_ */
