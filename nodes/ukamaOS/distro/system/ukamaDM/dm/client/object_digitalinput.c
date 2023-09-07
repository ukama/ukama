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
#include "objects/objects.h"
#include "objects/digital_input.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>



static uint8_t prv_exec(uint16_t instanceId, uint16_t resourceId,
		uint8_t * buffer, int length, lwm2m_object_t * objectP)
{
	int ret = 0;
	digital_input_t * targetP = NULL;
	void* data = NULL;
	size_t size = 0;
	targetP = (digital_input_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	switch (resourceId)
	{
	case RES_O_DIGITAL_INPUT_COUNTER_RESET:
		ret = ereg_exec_sensor(instanceId, OBJ_TYPE_DIP, resourceId, data, &size);
		if (ret){
			fprintf(stderr, "Failed to execute %d\r\n", resourceId);
			return COAP_500_INTERNAL_SERVER_ERROR;
		}
		return COAP_204_CHANGED;
	default:
		return COAP_405_METHOD_NOT_ALLOWED;
	}

}

static uint8_t prv_get_value(lwm2m_data_t * dataP,
		digital_input_t * targetP)
{
	switch (dataP->id)
	{
	case RES_M_DIGITAL_INPUT_STATE:
		lwm2m_data_encode_bool(targetP->data.digital_state, dataP);
		return COAP_205_CONTENT;
	case RES_O_DIGITAL_INPUT_COUNTER:
		lwm2m_data_encode_int(targetP->data.digital_counter, dataP);
		return COAP_205_CONTENT;
	case RES_O_DIGITAL_INPUT_POLARITY:
		lwm2m_data_encode_bool(targetP->data.digital_polarity, dataP);
		return COAP_205_CONTENT;
	case RES_O_DIGITAL_INPUT_DEBOUNCE:
		lwm2m_data_encode_int(targetP->data.digital_debounce, dataP);
		return COAP_205_CONTENT;
	case RES_O_DIGITIAL_INPUT_EDGE_SELECTION:
		lwm2m_data_encode_int(targetP->data.digitial_edge_selection, dataP);
		return COAP_205_CONTENT;
	case RES_O_APPLICATION_TYPE:
		lwm2m_data_encode_string(targetP->data.application_type, dataP);
		return COAP_205_CONTENT;
	case RES_O_SENSOR_TYPE:
		lwm2m_data_encode_string(targetP->data.sensor_type, dataP);
		return COAP_205_CONTENT;
	default:
		return COAP_404_NOT_FOUND;
	}
}

static uint8_t prv_set_int_value(lwm2m_data_t * dataArray, uint32_t * data) {
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

static uint8_t prv_send_data_ukama_edr(uint16_t instanceId, uint16_t rid, int objectType, void* data, size_t *size) {
	int ret = 0;
	ret = ereg_write_inst(instanceId, objectType, rid, data, size);
	if (ret) {
		return COAP_500_INTERNAL_SERVER_ERROR;
	} else {
		return COAP_204_CHANGED;
	}
}

static int read_dip_inst_data(uint16_t instanceId, digital_input_t** targetP) {
	int ret = 0;
	DipObjInfo* data = NULL;
	size_t szdip = 0;
	szdip = sizeof(DipObjInfo);
	/* Read Dip data */
	data = malloc(szdip);
	if (!data) {
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	ret = ereg_read_inst(instanceId, OBJ_TYPE_DIP, ALL_RESOURCE_ID, data, &szdip);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Dip data for instance %d\r\n", instanceId);
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	/* Copy the dip data read */
	(*targetP)->data.digital_state = data->digital_state;
	(*targetP)->data.digital_counter = data->digital_counter;
	(*targetP)->data.digital_polarity = data->digital_polarity;
	(*targetP)->data.digital_debounce = data->digital_debounce;
	(*targetP)->data.digitial_edge_selection = data->digitial_edge_selection;
	(*targetP)->data.instanceId = data->instanceId;
	strcpy((*targetP)->data.application_type, data->application_type);
	strcpy((*targetP)->data.sensor_type, data->sensor_type);

	cleanup:
	if(data) {
		free(data);
		data = NULL;
	}
	return ret;
}

static uint8_t prv_dip_info_read(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{

	uint8_t result = 0;
	int i = 0;
	digital_input_t * targetP = NULL;
	targetP = (digital_input_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	/* Read Dip data for instance */
	if (read_dip_inst_data(instanceId, &targetP)) {
		return COAP_500_INTERNAL_SERVER_ERROR;
	}

	// is the server asking for the full instance ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_DIGITAL_INPUT_STATE,
				RES_O_DIGITAL_INPUT_COUNTER,
				RES_O_DIGITAL_INPUT_POLARITY,
				RES_O_DIGITAL_INPUT_DEBOUNCE,
				RES_O_DIGITIAL_INPUT_EDGE_SELECTION,
				RES_O_APPLICATION_TYPE,
				RES_O_SENSOR_TYPE
		};
		int nbRes = sizeof(resList)/sizeof(uint16_t);
		*dataArrayP = lwm2m_data_new(nbRes);
		if (*dataArrayP == NULL) return COAP_500_INTERNAL_SERVER_ERROR;
		*numDataP = nbRes;
		for (i = 0 ; i < nbRes ; i++)
		{
			(*dataArrayP)[i].id = resList[i];
		}
	}

	i = 0;
	do
	{
		result = prv_get_value((*dataArrayP) + i, targetP);
		i++;
	} while (i < *numDataP && result == COAP_205_CONTENT);

	return result;
}

static uint8_t prv_dip_info_discover(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	digital_input_t * targetP;
	uint8_t result;
	int i;

	result = COAP_205_CONTENT;
	targetP = (digital_input_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full object ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_DIGITAL_INPUT_STATE,
				RES_O_DIGITAL_INPUT_COUNTER,
				RES_O_DIGITAL_INPUT_POLARITY,
				RES_O_DIGITAL_INPUT_DEBOUNCE,
				RES_O_DIGITIAL_INPUT_EDGE_SELECTION,
				RES_O_APPLICATION_TYPE,
				RES_O_SENSOR_TYPE
		};
		int nbRes = sizeof(resList) / sizeof(uint16_t);


		*dataArrayP = lwm2m_data_new(nbRes);
		if (*dataArrayP == NULL) return COAP_500_INTERNAL_SERVER_ERROR;
		*numDataP = nbRes;
		for (i = 0; i < nbRes; i++)
		{
			(*dataArrayP)[i].id = resList[i];
		}
	}
	else
	{
		for (i = 0; i < *numDataP && result == COAP_205_CONTENT; i++)
		{
			switch ((*dataArrayP)[i].id)
			{
			case RES_M_DIGITAL_INPUT_STATE:
			case RES_O_DIGITAL_INPUT_COUNTER:
			case RES_O_DIGITAL_INPUT_POLARITY:
			case RES_O_DIGITAL_INPUT_DEBOUNCE:
			case RES_O_DIGITIAL_INPUT_EDGE_SELECTION:
			case RES_O_APPLICATION_TYPE:
			case RES_O_SENSOR_TYPE:
				break;
			default:
				result = COAP_404_NOT_FOUND;
				break;
			}
		}
	}

	return result;
}

static uint8_t prv_dip_info_write(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	digital_input_t * targetP;
	int i;
	uint8_t result;
	size_t size = sizeof(DipObjInfo);

	targetP = (digital_input_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP)
	{
		return COAP_404_NOT_FOUND;
	}

	i = 0;
	do
	{
		switch (dataArray[i].id)
		{
		case RES_M_DIGITAL_INPUT_STATE:
		case RES_O_DIGITAL_INPUT_COUNTER:
		case RES_O_DIGITAL_INPUT_POLARITY:
		case RES_O_DIGITAL_INPUT_DEBOUNCE:
		case RES_O_DIGITIAL_INPUT_EDGE_SELECTION:
		case RES_O_APPLICATION_TYPE:
		case RES_O_SENSOR_TYPE:
			result = COAP_405_METHOD_NOT_ALLOWED;
			break;
		default:
			return COAP_404_NOT_FOUND;
		}
		i++;
	} while (i < numData && result == COAP_204_CHANGED);

	return result;
}

static uint8_t prv_dip_info_delete(uint16_t id,
		lwm2m_object_t * objectP)
{
	digital_input_t * dipInfoInstance = NULL;
	objectP->instanceList = lwm2m_list_remove(objectP->instanceList, id, (lwm2m_list_t **)&dipInfoInstance);
	if (NULL == dipInfoInstance) return COAP_404_NOT_FOUND;

	lwm2m_free(dipInfoInstance);

	return COAP_202_DELETED;
}

static uint8_t prv_dip_info_create(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	digital_input_t * dipInfoInstance;
	uint8_t result;

	dipInfoInstance = (digital_input_t *)lwm2m_malloc(sizeof(digital_input_t));
	if (NULL == dipInfoInstance) return COAP_500_INTERNAL_SERVER_ERROR;
	memset(dipInfoInstance, 0, sizeof(digital_input_t));

	dipInfoInstance->data.instanceId = instanceId;
	objectP->instanceList = LWM2M_LIST_ADD(objectP->instanceList, dipInfoInstance);


	result = prv_dip_info_write(instanceId, numData, dataArray, objectP);

	if (result != COAP_204_CHANGED)
	{
		(void)prv_dip_info_delete(instanceId, objectP);
	}
	else
	{
		result = COAP_201_CREATED;
	}

	return result;
}

void display_dip_info_object(lwm2m_object_t * object)
{
#ifdef WITH_LOGS
	fprintf(stdout, "  /%u: Dip Info object, instances:\r\n", object->objID);
	digital_input_t * dipInfoInstance = (digital_input_t *)object->instanceList;
	while (dipInfoInstance != NULL)
	{
		fprintf(stdout, "    /%u/%u: instanceId: %u, state value: %f",
				object->objID, dipInfoInstance->data.instanceId,
				dipInfoInstance->data.instanceId, dipInfoInstance->data.digital_state);
		fprintf(stdout, "\r\n");
		dipInfoInstance = (digital_input_t *)dipInfoInstance->next;
	}
#endif
}

lwm2m_object_t * create_dip_info_object()
{
	lwm2m_object_t * dipInfoObj = NULL;
	dipInfoObj = (lwm2m_object_t *)lwm2m_malloc(sizeof(lwm2m_object_t));
	if (NULL != dipInfoObj)
	{
		memset(dipInfoObj, 0, sizeof(lwm2m_object_t));
		dipInfoObj->objID = OBJECT_ID_DIGITAL_INPUT;
		/* set the callbacks. */
		dipInfoObj->readFunc = prv_dip_info_read;
		dipInfoObj->discoverFunc = prv_dip_info_discover;
		dipInfoObj->writeFunc = prv_dip_info_write;
		dipInfoObj->createFunc = NULL;
		dipInfoObj->deleteFunc = prv_dip_info_delete;
		dipInfoObj->executeFunc = prv_exec;
	}
	return dipInfoObj;
}

digital_input_t * create_dip_info_instance(uint16_t instance)
{
	digital_input_t * dipInfoInstance = NULL;

	/* allocate memory for module info object instance. */
	dipInfoInstance = (digital_input_t *)lwm2m_malloc(sizeof(digital_input_t));
	if (NULL == dipInfoInstance)
	{
		return NULL;
	}
	memset(dipInfoInstance, 0, sizeof(digital_input_t));

	/* Read Dip data for instance */
	if (read_dip_inst_data(instance, &dipInfoInstance)) {
		if (dipInfoInstance) {
			lwm2m_free(dipInfoInstance);
			dipInfoInstance = NULL;
		}
	}

	return dipInfoInstance;
}

lwm2m_object_t * get_dip_info_object()
{
	int ret = 0;
	lwm2m_object_t * dipInfoObj = create_dip_info_object();
	if (dipInfoObj == NULL)
	{
		fprintf(stderr, "Failed to create dip info object\r\n");
		return NULL;
	}

	int *count = lwm2m_malloc(sizeof(int));
	if (!count) {
		fprintf(stderr, "Failed to allocate memory for dip sensor count\r\n");
		lwm2m_free(dipInfoObj);
		dipInfoObj = NULL;
		goto cleanup;
	}

	size_t szcount = sizeof(int);
	ret = ereg_read_inst_count(OBJ_TYPE_DIP, count, &szcount);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Dip sensor count\r\n");
		lwm2m_free(dipInfoObj);
		dipInfoObj = NULL;
		goto cleanup;
	}

	/* Create instances for DipInfo Object. */
	for (uint16_t iter = 0; iter < *count; iter++)
	{
		digital_input_t * dipInfoInstance = create_dip_info_instance(iter);
		if (dipInfoInstance == NULL)
		{
			fprintf(stderr, "Failed to create Dip info instance\r\n");
			lwm2m_free(dipInfoObj);
			dipInfoObj = NULL;
			goto cleanup;
		}
		/* add the dip sensor instance to the Dip info object. */
		dipInfoObj->instanceList = LWM2M_LIST_ADD(dipInfoObj->instanceList, dipInfoInstance);
	}

	cleanup:
	if (count) {
		lwm2m_free(count);
		count = NULL;
	}

	return dipInfoObj;
}
