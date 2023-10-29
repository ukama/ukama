/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
