/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
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
