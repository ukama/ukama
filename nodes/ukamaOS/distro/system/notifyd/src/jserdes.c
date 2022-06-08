/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "jserdes.h"

#include "errorcode.h"
#include "json_types.h"
#include "node.h"
#include "notify.h"
#include "web_service.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

/* Parser to read real value from JSON object */
bool json_deserialize_real_value(const JsonObj *jObj, double *ivalue) {
    bool ret = USYS_FALSE;

    /* Check if object is number */
    if (json_is_real(jObj)) {
        *ivalue = json_real_value(jObj);
        ret = USYS_TRUE;
    } else if (json_is_integer(jObj)) {
        *ivalue = json_integer_value(jObj);
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read integer value from JSON object */
bool json_deserialize_integer_value(const JsonObj *jObj, int *ivalue) {
    bool ret = USYS_FALSE;

    /* Check if object is number */
    if (json_is_integer(jObj)) {
        *ivalue = json_integer_value(jObj);
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read integer value from JSON object */
bool json_deserialize_integer_object(const JsonObj *obj, const char *key,
                                int *ivalue) {
    bool ret = USYS_FALSE;

    /* Integer Json Object */
    const JsonObj *jIntObj = json_object_get(obj, key);

    /* Check if object is number */
    if (json_is_number(jIntObj)) {
        *ivalue = json_integer_value(jIntObj);
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read integer value from JSON object */
bool json_deserialize_uint32_object(const JsonObj *obj, const char *key,
                               uint32_t *ivalue) {
    bool ret = USYS_FALSE;

    /* Integer Json Object */
    const JsonObj *jIntObj = json_object_get(obj, key);

    /* Check if object is number */
    if (json_is_number(jIntObj)) {
        *ivalue = json_integer_value(jIntObj);
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read uint16_t value from JSON object */
bool json_deserialize_uint16_object(const JsonObj *obj, const char *key,
                               uint16_t *ivalue) {
    bool ret = USYS_FALSE;
    int value = 0;

    ret = json_deserialize_integer_object(obj, key, &value);
    if (ret) {
        *ivalue = (uint16_t)value;
    }

    return ret;
}

/* Parser to read uint8_t value from JSON object */
bool json_deserialize_uint8_object(const JsonObj *obj, const char *key,
                              uint8_t *ivalue) {
    bool ret = USYS_FALSE;
    int value = 0;

    ret = json_deserialize_integer_object(obj, key, &value);
    if (ret) {
        *ivalue = (uint8_t)value;
    }

    return ret;
}

/* Parser to read string value from JSON object */
bool json_deserialize_string_value(JsonObj *obj, char **svalue) {
    bool ret = USYS_FALSE;
    int len = 0;

    /* Check if object is string */
    if (json_is_string(obj)) {
        len = json_string_length(obj);
        svalue = usys_zmalloc(sizeof(char) * (len + 1));
        if (svalue) {
            const char *str = json_string_value(obj);
            usys_strcpy(*svalue, str);
            ret = USYS_TRUE;
        }
    }

    return ret;
}

/* Parser to read string value from JSON object */
bool json_deserialize_string_object(const JsonObj *obj, const char *key,
                               char **svalue) {
    bool ret = USYS_FALSE;

    /* String Json Object */
    JsonObj *jStrObj = json_object_get(obj, key);

    /* Check if object is string */
    if (json_is_string(jStrObj)) {
        int length = json_string_length(jStrObj);
        *svalue = usys_zmalloc(sizeof(char) * (length + 1));
        if (*svalue) {
            const char *str = json_string_value(jStrObj);
            usys_strcpy(*svalue, str);
            ret = USYS_TRUE;
        }
    }

    return ret;
}

/* Wrapper on top of parse_read_string */
bool json_deserialize_string_object_wrapper(const JsonObj *obj, const char *key,
                                       char *str) {
    bool ret = USYS_FALSE;
    char *tstr;
    if (json_deserialize_string_object(obj, key, &tstr)) {
        usys_strcpy(str, tstr);
        usys_free(tstr);
        tstr = NULL;
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read boolean value from JSON object */
bool json_deserialize_boolean_value(const JsonObj *jBoolObj, bool *bvalue) {
    bool ret = USYS_FALSE;

    /* Check if object is number */
    if (json_is_boolean(jBoolObj)) {
        *bvalue = json_boolean_value(jBoolObj);
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser to read boolean value from JSON object */
bool json_deserialize_boolean_object(const JsonObj *obj, const char *key,
                                bool *bvalue) {
    bool ret = USYS_FALSE;

    /* Integer Json Object */
    const JsonObj *jBoolObj = json_object_get(obj, key);

    /* Check if object is number */
    if (json_is_boolean(jBoolObj)) {
        *bvalue = json_boolean_value(jBoolObj);
        ret = USYS_TRUE;
    }

    return ret;
}

/* Parser Error */
void json_deserialize_error(JsonErrObj *jErr, char *msg) {
    if (jErr) {
        usys_log_error("%s. Error: %s ", msg, jErr->text);
    } else {
        usys_log_error("%s. No error info available", msg);
    }
}


json_t *json_encode_value(int type, void *data) {
    JsonObj *json = json_object();
    if (!json) {
        return NULL;
    }

    switch (type) {
    case TYPE_NULL: {
        json = json_null();
        break;
    }
    case TYPE_CHAR: {
        char *value = (char *)data;
        json = json_string(value);
        break;
    }
    case TYPE_BOOL: {
        bool value = *(bool *)data;
        json = json_boolean(value);
        break;
    }
    case TYPE_UINT8: {
        uint8_t value = *(uint8_t *)data;
        json = json_integer(value);
        break;
    }
    case TYPE_INT8: {
        int8_t value = *(int8_t *)data;
        json = json_integer(value);
        break;
    }
    case TYPE_UINT16: {
        uint16_t value = *(uint16_t *)data;
        json = json_integer(value);
        break;
    }
    case TYPE_INT16: {
        int16_t value = *(int16_t *)data;
        json = json_integer(value);
        break;
    }
    case TYPE_UINT32: {
        uint32_t value = *(uint32_t *)data;
        json = json_integer(value);
        break;
    }
    case TYPE_INT32: {
        int32_t value = *(int32_t *)data;
        json = json_integer(value);
        break;
    }
    case TYPE_INT: {
        int value = *(int *)data;
        json = json_integer(value);
        break;
    }
    case TYPE_FLOAT: {
        float value = *(float *)data;
        json = json_real(value);
        break;
    }
    case TYPE_ENUM: {
        int value = *(int *)data;
        json = json_integer(value);
        break;
    }
    case TYPE_DOUBLE: {
        double value = *(double *)data;
        json = json_real(value);
        break;
    }
    case TYPE_STRING: {
        char *value = (char *)data;
        json = json_string(value);
        break;
    }
    default: {
        json = json_null();
    }
    }

    return json;
}

void *json_decode_value(json_t *json, int type) {
    void *data = NULL;

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
        char *value = usys_zmalloc(sizeof(char) + 1);
        if (!value) {
            return NULL;
        }

        if (json_deserialize_string_value(json, &value)) {
            data = value;
        } else {
            usys_free(value);
            return NULL;
        }
        break;
    }
    case TYPE_BOOL: {
        data = usys_zmalloc(sizeof(bool));
        if (!data) {
            return NULL;
        }

        if (!json_deserialize_boolean_value(json, data)) {
            usys_free(data);
            data = NULL;
        }

        break;
    }
    case TYPE_UINT8: {
        int8_t *ndata = usys_zmalloc(sizeof(uint8_t));
        if (!ndata) {
            return NULL;
        }

        int value = 0;
        if (!json_deserialize_integer_value(json, &value)) {
            usys_free(ndata);
            return NULL;
        } else {
            *ndata = (uint8_t)value;
            data = ndata;
        }
        break;
    }
    case TYPE_INT8: {
        int8_t *ndata = usys_zmalloc(sizeof(int8_t));
        if (!ndata) {
            return NULL;
        }

        int value = 0;
        if (!json_deserialize_integer_value(json, &value)) {
            usys_free(ndata);
            return NULL;
        } else {
            *ndata = (int8_t)value;
            data = ndata;
        }
        break;
    }
    case TYPE_UINT16: {
        uint16_t *ndata = usys_zmalloc(sizeof(uint16_t));
        if (!ndata) {
            return NULL;
        }

        int value = 0;
        if (!json_deserialize_integer_value(json, &value)) {
            usys_free(ndata);
            return NULL;
        } else {
            *ndata = (uint16_t)value;
            data = ndata;
        }
        break;
    }
    case TYPE_INT16: {
        int16_t *ndata = usys_zmalloc(sizeof(int16_t));
        if (!ndata) {
            return NULL;
        }

        int value = 0;
        if (!json_deserialize_integer_value(json, &value)) {
            usys_free(ndata);
            return NULL;
        } else {
            *ndata = (int16_t)value;
            data = ndata;
        }
        break;
    }
    case TYPE_UINT32: {
        uint32_t *ndata = usys_zmalloc(sizeof(uint32_t));
        if (!ndata) {
            return NULL;
        }

        int value = 0;
        if (!json_deserialize_integer_value(json, &value)) {
            usys_free(ndata);
            return NULL;
        } else {
            *ndata = (uint32_t)value;
            data = ndata;
        }
        break;
    }
    case TYPE_INT32: {
        int32_t *ndata = usys_zmalloc(sizeof(int32_t));
        if (!ndata) {
            return NULL;
        }

        int value = 0;
        if (!json_deserialize_integer_value(json, &value)) {
            usys_free(ndata);
            return NULL;
        } else {
            *ndata = (int32_t)value;
            data = ndata;
        }
        break;
    }
    case TYPE_INT: {
        data = usys_zmalloc(sizeof(int));
        if (!data) {
            return NULL;
        }

        if (!json_deserialize_integer_value(json, data)) {
            usys_free(data);
            data = NULL;
        }

        break;
    }
    case TYPE_FLOAT: {
        float *ndata = usys_zmalloc(sizeof(float));
        if (!ndata) {
            return NULL;
        }

        double val;
        if (!json_deserialize_real_value(json, &val)) {
            usys_free(ndata);
            return NULL;
        } else {
            *ndata = (float)val;
            data = ndata;
        }
        break;
    }
    case TYPE_ENUM: {
        data = usys_zmalloc(sizeof(int));
        if (!data) {
            return NULL;
        }

        if (!json_deserialize_integer_value(json, data)) {
            usys_free(data);
            data = NULL;
        }

        break;
    }
    case TYPE_DOUBLE: {
        data = usys_zmalloc(sizeof(double));
        if (!data) {
            return NULL;
        }
        if (!json_deserialize_real_value(json, data)) {
            usys_free(data);
            data = NULL;
        }

        break;
    }
    case TYPE_STRING: {
        char *ndata = NULL;
        if (!json_deserialize_string_value(json, &ndata)) {
            data = NULL;
        } else {
            data = ndata;
        }
        break;
    }
    default: {
        json_object_set_new(json, JTAG_VALUE, json_null());
    }
    }

    return data;
}

void deserailize_node_type(int type, char** nodeType) {

    switch(type) {
        case TNODE:
            *nodeType = usys_strdup("TowerNode");
            break;
        case HNODE:
            *nodeType= usys_strdup("HomeNode");
            break;
        case ANODE:
            *nodeType= usys_strdup("AmplifierNode");
            break;
        default:
            *nodeType = NULL;
    }
}

/* Deserialize node Info */
bool json_deserialize_node_info(JsonObj *json, char* nodeId, char* nodeType) {

    bool ret = USYS_FALSE;

    if (!json){
        usys_log_error("No data to deserialize in node info");
        return ret;
    }

    JsonObj* jNodeInfo = json_object_get(json, JTAG_NODE_INFO);
    if (jNodeInfo == NULL) {
        usys_log_error("Missing mandatory %s from JSON", JTAG_NODE_INFO);
        return USYS_FALSE;
    }

    ret = json_deserialize_string_object_wrapper(jNodeInfo, JTAG_UUID, nodeId);
    if (!ret) {
        usys_log_error("Failed to parse Node ID %s in NodeInfo", JTAG_UUID);
        return ret;
    }

    int type = 0;
    ret = json_deserialize_integer_object(jNodeInfo, JTAG_TYPE, &type);
    if (!ret) {
        usys_log_error("Failed to parse Node Type %s in NodeInfo", JTAG_TYPE);
        return ret;
    } else {
        char *nType = NULL;
        deserailize_node_type(type, &nType);
        if (nType) {
            usys_strcpy(nodeType, nType);
            usys_free(nType);
        }
    }

    return ret;

}

/* Deserialize alert received from noded */
bool json_deserialize_noded_alerts(JsonObj *json, NodedNotifDetails* details ) {

    bool ret = USYS_FALSE;

    if (!json){
        usys_log_error("No data to deserialize alerts");
        return ret;
    }

    JsonObj* jNodeInfo = json_object_get(json, JTAG_NOTIFY);
    if (jNodeInfo == NULL) {
        usys_log_error("Missing mandatory %s from JSON", JTAG_NOTIFY);
        return USYS_FALSE;
    }

    ret = json_deserialize_string_object(jNodeInfo, JTAG_SERVICE_NAME,
                    &details->serviceName);
    if (!ret) {
        usys_log_error("Failed to parse mandatory tag %s from Node "
                        "notification", JTAG_SERVICE_NAME);
        return ret;
    }

    ret = json_deserialize_string_object(jNodeInfo, JTAG_SEVERITY,
                        &details->severity);
        if (!ret) {
            usys_log_error("Failed to parse mandatory tag %s from Node "
                            "notification", JTAG_SEVERITY);
            return ret;
        }

    ret = json_deserialize_uint32_object(jNodeInfo, JTAG_EPOCH_TIME,
                    &details->epcohTime);
    if (!ret) {
        usys_log_warn("Failed to parse %s from Node notification",
                        JTAG_EPOCH_TIME);
    }

    ret = json_deserialize_string_object(jNodeInfo, JTAG_UUID,
                    &details->moduleID);
    if (!ret) {
        usys_log_warn("Failed to parse %s from Node notification",
                        JTAG_UUID);
    }

    ret = json_deserialize_string_object(jNodeInfo, JTAG_NAME,
                    &details->deviceName);
    if (!ret) {
        usys_log_warn("Failed to parse %s from Node notification",
                        JTAG_NAME);
    }

    ret = json_deserialize_string_object(jNodeInfo, JTAG_DESCRIPTION,
                    &details->deviceDesc);
    if (!ret) {
        usys_log_warn("Failed to parse %s from Node notification",
                        JTAG_DESCRIPTION);
    }

    ret = json_deserialize_string_object(jNodeInfo, JTAG_PROPERTY_NAME,
                    &details->deviceAttr);
    if (!ret) {
        usys_log_warn("Failed to parse %s from Node notification",
                        JTAG_PROPERTY_NAME);
    }

    ret = json_deserialize_integer_object(jNodeInfo, JTAG_DATA_TYPE,
                    &details->dataType);
    if (!ret) {
        usys_log_warn("Failed to parse %s from Node notification",
                        JTAG_DATA_TYPE);
    }

    const JsonObj *jValue = json_object_get(jNodeInfo, JTAG_VALUE);
    if (!jValue){
        usys_log_warn("Failed to parse %s from Node notification",
                        JTAG_VALUE);
        return ret;
    }

    details->deviceAttrValue = (double*)usys_calloc(1, sizeof(double));
    if (details->deviceAttrValue) {
        ret = json_deserialize_real_value(jValue, details->deviceAttrValue);
        if (!ret) {
            usys_log_error("Failed to parse %s from Node notification",
                            JTAG_VALUE);
        }
    }

    ret = json_deserialize_string_object(jNodeInfo, JTAG_UNITS,
                    &details->units);
    if (!ret) {
        usys_log_error("Failed to parse %s from Node notification",
                        JTAG_UNITS);
    }

    return ret;

}


int json_serialize_error(JsonObj **json, int code, const char *str) {
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return ERR_JSON_CRETATION_ERR;
    }

    json_object_set_new(*json, JTAG_ERROR, json_object());

    JsonObj *jError = json_object_get(*json, JTAG_ERROR);
    if (jError) {
        json_object_set_new(jError, JTAG_ERROR_CODE, json_integer(code));

        json_object_set_new(jError, JTAG_ERROR_CSTRING, json_string(str));

    } else {
        return ERR_JSON_CRETATION_ERR;
    }

    return ret;
}

int json_serialize_api_list(JsonObj **json, WebServiceAPI *apiList,
                            uint16_t count) {
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return ERR_JSON_CRETATION_ERR;
    }

    if (!apiList) {
        return ERR_JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_API_LIST, json_array());

    JsonObj *jApiArr = json_object_get(*json, JTAG_API_LIST);
    if (jApiArr) {
        for (int iter = 0; iter < count; iter++) {
            json_t *jApi = json_object();

            json_object_set_new(jApi, JTAG_METHOD,
                                json_string(apiList[iter].method));

            json_object_set_new(jApi, JTAG_URL_EP,
                                json_string(apiList[iter].endPoint));

            /* Add element to array */
            json_array_append(jApiArr, jApi);
            json_decref(jApi);
        }

    } else {
        return ERR_JSON_CRETATION_ERR;
    }

    return ret;
}

/* Serialize alert details from the noded
 * This section of tt=he notification is specific to each service */
int json_serialize_noded_alert_details(JsonObj **json,
                NodedNotifDetails* details ) {
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return ERR_JSON_CRETATION_ERR;
    }

    if (!details) {
        return ERR_JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_UUID,
                    json_string(details->moduleID));

    json_object_set_new(*json, JTAG_NAME,
                    json_string(details->deviceName));

    json_object_set_new(*json, JTAG_DESCRIPTION,
                    json_string(details->deviceDesc));

    json_object_set_new(*json, JTAG_PROPERTY_NAME,
                    json_string(details->deviceAttr));

    json_object_set_new(*json, JTAG_EPOCH_TIME,
                    json_integer(details->dataType));

    //TODO: Remove hard coding
    json_object_set_new(*json, JTAG_VALUE,
                    json_encode_value(TYPE_DOUBLE, details->deviceAttrValue));

    json_object_set_new(*json, JTAG_UNITS, json_string(details->units));


    return ret;
}

/* Serialize notification to be forwaded to the remote server */
int json_serialize_notification(JsonObj **json, JsonObj* details,
                Notification* notif) {
    int ret = JSON_ENCODING_OK;

    *json = json_object();
    if (!json) {
        return ERR_JSON_CRETATION_ERR;
    }

    if (!details) {
        return ERR_JSON_NO_VAL_TO_ENCODE;
    }

    json_object_set_new(*json, JTAG_SERVICE_NAME,
                    json_string(notif->serviceName));

    json_object_set_new(*json, JTAG_NOTIFICATION_TYPE,
                    json_string(notif->notificationType));

    json_object_set_new(*json, JTAG_NODE_ID,
                    json_string(notif->nodeId));

    json_object_set_new(*json, JTAG_NODE_TYPE,
                    json_string(notif->nodeType));

    json_object_set_new(*json, JTAG_NOTIF_SEVERITY,
                    json_string(notif->severity));

    json_object_set_new(*json, JTAG_DESCRIPTION,
                        json_string(notif->description));

    json_object_set_new(*json, JTAG_EPOCH_TIME,
                    json_integer(notif->epcohTime));

    json_object_set_new(*json, JTAG_NOTIF_DETAILS,
                    details);

    return ret;
}




