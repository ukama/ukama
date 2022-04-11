/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "web_service.h"

#include "device.h"
#include "errorcode.h"
#include "inventory.h"
#include "jserdes.h"
#include "ledger.h"
#include "property.h"
#include "service.h"
#include "ulfius.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

UInst serverInst;
static char gNotifServer[128] = { 0 };
static uint16_t endPointCount = 0;
WebServiceAPI gApi[MAX_END_POINTS] = { 0 };

/**
 * @fn      int alert_details(DevObj*, Property*, uint16_t*)
 * @brief   Read the property table for the requested device object
 *
 * @param   obj
 * @param   prop
 * @param   pCount
 * @return  On Success STATUS_OK
 *          On failure STATUS_NOTOK
 */
int device_object_property_table(DevObj* obj, Property **prop, uint16_t* pCount) {
    int ret = STATUS_OK;

    /* read property count */
    ret = ldgr_read_prop_count(obj, pCount);
    if (ret != STATUS_OK) {
        return STATUS_NOK;
    }

    if (pCount == 0)  {
        return STATUS_NOK;
    }

    *prop = usys_zmalloc(sizeof(Property) * (*pCount));
    if (*prop) {
        /* read properties */
        int ret = ldgr_read_prop(obj, *prop);
        if (ret) {
            usys_free(*prop);
            *prop = NULL;
            ret = STATUS_NOK;
        }
    } else {

        usys_log_error("Web Service Failed to allocate memory for property "
                        "table. Error %s",
                              usys_error(errno));
        ret =  STATUS_NOK;
    }

    return ret;
}

/**
 * @fn      void web_service_alert_cb(DevObj*, AlertCallBackData**, int*)
 * @brief   This is called from the ISR context should be released as soon as
 *          possible to start monitor again. In this ISR context  we print alert
 *          info which is for demo purpose only. We have to add web client here
 *          to pass it to notification service
 *
 * @param   obj
 * @param   alertCbData
 * @param   count
 */
void web_service_alert_cb(DevObj *obj, AlertCallBackData **alertCbData, int *count) {
    if (*alertCbData) {
        Property* prop = NULL;
        uint16_t pCount = 0;
        JsonObj *json = NULL;
        URequest alertNotification;
        UResponse alertNotificationResp;
        AlertCallBackData *adata = *alertCbData;
        usys_log_debug(
            "Web service alert callback function for Device name %s Disc: %s "
            "Module UUID %s called with %d alerts.",
            obj->name, obj->desc, obj->modUuid, *count);

        int ret = device_object_property_table(obj, &prop, &pCount);
        if (ret != STATUS_OK) {
            usys_log_error("Web Service Failed to get  alert property details "
                            "from property table.");
            goto cleanup;
        }

        if (!prop) {
            goto cleanup;
        }

        int pidx = adata->pidx;
        int didx = prop[pidx].depProp->currIdx;
        uint8_t alertstate = adata->alertState;

        /* Considered double but need to be read based on type in
         * g_prop[pidx].data_type
         * int size = get_sizeof(prop[dep_idx].data_type);
         *  */
        double value = *(double *)adata->sValue;

        usys_log_debug("Alert %d received for Property[%d], Name: %s ,"
                        " Value %lf %s.",
                        alertstate, pidx, prop[pidx].name,
                        value, prop[didx].units);

        /* Report alert */
        ret = json_serialize_alert_data(&json, obj->modUuid, obj->name,
                        obj->desc, prop[pidx].name, prop[pidx].dataType,
                        adata->sValue, prop[didx].units);
        if (ret != JSON_ENCODING_OK) {
            usys_log_error( "Web service alert callback function for "
                            "Device name %s Disc: %s Module UUID %s called "
                            "with %d alerts failed to serialize alert",
                            obj->name, obj->desc, obj->modUuid, *count);
            goto cleanup;
        } else {
            ulfius_init_request(&alertNotification);
            ulfius_init_response(&alertNotificationResp);
            ulfius_set_request_properties(&alertNotification,
                            U_OPT_HTTP_VERB, "PUT",
                            U_OPT_HTTP_URL, gNotifServer,
                            U_OPT_HTTP_URL_APPEND, "alert",
                            U_OPT_TIMEOUT, 20,
                            U_OPT_JSON_BODY, json,
                            U_OPT_NONE);
            ret = ulfius_send_http_request(&alertNotification,
                            &alertNotificationResp);
            if (ret != STATUS_OK) {
                usys_log_error( "Web service alert callback function not able "
                                "to notify notification service.");
            } else {
                usys_log_debug( "Web service alert callback function"
                                " notification  response is %d.",
                                alertNotificationResp.status);
            }
            json_decref(json);
            ulfius_clean_request(&alertNotification);
        }

        usys_free(prop);

        cleanup:
        usys_free(adata->sValue);
        usys_free(adata);

    } else {

        usys_log_error( "Web service alert callback function for Device name %s"
                        " Disc: %s Module UUID %s called with %d alerts but with"
                        "any alert details.",
                        obj->name, obj->desc, obj->modUuid, *count);

    }
}

