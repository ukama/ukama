/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef IDB_CS_H_
#define IDB_CS_H_

#include <stdio.h>
#include <stdint.h>
#include <string.h>

int cs_init(void* data);
int cs_parse(void* data);
void* cs_fetch_header(char* uuid, uint16_t* size);
void* cs_fetch_index( char* uuid, uint16_t* size);
void* cs_fetch_unit_info(char* uuid, uint16_t* size);
void* cs_fetch_unit_cfg(char* uuid, uint16_t* size, uint8_t count);
void* cs_fetch_module_info(char* uuid, uint16_t* size, uint8_t idx);
void* cs_fetch_module_info_by_uuid(char *uuid, uint16_t* size, uint8_t count);
void* cs_fetch_module_cfg( char* uuid, uint16_t* size, uint8_t count);
void* cs_fetch_fact_config(char* uuid, uint16_t* size);
void* cs_fetch_user_config(char* uuid, uint16_t* size);
void* cs_fetch_fact_calib(char* uuid, uint16_t* size);
void* cs_fetch_user_calib(char* uuid, uint16_t* size);
void* cs_fetch_bs_certs(char* uuid, uint16_t* size);
void* cs_fetch_lwm2m_certs(char* uuid, uint16_t* size);


#endif /* IDB_CS_H_ */
