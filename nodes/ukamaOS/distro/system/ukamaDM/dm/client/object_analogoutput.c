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
#include "objects/analog_output.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>


static uint8_t prv_get_value(lwm2m_data_t * dataP,
		analog_output_info_t * targetP)
{
	switch (dataP->id)
	{
	case RES_M_OUT_CURR_VALUE:
		lwm2m_data_encode_float(targetP->data.outputcurr, dataP);
		return COAP_205_CONTENT;
	case RES_O_MIN_RANGE_VALUE:
		lwm2m_data_encode_float(targetP->data.minrange, dataP);
		return COAP_205_CONTENT;
	case RES_O_MAX_RANGE_VALUE:
		lwm2m_data_encode_float(targetP->data.maxrange, dataP);
		return COAP_205_CONTENT;
	case RES_O_APPLICATION_TYPE:
		lwm2m_data_encode_string(targetP->data.application_type, dataP);
		return COAP_205_CONTENT;
	default:
		return COAP_404_NOT_FOUND;
	}
}

static int read_adc_inst_data(uint16_t instanceId, analog_output_info_t** targetP) {
	int ret = 0;
	AdcObjInfo* data = NULL;
	size_t szadc = 0;
	szadc = sizeof(AdcObjInfo);
	/* Read Adc data */
	data = malloc(szadc);
	if (!data) {
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	ret = ereg_read_inst(instanceId, OBJ_TYPE_ADC, ALL_RESOURCE_ID, data, &szadc);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Adc data for instance %d\r\n", instanceId);
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	/* Copy the adc data read */
	(*targetP)->data.outputcurr = data->outputcurr;
	(*targetP)->data.minrange = data->minrange;
	(*targetP)->data.maxrange = data->maxrange;
	(*targetP)->data.instanceId = data->instanceId;
	strcpy((*targetP)->data.application_type, data->application_type);

	cleanup:
	if(data) {
		free(data);
		data = NULL;
	}
	return ret;
}

static uint8_t prv_adc_info_read(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{

	uint8_t result = 0;
	int i = 0;
	analog_output_info_t * targetP = NULL;
	targetP = (analog_output_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	/* Read Adc data for instance */
	if (read_adc_inst_data(instanceId, &targetP)) {
		return COAP_500_INTERNAL_SERVER_ERROR;
	}

	// is the server asking for the full instance ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_OUT_CURR_VALUE,
				RES_O_MIN_RANGE_VALUE,
				RES_O_MAX_RANGE_VALUE,
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

static uint8_t prv_adc_info_discover(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	analog_output_info_t * targetP;
	uint8_t result;
	int i;

	result = COAP_205_CONTENT;
	targetP = (analog_output_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full object ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_OUT_CURR_VALUE,
				RES_O_MIN_RANGE_VALUE,
				RES_O_MAX_RANGE_VALUE,
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
			case RES_M_OUT_CURR_VALUE:
			case RES_O_MIN_RANGE_VALUE:
			case RES_O_MAX_RANGE_VALUE:
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

static uint8_t prv_adc_info_write(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	analog_output_info_t * targetP;
	int i;
	uint8_t result;

	targetP = (analog_output_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP)
	{
		return COAP_404_NOT_FOUND;
	}

	i = 0;
	do
	{
		switch (dataArray[i].id)
		{
		case RES_M_OUT_CURR_VALUE:
		case RES_O_MIN_RANGE_VALUE:
		case RES_O_MAX_RANGE_VALUE:
		case RES_O_APPLICATION_TYPE:
			result = COAP_405_METHOD_NOT_ALLOWED;
			break;
		default:
			return COAP_404_NOT_FOUND;
		}
		i++;
	} while (i < numData && result == COAP_204_CHANGED);

	return result;
}

static uint8_t prv_adc_info_delete(uint16_t id,
		lwm2m_object_t * objectP)
{
	analog_output_info_t * adcInfoInstance = NULL;
	objectP->instanceList = lwm2m_list_remove(objectP->instanceList, id, (lwm2m_list_t **)&adcInfoInstance);
	if (NULL == adcInfoInstance) return COAP_404_NOT_FOUND;

	lwm2m_free(adcInfoInstance);

	return COAP_202_DELETED;
}

static uint8_t prv_adc_info_create(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	analog_output_info_t * adcInfoInstance;
	uint8_t result;

	adcInfoInstance = (analog_output_info_t *)lwm2m_malloc(sizeof(analog_output_info_t));
	if (NULL == adcInfoInstance) return COAP_500_INTERNAL_SERVER_ERROR;
	memset(adcInfoInstance, 0, sizeof(analog_output_info_t));

	adcInfoInstance->data.instanceId = instanceId;
	objectP->instanceList = LWM2M_LIST_ADD(objectP->instanceList, adcInfoInstance);

	//todo: will not be able to create objects as adc info resources are all read only. not sure if object instance can be created using coap calls.
			result = prv_adc_info_write(instanceId, numData, dataArray, objectP);

			if (result != COAP_204_CHANGED)
			{
				(void)prv_adc_info_delete(instanceId, objectP);
			}
			else
			{
				result = COAP_201_CREATED;
			}

			return result;
}

void display_adc_info_object(lwm2m_object_t * object)
{
#ifdef WITH_LOGS
	fprintf(stdout, "  /%u: Adc Info object, instances:\r\n", object->objID);
	analog_output_info_t * adcInfoInstance = (analog_output_info_t *)object->instanceList;
	while (adcInfoInstance != NULL)
	{
		fprintf(stdout, "    /%u/%u: instanceId: %u, sensor value: %f",
				object->objID, adcInfoInstance->data.instanceId,
				adcInfoInstance->data.instanceId, adcInfoInstance->data.outputcurr);
		fprintf(stdout, "\r\n");
		adcInfoInstance = (analog_output_info_t *)adcInfoInstance->next;
	}
#endif
}

lwm2m_object_t * create_adc_info_object()
{
	lwm2m_object_t * adcInfoObj = NULL;
	adcInfoObj = (lwm2m_object_t *)lwm2m_malloc(sizeof(lwm2m_object_t));
	if (NULL != adcInfoObj)
	{
		memset(adcInfoObj, 0, sizeof(lwm2m_object_t));
		adcInfoObj->objID = OBJECT_ID_ANALOG_OUTPUT;
		/* set the callbacks. */
		adcInfoObj->readFunc = prv_adc_info_read;
		adcInfoObj->discoverFunc = prv_adc_info_discover;
		adcInfoObj->writeFunc = prv_adc_info_write;
		adcInfoObj->createFunc = NULL;
		adcInfoObj->deleteFunc = prv_adc_info_delete;
		adcInfoObj->executeFunc = NULL;
	}
	return adcInfoObj;
}

analog_output_info_t * create_adc_info_instance(uint16_t instance)
{
	analog_output_info_t * adcInfoInstance = NULL;

	/* allocate memory for module info object instance. */
	adcInfoInstance = (analog_output_info_t *)lwm2m_malloc(sizeof(analog_output_info_t));
	if (NULL == adcInfoInstance)
	{
		return NULL;
	}
	memset(adcInfoInstance, 0, sizeof(analog_output_info_t));

	/* Read Adc data for instance */
	if (read_adc_inst_data(instance, &adcInfoInstance)) {
		if (adcInfoInstance) {
			lwm2m_free(adcInfoInstance);
			adcInfoInstance = NULL;
		}
	}

	return adcInfoInstance;
}

lwm2m_object_t * get_adc_info_object()
{
	int ret = 0;
	lwm2m_object_t * adcInfoObj = create_adc_info_object();
	if (adcInfoObj == NULL)
	{
		fprintf(stderr, "Failed to create adc info object\r\n");
		return NULL;
	}

	int *count = lwm2m_malloc(sizeof(int));
	if (!count) {
		fprintf(stderr, "Failed to allocate memory for adc sensor count\r\n");
		lwm2m_free(adcInfoObj);
		adcInfoObj = NULL;
		goto cleanup;
	}

	size_t szcount = sizeof(int);
	ret = ereg_read_inst_count(OBJ_TYPE_ADC, count, &szcount);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Adc sensor count\r\n");
		lwm2m_free(adcInfoObj);
		adcInfoObj = NULL;
		goto cleanup;
	}

	/* Create instances for AdcInfo Object. */
	for (uint16_t iter = 0; iter < *count; iter++)
	{
		analog_output_info_t * adcInfoInstance = create_adc_info_instance(iter);
		if (adcInfoInstance == NULL)
		{
			fprintf(stderr, "Failed to create Adc info instance\r\n");
			lwm2m_free(adcInfoObj);
			adcInfoObj = NULL;
			goto cleanup;
		}
		/* add the adc sensor instance to the Adc info object. */
		adcInfoObj->instanceList = LWM2M_LIST_ADD(adcInfoObj->instanceList, adcInfoInstance);
	}

	cleanup:
	if (count) {
		lwm2m_free(count);
		count = NULL;
	}

	return adcInfoObj;
}
