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
#include <uuid/uuid.h>

#include "mesh.h"
#include "jserdes.h"
#include "initClient.h"

static void log_json(json_t *json);
static int get_json_entry(json_t *json, char *key, json_type type,
						  char **strValue, int *intValue);

static void log_json(json_t *json) {

	char *str = NULL;

	str = json_dumps(json, JSON_ENCODE_ANY);
	if (str) {
		log_debug("json str: %s", str);
		free(str);
	}
}

static int get_json_entry(json_t *json, char *key, json_type type,
						  char **strValue, int *intValue) {

	json_t *jEntry=NULL;

	if (json == NULL || key == NULL) return FALSE;

	jEntry = json_object_get(json, key);
	if (jEntry == NULL) {
		log_error("Missing %s key in json", key);
		return FALSE;
	}

	if (type == JSON_STRING) {
		*strValue = strdup(json_string_value(jEntry));
	} else if (type == JSON_INTEGER) {
		*intValue = json_integer_value(jEntry);
	} else {
		log_error("Invalid type for json key-value: %d", type);
		return FALSE;
	}

	return TRUE;
}

static void add_map_to_request(json_t **json, UMap *map, int mapType) {

	json_t *jMap=NULL, *jArray=NULL;
	int i;

	if (map == NULL) {
		return;
	}

	if (map->nb_values == 0) {
		return;
	}

	json_object_set_new(*json, JSON_MAP, json_object());
	jMap = json_object_get(*json, JSON_MAP);

	if (mapType == MESH_MAP_TYPE_URL) {
		json_object_set_new(jMap, JSON_TYPE,
							json_string(MESH_MAP_TYPE_URL_STR));
	} else if (mapType == MESH_MAP_TYPE_HDR) {
		json_object_set_new(jMap, JSON_TYPE,
							json_string(MESH_MAP_TYPE_HDR_STR));
	} else if (mapType == MESH_MAP_TYPE_POST) {
		json_object_set_new(jMap, JSON_TYPE,
							json_string(MESH_MAP_TYPE_HDR_STR));
	} else if (mapType == MESH_MAP_TYPE_COOKIE) {
		json_object_set_new(jMap, JSON_TYPE,
							json_string(MESH_MAP_TYPE_COOKIE_STR));
	}

	/* For array of key/value pair. */
	json_object_set_new(jMap, JSON_DATA, json_array());
	jArray = json_object_get(jMap, JSON_DATA);

	if (jArray) {

		for (i=0; i < map->nb_values; i++) {

			json_t *entry = json_object();

			json_object_set_new(entry, JSON_KEY,
								json_string(map->keys[i]));
			json_object_set_new(entry, JSON_VALUE,
								json_string(map->values[i]));
			json_object_set_new(entry, JSON_LEN,
								json_integer((int)map->lengths[i]));

			json_array_append_new(jArray, entry);
		}
	}
}

int serialize_system_response(char **response, Message *message,
                              int retCode, int len, char *data) {

    json_t *json, *obj;

    /* basic sanity check */
	if (len == 0 || data == NULL || message == NULL)
		return FALSE;

	json = json_object();
	if (json == NULL) {
		return FALSE;
	}

	json_object_set_new(json, JSON_TYPE, json_string(UKAMA_NODE_RESPONSE));
    json_object_set_new(json, JSON_UUID, json_string(message->seqNo));

	/* Add response info. */
	json_object_set_new(json, JSON_MESSAGE, json_object());
	obj = json_object_get(json, JSON_MESSAGE);
    json_object_set_new(obj, JSON_CODE, json_integer(retCode));
	json_object_set_new(obj, JSON_LENGTH, json_integer(len));
	json_object_set_new(obj, JSON_DATA, json_string(data));

    *response = json_dumps(json, JSON_ENCODE_ANY);
    json_decref(json);

	return TRUE;
}

static void serialize_message_data(URequest *request, char **data) {

    json_t *json, *jRaw;

    json = json_object();
    if (json == NULL) return;
    
	json_object_set_new(json, JSON_PROTOCOL, json_string(request->http_protocol));
	json_object_set_new(json, JSON_METHOD,   json_string(request->http_verb));
	json_object_set_new(json, JSON_URL,      json_string(request->http_url));
	json_object_set_new(json, JSON_PATH,     json_string(request->url_path));

	/* Add maps if they exists. */
	if (request->map_url) {
		add_map_to_request(&json, request->map_url, MESH_MAP_TYPE_URL);
	}

	if (request->map_header) {
		add_map_to_request(&json, request->map_header, MESH_MAP_TYPE_HDR);
	}

	if (request->map_cookie) {
		add_map_to_request(&json, request->map_cookie, MESH_MAP_TYPE_COOKIE);
	}

	if (request->map_post_body) {
		add_map_to_request(&json, request->map_post_body, MESH_MAP_TYPE_POST);
	}

	/* And finally add raw binary data. Currently we assume raw is char* */
	if (request->binary_body_length > 0 && request->binary_body != NULL ){
		json_object_set_new(json, JSON_RAW_DATA, json_object());
		jRaw = json_object_get(json, JSON_RAW_DATA);
		json_object_set_new(jRaw, JSON_LENGTH,
							json_integer((int)request->binary_body_length));
		json_object_set_new(jRaw, JSON_DATA,
							json_string((char *)request->binary_body));
	}

    *data = json_dumps(json, JSON_ENCODE_ANY);
    json_decref(json);
}

