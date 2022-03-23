/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "jserdes.h"

#include "property.h"
#include "web_service.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

int json_serialize_error(JsonObj** json, int code, const char* str ) {
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    json_object_set_new(*json, JTAG_ERROR, json_object());

    JsonObj* jError = json_object_get(*json, JTAG_ERROR);
    if (jError) {

        json_object_set_new(jError, JTAG_ERROR_CODE, json_integer(code));

        json_object_set_new(jError, JTAG_ERROR_CSTRING, json_string(str));

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_version(JsonObj** json, Version* ver) {
    int ret = JSON_ENCODING_OK;

    if (*json) {

        json_object_set_new(*json, JTAG_MAJOR_VERSION, json_integer(ver->major));

        json_object_set_new(*json, JTAG_MINOR_VERSION, json_integer(ver->minor));

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_module_cfg(JsonObj** json, ModuleCfg* mCfg, uint8_t count){
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    if (!mCfg) {
        return JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_MODULE_CONFIG, json_array());

    JsonObj* jMCfgArr = json_object_get(*json, JTAG_MODULE_CONFIG);
    if (jMCfgArr) {

        for(int iter = 0; iter < count; iter++) {
            json_t* jMCfg = json_object();

            json_object_set_new(jMCfg, JTAG_NAME, json_string(mCfg[iter].devName));

            json_object_set_new(jMCfg, JTAG_DESCRIPTION,
                            json_string(mCfg[iter].devDesc));

            json_object_set_new(jMCfg, JTAG_TYPE,
                            json_integer(mCfg[iter].devType));

            json_object_set_new(jMCfg, JTAG_CLASS,
                            json_integer(mCfg[iter].devClass));

            json_object_set_new(jMCfg, JTAG_DEV_SYSFS_FILE,
                            json_string(mCfg[iter].sysFile));
            /* Add element to array */
            json_array_append(jMCfgArr, jMCfg);
            json_decref(jMCfg);

        }
    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_module_info(JsonObj** json, ModuleInfo* mInfo){
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    if (!mInfo) {
        return JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_UNIT_INFO, json_object());

    JsonObj* jMInfo = json_object_get(*json, JTAG_UNIT_INFO);
    if (jMInfo) {

        json_object_set_new(jMInfo, JTAG_UUID, json_string(mInfo->uuid));

        json_object_set_new(jMInfo, JTAG_NAME, json_string(mInfo->name));

        json_object_set_new(jMInfo, JTAG_TYPE, json_integer(mInfo->module));

        json_object_set_new(jMInfo, JTAG_PART_NUMBER, json_string(mInfo->partNo));

        json_object_set_new(jMInfo, JTAG_HW_VERSION, json_string(mInfo->hwVer));

        json_object_set_new(jMInfo, JTAG_MAC, json_string(mInfo->mac));

        json_object_set_new(jMInfo, JTAG_PROD_SW_VERSION, json_object());

        JsonObj* jPVer = json_object_get(jMInfo, JTAG_PROD_SW_VERSION);
        if (jPVer) {
            ret = json_serialize_version(&jPVer, &mInfo->pSwVer);
            if (ret != JSON_ENCODING_OK) {
                return ret;
            }
        }

        json_object_set_new(jMInfo, JTAG_SW_VERISION, json_object());

        JsonObj* jSVer = json_object_get(jMInfo, JTAG_SW_VERISION);
        if (jSVer) {
            ret = json_serialize_version(&jSVer, &mInfo->pSwVer);
            if (ret != JSON_ENCODING_OK) {
                return ret;
            }
        }

        json_object_set_new(jMInfo, JTAG_MFG_DATE, json_string(mInfo->mfgDate));

        json_object_set_new(jMInfo, JTAG_MFG_NAME, json_string(mInfo->mfgName));

        json_object_set_new(jMInfo, JTAG_DEVICE_COUNT, json_integer(mInfo->devCount));

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_unit_cfg(JsonObj** json, UnitCfg* uCfg, uint8_t count){
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    if (!uCfg) {
        return JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_UNIT_CONFIG, json_array());

    JsonObj* jUCfgArr = json_object_get(*json, JTAG_UNIT_CONFIG);
    if (jUCfgArr) {

        for(int iter = 0; iter < count; iter++) {
            json_t* jUCfg = json_object();

            json_object_set_new(jUCfg, JTAG_UUID, json_string(uCfg->modUuid));

            json_object_set_new(jUCfg, JTAG_NAME, json_string(uCfg->modName));

            json_object_set_new(jUCfg, JTAG_INVT_SYSFS_FILE,
                            json_string(uCfg->sysFs));

            /* Add element to array */
            json_array_append(jUCfgArr, jUCfg);
            json_decref(jUCfg);

        }

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_unit_info(JsonObj** json, UnitInfo* uInfo){
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    if (!uInfo) {
        return JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_UNIT_INFO, json_object());

    JsonObj* jUInfo = json_object_get(*json, JTAG_UNIT_INFO);
    if (jUInfo) {

        json_object_set_new(jUInfo, JTAG_UUID, json_string(uInfo->uuid));

        json_object_set_new(jUInfo, JTAG_NAME, json_string(uInfo->name));

        json_object_set_new(jUInfo, JTAG_TYPE, json_integer(uInfo->unit));

        json_object_set_new(jUInfo, JTAG_PART_NUMBER,
                        json_string(uInfo->partNo));

        json_object_set_new(jUInfo, JTAG_SKEW, json_string(uInfo->skew));

        json_object_set_new(jUInfo, JTAG_MAC, json_string(uInfo->mac));

        json_object_set_new(jUInfo, JTAG_PROD_SW_VERSION, json_object());

        JsonObj* jPVer = json_object_get(jUInfo, JTAG_PROD_SW_VERSION);
        if (jPVer) {
            ret = json_serialize_version(&jPVer, &uInfo->swVer);
            if (ret != JSON_ENCODING_OK) {
                return ret;
            }
        }

        json_object_set_new(jUInfo, JTAG_SW_VERISION, json_object());

        JsonObj* jSVer = json_object_get(jUInfo, JTAG_SW_VERISION);
        if (jSVer) {
            ret = json_serialize_version(&jSVer, &uInfo->pSwVer);
            if (ret != JSON_ENCODING_OK) {
                return ret;
            }
        }

        json_object_set_new(jUInfo, JTAG_ASM_DATE,
                        json_string(uInfo->assmDate));

        json_object_set_new(jUInfo, JTAG_OEM_NAME,
                        json_string(uInfo->oemName));

        json_object_set_new(jUInfo, JTAG_MODULE_COUNT,
                        json_integer(uInfo->modCount));

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

int json_serialize_api_list(JsonObj** json, WebServiceAPI* apiList, uint16_t count) {
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    if (!apiList) {
        return JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_API_LIST, json_array());

    JsonObj* jApiArr = json_object_get(*json, JTAG_API_LIST);
    if (jApiArr) {

        for(int iter = 0; iter < count; iter++) {
            json_t* jApi = json_object();

            json_object_set_new(jApi, JTAG_METHOD, json_string(apiList[iter].method));

            json_object_set_new(jApi, JTAG_URL_EP, json_string(apiList[iter].endPoint));

            /* Add element to array */
            json_array_append(jApiArr, jApi);
            json_decref(jApi);

        }

    } else {
        return JSON_CREATION_ERR;
    }

    return ret;
}

json_t* json_encode_value(int type, void* data) {

    JsonObj *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    switch (type) {
        case TYPE_NULL: {
            json_object_set_new(json, JTAG_VALUE, json_null());
            break;
        }
        case TYPE_CHAR: {
            char* value = (char*)data;
            json_object_set_new(json, JTAG_VALUE, json_string(value));
            break;
        }
        case TYPE_BOOL: {
            bool value = *(bool*)data;
            json_object_set_new(json, JTAG_VALUE, json_boolean(value));
            break;
        }
        case TYPE_UINT8: {
            uint8_t value = *(uint8_t*)data;
            json_object_set_new(json, JTAG_VALUE, json_integer(value));
            break;
        }
        case TYPE_INT8: {
            int8_t value = *(int8_t*)data;
            json_object_set_new(json, JTAG_VALUE, json_integer(value));
            break;
        }
        case TYPE_UINT16: {
            uint16_t value = *(uint16_t*)data;
            json_object_set_new(json, JTAG_VALUE, json_integer(value));
            break;
        }
        case TYPE_INT16: {
            int16_t value = *(int16_t*)data;
            json_object_set_new(json, JTAG_VALUE, json_integer(value));
            break;
            break;
        }
        case TYPE_UINT32: {
            uint32_t value = *(uint32_t*)data;
            json_object_set_new(json, JTAG_VALUE, json_integer(value));
            break;
        }
        case TYPE_INT32: {
            int32_t value = *(int32_t*)data;
            json_object_set_new(json, JTAG_VALUE, json_integer(value));
            break;
        }
        case TYPE_INT: {
            int value = *(int*)data;
            json_object_set_new(json, JTAG_VALUE, json_integer(value));
            break;
        }
        case TYPE_FLOAT: {
            float value = *(float*)data;
            json_object_set_new(json, JTAG_VALUE, json_real(value));
            break;
        }
        case TYPE_ENUM: {
            int value = *(int*)data;
            json_object_set_new(json, JTAG_VALUE, json_integer(value));
            break;
        }
        case TYPE_DOUBLE: {
            double value = *(double*)data;
            json_object_set_new(json, JTAG_VALUE, json_real(value));
            break;
        }
        case TYPE_STRING: {
            char* value = (char*)data;
            json_object_set_new(json, JTAG_VALUE, json_string(value));
            break;
        }
        default: {
            json_object_set_new(json, JTAG_VALUE, json_null());
        }
    }

    return json;
}

void* json_decode_value(json_t* json, int type) {
    void* data = NULL;

    if (!json) {
        return data;
    }

    switch (type) {
        case TYPE_NULL: {
            data = NULL;
            break;
        }
        case TYPE_CHAR: {
            /* Allocating extar one byte beacuse of '/0' */
            char *value = usys_zmalloc(sizeof(char)+1);
            if (!data){
                data = NULL;
            }

            if (parser_read_string_value(json, &value)) {
                data = value;
            } else {
               usys_free(data);
               data = NULL;
            }
            break;
        }
        case TYPE_BOOL: {
            bool *data = usys_zmalloc(sizeof(bool));
            if (!data){
                data = NULL;
            }

            if (!parser_read_boolean_value(json, data)) {
                usys_free(data);
                data = NULL;
            }

            break;
        }
        case TYPE_UINT8: {
            uint8_t *data = usys_zmalloc(sizeof(uint8_t));
            if (!data){
                data = NULL;
            }

            int value = 0;
            if (!parser_read_integer_value(json, &value)) {
                usys_free(data);
                data = NULL;
            } else {
                data = (uint8_t)value;
            }
            break;
        }
        case TYPE_INT8: {
            int8_t *data = usys_zmalloc(sizeof(int8_t));
            if (!data){
                data = NULL;
            }

            int value = 0;
            if (!parser_read_integer_value(json, &value)) {
                usys_free(data);
                data = NULL;
            } else {
                data = (int8_t)value;
            }
            break;
        }
        case TYPE_UINT16: {
            uint16_t *data = usys_zmalloc(sizeof(uint16_t));
            if (!data){
                return NULL;
            }

            int value = 0;
            if (!parser_read_integer_value(json, &value)) {
                usys_free(data);
                data = NULL;
            } else {
                data = (uint16_t)value;
            }
            break;
        }
        case TYPE_INT16: {
            int16_t *data = usys_zmalloc(sizeof(int16_t));
            if (!data){
                return NULL;
            }

            int value = 0;
            if (!parser_read_integer_value(json, &value)) {
                usys_free(data);
                data = NULL;
            } else {
                data = (int16_t)value;
            }
            break;
        }
        case TYPE_UINT32: {
            uint32_t *data = usys_zmalloc(sizeof(uint32_t));
            if (!data){
                return NULL;
            }

            int value = 0;
            if (!parser_read_integer_value(json, &value)) {
                usys_free(data);
                data = NULL;
            } else {
                data = (uint32_t)value;
            }
            break;
        }
        case TYPE_INT32: {
            int32_t *data = usys_zmalloc(sizeof(int32_t));
            if (!data){
                return NULL;
            }

            int value = 0;
            if (!parser_read_integer_value(json, &value)) {
                usys_free(data);
                data = NULL;
            } else {
                data = (int32_t)value;
            }
            break;
        }
        case TYPE_INT: {
            int *data = usys_zmalloc(sizeof(int));
            if (!data){
                return NULL;
            }

            if (!parser_read_integer_value(json, data)) {
                usys_free(data);
                data = NULL;
            }

            break;
        }
        case TYPE_FLOAT: {
            float *data = usys_zmalloc(sizeof(float));
            if (!data){
                return NULL;
            }

            double val;
            if (!parser_read_real_value(json, &val)) {
                usys_free(data);
                data = NULL;
            } else {
                *data = (float)val;
            }
            break;
        }
        case TYPE_ENUM: {
            int *data = usys_zmalloc(sizeof(int));
            if (!data){
                return NULL;
            }

            if (!parser_read_integer_value(json, data)) {
                usys_free(data);
                data = NULL;
            }

            break;
        }
        case TYPE_DOUBLE: {
            double *data = usys_zmalloc(sizeof(double));
            if (!data){
                return NULL;
            }

            if (!parser_read_real_value(json, &data)) {
                usys_free(data);
                data = NULL;
            }

            break;
        }
        case TYPE_STRING: {
            char* data = NULL;
            if (!parser_read_string_value(json, &data)) {
                data = NULL;
            }
            break;
        }
        default: {
            json_object_set_new(json, JTAG_VALUE, json_null());
        }
    }

    return data;
}

int json_serialize_sensor_data( JsonObj** json, const char* name,
                const char* desc, int type, void* data) {
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return JSON_CREATION_ERR;
    }

    if (!data) {
        return JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_NAME, json_string(name));

    json_object_set_new(*json, JTAG_DESCRIPTION, json_string(name));

    json_object_set_new(*json, JTAG_TYPE, json_integer(type));

    json_object_set_new(*json, JTAG_VALUE, json_encode_value(type, data));

    return ret;
}

int json_deserialize_sensor_data( JsonObj* json, char* name, char* desc,  int type, void* data) {
    int ret = JSON_DECODING_OK;

    if (!json) {
        return JSON_PARSER_ERR;
    }

   if (!parser_read_string_object(json, JTAG_NAME, &name) ) {
       return JSON_PARSER_ERR;
   }

   if (!parser_read_integer_object(json, JTAG_TYPE, &type) ) {
       return JSON_PARSER_ERR;
   }

   if (!parser_read_string_object(json, JTAG_DESCRIPTION, &desc) ) {
       return JSON_PARSER_ERR;
   }

   JsonObj* jData = json_object_get(json, JTAG_VALUE);
   data = json_decode_value(jData, type);
   if(!data) {
       return JSON_PARSER_ERR;
   }

   return ret;
}



