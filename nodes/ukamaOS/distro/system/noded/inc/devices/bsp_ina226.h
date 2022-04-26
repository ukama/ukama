/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */


#ifndef BSP_INA226_H_
#define BSP_INA226_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "device.h"

/**
 * @fn      int bsp_ina226_configure(void*, void*, void*)
 * @brief   passes configuration request to driver wrapper layer.
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_configure(void* p_dev, void* prop, void* data );

/**
 * @fn      int bsp_ina226_confirm_irq(Device*, AlertCallBackData**, char*, int*)
 * @brief   passes request to confirm IRQ to IRQ helper.
 *
 * @param   dev
 * @param   acbdata
 * @param   fpath
 * @param   evt
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_confirm_irq(Device *dev, AlertCallBackData** acbdata,
    char* fpath, int* evt);
/**
 * @fn      int bsp_ina226_disable(void*, void*, void*)
 * @brief   passes sensor disable request to driver wrapper layer
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_disable(void* p_dev, void* prop, void* data);

/**
 * @fn      int bsp_ina226_disable_irq(void*, void*, void*)
 * @brief   passes disable interrupts for sensor to IRQ helper.
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_disable_irq(void* p_dev, void* prop, void* data);

/**
 * @fn      int bsp_ina226_dreg_cb(void*, SensorCallbackFxn)
 * @brief   de-register IRQ call back function
 *
 * @param   p_dev
 * @param   fun
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_dreg_cb(void* p_dev, SensorCallbackFxn fun);

/**
 * @fn      int bsp_ina226_enable(void*, void*, void*)
 * @brief   passes sensor enable request to driver wrapper layer
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_enable(void* p_dev, void* prop, void* data);

/**
 * @fn      int bsp_ina226_enable_irq(void*, void*, void*)
 * @brief   passes enable interrupts for the sensor to IRQ helper
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_enable_irq(void* p_dev, void* prop, void* data);

/**
 * @fn      int bsp_ina226_get_irq_type(int, uint8_t*)
 * @brief   gets the type of IRQ reported by sensor.
 *
 * @param   bsp_pidx
 * @param   alertstate
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_get_irq_type(int bsp_pidx, uint8_t* alertstate);

/**
 * @fn      int bsp_ina226_init()
 * @brief   Read sensor properties from the property config parsed during
 *          startup and then initialization request to driver wrapper layer.
 *
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_init ();

/**
 * @fn      int bsp_ina226_read(void*, void*, void*)
 * @brief   passes read request to driver wrapper layer.
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_read(void* p_dev, void* prop, void* data);

/**
 * @fn      int bsp_ina226_read_properties(Device*, void*)
 * @brief   Reads the properties available for sensor.
 *
 * @param   p_dev
 * @param   prop
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_read_properties(Device* p_dev, void* prop);

/**
 * @fn      int bsp_ina226_read_prop_count(Device*, uint16_t*)
 * @brief   Reads the number properties available for sensor.
 *
 * @param   p_dev
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_read_prop_count(Device* p_dev, uint16_t * count);

/**
 * @fn      int bsp_ina226_registration(Device*)
 * @brief   passes registration request for driver wrapper layer.
 *
 * @param   p_dev
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_registration(Device* p_dev);

/**
 * @fn      int bsp_ina226_reg_cb(void*, SensorCallbackFxn)
 * @brief   register IRQ call back function
 *
 * @param   p_dev
 * @param   fun
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_reg_cb(void* p_dev, SensorCallbackFxn fun);

/**
 * @fn      int bsp_ina226_write(void*, void*, void*)
 * @brief   passes write request to driver wrapper layer.
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_ina226_write(void* p_dev, void* prop, void* data);

#ifdef __cplusplus
}
#endif

#endif /* BSP_INA226_H_ */
