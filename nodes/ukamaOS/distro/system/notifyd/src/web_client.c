/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "web_client.h"

#include "jserdes.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

//TODO
char* noded_host = "localhost";
int noded_port = 8095;
char* node_info_ep = "/noded/v1/nodeinfo";


int wc_send_http_request( URequest* httpReq , UResponse** httpResp) {
    int ret = STATUS_NOK;

    *httpResp = (UResponse *)usys_calloc(1, sizeof(UResponse));
    if (! (*httpResp)) {
      usys_log_error("Error allocating memory of size: %lu for http response",
                      sizeof(UResponse));
      return STATUS_NOK;
    }

    if (ulfius_init_response(*httpResp)) {
        usys_log_error("Error initializing new http response.");
        return STATUS_NOK;
    }

    ret = ulfius_send_http_request(httpReq,
                    *httpResp);
    if (ret != STATUS_OK) {
        usys_log_error( "Web client failed to send %s web request to %s",
                        httpReq->http_verb, httpReq->http_url);
    }
    return ret;
}

URequest* wc_create_http_request(char* url,
                char* method, JsonObj* body) {

    /* Preparing Request */
    URequest* httpReq = (URequest *)usys_calloc(1, sizeof(URequest));
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

    if(body) {
        ulfius_set_request_properties(httpReq,
                        U_OPT_JSON_BODY, body,
                        U_OPT_NONE);
    }

    return httpReq;
}

int wc_send_node_info_request(char* url, char* method,
                char* nodeID, char* nodeType) {
    int ret = STATUS_NOK;
    JsonObj *json = NULL;
    JsonErrObj jErr;

    UResponse *httpResp = NULL;

    URequest* httpReq = wc_create_http_request(url, method, NULL);
    if (!httpReq) {
        return ret;
    }

    ret = wc_send_http_request(httpReq, &httpResp);
    if (ret != STATUS_OK) {
        usys_log_error("Failed to send http request.");
       goto cleanup;
    }

    if (httpResp->status >= 200 && httpResp->status <= 300) {

        json = ulfius_get_json_body_response(httpResp, &jErr);
        if (json) {
            /* Parse response */
            ret = json_deserialize_node_info(json, nodeID, nodeType);
            if (!ret) {
                usys_log_error("Failed to parse NodeInfo response from noded.");
                return STATUS_NOK;
            }
            ret = STATUS_OK;
        }

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

int wc_forward_notification(char* url, char* method,
                JsonObj* body ) {
    int ret = STATUS_NOK;
    JsonObj *json = NULL;
    JsonErrObj jErr;

    UResponse *httpResp = NULL;

    URequest* httpReq = wc_create_http_request(url, method, body);
    if (!httpReq) {
        return ret;
    }

    char *logbody = json_dumps(body, (JSON_INDENT(4)|JSON_COMPACT|JSON_ENCODE_ANY));
    usys_log_trace("Body is :\n %s", logbody);

    ret = wc_send_http_request(httpReq, &httpResp);
    if (ret != STATUS_OK) {
        usys_log_error("Failed to send http request.");
        goto cleanup;
    }

    if (httpResp->status >= 200 && httpResp->status <= 300) {
        ret = STATUS_OK;
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

int wc_read_node_info(char* nodeID, char* nodeType, char* host, int port) {
    int ret = STATUS_NOK;
    /* Send HTTP request */
    char url[128]={0};

    sprintf(url,"http://%s:%d%s", host, port, node_info_ep);

    ret = wc_send_node_info_request(url, "GET", nodeID, nodeType);
    if (ret) {
        usys_log_error("Failed to parse NodeInfo response from noded.");
        return ret;
    }

    return ret;
}

int web_client_init(char* nodeID, char* nodeType) {

    int ret = wc_read_node_info(nodeID, nodeType, noded_host, noded_port);
    if (ret) {
        usys_log_error("Error reading NodeID from noded.d");
        return STATUS_NOK;
    }

    usys_log_info("NotifyD: Identified unit ID %s and type %s",
                    nodeID, nodeType);

    return STATUS_OK;

}
