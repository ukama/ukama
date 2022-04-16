/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_JDATA_H_
#define INC_JDATA_H_

#ifdef __cplusplus
extern "C" {
#endif


#include "schema.h"
/**
 * @fn      int jdata_init(void*)
 * @brief   Initializes schema parser.
 *
 * @param   ip
 * @return  On success, 0
 *          On failure, non zero value
 */
int jdata_init(void* ip);

/**
 * @fn      void jdata_exit()
 * @brief   Calls exit for schema parser
 *
 */
void jdata_exit();
/**
 * @fn      void jdata_fetch_bs_certs*(char*, uint16_t*)
 * @brief   reads the bootstrap certificates from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_bs_certs(char *puuid, uint16_t *size);

/**
 * @fn      void jdata_fetch_cloud_certs*(char*, uint16_t*)
 * @brief   reads the cloud certificates from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_cloud_certs(char *puuid, uint16_t *size);

/**
 * @fn      void jdata_fetch_fact_calib*(char*, uint16_t*)
 * @brief   reads the factory calibration from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_fact_calib(char *puuid, uint16_t *size);

/**
 * @fn      void jdata_fetch_fact_cfg*(char*, uint16_t*)
 * @brief   reads the factory configuration from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_fact_cfg(char *puuid, uint16_t *size);
/**
 * @fn      void jdata_fetch_header*(char*, uint16_t*)
 * @brief   reads the schema header from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_header(char *puuid, uint16_t *size);

/**
 * @fn      void jdata_fetch_idx*(char*, uint16_t*)
 * @brief   reads the index table from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_idx(char* puuid, uint16_t *size);

/**
 * @fn      void jdata_fetch_module_cfg*(char*, uint16_t*, uint8_t)
 * @brief   reads the module config from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @param   count
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_module_cfg(char *puuid, uint16_t *size, uint8_t count);

/**
 * @fn      void jdata_fetch_module_info_by_uuid*(char*, uint16_t*, uint8_t)
 * @brief   reads the module info from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @param   count
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_module_info_by_uuid(char *puuid, uint16_t *size, uint8_t count);

/**
 * @fn      void jdata_fetch_unit_cfg*(char*, uint16_t*, uint8_t)
 * @brief   reads the unit configuration from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @param   count
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_unit_cfg(char *puuid, uint16_t *size, uint8_t count);

/**
 * @fn      void jdata_fetch_node_info*(char*, uint16_t*)
 * @brief   reads the unit info from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_node_info(char *puuid, uint16_t *size);

/**
 * @fn      void jdata_fetch_user_calib*(char*, uint16_t*)
 * @brief   reads the user calibration data from the schema parser for the module
 *          specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_user_calib(char *puuid, uint16_t *size);

/**
 * @fn      void jdata_fetch_user_cfg*(char*, uint16_t*)
 * @brief   reads the user configuration data from the schema parser for the
 *          module specified by UUID.
 *
 * @param   puuid
 * @param   size
 * @return  On success, Information read
 *          On failure, NULL
 */
void *jdata_fetch_user_cfg(char *puuid, uint16_t *size);

#ifdef __cplusplus
}
#endif

#endif /* INC_JDATA_H_ */
