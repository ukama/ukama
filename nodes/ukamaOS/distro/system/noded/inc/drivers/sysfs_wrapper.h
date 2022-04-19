/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef SYSFS_WRAPPER_H_
#define SYSFS_WRAPPER_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "driver_ops.h"

#define IF_SYSFS_SUPPORT(file) 		((!strcmp(file, "") && \
                !strcmp(file, " "))?0:1)

/**
 * @fn      int sysfs_wrapper_configure(void*, void*, void*)
 * @brief   TBU
 *
 * @param   hwAttrs
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int sysfs_wrapper_configure(void* hwAttrs, void* prop , void* data);

/**
 * @fn      int sysfs_wrapper_enable(void*, void*, void*)
 * @brief   Wrapper to sysfs driver for enabling sensor device
 *
 * @param   hwAttrs
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int sysfs_wrapper_enable( void* hwAttrs, void* prop , void* data);

/**
 * @fn      int sysfs_wrapper_enable_irq(void*, void*, void*)
 * @brief   Wrapper to sysfs driver for enabling IRQ for sensor device
 *
 * @param   hwAttrs
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int sysfs_wrapper_enable_irq( void* hwAttrs, void* prop , void* data);

/**
 * @fn      int sysfs_wrapper_disable(void*, void*, void*)
 * @brief   Wrapper to sysfs driver for disable sensor device
 *
 * @param   hwAttrs
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int sysfs_wrapper_disable( void* hwAttrs, void* prop , void* data);

/**
 * @fn      int sysfs_wrapper_disable_irq(void*, void*, void*)
 * @brief    Wrapper to sysfs driver for disble IRQ for sensor device
 *
 * @param   hwAttrs
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int sysfs_wrapper_disable_irq( void* hwAttrs, void* prop , void* data);

/**
 * @fn      int sysfs_wrapper_dreg_cb(void*, SensorCallbackFxn)
 * @brief   de-register callback functions for the IRQ's for sensor device.
 *
 * @param   hwAttrs
 * @param   fun
 * @return  On success, 0
 *          On failure, non zero value
 */
int sysfs_wrapper_dreg_cb( void* hwAttrs, SensorCallbackFxn fun );

/**
 * @fn      int sysfs_wrapper_init()
 * @brief   Wrapper to sysfs driver for Initialization of driver and device.
 *
 * @return  On success, 0
 *          On failure, non zero value
 */
int sysfs_wrapper_init ();

/**
 * @fn      int sysfs_wrapper_read(void*, void*, void*)
 * @brief   Wrapper to sysfs driver for reading from device.
 *
 * @param   hwAttrs
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, -1
 */
int sysfs_wrapper_read(void* hwAttrs, void* prop , void* data );

/**
 * @fn      int sysfs_wrapper_registration(Device*)
 * @brief   TBU
 *
 * @param   dev
 * @return  On success, 0
 *          On failure, non zero value
 */
int sysfs_wrapper_registration(Device* dev);

/**
 * @fn      int sysfs_wrapper_reg_cb(void*, SensorCallbackFxn)
 * @brief   register callback functions for the IRQ's for sensor device.
 *
 * @param   hwAttrs
 * @param   fun
 * @return  On success, 0
 *          On failure, non zero value
 */
int sysfs_wrapper_reg_cb( void* hwAttrs, SensorCallbackFxn fun );

/**
 * @fn      int sysfs_wrapper_write(void*, void*, void*)
 * @brief    Wrapper to sysfs driver for writing to device.
 *
 * @param   hwAttrs
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, -1
 */
int sysfs_wrapper_write(void* hwAttrs, void* prop , void* data);

/**
 * @fn      const DrvrOps sysfs_wrapper_get_ops*()
 * @brief   Function to get sysfs wrapper operation supported.
 *
 * @return  pointer to the struct with supported operations.
 */
const DrvrOps* sysfs_wrapper_get_ops();

#ifdef __cplusplus
}
#endif

#endif /* SYSFS_WRAPPER_H_ */