/**
 * @fn      void report_failure_with_response_code(UResponse*, int, int,
 *           char*)
 * @brief   Reports the failure to client using json object with HTTP repsonse
 *          provided in respcode.
 *
 * @param   response
 * @param   responsecode
 * @param   ret
 * @param   msg
 */
static void report_failure_with_response_code(UResponse *response, int respcode,
                int ret, char *msg) {
    JsonObj *json = NULL;
    ret = json_serialize_error(&json, ret, msg);
    if (ret != JSON_ENCODING_OK) {
        ulfius_set_empty_body_response(response, respcode);
    }
    ulfius_set_json_body_response(response, respcode, json);
}

/**
 * @fn      void report_failure(UResponse*, int, char*)
 * @brief   Reports a generic failure to the client using JSON
 *          with HTTP response code 500.
 *
 * @param response
 * @param ret
 * @param msg
 */
static void report_failure(UResponse *response, int ret, char *msg) {
    JsonObj *json = NULL;
    ret = json_serialize_error(&json, ret, msg);
    if (ret != JSON_ENCODING_OK) {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }
    ulfius_set_json_body_response(response, RESP_CODE_SERVER_FAILURE, json);
}

/**
 * @fn      void report_memory_failure(UResponse*, int)
 * @brief   Report memory related failure to client using json.
 *          with HTTP response code 500.
 *
 * @param   response
 * @param   errnum
 */
static void report_memory_failure(UResponse *response, int errnum) {
    JsonObj *json = NULL;
    int ret = json_serialize_error(&json, errnum, usys_error(errnum));
    if (ret != JSON_ENCODING_OK) {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }
    ulfius_set_json_body_response(response, RESP_CODE_SERVER_FAILURE, json);
}

/**
 * @fn      int web_service_cb_ping(const URequest*, UResponse*, void*)
 * @brief   reports ping response to client
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return
 */
