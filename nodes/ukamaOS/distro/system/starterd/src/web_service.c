/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <signal.h>

#include "web_service.h"
#include "http_status.h"
#include "json_types.h"
#include "config.h"

#include "starter.h"

#include "usys_log.h"
#include "usys_mem.h"

#include "version.h"

extern SpaceList *gSpaceList;
extern void json_free(JsonObj** json);
extern bool json_serialize_add_capp_to_array(JsonObj **json,
                                             char *space,
                                             char *name,
                                             char *tag,
                                             char *status,
                                             int  pid);

extern bool find_matching_space(SpaceList **spaceList, char *name, Space **space);

static char *capp_status_str(int status) {

    char *str;

    switch(status) {
    case CAPP_RUNTIME_PEND:
        str = "Pending";
        break;
    case CAPP_RUNTIME_EXEC:
        str = "Active";
        break;
    case CAPP_RUNTIME_DONE:
        str = "Done";
        break;
    case CAPP_RUNTIME_FAILURE:
        str = "Failure";
        break;
    case CAPP_RUNTIME_UNKNOWN:
        str = "Unknown";
        break;
    default:
        str = "Unknown";
        break;
    }

    return str;
}

static Capp *find_matching_capp(char *spaceName, char *cappName, char *tag) {

    SpaceList *spacePtr = NULL;
    CappList  *cappList  = NULL;

    for (spacePtr = gSpaceList;
         spacePtr;
         spacePtr = spacePtr->next) {

        if (strcmp(spacePtr->space->name, spaceName) != 0)
            continue;

        for (cappList=spacePtr->space->cappList;
             cappList;
             cappList=cappList->next) {

            if (strcmp(cappList->capp->name, cappName) == 0) {

                if (tag != NULL) {
                    if (strcmp(cappList->capp->tag, tag) == 0)
                        return cappList->capp;
                    else
                        continue;
                } else {
                    return cappList->capp;
                }
            }
        }
    }

    return NULL;
}

