/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

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

    char      *cappName=NULL, *space=NULL;
    SpaceList *spacePtr = NULL;
    CappList  *cappList  = NULL;
    int       status=-1;

    cappName = u_map_get(request->map_url, "name");
    space    = u_map_get(request->map_url, "space");

    for (spacePtr = gSpaceList;
         spacePtr;
         spacePtr = spacePtr->next) {

        if (strcmp(spacePtr->space->name, space) != 0)
            continue;

        for (cappList=spacePtr->space->cappList;
             cappList;
             cappList=cappList->next) {

            if (strcmp(cappList->capp->name, cappName) != 0)
                continue;

            if (cappList->capp->runtime)
                status = cappList->capp->runtime->status;
            else
                status = CAPP_RUNTIME_NO_EXEC;
        }
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

    char *name=NULL;

    name = u_map_get(request->map_url, "name");
    ulfius_set_empty_body_response(response, HttpStatus_Accepted);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_restart(const URequest *request,
                                UResponse *response,
                                void *epConfig) {

    char *name=NULL;

    name = u_map_get(request->map_url, "name");
    ulfius_set_empty_body_response(response, HttpStatus_Accepted);

    return U_CALLBACK_CONTINUE;
}




