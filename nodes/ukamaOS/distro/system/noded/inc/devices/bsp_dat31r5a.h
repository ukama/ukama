/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */


#ifndef BSP_DAT31R5A_H_
#define BSP_DAT31R5A_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "device.h"

/**
 * @fn      int bsp_dat31r5a_disable(void*, void*, void*)
 * @brief   passes sensor disable request to driver wrapper layer
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_dat31r5a_disable(void* p_dev, void* prop, void* data);

/**
 * @fn      int bsp_dat31r5a_enable(void*, void*, void*)
 * @brief   passes sensor enable request to driver wrapper layer
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_dat31r5a_enable(void* p_dev, void* prop, void* data);

/**
 * @fn      int bsp_dat31r5a_configure(void*, void*, void*)
 * @brief   passes configuration request to driver wrapper layer.
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_dat31r5a_configure(void* p_dev, void* prop, void* data );

/**
 * @fn      int bsp_dat31r5a_init()
 * @brief   Read sensor properties from the property config parsed during
 *          startup and then initialization request to driver wrapper layer.
 *
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_dat31r5a_init ();

/**
 * @fn      int bsp_dat31r5a_registration(Device*)
 * @brief   passes registration request for driver wrapper layer.
 *
 * @param   p_dev
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_dat31r5a_registration(Device* p_dev);

/**
 * @fn      int bsp_dat31r5a_read_prop_count(Device*, uint16_t*)
 * @brief   Reads the number properties available for sensor.
 *
 * @param   p_dev
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_dat31r5a_read_prop_count(Device* p_dev, uint16_t * count);

/**
 * @fn      int bsp_dat31r5a_read_properties(Device*, void*)
 * @brief   Reads the properties available for sensor.
 *
 * @param   p_dev
 * @param   prop
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_dat31r5a_read_properties(Device* p_dev, void* prop);

/**
 * @fn      int bsp_dat31r5a_read(void*, void*, void*)
 * @brief   passes read request to driver wrapper layer.
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_dat31r5a_read(void* p_dev, void* prop, void* data);

/**
 * @fn      int bsp_dat31r5a_write(void*, void*, void*)
 * @brief   passes write request to driver wrapper layer.
 *
 * @param   p_dev
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int bsp_dat31r5a_write(void* p_dev, void* prop, void* data);

#ifdef __cplusplus
}
#endif

#endif /* BSP_DAT31R5A_H_ */
