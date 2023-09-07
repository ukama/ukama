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
#include "objects/objects.h"
#include "objects/temperature.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>



static uint8_t prv_exec(uint16_t instanceId, uint16_t resourceId,
		uint8_t * buffer, int length, lwm2m_object_t * objectP)
{
	int ret = 0;
	temp_info_t * targetP = NULL;
	void* data = NULL;
	size_t size = 0;
	targetP = (temp_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	switch (resourceId)
	{
	case RES_O_RESET_MIN_AND_MAX_MEASURED_VALUE:
		ret = ereg_exec_sensor(instanceId, OBJ_TYPE_TMP, resourceId, data, &size);
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
		temp_info_t * targetP)
{
	switch (dataP->id)
	{
	case RES_M_SENSOR_UNITS:
		lwm2m_data_encode_string(targetP->data.sensor_units, dataP);
		return COAP_205_CONTENT;
	case RES_M_SENSOR_VALUE:
		lwm2m_data_encode_float(targetP->data.sensor_value, dataP);
		return COAP_205_CONTENT;
	case RES_O_MIN_MEASURED_VALUE:
		lwm2m_data_encode_float(targetP->data.min_measured_value, dataP);
		return COAP_205_CONTENT;
	case RES_O_MAX_MEASURED_VALUE:
		lwm2m_data_encode_float(targetP->data.max_measured_value, dataP);
		return COAP_205_CONTENT;
	case RES_O_MIN_RANGE_VALUE:
		lwm2m_data_encode_float(targetP->data.min_range_value, dataP);
		return COAP_205_CONTENT;
	case RES_O_MAX_RANGE_VALUE:
		lwm2m_data_encode_float(targetP->data.max_range_value, dataP);
		return COAP_205_CONTENT;
	case RES_O_APPLICATION_TYPE:
		lwm2m_data_encode_string(targetP->data.application_type, dataP);
		return COAP_205_CONTENT;
	default:
		return COAP_404_NOT_FOUND;
	}
}

static int read_temp_inst_data(uint16_t instanceId, temp_info_t** targetP) {
	int ret = 0;
	TempObjInfo* data = NULL;
	size_t sztemp = 0;
	sztemp = sizeof(TempObjInfo);
	/* Read Temp data */
	data = malloc(sztemp);
	if (!data) {
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	ret = ereg_read_inst(instanceId, OBJ_TYPE_TMP, ALL_RESOURCE_ID, data, &sztemp);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Temp data for instance %d\r\n", instanceId);
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	/* Copy the temp data read */
	(*targetP)->data.max_measured_value = data->max_measured_value;
	(*targetP)->data.max_range_value = data->max_range_value;
	(*targetP)->data.min_measured_value = data->min_measured_value;
	(*targetP)->data.min_range_value = data->min_range_value;
	(*targetP)->data.sensor_value = data->sensor_value;
	(*targetP)->data.instanceId = data->instanceId;
	strcpy((*targetP)->data.sensor_units, data->sensor_units);
	strcpy((*targetP)->data.application_type, data->application_type);

	cleanup:
	if(data) {
		free(data);
		data = NULL;
	}
	return ret;
}

static uint8_t prv_temp_info_read(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{

	uint8_t result = 0;
	int i = 0;
	temp_info_t * targetP = NULL;
	targetP = (temp_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	/* Read Temp data for instance */
	if (read_temp_inst_data(instanceId, &targetP)) {
		return COAP_500_INTERNAL_SERVER_ERROR;
	}

	// is the server asking for the full instance ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_SENSOR_VALUE,
				RES_O_MIN_MEASURED_VALUE,
				RES_O_MAX_MEASURED_VALUE,
				RES_O_MIN_RANGE_VALUE,
				RES_O_MAX_RANGE_VALUE,
				RES_M_SENSOR_UNITS,
				RES_O_APPLICATION_TYPE
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

static uint8_t prv_temp_info_discover(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	temp_info_t * targetP;
	uint8_t result;
	int i;

	result = COAP_205_CONTENT;
	targetP = (temp_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full object ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_SENSOR_VALUE,
				RES_O_MIN_MEASURED_VALUE,
				RES_O_MAX_MEASURED_VALUE,
				RES_O_MIN_RANGE_VALUE,
				RES_O_MAX_RANGE_VALUE,
				RES_M_SENSOR_UNITS,
				RES_O_APPLICATION_TYPE
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
			case RES_M_SENSOR_VALUE:
			case RES_O_MIN_MEASURED_VALUE:
			case RES_O_MAX_MEASURED_VALUE:
			case RES_O_MIN_RANGE_VALUE:
			case RES_O_MAX_RANGE_VALUE:
			case RES_M_SENSOR_UNITS:
			case RES_O_APPLICATION_TYPE:
				break;
			default:
				result = COAP_404_NOT_FOUND;
				break;
			}
		}
	}

	return result;
}

static uint8_t prv_temp_info_write(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	temp_info_t * targetP;
	int i;
	uint8_t result;

	targetP = (temp_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP)
	{
		return COAP_404_NOT_FOUND;
	}

	i = 0;
	do
	{
		switch (dataArray[i].id)
		{
		case RES_M_SENSOR_VALUE:
		case RES_O_MIN_MEASURED_VALUE:
		case RES_O_MAX_MEASURED_VALUE:
		case RES_O_MIN_RANGE_VALUE:
		case RES_O_MAX_RANGE_VALUE:
		case RES_M_SENSOR_UNITS:
			result = COAP_405_METHOD_NOT_ALLOWED;
			break;
		default:
			return COAP_404_NOT_FOUND;
		}
		i++;
	} while (i < numData && result == COAP_204_CHANGED);

	return result;
}

static uint8_t prv_temp_info_delete(uint16_t id,
		lwm2m_object_t * objectP)
{
	temp_info_t * tempInfoInstance = NULL;
	objectP->instanceList = lwm2m_list_remove(objectP->instanceList, id, (lwm2m_list_t **)&tempInfoInstance);
	if (NULL == tempInfoInstance) return COAP_404_NOT_FOUND;

	lwm2m_free(tempInfoInstance);

	return COAP_202_DELETED;
}

static uint8_t prv_temp_info_create(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	temp_info_t * tempInfoInstance;
	uint8_t result;

	tempInfoInstance = (temp_info_t *)lwm2m_malloc(sizeof(temp_info_t));
	if (NULL == tempInfoInstance) return COAP_500_INTERNAL_SERVER_ERROR;
	memset(tempInfoInstance, 0, sizeof(temp_info_t));

	tempInfoInstance->data.instanceId = instanceId;
	objectP->instanceList = LWM2M_LIST_ADD(objectP->instanceList, tempInfoInstance);

	//todo: will not be able to create objects as temp info resources are all read only. not sure if object instance can be created using coap calls.
			result = prv_temp_info_write(instanceId, numData, dataArray, objectP);

			if (result != COAP_204_CHANGED)
			{
				(void)prv_temp_info_delete(instanceId, objectP);
			}
			else
			{
				result = COAP_201_CREATED;
			}

			return result;
}

void display_temp_info_object(lwm2m_object_t * object)
{
#ifdef WITH_LOGS
	fprintf(stdout, "  /%u: Temp Info object, instances:\r\n", object->objID);
	temp_info_t * tempInfoInstance = (temp_info_t *)object->instanceList;
	while (tempInfoInstance != NULL)
	{
		fprintf(stdout, "    /%u/%u: instanceId: %u, sensor value: %f",
				object->objID, tempInfoInstance->data.instanceId,
				tempInfoInstance->data.instanceId, tempInfoInstance->data.sensor_value);
		fprintf(stdout, "\r\n");
		tempInfoInstance = (temp_info_t *)tempInfoInstance->next;
	}
#endif
}

lwm2m_object_t * create_temp_info_object()
{
	lwm2m_object_t * tempInfoObj = NULL;
	tempInfoObj = (lwm2m_object_t *)lwm2m_malloc(sizeof(lwm2m_object_t));
	if (NULL != tempInfoObj)
	{
		memset(tempInfoObj, 0, sizeof(lwm2m_object_t));
		tempInfoObj->objID = OBJECT_ID_TMP;
		/* set the callbacks. */
		tempInfoObj->readFunc = prv_temp_info_read;
		tempInfoObj->discoverFunc = prv_temp_info_discover;
		tempInfoObj->writeFunc = prv_temp_info_write;
		tempInfoObj->createFunc = NULL;
		tempInfoObj->deleteFunc = prv_temp_info_delete;
		tempInfoObj->executeFunc = prv_exec;
	}
	return tempInfoObj;
}

temp_info_t * create_temp_info_instance(uint16_t instance)
{
	temp_info_t * tempInfoInstance = NULL;

	/* allocate memory for module info object instance. */
	tempInfoInstance = (temp_info_t *)lwm2m_malloc(sizeof(temp_info_t));
	if (NULL == tempInfoInstance)
	{
		return NULL;
	}
	memset(tempInfoInstance, 0, sizeof(temp_info_t));

	/* Read Temp data for instance */
	if (read_temp_inst_data(instance, &tempInfoInstance)) {
		if (tempInfoInstance) {
			lwm2m_free(tempInfoInstance);
			tempInfoInstance = NULL;
		}
	}

	return tempInfoInstance;
}

lwm2m_object_t * get_temp_info_object()
{
	int ret = 0;
	lwm2m_object_t * tempInfoObj = create_temp_info_object();
	if (tempInfoObj == NULL)
	{
		fprintf(stderr, "Failed to create temp info object\r\n");
		return NULL;
	}

	int *count = lwm2m_malloc(sizeof(int));
	if (!count) {
		fprintf(stderr, "Failed to allocate memory for temp sensor count\r\n");
		lwm2m_free(tempInfoObj);
		tempInfoObj = NULL;
		goto cleanup;
	}

	size_t szcount = sizeof(int);
	ret = ereg_read_inst_count(OBJ_TYPE_TMP, count, &szcount);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Temp sensor count\r\n");
		lwm2m_free(tempInfoObj);
		tempInfoObj = NULL;
		goto cleanup;
	}

	/* Create instances for TempInfo Object. */
	for (uint16_t iter = 0; iter < *count; iter++)
	{
		temp_info_t * tempInfoInstance = create_temp_info_instance(iter);
		if (tempInfoInstance == NULL)
		{
			fprintf(stderr, "Failed to create Temp info instance\r\n");
			lwm2m_free(tempInfoObj);
			tempInfoObj = NULL;
			goto cleanup;
		}
		/* add the temp sensor instance to the Temp info object. */
		tempInfoObj->instanceList = LWM2M_LIST_ADD(tempInfoObj->instanceList, tempInfoInstance);
	}

	cleanup:
	if (count) {
		lwm2m_free(count);
		count = NULL;
	}

	return tempInfoObj;
}
