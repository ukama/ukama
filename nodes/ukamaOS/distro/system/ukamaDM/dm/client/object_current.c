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
#include "objects/current.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>



static uint8_t prv_exec(uint16_t instanceId, uint16_t resourceId,
		uint8_t * buffer, int length, lwm2m_object_t * objectP)
{
	int ret = 0;
	curr_info_t * targetP = NULL;
	void* data = NULL;
	size_t size = 0;
	targetP = (curr_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	switch (resourceId)
	{
	case RES_O_RESET_MIN_AND_MAX_MEASURED_VALUE:
		ret = ereg_exec_sensor(instanceId, OBJ_TYPE_CURR, resourceId, data, &size);
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
		curr_info_t * targetP)
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
	case RES_O_CURR_CALIBRATION_VALUE:
		lwm2m_data_encode_float(targetP->data.calibration_value, dataP);
		return COAP_205_CONTENT;
	case RES_O_APPLICATION_TYPE:
		lwm2m_data_encode_string(targetP->data.application_type, dataP);
		return COAP_205_CONTENT;
	default:
		return COAP_404_NOT_FOUND;
	}
}


static int read_curr_inst_data(uint16_t instanceId, curr_info_t** targetP) {
	int ret = 0;
	CurrObjInfo* data = NULL;
	size_t szcurr = 0;
	szcurr = sizeof(CurrObjInfo);

	/* Read Curr data */
	data = malloc(szcurr);
	if (!data) {
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	ret = ereg_read_inst(instanceId, OBJ_TYPE_CURR, ALL_RESOURCE_ID, data, &szcurr);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Curr data for instance %d\r\n", instanceId);
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	/* Copy the curr data read */
	(*targetP)->data.max_measured_value = data->max_measured_value;
	(*targetP)->data.max_range_value = data->max_range_value;
	(*targetP)->data.min_measured_value = data->min_measured_value;
	(*targetP)->data.min_range_value = data->min_range_value;
	(*targetP)->data.avg_measured_value = data->avg_measured_value;
	(*targetP)->data.calibration_value = data->calibration_value;
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

static uint8_t prv_curr_info_read(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{

	uint8_t result = 0;
	int i = 0;
	curr_info_t * targetP = NULL;
	targetP = (curr_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	/* Read Curr data for instance */
	if (read_curr_inst_data(instanceId, &targetP)) {
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
				RES_O_CURR_CALIBRATION_VALUE,
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

static uint8_t prv_curr_info_discover(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	curr_info_t * targetP;
	uint8_t result;
	int i;

	result = COAP_205_CONTENT;
	targetP = (curr_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
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
				RES_O_CURR_CALIBRATION_VALUE,
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
			case RES_O_CURR_CALIBRATION_VALUE:
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

static uint8_t prv_curr_info_write(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	curr_info_t * targetP;
	int i;
	uint8_t result;
	size_t size = sizeof(CurrObjInfo);

	targetP = (curr_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
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
		case RES_O_CURR_CALIBRATION_VALUE:
			result = objh_set_double_value(dataArray + i, (double *)&(targetP->data.calibration_value));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_CURR, &targetP->data, &size);
			}
			break;
		default:
			return COAP_404_NOT_FOUND;
		}
		i++;
	} while (i < numData && result == COAP_204_CHANGED);

	return result;
}

static uint8_t prv_curr_info_delete(uint16_t id,
		lwm2m_object_t * objectP)
{
	curr_info_t * currInfoInstance = NULL;
	objectP->instanceList = lwm2m_list_remove(objectP->instanceList, id, (lwm2m_list_t **)&currInfoInstance);
	if (NULL == currInfoInstance) return COAP_404_NOT_FOUND;

	lwm2m_free(currInfoInstance);

	return COAP_202_DELETED;
}

static uint8_t prv_curr_info_create(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	curr_info_t * currInfoInstance;
	uint8_t result;

	currInfoInstance = (curr_info_t *)lwm2m_malloc(sizeof(curr_info_t));
	if (NULL == currInfoInstance) return COAP_500_INTERNAL_SERVER_ERROR;
	memset(currInfoInstance, 0, sizeof(curr_info_t));

	currInfoInstance->data.instanceId = instanceId;
	objectP->instanceList = LWM2M_LIST_ADD(objectP->instanceList, currInfoInstance);


	result = prv_curr_info_write(instanceId, numData, dataArray, objectP);

	if (result != COAP_204_CHANGED)
	{
		(void)prv_curr_info_delete(instanceId, objectP);
	}
	else
	{
		result = COAP_201_CREATED;
	}

	return result;
}

void display_curr_info_object(lwm2m_object_t * object)
{
#ifdef WITH_LOGS
	fprintf(stdout, "  /%u: Curr Info object, instances:\r\n", object->objID);
	curr_info_t * currInfoInstance = (curr_info_t *)object->instanceList;
	while (currInfoInstance != NULL)
	{
		fprintf(stdout, "    /%u/%u: instanceId: %u, sensor value: %f",
				object->objID, currInfoInstance->data.instanceId,
				currInfoInstance->data.instanceId, currInfoInstance->data.sensor_value);
		fprintf(stdout, "\r\n");
		currInfoInstance = (curr_info_t *)currInfoInstance->next;
	}
#endif
}

lwm2m_object_t * create_curr_info_object()
{
	lwm2m_object_t * currInfoObj = NULL;
	currInfoObj = (lwm2m_object_t *)lwm2m_malloc(sizeof(lwm2m_object_t));
	if (NULL != currInfoObj)
	{
		memset(currInfoObj, 0, sizeof(lwm2m_object_t));
		currInfoObj->objID = OBJECT_ID_CURR;
		/* set the callbacks. */
		currInfoObj->readFunc = prv_curr_info_read;
		currInfoObj->discoverFunc = prv_curr_info_discover;
		currInfoObj->writeFunc = prv_curr_info_write;
		currInfoObj->createFunc = NULL;
		currInfoObj->deleteFunc = prv_curr_info_delete;
		currInfoObj->executeFunc = prv_exec;
	}
	return currInfoObj;
}

curr_info_t * create_curr_info_instance(uint16_t instance)
{
	curr_info_t * currInfoInstance = NULL;

	/* allocate memory for module info object instance. */
	currInfoInstance = (curr_info_t *)lwm2m_malloc(sizeof(curr_info_t));
	if (NULL == currInfoInstance)
	{
		return NULL;
	}
	memset(currInfoInstance, 0, sizeof(curr_info_t));

	/* Read Curr data for instance */
	if (read_curr_inst_data(instance, &currInfoInstance)) {
		if (currInfoInstance) {
			lwm2m_free(currInfoInstance);
			currInfoInstance = NULL;
		}
	}

	return currInfoInstance;
}

lwm2m_object_t * get_curr_info_object()
{
	int ret = 0;
	lwm2m_object_t * currInfoObj = create_curr_info_object();
	if (currInfoObj == NULL)
	{
		fprintf(stderr, "Failed to create curr info object\r\n");
		return NULL;
	}

	int *count = lwm2m_malloc(sizeof(int));
	if (!count) {
		fprintf(stderr, "Failed to allocate memory for curr sensor count\r\n");
		lwm2m_free(currInfoObj);
		currInfoObj = NULL;
		goto cleanup;
	}

	size_t szcount = sizeof(int);
	ret = ereg_read_inst_count(OBJ_TYPE_CURR, count, &szcount);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Curr sensor count\r\n");
		lwm2m_free(currInfoObj);
		currInfoObj = NULL;
		goto cleanup;
	}

	/* Create instances for CurrInfo Object. */
	for (uint16_t iter = 0; iter < *count; iter++)
	{
		curr_info_t * currInfoInstance = create_curr_info_instance(iter);
		if (currInfoInstance == NULL)
		{
			fprintf(stderr, "Failed to create Curr info instance\r\n");
			lwm2m_free(currInfoObj);
			currInfoObj = NULL;
			goto cleanup;
		}
		/* add the curr sensor instance to the Curr info object. */
		currInfoObj->instanceList = LWM2M_LIST_ADD(currInfoObj->instanceList, currInfoInstance);
	}

	cleanup:
	if (count) {
		lwm2m_free(count);
		count = NULL;
	}

	return currInfoObj;
}