static int web_service_cb_ping(const URequest *request, UResponse *response,
                void *epConfig) {
    int respCode = RESP_CODE_SUCCESS;

    ulfius_set_string_body_response(response, respCode,
                    "NodeD Service: Hi, there..!!");

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int web_service_cb_default(const URequest*, UResponse*, void*)
 * @brief   default callback used by REST framework if valid endpoint is not
 *          requested.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_default(const URequest *request, UResponse *response,
                void *epConfig) {
    int respCode = RESP_CODE_SUCCESS;

    char *msg;

    asprintf(&msg, "URL endpoint %s not implemented.", request->http_url);

    ulfius_set_string_body_response(response, respCode, msg);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int web_service_cb_discover_api(const URequest*, UResponse*, void*)
 * @brief   HTTP callback used by REST framework on discovery request.
 *          This list all the available endpoints from the service.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_discover_api(const URequest *request,
                UResponse *response, void *epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    uint16_t size = 0;
    UnitCfg *uCfg = NULL;
    usys_log_trace("NodeD:: Received a discover api request.");

    ret = json_serialize_api_list(&json, gApi, endPointCount);
    if (ret != JSON_ENCODING_OK) {
        report_failure(response, ret, "Failed serializing endpoints.");
        goto completed;
    }

    /* Send response */
    if (json) {
        ulfius_set_json_body_response(response, respCode, json);
    } else {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }

    completed:
    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int web_service_cb_get_unit_cfg(const URequest*, UResponse*, void*)
 * @brief   Callback function for reading unit config.It reads unit info and
 *          creates a json body for response.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_get_unit_cfg(const URequest *request,
                UResponse *response, void *epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SERVER_FAILURE;
    int ret = STATUS_NOK;
    uint16_t size = 0;
    UnitCfg *uCfg = NULL;
    usys_log_trace("NodeD:: Received a get unit config request.");

    /* Allocate memory for unit info */
    UnitInfo *uInfo = usys_zmalloc(sizeof(UnitInfo));
    if (!uInfo) {
        usys_log_error("Web Service Failed to allocate memory. Error %s",
                        usys_error(errno));
        report_memory_failure(response, errno);
        goto completed;
    }

    /* Read Unit info */
    ret = invt_read_unit_info("", uInfo, &size);
    if (ret) {
        usys_log_error("Web Service Failed to read unit info prior to config."
                        " Error Code %d",
                        ret);
        report_failure(response, ret,
                        "Failed while fetching unit info prior to config.");
        goto completed;
    }

    /* Allocate memory for unit config */
    uCfg = invt_alloc_unit_cfg(uInfo->modCount);
    if (!uCfg) {
        usys_log_error("Web Service Failed to allocate memory. Error %s",
                        usys_error(errno));
        report_memory_failure(response, errno);
        goto completed;
    }

    /* Read unit config */
    ret = invt_read_unit_cfg("", uCfg, uInfo->modCount, &size);
    if (!ret) {

        /* serialize unit config */
        ret = json_serialize_unit_cfg(&json, uCfg, uInfo->modCount);
        if (ret != JSON_ENCODING_OK) {
            report_failure(response, ret,
                            "Failed serializing unit config.");
            goto completed;
        }

    } else {
        report_failure(response, ret,
                        "Failed while fetching unit config.");
        goto completed;
    }

    /* Send response */
    if (json) {
        ulfius_set_json_body_response(response, respCode, json);
    } else {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }

    completed:
    /* Free memory */
    if (uInfo) {
        invt_free_unit_cfg(uCfg, uInfo->modCount);
        usys_free(uInfo);
        uInfo = NULL;
    }
    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int web_service_cb_get_unit_info(const URequest*, UResponse*, void*)
 * @brief   Callback function for reading unit info.It reads unit info and
 *          creates a json body for response.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE which is 0
 */
static int web_service_cb_get_unit_info(const URequest *request,
                UResponse *response, void *epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SERVER_FAILURE;
    int ret = STATUS_NOK;
    uint16_t size = 0;
    usys_log_trace("NodeD:: Received a get unit info request.");

    /* Allocate memory */
    UnitInfo *uInfo = usys_zmalloc(sizeof(UnitInfo));
    if (!uInfo) {
        usys_log_error("Web Service Failed to allocate memory. Error %s",
                        usys_error(errno));
        report_memory_failure(response, errno);
        goto completed;
    }

    /* Reads unit info */
    ret = invt_read_unit_info("", uInfo, &size);
    if (!ret) {
        ret = json_serialize_unit_info(&json, uInfo);
        /* if every thing id ok set code to success */
        if (ret != JSON_ENCODING_OK) {
            report_failure(response, ret, "Failed serializing unit info.");
            goto completed;
        }
    } else {
        usys_log_error("Web Service Failed to read unit info."
                        " Error Code %d",
                        ret);
        report_failure(response, ret, "Failed while fetching unit info.");
        goto completed;
    }

    /* Send response */
    if (json) {
        ulfius_set_json_body_response(response, respCode, json);
    } else {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }

    completed:
    /* Free memory */
    if (uInfo) {
        usys_free(uInfo);
        uInfo = NULL;
    }

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int web_service_cb_get_module_cfg(const URequest*, UResponse*, void*)
 * @brief    HTTP callback used by REST framework on module config request.
 *          This function reads the module config from the inventory db.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_get_module_cfg(const URequest *request,
                UResponse *response, void *epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    uint16_t size = 0;
    ModuleCfg *mCfg = NULL;
    ModuleInfo *mInfo = NULL;

    char *moduleId = (char*)u_map_get(request->map_url, UUID);
    if (!moduleId) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "no module UUID present");
        goto completed;
    }
    usys_log_trace("NodeD:: Received a get module config request for UUID %s.",
                    moduleId);

    /* Allocate memory for Module Info */
    mInfo = usys_zmalloc(sizeof(ModuleInfo));
    if (!mInfo) {
        usys_log_error("Web Service Failed to allocate memory for module info."
                        "Error %s",
                        usys_error(errno));
        report_memory_failure(response, errno);
        goto completed;
    }

    /* read module info */
    ret = invt_read_module_info(moduleId, mInfo, &size);
    if (ret) {
        /* Module info failure */
        usys_log_error("Web Service Failed to read module info for %s."
                        " Error Code %d",
                        moduleId, ret);
        report_failure(response, ret,
                        "Failed while fetching module info prior to"
                        " module config.");
        goto completed;
    }

    /* Allocate memory for module config */
    mCfg = invt_alloc_module_cfg(mInfo->devCount);
    if (mCfg) {
        /*  Read Module config */
        ret = invt_read_module_cfg(moduleId, mCfg, mInfo->devCount, &size);
        if (!ret) {
            ret = json_serialize_module_cfg(&json, mCfg, mInfo->devCount);
            if (ret != JSON_ENCODING_OK) {
                report_failure(response, ret,
                                "Failed serializing module config.");
                goto completed;
            }
        } else {
            usys_log_error("Web Service Failed to read module config."
                            " Error Code %d",
                            ret);
            report_failure(response, ret,
                            "Failed while fetching module config");
            goto completed;
        }

    } else {
        usys_log_error("Web Service Failed to allocate memory for module"
                        "config. Error %s",
                        usys_error(errno));
        report_memory_failure(response, errno);
        goto completed;
    }

    /* Send response */
    if (json) {
        ulfius_set_json_body_response(response, respCode, json);
    } else {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }

    completed:
    /* Free memory */
    if (mInfo) {
        invt_free_module_cfg(mCfg, mInfo->devCount);
        mCfg = NULL;
        usys_free(mInfo);
        mInfo = NULL;
    }

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int web_service_cb_get_module_info(const URequest*, UResponse*,
 *          void*)
 * @brief   HTTP callback used by REST framework on module info request.
 *          This function reads the module information from the inventory db.
 *
 * @param request
 * @param response
 * @param epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_get_module_info(const URequest *request,
                UResponse *response, void *epConfig) {
    JsonObj *json = NULL;
    uint16_t size = 0;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;

    char *moduleId = (char*) u_map_get(request->map_url, UUID);
    if (!moduleId) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "no module UUID present");
        goto completed;
    }
    usys_log_trace("NodeD:: Received a get module info request for UUID %s.",
                    moduleId);

    /* Allocate memory for Module Info */
    ModuleInfo *mInfo = usys_zmalloc(sizeof(ModuleInfo));
    if (!mInfo) {
        usys_log_error("Web Service Failed to allocate memory for module info."
                        "Error %s",
                        usys_error(errno));
        report_memory_failure(response, errno);
        goto completed;
    }

    /* Read Module info */
    ret = invt_read_module_info(moduleId, mInfo, &size);
    if (!ret) {
        ret = json_serialize_module_info(&json, mInfo);
        if (ret != JSON_ENCODING_OK) {
            report_failure(response, ret, "Failed serializing module info.");
            goto completed;
        }
    } else {
        usys_log_error("Web Service Failed to read module info for %s."
                        " Error Code %d",
                        moduleId, ret);
        report_failure(response, ret, "Failed while fetching module info");
        goto completed;
    }

    /* Send response */
    if (json) {
        ulfius_set_json_body_response(response, respCode, json);
    } else {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }

    completed:
    /* Free memory */
    if (mInfo) {
        usys_free(mInfo);
        mInfo = NULL;
    }

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      DevObj prepare_object_for_request*(UResponse*, const char*,
 *          const char*, const char*, int*, const char*, const char*, void**,
 *          int*, bool)
 * @brief   This function creates a device object we need to access and set
 *          the property index to appropriate value of property.
 *
 * @param   response
 * @param   devName
 * @param   devDesc
 * @param   moduleId
 * @param   propId
 * @param   devType
 * @param   propName
 * @param   dataMem
 * @param   dataType
 * @param   isGetReq
 * @return  On Success, return the DeviceObject we are trying to access.
 *          On failure, NULL.
 */
DevObj *prepare_object_for_request(UResponse *response, const char *devName,
                const char *devDesc, const char *moduleId,
                int *propId, const char *devType,
                const char *propName, void **dataMem,
                int *dataType, bool isGetReq) {
    Property *prop = NULL;
    /*  Identify device */
    DevObj *obj = usys_zmalloc(sizeof(DevObj));
    if (obj) {
        usys_strcpy(obj->name, devName);
        usys_strcpy(obj->desc, devDesc);
        usys_strcpy(obj->modUuid, moduleId);
        obj->type = atoi(devType);
    } else {
        return obj;
    }

    /* Identify the property we have to work on */
    uint16_t pCount = 0;
    int pIdx = -1;
    void *data = NULL;

    /* read property count */
    ldgr_read_prop_count(obj, &pCount);
    if (pCount == 0) {
        usys_free(obj);
        obj = NULL;
        return obj;
    }

    prop = usys_zmalloc(sizeof(Property) * pCount);
    if (prop) {
        int ret = ldgr_read_prop(obj, prop);
        if (ret) {
            usys_free(prop);
            usys_free(obj);
            return NULL;
        }
        for (int iter = 0; iter < pCount; iter++) {
            if (!usys_strcasecmp(prop[iter].name, propName)) {
                pIdx = iter;
                *dataType = prop[iter].dataType;
                break;
            }
        }

        if (pIdx == -1) {
            report_failure_with_response_code(response,
                            RESP_CODE_RESOURCE_NOT_FOUND,
                            RESP_CODE_RESOURCE_NOT_FOUND,
                            "no such property found.");
            usys_free(obj);
            usys_free(prop);
            obj = NULL;
            return obj;
        }
    }

    if (isGetReq) {
        /* Allocate memory based on type of property data
         * This is required in case of get request Only*/
        uint8_t dataSize = get_sizeof(*dataType);
        data = usys_zmalloc(dataSize);
        if (!data) {
            report_memory_failure(response, errno);
            usys_free(obj);
            usys_free(prop);
            obj = NULL;
            return obj;
        }
        *dataMem = data;
    }

    *propId = pIdx;
    usys_free(prop);
    return obj;
}

/**
 * @fn      int web_service_cb_put_dev_property(const URequest*, UResponse*,
 *           void*)
 * @brief   Writes the data provided in body to the appropriate  property of
 *          the device requested.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_put_dev_property(const URequest *request,
                UResponse *response,
                void *epConfig) {
    unsigned int respCode = RESP_CODE_SERVER_FAILURE;
    int ret = STATUS_NOK;
    void *data = NULL;
    int pIdx = 0;
    int dataType = 0;
    DevObj *obj = NULL;
    usys_log_trace("NodeD:: Received a read request to device property.");

    JsonObj *json = ulfius_get_json_body_request(request, NULL);

    char *moduleId = (char*) u_map_get(request->map_url, UUID);
    if (!moduleId) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "no module UUID present");
        goto completed;
    }
    usys_log_trace("NodeD:: Received a get module info request for UUID %s.",
                    moduleId);

    const char *devType = u_map_get(request->map_url, DEVTYPE);
    const char *devName = u_map_get(request->map_url, DEVNAME);
    const char *devDesc = u_map_get(request->map_url, DEVDESC);
    const char *propName = u_map_get(request->map_url, PROPNAME);

    usys_log_trace("NodeD:: Received a get device property data request "
                    "for UUID %s .",
                    moduleId);

    if (!((devType) && (devName) && (devDesc) && (propName))) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "missing info in request");
        goto completed;
    }

    /* Deserialize data */
    ret = json_deserialize_sensor_data(json, &devName, &devDesc, &dataType,
                    &data);
    if (ret != JSON_DECODING_OK) {
        report_failure(response, ret, "failed to decode json request");
        goto completed;
    }

    obj = prepare_object_for_request(response, devName, devDesc,
                    moduleId, &pIdx,
                    devType, propName, &data, &dataType, false);
    if (!obj) {
        report_failure(response, ret,
                        "failed to prepare read request to ledger.");
        goto completed;
    }

    /* Read data */
    ret = ldgr_write(obj, &pIdx, data);
    if (ret) {
        report_failure(response, ret, "failed to update device property.");
        goto completed;
    } else {
        respCode = RESP_CODE_SUCCESS;
    }

    /* Send response */
    ulfius_set_empty_body_response(response, respCode);

    completed:
    usys_free(obj);
    usys_free(data);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int web_service_cb_get_dev_property(const URequest*, UResponse*, void*)
 * @brief   Reads the appropriate property of the sensor requested by client.
 *          Info related to the sensor is provided in the URL parameters.
 *
 * @param request
 * @param response
 * @param epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_get_dev_property(const URequest *request,
                UResponse *response,
                void *epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    void *data = NULL;
    int pIdx = 0;
    int dataType = 0;
    DevObj *obj = NULL;
    char *moduleId = (char*) u_map_get(request->map_url, UUID);
    if (!moduleId) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "no module UUID present");
        goto completed;
    }
    usys_log_trace("NodeD:: Received a get device property read request"
                    " for UUID %s.",
                    moduleId);

    const char *devType = u_map_get(request->map_url, DEVTYPE);
    const char *devName = u_map_get(request->map_url, DEVNAME);
    const char *devDesc = u_map_get(request->map_url, DEVDESC);
    const char *propName = u_map_get(request->map_url, PROPNAME);

    usys_log_trace("NodeD:: Received a get module manufacturing data request "
                    "for UUID %s .",
                    moduleId);

    if (!((devType) && (devName) && (devDesc) && (propName))) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "missing info in request");
        goto completed;
    }

    obj = prepare_object_for_request(response, devName, devDesc, moduleId,
                    &pIdx,devType, propName, &data, &dataType, true);
    if (!obj) {
        report_failure(response, ret,
                        "failed to prepare read request to ledger.");
        goto completed;
    }

    /* Read data */
    ret = ldgr_read(obj, &pIdx, data);
    if (ret) {
        report_failure(response, ret, "failed to read device property.");
        goto completed;
    } else {
        ret = json_serialize_sensor_data(&json, devName, devDesc, dataType,
                        data);
        if (ret != JSON_ENCODING_OK) {
            report_failure(response, ret,
                            "failed to serialize device property response.");
            goto completed;
        }
    }

    /* Send response */
    if (json) {
        ulfius_set_json_body_response(response, respCode, json);
    } else {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }

    completed:
    usys_free(obj);
    usys_free(data);

    return U_CALLBACK_CONTINUE;
}


