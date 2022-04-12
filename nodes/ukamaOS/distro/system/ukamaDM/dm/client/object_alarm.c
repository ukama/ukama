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
#include "objects/alarm.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>



static uint8_t prv_exec(uint16_t instanceId, uint16_t resourceId,
		uint8_t * buffer, int length, lwm2m_object_t * objectP)
{
	int ret = 0;
	alarm_info_t * targetP = NULL;
	void* data = NULL;
	size_t size = 0;
	targetP = (alarm_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	switch (resourceId)
	{
	case RES_M_AL_CLEAR:
	case RES_M_AL_ENABLE:
	case RES_M_AL_DISABLE:
		ret = ereg_exec_sensor(instanceId, OBJ_TYPE_ALARM, resourceId, data, &size);
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
		alarm_info_t * targetP)
{
	switch (dataP->id)
	{
	case RES_M_AL_EVENTTYPE:
		lwm2m_data_encode_int(targetP->data.eventtype, dataP);
		return COAP_205_CONTENT;
	case RES_M_AL_REALTIME:
		lwm2m_data_encode_bool(targetP->data.realtime, dataP);
		return COAP_205_CONTENT;
	case RES_M_AL_STATE:
		lwm2m_data_encode_uint(targetP->data.state, dataP);
		return COAP_205_CONTENT;
	case RES_O_AL_DESCRIPTION:
		lwm2m_data_encode_string(targetP->data.disc, dataP);
		return COAP_205_CONTENT;
	case RES_M_AL_LOW_THRESHOLD:
		lwm2m_data_encode_float(targetP->data.lowthreshold, dataP);
		return COAP_205_CONTENT;
	case RES_M_AL_HIGH_THRESHOLD:
		lwm2m_data_encode_float(targetP->data.highthreshold, dataP);
		return COAP_205_CONTENT;
	case RES_M_AL_CRIT_THRESHOLD:
		lwm2m_data_encode_float(targetP->data.crithreshold, dataP);
		return COAP_205_CONTENT;
	case RES_M_AL_EVT_COUNT:
		lwm2m_data_encode_int(targetP->data.eventcount, dataP);
		return COAP_205_CONTENT;
	case RES_M_AL_RECRD_TIME:
		lwm2m_data_encode_uint(targetP->data.time, dataP);
		return COAP_205_CONTENT;
	case RES_M_AL_OBJ_ID:
		lwm2m_data_encode_uint(targetP->data.sobjid, dataP);
		return COAP_205_CONTENT;
	case RES_M_AL_INST_ID:
		lwm2m_data_encode_uint(targetP->data.sinstid, dataP);
		return COAP_205_CONTENT;
	case RES_M_AL_RSRC_ID:
		lwm2m_data_encode_uint(targetP->data.srsrcid, dataP);
		return COAP_205_CONTENT;
	case RES_M_SENSOR_VALUE:
		lwm2m_data_encode_float(targetP->data.sensorvalue, dataP);
		return COAP_205_CONTENT;
	case RES_M_SENSOR_UNITS:
		lwm2m_data_encode_string(targetP->data.sensorunits, dataP);
		return COAP_205_CONTENT;
	case RES_O_APPLICATION_TYPE:
		lwm2m_data_encode_string(targetP->data.applicationtype, dataP);
		return COAP_205_CONTENT;
	default:
		return COAP_404_NOT_FOUND;
	}
}

static int update_instance() {

}
static int read_alarm_inst_data(uint16_t instanceId, alarm_info_t** targetP) {
	int ret = 0;
	AlarmObjInfo* data = NULL;
	size_t szalarm = 0;
	szalarm = sizeof(AlarmObjInfo);
	/* Read Alarm data */
	data = malloc(szalarm);
	if (!data) {
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	ret = ereg_read_inst(instanceId, OBJ_TYPE_ALARM, ALL_RESOURCE_ID, data, &szalarm);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Alarm data for instance %d\r\n", instanceId);
		ret = COAP_500_INTERNAL_SERVER_ERROR;
		goto cleanup;
	}

	/* Copy the alarm data read */
	(*targetP)->data.eventtype = data->eventtype;
	(*targetP)->data.realtime = data->realtime;
	(*targetP)->data.state = data->state;
	(*targetP)->data.lowthreshold = data->lowthreshold;
	(*targetP)->data.highthreshold = data->highthreshold;
	(*targetP)->data.crithreshold = data->crithreshold;
	(*targetP)->data.eventcount = data->eventcount;
	(*targetP)->data.time = data->time;
	(*targetP)->data.sobjid = data->sobjid;
	(*targetP)->data.sinstid = data->sinstid;
	(*targetP)->data.srsrcid = data->srsrcid;
	(*targetP)->data.sensorvalue = data->sensorvalue;
	(*targetP)->data.instanceId = data->instanceId;
	if(!strcmp(data->disc,"")){
		strcpy((*targetP)->data.disc, "No Data.");
	} else {
		strcpy((*targetP)->data.disc, data->disc);
	}
	strcpy((*targetP)->data.sensorunits, data->sensorunits);
	strcpy((*targetP)->data.applicationtype, data->applicationtype);

	cleanup:
	if(data) {
		free(data);
		data = NULL;
	}
	return ret;
}

static uint8_t prv_alarm_info_read(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{

	uint8_t result = 0;
	int i = 0;
	alarm_info_t * targetP = NULL;
	targetP = (alarm_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	/* Read Alarm data for instance */
	if (read_alarm_inst_data(instanceId, &targetP)) {
		return COAP_500_INTERNAL_SERVER_ERROR;
	}

	// is the server asking for the full instance ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_AL_EVENTTYPE,
				RES_M_AL_REALTIME,
				RES_M_AL_STATE,
				RES_M_AL_LOW_THRESHOLD,
				RES_M_AL_HIGH_THRESHOLD,
				RES_M_AL_CRIT_THRESHOLD,
				RES_M_AL_EVT_COUNT,
				RES_M_AL_RECRD_TIME,
				RES_M_AL_OBJ_ID,
				RES_M_AL_INST_ID,
				RES_M_AL_RSRC_ID,
				RES_O_AL_DESCRIPTION,
				RES_M_SENSOR_VALUE,
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

	fprintf(stdout, "Return for reading alarm value is %d\r\n", result);
	fflush(stdout);
	return result;
}

static uint8_t prv_alarm_info_discover(uint16_t instanceId,
		int * numDataP,
		lwm2m_data_t ** dataArrayP,
		lwm2m_object_t * objectP)
{
	alarm_info_t * targetP;
	uint8_t result;
	int i;

	result = COAP_205_CONTENT;
	targetP = (alarm_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) return COAP_404_NOT_FOUND;

	// is the server asking for the full object ?
	if (*numDataP == 0)
	{
		uint16_t resList[] = {
				RES_M_AL_EVENTTYPE,
				RES_M_AL_REALTIME,
				RES_M_AL_STATE,
				RES_M_AL_LOW_THRESHOLD,
				RES_M_AL_HIGH_THRESHOLD,
				RES_M_AL_CRIT_THRESHOLD,
				RES_M_AL_EVT_COUNT,
				RES_M_AL_RECRD_TIME,
				RES_M_AL_OBJ_ID,
				RES_M_AL_INST_ID,
				RES_M_AL_RSRC_ID,
				RES_O_AL_DESCRIPTION,
				RES_M_SENSOR_VALUE,
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
			case RES_M_AL_EVENTTYPE:
			case RES_M_AL_REALTIME:
			case RES_M_AL_STATE:
			case RES_M_AL_LOW_THRESHOLD:
			case RES_M_AL_HIGH_THRESHOLD:
			case RES_M_AL_CRIT_THRESHOLD:
			case RES_M_AL_EVT_COUNT:
			case RES_M_AL_RECRD_TIME:
			case RES_M_AL_OBJ_ID:
			case RES_M_AL_INST_ID:
			case RES_M_AL_RSRC_ID:
			case RES_O_AL_DESCRIPTION:
			case RES_M_SENSOR_VALUE:
			case RES_M_SENSOR_UNITS:
			case RES_O_APPLICATION_TYPE:
			case RES_M_AL_CLEAR:
			case RES_M_AL_ENABLE:
			case RES_M_AL_DISABLE:
				break;
			default:
				result = COAP_404_NOT_FOUND;
				break;
			}
		}
	}

	return result;
}

static uint8_t prv_alarm_info_write(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	alarm_info_t * targetP;
	int i;
	uint8_t result;
	size_t size = sizeof(AlarmObjInfo);

	targetP = (alarm_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP)
	{
		return COAP_404_NOT_FOUND;
	}

	i = 0;
	do
	{
		switch (dataArray[i].id)
		{
		case RES_M_AL_STATE:
		case RES_M_AL_LOW_LIMIT:
		case RES_M_AL_HIGH_LIMIT:
		case RES_M_AL_CRIT_LIMIT:
		case RES_M_AL_EVT_COUNT:
		case RES_M_AL_RECRD_TIME:
		case RES_M_AL_OBJ_ID:
		case RES_M_AL_INST_ID:
		case RES_M_AL_RSRC_ID:
		case RES_O_AL_DESCRIPTION:
		case RES_M_SENSOR_VALUE:
		case RES_M_SENSOR_UNITS:
		case RES_O_APPLICATION_TYPE:
		case RES_M_AL_CLEAR:
		case RES_M_AL_ENABLE:
		case RES_M_AL_DISABLE:
			result = COAP_405_METHOD_NOT_ALLOWED;
			break;
		case RES_M_AL_EVENTTYPE:
			result = objh_set_int_value(dataArray + i, (int *)&(targetP->data.eventtype));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_PWR, &targetP->data, &size);
			}
			break;
		case RES_M_AL_REALTIME:
			result = objh_set_bool_value(dataArray + i, (bool *)&(targetP->data.realtime));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_PWR, &targetP->data, &size);
			}
			break;
		case RES_M_AL_LOW_THRESHOLD:
			result = objh_set_double_value(dataArray + i, (double *)&(targetP->data.lowthreshold));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_PWR, &targetP->data, &size);
			}
			break;
		case RES_M_AL_HIGH_THRESHOLD:
			result = objh_set_double_value(dataArray + i, (double *)&(targetP->data.highthreshold));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_PWR, &targetP->data, &size);
			}
			break;
		case RES_M_AL_CRIT_THRESHOLD:
			result = objh_set_double_value(dataArray + i, (double *)&(targetP->data.crithreshold));
			if (result == COAP_204_CHANGED) {
				result = objh_send_data_ukama_edr(instanceId, (dataArray[i].id), OBJ_TYPE_PWR, &targetP->data, &size);
			}
			break;
		default:
			return COAP_404_NOT_FOUND;
		}
		i++;
	} while (i < numData && result == COAP_204_CHANGED);

	return result;
}

