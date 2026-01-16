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

static bool is_container_environment(void) {

    /* podman / docker / k8s */
    if (getenv("container") != NULL) {
        return USYS_TRUE;
    }

    /* podman guarantee */
    if (access("/run/.containerenv", F_OK) == 0) {
        return USYS_TRUE;
    }

    return USYS_FALSE;
}

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

static URequest* wc_create_http_request(char* url,
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

static int wc_send_request_raw(char *url,
                               char *method,
                               char *body,
                               long *httpStatus,
                               char **respBody) {

    int ret = STATUS_NOK;

    UResponse *httpResp = NULL;
    URequest  *httpReq  = NULL;

    if (httpStatus) *httpStatus = 0;
    if (respBody)   *respBody   = NULL;

    httpReq = wc_create_http_request(url, method, body);
    if (!httpReq) {
        return STATUS_NOK;
    }

    ret = wc_send_http_request(httpReq, &httpResp);
    if (ret != STATUS_OK) {
        usys_log_error("Failed to send http request.");
        goto cleanup;
    }

    if (httpStatus) {
        *httpStatus = httpResp->status;
    }

    if (respBody && httpResp->binary_body && httpResp->binary_body_length > 0) {
        *respBody = (char *)usys_calloc(1, httpResp->binary_body_length + 1);
        if (*respBody) {
            memcpy(*respBody, httpResp->binary_body, httpResp->binary_body_length);
            (*respBody)[httpResp->binary_body_length] = '\0';
        }
    }

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

static int get_capps_from_supervisord(Config *config, CappList **cappList) {

    FILE *fp = NULL;
    char line[512];
    char procName[128];
    char procState[32];
    int  pid = 0;

    const char *ukamaStatus = "Unknown";

    fp = popen("supervisorctl status", "r");
    if (fp == NULL) {
        usys_log_error("Failed to run supervisorctl");
        return STATUS_NOK;
    }

    while (fgets(line, sizeof(line), fp)) {

        memset(procName, 0, sizeof(procName));
        memset(procState, 0, sizeof(procState));
        pid = 0;

        /* Examples:
         * meshd_latest RUNNING pid 456, uptime 0:10:12
         * bootstrap_latest STOPPED Not started
         */
        if (sscanf(line, "%127s %31s pid %d",
                   procName, procState, &pid) < 2) {
            continue;
        }

        /* Only Ukama apps */
        if (strstr(procName, "_latest") == NULL) {
            continue;
        }

        if (strcmp(procState, "RUNNING") == 0) {
            ukamaStatus = "Active";
        } else if (strcmp(procState, "STARTING") == 0) {
            ukamaStatus = "Pending";
        } else if (strcmp(procState, "STOPPED") == 0) {
            ukamaStatus = "Pending";
            pid = 0;
        } else if (strcmp(procState, "EXITED") == 0 ||
                   strcmp(procState, "FATAL") == 0) {
            ukamaStatus = "Failure";
            pid = 0;
        }

        /*
         * Naming convention:
         * <app>_latest
         * space = "system"
         * tag   = "latest"
         */
        add_capp_to_list(cappList,
                         "system",
                         procName,
                         "latest",
                         ukamaStatus,
                         pid);
    }

    pclose(fp);
    usys_log_debug("Received capps from supervisord");

    return STATUS_OK;
}

static int get_capps_from_starterd(Config *config, CappList **cappList) {

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

    return ret;
}

static void wc_free_gps_data(GPSClientData *gps) {

    if (!gps) return;

    usys_free(gps->coordinates);
    usys_free(gps->gpsTime);
    gps->coordinates = NULL;
    gps->gpsTime     = NULL;
}

static int get_gps_data(GPSClientData *gps) {

    int  ret = STATUS_NOK;
    int  port = 0;
    long status = 0;

    char url[MAX_BUFFER] = {0};
    char *body = NULL;

    if (gps == NULL) {
        return STATUS_NOK;
    }

    gps->gpsLock     = USYS_FALSE;
    gps->coordinates = NULL;
    gps->gpsTime     = NULL;

    /* resolve local GPS service port */
    port = usys_find_service_port(SERVICE_GPS);
    if (port <= 0) {
        usys_log_error("Failed to resolve port for %s", SERVICE_GPS);
        return STATUS_NOK;
    }

    /* lock: GET http://localhost:<port>/v1/lock  -> 200 empty if locked else 404 */
    snprintf(url, sizeof(url), "http://localhost:%d/v1/lock", port);
    ret = wc_send_request_raw(url, "GET", NULL, &status, &body);
    usys_free(body);
    body = NULL;

    if (ret != STATUS_OK) {
        usys_log_error("Failed to read gps lock from %s", url);
        return STATUS_NOK;
    }

    if (status != HttpStatus_OK) {
        gps->gpsLock = USYS_FALSE;
        return STATUS_OK; /* not locked is not a client failure */
    }

    gps->gpsLock = USYS_TRUE;

    /* coordinates: GET http://localhost:<port>/v1/coordinates -> 200 "lon,lat" else 404 */
    snprintf(url, sizeof(url), "http://localhost:%d/v1/coordinates", port);
    ret = wc_send_request_raw(url, "GET", NULL, &status, &body);
    if (ret == STATUS_OK && status == HttpStatus_OK && body && body[0] != '\0') {
        gps->coordinates = strdup(body);
    }
    usys_free(body);
    body = NULL;

    /* time: GET http://localhost:<port>/v1/time -> 200 "<time>" else 404 */
    snprintf(url, sizeof(url), "http://localhost:%d/v1/time", port);
    ret = wc_send_request_raw(url, "GET", NULL, &status, &body);
    if (ret == STATUS_OK && status == HttpStatus_OK && body && body[0] != '\0') {
        gps->gpsTime = strdup(body);
    }
    usys_free(body);
    body = NULL;

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

    int ret    = USYS_TRUE;
    int status = STATUS_NOK;

    GPSClientData gps;
    memset(&gps, 0, sizeof(GPSClientData));

    /* Get capps */
    if (is_container_environment()) {
        status = get_capps_from_supervisord(config, &cappList);
        if (status != STATUS_OK) {
            usys_log_error("Unable to get capps from supervisord");
        }
    } else {
        status = get_capps_from_starterd(config, &cappList);
        if (status != STATUS_OK) {
            usys_log_error("Unable to get capps from starterd");
        }
    }

    /* Collect resource usage only if we got capps */
    if (status == STATUS_OK) {
        for (ptr = cappList; ptr; ptr = ptr->next) {
            runtime = ptr->capp->runtime;
            if (runtime == NULL || runtime->pid <= 0) {
                continue;
            }

            runtime->memory = get_memory_usage(runtime->pid);
            runtime->disk   = get_disk_usage(runtime->pid);
            runtime->cpu    = get_cpu_usage(runtime->pid);
        }
    }

    /* GPS data (best-effort) */
    if (get_gps_data(&gps) != STATUS_OK) {
        gps.gpsLock = USYS_FALSE;
        gps.coordinates = NULL;
        gps.gpsTime = NULL;
    }

    if (!json_serialize_health_report(&json,
                                      config->nodeID,
                                      cappList,
                                      &gps)) {
        usys_log_error("Error serializing health report. Ignoring");
        wc_free_gps_data(&gps);
        return USYS_FALSE;
    }

    wc_free_gps_data(&gps);

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

static void add_capp_to_list(CappList **list,
                             const char *space,
                             const char *name,
                             const char *tag,
                             const char *status,
                             int pid) {

    CappList *newEntry = NULL;
    CappList *tail = NULL;

    if (list == NULL || space == NULL || name == NULL || tag == NULL || status == NULL) {
        return;
    }

    newEntry = (CappList *)calloc(1, sizeof(CappList));
    if (newEntry == NULL) {
        return;
    }

    newEntry->capp = (Capp *)calloc(1, sizeof(Capp));
    if (newEntry->capp == NULL) {
        free(newEntry);
        return;
    }

    newEntry->capp->runtime = (CappRuntime *)calloc(1, sizeof(CappRuntime));
    if (newEntry->capp->runtime == NULL) {
        free(newEntry->capp);
        free(newEntry);
        return;
    }

    newEntry->capp->name  = strdup(name);
    newEntry->capp->tag   = strdup(tag);
    newEntry->capp->space = strdup(space);

    newEntry->capp->runtime->status = strdup(status);
    newEntry->capp->runtime->pid    = pid;

    newEntry->capp->runtime->memory = -1;
    newEntry->capp->runtime->disk   = -1;
    newEntry->capp->runtime->cpu    = -1;

    newEntry->next = NULL;

    /* First entry */
    if (*list == NULL) {
        *list = newEntry;
        return;
    }

    /* Append to tail */
    tail = *list;
    while (tail->next != NULL) {
        tail = tail->next;
    }

    tail->next = newEntry;
}
