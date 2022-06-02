/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef JSERDES_H
#define JSERDES_H

#include <jansson.h>

#include "json_types.h"

#include "usys_types.h"

#define JSON_OK                        STATUS_OK
#define JSON_FAILURE                   STATUS_NOTOK
#define JSON_ENCODING_OK               JSON_OK
#define JSON_DECODING_OK               JSON_OK

#define JSON_STRING  1
#define JSON_INTEGER 2

/**
 * @fn      int json_serialize_error(JsonObj**, int, const char*)
 * @brief   Serialize error to report to client in JSON body .
 *
 * @param   obj
 * @param   code
 * @param   str
 * @return  On Success, JSON_ENCODING_OK (STATUS_OK)
 *          On Failure, NodeD JSON error code
 */
int json_serialize_error(JsonObj** obj, int code , const char* str );


/**
 * @fn      bool parser_read_boolean_object(const JsonObj*, const char*, bool*)
 * @brief   Parses the object which contain boolean object with key
 *          supplied in argument.
 *
 * @param   obj
 * @param   key
 * @param   bvalue
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_boolean_object(const JsonObj *obj, const char* key,
                bool *bvalue);

/**
 * @fn      bool parser_read_boolean_value(const JsonObj*, bool*)
 * @brief   Reads the boolean value from the json object.
 *
 * @param   jBoolObj
 * @param   bvalue
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_boolean_value(const JsonObj *jBoolObj,
                bool *bvalue);

/**
 * @fn      bool parser_read_integer_object(const JsonObj*, const char*, int*)
 * @brief   Parses the object which contain integer object with key
 *          supplied in argument.
 *
 * @param   obj
 * @param   key
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_integer_object(const JsonObj *obj, const char* key,
                int *ivalue);

/**
 * @fn      bool parser_read_integer_value(const JsonObj*, int*)
 * @brief   Reads the integer value from the json object.
 *
 * @param   obj
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_integer_value(const JsonObj *obj, int *ivalue);

/**
 * @fn      bool parser_read_real_value(const JsonObj*, double*)
 * @brief   Reads the real (double) value from the json object.
 *
 * @param   jObj
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_real_value(const JsonObj *jObj, double *ivalue);

/**
 * @fn      bool parser_read_string_object(const JsonObj*, const char*, char**)
 * @brief   Parses the object which contain string object with key
 *          supplied in argument and return the memory pointer to string
 *          in svalue. It's user responsibility to free the memory later.
 *
 * @param   obj
 * @param   key
 * @param   svalue
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_string_object(const JsonObj *obj, const char* key,
                char **svalue);

/**
 * @fn      bool parser_read_string_object_wrapper(const JsonObj*, const char*, char*)
 * @brief   Parses the object which contain string object with key
 *          supplied in argument. After reading the value it copies the string
 *          to the str.
 *
 * @param   obj
 * @param   key
 * @param   str
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_string_object_wrapper(const JsonObj *obj, const char* key,
                char* str);

/**
 * @fn      bool parser_read_string_value(JsonObj*, char**)
 * @brief   reads the value from the string object and return the memory
 *          pointer to it. It's caller responsibility to free the memory after
 *          use.
 *
 * @param   obj
 * @param   svalue
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_string_value(JsonObj *obj, char **svalue);

/**
 * @fn      bool parser_read_uint16_object(const JsonObj*, const char*, uint16_t*)
 * @brief   Parses the object which contain integer object with key
 *          supplied in argument abd typecast value to uint16
 *
 * @param   obj
 * @param   key
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_uint16_object(const JsonObj *obj, const char* key,
                uint16_t *ivalue);
/**
 * @fn      bool parser_read_uint32_object(const JsonObj*, const char*, uint32_t*)
 * @brief   Parses the object which contain integer object with key
 *          supplied in argument abd typecast value to uint32
 *
 * @param   obj
 * @param   key
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_uint32_object(const JsonObj *obj, const char* key,
                uint32_t *ivalue);

/**
 * @fn      bool parser_read_uint8_object(const JsonObj*, const char*, uint8_t*)
 * @brief   Parses the object which contain integer object with key
 *          supplied in argument abd typecast value to uint8
 *
 * @param   obj
 * @param   key
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool parser_read_uint8_object(const JsonObj *obj, const char* key,
                uint8_t *ivalue);
/**
 * @fn      void parser_error(JsonErrObj*, char*)
 * @brief   Logs the parser error occurred while parsing json file.
 *
 * @param   jErr
 * @param   msg
 */
void parser_error(JsonErrObj *jErr, char* msg);

#endif /* JSERDES_H_ */
