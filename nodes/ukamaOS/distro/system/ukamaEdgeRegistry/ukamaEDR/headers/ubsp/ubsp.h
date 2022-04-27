/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef HEADERS_UBSP_H_
#define HEADERS_UBSP_H_

#include "headers/errorcode.h"
#include "headers/ubsp/devices.h"
#include "headers/ubsp/property.h"
#include "headers/ubsp/ukdblayout.h"

#include <stdint.h>


/*
 * For Usage of these functions have a look at test.c
 * Each function is having one or more test related to it.
 */

/******************************************************************************
 *
 * Description:
 *      Allocates memory of size mentioned by argument in bytes.
 * Args:
 *       size_t               : size in bytes
 *
 * Return Value:
 *      void*                 : On Success
 *      NULL                  : On Failure.
 *
 *****************************************************************************/
void* ubsp_alloc(size_t size);

/******************************************************************************
 *
 * Description:
 *      Free memory.
 * Args:
 *      void*                : Memory
 * Return Value:
 *
 *****************************************************************************/
void ubsp_free(void* mem);

/******************************************************************************
 *
 * Description:
 *      Allocates memory for unit config.
 * Args:
 *uint8_t                    : Device count
 *
 * Return Value:
 *      UnitCfg*             : On Success
 *      NULL                 : On Failure.
 *
 *****************************************************************************/
UnitCfg* ubsp_alloc_unit_cfg(uint8_t module_count);

/******************************************************************************
 *
 * Description:
 *      Free memory for Unit config.
 * Args:
 *   UnitCfg*                : Pointer to Unit Cfg memory
 *   Uint8_t                 : Module count in Unit.
 * Return Value:
 *
 *****************************************************************************/
void ubsp_free_unit_cfg(UnitCfg *cfg, uint8_t module_count);

/******************************************************************************
 *
 * Description:
 *      Allocates memory for Module config.
 * Args:
 *      uint8_t             : Device count
 *
 * Return Value:
 *      ModuleCfg*           : On Success
 *      NULL                 : On Failure.
 *
 *****************************************************************************/
ModuleCfg* ubsp_alloc_module_cfg( uint8_t dev_count);

/******************************************************************************
 *
 * Description:
 *      Free memory for Module config.
 * Args:
 *   ModuleCfg*           :  Pointer to Module cfg memory.
 *   Uint8_t              : Device count in module.
 * Return Value:
 *
 *****************************************************************************/
void ubsp_free_module_cfg(ModuleCfg *cfg, uint8_t dev_count);

/******************************************************************************
 *
 * Description:
 * 		Initializes ubsp lists for ModuleDb.
 * Args:
 *      char*               : Pointer to database location in file system.
 *                            Default argument path can be /systemdb
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_ukdb_init(char* sys_db_path);

/******************************************************************************
 *
 * Description:
 *      Initializes ubsp lists for deviceDb and irqDb
 * Args:
 *      void*               : Typically it has to be JSONInput.
 * Return Value:
 *      0                   : On Success
 *      Negative values     : On Failure.
 *
 *****************************************************************************/
int ubsp_devdb_init(void *data);

/******************************************************************************
 *
 * Description:
 *      Initializes ubsp mfg specfic parameter for DB creation
 * Args:
 *      JSONInput            : structure of list of file name and a count.
 * Return Value:
 *      0                   : On Success
 *      Negative values     : On Failure.
 *
 *****************************************************************************/
int ubsp_idb_init(void* data);

/******************************************************************************
 *
 * Description:
 *      Clears the structs allocated for storing mfg data.
 * Args:
 *
 * Return Value:
 *
 *****************************************************************************/
void ubsp_idb_exit();

/******************************************************************************
 *
 * Description:
 *      Remove ubsp lists for deviceDb, ModuleDb and irqDb
 * Args:
 *
 * Return Value:
 *      0                   : On Success
 *      Negative values     : On Failure.
 *
 *****************************************************************************/
int ubsp_exit();

/******************************************************************************
 *
 * Description:
 *      Register module with the UUID.
 * Args:
 *      UnitCfg:            : Unit Config for Module UUID.
 *
 * Return Value:
 *      0                   : On Success
 *      Negative values     : On Failure.
 *
 *****************************************************************************/
int ubsp_register_module(UnitCfg* cfg);

