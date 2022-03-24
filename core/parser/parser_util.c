/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "json_types.h"
#include "schema.h"

#include "usys_types.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

/* Parser to read real value from JSON object */
bool parser_read_real_value(const JsonObj *jObj, double *ivalue) {
    bool ret = USYS_FALSE;

    /* Check if object is number */
    if (json_is_real(jObj)) {
        *ivalue = json_real_value(jObj) ;
        ret = USYS_TRUE;
    } else if (json_is_integer(jObj)) {
        *ivalue = json_integer_value(jObj) ;
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read integer value from JSON object */
bool parser_read_integer_value(const JsonObj *jObj, int *ivalue) {
    bool ret = USYS_FALSE;

    /* Check if object is number */
    if (json_is_integer(jObj)) {
        *ivalue = json_integer_value(jObj) ;
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
    if (json_is_number(jIntObj)) {
        *ivalue = json_integer_value(jIntObj) ;
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read integer value from JSON object */
bool parser_read_uint32_object(const JsonObj *obj, const char* key,
                uint32_t *ivalue) {
    bool ret = USYS_FALSE;

    /* Integer Json Object */
    const JsonObj *jIntObj = json_object_get(obj, key);

    /* Check if object is number */
    if (json_is_number(jIntObj)) {
        *ivalue = json_integer_value(jIntObj) ;
        ret = USYS_TRUE;
    }

    return ret;
}


/* Parser to read uint16_t value from JSON object */
bool parser_read_uint16_object(const JsonObj *obj, const char* key,
                uint16_t *ivalue) {
    bool ret = USYS_FALSE;
    int value = 0;

    ret = parser_read_integer_object(obj, key, &value);
    if (ret) {
        *ivalue = (uint16_t) value;
    }

    return ret;
}

/* Parser to read uint8_t value from JSON object */
bool parser_read_uint8_object(const JsonObj *obj, const char* key,
                uint8_t *ivalue) {
    bool ret = USYS_FALSE;
    int value = 0;

    ret = parser_read_integer_object(obj, key, &value);
    if (ret) {
        *ivalue = (uint8_t)value;
    }

    return ret;
}


/* Parser to read string value from JSON object */
bool parser_read_string_value(JsonObj *obj, char **svalue) {
    bool ret = USYS_FALSE;
    int len = 0;

    /* Check if object is string */
    if (json_is_string(obj)) {

        len = json_string_length(obj);
        svalue = usys_zmalloc(sizeof(char) * (len+1));
        if (svalue) {
            const char *str = json_string_value(obj);
            usys_strcpy(*svalue, str);
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
    JsonObj *jStrObj = json_object_get(obj, key);

    /* Check if object is number */
    if (jStrObj && json_is_string(jStrObj)) {

        int length = json_string_length(jStrObj);
        *svalue = usys_zmalloc(sizeof(char) * (length+1));
        if (*svalue) {
            const char *str = json_string_value(jStrObj);
            usys_strcpy(*svalue, str);
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
        tstr = NULL;
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read boolean value from JSON object */
bool parser_read_boolean_value(const JsonObj *jBoolObj,
                bool *bvalue) {
    bool ret = USYS_FALSE;

    /* Check if object is number */
    if (json_is_boolean(jBoolObj)) {
        *bvalue = json_boolean_value(jBoolObj) ;
        ret = USYS_TRUE;
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
    if (json_is_boolean(jBoolObj)) {
        *bvalue = json_boolean_value(jBoolObj) ;
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser Error */
void parser_error(JsonErrObj *jErr, char* msg) {
    if (jErr) {
        usys_log_error("%s. Error: %s ", msg, jErr->text);
    } else {
        usys_log_error("%s. No error info available", msg);
    }
}

/* Parse version */
Version *parse_version(const JsonObj *jVersion) {
    const JsonObj *jMajor = NULL;
    const JsonObj *jMinor = NULL;

    Version *pversion = usys_zmalloc(sizeof(Version));
    if (pversion) {

        /* Major version */
        uint8_t major = 0;
        if (!parser_read_uint8_object(jVersion, JTAG_MAJOR_VERSION,
                        &major)) {
            usys_free(pversion);
            return NULL;
        } else {
            pversion->major = (uint8_t) major;
        }

        /* Minor version */
        uint8_t minor = 0;
        if (!parser_read_uint8_object(jVersion, JTAG_MINOR_VERSION,
                        &minor)) {
            usys_free(pversion);
            return NULL;
        } else {
            pversion->minor = (uint8_t) minor;
        }

    }

    return pversion;
}

