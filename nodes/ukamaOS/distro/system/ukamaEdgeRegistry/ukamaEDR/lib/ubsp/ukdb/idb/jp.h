/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef IDB_JP_H_
#define IDB_JP_H_

#include "headers/ubsp/ukdblayout.h"

#include <stdio.h>
#include <stdint.h>
#include <string.h>

int jp_init(void* ip);
void jp_exit();
void *jp_fetch_header(char *puuid, uint16_t *size);
void* jp_fetch_index(char* puuid, uint16_t *size);
void *jp_fetch_unit_info(char *puuid, uint16_t *size);
void *jp_fetch_unit_cfg(char *puuid, uint16_t *size, uint8_t count);
void *jp_fetch_module_info_by_uuid(char *puuid, uint16_t *size, uint8_t count);
void *jp_fetch_module_cfg(char *puuid, uint16_t *size, uint8_t count);
void *jp_fetch_fact_config(char *puuid, uint16_t *size);
void *jp_fetch_user_config(char *puuid, uint16_t *size);
void *jp_fetch_fact_calib(char *puuid, uint16_t *size);
void *jp_fetch_user_calib(char *puuid, uint16_t *size);
void *jp_fetch_bs_certs(char *puuid, uint16_t *size);
void *jp_fetch_lwm2m_certs(char *puuid, uint16_t *size);

#endif /* IDB_JP_H_ */
