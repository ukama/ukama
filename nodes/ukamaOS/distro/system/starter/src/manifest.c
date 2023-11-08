/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>

#include "manifest.h"
#include "config.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"


void json_log(json_t *json) {

    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        log_debug("json str: %s", str);
        free(str);
    }
}

static bool get_json_entry(json_t *json, char *key, json_type type,
                           char **strValue,
                           int *intValue,
                           double *doubleValue,
                           json_t **jsonObj) {

    json_t *jEntry=NULL;

    if (json == NULL || key == NULL) return USYS_FALSE;

    jEntry = json_object_get(json, key);
    if (jEntry == NULL) {
        log_error("Missing %s key in json", key);
        return USYS_FALSE;
    }

    switch(type) {
    case (JSON_STRING):
        *strValue = strdup(json_string_value(jEntry));
        break;
    case (JSON_INTEGER):
        *intValue = json_integer_value(jEntry);
        break;
    case (JSON_REAL):
        *doubleValue = json_real_value(jEntry);
        break;
    case (JSON_OBJECT):
        *jsonObj = json_object_get(json, key);
        break;
    default:
        log_error("Invalid type for json key-value: %d", type);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static void free_capps(CappsManifest *ptr) {

    usys_free(ptr->name);
    usys_free(ptr->tag);
    usys_free(ptr->space);
}

void free_manifest(Manifest *ptr) {

    CappsManifest *cPtr, *oldPtr;

    if (ptr == NULL) return;

    cPtr = ptr->cappsManifest;
    usys_free(ptr->version);
    usys_free(ptr->target);

    while(cPtr) {

        oldPtr = cPtr;
        free_capps(cPtr);
        cPtr = cPtr->next;
        usys_free(oldPtr);
    }

    usys_free(ptr);
    return;
}

static bool deserialize_spaces(Manifest **manifest, json_t *json) {

    int ret=USYS_FALSE;
    SpacesManifest *spaces, *ptr;
    
    spaces = (SpacesManifest *)calloc(1, sizeof(SpacesManifest));
    if (spaces == NULL) {
        usys_log_error("Error allocating memory of size: %d",
                       sizeof(SpacesManifest));
        return USYS_FALSE;
    }

    ret |= get_json_entry(json, JTAG_NAME, JSON_STRING,
                          &spaces->name, NULL, NULL, NULL);
    if (ret == USYS_FALSE) {
          usys_log_error("Error deserializing the spaces json.");
          return USYS_FALSE;
    }

    if ((*manifest)->spacesManifest == NULL ){
        (*manifest)->spacesManifest = spaces;
    } else {
        for (ptr=(*manifest)->spacesManifest; ptr->next; ptr=ptr->next);
        ptr->next = spaces;
    }

    return USYS_TRUE;
}

static int deserialize_capps(Manifest **manifest, json_t *json) {

    int ret=USYS_FALSE;
    CappsManifest *capp, *ptr;
    
    capp = (CappsManifest *)calloc(1, sizeof(CappsManifest));
    if (capp == NULL) {
        usys_log_error("Error allocating memory of size: %d",
                       sizeof(CappsManifest));
        return USYS_FALSE;
    }

    ret |= get_json_entry(json, JTAG_NAME, JSON_STRING,
                          &capp->name, NULL, NULL, NULL);
    ret |= get_json_entry(json, JTAG_TAG, JSON_STRING,
                          &capp->tag, NULL, NULL, NULL);
    ret |= get_json_entry(json, JTAG_RESTART, JSON_INTEGER,
                          NULL, &capp->restart, NULL, NULL);
    ret |= get_json_entry(json, JTAG_SPACE, JSON_STRING,
                          &capp->space, NULL, NULL, NULL);

    if (ret == USYS_FALSE) {
          usys_log_error("Error deserializing the capp json.");
          free_capps(capp);

          return USYS_FALSE;
    }

    if ((*manifest)->cappsManifest == NULL ){
        (*manifest)->cappsManifest = capp;
    } else {
        for (ptr=(*manifest)->cappsManifest; ptr->next; ptr=ptr->next);
        ptr->next = capp;
    }

    return USYS_TRUE;
}

static bool deserialize_manifest_file(Manifest **manifest,
                                      json_t *json) {

    int ret=USYS_TRUE, count=0, i;
    json_t *jCappsArray=NULL;
    json_t *jEntry=NULL;
    json_t *jSpacesArray=NULL;
    
    *manifest = (Manifest *)calloc(1, sizeof(Manifest));
    if (*manifest == NULL) {
        usys_log_error("Error allocating memory of size: %ld",
                       sizeof(Manifest));
        return USYS_FALSE;
    }

    ret |= get_json_entry(json, JTAG_VERSION, JSON_STRING,
                          &(*manifest)->version, NULL, NULL, NULL);
    ret |= get_json_entry(json, JTAG_TARGET, JSON_STRING,
                          &(*manifest)->target, NULL, NULL, NULL);
    ret |= get_json_entry(json, JTAG_SPACES, JSON_OBJECT,
                          NULL, NULL, NULL, &jSpacesArray);
    ret |= get_json_entry(json, JTAG_CAPPS, JSON_OBJECT,
                          NULL, NULL, NULL, &jCappsArray);

    if (ret == USYS_FALSE) {
        usys_log_error("Error deserializing the manifest json");
        free_manifest(*manifest);
        json_log(json);
        return USYS_FALSE;
    }

    /* spaces */
    count = json_array_size(jSpacesArray);
    if (count == 0) {
        usys_log_error("No space defined!");
        json_log(json);
        return USYS_FALSE;
    } 
    for (i=0; i<count; i++) {
        jEntry = json_array_get(jSpacesArray, i);
        deserialize_spaces(manifest, jEntry);
    }

    /* capps */
    count = json_array_size(jCappsArray);
    if (count == 0) {
        usys_log_error("No capps to run!");
        json_log(json);
        return USYS_FALSE;
    }
    
    for (i=0; i<count; i++) {
        jEntry = json_array_get(jCappsArray, i);
        deserialize_capps(manifest, jEntry);
    }

    return USYS_TRUE;
}

bool read_manifest_file(Manifest **manifest, char *fileName) {

    int ret=USYS_FALSE;
    FILE *fp;
    char *buffer=NULL;
    long size=0;
    json_t *json;
    json_error_t jerror;

    /* Sanity check */
    if (fileName == NULL || manifest == NULL) return USYS_FALSE;

    if ((fp = fopen(fileName, "rb")) == NULL) {
        log_error("Error opening manifest file: %s Error %s", fileName,
                  strerror(errno));
        return USYS_FALSE;
    }

    /* Read everything into buffer */
    fseek(fp, 0, SEEK_END);
    size = ftell(fp);
    fseek(fp, 0, SEEK_SET);

    if (size > MAX_MANIFEST_FILE_SIZE) {
        usys_log_error("Error opening manifest file: %s "
                       "Error: File size too big: %ld",
                       fileName, size);
        fclose(fp);
        return USYS_FALSE;
    }

    buffer = (char *)calloc(1, size+1);
    if (buffer == NULL) {
        log_error("Error allocating memory of size: %ld", size+1);
        fclose(fp);
        return USYS_FALSE;
    }
    memset(buffer, 0, size+1);
    fread(buffer, 1, size, fp); /* Read everything into buffer */

    /* Trying loading it as JSON */
    json = json_loads(buffer, 0, &jerror);
    if (json == NULL) {
        usys_log_error("Error loading manifest into JSON format."
                       "File: %s Size: %ld",
                       fileName, size);
        usys_log_error("JSON error on line: %d: %s", jerror.line, jerror.text);
    } else {
        /* Now convert JSON into internal struct */
        ret = deserialize_manifest_file(manifest, json);
    }

done:
    if (buffer) free(buffer);

    fclose(fp);
    json_decref(json);

    return ret;
}
