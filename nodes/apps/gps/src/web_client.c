/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <ulfius.h>

#include "web_client.h"
#include "json_types.h"
#include "http_status.h"
#include "gpsd.h"
#include "config.h"
#include "jserdes.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

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

URequest* wc_create_http_request(char* httpURL,
                                 char *urlPath,
                                 char* method,
                                 JsonObj* jBody) {

    char urlWithEp[MAX_URL_LENGTH] = {0};
    URequest* httpReq;

    httpReq = (URequest *)usys_calloc(1, sizeof(URequest));
    if (!httpReq) {
        usys_log_error("Error allocating memory of size: %lu for http Request",
                       sizeof(URequest));
        return NULL;
    }

    if (ulfius_init_request(httpReq)) {
        usys_log_error("Error initializing new http request.");
        return NULL;
    }

    sprintf(urlWithEp, "%s/%s", httpURL, urlPath);
    ulfius_set_request_properties(httpReq,
                                  U_OPT_HTTP_VERB, method,
                                  U_OPT_HTTP_URL, urlWithEp,
                                  U_OPT_HEADER_PARAMETER, "User-Agent", SERVICE_NAME,
                                  U_OPT_TIMEOUT, 20,
                                  U_OPT_NONE);

    if (urlPath) {
        httpReq->url_path = strdup(urlPath);
        if (httpReq->url_path == NULL) {
            usys_log_error("Error allocating memory for URL path");
            ulfius_clean_request(httpReq);
            usys_free(httpReq);
            return NULL;
        }
    }

    if (jBody) {
       if (STATUS_OK != ulfius_set_json_body_request(httpReq, jBody)) {
           ulfius_clean_request(httpReq);
           usys_free(httpReq);
           httpReq = NULL;
       }
    }

    return httpReq;
}

static int wc_send_node_info_request(char *httpURL,
                                     char *urlPath,
                                     char *method,
                                     char **nodeID) {

    int ret = STATUS_NOK;
    JsonObj *json = NULL;
    JsonErrObj jErr;
    UResponse *httpResp = NULL;
    URequest *httpReq = NULL;

    httpReq = wc_create_http_request(httpURL, urlPath, method, NULL);
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
            ret = json_deserialize_node_id(nodeID, json);
            if (!ret) {
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


static int wc_read_node_info(Config* config) {

    int ret = STATUS_NOK;
    char httpURL[128]={0};

    sprintf(httpURL,"http://%s:%d", config->nodedHost, config->nodedPort);
    ret = wc_send_node_info_request(httpURL, config->nodedEP, "GET", &config->nodeID);
    if (ret) {
        usys_log_error("Failed to parse NodeInfo response from noded.");
        return ret;
    }

    return ret;
}

int get_nodeid_from_noded(Config *config) {

    if (wc_read_node_info(config)) {
        usys_log_error("Error reading NodeID from noded.d");
        return STATUS_NOK;
    }

    usys_log_info("notify.d: Node ID: %s", config->nodeID);

    return STATUS_OK;
}
