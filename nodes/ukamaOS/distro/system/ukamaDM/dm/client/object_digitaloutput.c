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
#include "objects/digital_output.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>

static uint8_t prv_get_value(lwm2m_data_t * dataP,
		digital_output_t * targetP)
{
	switch (dataP->id)
	{
	case RES_M_DIGITAL_OUTPUT_STATE:
		lwm2m_data_encode_bool(targetP->data.digital_state, dataP);
		return COAP_205_CONTENT;
	case RES_O_DIGITAL_OUTPUT_POLARITY:
		lwm2m_data_encode_bool(targetP->data.digital_polarity, dataP);
		return COAP_205_CONTENT;
	case RES_O_APPLICATION_TYPE:
		lwm2m_data_encode_string(targetP->data.application_type, dataP);
		return COAP_205_CONTENT;
	default:
		return COAP_404_NOT_FOUND;
	}
}

static int read_dop_inst_data(uint16_t instanceId, digital_output_t** targetP) {
	int ret = 0;
	DopObjInfo* data = NULL;
	size_t szdop = 0;
	szdop = sizeof(DopObjInfo);
	/* Read Dop data */
	data = malloc(szdop);
	if (!data) {
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	ret = ereg_read_inst(instanceId, OBJ_TYPE_DOP, ALL_RESOURCE_ID, data, &szdop);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Dop data for instance %d\r\n", instanceId);
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	/* Copy the dop data read */
	(*targetP)->data.digital_state = data->digital_state;
	(*targetP)->data.digital_polarity = data->digital_polarity;
	(*targetP)->data.instanceId = data->instanceId;
	strcpy((*targetP)->data.application_type, data->application_type);

	cleanup:
	if(data) {
		free(data);
		data = NULL;
	}
	return ret;
}

static uint8_t prv_dop_info_read(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{

	uint8_t result = 0;
	int i = 0;
	digital_output_t * targetP = NULL;
	targetP = (digital_output_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	/* Read Dop data for instance */
	if (read_dop_inst_data(instanceId, &targetP)) {
		return COAP_500_INTERNAL_SERVER_ERROR;
	}

	// is the server asking for the full instance ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_DIGITAL_OUTPUT_STATE,
				RES_O_DIGITAL_OUTPUT_POLARITY,
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

static uint8_t prv_dop_info_discover(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	digital_output_t * targetP;
	uint8_t result;
	int i;

	result = COAP_205_CONTENT;
	targetP = (digital_output_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full object ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_DIGITAL_OUTPUT_STATE,
				RES_O_DIGITAL_OUTPUT_POLARITY,
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
			case RES_M_DIGITAL_OUTPUT_STATE:
			case RES_O_DIGITAL_OUTPUT_POLARITY:
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

static uint8_t prv_dop_info_write(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	digital_output_t * targetP;
	int i;
	uint8_t result;
	size_t size = sizeof(DopObjInfo);

	targetP = (digital_output_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP)
	{
		return COAP_404_NOT_FOUND;
	}

	i = 0;
	do
	{
		switch (dataArray[i].id)
		{
		case RES_M_DIGITAL_OUTPUT_STATE:
			result = objh_set_bool_value(dataArray + i, (bool *)&(targetP->data.digital_state));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_DOP, &targetP->data, &size);
			}
			break;
		case RES_O_DIGITAL_OUTPUT_POLARITY:
			result = objh_set_bool_value(dataArray + i, (bool *)&(targetP->data.digital_polarity));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_DOP, &targetP->data, &size);
			}
			break;
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

static uint8_t prv_dop_info_delete(uint16_t id,
		lwm2m_object_t * objectP)
{
	digital_output_t * dopInfoInstance = NULL;
	objectP->instanceList = lwm2m_list_remove(objectP->instanceList, id, (lwm2m_list_t **)&dopInfoInstance);
	if (NULL == dopInfoInstance) return COAP_404_NOT_FOUND;

	lwm2m_free(dopInfoInstance);

	return COAP_202_DELETED;
}

static uint8_t prv_dop_info_create(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	digital_output_t * dopInfoInstance;
	uint8_t result;

	dopInfoInstance = (digital_output_t *)lwm2m_malloc(sizeof(digital_output_t));
	if (NULL == dopInfoInstance) return COAP_500_INTERNAL_SERVER_ERROR;
	memset(dopInfoInstance, 0, sizeof(digital_output_t));

	dopInfoInstance->data.instanceId = instanceId;
	objectP->instanceList = LWM2M_LIST_ADD(objectP->instanceList, dopInfoInstance);


	result = prv_dop_info_write(instanceId, numData, dataArray, objectP);

	if (result != COAP_204_CHANGED)
	{
		(void)prv_dop_info_delete(instanceId, objectP);
	}
	else
	{
		result = COAP_201_CREATED;
	}

	return result;
}

void display_dop_info_object(lwm2m_object_t * object)
{
#ifdef WITH_LOGS
	fprintf(stdout, "  /%u: Dop Info object, instances:\r\n", object->objID);
	digital_output_t * dopInfoInstance = (digital_output_t *)object->instanceList;
	while (dopInfoInstance != NULL)
	{
		fprintf(stdout, "    /%u/%u: instanceId: %u, state value: %f",
				object->objID, dopInfoInstance->data.instanceId,
				dopInfoInstance->data.instanceId, dopInfoInstance->data.digital_output_state);
		fprintf(stdout, "\r\n");
		dopInfoInstance = (digital_output_t *)dopInfoInstance->next;
	}
#endif
}

lwm2m_object_t * create_dop_info_object()
{
	lwm2m_object_t * dopInfoObj = NULL;
	dopInfoObj = (lwm2m_object_t *)lwm2m_malloc(sizeof(lwm2m_object_t));
	if (NULL != dopInfoObj)
	{
		memset(dopInfoObj, 0, sizeof(lwm2m_object_t));
		dopInfoObj->objID = OBJECT_ID_DIGITAL_OUTPUT;
		/* set the callbacks. */
		dopInfoObj->readFunc = prv_dop_info_read;
		dopInfoObj->discoverFunc = prv_dop_info_discover;
		dopInfoObj->writeFunc = prv_dop_info_write;
		dopInfoObj->createFunc = NULL;
		dopInfoObj->deleteFunc = prv_dop_info_delete;
		dopInfoObj->executeFunc = NULL;
	}
	return dopInfoObj;
}

digital_output_t * create_dop_info_instance(uint16_t instance)
{
	digital_output_t * dopInfoInstance = NULL;

	/* allocate memory for module info object instance. */
	dopInfoInstance = (digital_output_t *)lwm2m_malloc(sizeof(digital_output_t));
	if (NULL == dopInfoInstance)
	{
		return NULL;
	}
	memset(dopInfoInstance, 0, sizeof(digital_output_t));

	/* Read Dop data for instance */
	if (read_dop_inst_data(instance, &dopInfoInstance)) {
		if (dopInfoInstance) {
			lwm2m_free(dopInfoInstance);
			dopInfoInstance = NULL;
		}
	}

	return dopInfoInstance;
}

lwm2m_object_t * get_dop_info_object()
{
	int ret = 0;
	lwm2m_object_t * dopInfoObj = create_dop_info_object();
	if (dopInfoObj == NULL)
	{
		fprintf(stderr, "Failed to create dop info object\r\n");
		return NULL;
	}

	int *count = lwm2m_malloc(sizeof(int));
	if (!count) {
		fprintf(stderr, "Failed to allocate memory for dop sensor count\r\n");
		lwm2m_free(dopInfoObj);
		dopInfoObj = NULL;
		goto cleanup;
	}

	size_t szcount = sizeof(int);
	ret = ereg_read_inst_count(OBJ_TYPE_DOP, count, &szcount);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Dop sensor count\r\n");
		lwm2m_free(dopInfoObj);
		dopInfoObj = NULL;
		goto cleanup;
	}

	/* Create instances for DopInfo Object. */
	for (uint16_t iter = 0; iter < *count; iter++)
	{
		digital_output_t * dopInfoInstance = create_dop_info_instance(iter);
		if (dopInfoInstance == NULL)
		{
			fprintf(stderr, "Failed to create Dop info instance\r\n");
			lwm2m_free(dopInfoObj);
			dopInfoObj = NULL;
			goto cleanup;
		}
		/* add the dop sensor instance to the Dop info object. */
		dopInfoObj->instanceList = LWM2M_LIST_ADD(dopInfoObj->instanceList, dopInfoInstance);
	}

	cleanup:
	if (count) {
		lwm2m_free(count);
		count = NULL;
	}

	return dopInfoObj;
}
