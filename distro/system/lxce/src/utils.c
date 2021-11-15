/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#define _GNU_SOURCE
#include <sched.h>
#include <stdlib.h>
#include <jansson.h>
#include <string.h>

#include "log.h"
#include "utils.h"

#define TRUE  1
#define FALSE 0
/*
 * set_integer_object_value --
 *
 */
int set_integer_object_value(json_t *json, int *param, char *objName,
			     int mandatory, int defValue) {

  json_t *obj;

  obj = json_object_get(json, objName);
  if (obj==NULL) {
    if (mandatory) {
      log_error("Missing Mandatory JSON field: %s Setting to default: %d",
                objName, defValue);
      if (defValue)  {
        *param = defValue;
      } else {
        return FALSE;
      }
    } else {
      log_debug("Missing JSON field: %s. Ignored.", objName);
      *param = 0;
    }
  } else {
    *param = json_integer_value(obj);
  }

  return TRUE;
}

/*
 * set_str_object_value --
 *
 */
int set_str_object_value(json_t *json, char **param, char *objName,
			 int mandatory, char *defValue) {

  json_t *obj;

  obj = json_object_get(json, objName);
  if (obj==NULL) {
    if (mandatory) {
      log_error("Missing Mandatory JSON field: %s Setting to default: %s",
                objName, defValue);
      if (defValue)  {
        *param = strdup(defValue);
      } else {
        return FALSE;
      }
    } else {
      log_debug("Missing JSON field: %s. Ignored.", objName);
      *param = NULL;
    }
  } else {
    *param = strdup(json_string_value(obj));
  }

  return TRUE;
}

/*
 * namespace_flag --
 *
 */
int namespaces_flag(char *ns) {

  if (strcmp(ns, "pid")==0) {
    return CLONE_NEWPID;
  } else if (strcmp(ns, "uts")==0) {
    return CLONE_NEWUTS;
  } else if (strcmp(ns, "network")==0) {
    return CLONE_NEWNET;
  } else if (strcmp(ns, "mount")==0) {
    return CLONE_NEWNS;
  } else if (strcmp(ns, "user")==0) {
    return CLONE_NEWUSER;
  } else {
    log_error("Unsupported namespace type detecetd: %s", ns);
    return 0;
  }

  return 0;
}