static int add_new_capp_to_space(char *spaceName,
                                 char *cappName,
                                 char *cappTag) {

    SpaceList *currentSpaceList = NULL, *newSpaceList = NULL;
    CappList  *newCappList = NULL;
    bool      addSpace = USYS_TRUE;

    for (currentSpaceList = gSpaceList;
         currentSpaceList;
         currentSpaceList = currentSpaceList->next) {

        if (strcmp(currentSpaceList->space->name, spaceName) == 0) {
            addSpace = USYS_FALSE;
        }
    }

    if (addSpace) {
        /* add new space */
        newSpaceList        = (SpaceList *) calloc(1, sizeof(SpaceList));
        newSpaceList->space = (Space *) calloc(1, sizeof(Space));

        newSpaceList->space->name     = strdup(spaceName);
        newSpaceList->space->rootfs   = NULL;
        newSpaceList->space->cappList = NULL;
        newSpaceList->next            = NULL;

        /* Forward to last spot on the list*/
        for (currentSpaceList = gSpaceList;
             currentSpaceList->next;
             currentSpaceList = currentSpaceList->next) ;

        currentSpaceList->next = newSpaceList;
    }

    /* Now find the matching space and add to it */
    for (currentSpaceList = gSpaceList;
         currentSpaceList;
         currentSpaceList = currentSpaceList->next) {

        if (strcmp(currentSpaceList->space->name, spaceName) == 0) {

            newCappList = (CappList *)calloc(1, sizeof(CappList));

            newCappList->capp          = (Capp *)calloc(1, sizeof(Capp));
            newCappList->capp->name    = strdup(cappName);
            newCappList->capp->tag     = strdup(cappTag);
            newCappList->capp->rootfs  = NULL;
            newCappList->capp->space   = strdup(spaceName);
            newCappList->capp->restart = USYS_FALSE;
            newCappList->capp->fetch   = CAPP_PKG_NOT_FOUND;

            newCappList->next = currentSpaceList->space->cappList;
            currentSpaceList->space->cappList = newCappList;

            return USYS_TRUE;
        }
    }

    return USYS_FALSE;
}

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *data) {

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {
    
    ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_status(const URequest *request,
                              UResponse *response,
                              void *epConfig) {

    const char   *cappName=NULL, *spaceName=NULL;
    Capp *capp =NULL;
    int  status=-1;

    cappName  = u_map_get(request->map_url, "name");
    spaceName = u_map_get(request->map_url, "space");

    capp = find_matching_capp(spaceName, cappName, NULL);
    if (capp == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
        return U_CALLBACK_CONTINUE;
    }

    if (capp->runtime) {
            status = capp->runtime->status;
    } else {
            status = CAPP_RUNTIME_PEND;
    }

    if (status == -1) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
    } else {
        ulfius_set_string_body_response(response, HttpStatus_OK,
                                        capp_status_str(status));
    }

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_all_capps_status(const URequest *request,
                                        UResponse *response,
                                        void *epConfig) {

    SpaceList *spacePtr = NULL;
    CappList  *cappList = NULL;
    JsonObj   *json     = NULL;
    char      *status   = NULL;
    char      *jStr     = NULL;
    int       pid       = 0;

    json = json_object();
    json_object_set_new(json, JTAG_CAPPS, json_array());
    if (json == NULL) {
        ulfius_set_string_body_response(response,
                               HttpStatus_InternalServerError,
                               HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    for (spacePtr = gSpaceList;
         spacePtr;
         spacePtr = spacePtr->next) {

        /* for each space find all the capps */
        for (cappList=spacePtr->space->cappList;
             cappList;
             cappList=cappList->next) {

            if (cappList->capp->runtime) {
                status = capp_status_str(cappList->capp->runtime->status);
                pid    = cappList->capp->runtime->pid;
            } else {
                status = capp_status_str(CAPP_RUNTIME_PEND);
                pid    = 0;
            }

            json_serialize_add_capp_to_array(&json,
                                             spacePtr->space->name,
                                             cappList->capp->name,
                                             cappList->capp->tag,
                                             status, pid);
        }
    }

    jStr = json_dumps(json, 0);
    if (jStr == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
    } else {
        ulfius_set_json_body_response(response, HttpStatus_OK, json);
    }

    usys_free(jStr);
    json_free(&json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_update(const URequest *request,
                               UResponse *response,
                               void *epConfig) {

    const char *cappName = NULL, *tag = NULL;
    const char *spaceName = NULL;
    Capp *capp = NULL;
    int  status;

    spaceName = u_map_get(request->map_url, "space");
    cappName  = u_map_get(request->map_url, "name");
    tag       = u_map_get(request->map_url, "tag");

    if (strcmp(spaceName, SPACE_BOOT) == 0 ||
        strcmp(spaceName, SPACE_REBOOT) == 0) {

        ulfius_set_string_body_response(response,
                                        HttpStatus_Forbidden,
                                        HttpStatusStr(HttpStatus_Forbidden));
        return U_CALLBACK_CONTINUE;
    }

    capp = find_matching_capp(spaceName, cappName, tag);
    if (capp == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
        return U_CALLBACK_CONTINUE;
    }

    /* Terminate and set fetch flag */
    /* Only if the capp is running */
    if (capp->runtime != NULL) {
        if (capp->runtime->status == CAPP_RUNTIME_EXEC) {
            status = killpg(capp->runtime->pid, SIGTERM);
            if ( status == 0 ){
                usys_log_debug("Capp update - %s:%s", cappName, tag);
                usys_log_debug("SIGTERM send to capp: %s:%s", cappName, tag);
            } else {
                usys_log_debug("Unable to kill capp: %s:%s",
                               capp->name, capp->tag);
                ulfius_set_string_body_response(response,
                                HttpStatus_InternalServerError,
                                HttpStatusStr(HttpStatus_InternalServerError));
                return U_CALLBACK_CONTINUE;
            }
        }
    }

    /* set the flag */
    capp->fetch = CAPP_PKG_NOT_FOUND;

    ulfius_set_empty_body_response(response, HttpStatus_Accepted);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_terminate(const URequest *request,
                                  UResponse *response,
                                  void *epConfig) {

    char   *cappName=NULL, *spaceName=NULL;
    Capp   *capp = NULL;
    int    status;

    cappName  = u_map_get(request->map_url, "name");
    spaceName = u_map_get(request->map_url, "space");

    if (strcmp(spaceName, SPACE_BOOT) == 0 ||
        strcmp(spaceName, SPACE_REBOOT) == 0) {

        ulfius_set_string_body_response(response,
                                        HttpStatus_Forbidden,
                                        HttpStatusStr(HttpStatus_Forbidden));
        return U_CALLBACK_CONTINUE;
    }

    capp = find_matching_capp(spaceName, cappName, NULL);
    if (capp == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
        return U_CALLBACK_CONTINUE;
    }

    if (capp->runtime == NULL) {
        /* capp is not yet gone through the runtime setup */
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    } else {
        /* already done or not executing */
        if (capp->runtime->status == CAPP_RUNTIME_PEND ||
            capp->runtime->status == CAPP_RUNTIME_FAILURE ||
            capp->runtime->status == CAPP_RUNTIME_DONE) {
            ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
            return U_CALLBACK_CONTINUE;
        }
    }

    /* Only if the capp is running */
    status = killpg(capp->runtime->pid, SIGTERM);
    if (status == 0){
        usys_log_debug("SIGTERM send to capp: %s:%s", capp->name, capp->tag);
        ulfius_set_string_body_response(response, HttpStatus_Accepted,
                                        HttpStatusStr(HttpStatus_Accepted));
    } else {
        usys_log_debug("Unable to kill capp: %s:%s", capp->name, capp->tag);
        ulfius_set_string_body_response(response,
                              HttpStatus_InternalServerError,
                              HttpStatusStr(HttpStatus_InternalServerError));
    }

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_exec(const URequest *request,
                             UResponse *response,
                             void *epConfig) {

    char *cappName = NULL, *tag = NULL;
    char *spaceName = NULL;
    Capp *capp = NULL;

    spaceName = u_map_get(request->map_url, "space");
    cappName  = u_map_get(request->map_url, "name");
    tag       = u_map_get(request->map_url, "tag");

    if (strcmp(spaceName, SPACE_BOOT) == 0 ||
        strcmp(spaceName, SPACE_REBOOT) == 0) {

        ulfius_set_string_body_response(response,
                                        HttpStatus_Forbidden,
                                        HttpStatusStr(HttpStatus_Forbidden));
        return U_CALLBACK_CONTINUE;
    }

    capp = find_matching_capp(spaceName, cappName, tag);
    if (capp != NULL) {
        if (capp->runtime) {
            if (capp->runtime->status == CAPP_RUNTIME_EXEC) {
                usys_log_debug("Can't exec already running capp %s:%s:%s",
                               spaceName, cappName, tag);
                ulfius_set_string_body_response(response,
                                         HttpStatus_Forbidden,
                                         HttpStatusStr(HttpStatus_Forbidden));
                return U_CALLBACK_CONTINUE;
            }
        }

        /* Set the fetch flag so it can automatically start in next cycle */
        capp->fetch = CAPP_PKG_NOT_FOUND;

        ulfius_set_empty_body_response(response, HttpStatus_Accepted);
        return U_CALLBACK_CONTINUE;
    }

    /* Add new capp */
    if (add_new_capp_to_space(spaceName, cappName, tag)) {
        ulfius_set_empty_body_response(response, HttpStatus_Accepted);
    } else {
        ulfius_set_string_body_response(response,
                               HttpStatus_InternalServerError,
                               HttpStatusStr(HttpStatus_InternalServerError));
    }

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_MethodNotAllowed,
                                    HttpStatusStr(HttpStatus_MethodNotAllowed));
    return U_CALLBACK_CONTINUE;
}
