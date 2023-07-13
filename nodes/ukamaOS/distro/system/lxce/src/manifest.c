/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <errno.h>
#include <string.h>

#include "manifest.h"
#include "log.h"
#include "lxce_config.h"
#include "cspace.h"

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
                           char **strValue, int *intValue,
                           double *doubleValue, json_t **jsonObj) {

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

static void free_manifest_capps(ArrayElem *ptr) {

    if (ptr == NULL) return;

    usys_free(ptr->name);
    usys_free(ptr->tag);
    usys_free(ptr->rootfs);

    return free_manifest_capps(ptr->next);
}

void free_manifest(Manifest *ptr) {

    if (ptr == NULL) return;

    usys_free(ptr->version);
    usys_free(ptr->target);

    free_manifest_capps(ptr->boot);
    free_manifest_capps(ptr->services);
    free_manifest_capps(ptr->reboot);

    return;
}

static int deserialize_capp_entries(ArrayElem **capp, json_t *json) {

    int ret=USYS_FALSE;
    int entries=0, i=0;
    json_t *jEntry = NULL;
    ArrayElem **ptr;

    if (json == NULL) {
        usys_log_error("No data to deserialize");
        return USYS_FALSE;
    }

    entries = json_array_size(json);
    if (!entries) {
        usys_log_error("No array elements");
        return USYS_FALSE;
    }

    ptr = capp;
    
    for (i=0; i<entries; i++) {

      jEntry = json_array_get(json, i);

      *ptr = (ArrayElem *)calloc(1, sizeof(ArrayElem));
      if (*ptr == NULL) {
        usys_log_error("Error allocating memory of size: %d",
                       sizeof(ArrayElem));
        return USYS_FALSE;
      }

      ret |= get_json_entry(jEntry, JTAG_NAME, JSON_STRING,
                            &(*ptr)->name, NULL, NULL, NULL);
      ret |= get_json_entry(jEntry, JTAG_TAG, JSON_STRING,
                            &(*ptr)->tag, NULL, NULL, NULL);
      ret |= get_json_entry(jEntry, JTAG_RESTART, JSON_INTEGER,
                            NULL, &(*ptr)->restart, NULL, NULL);
      if (ret == USYS_FALSE) {
          usys_log_error("Error deserializing the capp json.");
          json_log(json);
          return USYS_FALSE;
      }

      ptr = &(*ptr)->next;
    }
    
    return USYS_TRUE;
}

static bool deserialize_manifest_file(Manifest **manifest,
                                     CSpace *spaces,
                                     json_t *json) {

    int ret;
    json_t *boot=NULL, *services=NULL, *reboot=NULL;

    if (json == NULL) {
        usys_log_error("No data to deserialize");
        return USYS_FALSE;
    }

    *manifest = (Manifest *)calloc(1, sizeof(Manifest));
    if (*manifest == NULL) {
        usys_log_error("Error allocating memory of size: %d",
                       sizeof(Manifest));
        return USYS_FALSE;
    }

    ret |= get_json_entry(json, JTAG_VERSION, JSON_STRING,
                          &(*manifest)->version, NULL, NULL, NULL);
    ret |= get_json_entry(json, JTAG_TARGET, JSON_STRING,
                          &(*manifest)->target, NULL, NULL, NULL);

    ret |= get_json_entry(json, JTAG_BOOT, JSON_OBJECT,
                          NULL, NULL, NULL, &boot);
    ret |= get_json_entry(json, JTAG_SERVICES, JSON_OBJECT,
                          NULL, NULL, NULL, &services);
    ret |= get_json_entry(json, JTAG_REBOOT, JSON_OBJECT,
                          NULL, NULL, NULL, &reboot);

    if (ret == USYS_FALSE) {
        usys_log_error("Error deserializing the json.");
        json_log(json);
        free_manifest(manifest);
        return USYS_FALSE;
    }

    ret |= deserialize_capp_entries(&(*manifest)->boot,     boot);
    ret |= deserialize_capp_entries(&(*manifest)->services, services);
    ret |= deserialize_capp_entries(&(*manifest)->reboot,   reboot);

    if (ret == USYS_FALSE) {
        usys_log_error("Error deserializing the capp json");
        json_log(json);
        free_manifest(manifest);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

int process_manifest(Manifest **manifest, char *fileName, void *arg) {

    int ret=USYS_FALSE;
    FILE *fp;
    char *buffer=NULL;
    long size=0;
    json_t *json;
    json_error_t jerror;
    CSpace *spaces = (CSpace *)arg;

    /* Sanity check */
    if (fileName == NULL || manifest == NULL) return FALSE;

    if ((fp = fopen(fileName, "rb")) == NULL) {
        log_error("Error opening manifest file: %s Error %s", fileName,
                  strerror(errno));
        return USYS_FALSE;
    }

    /* Read everything into buffer */
    fseek(fp, 0, SEEK_END);
    size = ftell(fp);
    fseek(fp, 0, SEEK_SET);

    if (size > MANIFEST_MAX_SIZE) {
        usys_log_error("Error opening manifest file: %s "
                       "Error: File size too big: %ld",
                       fileName, size);
        fclose(fp);
        return USYS_FALSE;
    }

    buffer = (char *)malloc(size+1);
    if (buffer==NULL) {
        log_error("Error allocating memory of size: %ld", size+1);
        fclose(fp);
        return USYS_FALSE;
    }
    memset(buffer, 0, size+1);
    fread(buffer, 1, size, fp); /* Read everything into buffer */

    /* Trying loading it as JSON */
    json = json_loads(buffer, 0, &jerror);
    if (json==NULL) {
        usys_log_error("Error loading manifest into JSON format."
                       "File: %s Size: %ld",
                       fileName, size);
        usys_log_error("JSON error on line: %d: %s", jerror.line, jerror.text);
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
 * copy_capps_to_cspace_rootfs --
 *                             sPath: /capps/pkgs
 *                             dPath: /capps/rootfs
 */
static void copy_capps_to_cspace_rootfs(char *spaceName,
                                        ArrayElem *ptr,
                                        char *sPath,
                                        char *dPath) {
    int ret=0;
    char src[CSPACE_MAX_BUFFER]   = {0};
    char dest[CSPACE_MAX_BUFFER]  = {0};
    char runMe[CSPACE_MAX_BUFFER] = {0};

    if (ptr == NULL) return;

    for (; ptr; ptr=ptr->next) {
        if (ptr->name == NULL || ptr->tag == NULL ) {
            usys_log_error("Invalid manifest entry. Ignoring");
            continue;
        }

        /* Create dest dir if needed */
        sprintf(runMe, "/bin/mkdir -p %s/%s/%s", dPath, spaceName, sPath);
        log_debug("Running command: %s", runMe);
        if ((ret = system(runMe)) < 0) {
            usys_log_error("Unable to execute cmd %s for space: %s Code: %d",
                           runMe, spaceName, ret);
            continue;
        }

        /* Copy the pkg from /capps/pkgs/[name]_[tag].tar.gz to
         * /capps/rootfs/[spaceName]/capps/pkgs/[name]_[tag].tar.gz
         */
        sprintf(src,  "%s/%s_%s.tar.gz", sPath, ptr->name, ptr->tag);
        sprintf(dest, "%s/%s/%s/%s_%s.tar.gz", dPath, spaceName, sPath,
                ptr->name, ptr->tag);
        sprintf(runMe, "/bin/cp %s %s", src, dest);
        usys_log_debug("Running command: %s", runMe);
        if ((ret = system(runMe)) < 0) {
            usys_log_error("Unable to copy file src: %s dest: %s Code: %d",
                      src, dest, ret);
            continue;
        }
    }
}

void copy_capps_to_rootfs(Manifest *manifest) {
    
    if (manifest == NULL) return;

    /* boot */
    copy_capps_to_cspace_rootfs(JTAG_BOOT,
                                manifest->boot,
                                DEF_CAPP_PATH,
                                DEF_CSPACE_ROOTFS_PATH);

    /* services */
    copy_capps_to_cspace_rootfs(JTAG_SERVICES,
                                manifest->services,
                                DEF_CAPP_PATH,
                                DEF_CSPACE_ROOTFS_PATH);

    /* reboot */
    copy_capps_to_cspace_rootfs(JTAG_REBOOT,
                                manifest->reboot,
                                DEF_CAPP_PATH,
                                DEF_CSPACE_ROOTFS_PATH);

}
