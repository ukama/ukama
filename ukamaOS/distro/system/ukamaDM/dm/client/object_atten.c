/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include "liblwm2m.h"
#include "object_helper.h"
#include "inc/ereg.h"
#include "objects/objects.h"
#include "objects/atten.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>

static uint8_t prv_get_value(lwm2m_data_t * dataP,
		atten_info_t * targetP)
{
	switch (dataP->id)
	{
	case RES_M_ATTVALUE:
		lwm2m_data_encode_int(targetP->data.attvalue, dataP);
		return COAP_205_CONTENT;
	case RES_M_MINRANGE:
		lwm2m_data_encode_int(targetP->data.minrange, dataP);
		return COAP_205_CONTENT;
	case RES_M_MAXRANGE:
		lwm2m_data_encode_int(targetP->data.maxrange, dataP);
		return COAP_205_CONTENT;
	case RES_M_LATCH:
		lwm2m_data_encode_int(targetP->data.latchenable, dataP);
		return COAP_205_CONTENT;
	case RES_M_SENSOR_UNITS:
		lwm2m_data_encode_string(targetP->data.sensor_units, dataP);
		return COAP_205_CONTENT;
	case RES_O_APPLICATION_TYPE:
		lwm2m_data_encode_string(targetP->data.application_type, dataP);
		return COAP_205_CONTENT;
	default:
		return COAP_404_NOT_FOUND;
	}
}

