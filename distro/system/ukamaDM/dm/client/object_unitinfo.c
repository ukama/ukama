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
#include "objects/unit.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>

static uint8_t prv_get_value(lwm2m_data_t * dataP,
		unit_info_t * targetP)
{
	switch (dataP->id)
	{
	case RES_M_UNIT_UUID:
		lwm2m_data_encode_string(targetP->data.uuid, dataP);
		return COAP_205_CONTENT;
	case RES_M_UNIT_NAME:
		lwm2m_data_encode_string(targetP->data.name, dataP);
		return COAP_205_CONTENT;
	case RES_M_UNIT_CLASS:
		lwm2m_data_encode_int(targetP->data.class, dataP);
		return COAP_205_CONTENT;
	case RES_M_SKEW:
		lwm2m_data_encode_string(targetP->data.skew, dataP);
		return COAP_205_CONTENT;
	case RES_M_UNIT_OEMNAME:
		lwm2m_data_encode_string(targetP->data.oemname, dataP);
		return COAP_205_CONTENT;
	case RES_M_UNIT_ASMDATE:
		lwm2m_data_encode_string(targetP->data.asmdate, dataP);
		return COAP_205_CONTENT;
	case RES_M_UNIT_MAC:
		lwm2m_data_encode_string(targetP->data.mac, dataP);
		return COAP_205_CONTENT;
	case RES_M_UNIT_SW_VERSION:
		lwm2m_data_encode_string(targetP->data.sw_version, dataP);
		return COAP_205_CONTENT;
	case RES_M_UNIT_PSW_VERSION:
		lwm2m_data_encode_string(targetP->data.psw_version, dataP);
		return COAP_205_CONTENT;
	case RES_M_UNIT_MOD_COUNT:
		lwm2m_data_encode_int(targetP->data.module_count, dataP);
		return COAP_205_CONTENT;
	default:
		return COAP_404_NOT_FOUND;
	}
}

static uint8_t prv_unit_info_read(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	unit_info_t * targetP;
	uint8_t result;
	int i;

	targetP = (unit_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full instance ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_UNIT_UUID,
				RES_M_UNIT_NAME,
				RES_M_UNIT_CLASS,
				RES_M_SKEW,
				RES_M_UNIT_OEMNAME,
				RES_M_UNIT_ASMDATE,
				RES_M_UNIT_MAC,
				RES_M_UNIT_SW_VERSION,
				RES_M_UNIT_PSW_VERSION,
				RES_M_UNIT_MOD_COUNT
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

static uint8_t prv_unit_info_discover(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	unit_info_t * targetP;
	uint8_t result;
	int i;

	result = COAP_205_CONTENT;

	targetP = (unit_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full object ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_UNIT_UUID,
				RES_M_UNIT_NAME,
				RES_M_UNIT_CLASS,
				RES_M_SKEW,
				RES_M_UNIT_OEMNAME,
				RES_M_UNIT_ASMDATE,
				RES_M_UNIT_MAC,
				RES_M_UNIT_SW_VERSION,
				RES_M_UNIT_PSW_VERSION,
				RES_M_UNIT_MOD_COUNT
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
			case RES_M_UNIT_UUID:
			case RES_M_UNIT_NAME:
			case RES_M_UNIT_CLASS:
			case RES_M_SKEW:
			case RES_M_UNIT_OEMNAME:
			case RES_M_UNIT_ASMDATE:
			case RES_M_UNIT_MAC:
			case RES_M_UNIT_SW_VERSION:
			case RES_M_UNIT_PSW_VERSION:
			case RES_M_UNIT_MOD_COUNT:
				break;

			default:
				result = COAP_404_NOT_FOUND;
				break;
			}
		}
	}

	return result;
}

static uint8_t prv_unit_info_write(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	unit_info_t * targetP;
	int i;
	uint8_t result;

	targetP = (unit_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP)
	{
		return COAP_404_NOT_FOUND;
	}

	i = 0;
	do
	{
		switch (dataArray[i].id)
		{
		case RES_M_UNIT_UUID:
		case RES_M_UNIT_NAME:
		case RES_M_UNIT_CLASS:
		case RES_M_SKEW:
		case RES_M_UNIT_OEMNAME:
		case RES_M_UNIT_ASMDATE:
		case RES_M_UNIT_MAC:
		case RES_M_UNIT_SW_VERSION:
		case RES_M_UNIT_PSW_VERSION:
		case RES_M_UNIT_MOD_COUNT:
			result = COAP_405_METHOD_NOT_ALLOWED;
			break;
		default:
			return COAP_404_NOT_FOUND;
		}
		i++;
	} while (i < numData && result == COAP_204_CHANGED);

	return result;
}

static uint8_t prv_unit_info_delete(uint16_t id,
		lwm2m_object_t * objectP)
{
	unit_info_t * unitInfoInstance;

	objectP->instanceList = lwm2m_list_remove(objectP->instanceList, id, (lwm2m_list_t **)&unitInfoInstance);
	if (NULL == unitInfoInstance) return COAP_404_NOT_FOUND;

	lwm2m_free(unitInfoInstance);

	return COAP_202_DELETED;
}

