/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string.h>
#include <sys/stat.h>

#include "starter.h"
#include "manifest.h"

static int need_to_fetch_capp_from_hub(char *spaceName, char *cappPkg) {

    struct stat stats;

    if (stat(cappPkg, &stats) == 0) {
        return CAPP_PKG_FOUND;
    }

    if (strcmp(spaceName, SPACE_BOOT) == 0 ||
        strcmp(spaceName, SPACE_REBOOT) == 0) {
        usys_log_error("Mandatory capp not found. space: %s capp: %s",
                       spaceName, cappPkg);
    }

    return CAPP_PKG_NOT_FOUND;
}

bool find_matching_space(SpaceList **spaceList, char *name, Space **space) {

    SpaceList *ptr;

    for (ptr = *spaceList; ptr; ptr=ptr->next) {
        if (strcmp(ptr->space->name, name) == 0) {
            *space = ptr->space;
            return USYS_TRUE;
        }
    }

    return USYS_FALSE;
}

static void add_capp_to_space_list(SpaceList **spaceList,
                                   CappsManifest *cPtr,
                                   char *spaceName) {
    
    SpaceList *currentSpaceList, *newSpaceList;
    CappList *newCappList;
    Space *newSpace;
    
    if (*spaceList == NULL) {
        /* SpaceList is empty, create a new SpaceList */
        SpaceList *newSpaceList = (SpaceList *) calloc(1, sizeof(SpaceList));
        newSpaceList->space     = (Space *) calloc(1, sizeof(Space));

        newSpaceList->space->name     = strdup(spaceName);
        newSpaceList->space->rootfs   = NULL;
        newSpaceList->space->cappList = NULL;
        newSpaceList->next            = NULL;
        
        *spaceList = newSpaceList;
    }
    
    currentSpaceList = *spaceList;
    while (currentSpaceList != NULL) {

        if (strcmp(currentSpaceList->space->name, spaceName) == 0) {
            /* Space found, add capp to the cappList of the space */
            newCappList = (CappList *)calloc(1, sizeof(CappList));

            newCappList->capp          = (Capp *)calloc(1, sizeof(Capp));
            newCappList->capp->name    = strdup(cPtr->name);
            newCappList->capp->tag     = strdup(cPtr->tag);
            newCappList->capp->rootfs  = NULL;
            newCappList->capp->space   = strdup(cPtr->space);
            newCappList->capp->restart = cPtr->restart;
            newCappList->capp->runtime = NULL;
            
            newCappList->next = currentSpaceList->space->cappList;
            currentSpaceList->space->cappList = newCappList;
            
            return;
        }
        currentSpaceList = currentSpaceList->next;
    }
    
    /* Space not found, create a new space and add capp to the 
     * cappList of the new space
     */
    newSpace         = (Space *)malloc(sizeof(Space));
    newSpace->name   = strdup(spaceName);
    newSpace->rootfs = NULL;
    
    newSpace->cappList  = (CappList *)malloc(sizeof(CappList));
    newCappList         = newSpace->cappList;
    newCappList->capp   = (Capp *)malloc(sizeof(Capp));
    
    newCappList->capp->name    = strdup(cPtr->name);
    newCappList->capp->tag     = strdup(cPtr->tag);
    newCappList->capp->rootfs  = NULL;
    newCappList->capp->space   = strdup(cPtr->space);
    newCappList->capp->restart = cPtr->restart;
    newCappList->capp->runtime = NULL;
    newSpace->cappList->next   = NULL;
    
    newSpaceList        = (SpaceList *)malloc(sizeof(SpaceList));
    newSpaceList->space = newSpace;
    newSpaceList->next  = *spaceList;
    
    *spaceList = newSpaceList;
}

void process_manifest_file(SpaceList **spaceList, Manifest *manifest) {

    SpacesManifest *sPtr = NULL;
    CappsManifest  *cPtr = NULL;

    sPtr = manifest->spacesManifest;
    cPtr = manifest->cappsManifest;

    /* For each space, find its capps and add them to the list */
    while (sPtr) {
        for (cPtr = manifest->cappsManifest; cPtr; cPtr = cPtr->next) {
            if (strcmp(sPtr->name, cPtr->space) == 0) {
                add_capp_to_space_list(spaceList, cPtr, sPtr->name);
            }
        }
        sPtr = sPtr->next;
    }
}

void print_spaces_list(SpaceList *spaceList) {

    SpaceList* currentSpaceList;
    CappList* currentCappList;
    Capp*     currentCapp;

    currentSpaceList = spaceList;
    
    while (currentSpaceList != NULL) {
        usys_log_debug("Space: %s", currentSpaceList->space->name);
        
        currentCappList = currentSpaceList->space->cappList;
        while (currentCappList != NULL) {
            currentCapp = currentCappList->capp;
            usys_log_debug("      cApp Name: %s", currentCapp->name);
            usys_log_debug("      cApp Tag: %s", currentCapp->tag);
            usys_log_debug("      cApp RootFS: %s", currentCapp->rootfs);
            usys_log_debug("      cApp Restart: %s",
                           (currentCapp->restart == 1 ? "true" : "false"));
            
            currentCappList = currentCappList->next;
        }
        currentSpaceList = currentSpaceList->next;
    }
}

/*
 * copy_capps_to_space_rootfs --
 *                             sPath: /capps/pkgs
 *                             dPath: /capps/rootfs
 */
static void copy_capps_to_space_rootfs(char *spaceName,
                                       CappList *cappList,
                                       char *sPath,
                                       char *dPath) {

    int ret=0;
    CappList *currentCapp = NULL;
    Capp *capp = NULL;

    char src[SPACE_MAX_BUFFER]  = {0};
    char dest[SPACE_MAX_BUFFER] = {0};
    char runMe[MAX_BUFFER]      = {0};

    if (cappList == NULL) return;

    for (currentCapp = cappList;
         currentCapp;
         currentCapp = currentCapp->next) {

        capp = currentCapp->capp;
        
        if (capp->name == NULL || capp->tag == NULL ) {
            usys_log_error("Invalid capp entry. Ignoring");
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
        sprintf(src,  "%s/%s_%s.tar.gz", sPath, capp->name, capp->tag);

        /* Flag if capp needs to be fetched from hub via wimc */
        capp->fetch = need_to_fetch_capp_from_hub(spaceName, src);
        if (capp->fetch == CAPP_PKG_NOT_FOUND) {
            continue;
        }

        sprintf(dest, "%s/%s/%s/%s_%s.tar.gz", dPath, spaceName, sPath,
                capp->name, capp->tag);
        sprintf(runMe, "/bin/cp %s %s", src, dest);
        usys_log_debug("Running command: %s", runMe);
        if ((ret = system(runMe)) < 0) {
            usys_log_error("Unable to copy file src: %s dest: %s code: %d",
                      src, dest, ret);
            continue;
        }
    }
}

void copy_capps_to_rootfs(SpaceList *spaceList) {

    SpaceList *currentSpace;

    /* for each space, copy its "capps rootfs" to its own rootfs" */
    for (currentSpace = spaceList;
         currentSpace;
         currentSpace = currentSpace->next) {

        if (currentSpace->space) {

            copy_capps_to_space_rootfs(currentSpace->space->name,
                                       currentSpace->space->cappList,
                                       DEF_CAPP_PATH,
                                       DEF_SPACE_ROOTFS_PATH);
        }
    }
}
