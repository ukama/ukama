/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef CLIENT_OBJECT_HELPER_H_
#define CLIENT_OBJECT_HELPER_H_

uint8_t objh_send_data_ukama_edr(uint16_t instanceId, uint16_t rid, int objectType,
		void* data, size_t *size);
uint8_t objh_set_bool_value(lwm2m_data_t * dataArray, bool * data);
uint8_t objh_set_double_value(lwm2m_data_t * dataArray, double * data);
uint8_t objh_set_int_value(lwm2m_data_t * dataArray, uint32_t * data);
uint8_t objh_set_str_value(lwm2m_data_t * dataArray, char* data);

int objh_store_data(char* filename, char* data, int size);
int objh_parse_addr(char* data, int size, char** addr);

#endif /* CLIENT_OBJECT_HELPER_H_ */
