/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "web_service.h"

#include "jserdes.h"
#include "service.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

UInst serverInst;

static uint16_t endPointCount = 0;
WebServiceAPI gApi[MAX_END_POINTS] = { 0 };



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
                    "NotifyD Service: Hi, there..!!");

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

    usys_log_trace("NotifyD:: Received a discover api request.");

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
 * @fn      int web_service_cb_post_event(const URequest*, UResponse*, void*)
 * @brief   Receive a new event reported from service.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_post_event(const URequest *request,
                UResponse *response, void *epConfig) {
    int ret = STATUS_NOK;
    unsigned int respCode = RESP_CODE_SERVER_FAILURE;
    const char *service = u_map_get(request->map_url, "service");
    usys_log_trace("NotifyD:: Received a post request to event from %s.",
                    service);

    JsonObj *json = ulfius_get_json_body_request(request, NULL);

    ret = notify_process_incoming_notification(service, NOTIFICATION_EVENT, json);
    if (ret == STATUS_OK) {
        respCode = RESP_CODE_ACCEPTED;
    }

    /* Send response */
    ulfius_set_empty_body_response(response, respCode);

    return U_CALLBACK_CONTINUE;
}

/**
 * @fn      int web_service_cb_post_alert(const URequest*, UResponse*, void*)
 * @brief   Receive a new event reported from service.
 *
 * @param   request
 * @param   response
 * @param   epConfig
 * @return  U_CALLBACK_CONTINUE is returned to REST framework.
 */
static int web_service_cb_post_alert(const URequest *request,
                UResponse *response, void *epConfig) {
    int ret = STATUS_NOK;
    unsigned int respCode = RESP_CODE_SERVER_FAILURE;
    const char *service = u_map_get(request->map_url, "service");
    usys_log_trace("NotifyD:: Received a post request to alert from %s.",
                    service);

    JsonObj *json = ulfius_get_json_body_request(request, NULL);

    ret = notify_process_incoming_notification(service, NOTIFICATION_ALERT, json);
    if (ret == STATUS_OK) {
        respCode = RESP_CODE_ACCEPTED;
    }

    /* Send response */
    ulfius_set_empty_body_response(response, respCode);

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
    ulfius_add_endpoint_by_val(&serverInst, method, URL_PREFIX, endPoint, 0, cb,
                    config);
    usys_strcpy(gApi[endPointCount].method, method);
    char api[128] = {0};
    usys_strcpy(api, (URL_PREFIX));
    usys_strcat(api, endPoint);
    usys_strcpy(gApi[endPointCount].endPoint, api);
    usys_log_trace("Added api[%d] Method %s Endpoint: %s.", endPointCount,
                    "Get", gApi[endPointCount].endPoint);
    endPointCount++;
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
 * @fn      void web_service_add_notify_event_endpoints()
 * @brief   Add REST end point for reporting new events.
 *
 */
void web_service_add_notify_event_endpoints() {
    web_service_add_end_point("POST", API_RES_EP("event/:service"), NULL,
                    web_service_cb_post_event);
}

/**
 * @fn      void web_service_add_notify_alert_endpoints()
 * @brief   Add REST end point for reporting new events.
 *
 */
void web_service_add_notify_alert_endpoints() {
    web_service_add_end_point("POST", API_RES_EP("alert/:service"), NULL,
                    web_service_cb_post_alert);
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

    web_service_add_notify_event_endpoints();

    web_service_add_notify_alert_endpoints();

    web_service_add_discover_endpoints();

    /* open connection for web_services */
    if (start_framework(&serverInst)) {
        usys_log_error("Failed to start web_services for notifyd.");
        return STATUS_NOK;
    }

    usys_log_info("WebService on notifyd started.");

    return STATUS_OK;
}

/**
 * @fn      int web_service_start()
 * @brief   Initializes the ulfius framework for REST server.
 * @param   port
 * @return  On success STATUS_OK
 *          On failure STATUS_NOK
 */
int web_service_init(int port) {

    /* Initialize the web_services framework. */
    if (init_framework(&serverInst, port) != STATUS_OK) {
        usys_log_error("Error initializing web_service framework");
        return STATUS_NOK;
    }
    usys_log_info("WebService on notifyd initialized at port %d.", port);
    return STATUS_OK;
}