/******************************************************************************
 *
 * Description:
 * 		Read the EEPROM header for the Module UUID
 *
 * Args:
 *		UUID				: Module UUID.
 *		Header				: A pointer to header info.
 *
 * Return Value:
 * 		0 					: On Success
 * 	     Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_header(char* puuid, UKDBHeader* pheader);
//TODO: Try to allocate memory and return the pointer in UKDB similar to unit info.

/******************************************************************************
 *
 * Description:
 * 		Validate the Magic work for the Ukama Db.
 *
 * Args:
 *		UUID				: Module UUID.
 *		Version				: A pointer to Version info.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_validating_magicword(char* puuid);

/******************************************************************************
 *
 * Description:
 * 		Read the Ukama Db version for the Module UUID
 *
 * Args:
 *		UUID				: Module UUID.
 *		Version				: A pointer to Version info.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_dbversion(char* puuid, Version* pver);

/******************************************************************************
 *
 * Description:
 * 		Updates the Ukama Db version for the Module UUID
 *
 * Args:
 *		UUID				: Module UUID.
 *		Version				: Version info struct.
 *
 * Return Value:
 * 		0 					: On Success
 * 	     Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_update_dbversion(char* puuid, Version ver);

/******************************************************************************
 *
 * Description:
 * 		Read Unit Information stored in the Module UIID
 *
 * Args:
 *		UUID				: Module UUID.
 *		UnitInfo			: A pointer to Unit Info.
 * 		Size				: A pointer to size of Unit Info.
 *
 * Return Value:
 * 		0 					: On Success
 * 	     Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_unit_info(char* puuid, UnitInfo* pdata, uint16_t* psize);

/******************************************************************************
 *
 * Description:
 * 		Read Unit Configuration stored in the Module UIID
 *
 * Args:
 * 		UIID				: Module UUID.
 * 		UnitCfg				: Unit Config
 * 		Count				: Number of module in the Unit.Read from unit info.
 * 		Size				: A pointer to size of Unit Config.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_unit_cfg(char* puuid, UnitCfg* pucfg, uint8_t count,
		uint16_t* psize);

/******************************************************************************
 *
 * Description:
 * 		Read Module Information for the Module UIID
 *
 * Args:
 *		UUID				: Module UUID.
 *		ModuleInfo			: A pointer to Module Info.
 * 		Size				: A pointer to size of Module Info.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_module_info(char* puuid, ModuleInfo* pinfo, uint16_t* psize);

/******************************************************************************
 *
 * Description:
 * 		Read Module Configuration for the Module UIID
 *
 * Args:
 * 		UIID				: Module UUID.
 * 		UnitCfg				: Module Config
 * 		Count				: Number of devices under the Module.
 * 							  Read from module info.
 * 		Size				: A pointer to size of Module Config.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_module_cfg(char* puuid, ModuleCfg* pcfg, uint8_t count,
		uint16_t* psize);

/******************************************************************************
 *
 * Description:
 * 		Create a unit schema from Unit Info and Unit Config.
 *
 * Args:
 *		UnitInfo*			: Unit Info.
 *		UnitCfg*			: Unit config.
 * 		char*				: A string to json data is returned.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_create_unit_schema(UnitInfo *unit_info, UnitCfg *unit_cfg, char* junit_schema);

/******************************************************************************
 *
 * Description:
 * 		Create a module schema from Module Info and Module Config.
 *
 * Args:
 *		ModuleInfo*			: Module Info.
 * 		char*				: A string to json data is returned.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_create_module_schema( ModuleInfo* pinfo, char* jmod_schema);

/******************************************************************************
 *
 * Description:
 * 		Create a schema from Unit with Unit info, unit config, module info and module Config.
 *
 * Args:
 *      UnitInfo*			: Unit Info.
 *		UnitCfg*			: Unit config.
 *		ModuleInfo*			: Module Info.
 * 		char**				: A string to json data is returned.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_create_schema( UnitInfo *unit_info, UnitCfg *unit_cfg,
		ModuleInfo* minfo, char** junit_schema);

/******************************************************************************
 *
 * Description:
 * 		Read factory config for the Module UIID
 *
 * Args:
 * 		UIID				: Module UUID.
 * 		Data				: void pointer
 * 		Size				: A pointer to size of factory config.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_fact_config(char* puuid, void* pdata, uint16_t* psize);

/******************************************************************************
 *
 * Description:
 * 		Read user config for the Module UIID
 *
 * Args:
 * 		UIID				: Module UUID.
 * 		Data				: void pointer
 * 		Size				: A pointer to size of user config.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_user_config(char* puuid, void* pdata, uint16_t* psize);

/******************************************************************************
 *
 * Description:
 * 		Read factory calibration for the Module UIID
 *
 * Args:
 * 		UIID				: Module UUID.
 * 		Data				: void pointer
 * 		Size				: A pointer to size of factory calibration.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_fact_calib(char* puuid, void* pdata, uint16_t* psize);

/******************************************************************************
 *
 * Description:
 * 		Read user calibration for the Module UIID
 *
 * Args:
 * 		UIID				: Module UUID.
 * 		Data				: void pointer
 * 		Size				: A pointer to size of user calibration.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_user_calib(char* puuid, void* pdata, uint16_t* psize);

/******************************************************************************
 *
 * Description:
 * 		Read bootstrap certs for the Module UIID
 *
 * Args:
 * 		UIID				: Module UUID.
 * 		Data				: void pointer
 * 		Size				: A pointer to size of bootstrap certs.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_bs_certs(char* puuid, void* pdata, uint16_t* psize);


/******************************************************************************
 *
 * Description:
 * 		Read lwm2m certs for the Module UIID
 *
 * Args:
 * 		UIID				: Module UUID.
 * 		Data				: void pointer
 * 		Size				: A pointer to size of lwm2m certs.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_lwm2m_certs(char* puuid, void* pdata, uint16_t* psize);

/******************************************************************************
 *
 * Description:
 *      Used in test scenarios to create a module DB for slave modules.
 *
 * Args:
 *      UUID                : Module UUID.
 *
 * Return Value:
 *      0                   : On Success
 *      Negative values     : On Failure.
 *
 *****************************************************************************/
