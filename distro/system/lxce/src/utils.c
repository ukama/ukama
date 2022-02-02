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
#include <sys/capability.h>

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

/*
 * get_json_array_elems --
 *
 */
void get_json_arrary_elems(json_t *json, int *argc, char ***argv,
			   char *objName) {

  json_t *jArray, *jElem;
  int i;

  jArray = json_object_get(json, objName);

  if (jArray != NULL) {
    *argc = json_array_size(jArray);

    if (*argc == 0) {
      *argv = NULL;
      return;
    }

    *argv = (char **)calloc(*argc, sizeof(char *));
    if (argv==NULL) return;

    for (i=0; i<(*argc); i++) {
      jElem = json_array_get(jArray, i);
      if (jElem) {
	(*argv)[i] = strdup(json_string_value(jElem));
	log_debug("argv: %d %s", i, (*argv)[i]);
      }
    }
  } else {
    log_error("Array not found: %s", objName);
  }
}

/*
 * str_to_cap --
 *
 */
int str_to_cap(const char *str) {

  if (strcmp(str, "CAP_BLOCK_SUSPEND")==0) {
    return CAP_BLOCK_SUSPEND;
  } else if (strcmp(str, "CAP_IPC_LOCK")==0) {
    return CAP_IPC_LOCK;
  } else if (strcmp(str, "CAP_MAC_ADMIN")==0) {
    return CAP_MAC_ADMIN;
  } else if (strcmp(str, "CAP_MAC_OVERRIDE")==0) {
    return CAP_MAC_OVERRIDE;
  }

  log_error("Invalid capabilities: %s", str);
  return 0;
}
