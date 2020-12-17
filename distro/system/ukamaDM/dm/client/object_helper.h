/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef CLIENT_OBJECT_HELPER_H_
#define CLIENT_OBJECT_HELPER_H_

uint8_t objh_send_data_ukama_edr(uint16_t instanceId, uint16_t rid, int objectType,
		void* data, size_t *size);
uint8_t objh_set_bool_value(lwm2m_data_t * dataArray, bool * data);
uint8_t objh_set_double_value(lwm2m_data_t * dataArray, double * data);
uint8_t objh_set_int_value(lwm2m_data_t * dataArray, uint32_t * data);
uint8_t objh_set_str_value(lwm2m_data_t * dataArray, char* data);

#endif /* CLIENT_OBJECT_HELPER_H_ */
