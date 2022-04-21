/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>

#include "jserdes.h"
#include "log.h"

static void log_json(json_t *json);
static void free_service(Service **service);
static int add_key_value_to_pattern(Service **service, const char *key,
				    json_t *jValue);

/*
 * deserialize_post_route_request --
 *
 *  {
 *      "pattern": {
 *		"key1": "value1",
 *		"key2": "value2"
 *	},
 *
 *	"forward": {
 *		"ip": "192.168.0.1",
 *		"port": "8080"
 *	}
 *  }
 */
int deserialize_post_route_request(Service **service, json_t *json) {

  json_t *jPattern=NULL, *jForward=NULL;
  json_t *value, *jIP, *jPort;
  void *iter;
  const char *key;

  if (json == NULL) return FALSE;

  jPattern = json_object_get(json, JSON_PATTERN);
  jForward = json_object_get(json, JSON_FORWARD);

  if (jPattern == NULL || jForward == NULL) {
    log_error("Missing mandatory %s or %s from recvd json request",
	      JSON_PATTERN, JSON_FORWARD);
    log_json(json);
    return FALSE;
  }

  /* forward: ip and port */
  jIP   = json_object_get(jForward, JSON_IP);
  jPort = json_object_get(jForward, JSON_PORT);

  if (jIP == NULL || jPort == NULL) {
    log_error("Missing %s or %s from recvd json request", JSON_IP, JSON_PORT);
    log_json(json);
    return FALSE;
  }

  *service = (Service *)calloc(1, sizeof(Service));
  if (*service == NULL) {
    log_error("Error allocating memory of size: %lu", sizeof(Service));
    return FALSE;
  }

  uuid_clear((*service)->uuid);
  (*service)->forward = (Forward *)calloc(1, sizeof(Forward));
  if ((*service)->forward == NULL) {
    log_error("Error allocating memory of size: %ls or %lu",
	      sizeof(Pattern), sizeof(Forward));
    goto failure;
  }

  /* Iterate to get all key-value pairs for pattern json object */
  iter = json_object_iter(jPattern);
  while (iter) {

    key   = json_object_iter_key(iter);
    value = json_object_iter_value(iter);

    add_key_value_to_pattern(service, key, value);

    /* iterate to next one */
    iter = json_object_iter_next(jPattern, iter);
  }

  (*service)->forward->ip   = strdup(json_string_value(jIP));
  (*service)->forward->port = strdup(json_string_value(jPort));

  return TRUE;

 failure:
  if (*service) {
    free_service(service);
    free(*service);
    *service=NULL;
  }
  return FALSE;
}

/*
 * deserialize_delete_router_request
 *
 */
int deserialize_delete_route_request(char **uuidStr, json_t *json) {

  json_t *jID;

  if (json == NULL) return FALSE;

  jID = json_object_get(json, JSON_UUID);
  if (jID == NULL) {
    log_error("Unable to find %s as key for DELETE request", JSON_UUID);
    return FALSE;
  }

  *uuidStr = strdup(json_string_value(jID));

  return TRUE;
}

/*
 * serialize_post_route_response --
 *
 */
int serialize_post_route_response(json_t **json, int respCode, uuid_t uuid,
				  char *errStr) {

  char idStr[36+1];

  *json = json_object();
  if (*json == NULL) {
    return FALSE;
  }

  switch (respCode) {
  case UUID:
    if (uuid_is_null(uuid)) return FALSE;

    uuid_unparse(uuid, &idStr[0]);
    json_object_set_new(*json, JSON_UUID, json_string(idStr));
    break;

  case ERROR:
    if (errStr==NULL) return FALSE;

    json_object_set_new(*json, JSON_ERROR, json_string(errStr));
    break;

  default:
    break;
  }

  return TRUE;
}

/*
 * add_key_value_to_pattern --
 *
 */
static int add_key_value_to_pattern(Service **service, const char *key,
				    json_t *jValue) {

  Pattern **ptr=NULL;
  Pattern *tmp=NULL;

  if (*service == NULL || jValue == NULL) return FALSE;

  if (!(*service)->pattern) {
    ptr = &((*service)->pattern);
  } else {
    for (tmp=(*service)->pattern; tmp->next; tmp=tmp->next);
    ptr = &(tmp->next);
  }

  *ptr = (Pattern *)calloc(1, sizeof(Pattern));
  if (*ptr == NULL) {
    log_error("Error allocating memory of size: %lu", sizeof(Pattern));
    return FALSE;
  }

  (*ptr)->key   = strdup(key);
  (*ptr)->value = strdup(json_string_value(jValue));
  (*ptr)->next  = NULL;

  return TRUE;
}

/*
 * log_json --
 *
 */
static void log_json(json_t *json) {

  char *str = NULL;

  str = json_dumps(json, 0);
  if (str) {
    log_debug("json str: %s", str);
    free(str);
  }
}

/*
 * free_service --
 *
 */

static void free_service(Service **service) {

  Pattern *pattern, *tmp;
  Forward *forward;

  if (*service == NULL) return;

  pattern = (*service)->pattern;
  forward = (*service)->forward;

  while(pattern) {
    if (pattern->key)   free(pattern->key);
    if (pattern->value) free(pattern->value);

    tmp = pattern->next;
    free(pattern);
    pattern = tmp;
  }

  if (forward) {
    if (forward->ip)   free(forward->ip);
    if (forward->port) free(forward->port);
  }
}
