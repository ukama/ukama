/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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

