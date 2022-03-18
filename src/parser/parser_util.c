/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "json_types.h"
#include "usys_types.h"
#include "usys_mem.h"

/* Parser to read integer value from JSON object */
bool parser_read_integer_value(const JsonObj *obj, int *ivalue) {
    bool ret = USYS_FALSE;

    /* Check if object is number */
    if (json_is_number(obj)) {
        *ivalue = json_integer_value(obj) ;
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read integer value from JSON object */
bool parser_read_integer_object(const JsonObj *obj, const char* key,
                int *ivalue) {
    bool ret = USYS_FALSE;

    /* Integer Json Object */
    const JsonObj *jIntObj = json_object_get(obj, key);

    /* Check if object is number */
    if (jIntObj && json_is_number(obj)) {
        *ivalue = json_integer_value(obj) ;
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read string value from JSON object */
bool parser_read_string_value(const JsonObj *obj, char *svalue) {
    bool ret = USYS_FALSE;
    int len = 0;

    /* Check if object is string */
    if (json_is_string(obj)) {
        len = json_string_length(obj);

        svalue = usys_malloc(sizeof(char) * len);
        if (svalue) {
            usys_memset(svalue, '\0', sizeof(char) * len);
            char *str = json_string_value(obj);
            usys_strcpy(svalue, str);
            json_decref(obj);
            ret = USYS_TRUE;
        }

    }

    return ret;
}

/* Parser to read string value from JSON object */
bool parser_read_string_object(const JsonObj *obj, const char* key,
                char **svalue) {
    bool ret = USYS_FALSE;

    /* String Json Object */
    const JsonObj *jStrObj = json_object_get(obj, key);

    /* Check if object is number */
    if (jStrObj && json_is_string(obj)) {
        int length = json_string_length(obj);

        *svalue = usys_malloc(sizeof(char) * length);
        if (*svalue) {
            usys_memset(*svalue, '\0', sizeof(char) * length);
            char *str = json_string_value(obj);
            usys_strcpy(*svalue, str);
            json_decref(obj);
            ret = USYS_TRUE;
        }
    }

    return ret;
}

/* Wrapper on top of parse_read_string */
bool parser_read_string_object_wrapper(const JsonObj *obj, const char* key,
                char* str) {
    bool ret = USYS_FALSE;
    char *tstr;
    if (parser_read_string_object(obj, key, &tstr)) {
        usys_strcpy(str, tstr);
        usys_free(tstr);
    }

    return ret;
}

/* Parser to read boolean value from JSON object */
bool parser_read_boolean_object(const JsonObj *obj, const char* key,
                bool *bvalue) {
    bool ret = USYS_FALSE;

    /* Integer Json Object */
    const JsonObj *jBoolObj = json_object_get(obj, key);

    /* Check if object is number */
    if (jBoolObj && json_is_boolean(obj)) {
        *bvalue = json_boolean_value(obj) ;
        ret = USYS_TRUE;
    }

    return ret;
}

