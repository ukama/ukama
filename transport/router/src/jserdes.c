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
#include "pattern.h"
#include "log.h"

static void log_json(json_t *json);
static int add_key_value_to_pattern(Patterns **patterns, const char *key,
				    json_t *jValue);

/*
 * deserialize_post_route_request --
 *
 * {
 *      "name": "service_name",
 *	"patterns": [{
 *			"key1": "value1",
 *			"key1": "value2",
 *			"path": "/abc"
 *		},
 *		{
 *			"key1": "value1",
 *			"path": "/abv/xcv"
 *		}
 *	],
 *	"forward": {
 *		"ip": "10.0.0.1",
 *		"port": 8080,
 *		"default_path": "/abc"
 *	}
 * }
 *
 */
int deserialize_post_route_request(Service **service, json_t *json) {

  json_t *jName=NULL, *jPatterns=NULL;
  json_t *jPattern=NULL, *jForward=NULL;
  json_t *value, *jIP, *jPort;
  void *iter;
  const char *key;
  Patterns *ptr=NULL;
  int i, count=0;

  if (json == NULL) return FALSE;

  jName     = json_object_get(json, JSON_NAME);
  jPatterns = json_object_get(json, JSON_PATTERNS);
  jForward  = json_object_get(json, JSON_FORWARD);

  if (jName == NULL || jPatterns == NULL || jForward == NULL) {
    log_error("Missing mandatory %s or %s or %s from recvd json request",
	      JSON_NAME, JSON_PATTERNS, JSON_FORWARD);
    log_json(json);
    return FALSE;
  }

  if (!json_is_array(jPatterns)) {
    log_error("Expecting %s Array but missing", JSON_PATTERNS);
    log_json(json);
    return FALSE;
  }

  /* Non-empty */
  count = json_array_size(jPatterns);
  if (count == 0) {
    log_error("%s array with no element.", JSON_PATTERNS);
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
  (*service)->name = strdup(json_string_value(jName));

  (*service)->forward = (Forward *)calloc(1, sizeof(Forward));
  if ((*service)->forward == NULL) {
    log_error("Error allocating memory of size: %lu", sizeof(Forward));
    goto failure;
  }

  (*service)->patterns = (Patterns *)calloc(1, sizeof(Patterns));
  if ((*service)->patterns == NULL) {
    log_error("Error allocating memory of size: %lu", sizeof(Patterns));
    goto failure;
  }

  ptr = (*service)->patterns;

  /* Patterns is an array, iterate over each element */
  for (i=0; i<count; i++) {
    jPattern = json_array_get(jPatterns, i);

    if (jPattern == NULL) {
      goto failure;
    }

    /* Iterate to get all key-value pairs for pattern json object */
    iter = json_object_iter(jPattern);
    while (iter) {

      key   = json_object_iter_key(iter);
      value = json_object_iter_value(iter);

      add_key_value_to_pattern(&ptr, key, value);

      /* iterate to next one */
      iter = json_object_iter_next(jPattern, iter);
    }

    /* if path wasn't specified, log it as info and go with default */
    if (ptr->path == NULL) {
      log_info("Path isn't defined for the pattern. Going default: %s",
	       DEFAULT_PATTERN_PATH);
      ptr->path = strdup(DEFAULT_PATTERN_PATH);
    }

    if (i+1 != count) {
      ptr->next =  (Patterns *)calloc(1, sizeof(Patterns));
      if (ptr->next == NULL) {
	log_error("Error allocating memory of size: %lu", sizeof(Patterns));
	goto failure;
      }
    }

    ptr = ptr->next;
  }

  (*service)->forward->ip   = strdup(json_string_value(jIP));
  (*service)->forward->port = json_integer_value(jPort);

  return TRUE;

 failure:
  if (*service) {
    free_service(*service);
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
 * serialize_get_routes_request --
 *
 * [{
 *	"name": "service_name",
 *	"patterns": [{
 *			"key1": "value1",
 *			"key2": "value2",
 *			"path": "/abc"
 *		},
 *		{
 *			"key1": "value1",
 *			"path": "/abv/xcv"
 *		}
 *	],
 *	"forward": {
 *		"ip": "10.0.0.1",
 *		"port": 8080,
 *		"default_path": "/abc"
 *	}
 *    }
 *  ]
 *
 */
int serialize_get_routes_request(json_t **json, Router *router) {

  json_t *jArray=NULL, *jService=NULL, *jPArray=NULL;
  json_t *jPattern=NULL, *jForward=NULL;
  Patterns *patterns=NULL;
  Pattern *pattern=NULL;
  Service *service=NULL;

  if (router == NULL) return FALSE;

  *json = json_object();
  if (*json == NULL) {
    return FALSE;
  }

  json_object_set_new(*json, JSON_ROUTES, json_array());
  jArray = json_object_get(*json, JSON_ROUTES);

  service = router->services;

  while (service) {

    jService = json_object();
    json_object_set_new(jService, JSON_NAME, json_string(service->name));

    json_object_set_new(jService, JSON_PATTERNS, json_array());
    jPArray = json_object_get(jService, JSON_PATTERNS);

    patterns = service->patterns;

    while (patterns) {

      pattern = patterns->pattern;
      jPattern = json_object();

      while (pattern) {
	json_object_set_new(jPattern, pattern->key,
			    json_string(pattern->value));
	pattern = pattern->next;
      }

      json_object_set_new(jPattern, JSON_PATH,
			  json_string(patterns->path));

      json_array_append(jPArray, jPattern);
      patterns = patterns->next;
      json_decref(jPattern);
    }

    json_object_set_new(jService, JSON_FORWARD, json_object());
    jForward = json_object_get(jService, JSON_FORWARD);

    if (service->forward) {
      json_object_set_new(jForward, JSON_IP,
			  json_string(service->forward->ip));
      json_object_set_new(jForward, JSON_PORT,
			  json_integer(service->forward->port));
    }

    json_array_append(jArray, jService);
    service = service->next;

    json_decref(jService);
  }

  return TRUE;
}

/*
 * add_key_value_to_pattern --
 *
 */
static int add_key_value_to_pattern(Patterns **patterns, const char *key,
				    json_t *jValue) {

  Pattern **ptr=NULL;
  Pattern *tmp=NULL;

  if (*patterns == NULL || jValue == NULL) return FALSE;

  if (strcmp(key, JSON_PATH) == 0){
    (*patterns)->path = strdup(json_string_value(jValue));
    return TRUE;
  }

  if (!(*patterns)->pattern) {
    ptr = &((*patterns)->pattern);
  } else {
    for (tmp=(*patterns)->pattern; tmp->next; tmp=tmp->next);
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