/**
 * @fn      int web_service_enable_alert(UResponse*, int, DevObj*)
 * @brief   Register callback and enable alerts on device property
 *
 * @param   response
 * @param   pIdx
 * @param   obj
 * @return  On success 0
 *          On failure -1
 */
static int  web_service_enable_alert(UResponse *response, int pIdx,
                DevObj* obj) {
    int ret = 0;

    /* Register callback */
    ret = ldgr_reg_app_cb(obj, &pIdx, &web_service_alert_cb);
    if (ret) {
        report_failure(response, ret, "failed to register alert callback.");
        ret = -1;
    }

    /* Enable alert */
    ret = ldgr_enable_irq(obj, &pIdx, NULL);
    if (ret) {
        report_failure(response, ret, "failed to enable alert on "
                        "device property.");
        ret = -1;
    }

    return ret;
}

/**
 * @fn      int web_service_disable_alert(UResponse*, int, DevObj*)
 * @brief   de-register callback and enable alerts on device property
 *
 * @param   response
 * @param   pIdx
 * @param   obj
 * @return  On success 0
 *          On failure -1
 */
static int web_service_disable_alert(UResponse *response, int pIdx, DevObj* obj) {
    int ret = 0;

    /* De-register callback */
    ret = ldgr_dereg_app_cb(obj, &pIdx, &web_service_alert_cb);
    if (ret) {
        report_failure(response, ret, "failed to register alert callback.");
        ret = -1;
    }

    /* Disable alert */
    ret = ldgr_disable_irq(obj, &pIdx, NULL);
    if (ret) {
        report_failure(response, ret, "failed to enable alert on "
                        "device property.");
        ret = -1;
    }

    return ret;
}

