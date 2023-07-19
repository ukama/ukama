/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <fcntl.h>
#include <string.h>
#include <errno.h>
#include <sys/stat.h>
#include <sys/types.h>

#include "starter.h"

static bool is_valid_capp(char *path) {

    char fileName[CAPP_MAX_BUFFER] = {0};
    struct stat stats;

    if (path == NULL) return USYS_FALSE;

    sprintf(fileName, "%s/%s", path, DEF_CAPP_CONFIG_FILE);

    if (stat(fileName, &stats) == 0) {
        usys_log_debug("Valid config file at: %s", fileName);
        return USYS_TRUE;
    } else {
        usys_log_error("Default config file not found: %s", fileName);
    }

    return USYS_FALSE;
}

/*
 * capp_unpack -- unpack the capp
 *                Currently using tar eventually should be done via
 *                libarchive (TODO)
 */
static bool capp_unpack(char *name, char *tag, char *unpackPath) {

    char runMe[CAPP_MAX_BUFFER] = {0};
    char path[CAPP_MAX_BUFFER]  = {0};
    struct stat stats;

    if (name == NULL || tag == NULL || unpackPath == NULL) return USYS_FALSE;

    stat(unpackPath, &stats);
    if (!S_ISDIR(stats.st_mode)) {
        if(mkdir(unpackPath, 0700) < 0) {
            usys_log_error("Error creating default unpack dir: %s Error: %s",
                           unpackPath, strerror(errno));
            return USYS_FALSE;
        }
    }

    usys_log_debug("Unpacking the capp: %s/%s_%s.tar.gz",
                   DEF_CAPP_PATH, name, tag);
    sprintf(runMe, "/bin/tar xfz %s/%s_%s.tar.gz -C %s",
            DEF_CAPP_PATH,  name, tag, unpackPath);
    if (system(runMe) != 0) {
        usys_log_error("Unable to unpack the capp: %s_%s.tar.gz to %s",
                       name, tag, unpackPath);
        return USYS_FALSE;
    }

    /* check to see if config.json exists */
    sprintf(path, "%s/%s_%s/", unpackPath, name, tag);
    if (is_valid_capp(path)) {
        return USYS_TRUE;
    }

    return USYS_FALSE;
}

static bool unpack_all_capps_to_space_rootfs(CappList *cappList,
                                             char *rootfsPath,
                                             char *cappPath) {

    int ret=0;
    CappList *currentCapp = NULL;
    Capp *capp = NULL;

    char unpackPath[SPACE_MAX_BUFFER] = {0};
    char runMe[MAX_BUFFER]            = {0};

    if (cappList == NULL) return USYS_FALSE;

    /* cleanup existing capps */
    for (currentCapp = cappList;
         currentCapp;
         currentCapp = currentCapp->next) {

        capp = currentCapp->capp;

        if (capp->name == NULL || capp->tag == NULL) continue;

        if (capp->fetch == CAPP_PKG_NOT_FOUND) continue;

        /* existing capps will be at:
         * /capps/rootfs/[contained]/capps/pkgs/unpack
         * delete the whole directory for the space
         */
        sprintf(unpackPath, "%s/%s/%s/unpack",
                rootfsPath,
                capp->space,
                cappPath);

        sprintf(runMe, "/bin/rm -rf %s", unpackPath);
        log_debug("Running command: %s", runMe);

        if ((ret = system(runMe)) != 0) {
            log_error("Unable to execute cmd %s for space: %s Code: %d",
                      runMe, capp->space, ret);
            continue;
        }
    }

    /* unpack the capps to: /capps/rootfs/[contained]/capps/pkgs/unpack/ */
    for (currentCapp = cappList;
         currentCapp;
         currentCapp = currentCapp->next) {

        capp = currentCapp->capp;

        if (capp->name == NULL || capp->tag == NULL) continue;

        if (capp->fetch == CAPP_PKG_NOT_FOUND) continue;

        /* create unpack directory */
        sprintf(unpackPath, "%s/%s/%s/unpack",
                rootfsPath,
                capp->space,
                cappPath);

        sprintf(runMe, "/bin/mkdir -p %s", unpackPath);
        log_debug("Running command: %s", runMe);
        if ((ret = system(runMe)) != 0) {
            log_error("Unable to execute cmd %s for space: %s Code: %d",
                      runMe, capp->space, ret);
            continue;
        }

        /* unpack and verify */
        if (!capp_unpack(capp->name, capp->tag, unpackPath)) {
            return USYS_FALSE;
        } else {
            /* 4: 2 for /, 1 for _ and 1 for NULL */
            capp->rootfs = (char *)malloc(strlen(unpackPath) +
                                          strlen(capp->name) +
                                          strlen(capp->tag)  + 4);

            sprintf(capp->rootfs, "%s/%s/%s/unpack/%s_%s",
                    rootfsPath,
                    capp->space,
                    cappPath,
                    capp->name,
                    capp->tag);
        }
    }

    return USYS_TRUE;
}

bool unpack_all_capps(SpaceList *spaceList) {

    SpaceList *currentSpace;
    
    /* for each space, copy its "capps rootfs" to its own rootfs" */
    for (currentSpace = spaceList;
         currentSpace;
         currentSpace = currentSpace->next) {

        if (currentSpace->space) {
            if (unpack_all_capps_to_space_rootfs(currentSpace->space->cappList,
                                                 DEF_SPACE_ROOTFS_PATH,
                                                 DEF_CAPP_PATH) != USYS_TRUE) {
                usys_log_error("Unable to unpack for space: %s rootfs: %s",
                               currentSpace->space->name,
                               currentSpace->space->rootfs);
                return USYS_FALSE;
            }
        }
    }

    return USYS_TRUE;
}
