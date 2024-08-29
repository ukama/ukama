/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "lookout.h"
#include "web_client.h"
#include "jserdes.h"
#include "json_types.h"
#include "http_status.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_file.h"

/* implemented in resources.c */
extern int    get_memory_usage(int pid);
extern int    get_disk_usage(int pid);
extern double get_cpu_usage(int pid);

static int wc_send_http_request(URequest* httpReq, UResponse** httpResp) {

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

    ret = ulfius_send_http_request(httpReq, *httpResp);
    if (ret != STATUS_OK) {
        usys_log_error( "Web client failed to send %s web request to %s",
                        httpReq->http_verb, httpReq->http_url);
    }
    return ret;
}

URequest* wc_create_http_request(char* url,
                                 char* method,
                                 char* body) {

    JsonObj *jBody = NULL;

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
                                  U_OPT_HEADER_PARAMETER, "User-Agent", SERVICE_NAME,
                                  U_OPT_TIMEOUT, 20,
                                  U_OPT_NONE);

    if (body) {
        jBody = json_loads(body, JSON_DECODE_ANY, NULL);
        if (STATUS_OK != ulfius_set_json_body_request(httpReq, jBody)) {
            ulfius_clean_request(httpReq);
            json_decref(jBody);
            usys_free(httpReq);
            return NULL;
        }
    }

    json_decref(jBody);
    return httpReq;
}

static int wc_send_request(char *url,
                           char *method,
                           char *body,
                           char **buffer) {

    int ret = STATUS_NOK;

    JsonObj   *json = NULL;
    JsonErrObj jErr;
    UResponse *httpResp = NULL;
    URequest  *httpReq = NULL;

    httpReq = wc_create_http_request(url, method, body);
    if (!httpReq) {
        return STATUS_NOK;
    }

    ret = wc_send_http_request(httpReq, &httpResp);
    if (ret != STATUS_OK) {
        usys_log_error("Failed to send http request.");
        goto cleanup;
    }

    if (httpResp->status == HttpStatus_OK ||
        httpResp->status == HttpStatus_Created) {
        json = ulfius_get_json_body_response(httpResp, &jErr);
        if (json) {
            *buffer = json_dumps(json, 0);
            ret     = STATUS_OK;
        }
    } else {
        *buffer = NULL;
        ret     = STATUS_NOK;
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

static int wc_read_from_local_service(Config* config,
                                      int service,
                                      char **buffer) {

    int ret = STATUS_NOK;
    char url[MAX_BUFFER] = {0};

    if (service == SERVICE_NODED) {
        sprintf(url,"http://%s:%d/%s",
                DEF_NODED_HOST,
                config->nodedPort,
                DEF_NODED_EP);
    } else if (service == SERVICE_STARTERD) {
        sprintf(url,"http://%s:%d/%s",
                DEF_STARTERD_HOST,
                config->starterdPort,
                DEF_STARTERD_EP);
    }

    ret = wc_send_request(url, "GET", NULL, buffer);

    if (ret) {
        usys_log_error("failed to parse response from local service");
        return ret;
    }

    return ret;
}

int get_nodeid_from_noded(Config *config) {

    int     ret     = STATUS_OK;
    char    *buffer = NULL;
    JsonObj *json   = NULL;

    if (wc_read_from_local_service(config, SERVICE_NODED, &buffer)) {
        usys_log_error("Error reading NodeID from noded.d");
        return STATUS_NOK;
    }

    json = json_loads(buffer, JSON_DECODE_ANY, NULL);
    if (json_deserialize_node_id(&config->nodeID, json) == USYS_FALSE) {
        usys_log_error("Failed to parse NodeInfo response from noded.");
        ret = STATUS_NOK;
    }

    usys_log_info("lookout.d: Node ID: %s", config->nodeID);

    json_decref(json);
    usys_free(buffer);

    return ret;
}

int get_capps_from_starterd(Config *config, CappList **cappList) {

    int     ret     = STATUS_OK;
    char    *buffer = NULL;
    JsonObj *json   = NULL;

    if (wc_read_from_local_service(config, SERVICE_STARTERD, &buffer)) {
        usys_log_error("Error reading capps from starter.d");
        return STATUS_NOK;
    }

    usys_log_debug("%s: capps: %s", SERVICE_NAME, buffer);

    json = json_loads(buffer, JSON_DECODE_ANY, NULL);
    if (json_deserialize_capps(cappList, json) == USYS_FALSE) {
        usys_log_error("Failed to parse capps response from starterd");
        ret = STATUS_NOK;
    }

    json_decref(json);
    usys_free(buffer);

    return STATUS_OK;
}

int send_health_report(Config *config) {

    CappList    *cappList = NULL;
    CappList    *ptr      = NULL;
    CappRuntime *runtime  = NULL;
    JsonObj     *json     = NULL;

    char url[MAX_BUFFER] = {0};
    char *ukama  = NULL;
    char *buffer = NULL;
    char *report = NULL;

    int ret = USYS_TRUE;

    /* Get capps from starterd; for each get its resource usage */
    if (get_capps_from_starterd(config, &cappList) == STATUS_OK) {
        for (ptr = cappList; ptr; ptr = ptr->next) {

            runtime         = ptr->capp->runtime;

            runtime->memory = get_memory_usage(runtime->pid);
            runtime->disk   = get_disk_usage(runtime->pid);
            runtime->cpu    = get_cpu_usage(runtime->pid);
        }
    } else {
        usys_log_error("Unable to get capp status");
    }

    if (!json_serialize_health_report(&json,
                                      config->nodeID,
                                      cappList)) {
        usys_log_error("Error serializing health report. Ignoring");
        return USYS_FALSE;
    }

    usys_find_ukama_service_address(&ukama);
    sprintf(url,"%s/node/v1/health/nodes/%s/performance",
            ukama, config->nodeID);
    report = json_dumps(json, 0);

    usys_log_debug("Sending to URL: %s the health report %s",
                   url, report);

    if (wc_send_request(url, "POST", report, &buffer) == STATUS_NOK) {
        usys_log_error("failed to parse response from local service");
        ret = USYS_FALSE;
    }

    json_decref(json);
    usys_free(report);
    usys_free(ukama);
    return ret;
}
