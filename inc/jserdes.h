/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_JSERDES_H_
#define INC_JSERDES_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "json_types.h"
#include "web_service.h"

#include "usys_types.h"

#define JSON_OK                        STATUS_OK
#define JSON_FAILURE                   STATUS_NOTOK
#define JSON_ENCODING_OK               JSON_OK
#define JSON_DECODING_OK               JSON_OK

/**
 * @fn      int json_deserialize_sensor_data(JsonObj*, const char**, const char**, int*, void**)
 * @brief   Deserialize the sensor data supplied by  user in HTTP body for
 *          configuring sensor.
 *
 * @param   json
 * @param   name
 * @param   desc
 * @param   dataType
 * @param   data
 * @return  On Success, JSON_DECODING_OK (STATUS_OK)
 *          On Failure, NodeD JSON error code
 */
int json_deserialize_sensor_data( JsonObj* json, const char** name,
                const char** desc, int* dataType, void** data);

/**
 * @fn      int json_serialize_alert_data(JsonObj**, const char*, const char*, const char*, const char*, int, void*, char*)
 * @brief    Serialize alert info into the JSON body.
 *
 * @param   json
 * @param   modUuid
 * @param   devName
 * @param   devDesc
 * @param   propName
 * @param   type
 * @param   data
 * @param   units
 * @return  On Success, JSON_ENCODING_OK (STATUS_OK)
 *          On Failure, NodeD JSON error code
 */
int json_serialize_alert_data(JsonObj **json, const char* modUuid,
                const char *devName, const char *devDesc, const char *propName,
                int type, void *data, char* units);
/**
 * @fn      int json_serialize_api_list(JsonObj**, WebServiceAPI*, uint16_t)
 * @brief   Serialize API list into the JSON body.
 *
 * @param   json
 * @param   apiList
 * @param   count
 * @return  On Success, JSON_ENCODING_OK (STATUS_OK)
 *          On Failure, NodeD JSON error code
 */
int json_serialize_api_list(JsonObj** json, WebServiceAPI* apiList,
                uint16_t count);
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
 * @fn      int json_serialize_module_cfg(JsonObj**, ModuleCfg*, uint8_t)
 * @brief   Serialize the module configuration into JSON body.
 *
 * @param   obj
 * @param   uCfg
 * @param   count
 * @return  On Success, JSON_ENCODING_OK (STATUS_OK)
 *          On Failure, NodeD JSON error code
 */
int json_serialize_module_cfg(JsonObj** obj, ModuleCfg* uCfg, uint8_t count);
/**
 * @fn      int json_serialize_module_info(JsonObj**, ModuleInfo*)
 * @brief   Serialize the module information into JSON body.
 *
 * @param   obj
 * @param   uInfo
 * @return  On Success, JSON_ENCODING_OK (STATUS_OK)
 *          On Failure, NodeD JSON error code
 */
int json_serialize_module_info(JsonObj** obj, ModuleInfo* uInfo);
/**
 * @fn      int json_serialize_sensor_data(JsonObj**, const char*, const char*,
 *          int, void*)
 * @brief   Serialize the read value and its data type for sensor into JSON body
 *
 * @param   json
 * @param   name
 * @param   desc
 * @param   type
 * @param   data
 * @return  On Success, JSON_ENCODING_OK (STATUS_OK)
 *          On Failure, NodeD JSON error code
 */
int json_serialize_sensor_data(JsonObj** json, const char* name,
                const char* desc, int type, void* data);
/**
 * @fn      int json_serialize_unit_cfg(JsonObj**, UnitCfg*, uint8_t)
 * @brief   Serialize the unit configuration into JSON body.
 *
 * @param   obj
 * @param   uCfg
 * @param    count
 * @return  On Success, JSON_ENCODING_OK (STATUS_OK)
 *          On Failure, NodeD JSON error code
 */
int json_serialize_unit_cfg(JsonObj** obj, UnitCfg* uCfg, uint8_t count);

/**
 * @fn      int json_serialize_unit_info(JsonObj**, UnitInfo*)
 * @brief   Serialize the unit information into JSON body.
 *
 * @param   obj
 * @param   uInfo
 * @return  On Success, JSON_ENCODING_OK (STATUS_OK)
 *          On Failure, NodeD JSON error code
 */
int json_serialize_unit_info(JsonObj** obj, UnitInfo* uInfo);

#ifdef __cplusplus
}
#endif

#endif /* INC_JSERDES_H_ */
