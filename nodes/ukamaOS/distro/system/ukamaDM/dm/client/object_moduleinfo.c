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
#include "objects/module.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>

static uint8_t prv_get_value(lwm2m_data_t * dataP,
		module_info_t * targetP)
{
	switch (dataP->id)
	{
	case RES_M_MOD_UUID:
		lwm2m_data_encode_string(targetP->data.uuid, dataP);
		return COAP_205_CONTENT;
	case RES_M_MOD_NAME:
		lwm2m_data_encode_string(targetP->data.name, dataP);
		return COAP_205_CONTENT;
	case RES_M_MOD_CLASS:
		lwm2m_data_encode_int(targetP->data.class, dataP);
		return COAP_205_CONTENT;
	case RES_M_PART_NUMBER:
		lwm2m_data_encode_string(targetP->data.partnumber, dataP);
		return COAP_205_CONTENT;
	case RES_M_MOD_MFGNAME:
		lwm2m_data_encode_string(targetP->data.mfgname, dataP);
		return COAP_205_CONTENT;
	case RES_M_MOD_MFGDATE:
		lwm2m_data_encode_string(targetP->data.mfgdate, dataP);
		return COAP_205_CONTENT;
	case RES_M_MOD_MAC:
		lwm2m_data_encode_string(targetP->data.mac, dataP);
		return COAP_205_CONTENT;
	case RES_M_MOD_SW_VERSION:
		lwm2m_data_encode_string(targetP->data.sw_version, dataP);
		return COAP_205_CONTENT;
	case RES_M_MOD_PSW_VERSION:
		lwm2m_data_encode_string(targetP->data.psw_version, dataP);
		return COAP_205_CONTENT;
	case RES_M_MOD_HW_VERSION:
		lwm2m_data_encode_string(targetP->data.hw_version, dataP);
		return COAP_205_CONTENT;
	case RES_M_MOD_DEV_COUNT:
		lwm2m_data_encode_int(targetP->data.device_count, dataP);
		return COAP_205_CONTENT;
	default:
		return COAP_404_NOT_FOUND;
	}
}

static uint8_t prv_module_info_read(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	module_info_t * targetP;
	uint8_t result;
	int i;

	targetP = (module_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full instance ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_MOD_UUID,
				RES_M_MOD_NAME,
				RES_M_MOD_CLASS,
				RES_M_PART_NUMBER,
				RES_M_MOD_MFGNAME,
				RES_M_MOD_MFGDATE,
				RES_M_MOD_MAC,
				RES_M_MOD_SW_VERSION,
				RES_M_MOD_PSW_VERSION,
				RES_M_MOD_HW_VERSION,
				RES_M_MOD_DEV_COUNT
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

static uint8_t prv_module_info_discover(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	module_info_t * targetP;
	uint8_t result;
	int i;

	result = COAP_205_CONTENT;

	targetP = (module_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full object ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_MOD_UUID,
				RES_M_MOD_NAME,
				RES_M_MOD_CLASS,
				RES_M_PART_NUMBER,
				RES_M_MOD_MFGNAME,
				RES_M_MOD_MFGDATE,
				RES_M_MOD_MAC,
				RES_M_MOD_SW_VERSION,
				RES_M_MOD_PSW_VERSION,
				RES_M_MOD_HW_VERSION,
				RES_M_MOD_DEV_COUNT
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
			case RES_M_MOD_UUID:
			case RES_M_MOD_NAME:
			case RES_M_MOD_CLASS:
			case RES_M_PART_NUMBER:
			case RES_M_MOD_MFGNAME:
			case RES_M_MOD_MFGDATE:
			case RES_M_MOD_MAC:
			case RES_M_MOD_SW_VERSION:
			case RES_M_MOD_PSW_VERSION:
			case RES_M_MOD_HW_VERSION:
			case RES_M_MOD_DEV_COUNT:
				break;

			default:
				result = COAP_404_NOT_FOUND;
				break;
			}
		}
	}

	return result;
}

