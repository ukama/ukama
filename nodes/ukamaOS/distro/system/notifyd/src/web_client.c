/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "notify_macros.h"
#include "jserdes.h"
#include "web_service.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

char gNodeID[32];
char gNodeType[32];

char* noded_host = "localhost";
int noded_port = 8080;
char* node_info_ep = "/noded/v1/unitinfo";


int wc_send_node_info_request(char* url, char* ep, char* method,
                char* nodeID, char* nodeType) {
    int ret = STATUS_NOK;
    JsonObj *json = NULL;
    JsonErrObj jErr;
    URequest httpReq;
    UResponse httpResp;
    ulfius_init_request(&httpReq);
    ulfius_init_response(&httpResp);
    ulfius_set_request_properties(&httpReq,
                    U_OPT_HTTP_VERB, method,
                    U_OPT_HTTP_URL, url,
                    U_OPT_HTTP_URL_APPEND, ep,
                    U_OPT_TIMEOUT, 20,
                    U_OPT_NONE);

    ret = ulfius_send_http_request(&httpReq,
                    &httpResp);
    if (ret != STATUS_OK) {
        usys_log_error( "Web service alert callback function not able "
                        "to notify notification service.");
    }

    usys_log_debug( "Web service alert callback function"
                    " notification  response is %d.",
                    httpResp.status);

    if (httpResp.status >= 200 && httpResp.status >= 300) {

        json = ulfius_get_json_body_response(&httpResp, &jErr);
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
    ulfius_clean_request(&httpReq);
    ulfius_clean_response(&httpReq);
}

int wc_read_node_info(char* nodeID, char* nodeType, char* host, int port) {
    int ret = STATUS_NOK;
    /* Send HTTP request */
    char url[128]={0};

    sprintf(url,"%s:%d", host, port);

    ret = wc_send_node_info_request(url, node_info_ep, "GET", nodeID, nodeType);
    if (ret) {
        usys_log_error("Failed to parse NodeInfo response from noded.");
        return ret;
    }

    return ret;
}

int web_client_init() {
    char* nodeID = NULL;
    char* nodeType = NULL;

    int ret = wc_read_node_info(nodeID, nodeType, noded_host, noded_port);
    if (!ret) {
        usys_log_error("Error reading NodeID from noded.d");
        return STATUS_NOK;
    } else {

        usys_memcpy(gNodeID, nodeID, strlen(nodeID));
        usys_memcpy(gNodeType, nodeType, strlen(nodeType));

    }

    return STATUS_OK;

}