int ubsp_pre_create_ukdb_hook(char* mod_uuid);

/******************************************************************************
 *
 * Description:
 *      Create a UKDB for module. Used on targets. This will enumerate modules
 *      from the master module DB and register devices to device DB.
 *
 * Args:
 *      UUID                : Module UUID.
 *
 * Return Value:
 *      0                   : On Success
 *      Negative values     : On Failure.
 *
 *****************************************************************************/
int ubsp_create_ukdb(char* mod_uuid);

/******************************************************************************
 *
 * Description:
 *      Delete DB for the module.
 *
 * Args:
 *      UUID                : Module UUID.
 *
 * Return Value:
 *      0                   : On Success
 *      Negative values     : On Failure.
 *
 *****************************************************************************/
int ubsp_remove_ukdb(char* mod_uuid);

/******************************************************************************
 *
 * Description:
 * 		Read the count of registered devices under mentioned device type
 *
 * Args:
 * 		Type				: DEVICE TYPE.
 * 		Count				: A pointer to number of devices registered
 * 							  for device the provided device type.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_registered_dev_count(DeviceType type, uint16_t* count);

/******************************************************************************
 *
 * Description:
 * 		Read the registered devices info under mentioned device type
 *
 * Args:
 * 		Type				: DEVICE TYPE.
 * 		Device				: A pointer to array of registered device
 * 							  the provided device type.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_registered_dev(DeviceType type, Device* dev );


/******************************************************************************
 *
 * Description:
 * 		Read the count of the device property for particular device represented
 * 		by device object
 *
 * Args:
 * 		DevObj				: A pointerDevice Object.
 * 		Count				: A pointer to number of devices registered
 * 							  for device the provided device type.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int upsb_read_dev_prop_count(DevObj* obj, uint16_t* count);

/******************************************************************************
 *
 * Description:
 * 		Read the device properties for particular device represented by
 * 		device object
 *
 * Args:
 * 		DevObj				: A pointer to device object .
 * 		Property			: A pointer to array of properties.
 * 							  the provided device type.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_dev_props(DevObj* obj, Property* prop );


/******************************************************************************
 *
 * Description:
 * 		Read the property for particular device represented by
 * 		device object and property arguments in function.
 *
 * Args:
 * 		DevObj				: A pointer to device object .
 * 		int					: A pointer property index.
 *		data				: A pointer to data of type void.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_read_from_prop(DevObj* obj, int* prop, void* data);

/******************************************************************************
 *
 * Description:
 * 		Write the property for particular device represented by
 * 		device object and property arguments in function.
 *
 * Args:
 * 		DevObj				: A pointer to device object .
 * 		int					: A pointer property index.
 *		data				: A pointer to data of type void.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_write_to_prop(DevObj* obj, int* prop, void* data);

/******************************************************************************
 * TODO:
 * Description:
 * 		Enable device mentioned in arguments of the function.
 *
 * Args:
 * 		DevObj				: A pointer to device object .
 *		data				: A void pinter.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_enable(DevObj* obj, void* data);

/******************************************************************************
 * TODO:
 * Description:
 * 		Disable device mentioned in arguments of the function.
 *
 * Args:
 * 		DevObj				: A pointer to device object .
 *		data				: A void pinter.
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_disable(DevObj* obj, void* data);


/******************************************************************************
 * Description:
 * 		Enable the IRQ thread for all the device alerts which are configured.
 *
 * Args:
 * 		DevObj				: A pointer to device object.
 * 		int					: A pointer to property index  from the property array
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_enable_irq(DevObj* obj, int* idx);

/******************************************************************************
 * Description:
 * 		Disable the IRQ thread for all the device alerts which are configured.
 *
 * Args:
 * 		DevObj				: A pointer to device object .
 * 		int					: A pointer to property index from the property array
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_disable_irq(DevObj* obj, int* idx);


/******************************************************************************
 * Description:
 * 		Register the application callback for a device. Callback is per device type.
 *		All temp sensors will have one callback.
 *
 * Args:
 * 		DevObj				: A pointer to device object .
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_register_app_cb(DevObj* obj, void* prop, CallBackFxn fn);

/******************************************************************************
 * Description:
 * 		De-register the application callback for a device. Callback is per device type.
 *		All temp sensors will have one callback.
 *
 * Args:
 * 		DevObj				: A pointer to device object .
 *
 * Return Value:
 * 		0 					: On Success
 * 	    Negative values 	: On Failure.
 *
 *****************************************************************************/
int ubsp_deregister_app_cb(DevObj* obj, void* prop, CallBackFxn fn);

#endif /* HEADERS_UBSP_H_ */
