/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_MFG_H_
#define INC_MFG_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "schema.h"

#include "usys_types.h"

typedef int (*MfgInitFxn)(void* data);
typedef void (*MfgExitFxn)();
typedef int (*MfgParseFxn)(void* data);
typedef void* (*MfgReadHeaderFxn)(char* uuid, uint16_t* size);
typedef void* (*MfgReadIndexFxn)( char* uuid, uint16_t* size);
typedef void* (*MfgReadUnitInfoFxn)( char* uuid, uint16_t* size);
typedef void* (*MfgReadUnitConfigFxn)( char* uuid, uint16_t* size, uint8_t count);
typedef void* (*MfgReadModuleInfoFxn)( char* uuid, uint16_t* size, uint8_t idx);
typedef void* (*MfgReadModuleInfoByUUUIDFxn)( char *uuid, uint16_t* size, uint8_t count);
typedef void* (*MfgReadModuleConfigFxn)( char* uuid, uint16_t* size, uint8_t count);
typedef void* (*MfgReadFactConfigFxn)( char* uuid, uint16_t* size);
typedef void* (*MfgReadUserConfigFxn)( char* uuid, uint16_t* size);
typedef void* (*MfgReadFactCalibFxn)( char* uuid, uint16_t* size);
typedef void* (*MfgReadUserCalibFxn)( char* uuid, uint16_t* size);
typedef void* (*MfgReadBsCertsFxn)( char* uuid, uint16_t* size);
typedef void* (*MfgReadCloudCertsFxn)( char* uuid, uint16_t* size);

/* Basic operations needs to be supported by schema parser
 * As of now we only have JSON schema's but this could be extended to any
 * type and parser has to implement the below listed functions to be compatible
 * with inventory. Mfg provides abstracts type of parser required from higher
 * layers.
 */
typedef struct  {
  MfgInitFxn init;
  MfgExitFxn  exit;
  MfgReadHeaderFxn readHeader;
  MfgReadIndexFxn readIndex;
  MfgReadUnitInfoFxn readUnitInfo;
  MfgReadUnitConfigFxn readNodeCfg;
  MfgReadModuleInfoFxn readModuleInfo;
  MfgReadModuleInfoByUUUIDFxn readModuleInfoByUuid;
  MfgReadModuleConfigFxn readModuleCfg;
  MfgReadFactConfigFxn readFactCfg;
  MfgReadUserConfigFxn readUserCfg;
  MfgReadFactCalibFxn readFactCalib;
  MfgReadUserCalibFxn readUserCalib;
  MfgReadBsCertsFxn readBsCerts;
  MfgReadCloudCertsFxn readCloudCerts;
} MfgOperations;

/**
 * @fn      int mfg_init(void*)
 * @brief   Abstracts parser initialization.
 *
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_init(void* data);

/**
 * @fn      int mfg_fetch_bs_certs(void**, char*, uint16_t*)
 * @brief   Abstracts fetching of bootstrap certificates.
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_bs_certs(void** data, char* uuid, uint16_t* size);

/**
 * @fn      int mfg_fetch_cloud_certs(void**, char*, uint16_t*)
 * @brief   Abstracts fetching of cloud certificates.
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_cloud_certs(void** data, char* uuid, uint16_t* size);

/**
 * @fn      int mfg_fetch_fact_calib(void**, char*, uint16_t*)
 * @brief    Abstracts fetching of factory calibaration.
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_fact_calib(void** data, char* uuid, uint16_t* size);

/**
 * @fn      int mfg_fetch_fact_cfg(void**, char*, uint16_t*)
 * @brief    Abstracts fetching of factory configuration.
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_fact_cfg(void** data, char* uuid, uint16_t* size);

/**
 * @fn      int mfg_fetch_header(void**, char*, uint16_t*)
 * @brief    Abstracts fetching of header.
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_header(SchemaHeader** data, char* uuid, uint16_t* size);

/**
 * @fn      int mfg_fetch_idx(void**, char*, uint16_t*)
 * @brief    Abstracts fetching of index table.
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_idx(SchemaIdxTuple** data, char* uuid, uint16_t* size);

/**
 * @fn      int mfg_fetch_module_info(void**, char*, uint16_t*, uint8_t)
 * @brief   Abstracts fetching of module info from module with uuid.
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @param   idx
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_module_info(ModuleInfo** data, char* uuid, uint16_t* size, uint8_t idx);

/**
 * @fn      int mfg_fetch_module_info_by_uuid(void**, char*, uint16_t*, uint8_t)
 * @brief   Abstracts fetching of module info with uuid.
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_module_info_by_uuid(ModuleInfo** data, char* uuid, uint16_t* size, uint8_t count);

/**
 * @fn      int mfg_fetch_module_cfg(void**, char*, uint16_t*, uint8_t)
 * @brief   Abstracts fetching of module config
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_module_cfg(void** data, char* uuid, uint16_t* size, uint8_t count);

/**
 * @fn      int mfg_fetch_payload_from_mfg_data(void**, char*, uint16_t*, uint16_t)
 * @brief   Abstracts fetching of manufacturing payloads.
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @param   id
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_payload_from_mfg_data(void** data, char* uuid, uint16_t* size, uint16_t id);

/**
 * @fn      int mfg_fetch_node_info(void**, char*, uint16_t*)
 * @brief   Abstracts fetching of unit info
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_node_info(NodeInfo** data, char* uuid, uint16_t* size);

/**
 * @fn      int mfg_fetch_node_cfg(void**, char*, uint16_t*, uint8_t)
 * @brief   Abstracts fetching of unit config
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @param   count
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_node_cfg(NodeCfg** data, char* uuid, uint16_t* size, uint8_t count);

/**
 * @fn      int mfg_fetch_user_calib(void**, char*, uint16_t*)
 * @brief   Abstracts fetching of user calibration data
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_user_calib(void** data, char* uuid, uint16_t* size);

/**
 * @fn      int mfg_fetch_user_cfg(void**, char*, uint16_t*)
 * @brief   Abstracts fetching of user configuration data
 *
 * @param   data
 * @param   uuid
 * @param   size
 * @return  On success, 0
 *          On failure, non zero value
 */
int mfg_fetch_user_cfg(void** data, char* uuid, uint16_t* size);

/**
 * @fn      void mfg_exit()
 * @brief   Wrapper to free all memory used by parser.
 *
 */
void mfg_exit();

#ifdef __cplusplus
}
#endif

#endif /* INC_MFG_H_ */