/**
 * @fn      int web_service_cb_put_dev_alert_state(const URequest*,
 *          UResponse*, void*)
 * @brief   Enable/Disable alert for the device for the property mentioned
 *          in query.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_put_dev_alert_state(const URequest *request,
                UResponse *response,
                void *epConfig) {
    unsigned int respCode = RESP_CODE_SERVER_FAILURE;
    int ret = STATUS_NOK;
    void *data = NULL;
    int pIdx = 0;
    int dataType = 0;
    DevObj *obj = NULL;
    usys_log_trace("NodeD:: Received a enable/disable alert request.");

    JsonObj *json = ulfius_get_json_body_request(request, NULL);

    char *moduleId = (char*) u_map_get(request->map_url, UUID);
    if (!moduleId) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "no module UUID present");
        goto completed;
    }
    usys_log_trace("NodeD:: Received a enable/disable alert request for UUID %s.",
                    moduleId);

    const char *devType = u_map_get(request->map_url, DEVTYPE);
    const char *devName = u_map_get(request->map_url, DEVNAME);
    const char *devDesc = u_map_get(request->map_url, DEVDESC);
    const char *propName = u_map_get(request->map_url, PROPNAME);
    const char *state = u_map_get(request->map_url, ALERTSTATE);

    usys_log_trace("NodeD:: Received a get device property data request "
                    "for UUID %s .",
                    moduleId);

    if (!((devType) && (devName) && (devDesc) && (propName) && (state))) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "missing info in request");
        goto completed;
    }

//    /* Deserialize data */
//    ret = json_deserialize_sensor_data(json, &devName, &devDesc, &dataType,
//                    &data);
//    if (ret != JSON_DECODING_OK) {
//        report_failure(response, ret, "failed to decode json request");
//        goto completed;
//    }

    obj = prepare_object_for_request(response, devName, devDesc,
                    moduleId, &pIdx,
                    devType, propName, &data, &dataType, false);
    if (!obj) {
        report_failure(response, ret,
                        "failed to prepare read request to ledger.");
        goto completed;
    }

    if (usys_strcasecmp(state, ENABLE)==0) {
        ret = web_service_enable_alert(response, pIdx, obj);
    } else {
        ret = web_service_disable_alert(response, pIdx, obj);
    }

    if (ret) {
        goto completed;
    } else {
        respCode = RESP_CODE_SUCCESS;
    }

    /* Send response */
    ulfius_set_empty_body_response(response, respCode);

    completed:
    usys_free(obj);
    usys_free(data);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int field_name_to_id(const char*, uint16_t*)
 * @brief   Translates the URL parameter into the field Id
 *          for inventory database.
 *
 * @param   fieldName
 * @param   fieldId
 * @return  On success STATUS_OK
 *          On failure STATUS_NOK
 */
