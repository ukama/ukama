/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>
#include <uuid/uuid.h>

#include "mesh.h"
#include "jserdes.h"

/*
 * add_map_to_request --
 *
 */
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

			json_object_set_new(entry, JSON_KEY, json_string(map->keys[i]));
			json_object_set_new(entry, JSON_VALUE, json_string(map->values[i]));
			json_object_set_new(entry, JSON_LEN,
								json_integer((int)map->lengths[i]));

			json_array_append_new(jArray, entry);
		}
	}
}

/*
 * serialize_device_info --
 *
 */
int serialize_device_info(json_t **json, NodeInfo *device) {

	json_t *jDevice=NULL;

	if (device == NULL) {
		return FALSE;
	}

	if (device->nodeID == NULL) {
		return FALSE;
	}

	*json = json_object();
	if (*json == NULL) {
		return FALSE;
	}

	/* Add Device info. Currently only is the UUID. */
	json_object_set_new(*json, JSON_NODE_INFO, json_object());
	jDevice = json_object_get(*json, JSON_NODE_INFO);

	if (jDevice == NULL) {
		json_decref(*json);
		*json=NULL;
		return FALSE;
	}

	json_object_set_new(jDevice, JSON_NODE_ID, json_string(device->nodeID));

	return TRUE;
}

/*
 * serialize_local_service_response --
 *
 */
int serialize_local_service_response(char **response, Message *message, int len,
                                     char *data) {

    json_t *json, *obj;
    
    /* basic sanity check */
	if (len == 0 || data == NULL || message == NULL)
		return FALSE;

	json = json_object();
	if (json == NULL) {
		return FALSE;
	}

	json_object_set_new(json, JSON_TYPE, json_string(MESH_SERVICE_RESPONSE));
    json_object_set_new(json, JSON_SEQ, json_integer(message->seqNo));

    /* Add node info. */
	json_object_set_new(json, JSON_NODE_INFO, json_object());
	obj = json_object_get(json, JSON_NODE_INFO);
	json_object_set_new(obj, JSON_NODE_ID,
                        json_string(message->nodeInfo->nodeID));
    json_object_set_new(obj, JSON_PORT, json_string(message->nodeInfo->port));

    /* Add service info. */
	json_object_set_new(json, JSON_SERVICE_INFO, json_object());
	obj = json_object_get(json, JSON_SERVICE_INFO);
	json_object_set_new(obj, JSON_NAME,
                        json_string(message->serviceInfo->name));
    
	/* Add response info. */
	json_object_set_new(json, JSON_MESSAGE, json_object());
	obj = json_object_get(json, JSON_MESSAGE);
	json_object_set_new(obj, JSON_LENGTH, json_integer(len));
	json_object_set_new(obj, JSON_DATA, json_string(data));

    *response = json_dumps(json, 0);
    json_decref(json);

	return TRUE;
}

/*
 * serialize_message_data --
 *
 */
static void serialize_message_data(URequest *request, char **data) {

    json_t *json, *jRaw;

    json = json_object();
    if (json == NULL) return;
    
	json_object_set_new(json, JSON_PROTOCOL,
						json_string(request->http_protocol));
	json_object_set_new(json, JSON_METHOD, json_string(request->http_verb));
	json_object_set_new(json, JSON_URL, json_string(request->http_url));
	json_object_set_new(json, JSON_PATH, json_string(request->url_path));

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

    *data = json_dumps(json, 0);
    json_decref(json);
}

/*
 * serialize_websocket_message --
 *
 */
int serialize_websocket_message(char **str, URequest *request, char *nodeID,
                                char *port, char *agent, char *sourcePort) {

    json_t *json=NULL, *jDevice=NULL, *jService=NULL;
	json_t *jRequest=NULL;
    char *data=NULL;

	json = json_object();
	if (json == NULL) {
		return FALSE;
	}

	json_object_set_new(json, JSON_TYPE, json_string(MESH_NODE_REQUEST));
	json_object_set_new(json, JSON_SEQ, json_integer(123456));

	/* Node info */
	json_object_set_new(json, JSON_NODE_INFO, json_object());
	jDevice = json_object_get(json, JSON_NODE_INFO);
	json_object_set_new(jDevice, JSON_NODE_ID, json_string(nodeID));
    json_object_set_new(jDevice, JSON_PORT, json_string(port));

    /* Service info */
	json_object_set_new(json, JSON_SERVICE_INFO, json_object());
	jService = json_object_get(json, JSON_SERVICE_INFO);
	json_object_set_new(jService, JSON_NAME, json_string(agent));
    json_object_set_new(jService, JSON_PORT, json_string(sourcePort));

	/* Serialize and add request info */
    serialize_message_data(request, &data);
	json_object_set_new(json, JSON_MESSAGE, json_object());
	jRequest = json_object_get(json, JSON_MESSAGE);
	json_object_set_new(jRequest, JSON_LENGTH, json_integer(strlen(data)));
	json_object_set_new(jRequest, JSON_DATA, json_string(data));

    *str = json_dumps(json, 0);

	return TRUE;
}

/*
 * deserialize_node_info --
 *
 */
int deserialize_node_info(NodeInfo **node, json_t *json) {

	json_t *id, *port;

	if (json == NULL && node == NULL)
		return FALSE;
    
    id = json_object_get(json, JSON_NODE_ID);
    port = json_object_get(json, JSON_PORT);

    if (id == NULL || port == NULL) {
        return FALSE;
    }
    
	*node = (NodeInfo *)calloc(1, sizeof(NodeInfo));
	if (*node == NULL) {
        log_error("Error allocating memory of size: %d", sizeof(NodeInfo));
		return FALSE;
    }

    (*node)->nodeID = strdup(json_string_value(id));
    (*node)->port   = strdup(json_string_value(port));

	return TRUE;
}

