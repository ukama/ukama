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
#include "objects/led.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>

static uint8_t prv_get_value(lwm2m_data_t * dataP,
		led_info_t * targetP)
{
	switch (dataP->id)
	{
	case RES_M_ONOFF_VALUE:
		lwm2m_data_encode_bool(targetP->data.onoff, dataP);
		return COAP_205_CONTENT;
	case RES_M_DIMMER_VALUE:
		lwm2m_data_encode_int(targetP->data.dimmer, dataP);
		return COAP_205_CONTENT;
	case RES_M_COLOUR_VALUE:
		lwm2m_data_encode_string(targetP->data.colour, dataP);
		return COAP_205_CONTENT;
	case RES_O_APPLICATION_TYPE:
		lwm2m_data_encode_string(targetP->data.application_type, dataP);
		return COAP_205_CONTENT;
	default:
		return COAP_404_NOT_FOUND;
	}
}

static int read_led_inst_data(uint16_t instanceId, led_info_t** targetP) {
	int ret = 0;
	LedObjInfo* data = NULL;
	size_t szled = 0;
	szled = sizeof(LedObjInfo);
	/* Read Led data */
	data = malloc(szled);
	if (!data) {
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	ret = ereg_read_inst(instanceId, OBJ_TYPE_LED, ALL_RESOURCE_ID, data, &szled);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Led data for instance %d\r\n", instanceId);
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	/* Copy the led data read */
	(*targetP)->data.onoff = data->onoff;
	(*targetP)->data.dimmer = data->dimmer;
	(*targetP)->data.instanceId = data->instanceId;
	(*targetP)->data.ontime = data->ontime;
	(*targetP)->data.cumm_active_pwr = data->cumm_active_pwr;
	(*targetP)->data.pwr_factor = data->pwr_factor;
	strcpy((*targetP)->data.sensor_units, data->sensor_units);
	strcpy((*targetP)->data.colour, data->colour);
	strcpy((*targetP)->data.application_type, data->application_type);

	cleanup:
	if(data) {
		free(data);
		data = NULL;
	}
	return ret;
}

static uint8_t prv_led_info_read(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{

	uint8_t result = 0;
	int i = 0;
	led_info_t * targetP = NULL;
	targetP = (led_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	/* Read Led data for instance */
	if (read_led_inst_data(instanceId, &targetP)) {
		return COAP_500_INTERNAL_SERVER_ERROR;
	}

	// is the server asking for the full instance ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_ONOFF_VALUE,
				RES_M_DIMMER_VALUE,
				RES_M_COLOUR_VALUE,
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

static uint8_t prv_led_info_discover(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	led_info_t * targetP;
	uint8_t result;
	int i;

	result = COAP_205_CONTENT;
	targetP = (led_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full object ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_ONOFF_VALUE,
				RES_M_DIMMER_VALUE,
				RES_M_COLOUR_VALUE,
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
			case RES_M_ONOFF_VALUE:
			case RES_M_DIMMER_VALUE:
			case RES_M_COLOUR_VALUE:
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

static uint8_t prv_led_info_write(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	led_info_t * targetP;
	int i;
	uint8_t result;
	size_t size = sizeof(LedObjInfo);

	targetP = (led_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP)
	{
		return COAP_404_NOT_FOUND;
	}

	i = 0;
	do
	{
		switch (dataArray[i].id)
		{
		case RES_M_ONOFF_VALUE:
			result = objh_set_bool_value(dataArray + i, (bool *)&(targetP->data.onoff));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_LED, &targetP->data, &size);
			}
			break;
		case RES_M_DIMMER_VALUE:
			/*
			result = prv_set_int_value(dataArray + i, (bool *)&(targetP->data.dimmer));
			if (result == COAP_204_CHANGED) {
				result = prv_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_LED, &targetP->data, &size);
			}*/
			result = COAP_405_METHOD_NOT_ALLOWED;
			break;
			/*TODO::TBD We nay have to modify LED working from driver till lwm2m2 client*/
		case RES_M_COLOUR_VALUE:
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

static uint8_t prv_led_info_delete(uint16_t id,
		lwm2m_object_t * objectP)
{
	led_info_t * ledInfoInstance = NULL;
	objectP->instanceList = lwm2m_list_remove(objectP->instanceList, id, (lwm2m_list_t **)&ledInfoInstance);
	if (NULL == ledInfoInstance) return COAP_404_NOT_FOUND;

	lwm2m_free(ledInfoInstance);

	return COAP_202_DELETED;
}

static uint8_t prv_led_info_create(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	led_info_t * ledInfoInstance;
	uint8_t result;

	ledInfoInstance = (led_info_t *)lwm2m_malloc(sizeof(led_info_t));
	if (NULL == ledInfoInstance) return COAP_500_INTERNAL_SERVER_ERROR;
	memset(ledInfoInstance, 0, sizeof(led_info_t));

	ledInfoInstance->data.instanceId = instanceId;
	objectP->instanceList = LWM2M_LIST_ADD(objectP->instanceList, ledInfoInstance);


	result = prv_led_info_write(instanceId, numData, dataArray, objectP);

	if (result != COAP_204_CHANGED)
	{
		(void)prv_led_info_delete(instanceId, objectP);
	}
	else
	{
		result = COAP_201_CREATED;
	}

	return result;
}

void display_led_info_object(lwm2m_object_t * object)
{
#ifdef WITH_LOGS
	fprintf(stdout, "  /%u: Led Info object, instances:\r\n", object->objID);
	led_info_t * ledInfoInstance = (led_info_t *)object->instanceList;
	while (ledInfoInstance != NULL)
	{
		fprintf(stdout, "    /%u/%u: instanceId: %u, state value: %f",
				object->objID, ledInfoInstance->data.instanceId,
				ledInfoInstance->data.instanceId, ledInfoInstance->data.onoff);
		fprintf(stdout, "\r\n");
		ledInfoInstance = (led_info_t *)ledInfoInstance->next;
	}
#endif
}

lwm2m_object_t * create_led_info_object()
{
	lwm2m_object_t * ledInfoObj = NULL;
	ledInfoObj = (lwm2m_object_t *)lwm2m_malloc(sizeof(lwm2m_object_t));
	if (NULL != ledInfoObj)
	{
		memset(ledInfoObj, 0, sizeof(lwm2m_object_t));
		ledInfoObj->objID = OBJECT_ID_LED;
		/* set the callbacks. */
		ledInfoObj->readFunc = prv_led_info_read;
		ledInfoObj->discoverFunc = prv_led_info_discover;
		ledInfoObj->writeFunc = prv_led_info_write;
		ledInfoObj->createFunc = NULL;
		ledInfoObj->deleteFunc = prv_led_info_delete;
		ledInfoObj->executeFunc = NULL;
	}
	return ledInfoObj;
}

led_info_t * create_led_info_instance(uint16_t instance)
{
	led_info_t * ledInfoInstance = NULL;

	/* allocate memory for module info object instance. */
	ledInfoInstance = (led_info_t *)lwm2m_malloc(sizeof(led_info_t));
	if (NULL == ledInfoInstance)
	{
		return NULL;
	}
	memset(ledInfoInstance, 0, sizeof(led_info_t));

	/* Read Led data for instance */
	if (read_led_inst_data(instance, &ledInfoInstance)) {
		if (ledInfoInstance) {
			lwm2m_free(ledInfoInstance);
			ledInfoInstance = NULL;
		}
	}

	return ledInfoInstance;
}

lwm2m_object_t * get_led_info_object()
{
	int ret = 0;
	lwm2m_object_t * ledInfoObj = create_led_info_object();
	if (ledInfoObj == NULL)
	{
		fprintf(stderr, "Failed to create led info object\r\n");
		return NULL;
	}

	int *count = lwm2m_malloc(sizeof(int));
	if (!count) {
		fprintf(stderr, "Failed to allocate memory for led sensor count\r\n");
		lwm2m_free(ledInfoObj);
		ledInfoObj = NULL;
		goto cleanup;
	}

	size_t szcount = sizeof(int);
	ret = ereg_read_inst_count(OBJ_TYPE_LED, count, &szcount);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Led sensor count\r\n");
		lwm2m_free(ledInfoObj);
		ledInfoObj = NULL;
		goto cleanup;
	}

	/* Create instances for LedInfo Object. */
	for (uint16_t iter = 0; iter < *count; iter++)
	{
		led_info_t * ledInfoInstance = create_led_info_instance(iter);
		if (ledInfoInstance == NULL)
		{
			fprintf(stderr, "Failed to create Led info instance\r\n");
			lwm2m_free(ledInfoObj);
			ledInfoObj = NULL;
			goto cleanup;
		}
		/* add the led sensor instance to the Led info object. */
		ledInfoObj->instanceList = LWM2M_LIST_ADD(ledInfoObj->instanceList, ledInfoInstance);
	}

	cleanup:
	if (count) {
		lwm2m_free(count);
		count = NULL;
	}

	return ledInfoObj;
}
