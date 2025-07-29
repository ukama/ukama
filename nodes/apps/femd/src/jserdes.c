/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include "jserdes.h"
#include "femd.h"
#include "json_types.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

void json_log(json_t *json) {
    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        usys_log_debug("json str: %s", str);
        free(str);
    }
}

static bool get_json_entry(json_t *json, char *key, json_type type,
                           char **strValue, int *intValue,
                           double *doubleValue) {
    json_t *jEntry = NULL;

    if (json == NULL || key == NULL) return USYS_FALSE;

    jEntry = json_object_get(json, key);
    if (jEntry == NULL) {
        usys_log_error("Missing %s key in json", key);
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
    default:
        usys_log_error("Invalid type for json key-value: %d", type);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

int json_serialize_error(JsonObj **json, int code, const char *str) {
    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_ERROR, json_object());
    JsonObj *jError = json_object_get(*json, JTAG_ERROR);
    if (jError) {
        json_object_set_new(jError, JTAG_ERROR_CODE, json_integer(code));
        json_object_set_new(jError, JTAG_ERROR_CSTRING, json_string(str));
    }

    return JSON_OK;
}

int json_serialize_success(JsonObj **json, const char *message) {
    *json = json_object();
    if (!*json) {
        return ERR_FEMD_JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_STATUS, json_string("success"));
    if (message) {
        json_object_set_new(*json, JTAG_MESSAGE, json_string(message));
    }

    return JSON_OK;
}

void json_free(JsonObj **json) {
    if (*json) {
        json_decref(*json);
        *json = NULL;
    }
}