/* 
 * deserialize_service_info --
 *
 */
static int deserialize_service_info(ServiceInfo **service, json_t *json) {

	json_t *name, *port;
  
	if (json == NULL && service == NULL) return FALSE;

    name = json_object_get(json, JSON_NAME);
    port = json_object_get(json, JSON_PORT);

    if (name == NULL || port == NULL) return FALSE;

	*service = (ServiceInfo *)calloc(1, sizeof(ServiceInfo));
	if (*service == NULL) return FALSE;

    (*service)->name = strdup(json_string_value(name));
    (*service)->port = strdup(json_string_value(port));

	return TRUE;
}

/* 
 * deserialize_map_array --
 *
 */
static void deserialize_map_array(UMap **map, json_t *json) {

	json_t *jArray;
	json_t *elem, *key, *val, *len;
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
			len = json_object_get(elem, JSON_LEN);

			u_map_put(*map, json_string_value(key), json_string_value(val));
		}
	}
}

/*
 * deserialize_map --
 *
 */
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

/*
 * deserialize_request_info --
 *
 */
int deserialize_request_info(URequest **request, char *str) {

	json_t *json, *obj, *jRaw;
	int size, i;

	if (str == NULL) return FALSE;

    json = json_loads(str, JSON_DECODE_ANY, NULL);
    if (json == NULL) return FALSE;

	*request = (URequest *)calloc(1, sizeof(URequest));
	if (*request == NULL) {
        log_error("Error allocating memory of size: %d", sizeof(URequest));
		return FALSE;
    }

	obj = json_object_get(json, JSON_PROTOCOL);
	if (obj) {
		(*request)->http_protocol = strdup(json_string_value(obj));
	} else {
		(*request)->http_protocol = strdup("");
	}

	obj = json_object_get(json, JSON_METHOD);
	if (obj) {
		(*request)->http_verb = strdup(json_string_value(obj));
	} else {
		(*request)->http_verb = strdup("");
	}

	obj = json_object_get(json, JSON_URL);
	if (obj) {
		(*request)->http_url = strdup(json_string_value(obj));
	} else {
		(*request)->http_url = strdup("");
	}

	obj = json_object_get(json, JSON_PATH);
	if (obj) {
		(*request)->url_path = strdup(json_string_value(obj));
	} else {
		(*request)->url_path = strdup("");
	}

	/* de-ser the various map. URL, Header, POST, Cookie */
	for (i=0; i < 4; i++) {
		obj = json_object_get(json, JSON_MAP);
		if (obj) {
			deserialize_map(request, obj);
		}
	}

	/* Lastly, de-serialize raw binary data. */
	jRaw = json_object_get(json, JSON_RAW_DATA);
	if (jRaw) {

		obj = json_object_get(jRaw, JSON_LENGTH);
		size = json_integer_value(obj);
		(*request)->binary_body_length = size;

		/* Get the actual data now */
		obj = json_object_get(jRaw, JSON_DATA);

		if (obj) {

			(*request)->binary_body = (void *)calloc(1, size);
			if ((*request)->binary_body == NULL)
				return FALSE;

			memcpy((*request)->binary_body, (void *)json_string_value(obj),
				   size);
		}
	}

	return TRUE;
}

/*
 * deserialize_websocket_message --
 *
 */
int deserialize_websocket_message(Message **message, json_t *json) {

    json_t *jType, *jSeq, *jNodeInfo, *jServiceInfo, *jMessage;
    json_t *jLength, *jData;
	char *jStr=NULL;

	/* Sanity check */
	if (json == NULL) {
		return FALSE;
	}

    jType        = json_object_get(json, JSON_TYPE);
    jSeq         = json_object_get(json, JSON_SEQ);
    jNodeInfo    = json_object_get(json, JSON_NODE_INFO);
    jServiceInfo = json_object_get(json, JSON_SERVICE_INFO);
    jMessage     = json_object_get(json, JSON_MESSAGE);

    if (jType == NULL || jSeq == NULL || jNodeInfo == NULL ||
        jServiceInfo == NULL || jMessage == NULL) {
        jStr = json_dumps(json, 0);
        log_error("Error decoding JSON: %s", jStr);
        free(jStr);
        return FALSE;
    }

    jLength = json_object_get(jMessage, JSON_LENGTH);
    jData    = json_object_get(jMessage, JSON_DATA);

    if (jLength == NULL || jData == NULL) {
        jStr = json_dumps(json, 0);
        log_error("Error decoding JSON: %s", jStr);
        free(jStr);
        return FALSE;
    }

    *message = (Message *)calloc(1, sizeof(Message));
	if (*message == NULL) {
        log_error("Unable to allocate memory of size: %d", sizeof(Message));
		return FALSE;
	}

    (*message)->reqType  = strdup(json_string_value(jType));
    (*message)->seqNo    = json_integer_value(jSeq);
    (*message)->dataSize = json_integer_value(jLength);
    (*message)->data     = strdup(json_string_value(jData));
    
    deserialize_node_info(&(*message)->nodeInfo, jNodeInfo);
	deserialize_service_info(&(*message)->serviceInfo, jServiceInfo);

    /* deserialize the data */
    if (strcmp((*message)->reqType, MESH_SERVICE_REQUEST) == 0) {
        deserialize_request_info((URequest **)&(*message)->data, jData);
    }

	return TRUE;
}