int serialize_websocket_message(char **str,
                                URequest *request,
                                char *uuid) {

    json_t *json=NULL;
	json_t *jRequest=NULL;
    char *data=NULL;

	json = json_object();
	if (json == NULL) {
        *str = NULL;
		return FALSE;
	}

	json_object_set_new(json, JSON_TYPE, json_string(UKAMA_SERVICE_REQUEST));
	json_object_set_new(json, JSON_UUID, json_string(uuid));

    serialize_message_data(request, &data);
	json_object_set_new(json, JSON_MESSAGE, json_object());

	jRequest = json_object_get(json, JSON_MESSAGE);
	json_object_set_new(jRequest, JSON_LENGTH, json_integer(strlen(data)));
	json_object_set_new(jRequest, JSON_DATA,   json_string(data));

    *str = json_dumps(json, JSON_ENCODE_ANY);

    json_decref(json);
    free(data);

	return TRUE;
}

int deserialize_system_info(SystemInfo **systemInfo, json_t *json) {

	int ret=TRUE;

	if (json == NULL) return FALSE;
	log_json(json);

	*systemInfo = (SystemInfo *)calloc(1, sizeof(SystemInfo));
	if (*systemInfo == NULL) {
		log_error("Error allocating memory of size: %lu", sizeof(SystemInfo));
		return FALSE;
	}

	ret |= get_json_entry(json, JSON_SYSTEM_NAME, JSON_STRING,
						  &(*systemInfo)->systemName, NULL);
	ret |= get_json_entry(json, JSON_SYSTEM_ID, JSON_STRING,
						  &(*systemInfo)->systemId, NULL);
	ret |= get_json_entry(json, JSON_CERTIFICATE, JSON_STRING,
						  &(*systemInfo)->certificate, NULL);
	ret |= get_json_entry(json, JSON_NODE_GW_IP, JSON_STRING,
						  &(*systemInfo)->ip, NULL);
	ret |= get_json_entry(json, JSON_NODE_GW_PORT, JSON_INTEGER,
						  NULL, &(*systemInfo)->port);
	ret |= get_json_entry(json, JSON_HEALTH, JSON_INTEGER,
						  NULL, &(*systemInfo)->health);

	if (ret == FALSE) {
		log_error("Error deserializing node info");
		log_json(json);
		free_system_info(*systemInfo);
		*systemInfo = NULL;
	}

	return ret;
}

int deserialize_websocket_message(Message **message, char *data) {

    json_t *json;
    json_t *jType, *jSeq, *jMessage;
    json_t *jLength, *jData, *jCode;

    json = json_loads(data, JSON_DECODE_ANY, NULL);
	if (json == NULL) {
		return FALSE;
	}

    jType        = json_object_get(json, JSON_TYPE);
    jSeq         = json_object_get(json, JSON_UUID);
    jMessage     = json_object_get(json, JSON_MESSAGE);
    if (jType == NULL || jSeq == NULL || jMessage == NULL) {
        json_decref(json);
        log_error("Error decoding JSON. Missing fields. %s", data);
        return FALSE;
    }

    jLength  = json_object_get(jMessage,     JSON_LENGTH);
    jData    = json_object_get(jMessage,     JSON_DATA);
    jCode    = json_object_get(jMessage,     JSON_CODE);

    if (jLength ==NULL || jData == NULL || jCode == NULL) {
        json_decref(json);
        log_error("Error decoding JSON. Missing fields. %s", data);
        return FALSE;
    }

    *message = (Message *)calloc(1, sizeof(Message));
	if (*message == NULL) {
        log_error("Unable to allocate memory of size: %d", sizeof(Message));
		return FALSE;
	}

    (*message)->reqType  = strdup(json_string_value(jType));
    (*message)->seqNo    = strdup(json_string_value(jSeq));
    (*message)->code     = json_integer_value(jCode);
    (*message)->dataSize = json_integer_value(jLength);
    (*message)->data     = strdup(json_string_value(jData));

    json_decref(json);
    
	return TRUE;
}

static void deserialize_map_array(UMap **map, json_t *json) {

	json_t *jArray;
	json_t *elem, *key, *val;
	int i, size=0;

	*map = (UMap *)calloc(1, sizeof(UMap));
	if (*map==NULL)
		return;

	u_map_init(*map);

	jArray = json_object_get(json, JSON_DATA);

	if (json_is_array(jArray)) {
		size = json_array_size(jArray);

		for (i=0; i<size; i++) {
			elem = json_array_get(jArray, i);

			key = json_object_get(elem, JSON_KEY);
			val = json_object_get(elem, JSON_VALUE);

			u_map_put(*map, json_string_value(key), json_string_value(val));
		}
	}
}

static void deserialize_map(URequest **request, json_t *json) {

	json_t *obj;
	char *str;

	/* Determine the type of map. */
	obj = json_object_get(json, JSON_TYPE);
	if (obj==NULL) {
		return;
	}

	str = json_string_value(obj);

	if (strcasecmp(str, MESH_MAP_TYPE_URL_STR)==0) {
		deserialize_map_array(&(*request)->map_url, json);
	} else if (strcasecmp(str, MESH_MAP_TYPE_HDR_STR)==0) {
		deserialize_map_array(&(*request)->map_header, json);
	} else if (strcasecmp(str, MESH_MAP_TYPE_POST_STR)==0) {
		deserialize_map_array(&(*request)->map_post_body, json);
	} else if (strcasecmp(str, MESH_MAP_TYPE_COOKIE_STR)==0) {
		deserialize_map_array(&(*request)->map_cookie, json);
	}
}
