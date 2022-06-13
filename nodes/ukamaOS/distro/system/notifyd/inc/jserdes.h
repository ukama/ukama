/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef JSERDES_H
#define JSERDES_H

#ifdef __cplusplus
extern "C" {
#endif

#include <jansson.h>

#include "json_types.h"
#include "notify.h"
#include "web_service.h"
#include "usys_types.h"

#define JSON_OK                        STATUS_OK
#define JSON_FAILURE                   STATUS_NOTOK
#define JSON_ENCODING_OK               JSON_OK
#define JSON_DECODING_OK               JSON_OK

/**
 * @fn      bool json_deserialize_boolean_object(const JsonObj*, const char*, bool*)
 * @brief   Parses the object which contain boolean object with key
 *          supplied in argument.
 *
 * @param   obj
 * @param   key
 * @param   bvalue
 * @return  On success, true
 *          On failure, false
 */
bool json_deserialize_boolean_object(const JsonObj *obj, const char* key,
                bool *bvalue);

/**
 * @fn      bool json_deserialize_boolean_value(const JsonObj*, bool*)
 * @brief   Reads the boolean value from the json object.
 *
 * @param   jBoolObj
 * @param   bvalue
 * @return  On success, true
 *          On failure, false
 */
bool json_deserialize_boolean_value(const JsonObj *jBoolObj,
                bool *bvalue);

/**
 * @fn      bool json_deserialize_integer_object(const JsonObj*, const char*, int*)
 * @brief   Parses the object which contain integer object with key
 *          supplied in argument.
 *
 * @param   obj
 * @param   key
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool json_deserialize_integer_object(const JsonObj *obj, const char* key,
                int *ivalue);

/**
 * @fn      bool json_deserialize_integer_value(const JsonObj*, int*)
 * @brief   Reads the integer value from the json object.
 *
 * @param   obj
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool json_deserialize_integer_value(const JsonObj *obj, int *ivalue);

/**
 * @fn      bool json_deserialize_real_value(const JsonObj*, double*)
 * @brief   Reads the real (double) value from the json object.
 *
 * @param   jObj
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool json_deserialize_real_value(const JsonObj *jObj, double *ivalue);

/**
 * @fn      bool json_deserialize_string_object(const JsonObj*, const char*, char**)
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
bool json_deserialize_string_object(const JsonObj *obj, const char* key,
                char **svalue);

/**
 * @fn      bool json_deserialize_string_object_wrapper(const JsonObj*, const char*, char*)
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
bool json_deserialize_string_object_wrapper(const JsonObj *obj, const char* key,
                char* str);

/**
 * @fn      bool json_deserialize_string_value(JsonObj*, char**)
 * @brief   reads the value from the string object and return the memory
 *          pointer to it. It's caller responsibility to free the memory after
 *          use.
 *
 * @param   obj
 * @param   svalue
 * @return  On success, true
 *          On failure, false
 */
bool json_deserialize_string_value(JsonObj *obj, char **svalue);

/**
 * @fn      bool json_deserialize_uint16_object(const JsonObj*, const char*, uint16_t*)
 * @brief   Parses the object which contain integer object with key
 *          supplied in argument abd typecast value to uint16
 *
 * @param   obj
 * @param   key
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool json_deserialize_uint16_object(const JsonObj *obj, const char* key,
                uint16_t *ivalue);
/**
 * @fn      bool json_deserialize_uint32_object(const JsonObj*, const char*, uint32_t*)
 * @brief   Parses the object which contain integer object with key
 *          supplied in argument abd typecast value to uint32
 *
 * @param   obj
 * @param   key
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool json_deserialize_uint32_object(const JsonObj *obj, const char* key,
                uint32_t *ivalue);

/**
 * @fn      bool json_deserialize_uint8_object(const JsonObj*, const char*, uint8_t*)
 * @brief   Parses the object which contain integer object with key
 *          supplied in argument abd typecast value to uint8
 *
 * @param   obj
 * @param   key
 * @param   ivalue
 * @return  On success, true
 *          On failure, false
 */
bool json_deserialize_uint8_object(const JsonObj *obj, const char* key,
                uint8_t *ivalue);
/**
 * @fn      void json_deserialize_error(JsonErrObj*, char*)
 * @brief   Logs the parser error occurred while parsing json file.
 *
 * @param   jErr
 * @param   msg
 */
void json_deserialize_error(JsonErrObj *jErr, char* msg);

/**
 * @fn      bool json_deserialize_node_info(JsonObj*, char*, char*)
 * @brief   Deserialize Node Serial Id and Type from the node info.
 *
 * @param   json
 * @param   nodeId
 * @param   nodeType
 * @return  On success, true
 *          On failure, false
 */
bool json_deserialize_node_info(JsonObj *json, char* nodeId, char* nodeType);

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
 * @fn      int json_serialize_api_list(JsonObj**, WebServiceAPI*, uint16_t)
 * @brief   Serialize the API list exposed by notify.
 *
 * @param   json
 * @param   apiList
 * @param   count
 * @return
 */
int json_serialize_api_list(JsonObj **json, WebServiceAPI *apiList,
                            uint16_t count);

/**
 * @fn      int json_serialize_noded_notif_details(JsonObj**, NodedNotifDetails*)
 * @brief   Serializes the noded alert details into JSON.
 *
 * @param   json
 * @param   details
 * @return  On success, JSON_ENCODING_OK
 *          On failure, Non zero value
 */
int json_serialize_noded_notif_details(JsonObj **json,
                NodedNotifDetails* details );
/**
 * @fn      int json_serialize_notification(JsonObj**, JsonObj*, Notification*)
 * @brief   Serializes the notification data into JSON.
 *
 * @param   json
 * @param   details
 * @param   notif
 * @return  On success, JSON_ENCODING_OK
 *          On failure, Non zero value
 */
int json_serialize_notification(JsonObj **json, JsonObj* details,
                Notification* notif);

/**
 * @fn      bool json_deserialize_noded_notif(JsonObj*, NodedNotifDetails*)
 * @brief   Deserailize node alert details received by notify service.
 *
 * @param   json
 * @param   details
 * @return  On success, TRUE
 *          On failure, FALSE
 */
bool json_deserialize_noded_notif(JsonObj *json, NodedNotifDetails* details );

/**
 * @fn      bool json_deserialize_generic_notification(JsonObj*,
 *              ServiceNotifDetails*)
 * @brief   deserialize generic notifications.
 *
 * @param   json
 * @param   details
 * @return  On success, TRUE
 *          On failure, FALSE
 */
bool json_deserialize_generic_notification(JsonObj *json,
                ServiceNotifDetails* details );

/**
 * @fn      int json_serialize_generic_details(JsonObj**, ServiceNotifDetails*)
 * @brief   serialize generic notification details to send to remote server.
 *
 * @param   json
 * @param   details
 * @return  On success, TRUE
 *          On failure, FALSE
 */
int json_serialize_generic_details(JsonObj **json,
                ServiceNotifDetails* details );

/**
 * @fn      void json_free(JsonObj*)
 * @brief   Free the json object
 *
 * @param   json
 */
void json_free(JsonObj** json);

#ifdef __cplusplus
}
#endif
#endif /* JSERDES_H_ */