static uint8_t prv_alarm_info_delete(uint16_t id,
		lwm2m_object_t * objectP)
{
	alarm_info_t * alarmInfoInstance = NULL;
	objectP->instanceList = lwm2m_list_remove(objectP->instanceList, id, (lwm2m_list_t **)&alarmInfoInstance);
	if (NULL == alarmInfoInstance) return COAP_404_NOT_FOUND;

	lwm2m_free(alarmInfoInstance);

	return COAP_202_DELETED;
}

static uint8_t prv_alarm_info_create(uint16_t instanceId,
		int numData,
		lwm2m_data_t * dataArray,
		lwm2m_object_t * objectP)
{
	alarm_info_t * alarmInfoInstance;
	uint8_t result;

	alarmInfoInstance = (alarm_info_t *)lwm2m_malloc(sizeof(alarm_info_t));
	if (NULL == alarmInfoInstance) return COAP_500_INTERNAL_SERVER_ERROR;
	memset(alarmInfoInstance, 0, sizeof(alarm_info_t));

	alarmInfoInstance->data.instanceId = instanceId;
	objectP->instanceList = LWM2M_LIST_ADD(objectP->instanceList, alarmInfoInstance);


	result = prv_alarm_info_write(instanceId, numData, dataArray, objectP);

	if (result != COAP_204_CHANGED)
	{
		(void)prv_alarm_info_delete(instanceId, objectP);
	}
	else
	{
		result = COAP_201_CREATED;
	}

	return result;
}

