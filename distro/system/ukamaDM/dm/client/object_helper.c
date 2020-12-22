/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "liblwm2m.h"
#include "inc/ereg.h"
#include "object_helper.h"
#include "objects/atten.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>

uint8_t objh_set_bool_value(lwm2m_data_t * dataArray, bool * data) {
	int ret = 0;
	bool value;
	uint8_t result = COAP_400_BAD_REQUEST;
	if (1 == lwm2m_data_decode_bool(dataArray, &value))
	{
		if ((value == false) || (value == true) )
		{
			*data = value;
			result = COAP_204_CHANGED;
		}
		else
		{
			result = COAP_406_NOT_ACCEPTABLE;
		}
	}
	return result;
}

uint8_t objh_set_int_value(lwm2m_data_t * dataArray, uint32_t * data) {
	int ret = 0;
	int64_t value;
	uint8_t result = COAP_400_BAD_REQUEST;
	if (1 == lwm2m_data_decode_int(dataArray, &value))
	{
		if (value >= 0 && value <= 0xFFFFFFFF)
		{
			*data = value;
			result = COAP_204_CHANGED;
		}
		else
		{
			result = COAP_406_NOT_ACCEPTABLE;
		}
	}
	return result;
}

uint8_t objh_set_str_value(lwm2m_data_t * dataArray, char* data) {
	int64_t value;
	uint8_t result = COAP_400_BAD_REQUEST;
	if ( dataArray->type == LWM2M_TYPE_STRING
			&& dataArray->value.asBuffer.length > 0 ) {
		if (data) {
			lwm2m_free(data);
		}
		size_t szstr = dataArray->value.asBuffer.length + 1;
		data = (char *)lwm2m_malloc(szstr);
		if (data) {
			memset(data, 0, szstr);
			strncpy(data, (char*)dataArray->value.asBuffer.buffer, szstr);
			result = COAP_204_CHANGED;
		} else {
			result =  COAP_500_INTERNAL_SERVER_ERROR;
		}
	}
	if (data) {
		lwm2m_free(data);
	}
	return result;
}

uint8_t objh_set_double_value(lwm2m_data_t * dataArray, double * data) {
	int ret = 0;
	double value;
	uint8_t result = COAP_400_BAD_REQUEST;
	if (1 == lwm2m_data_decode_float(dataArray, &value))
	{
		if (value >= 0 && value <= 0x7FFFFFFFFFFF)
		{
			*data = value;
			result = COAP_204_CHANGED;
		}
		else
		{
			result = COAP_406_NOT_ACCEPTABLE;
		}
	}
	return result;
}

uint8_t objh_send_data_ukama_edr(uint16_t instanceId, uint16_t rid, int objectType, void* data, size_t *size) {
	int ret = 0;
	ret = ereg_write_inst(instanceId, objectType, rid, data, size);
	if (ret) {
		return COAP_500_INTERNAL_SERVER_ERROR;
	} else {
		return COAP_204_CHANGED;
	}
}


