/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef MANIFEST_H
#define MANIFEST_H

#include "usys_types.h"

#define MAX_MANIFEST_FILE_SIZE 1000000 /* 1MB max. file size */
#define MAX_CAPPS              64

#define JTAG_VERSION  "version"
#define JTAG_TARGET   "target"
#define JTAG_CAPPS    "capps"
#define JTAG_SPACES   "spaces"

#define JTAG_NAME     "name"
#define JTAG_TAG      "tag"
#define JTAG_RESTART  "restart"
#define JTAG_SPACE    "space"

typedef struct _cappsManifest {

    char *name;      /* Name of the cApp */
    char *tag;       /* cApp tag/version */
    char *space;     /* group it belongs to */
    int  restart;    /* 1: yes, always restart. 0: No */

    struct _cappsManifest *next; /* Next in the list */
} CappsManifest;

typedef struct _spacesManifest {
    
    char *name;

    struct _spacesManifest *next;
} SpacesManifest;

typedef struct {

    char           *version;
    char           *target;
    CappsManifest  *cappsManifest;
    SpacesManifest *spacesManifest;
} Manifest;

/* Function headers. */
bool read_manifest_file(Manifest **manifest, char *fileName);
void free_manifest(Manifest *ptr);

#endif /* MANIFEST_H */