void display_alarm_info_object(lwm2m_object_t * object)
{
#ifdef WITH_LOGS
	fprintf(stdout, "  /%u: Alarm Info object, instances:\r\n", object->objID);
	alarm_info_t * alarmInfoInstance = (alarm_info_t *)object->instanceList;
	while (alarmInfoInstance != NULL)
	{
		fprintf(stdout, "    /%u/%u: instanceId: %u, Alarm %d for sensor value: %f %s",
				object->objID, alarmInfoInstance->data.instanceId,
				alarmInfoInstance->data.instanceId, alarmInfoInstance->data.state,
				alarmInfoInstance->data.sensorvalue, alarmInfoInstance->data.sensorvalue);
		fprintf(stdout, "\r\n");
		alarmInfoInstance = (alarm_info_t *)alarmInfoInstance->next;
	}
#endif
}

lwm2m_object_t * create_alarm_info_object()
{
	lwm2m_object_t * alarmInfoObj = NULL;
	alarmInfoObj = (lwm2m_object_t *)lwm2m_malloc(sizeof(lwm2m_object_t));
	if (NULL != alarmInfoObj)
	{
		memset(alarmInfoObj, 0, sizeof(lwm2m_object_t));
		alarmInfoObj->objID = OBJECT_ID_DEV_ALARM;
		/* set the callbacks. */
		alarmInfoObj->readFunc = prv_alarm_info_read;
		alarmInfoObj->discoverFunc = prv_alarm_info_discover;
		alarmInfoObj->writeFunc = prv_alarm_info_write;
		alarmInfoObj->createFunc = NULL;
		alarmInfoObj->deleteFunc = prv_alarm_info_delete;
		alarmInfoObj->executeFunc = prv_exec;
	}
	return alarmInfoObj;
}