static uint8_t prv_module_info_write(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	module_info_t * targetP;
	int i;
	uint8_t result;

	targetP = (module_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP)
	{
		return COAP_404_NOT_FOUND;
	}

	i = 0;
	do
	{
		switch (dataArray[i].id)
		{
		case RES_M_MOD_UUID:
		case RES_M_MOD_NAME:
		case RES_M_MOD_CLASS:
		case RES_M_PART_NUMBER:
		case RES_M_MOD_MFGNAME:
		case RES_M_MOD_MFGDATE:
		case RES_M_MOD_MAC:
		case RES_M_MOD_SW_VERSION:
		case RES_M_MOD_PSW_VERSION:
		case RES_M_MOD_HW_VERSION:
		case RES_M_MOD_DEV_COUNT:
			result = COAP_405_METHOD_NOT_ALLOWED;
			break;
		default:
			return COAP_404_NOT_FOUND;
		}
		i++;
	} while (i < numData && result == COAP_204_CHANGED);

	return result;
}

static uint8_t prv_module_info_delete(uint16_t id,
		lwm2m_object_t * objectP)
{
	module_info_t * moduleInfoInstance;

	objectP->instanceList = lwm2m_list_remove(objectP->instanceList, id, (lwm2m_list_t **)&moduleInfoInstance);
	if (NULL == moduleInfoInstance) return COAP_404_NOT_FOUND;

	lwm2m_free(moduleInfoInstance);

	return COAP_202_DELETED;
}

static uint8_t prv_module_info_create(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	module_info_t * moduleInfoInstance;
	uint8_t result;

	moduleInfoInstance = (module_info_t *)lwm2m_malloc(sizeof(module_info_t));
	if (NULL == moduleInfoInstance) return COAP_500_INTERNAL_SERVER_ERROR;
	memset(moduleInfoInstance, 0, sizeof(module_info_t));

	moduleInfoInstance->data.instanceId = instanceId;
	objectP->instanceList = LWM2M_LIST_ADD(objectP->instanceList, moduleInfoInstance);

	//todo: will not be able to create objects as module info resources are all read only. not sure if object instance can be created using coap calls.
	result = prv_module_info_write(instanceId, numData, dataArray, objectP);

	if (result != COAP_204_CHANGED)
	{
		(void)prv_module_info_delete(instanceId, objectP);
	}
	else
	{
		result = COAP_201_CREATED;
	}

	return result;
}


void display_module_info_object(lwm2m_object_t * object)
{
#ifdef WITH_LOGS
	fprintf(stdout, "  /%u: Module Info object, instances:\r\n", object->objID);
	module_info_t * moduleInfoInstance = (module_info_t *)object->instanceList;
	while (moduleInfoInstance != NULL)
	{
		fprintf(stdout, "    /%u/%u: instanceId: %u, uuid: %s, manufacturer: %s",
				object->objID, moduleInfoInstance->data.instanceId,
				moduleInfoInstance->data.instanceId, moduleInfoInstance->data.uuid, moduleInfoInstance->data.manufacturer);
		fprintf(stdout, "\r\n");
		moduleInfoInstance = (module_info_t *)moduleInfoInstance->next;
	}
#endif
}

lwm2m_object_t * create_module_info_object()
{
	lwm2m_object_t * moduleInfoObj = NULL;
	moduleInfoObj = (lwm2m_object_t *)lwm2m_malloc(sizeof(lwm2m_object_t));
	if (NULL != moduleInfoObj)
	{
		memset(moduleInfoObj, 0, sizeof(lwm2m_object_t));
		/* ID as per device management document on google docs. */
		moduleInfoObj->objID = OBJECT_ID_MODULE;
		/* set the callbacks. */
		moduleInfoObj->readFunc = prv_module_info_read;
		moduleInfoObj->discoverFunc = prv_module_info_discover;
		moduleInfoObj->writeFunc = NULL;
		moduleInfoObj->createFunc = NULL;
		moduleInfoObj->deleteFunc = prv_module_info_delete;
		moduleInfoObj->executeFunc = NULL;
	}
	return moduleInfoObj;
}