int field_name_to_id(const char *fieldName, uint16_t *fieldId) {
    int ret = STATUS_NOK;

    if (!usys_strcmp("factorycalibration", fieldName)) {
        *fieldId = FIELD_ID_FACT_CALIB;
        ret = STATUS_OK;
    }

    if (!usys_strcmp("usercalibration", fieldName)) {
        *fieldId = FIELD_ID_USER_CALIB;
        ret = STATUS_OK;
    }

    if (!usys_strcmp("factoryconfig", fieldName)) {
        *fieldId = FIELD_ID_FACT_CFG;
        ret = STATUS_OK;
    }

    if (!usys_strcmp("userconfig", fieldName)) {
        *fieldId = FIELD_ID_USER_CFG;
        ret = STATUS_OK;
    }

    if (!usys_strcmp("bootstrapcerts", fieldName)) {
        *fieldId = FIELD_ID_BS_CERTS;
        ret = STATUS_OK;
    }

    if (!usys_strcmp("cloudcerts", fieldName)) {
        *fieldId = FIELD_ID_CLOUD_CERTS;
        ret = STATUS_OK;
    }

    return ret;
}

/**
 * @fn      int web_service_cb_get_module_mfg(const URequest*, UResponse*, void*)
 * @brief   Read manufacturing data for the module.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_get_module_mfg(const URequest *request,
                UResponse *response, void *epConfig) {
    JsonObj *json = NULL;
    int ret = STATUS_NOK;
    uint16_t size = 0;
    uint16_t fieldId = 0;
    char *data = NULL;
    char *moduleId = (char*) u_map_get(request->map_url, UUID);
    if (!moduleId) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "no module UUID present");
        goto completed;
    }
    usys_log_trace("NodeD:: Received a get module info request for UUID %s.",
                    moduleId);

    const char *fieldName = u_map_get(request->map_url, MFGDATA);
    if (!fieldName) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "no mfg data name present");
        goto completed;
    }
    usys_log_trace("NodeD:: Manufacturing data info request for %s.",
                    fieldName);

    ret = field_name_to_id(fieldName, &fieldId);
    if (ret) {
        report_failure(response, ret,
                        "data name provided is not matching "
                        "to any field.");
        goto completed;
    }

    /* Read data from request */
    ret = invt_read_payload_for_field_id(moduleId, (void **)&data, fieldId,
                    &size);
    if (ret) {
        report_failure(response, ret, "failed to read data.");
        goto completed;
    } else {
        ulfius_set_binary_body_response(response, RESP_CODE_SUCCESS, data,
                        size);
    }

    completed:
    usys_free(data);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int web_service_cb_put_module_mfg(const URequest*, UResponse*, void*)
 * @brief   Update module manufacturing data.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_put_module_mfg(const URequest *request,
                UResponse *response, void *epConfig) {
    int ret = STATUS_NOK;
    int size = 0;
    char *data = NULL;
    uint16_t fieldId = 0;
    char *moduleId = (char*) u_map_get(request->map_url, UUID);
    if (!moduleId) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "no module UUID present");
        goto completed;
    }
    usys_log_trace("NodeD:: Received a get module info request for UUID %s.",
                    moduleId);

    const char *fieldName = u_map_get(request->map_url, MFGDATA);
    if (!fieldName) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        RESP_CODE_INVALID_REQUEST,
                        "no mfg data name present");
        goto completed;
    }
    usys_log_trace("NodeD:: Manufacturing data info request for %s.",
                    fieldName);

    ret = field_name_to_id(fieldName, &fieldId);
    if (ret) {
        report_failure(response, ret,
                        "data name provided is not matching "
                        "to any field.");
        goto completed;
    }

    /* Check the size limit */
    if (request->binary_body_length > SCH_MAX_PAYLOAD_SIZE) {
        report_failure_with_response_code(response, RESP_CODE_INVALID_REQUEST,
                        ERR_NODED_EXCEED_MAX_SIZE,
                        "Payload exceeds max payload length.");
        goto completed;
    }

    /* Write data from request */
    ret = invt_update_payload(moduleId, request->binary_body, fieldId,
                    request->binary_body_length);
    if (ret) {
        report_failure(response, ret, "failed to update data.");
        goto completed;
    } else {
        ulfius_set_empty_body_response(response, RESP_CODE_SUCCESS);
    }

    completed:
    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      void web_service_add_end_point(char*, char*, void*, HttpCb)
 * @brief   Wrapper function on adding endpoint to REST framework. This
 *          also populates struct for endPoint discovery which can then
 *          be used for listing endpoints in later stages.
 *
 * @param   method
 * @param   endPoint
 * @param   config
 * @param   cb
 */
