/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_LEDGER_H_
#define INC_LEDGER_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "device.h"
#include "schema.h"

/**
 * @fn      int (*)(char*, char*, uint8_t, ModuleCfg*)
 * @brief   Alert callback function for the device.
 *
 * @param   p_uuid
 * @param   name
 * @param   count
 * @param   p_mcfg
 * @return  On success, 0
 *          On failure, non zero value
 */
typedef int (*RegisterDeviceCB)(char *p_uuid, char *name, uint8_t count, ModuleCfg *p_mcfg);

/**
 * @fn      int ldgr_configure(DevObj*, void*, void*)
 * @brief   TBU: Could be used for User config, default config and factory
 *          config where property can be a list of all properties and same goes
 *          for data which contain the value to be set for those.
 *
 * @param   obj
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_configure(DevObj* obj, void* prop, void* data);

/**
 * @fn      int ldgr_dereg_app_cb(DevObj*, void*, CallBackFxn)
 * @brief   De-register the application callback for a sensor.
 *
 * @param   obj
 * @param   prop
 * @param   fn
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_dereg_app_cb(DevObj* obj, void* prop, CallBackFxn fn);

/**
 * @fn      int ldgr_disable(DevObj*, void*, void*)
 * @brief   Disable the sensor.
 *
 * @param   obj
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_disable(DevObj* obj, void* prop, void* data);


/**
 * @fn      int ldgr_disable_irq(DevObj*, void*, void*)
 * @brief   Disable  the IRQ thread for all the sensor alerts which are
 *          configured.
 *
 * @param   obj
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_disable_irq(DevObj* obj, void* prop, void* data);

/* Enable the sensor. */
/**
 * @fn      int ldgr_enable(DevObj*, void*, void*)
 * @brief
 *
 * @param   obj
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_enable(DevObj* obj, void* prop, void* data);

/**
 * @fn      int ldgr_enable_irq(DevObj*, void*, void*)
 * @brief   Enable the IRQ thread for all the sensor alerts which are
 *          configured.
 *
 * @param   obj
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_enable_irq(DevObj* obj, void* prop, void* data);

/**
 * @fn      void ldgr_exit()
 * @brief   Remove ledger for each sensor class.
 *          Sensor classes are like ADC, TMP, POWER, GPIO, LED's
 *
 */
void ldgr_exit();

/**
 * @fn      int ldgr_init(void*)
 * @brief   Creates ledger for each sensor class.
 *          Sensor classes are like ADC, TMP, POWER, GPIO, LED's
 *
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_init(void *data);

/**
 * @fn      int ldgr_read(DevObj*, void*, void*)
 * @brief   Read the property value from sensor
 *
 * @param   obj
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_read(DevObj* obj, void* prop, void* data);

/* .*/
/**
 * @fn      int ldgr_read_prop(DevObj*, void*)
 * @brief   Read the sensor properties for particular sensor represented
 *          by device object from property config.
 *
 * @param   obj
 * @param   prop
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_read_prop(DevObj* obj, void* prop );

/**
 * @fn      int ldgr_read_prop_count(DevObj*, uint16_t*)
 * @brief   Read the sensor count of properties for particular sensor
 *          represented by device object from property config.
 *
 * @param   obj
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_read_prop_count(DevObj* obj, uint16_t* count);

/**
 * @fn      int ldgr_read_reg_dev(DeviceType, Device*)
 * @brief   Read the list of registered devices info under Sensor class
 *          mentioned by type
 *
 * @param   type
 * @param   dev
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_read_reg_dev(DeviceType type, Device* dev );

/**
 * @fn      int ldgr_read_reg_dev_count(DeviceType, uint16_t*)
 * @brief   Read the count of registered devices under Sensor class mentioned
 *          by type.
 *
 * @param   type
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_read_reg_dev_count(DeviceType type, uint16_t* count);

/**
 * @fn      int ldgr_register(char*, char*, uint8_t, ModuleCfg*)
 * @brief   Register individual sensors to their respective classes.
 *          like TMP464, SE98 and ADT to TMP class.
 *          INA226 to power class.
 *          ADS1015 to ADC  class.
 *
 * @param   p_uuid
 * @param   name
 * @param   count
 * @param   p_mcfg
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_register(char* p_uuid, char* name, uint8_t count, ModuleCfg* p_mcfg);

/**
 * @fn      int ldgr_reg_app_cb(DevObj*, void*, CallBackFxn)
 * @brief   Register the application callback for a sensor.
 *
 * @param   obj
 * @param   prop
 * @param   fn
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_reg_app_cb(DevObj* obj, void* prop, CallBackFxn fn);

/**
 * @fn      int ldgr_write(DevObj*, void*, void*)
 * @brief   Write the property value to sensor hardware.
 *
 * @param   obj
 * @param   prop
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int ldgr_write(DevObj* obj, void* prop, void* data);

#ifdef __cplusplus
}
#endif

#endif /* INC_LEDGER_H_ */