module_info_t * create_module_info_instance(uint16_t inst)
{
	int ret = 0;
	ModuleObjInfo* data = NULL;
	module_info_t * moduleInfoInstance = NULL;
	size_t szmod = sizeof(ModuleObjInfo);
	data = lwm2m_malloc(szmod);
	if (!data) {
		goto cleanup;
	}
	memset(data, '\0',szmod);

	/* Read Module Info from the UkamaEDR */
	ret = ereg_read_inst(inst, OBJ_TYPE_MOD, ALL_RESOURCE_ID, data, &szmod);
	if (ret) {
		fprintf(stderr, "Err(%d) Failed to read Unit Info from the UkamaEDR.", ret);
		goto cleanup;
	}

	/* allocate memory for module info object instance. */
	moduleInfoInstance = (module_info_t *)lwm2m_malloc(sizeof(module_info_t));
	if (NULL == moduleInfoInstance)
	{
		goto cleanup;
	}
	memset(moduleInfoInstance, 0, sizeof(module_info_t));

	/* Copy Data to Unit Info Instance */
	moduleInfoInstance->data.instanceId = data->instanceId;
	strcpy (moduleInfoInstance->data.uuid, data->uuid);
	strcpy (moduleInfoInstance->data.name, data->name);
	moduleInfoInstance->data.class = data->class;
	strcpy (moduleInfoInstance->data.partnumber, data->partnumber);
	strcpy (moduleInfoInstance->data.mfgname, data->mfgname);
	strcpy (moduleInfoInstance->data.mfgdate, data->mfgdate);
	strcpy (moduleInfoInstance->data.mac, data->mac);
	strcpy (moduleInfoInstance->data.sw_version, data->sw_version);
	strcpy (moduleInfoInstance->data.psw_version, data->psw_version);
	strcpy (moduleInfoInstance->data.hw_version, data->hw_version);
	moduleInfoInstance->data.device_count = data->device_count;

	cleanup:
	if(data) {
		lwm2m_free(data);
		data = NULL;
	}
	return moduleInfoInstance;
}

lwm2m_object_t * get_module_info_object()
{
	int ret = 0;
	lwm2m_object_t * moduleInfoObj = create_module_info_object();
	if (moduleInfoObj == NULL)
	{
		fprintf(stderr, "Failed to create module info object\r\n");
		lwm2m_free(moduleInfoObj);
		moduleInfoObj = NULL;
		goto cleanup;
	}

	int *count = lwm2m_malloc(sizeof(int));
	if (!count) {
		fprintf(stderr, "Failed to allocate memory for module count\r\n");
		lwm2m_free(moduleInfoObj);
		moduleInfoObj = NULL;
		goto cleanup;
	}

	size_t szcount = sizeof(int);
	ret = ereg_read_inst_count(OBJ_TYPE_MOD, count, &szcount);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve module count\r\n");
		lwm2m_free(moduleInfoObj);
		moduleInfoObj = NULL;
		goto cleanup;
	}

	/* Create instances for ModuleInfo Object. */
	for (uint16_t iter = 0; iter < *count; iter++)
	{
		module_info_t * moduleInfoInstance = create_module_info_instance(iter);
		if (moduleInfoInstance == NULL)
		{
			fprintf(stderr, "Failed to create module info instance\r\n");
			lwm2m_free(moduleInfoObj);
			moduleInfoObj = NULL;
			goto cleanup;
		}
		/* add the module info instance to the module info object. */
		moduleInfoObj->instanceList = LWM2M_LIST_ADD(moduleInfoObj->instanceList, moduleInfoInstance);
	}

	cleanup:
	if (count) {
		lwm2m_free(count);
		count = NULL;
	}
	return moduleInfoObj;
}


