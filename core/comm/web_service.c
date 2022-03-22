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
#include "inventory.h"
#include "jserdes.h"
#include "property.h"
#include "service.h"
#include "ulfius.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

UInst serverInst;

static uint16_t endPointCount = 0;
WebServiceAPI gApi[MAX_END_POINTS] = {0};

/**
 * @fn int web_service_cb_ping(const URequest*, UResponse*, void*)
 * @brief
 *
 * @param request
 * @param response
 * @param epConfig
 * @return
 */
static int web_service_cb_ping(const URequest * request, UResponse * response, void * epConfig) {

    int respCode = RESP_CODE_SUCCESS;

    ulfius_set_string_body_response(response, respCode, "NodeD Service: Hi, there..!!");

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn int web_service_cb_default(const URequest*, UResponse*, void*)
 * @brief
 *
 * @param request
 * @param response
 * @param epConfig
 * @return
 */
static int web_service_cb_default(const URequest * request, UResponse * response, void * epConfig) {

    int respCode = RESP_CODE_SUCCESS;

    char *msg;

    asprintf(&msg,"URL endpoint %s not implemented.", request->http_url);

    ulfius_set_string_body_response(response, respCode, msg);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn int web_service_cb_discover_api(const URequest*, UResponse*, void*)
 * @brief
 *
 * @param request
 * @param response
 * @param epConfig
 * @return
 */
static int web_service_cb_discover_api(const URequest * request, UResponse * response, void * epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    uint16_t size = 0;
    UnitCfg *uCfg = NULL;
    usys_log_trace("NodeD:: Received a discover api request.");

    ret = json_serialize_api_list(&json, gApi, endPointCount);
    if (ret != JSON_ENCODING_OK) {

        usys_log_error("Web Service Failed to serialize endpoint data."
                        " Error Code %d",
                        ret);
        ret = json_serialize_error(&json, ret, "Failed to create a json data.");
    }

    /* Send response */
    if (json) {
        ulfius_set_json_body_response(response, respCode, json);
    } else {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }


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
 * @return
 */
static int web_service_cb_get_unit_cfg(const URequest * request, UResponse * response, void * epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    uint16_t size = 0;
    UnitCfg *uCfg = NULL;
    usys_log_trace("NodeD:: Received a get unit config request.");

    /* Reads unit info */
    UnitInfo *uInfo = usys_zmalloc(sizeof(UnitInfo));
    if(uInfo) {

        ret = invt_read_unit_info("", uInfo, &size);
        if (!ret) {
            uCfg = invt_alloc_unit_cfg(uInfo->modCount);
            if(uCfg) {
                ret = invt_read_unit_cfg("", uCfg, uInfo->modCount, &size);
                if (!ret) {
                    ret = json_serialize_unit_cfg(&json, uCfg, uInfo->modCount);
                } else {
                    ret = json_serialize_error(&json, ret,
                                    "Failed to fetch unit config.");
                }
            }
        } else {
            usys_log_error("Web Service Failed to read unit info."
                            " Error Code %d",
                            ret);
            ret = json_serialize_error(&json, ret, "Failed to fetch unit info.");
        }

    } else {

        usys_log_error("Web Service Failed to allocate memory. Error %s",
                        usys_error(errno));
        ret = json_serialize_error(&json, errno, usys_error(errno));

    }

    /* Send response */
    if (json) {
        ulfius_set_json_body_response(response, respCode, json);
    } else {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }

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
static int web_service_cb_get_unit_info(const URequest * request, UResponse * response, void * epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    uint16_t size = 0;
    usys_log_trace("NodeD:: Received a get unit info request.");

    /* Reads unit info */
    UnitInfo *uInfo = usys_zmalloc(sizeof(UnitInfo));
    if(uInfo) {

        ret = invt_read_unit_info("", uInfo, &size);
        if (!ret) {
            ret = json_serialize_unit_info(&json, uInfo);
        } else {
            usys_log_error("Web Service Failed to read unit info."
                            " Error Code %d",
                            ret);
            ret = json_serialize_error(&json, ret, "Failed to fetch unit info.");
        }

    } else {

        usys_log_error("Web Service Failed to allocate memory. Error %s",
                        usys_error(errno));
        ret = json_serialize_error(&json, errno, usys_error(errno));

    }

    /* Send response */
    if (json) {
        ulfius_set_json_body_response(response, respCode, json);
    } else {
        ulfius_set_empty_body_response(response, RESP_CODE_SERVER_FAILURE);
    }

    /* Free memory */
    if (uInfo) {
        usys_free(uInfo);
        uInfo = NULL;
    }

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn int web_service_cb_get_module_cfg(const URequest*, UResponse*, void*)
 * @brief
 *
 * @param request
 * @param response
 * @param epConfig
 * @return
 */
static int web_service_cb_get_module_cfg(const URequest * request, UResponse * response, void * epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    const char *moduleId = u_map_get(request->map_url, UUID);
    usys_log_trace("NodeD:: Received a get module config request for UUID %s.", moduleId);

    /* Prepare response */
    json = json_object();
    if (json) {

    }
    ulfius_set_json_body_response(response, respCode, json);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn int web_service_cb_get_module_info(const URequest*, UResponse*, void*)
 * @brief
 *
 * @param request
 * @param response
 * @param epConfig
 * @return
 */
static int web_service_cb_get_module_info(const URequest * request, UResponse * response, void * epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    const char *moduleId = u_map_get(request->map_url, UUID);
    usys_log_trace("NodeD:: Received a get module info request for UUID %s.", moduleId);

    /* Prepare response */
    json = json_object();
    if (json) {

    }
    ulfius_set_json_body_response(response, respCode, json);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn int web_service_cb_put_dev_property(const URequest*, UResponse*, void*)
 * @brief
 *
 * @param request
 * @param response
 * @param epConfig
 * @return
 */
static int web_service_cb_put_dev_property(const URequest * request, UResponse * response, void * epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    usys_log_trace("NodeD:: Received a write request to device property.");

    /* Prepare response */
    json = json_object();
    if (json) {

    }
    ulfius_set_json_body_response(response, respCode, json);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn int web_service_cb_get_dev_property(const URequest*, UResponse*, void*)
 * @brief
 *
 * @param request
 * @param response
 * @param epConfig
 * @return
 */
static int web_service_cb_get_dev_property(const URequest * request, UResponse * response, void * epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    usys_log_trace("NodeD:: Received a read request to device property.");

    /* Prepare response */
    json = json_object();
    if (json) {

    }
    ulfius_set_json_body_response(response, respCode, json);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn int web_service_cb_put_module_mfg(const URequest*, UResponse*, void*)
 * @brief
 *
 * @param request
 * @param response
 * @param epConfig
 * @return
 */
static int web_service_cb_put_module_mfg(const URequest * request, UResponse * response, void * epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    usys_log_trace("NodeD:: Received a update request to manufacturing data.");

    /* Prepare response */
    json = json_object();
    if (json) {

    }
    ulfius_set_json_body_response(response, respCode, json);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn int web_service_cb_get_module_mfg(const URequest*, UResponse*, void*)
 * @brief
 *
 * @param request
 * @param response
 * @param epConfig
 * @return
 */
static int web_service_cb_get_module_mfg(const URequest * request, UResponse * response, void * epConfig) {
    JsonObj *json = NULL;
    unsigned int respCode = RESP_CODE_SUCCESS;
    int ret = STATUS_NOK;
    usys_log_trace("NodeD:: Received a read request to manufacturing data.");

    /* Prepare response */
    json = json_object();
    if (json) {

    }
    ulfius_set_json_body_response(response, respCode, json);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn void web_service_add_end_point(char*, char*, void*, HttpCb)
 * @brief
 *
 * @param method
 * @param endPoint
 * @param config
 * @param cb
 */
static void web_service_add_end_point(char* method, char* endPoint, void *config, HttpCb cb) {
    ulfius_add_endpoint_by_val(&serverInst, method, endPoint, NULL, 0, cb, config);
    usys_strcpy(gApi[endPointCount].method, method);
    usys_strcpy(gApi[endPointCount].endPoint, endPoint);
    usys_log_trace("Added api[%d] Method %s Endpoint: %s.", endPointCount, "Get", endPoint);
    endPointCount++;
}

/**
 * @fn void web_service_add_discover_endpoints()
 * @brief
 *
 */
static void web_service_add_discover_endpoints() {
    web_service_add_end_point("GET", API_RES_EP("discover"), NULL, web_service_cb_discover_api);
}

/**
 * @fn void web_service_add_unit_endpoints()
 * @brief
 *
 */
void web_service_add_unit_endpoints() {
    web_service_add_end_point("GET", API_RES_EP("unitinfo"), NULL, web_service_cb_get_unit_info);
    web_service_add_end_point("GET", API_RES_EP("unitconfig"), NULL, web_service_cb_get_unit_cfg);
}

/**
 * @fn void web_service_add_module_endpoints()
 * @brief
 *
 */
void web_service_add_module_endpoints() {
    web_service_add_end_point("GET", API_RES_EP("moduleinfo/:UUID"), NULL, web_service_cb_get_module_info);
    web_service_add_end_point("GET", API_RES_EP("moduleconfig/:UUID"), NULL, web_service_cb_get_module_cfg);
}

/**
 * @fn void setup_web_service_endpoints(UInst*, void*)
 * @brief
 *
 * @param instance
 * @param config
 */
static void setup_web_service_endpoints(UInst *instance, void * config) {

    /* Ping */
    web_service_add_end_point("GET", API_RES_EP("ping"), NULL, web_service_cb_ping);

    /* default endpoint. */
    ulfius_set_default_endpoint(instance, &web_service_cb_default, NULL);

}

/**
 * @fn int start_framework(UInst*)
 * @brief
 *
 * @param instance
 * @return
 */
static int start_framework(UInst *instance ) {

    int ret;

    ret = ulfius_start_framework(instance);
    if (ret != U_OK) {
        usys_log_error( "Error starting the web_service.");

        /* clean up. */
        ulfius_stop_framework(instance); /* don't think need this. XXX */
        ulfius_clean_instance(instance);

        return STATUS_NOK;
    }

    return STATUS_OK;
}

/**
 * @fn int init_framework(UInst*, int)
 * @brief
 *
 * @param inst
 * @param port
 * @return
 */
static int init_framework(UInst *inst, int port) {

    if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
        usys_log_error("Error initializing instance for websocket remote port %d", port);
        return STATUS_NOK;
    }

    /* Set few params. */
    u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");


    return STATUS_OK;
}

/**
 * @fn int web_service_start()
 * @brief
 *
 * @return
 */
int web_service_start() {

    /* setup endpoints and methods callback. */
    setup_web_service_endpoints(&serverInst, NULL);

    web_service_add_unit_endpoints();

    web_service_add_module_endpoints();

    web_service_add_discover_endpoints();

    /* open connection for both admin and client web_services */
    if (start_framework(&serverInst)) {
        usys_log_error( "Failed to start web_services for cld_ctrl: %d",
                        WEB_SERVICE_PORT);
        return STATUS_NOK;
    }

    usys_log_info( "Webservice on client port: %d started.", WEB_SERVICE_PORT);

    return STATUS_OK;
}

/**
 * @fn int web_service_init()
 * @brief
 *
 * @return
 */
int web_service_init() {
    /* Initialize the admin and client web_services framework. */
    if (init_framework(&serverInst, WEB_SERVICE_PORT) != STATUS_OK){
        usys_log_error( "Error initializing web_service framework");
        return STATUS_NOK;
    }
    return STATUS_OK;
}

