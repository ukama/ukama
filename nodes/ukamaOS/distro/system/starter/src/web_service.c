/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <signal.h>

#include "web_service.h"
#include "http_status.h"
#include "config.h"

#include "starter.h"
#include "usys_log.h"

extern SpaceList *gSpaceList;

static char *capp_status_str(int status) {

    char *str;

    switch(status) {
    case CAPP_RUNTIME_NO_EXEC:
        str = "Not running";
        break;
    case CAPP_RUNTIME_EXEC:
        str = "Running";
        break;
    case CAPP_RUNTIME_DONE:
        str = "Done";
        break;
    default:
        str = "Uknown";
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

            if (strcmp(cappList->capp->name, cappName) == 0)

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

    return NULL;
}

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {
    
    ulfius_set_string_body_response(response, HttpStatus_Unauthorized,
                                    HttpStatusStr(HttpStatus_Unauthorized));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_status(const URequest *request,
                              UResponse *response,
                              void *epConfig) {

    char   *cappName=NULL, *spaceName=NULL;
    Capp   *capp = NULL;
    int    status=-1;

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
            status = CAPP_RUNTIME_NO_EXEC;
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

int web_service_cb_post_update(const URequest *request,
                               UResponse *response,
                               void *epConfig) {

    char *cappName = NULL, *tag = NULL;
    char *spaceName = NULL;
    Capp *capp = NULL;
    int  status;

    spaceName = u_map_get(request->map_url, "space");
    cappName  = u_map_get(request->map_url, "name");
    tag       = u_map_get(request->map_url, "tag");

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
        if (capp->runtime->status == CAPP_RUNTIME_NO_EXEC ||
            capp->runtime->status == CAPP_RUNTIME_DONE) {
            ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
            return U_CALLBACK_CONTINUE;
        }
    }

    /* Only if the capp is running */
    status = killpg(capp->runtime->pid, SIGTERM);
    if ( status == 0 ){
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