static void web_service_add_end_point(char *method, char *endPoint,
                void *config, HttpCb cb) {
    ulfius_add_endpoint_by_val(&serverInst, method, endPoint, NULL, 0, cb,
                    config);
    usys_strcpy(gApi[endPointCount].method, method);
    usys_strcpy(gApi[endPointCount].endPoint, endPoint);
    usys_log_trace("Added api[%d] Method %s Endpoint: %s.", endPointCount,
                    "Get", endPoint);
    endPointCount++;
}

/**
 * @fn      void web_service_add_device_based_endpoint(int, void*, DevObj*)
 * @brief   Add REST endpoints for device config reading and writing.
 *
 * @param perm
 * @param config
 * @param devEp
 */
static void web_service_add_device_based_endpoint() {
    /* Write permissions */
    web_service_add_end_point("PUT", API_RES_EP("deviceconfig"), NULL,
                    web_service_cb_put_dev_property);

    /* Read permissions */
    web_service_add_end_point("GET", API_RES_EP("deviceconfig"), NULL,
                    web_service_cb_get_dev_property);

    /* Enable alerts */
    web_service_add_end_point("PUT", API_RES_EP("alertstate"), NULL,
                    web_service_cb_put_dev_alert_state);

}
/**
 * @fn      void web_service_add_discover_endpoints()
 * @brief   Add REST endpoints for endpoint discovery.
 *
 */
