/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string.h>
#include <jansson.h>

#include "wimc.h"
#include "jserdes.h"

/*
 * deserialize_wimc_response --
 *
 */

int deserialize_wimc_response(char *response, json_t *json) {

  int ret;
  json_t *root, *type, *str;
  char *respType;

  /* sanity check */
  if (!json) return FALSE;

  root = json_object_get(json, JSON_WIMC_RESPONSE);
  if (root) {
    type = json_object_get(root, JSON_TYPE);
    str  = json_object_get(root, JSON_VOID_STR);

    respType = strdup(json_string_value(type));

    if (strcmp(respType, WIMC_RESP_TYPE_RESULT)==0) {
      response = strdup(json_string_value(str));
      ret = TRUE;
    } else if (strcmp(respType, WIMC_RESP_TYPE_PROCESSING)==0) {
      response = NULL;
      ret = TRUE+1;
    } else if (strcmp(respType, WIMC_RESP_TYPE_ERROR)==0) {
      response = strdup(json_string_value(str));
      ret = FALSE;
    } else {
      ret = FALSE-1;
    }
  }

  free(respType);
  return ret;
}