static int read_atten_inst_data(uint16_t instanceId, atten_info_t** targetP) {
	int ret = 0;
	AttObjInfo* data = NULL;
	size_t szatten = 0;
	szatten = sizeof(AttObjInfo);
	/* Read Atten data */
	data = malloc(szatten);
	if (!data) {
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	ret = ereg_read_inst(instanceId, OBJ_TYPE_ATT, ALL_RESOURCE_ID, data, &szatten);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Atten data for instance %d\r\n", instanceId);
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;

	}

	/* Copy the Atten data read */
	(*targetP)->data.attvalue = data->attvalue;
	(*targetP)->data.latchenable = data->latchenable;
	(*targetP)->data.minrange = data->minrange;
	(*targetP)->data.maxrange = data->maxrange;
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

static uint8_t prv_atten_info_read(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{

	uint8_t result = 0;
	int i = 0;
	atten_info_t * targetP = NULL;
	targetP = (atten_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	/* Read Atten data for instance */
	if (read_atten_inst_data(instanceId, &targetP)) {
		return COAP_500_INTERNAL_SERVER_ERROR;
	}

	// is the server asking for the full instance ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_ATTVALUE,
				RES_M_MINRANGE,
				RES_M_MAXRANGE,
				RES_M_LATCH,
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

static uint8_t prv_atten_info_discover(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	atten_info_t * targetP;
	uint8_t result;
	int i;

	result = COAP_205_CONTENT;
	targetP = (atten_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full object ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_ATTVALUE,
				RES_M_MINRANGE,
				RES_M_MAXRANGE,
				RES_M_LATCH,
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
			case RES_M_ATTVALUE:
			case RES_M_MINRANGE:
			case RES_M_MAXRANGE:
			case RES_M_LATCH:
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

static uint8_t prv_atten_info_write(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	atten_info_t * targetP;
	int i;
	uint8_t result;
	size_t size = sizeof(AttObjInfo);

	targetP = (atten_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP)
	{
		return COAP_404_NOT_FOUND;
	}

	i = 0;
	do
	{
		switch (dataArray[i].id)
		{
		case RES_M_MINRANGE:
		case RES_M_MAXRANGE:
		case RES_M_SENSOR_UNITS:
		case RES_O_APPLICATION_TYPE:
			result = COAP_405_METHOD_NOT_ALLOWED;
			break;
		case RES_M_ATTVALUE:
			result = objh_set_int_value(dataArray + i, (uint32_t *)&(targetP->data.attvalue));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_ATT, &targetP->data, &size);
			}
			break;
		case RES_M_LATCH:
			result = objh_set_int_value(dataArray + i, (uint32_t *)&(targetP->data.latchenable));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_ATT, &targetP->data, &size);
			}
			break;
		default:
			return COAP_404_NOT_FOUND;
		}
		i++;
	} while (i < numData && result == COAP_204_CHANGED);

	return result;
}

static uint8_t prv_atten_info_delete(uint16_t id,
		lwm2m_object_t * objectP)
{
	atten_info_t * attenInfoInstance = NULL;
	objectP->instanceList = lwm2m_list_remove(objectP->instanceList, id, (lwm2m_list_t **)&attenInfoInstance);
	if (NULL == attenInfoInstance) return COAP_404_NOT_FOUND;

	lwm2m_free(attenInfoInstance);

	return COAP_202_DELETED;
}

static uint8_t prv_atten_info_create(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	atten_info_t * attenInfoInstance;
	uint8_t result;

	attenInfoInstance = (atten_info_t *)lwm2m_malloc(sizeof(atten_info_t));
	if (NULL == attenInfoInstance) return COAP_500_INTERNAL_SERVER_ERROR;
	memset(attenInfoInstance, 0, sizeof(atten_info_t));

	attenInfoInstance->data.instanceId = instanceId;
	objectP->instanceList = LWM2M_LIST_ADD(objectP->instanceList, attenInfoInstance);

	result = prv_atten_info_write(instanceId, numData, dataArray, objectP);

	if (result != COAP_204_CHANGED)
	{
		(void)prv_atten_info_delete(instanceId, objectP);
	}
	else
	{
		result = COAP_201_CREATED;
	}

	return result;
}

void display_atten_info_object(lwm2m_object_t * object)
{
#ifdef WITH_LOGS
fprintf(stdout, "  /%u: Atten Info object, instances:\r\n", object->objID);
atten_info_t * attenInfoInstance = (atten_info_t *)object->instanceList;
while (attenInfoInstance != NULL)
{
	fprintf(stdout, "    /%u/%u: instanceId: %u, sensor value: %f",
			object->objID, attenInfoInstance->data.instanceId,
			attenInfoInstance->data.instanceId, attenInfoInstance->data.attvalue);
	fprintf(stdout, "\r\n");
	attenInfoInstance = (atten_info_t *)attenInfoInstance->next;
}
#endif
}

lwm2m_object_t * create_atten_info_object()
{
	lwm2m_object_t * attenInfoObj = NULL;
	attenInfoObj = (lwm2m_object_t *)lwm2m_malloc(sizeof(lwm2m_object_t));
	if (NULL != attenInfoObj)
	{
		memset(attenInfoObj, 0, sizeof(lwm2m_object_t));
		attenInfoObj->objID = OBJECT_ID_ATTEN_OUTPUT;
		/* set the callbacks. */
		attenInfoObj->readFunc = prv_atten_info_read;
		attenInfoObj->discoverFunc = prv_atten_info_discover;
		attenInfoObj->writeFunc = prv_atten_info_write;
		attenInfoObj->createFunc = NULL;
		attenInfoObj->deleteFunc = prv_atten_info_delete;
		attenInfoObj->executeFunc = NULL;
	}
	return attenInfoObj;
}

atten_info_t * create_atten_info_instance(uint16_t instance)
{
	atten_info_t * attenInfoInstance = NULL;

	/* allocate memory for module info object instance. */
	attenInfoInstance = (atten_info_t *)lwm2m_malloc(sizeof(atten_info_t));
	if (NULL == attenInfoInstance)
	{
		return NULL;
	}
	memset(attenInfoInstance, 0, sizeof(atten_info_t));

	/* Read Atten data for instance */
	if (read_atten_inst_data(instance, &attenInfoInstance)) {
		if (attenInfoInstance) {
			lwm2m_free(attenInfoInstance);
			attenInfoInstance = NULL;
		}
	}

	return attenInfoInstance;
}

lwm2m_object_t * get_atten_info_object()
{
	int ret = 0;
	lwm2m_object_t * attenInfoObj = create_atten_info_object();
	if (attenInfoObj == NULL)
	{
		fprintf(stderr, "Failed to create Atten info object\r\n");
		return NULL;
	}

	int *count = lwm2m_malloc(sizeof(int));
	if (!count) {
		fprintf(stderr, "Failed to allocate memory for Atten sensor count\r\n");
		lwm2m_free(attenInfoObj);
		attenInfoObj = NULL;
		goto cleanup;
	}

	size_t szcount = sizeof(int);
	ret = ereg_read_inst_count(OBJ_TYPE_ATT, count, &szcount);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Atten sensor count\r\n");
		lwm2m_free(attenInfoObj);
		attenInfoObj = NULL;
		goto cleanup;
	}

	/* Create instances for AttenInfo Object. */
	for (uint16_t iter = 0; iter < *count; iter++)
	{
		atten_info_t * attenInfoInstance = create_atten_info_instance(iter);
		if (attenInfoInstance == NULL)
		{
			fprintf(stderr, "Failed to create Atten info instance\r\n");
			lwm2m_free(attenInfoObj);
			attenInfoObj = NULL;
			goto cleanup;
		}
		/* add the atten sensor instance to the Atten info object. */
		attenInfoObj->instanceList = LWM2M_LIST_ADD(attenInfoObj->instanceList, attenInfoInstance);
	}

	cleanup:
	if (count) {
		lwm2m_free(count);
		count = NULL;
	}

	return attenInfoObj;
}