alarm_info_t * create_alarm_info_instance(uint16_t instance)
{
	alarm_info_t * alarmInfoInstance = NULL;

	/* allocate memory for module info object instance. */
	alarmInfoInstance = (alarm_info_t *)lwm2m_malloc(sizeof(alarm_info_t));
	if (NULL == alarmInfoInstance)
	{
		return NULL;
	}
	memset(alarmInfoInstance, 0, sizeof(alarm_info_t));

	/* Read Alarm data for instance */
	if (read_alarm_inst_data(instance, &alarmInfoInstance)) {
		if (alarmInfoInstance) {
			lwm2m_free(alarmInfoInstance);
			alarmInfoInstance = NULL;
		}
	}
	return alarmInfoInstance;
}

lwm2m_object_t * get_alarm_info_object()
{
	int ret = 0;
	lwm2m_object_t * alarmInfoObj = create_alarm_info_object();
	if (alarmInfoObj == NULL)
	{
		fprintf(stderr, "Failed to create alarm info object\r\n");
		return NULL;
	}

	int *count = lwm2m_malloc(sizeof(int));
	if (!count) {
		fprintf(stderr, "Failed to allocate memory for alarm sensor count\r\n");
		lwm2m_free(alarmInfoObj);
		alarmInfoObj = NULL;
		goto cleanup;
	}

	size_t szcount = sizeof(int);
	ret = ereg_read_inst_count(OBJ_TYPE_ALARM, count, &szcount);
	if (ret)
	{
		fprintf(stderr, "Failed to retrieve Alarm sensor count\r\n");
		lwm2m_free(alarmInfoObj);
		alarmInfoObj = NULL;
		goto cleanup;
	}

	/* Create instances for AlarmInfo Object. */
	for (uint16_t iter = 0; iter < *count; iter++)
	{
		alarm_info_t * alarmInfoInstance = create_alarm_info_instance(iter);
		if (alarmInfoInstance == NULL)
		{
			fprintf(stderr, "Failed to create Alarm info instance\r\n");
			lwm2m_free(alarmInfoObj);
			alarmInfoObj = NULL;
			goto cleanup;
		}
		/* add the alarm sensor instance to the Alarm info object. */
		alarmInfoObj->instanceList = LWM2M_LIST_ADD(alarmInfoObj->instanceList, alarmInfoInstance);
	}

	cleanup:
	if (count) {
		lwm2m_free(count);
		count = NULL;
	}

	return alarmInfoObj;
}

int alarm_change(void * pdata,
		lwm2m_object_t * objectP, uint16_t instanceId)
{
	int  ret  = 0;
	alarm_info_t * targetP = NULL;
	targetP = (alarm_info_t *)lwm2m_list_find(objectP->instanceList, instanceId);
	if (NULL == targetP) {
		fprintf(stderr, "Failed to find instance id %d in ObjectP list.\r\n", instanceId);
		fflush(stderr);
		return COAP_404_NOT_FOUND;
	}

	if (pdata) {
		AlarmObjInfo* data = pdata;
		/* Copy the alarm data read */
		(targetP)->data.eventtype = data->eventtype;
		(targetP)->data.realtime = data->realtime;
		(targetP)->data.state = data->state;
		(targetP)->data.lowthreshold = data->lowthreshold;
		(targetP)->data.highthreshold = data->highthreshold;
		(targetP)->data.crithreshold = data->crithreshold;
		(targetP)->data.eventcount = data->eventcount;
		(targetP)->data.time = data->time;
		(targetP)->data.sobjid = data->sobjid;
		(targetP)->data.sinstid = data->sinstid;
		(targetP)->data.srsrcid = data->srsrcid;
		(targetP)->data.sensorvalue = data->sensorvalue;
		//(targetP)->data.instanceId = data->instanceId;
		if(!strcmp(data->disc,"")){
			strcpy((targetP)->data.disc, "No Data.");
		} else {
			strcpy((targetP)->data.disc, data->disc);
		}
		strcpy((targetP)->data.sensorunits, data->sensorunits);
		strcpy((targetP)->data.applicationtype, data->applicationtype);
		fprintf(stdout, "Received alarm for instance id %d\r\n", instanceId);
		fflush(stdout);
	} else {
		fprintf(stderr, "Failed to get alarm data\r\n");
		fflush(stderr);
		ret = COAP_500_INTERNAL_SERVER_ERROR;
	}
	return ret;
}