static uint8_t prv_unit_info_create(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	unit_info_t * unitInfoInstance;
	uint8_t result;

	unitInfoInstance = (unit_info_t *)lwm2m_malloc(sizeof(unit_info_t));
	if (NULL == unitInfoInstance) return COAP_500_INTERNAL_SERVER_ERROR;
	memset(unitInfoInstance, 0, sizeof(unit_info_t));

	unitInfoInstance->data.instanceId = instanceId;
	objectP->instanceList = LWM2M_LIST_ADD(objectP->instanceList, unitInfoInstance);

	//todo: will not be able to create objects as module info resources are all read only. not sure if object instance can be created using coap calls.
	result = prv_unit_info_write(instanceId, numData, dataArray, objectP);

	if (result != COAP_204_CHANGED)
	{
		(void)prv_unit_info_delete(instanceId, objectP);
	}
	else
	{
		result = COAP_201_CREATED;
	}

	return result;
}


void display_unit_info_object(lwm2m_object_t * object)
{
#ifdef WITH_LOGS
	fprintf(stdout, "  /%u: Module Info object, instances:\r\n", object->objID);
	unit_info_t * unitInfoInstance = (unit_info_t *)object->instanceList;
	while (unitInfoInstance != NULL)
	{
		fprintf(stdout, "    /%u/%u: instanceId: %u, uuid: %s, manufacturer: %s",
				object->objID, unitInfoInstance->data.instanceId,
				unitInfoInstance->data.instanceId, unitInfoInstance->data.uuid, unitInfoInstance->data.manufacturer);
		fprintf(stdout, "\r\n");
		unitInfoInstance = (unit_info_t *)unitInfoInstance->next;
	}
#endif
}


lwm2m_object_t * create_unit_info_object()
{
	lwm2m_object_t * unitInfoObj = NULL;
	unitInfoObj = (lwm2m_object_t *)lwm2m_malloc(sizeof(lwm2m_object_t));
	if (NULL != unitInfoObj)
	{
		memset(unitInfoObj, 0, sizeof(lwm2m_object_t));
		/* ID as per device management document on google docs. */
		unitInfoObj->objID = OBJECT_ID_UNIT;
		/* set the callbacks. */
		unitInfoObj->readFunc = prv_unit_info_read;
		unitInfoObj->discoverFunc = prv_unit_info_discover;
		unitInfoObj->writeFunc = NULL;
		unitInfoObj->createFunc = NULL;
		unitInfoObj->deleteFunc = prv_unit_info_delete;
		unitInfoObj->executeFunc = NULL;
	}
	return unitInfoObj;
}

unit_info_t * create_unit_info_instance(uint8_t inst)
{
	int ret = 0;
	UnitObjInfo *data = NULL;
	unit_info_t * unitInfoInstance =  NULL;
	size_t szunit = sizeof(UnitObjInfo);
	data = lwm2m_malloc(szunit);
	if (!data) {
		goto cleanup;
	}
	memset(data, '\0',szunit);

	/* Read Unit Info from the UkamaEDR */
	ret = ereg_read_inst(inst, OBJ_TYPE_UNIT, ALL_RESOURCE_ID, data, &szunit);
	if (ret) {
		fprintf(stderr, "Err(%d) Failed to read Unit Info from the UkamaEDR.", ret);
		goto cleanup;
	}

	/* allocate memory for module info object instance. */
	unitInfoInstance = (unit_info_t *)lwm2m_malloc(sizeof(unit_info_t));
	if (NULL == unitInfoInstance)
	{
		goto cleanup;
	}
	memset(unitInfoInstance, '\0', sizeof(unit_info_t));

	/* Copy Data to Unit Info Instance */
	unitInfoInstance->data.instanceId = data->instanceId;
	strcpy (unitInfoInstance->data.uuid, data->uuid);
	strcpy (unitInfoInstance->data.name, data->name);
	unitInfoInstance->data.class = data->class;
	strcpy (unitInfoInstance->data.skew, data->skew);
	strcpy (unitInfoInstance->data.oemname, data->oemname);
	strcpy (unitInfoInstance->data.asmdate, data->asmdate);
	strcpy (unitInfoInstance->data.mac, data->mac);
	strcpy (unitInfoInstance->data.sw_version, data->sw_version);
	strcpy (unitInfoInstance->data.psw_version, data->psw_version);
	unitInfoInstance->data.module_count = data->module_count;
	cleanup:
	if (data) {
		lwm2m_free(data);
		data = NULL;
	}
	return unitInfoInstance;
}

lwm2m_object_t * get_unit_info_object()
{
	lwm2m_object_t * unitInfoObj = create_unit_info_object();
	if (unitInfoObj == NULL)
	{
		fprintf(stderr, "Failed to create module info object\r\n");
		return NULL;
	}

	/* Unit Info is only having single instance */
	unit_info_t * unitInfoInstance = create_unit_info_instance(0);
	if (unitInfoInstance == NULL)
	{
		fprintf(stderr, "Failed to create module info instance\r\n");
		lwm2m_free(unitInfoObj);
		return NULL;
	}

	/* add the module info instance to the module info object. */
	unitInfoObj->instanceList = LWM2M_LIST_ADD(unitInfoObj->instanceList, unitInfoInstance);

	return unitInfoObj;
}
