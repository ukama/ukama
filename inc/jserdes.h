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

#include "json_types.h"
#include "web_service.h"

#include "usys_types.h"


#define JSON_OK                        0
#define JSON_ENCODING_OK               JSON_OK
#define JSON_DECODING_OK               JSON_OK


/* JSON Ser-Des Error Codes */
#define JSON_FAILURE                   -1
#define JSON_CREATION_ERR              -2000
#define JSON_NO_VAL_TO_ENCODE          -2001
#define JSON_INVALID                   -1000
#define JSON_PARSER_ERR                -1001
#define JSON_UNEXPECTED_TAG            -1002
#define JSON_BAD_REQ                   -1003

int json_deserialize_sensor_data( JsonObj* json, char** name, char** desc,
                int* type, void** data);
int json_serialize_api_list(JsonObj** json, WebServiceAPI* apiList, uint16_t count);
int json_serialize_error(JsonObj** obj, int code , const char* str );
int json_serialize_module_cfg(JsonObj** obj, ModuleCfg* uCfg, uint8_t count);
int json_serialize_module_info(JsonObj** obj, ModuleInfo* uInfo);
int json_serialize_sensor_data(JsonObj** json, const char* name,
                const char* desc, int type, void* data);
int json_serialize_unit_cfg(JsonObj** obj, UnitCfg* uCfg, uint8_t count);
int json_serialize_unit_info(JsonObj** obj, UnitInfo* uInfo);

#endif /* INC_JSERDES_H_ */
