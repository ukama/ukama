/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_DEVHELPER_H_
#define INC_DEVHELPER_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "driver_ops.h"

/**
 * @fn      int dhelper_configure(const DrvrOps*, Device*, Property*, int, void*)
 * @brief   Checks for the property if it's valid or not and then reads
 *          the sensor specific hw attributes and pass that info to the
 *          driver layer selected by bsp layer for further configuration.
 *
 * @param   drvr
 * @param   dev
 * @param   prop
 * @param   pidx
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_configure(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data);

/**
 * @fn      int dhelper_confirm_irq(const DrvrOps*, Device*, Property*, AlertCallBackData**, char*, int, int*)
 * @brief   Calls irqhelper to confirm IRQ reported recently.
 *
 * @param   drvr
 * @param   dev
 * @param   prop
 * @param   acbdata
 * @param   fpath
 * @param   maxpcount
 * @param   evt
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_confirm_irq(const DrvrOps *drvr,Device *dev, Property *prop, AlertCallBackData **acbdata, char *fpath,
                        int maxpcount, int *evt);

/**
 * @fn      int dhelper_disable(const DrvrOps*, Device*, Property*, int, void*)
 * @brief   Checks for the property if it's valid or not and then reads
 *          the sensor specific hw attributes and pass that info to the
 *          driver layer selected by bsp layer for disabling sensor.
 *
 * @param   drvr
 * @param   dev
 * @param   prop
 * @param   pidx
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_disable(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data);

/**
 * @fn      int dhelper_disable_irq(const DrvrOps*, Device*, Property*,
 *          int, void*)
 * @brief   Checks for the property if it's valid or not and then reads
 *          the sensor specific hw attributes and pass that info to the
 *          driver layer selected by bsp layer for further disabling IRQ.
 *
 * @param   drvr
 * @param   dev
 * @param   prop
 * @param   pidx
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_disable_irq(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data);

/**
 * @fn      int dhelper_enable(const DrvrOps*, Device*, Property*, int, void*)
 * @brief   Checks for the property if it's valid or not and then reads
 *          the sensor specific hw attributes and pass that info to the
 *          driver layer selected by bsp layer for enabling sensor.
 *
 * @param   drvr
 * @param   dev
 * @param   prop
 * @param   pidx
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_enable(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data);

/**
 * @fn      int dhelper_enable_irq(const DrvrOps*, SensorCallbackFxn, Device*, Property*, int, void*)
 * @brief   Checks for the property if it's valid or not and then reads
 *          the sensor specific hw attributes and pass that info to the
 *          driver layer selected by bsp layer for further enabling IRQ.
 *
 * @param   drvr
 * @param   sensor_cb
 * @param   dev
 * @param   prop
 * @param   pidx
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_enable_irq(const DrvrOps *drvr, SensorCallbackFxn sensor_cb,
    Device *dev, Property *prop, int pidx, void *data);

/**
 * @fn      int dhelper_init_driver(const DrvrOps*, Device*)
 * @brief    Checks for the property if it's valid or not and then reads
 *          the sensor specific hw attributes and pass that info to the
 *          driver layer selected by bsp layer for further initialization.
 *
 * @param   drvr
 * @param   dev
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_init_driver(const DrvrOps *drvr, Device *dev);

/**
 * @fn      int dhelper_init_property_from_parser(Device*, Property**, int*)
 * @brief   request the property read by parser for sensor device during startup
 *          process.
 *
 * @param   p_dev
 * @param   prop
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_init_property_from_parser(Device *p_dev, Property** prop,
    int* count);

/**
 * @fn      int dhelper_read(const DrvrOps*, Device*, Property*, int, void*)
 * @brief   Checks for the property if it's valid or not and then reads
 *          the sensor specific hw attributes and pass that info to the
 *          driver layer selected by bsp layer for further configuration.
 *
 * @param   drvr
 * @param   dev
 * @param   prop
 * @param   pidx
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_read(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data) ;
/**
 * @fn      int dhelper_registration(const DrvrOps*, Device*)
 * @brief   Checks for the property if it's valid or not and then reads
 *          the sensor specific hw attributes and pass that info to the
 *          driver layer selected by bsp layer for further configuration.
 *
 * @param   drvr
 * @param   p_dev
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_registration(const DrvrOps *drvr, Device *p_dev);

/**
 * @fn      int dhelper_write(const DrvrOps*, Device*, Property*, int, void*)
 * @brief   Checks for the property if it's valid or not and then reads
 *          the sensor specific hw attributes and pass that info to the
 *          driver layer selected by bsp layer for performing write operation.
 *
 * @param   drvr
 * @param   dev
 * @param   prop
 * @param   pidx
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int dhelper_write(const DrvrOps *drvr, Device *dev, Property *prop,
    int pidx, void *data);

/**
 * @fn      int dhelper_validate_property(Property*, int)
 * @brief   Validates if the property is valid or not by checking if it's
 *          marked available in property config.
 *
 * @param   prop
 * @param   pidx
 * @return  On success, 0
 *          On failure, -1
 */
int dhelper_validate_property(Property* prop, int pidx);

/**
 * @fn      int dhelper_validate_property_type_alert(Property*, int)
 * @brief   validate if property type is alert or not.
 *
 * @param   prop
 * @param   pidx
 * @return  On success, 0
 *          On failure, -1
 */
int dhelper_validate_property_type_alert(Property *prop, int pidx);

/**
 * @fn      int dhelper_validate_permissions(Property*, int)
 * @brief   Validates if the property has required permissions or not.
 * @param   prop
 * @param   pidx
 * @return  On success, 0
 *          On failure, -1
 */
int dhelper_validate_permissions(Property *prop, int pidx, uint16_t perm);
/**
 * @fn      void dhelper_irq_callback(DevObj*, void*, void*)
 * @brief   TBU intention is to set IRQ callback
 *
 * @param   obj
 * @param   prop
 * @param   data
 */
void dhelper_irq_callback(DevObj *obj, void *prop, void *data);

#ifdef __cplusplus
}
#endif

#endif /* INC_DEVHELPER_H_ */