static void web_service_add_discover_endpoints() {
    web_service_add_end_point("GET", API_RES_EP("discover"), NULL,
                    web_service_cb_discover_api);
}

/**
 * @fn      void web_service_add_unit_endpoints()
 * @brief   Add REST endpoints for unit info reading.
 *
 */
void web_service_add_unit_endpoints() {
    web_service_add_end_point("GET", API_RES_EP("unitinfo"), NULL,
                    web_service_cb_get_unit_info);
    web_service_add_end_point("GET", API_RES_EP("unitconfig"), NULL,
                    web_service_cb_get_unit_cfg);
}

/**
 * @fn      void web_service_add_module_endpoints()
 * @brief   Add REST endpoints for module info reading.
 *
 */
void web_service_add_module_endpoints() {
    web_service_add_end_point("GET", API_RES_EP("moduleinfo"), NULL,
                    web_service_cb_get_module_info);
    web_service_add_end_point("GET", API_RES_EP("moduleconfig"), NULL,
                    web_service_cb_get_module_cfg);
}

/**
 * @fn      void web_service_add_mfg_data_endpoints()
 * @brief   Add REST endpoints for manufacturing data read and write.
 *
 */
void web_service_add_mfg_data_endpoints() {
    web_service_add_end_point("GET", API_RES_EP("mfg"), NULL,
                    web_service_cb_get_module_mfg);
    web_service_add_end_point("PUT", API_RES_EP("mfg"), NULL,
                    web_service_cb_put_module_mfg);
}

/**
 * @fn      void setup_web_service_endpoints(UInst*, void*)
 * @brief   Add default endpoint and endpont for ping
 *
 * @param   instance
 * @param   config
 */
static void setup_web_service_endpoints(UInst *instance, void *config) {

    /* Ping */
    web_service_add_end_point("GET", API_RES_EP("ping"), NULL,
                    web_service_cb_ping);

    /* default endpoint. */
    ulfius_set_default_endpoint(instance, &web_service_cb_default, NULL);
}

/**
 * @fn      int start_framework(UInst*)
 * @brief   Initializes the REST server framework
 *
 * @param   instance
 * @return  On success STATUS_OK
 *          On failure STATUS_NOK
 */
static int start_framework(UInst *instance) {
    int ret;

    ret = ulfius_start_framework(instance);
    if (ret != U_OK) {
        usys_log_error("Error starting the web_service.");

        /* clean up. */
        ulfius_stop_framework(instance); /* don't think need this. XXX */
        ulfius_clean_instance(instance);

        return STATUS_NOK;
    }

    return STATUS_OK;
}

/**
 * @fn      int init_framework(UInst*, int)
 * @brief   Initializes the REST server framework
 *
 * @param   inst
 * @param   port
 * @return  On success STATUS_OK
 *          On failure STATUS_NOK
 */
static int init_framework(UInst *inst, int port) {
    if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
        usys_log_error(
                        "Error initializing instance for websocket"
                        " remote port %d", port);
        return STATUS_NOK;
    }

    /* Set few params. */
    u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

    return STATUS_OK;
}

/**
 * @fn      int web_service_start()
 * @brief   Add API endpoints and start the REST HTTP server
 *
 * @return  On success STATUS_OK
 *          On failure STATUS_NOK
 */
int web_service_start() {
    /* setup endpoints and methods callback. */
    setup_web_service_endpoints(&serverInst, NULL);

    web_service_add_unit_endpoints();

    web_service_add_module_endpoints();

    web_service_add_device_based_endpoint();

    web_service_add_mfg_data_endpoints();

    web_service_add_discover_endpoints();

    /* open connection for web_services */
    if (start_framework(&serverInst)) {
        usys_log_error("Failed to start web_services for cld_ctrl: %d",
                        WEB_SERVICE_PORT);
        return STATUS_NOK;
    }

    usys_log_info("Webservice on client port: %d started.", WEB_SERVICE_PORT);

    return STATUS_OK;
}

/**
 * @fn      int web_service_start()
 * @brief   Initializes the ulfius framework for REST server.
 *
 * @return  On success STATUS_OK
 *          On failure STATUS_NOK
 */
int web_service_init(char *notifServer) {
    /* Notification server */
    usys_strcpy(gNotifServer, notifServer);
    usys_log_info("Added notification server  %s", gNotifServer);

    /* Initialize the web_services framework. */
    if (init_framework(&serverInst, WEB_SERVICE_PORT) != STATUS_OK) {
        usys_log_error("Error initializing web_service framework");
        return STATUS_NOK;
    }
    return STATUS_OK;
}
