/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "web_client.h"
#include "json_types.h"
#include "http_status.h"
#include "deviced.h"
#include "config.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

/* jserdes.c */
extern bool json_deserialize_node_info(char **data, char *tag, json_t *json);
extern bool json_serialize_alarm_notification(JsonObj **json, Config *config);

static char *ukama_node_type_from_nodeid(const char *nodeID) {

    if (!nodeID || !*nodeID) return NULL;

    /* Match the embedded type token */
    if (strstr(nodeID, "-tnode-") != NULL) return strdup(UKAMA_TOWER_NODE);
    if (strstr(nodeID, "-anode-") != NULL) return strdup(UKAMA_AMPLIFIER_NODE);
    if (strstr(nodeID, "-pnode-") != NULL) return strdup(UKAMA_POWER_NODE);

    return NULL;
}

static bool ukama_config_set_node_type(Config *config) {

    if (!config) return USYS_FALSE;

    char *t = ukama_node_type_from_nodeid(config->nodeID);
    if (!t) return USYS_FALSE;

    free(config->nodeType);
    config->nodeType = t;

    return USYS_TRUE;
}

static int wc_send_http_request(URequest *httpReq, UResponse **httpResp) {

    *httpResp = (UResponse *)usys_calloc(1, sizeof(UResponse));
    if (*httpResp == NULL) {
        usys_log_error("Error allocating memory of size: %lu for http response",
                       sizeof(UResponse));
        return STATUS_NOK;
    }

    if (ulfius_init_response(*httpResp)) {
        usys_log_error("Error initializing new http response.");
        return STATUS_NOK;
    }

    if (ulfius_send_http_request(httpReq, *httpResp) != STATUS_OK) {
        usys_log_error( "Web client failed to send %s web request to %s",
                        httpReq->http_verb, httpReq->http_url);
        return STATUS_NOK;
    }

    return STATUS_OK;
}

static URequest* wc_create_http_request(char *url,
                                        char *method,
                                        JsonObj *body) {

    /* Preparing Request */
    URequest *httpReq = (URequest *)usys_calloc(1, sizeof(URequest));
    if (!httpReq) {
      usys_log_error("Error allocating memory of size: %lu for http Request",
                      sizeof(URequest));
      return NULL;
    }

    if (ulfius_init_request(httpReq)) {
        usys_log_error("Error initializing new http request.");
        return NULL;
    }

    ulfius_set_request_properties(httpReq,
                       U_OPT_HTTP_VERB, method,
                       U_OPT_HTTP_URL, url,
                       U_OPT_TIMEOUT, 20,
                       U_OPT_NONE);

    if (body) {
       if (STATUS_OK != ulfius_set_json_body_request(httpReq, body)) {
           ulfius_clean_request(httpReq);
           usys_free(httpReq);
           httpReq = NULL;
       }
    }

    return httpReq;
}

static int wc_send_node_info_request(char *url,
                                     char *method,
                                     char **nodeID) {

    int ret = STATUS_NOK;
    JsonObj *json = NULL;
    JsonErrObj jErr;
    UResponse *httpResp = NULL;
    URequest *httpReq = NULL;

    httpReq = wc_create_http_request(url, method, NULL);
    if (!httpReq) {
        return ret;
    }

    ret = wc_send_http_request(httpReq, &httpResp);
    if (ret != STATUS_OK) {
        usys_log_error("Failed to send http request.");
       goto cleanup;
    }

    if (httpResp->status == 200) {
        json = ulfius_get_json_body_response(httpResp, &jErr);
        if (json) {
            json_deserialize_node_info(nodeID,   JTAG_NODE_ID, json);
            if (nodeID == NULL) {
                usys_log_error("Failed to parse NodeInfo response from noded.");
                return STATUS_NOK;
            }
            ret = STATUS_OK;
        }
    } else {
        ret = STATUS_NOK;
    }

    json_decref(json);

cleanup:
    if (httpReq) {
        ulfius_clean_request(httpReq);
        usys_free(httpReq);
    }
    if (httpResp) {
        ulfius_clean_response(httpResp);
        usys_free(httpResp);
    }

    return ret;
}

