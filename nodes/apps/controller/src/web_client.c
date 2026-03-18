/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "web_client.h"
#include "json_types.h"
#include "http_status.h"
#include "config.h"
#include "jserdes.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static int wc_send_http_request(URequest *httpReq,
                                UResponse **httpResp) {

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
                                     char **nodeId,
                                     char **nodeType) {

    int ret = STATUS_NOK;
    JsonObj *json = NULL;
    JsonErrObj jErr;
    UResponse *httpResp = NULL;
    URequest  *httpReq = NULL;
    int type=0;

    httpReq = wc_create_http_request(url, method, NULL);
    if (!httpReq) {
        return ret;
    }

    ret = wc_send_http_request(httpReq, &httpResp);
    if (ret != STATUS_OK) {
        usys_log_error("Failed to send http request.");
        goto cleanup;
    }

    if (httpResp->status == HttpStatus_OK) {
        json = ulfius_get_json_body_response(httpResp, &jErr);
        if (json) {
            json_deserialize_node_info(nodeId, NULL,  JTAG_NODE_ID, JSON_STRING,  json);
            json_deserialize_node_info(NULL,   &type, JTAG_TYPE,    JSON_INTEGER, json);
            if (*nodeId == NULL || type == 0) {
                usys_log_error("Failed to parse NodeInfo response from noded.");
                return STATUS_NOK;
            }

            if (type!=4 ){ /* type as is returned from node.d xxx - fixme */
                *nodeType = strdup("Unknown");
            } else {
                *nodeType = strdup("Amplifier");
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
                                  &config->nodeId,
                                  &config->nodeType) == STATUS_NOK) {
        usys_log_error("Failed to parse NodeInfo response from noded.");
        return STATUS_NOK;
    }

    usys_log_info("%s: Node ID: %s Node Type: %s",
                  SERVICE_NAME,
                  config->nodeId,
                  config->nodeType);

    return STATUS_OK;
}

int alarms_send_notification(const Config *config,
                             AlarmType type,
                             Severity severity,
                             const char *message) {

    char url[256] = {0};
    JsonObj *json = NULL;
    URequest *httpReq = NULL;
    UResponse *httpResp = NULL;
    int ret = -1;

    if (config == NULL || config->enableNotify == 0) {
        return 0;
    }

    if (snprintf(url,
                 sizeof(url),
                 "http://%s:%d%s",
                 config->notifyHost,
                 config->notifyPort,
                 config->notifyPath) >= (int)sizeof(url)) {
        usys_log_error("alarms: notification URL too long");
        return -1;
    }

    json = json_serialize_alarm_notification(config, type, severity, message);
    if (json == NULL) {
        usys_log_error("alarms: failed to serialize notification");
        return -1;
    }

    httpReq = wc_create_http_request(url, "POST", json);
    json_decref(json);
    json = NULL;

    if (httpReq == NULL) {
        usys_log_error("alarms: failed to create HTTP request");
        return -1;
    }

    ret = wc_send_http_request(httpReq, &httpResp);
    if (ret != STATUS_OK) {
        usys_log_error("alarms: failed to send notification to %s", url);
        ret = -1;
        goto cleanup;
    }

    if (httpResp->status != HttpStatus_OK &&
        httpResp->status != HttpStatus_Accepted) {
        usys_log_error("alarms: notify.d returned status %d", httpResp->status);
        ret = -1;
        goto cleanup;
    }

    usys_log_info("alarms: notification sent - type=%s severity=%s",
                  alarm_type_str(type),
                  severity_str(severity));
    ret = 0;

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
