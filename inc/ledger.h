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

typedef int (*RegisterDeviceCB)(char *p_uuid, char *name, uint8_t count, ModuleCfg *p_mcfg);

 /* Initialize DataBase for each sensor class.
 * Sensor classes are like ADC, TMP, POWER, GPIO, LED's
 * */
int ldgr_init(void *data);

/* Remove DataBase for each sensor class.
* Sensor classes are like ADC, TMP, POWER, GPIO, LED's
* */
void ldgr_exit();

/*  Register individual sensors to their respective classes.
 *  like TMP464, SE98 and ADT to TMP class.
 *  INA226 to power class.
 *  ADS1015 to ADC  class.
 *  */
int ldgr_register(char* p_uuid, char* name, uint8_t count, ModuleCfg* p_mcfg);

/* Read the count of registered devices under Sensor class mentioned by type.*/
int ldgr_read_reg_dev_count(DeviceType type, uint16_t* count);

/* Read the registered devices info under Sensor class mentioned by type.*/
int ldgr_read_reg_dev(DeviceType type, Device* dev );

/* Read the sensor property for particular sensor represented by device object.*/
int ldgr_read_prop(DevObj* obj, void* prop );

/* Read the count of the sensor property for particular sensor represented by device object.*/
int ldgr_read_prop_count(DevObj* obj, uint16_t* count);

/* TODO: Read write properties will do most of the work. Not sure if this is required.
 * Could be used for User config, default config and factory config where
 * property can be a list of all properties and same goes for data
 * which contain the value to be set for those.
 */
int ldgr_configure(DevObj* obj, void* prop, void* data);

/* Read the property for sensor */
int ldgr_read(DevObj* obj, void* prop, void* data);

/* Write the property for sensor */
int ldgr_write(DevObj* obj, void* prop, void* data);

/* Enable the sensor. */
int ldgr_enable(DevObj* obj, void* prop, void* data);

/* Disable the sensor. */
int ldgr_disable(DevObj* obj, void* prop, void* data);

/* Enable the IRQ thread for all the sensor alerts which are configured. */
int ldgr_enable_irq(DevObj* obj, void* prop, void* data);

/* Disable  the IRQ thread for all the sensor alerts which are configured. */
int ldgr_disable_irq(DevObj* obj, void* prop, void* data);

/* Register the application callback for a sensor. */
int ldgr_reg_app_cb(DevObj* obj, void* prop, CallBackFxn fn);

/* De-register the application callback for a sensor. */
int ldgr_dereg_app_cb(DevObj* obj, void* prop, CallBackFxn fn);

#ifdef __cplusplus
}
#endif

#endif /* INC_LEDGER_H_ */
