/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <jansson.h>

#include "web_service.h"
#include "json_types.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

bool json_serialize_add_capp_to_array(JsonObj **json,
                                         char *name,
                                         char *tag,
                                         char *status,
                                         int  pid) {

    JsonObj *jArray = NULL;
    JsonObj *jCapp  = NULL;

    if (*json == NULL) return USYS_FALSE;

    jArray = json_object_get(*json, JTAG_CAPPS);

    if (jArray) {

        jCapp = json_object();
        if (jCapp == NULL) return USYS_FALSE;

        json_object_set_new(jCapp, JTAG_NAME,   json_string(name));
        json_object_set_new(jCapp, JTAG_TAG,    json_string(tag));
        json_object_set_new(jCapp, JTAG_STATUS, json_string(status));
        json_object_set_new(jCapp, JTAG_PID,    json_integer(pid));

        json_array_append_new(jArray, jCapp);

        return USYS_TRUE;
    }

    return USYS_FALSE;
}

/* Decrement json references */
void json_free(JsonObj** json) {
    if (*json){
        json_decref(*json);
        *json = NULL;
    }
}