int get_nodeid_and_type_from_noded(Config *config) {

    char url[128] = {0};

    sprintf(url,"http://%s:%d%s", DEF_NODED_HOST,
            config->nodedPort, DEF_NODED_EP);

    if (wc_send_node_info_request(url,
                                  "GET",
                                  &config->nodeID) == STATUS_NOK) {
        usys_log_error("Failed to parse NodeInfo response from noded.");
        return STATUS_NOK;
    }

    if (ukama_config_set_node_type(config) == USYS_FALSE) {
        usys_log_error("Failed to parse node type from nodeID: %s", config->nodeID);
        return STATUS_NOK;
    }

    usys_log_info("Node ID: %s Node Type: %s", config->nodeID, config->nodeType);

    return STATUS_OK;
}

int wc_send_alarm_to_notifyd(Config *config, int *retCode) {

    int ret = USYS_OK;
    char url[128] = {0};
    char *jsonStr=NULL;
    JsonObj *json = NULL;
    UResponse *httpResp = NULL;
    URequest *httpReq = NULL;

    sprintf(url,"http://%s:%d%s%s", DEF_NOTIFY_HOST,
            config->notifydPort, DEF_NOTIFY_EP, config->serviceName);

    if (json_serialize_alarm_notification(&json, config) == USYS_FALSE) {
        usys_log_error("Unable to serialize the notification");
        return USYS_NOK;
    }

    httpReq = wc_create_http_request(url, "POST", json);
    if (!httpReq) {
        json_decref(json);
        return USYS_NOK;
    }

    jsonStr = json_dumps(json, 0);
    usys_log_debug("Sending Notification. URL: %s method: POST, json: %s",
                   url, jsonStr);
    free(jsonStr);

    ret = wc_send_http_request(httpReq, &httpResp);
    if (ret != STATUS_OK || httpResp->status != HttpStatus_Accepted) {
        usys_log_error("Failed sending alarm to notiy.d: %s Code: %d Str: %s",
                       url, httpResp->status,
                       HttpStatusStr(httpResp->status));
        ret = USYS_NOK;
    }

    *retCode = httpResp->status;

    /* cleaup code */
    json_decref(json);
    if (httpReq) {
        ulfius_clean_request(httpReq);
        usys_free(httpReq);
    }

    if (httpResp) {
        ulfius_clean_response(httpResp);
        usys_free(httpResp);
    }

    return ret;
}

int wc_send_reboot_to_client(Config *config, int *retCode) {

    int ret = USYS_OK;
    char url[128] = {0};
    UResponse *httpResp = NULL;
    URequest *httpReq = NULL;

    sprintf(url,"http://%s:%d%s%s",
            config->clientHost,
            config->clientPort,
            URL_PREFIX,
            API_RES_EP("reboot/"));

    httpReq = wc_create_http_request(url, "POST", NULL);
    if (!httpReq) {
        return USYS_NOK;
    }

    usys_log_debug("Sending client reboot. URL: %s", url);

    ret = wc_send_http_request(httpReq, &httpResp);
    if (ret != STATUS_OK || httpResp->status != HttpStatus_Accepted) {
        usys_log_error("Failed sending reboot to client device.d");
        usys_log_error("URL: %s Code: %d Str: %s",
                       url, httpResp->status,
                       HttpStatusStr(httpResp->status));
        ret = USYS_NOK;
    }

    *retCode = httpResp->status;

    /* cleaup code */
    if (httpReq) {
        ulfius_clean_request(httpReq);
        usys_free(httpReq);
    }

    if (httpResp) {
        ulfius_clean_response(httpResp);
        usys_free(httpResp);
    }

    return ret;
}
