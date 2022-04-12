/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef IDB_IDB_H_
#define IDB_IDB_H_

#include <stdio.h>
#include <stdint.h>

typedef int (*IDbInitFxn)(void* data);
typedef void (*IDbExitFxn)();
typedef int (*IDbParseFxn)(void* data);
typedef void* (*IDbReadHeaderFxn)(char* uuid, uint16_t* size);
typedef void* (*IDbReadIndexFxn)( char* uuid, uint16_t* size);
typedef void* (*IDbReadUnitInfoFxn)( char* uuid, uint16_t* size);
typedef void* (*IDbReadUnitConfigFxn)( char* uuid, uint16_t* size, uint8_t count);
typedef void* (*IDbReadModuleInfoFxn)( char* uuid, uint16_t* size, uint8_t idx);
typedef void* (*IDbReadModuleInfoByUUUIDFxn)( char *uuid, uint16_t* size, uint8_t count);
typedef void* (*IDbReadModuleConfigFxn)( char* uuid, uint16_t* size, uint8_t count);
typedef void* (*IDbReadFactConfigFxn)( char* uuid, uint16_t* size);
typedef void* (*IDbReadUserConfigFxn)( char* uuid, uint16_t* size);
typedef void* (*IDbReadFactCalibFxn)( char* uuid, uint16_t* size);
typedef void* (*IDbReadUserCalibFxn)( char* uuid, uint16_t* size);
typedef void* (*IDbReadBsCertsFxn)( char* uuid, uint16_t* size);
typedef void* (*IDbReadLwm2mCertsFxn)( char* uuid, uint16_t* size);

/*basic read write operation to UKDB*/
typedef struct  {
	IDbInitFxn init;
	IDbExitFxn  exit;
	IDbReadHeaderFxn read_header;
	IDbReadIndexFxn read_index;
	IDbReadUnitInfoFxn read_unit_info;
	IDbReadUnitConfigFxn read_unit_cfg;
	IDbReadModuleInfoFxn read_module_info;
	IDbReadModuleInfoByUUUIDFxn read_module_info_by_uuid;
	IDbReadModuleConfigFxn read_module_cfg;
	IDbReadFactConfigFxn read_fact_config;
	IDbReadUserConfigFxn read_user_config;
	IDbReadFactCalibFxn read_fact_calib;
	IDbReadUserCalibFxn read_user_calib;
	IDbReadBsCertsFxn read_bs_certs;
	IDbReadLwm2mCertsFxn read_lwm2m_certs;
} IDBFxnTable;

int idb_init(void* data);
void idb_exit();
int idb_parse(void* data);
int idb_fetch_header(void** data, char* uuid, uint16_t* size);
int idb_fetch_index(void** data, char* uuid, uint16_t* size);
int idb_fetch_unit_info(void** data, char* uuid, uint16_t* size);
int idb_fetch_unit_cfg(void** data, char* uuid, uint16_t* size, uint8_t count);
int idb_fetch_module_info(void** data, char* uuid, uint16_t* size, uint8_t idx);
int idb_fetch_module_info_by_uuid(void** data, char* uuid, uint16_t* size, uint8_t count);
int idb_fetch_module_cfg(void** data, char* uuid, uint16_t* size, uint8_t count);
int idb_fetch_fact_config(void** data, char* uuid, uint16_t* size);
int idb_fetch_user_config(void** data, char* uuid, uint16_t* size);
int idb_fetch_fact_calib(void** data, char* uuid, uint16_t* size);
int idb_fetch_user_calib(void** data, char* uuid, uint16_t* size);
int idb_fetch_bs_certs(void** data, char* uuid, uint16_t* size);
int idb_fetch_lwm2m_certs(void** data, char* uuid, uint16_t* size);
int idb_fetch_payload_from_mfgdata(void** data, char* uuid, uint16_t* size, uint16_t id);
#endif /* IDB_IDB_H_ */
