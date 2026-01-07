/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>

#include "initClient.h"
#include "jserdes.h"
#include "log.h"

/* JSON (de)-serialization functions. */

int serialize_request(Request *request, json_t **json) {

	int ret=FALSE;
	Register *reg;
	char *str=NULL;
	log_info("serialize_request for %d", request->reqType);
	if ((request->reqType == (ReqType)REQ_REGISTER) || (request->reqType == (ReqType)REQ_UPDATE)) {

		*json = json_object();
		if (*json == NULL) {
			log_error("JSON object null");
			return FALSE;
		}

		reg = request->reg;

		json_object_set_new(*json, JSON_IP,          json_string(reg->ip));
		json_object_set_new(*json, JSON_PORT,        json_integer(atoi(reg->port)));
		json_object_set_new(*json, JSON_CERTIFICATE, json_string(reg->cert));
		json_object_set_new(*json, JSON_NODE_GW_IP,  json_string(reg->nodeGWip));
		json_object_set_new(*json, JSON_PORT,        json_integer(atoi(reg->nodeGWport)));

		str = json_dumps(*json, 0);
		if (str) {
			log_debug("Registration JSON: %s", str);
			free(str);
		}
		ret = TRUE;
	} else if (request->reqType == (ReqType)REQ_UNREGISTER) {
	  ret = TRUE;
	} else if (request->reqType == (ReqType)REQ_QUERY) {
	  ret = TRUE;
	} else if (request->reqType == (ReqType)REQ_QUERY_SYSTEM) {
		ret = TRUE;
	}

	return ret;
}

int deserialize_response(ReqType reqType, QueryResponse **queryResponse,
						 char *str) {

	int ret=TRUE;
	json_t *json=NULL;
	json_t *name, *id, *cert, *ip, *port, *health;

	if (str == NULL) return FALSE;

	json = json_loads(str, JSON_DECODE_ANY, NULL);
	if (!json) {
		log_error("Can no load str into JSON object. Str: %s", str);
		return FALSE;
	}

	name   = json_object_get(json, JSON_SYSTEM_NAME);
	id     = json_object_get(json, JSON_SYSTEM_ID);
	cert   = json_object_get(json, JSON_CERTIFICATE);
	ip     = json_object_get(json, JSON_IP);
	port   = json_object_get(json, JSON_PORT);

	if (reqType == (ReqType)REQ_QUERY) {
		health = json_object_get(json, JSON_HEALTH);
	}

	if (!name || !id || !cert || !ip || !port) {
		log_error("Error deserilaizing response");
		ret = FALSE;
		goto failure;
	}

	if (reqType == (ReqType)REQ_QUERY && !health) {
		log_error("Error deserilaizing response");
		ret = FALSE;
		goto failure;
	}

	*queryResponse = (QueryResponse *)calloc(1, sizeof(QueryResponse));
	if (*queryResponse == NULL) {
		log_error("Memory allocation error of size: %ld",
				  sizeof(QueryResponse));
		ret = FALSE;
		goto failure;
	}

	(*queryResponse)->systemName  = strdup(json_string_value(name));
	(*queryResponse)->systemID    = strdup(json_string_value(id));
	(*queryResponse)->certificate = strdup(json_string_value(cert));
	(*queryResponse)->ip          = strdup(json_string_value(ip));
	(*queryResponse)->port        = json_integer_value(port);

	if (reqType == (ReqType)REQ_QUERY) {
		(*queryResponse)->health      = json_integer_value(health);
	}

 failure:
	json_decref(json);
	return ret;
}

int serialize_uuids_from_file(SystemRegistrationId *sysReg, json_t **json) {

	char *str=NULL;
	if (!sysReg) {
		return FALSE;
	}

	*json = json_object();
	if (*json == NULL) {
		return FALSE;
	}

	if (sysReg->globalUUID) {
		json_object_set_new(*json, JSON_GLOBAL_UUID,   json_string(sysReg->globalUUID));
	}

	if (sysReg->localUUID) {
		json_object_set_new(*json, JSON_LOCAL_UUID, json_string(sysReg->localUUID));
	}

	str = json_dumps(*json, 0);
	if (str) {
		log_debug("System Registration JSON: %s", str);
		free(str);
	}


	return TRUE;

}

int deserialize_uuids_from_file(char* str, SystemRegistrationId** sysReg) {

	int ret=TRUE;
	json_t *json=NULL;
	json_t *gUUID, *lUUID;

	json = json_loads(str, JSON_DECODE_ANY, NULL);
	if (!json) {
		log_error("Can not load str into JSON object. Str: %s", str);
		return FALSE;
	}

	gUUID = json_object_get(json, JSON_GLOBAL_UUID);
	lUUID = json_object_get(json, JSON_LOCAL_UUID);

	*sysReg = (SystemRegistrationId *)calloc(1, sizeof(SystemRegistrationId));
	if (*sysReg == NULL) {
		log_error("Memory allocation error of size: %ld",
				  sizeof(SystemRegistrationId));
		ret = FALSE;
		goto failure;
	}

	if (gUUID) {
		(*sysReg)->globalUUID  = strdup(json_string_value(gUUID));
	}

	if (lUUID) {
		(*sysReg)->localUUID  = strdup(json_string_value(lUUID));
	}

failure:
    json_decref(json);
    return ret;
}
