/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Functions related to manifest.conf
 */

#include <stdio.h>
#include <errno.h>
#include <string.h>

#include "manifest.h"
#include "log.h"
#include "config.h"

/*
 * deserialize_container_cfg --
 *
 */
static int deserialize_container_cfg(Container *con, json_t *json) {

  json_t *order, *name, *ver, *active, *restart;

  order   = json_object_get(json, JSON_ORDER);
  name    = json_object_get(json, JSON_NAME);
  ver     = json_object_get(json, JSON_VERSION);
  active  = json_object_get(json, JSON_ACTIVE);
  restart = json_object_get(json, JSON_RESTART);

  if (name==NULL || ver==NULL) {
    log_error("Missing container cfg parameter in the Manifest file");
    return FALSE;
  }

  if (order) {
    con->order   = json_integer_value(order);
    if (con->order < 1) {
      con->order = 1; /* Order starts with 1 */
    }
  } else {
    con->order = 0; /* Execute before anything with order */
  }

  con->name    = strdup(json_string_value(name));
  con->version = strdup(json_string_value(ver));

  if (active) {
    con->active  = json_integer_value(active);
  } else {
    con->active = TRUE;
  }

  if (restart) {
    con->restart = json_integer_value(restart);
  } else {
    con->restart = FALSE;
  }

  con->path = NULL; /* will be set after querying WIMC */
  con->next = NULL;

  return TRUE;
}

/*
 * deserialize_manifest_file -- convert the json into internal struct.
 *
 */
static int deserialize_manifest_file(Manifest *manifest, json_t *json) {

  int i=0, j=0, size=0;
  json_t *obj;
  json_t *jArray, *elem;
  Container **cPtr=NULL;
  char *type;

  if (manifest == NULL) return FALSE;
  if (json == NULL) return FALSE;

  obj = json_object_get(json, JSON_VERSION);
  if (obj==NULL) {
    log_error("Missing mandatory JSON field: %s", JSON_VERSION);
    return FALSE;
  } else {
    manifest->version = strdup(json_string_value(obj));
  }

  obj = json_object_get(json, JSON_TARGET);
  if (obj==NULL) {
    log_error("Missing mandatory JSON field: %s", JSON_TARGET);
    return FALSE;
  } else {
    manifest->target = strdup(json_string_value(obj));
  }

  if (strcmp(manifest->target, MANIFEST_SERIAL)==0) {
    obj = json_object_get(json, JSON_SERIAL);
    if (obj==NULL) {
      log_error("Missing mandatory JSON field: %s", JSON_SERIAL);
      return FALSE;
    } else {
      manifest->serial = strdup(json_string_value(obj));
    }
  } else {
    manifest->serial = NULL;
  }

  /* deserialize containers cfg */
  do {

    if (i==0) { /* boot containers */
      jArray = json_object_get(json, JSON_BOOT);
      cPtr = &manifest->boot;
      type = JSON_BOOT;
    } else if (i==1) { /* service containers */
      jArray = json_object_get(json, JSON_SERVICE);
      cPtr = &manifest->service;
      type = JSON_SERVICE;
    } else if (i==2) { /* shutdown containers */
      jArray = json_object_get(json, JSON_SHUTDOWN);
      cPtr = &manifest->shutdown;
      type = JSON_SHUTDOWN;
    }

    if (jArray != NULL) {

      size = json_array_size(jArray);

      for (j=0; j<size; j++) {
	elem = json_array_get(jArray, j);
	if (elem) {
	  *cPtr = (Container *)calloc(1, sizeof(Container));
	  if (*cPtr == NULL) {
	    log_error("Memory allocation error. Size: %d", sizeof(Container));
	    return FALSE;
	  }

	  if (deserialize_container_cfg(*cPtr, elem)) {
	    cPtr = &((*cPtr)->next);
	  } else {
	    log_error("Error parsing Container cfg for %s", type);
	    return FALSE;
	  }
	}
      }
    } else {
      log_error("Error parsing Container cfg for %s", type);
      return FALSE;
    }

    i++;
  } while (i<3);

  return TRUE;
}

/*
 * process_manifest -- parse the manifest file.
 *
 */
int process_manifest(char *fileName, Manifest *manifest) {

  int ret=FALSE;
  FILE *fp;
  char *buffer=NULL;
  long size=0;
  json_t *json;
  json_error_t jerror;

  /* Sanity check */
  if (fileName==NULL) return FALSE;
  if (manifest==NULL) return FALSE;

  if ((fp = fopen(fileName, "rb")) == NULL) {
    log_error("Error opening manifest file: %s Error %s", fileName,
	      strerror(errno));
    return FALSE;
  }

  /* Read everything into buffer */
  fseek(fp, 0, SEEK_END);
  size = ftell(fp);
  fseek(fp, 0, SEEK_SET);

  if (size > MANIFEST_MAX_SIZE) {
    log_error("Error opening manifest file: %s Error: File size too big: %ld",
	      fileName, size);
    fclose(fp);
    return FALSE;
  }

  buffer = (char *)malloc(size+1);
  if (buffer==NULL) {
    log_error("Error allocating memory of size: %ld", size+1);
    fclose(fp);
    return FALSE;
  }

  fread(buffer, 1, size, fp); /* Read everything into buffer */

  /* Trying loading it as JSON */
  json = json_loads(buffer, 0, &jerror);
  if (json==NULL) {
    log_error("Error loading manifest into JSON format. File: %s Size: %ld",
	      fileName, size);
    log_error("JSON error on line: %d: %s", jerror.line, jerror.text);
    goto done;
  }

  /* Now convert JSON into internal struct */
  ret = deserialize_manifest_file(manifest, json);

 done:
  if (buffer) free(buffer);
  fclose(fp);

  json_decref(json);
  return ret;
}

/*
 * get_container_local_path --
 *
 */

void get_containers_local_path(Manifest *manifest, Config *config) {

  Container *ptr=NULL;
  int i;

  /* iterate over boot, service and shutdown container name:tag
   * and get each one's path from wimc.d
   */

  for (i=0; i<3; i++) {

    if (i==0 && manifest->boot) {
	ptr = manifest->boot;
    } else if (i==1 && manifest->service) {
      ptr = manifest->service;
    } else if (i==3 && manifest->shutdown) {
	ptr = manifest->shutdown;
    }

    while (ptr) {
      if (ptr->name && ptr->version) {
	get_container_path_from_wimc(ptr->name, ptr->version,
				     config->wimcHost, config->wimcPort,
				     ptr->path);
      }
      ptr = ptr->next;
    }
  }
}

/*
 * clear_con_cfg --
 *
 */
static void clear_con_cfg(Container *container) {

  Container *ptr, *prev;

  ptr = container;

  while (ptr) {

    free(ptr->name);
    free(ptr->version);
    if (ptr->path) free(ptr->path);

    prev = ptr;
    ptr = ptr->next;

    free(prev);
  }

  return;
}

/*
 * clear_manifest --
 *
 */
void clear_manifest(Manifest *manifest) {

  if (manifest==NULL) return;

  free(manifest->version);
  free(manifest->serial);
  free(manifest->target);

  clear_con_cfg(manifest->boot);
  clear_con_cfg(manifest->service);
  clear_con_cfg(manifest->shutdown);

  return;
}
