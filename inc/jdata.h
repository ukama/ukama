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

#include "schema.h"

int jdata_init(void* ip);
void jdata_exit();
void *jdata_fetch_bs_certs(char *puuid, uint16_t *size);
void *jdata_fetch_cloud_certs(char *puuid, uint16_t *size);
void *jdata_fetch_fact_calib(char *puuid, uint16_t *size);
void *jdata_fetch_fact_cfg(char *puuid, uint16_t *size);
void *jdata_fetch_header(char *puuid, uint16_t *size);
void *jdata_fetch_index(char* puuid, uint16_t *size);
void *jdata_fetch_module_cfg(char *puuid, uint16_t *size, uint8_t count);
void *jdata_fetch_module_info_by_uuid(char *puuid, uint16_t *size, uint8_t count);
void *jdata_fetch_unit_cfg(char *puuid, uint16_t *size, uint8_t count);
void *jdata_fetch_unit_info(char *puuid, uint16_t *size);
void *jdata_fetch_user_calib(char *puuid, uint16_t *size);
void *jdata_fetch_user_cfg(char *puuid, uint16_t *size);

#endif /* IDB_JP_H_ */
