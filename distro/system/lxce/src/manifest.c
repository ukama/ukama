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
#include "lxce_config.h"
#include "cspace.h"

/*
 * is_valid_cspace --
 *
 */
static int is_valid_cspace(char *name, CSpace *space) {

  CSpace *ptr=space;

  if (!ptr) return FALSE;

  while (ptr) {
    if (strcmp(ptr->name, name)==0) {
      return TRUE;
    }
    ptr = ptr->next;
  }

  return FALSE;
}

/*
 * deserialize_cApp --
 *
 */
static int deserialize_cApp(ArrayElem *elem, json_t *json, CSpace *spaces) {

  json_t *name, *tag, *restart, *contained;
  char *tmp;

  name      = json_object_get(json, JSON_NAME);
  tag       = json_object_get(json, JSON_TAG);
  restart   = json_object_get(json, JSON_RESTART);
  contained = json_object_get(json, JSON_CONTAINED);

  if (name==NULL || tag==NULL || contained==NULL) {
    log_error("Missing cAPP cfg parameter in the Manifest file");
    return FALSE;
  }

  tmp = json_string_value(contained);

  if (is_valid_cspace(tmp, spaces)==FALSE) {
    log_error("Invalid cSpace \"%s\" in the config. Ignoring", tmp);
    return FALSE;
  }

  elem->name      = strdup(json_string_value(name));
  elem->tag       = strdup(json_string_value(tag));
  elem->contained = strdup(json_string_value(contained));

  if (restart) {
    elem->restart = json_integer_value(restart);
  } else {
    elem->restart = FALSE;
  }

  elem->next = NULL;

  return TRUE;
}

/*
 * deserialize_manifest_file -- convert the json into internal struct
 *
 */
static int deserialize_manifest_file(Manifest *manifest, CSpace *spaces,
				     json_t *json) {

  int j=0, size=0;
  json_t *obj;
  json_t *jArray, *jElem;
  ArrayElem *elem=NULL;

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

  /* deserialize ukama contained Apps */
  jArray = json_object_get(json, JSON_CAPP);
  if (jArray != NULL) {

    size = json_array_size(jArray);

    manifest->arrayElem = (ArrayElem *)calloc(size, sizeof(ArrayElem));
    if (manifest->arrayElem==NULL) {
      log_error("Memory allocation error. Size: %d", size*sizeof(ArrayElem));
      return FALSE;
    }

    elem = manifest->arrayElem;

    for (j=0; j<size; j++) {

      jElem = json_array_get(jArray, j);
      if (jElem) {

	if (elem==NULL) {
	  elem = (ArrayElem *)calloc(1, sizeof(ArrayElem));
	  if (elem == NULL) {
	    log_error("Memory allocation error. Size: %d", sizeof(ArrayElem));
	    return FALSE;
	  }
	}

	if (deserialize_cApp(elem, jElem, spaces)) {
	  elem = elem->next;
	}
      }
    } /* for loop */
  } else {
    log_error("Error parsing %s", JSON_CAPP);
    return FALSE;
  }

  return TRUE;
}

/*
 * process_manifest -- parse the manifest file.
 *
 */
int process_manifest(char *fileName, Manifest *manifest, void *arg) {

  int ret=FALSE;
  FILE *fp;
  char *buffer=NULL;
  long size=0;
  json_t *json;
  json_error_t jerror;
  CSpace *spaces = (CSpace *)arg;

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
  memset(buffer, 0, size+1);
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
  ret = deserialize_manifest_file(manifest, spaces, json);

 done:
  if (buffer) free(buffer);
  fclose(fp);

  json_decref(json);
  return ret;
}

/*
 * clear_cApp_cfg --
 *
 */
static void clear_cApp_cfg(ArrayElem *elem) {

  ArrayElem *ptr, *prev;

  ptr = elem;

  while (ptr) {

    free(ptr->name);
    free(ptr->tag);
    free(ptr->contained);

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

  clear_cApp_cfg(manifest->arrayElem);

  return;
}
