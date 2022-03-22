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


static int web_service_cb_ping(const URequest * request, UResponse * response, void * epConfig) {

    int respCode = RESP_CODE_SUCCESS;

    ulfius_set_string_body_response(response, respCode, "NodeD Service: Hi, there..!!");

    return U_CALLBACK_CONTINUE;
}

static int web_service_cb_default(const URequest * request, UResponse * response, void * epConfig) {

    int respCode = RESP_CODE_SUCCESS;

    char *msg;

    asprintf(&msg,"URL endpoint %s not implemented.", request->http_url);

    ulfius_set_string_body_response(response, respCode, msg);

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
#if 0
static char* create_endpoint(DevObj *ep) {
    char *epoint =  usys_zmalloc(usys_strlen(API_RES_EP)
                    + usys_strlen(ep->modUuid)
                    + usys_strlen(usys_itoa(ep->type))
                    + usys_strlen(ep->name)
                    + usys_strlen(ep->desc) + 5);
    if (epoint) {
        usys_memcpy(epoint, API_RES_EP, strlen(API_RES_EP));
        usys_strcat(epoint, ep->modUuid);
        usys_strcat(epoint, EP_PS);
        usys_strcat(epoint, usys_itoa(ep->type));
        usys_strcat(epoint, EP_PS);
        usys_strcat(epoint, ep->name);
        usys_strcat(epoint, EP_PS);
        usys_strcat(epoint, ep->desc);
    }

    return epoint;
}

void char* create_new_end_point(int num, ...) {
    va_list valist;
    char *endPoint = NULL;
    char *str = NULL;
    int len = usys_strlen(usys_strlen(API_RES_EP));

    char *endPoint = usys_zmalloc(len);
    if (endPoint) {
        /* initialize valist for num number of arguments */
        va_start(valist, num);


        /* access all the arguments assigned to valist */
        for (int i = 0; i < num; i++) {
            str = va_arg(valist, char*);
            len += usys_strlen(str);
            endPoint = usys_realloc(endPoint, len+2);
            usys_strcat(endPoint, EP_PS"/0");
            usys_strcat(endPoint, str);
            endPoint[len] = '\0';
        }

        /* clean memory reserved for valist */
        va_end(valist);
    }
    return endPoint;
}
#endif

static void web_service_add_device_based_endpoint(int perm, void *config, DevObj* devEp) {
    char *endPoint = API_RES_EP(":UUID/:devType/:devName/:devDesc");
    //    char* endPoint[] = API_RES_EP(":UUID/:devType/:devName/:devDesc")
    //    usys_strcat(epoint, ep->modUuid);
    //    usys_strcat(epoint, EP_PS);
    //    usys_strcat(epoint, usys_itoa(ep->type));
    //    usys_strcat(epoint, EP_PS);
    //    usys_strcat(epoint, ep->name);
    //    usys_strcat(epoint, EP_PS);
    //    usys_strcat(epoint, ep->desc);
    //    if (endPoint) {

    /* Write permissions */
    if (perm & PERM_WR) {
        ulfius_add_endpoint_by_val(&serverInst, "PUT", endPoint, config, 0, &web_service_cb_put_dev_property, config);
        usys_log_trace("Added Method %s Endpoint: %s." "Put", endPoint);

    }

    /* Read permissions */
    if (perm & PERM_RD) {
        ulfius_add_endpoint_by_val(&serverInst, "GET", endPoint, config, 0, &web_service_cb_get_dev_property, config);
        usys_log_trace("Added Method %s Endpoint: %s." "Get", endPoint);
    }

    //        usys_free(endPoint);
    //    }
}

static void web_service_add_unit_info_endpoint() {
    char *endPoint = API_RES_EP("unitinfo");
    ulfius_add_endpoint_by_val(&serverInst, "GET", endPoint, NULL, 0, &web_service_cb_get_unit_info, NULL);
    usys_log_trace("Added Method %s Endpoint: %s.", "Get", endPoint);
}

static void web_service_add_unit_cfg_endpoint() {
    char *endPoint = API_RES_EP("unitconfig");
    ulfius_add_endpoint_by_val(&serverInst, "GET", endPoint, NULL, 0, &web_service_cb_get_unit_cfg, NULL);
    usys_log_trace("Added Method %s Endpoint: %s.", "Get", endPoint);
}

void web_service_add_unit_endpoints() {
    web_service_add_unit_info_endpoint();
    web_service_add_unit_cfg_endpoint();
}

static void web_service_add_module_info_endpoint() {
    char *endPoint = API_RES_EP(":UUID/moduleinfo");
    ulfius_add_endpoint_by_val(&serverInst, "GET", endPoint, NULL, 0, &web_service_cb_get_module_info, NULL);
    usys_log_trace("Added Method %s Endpoint: %s.", "Get", endPoint);

}

static void web_service_add_module_cfg_endpoint() {
    char *endPoint = API_RES_EP(":UUID/moduleconfig");
    ulfius_add_endpoint_by_val(&serverInst, "GET", endPoint, NULL, 0, &web_service_cb_get_module_cfg, NULL);
    usys_log_trace("Added Method %s Endpoint: %s.", "Get", endPoint);

}

void web_service_add_module_endpoints() {
    web_service_add_module_info_endpoint();
    web_service_add_module_cfg_endpoint();
}

static void web_service_add_module_mfg_endpoint(char* moduleId, char* fieldName, uint8_t* fieldId, void *config) {
    //char *endPoint = create_new_endpoint(moduleId, fieldName);
    char *endPoint = API_RES_EP(":UUID/:fieldName");
    ulfius_add_endpoint_by_val(&serverInst, "GET", fieldId, NULL, 0, &web_service_cb_get_module_mfg, NULL);
    usys_log_trace("Added Method %s Endpoint: %s." "Get", endPoint);

    ulfius_add_endpoint_by_val(&serverInst, "PUT", fieldId, NULL, 0, &web_service_cb_put_module_mfg, NULL);
    usys_log_trace("Added Method %s Endpoint: %s." "PUT", endPoint);
    usys_free(endPoint);
}

//static void web_service_add_module_based_endpoint(int perm, void *config, DevObj* devEp) {
//    char *endPoint = create_endpoint(devEp);
//    if (endPoint) {
//
//        /* Write permissions */
//        if (perm & PERM_WR) {
//            ulfius_add_endpoint_by_val(&serverInst, "PUT", endPoint, config, 0, &web_service_cb_put_dev_property, config);
//            usys_log_trace("Added Method %s Endpoint: %s." "Put", endPoint);
//
//        }
//
//        /* Read permissions */
//        if (perm & PERM_RD) {
//            ulfius_add_endpoint_by_val(serverInst, "GET", config, NULL, 0, &web_service_cb_get_dev_property, config);
//            usys_log_trace("Added Method %s Endpoint: %s." "Get", endPoint);
//        }
//
//        usys_free(endPoint);
//    }
//}

/*
 * setup_web_service_endpoints --
 *
 */
static void setup_web_service_endpoints(UInst *instance, void * config) {

    /* Endpoint list declaration. */
    ulfius_add_endpoint_by_val(instance, "GET", API_RES_EP("ping"), NULL, 0,
            &web_service_cb_ping, config);

    /* default endpoint. */
    ulfius_set_default_endpoint(instance, &web_service_cb_default, NULL);

}

/*
 * start_framework --
 *
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

static int init_framework(UInst *inst, int port) {

    if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
        usys_log_error("Error initializing instance for websocket remote port %d", port);
        return STATUS_NOK;
    }

    /* Set few params. */
    u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");


    return STATUS_OK;
}

/*
 * start_web_services -- start accepting REST clients on 127.0.0.1:port
 *
 */
int web_service_start() {

    /* setup endpoints and methods callback. */
    setup_web_service_endpoints(&serverInst, NULL);

    web_service_add_unit_endpoints();

    web_service_add_module_endpoints();

    /* open connection for both admin and client web_services */
    if (start_framework(&serverInst)) {
        usys_log_error( "Failed to start web_services for cld_ctrl: %d",
                WEB_SERVICE_PORT);
        return STATUS_NOK;
    }

    usys_log_info( "Webservice on client port: %d started.", WEB_SERVICE_PORT);

    return STATUS_OK;
}

int web_service_init() {
    /* Initialize the admin and client web_services framework. */
    if (init_framework(&serverInst, WEB_SERVICE_PORT) != STATUS_OK){
        usys_log_error( "Error initializing web_service framework");
        return STATUS_NOK;
    }
    return STATUS_OK;
}

