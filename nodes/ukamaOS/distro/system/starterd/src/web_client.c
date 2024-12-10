/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <jansson.h>
#include <ulfius.h>
#include <curl/curl.h>
#include <string.h>

#include "starter.h"
#include "config.h"
#include "web_client.h"
#include "http_status.h"

#include "usys_log.h"
#include "usys_types.h"
#include "usys_mem.h"
#include "usys_api.h"
#include "usys_file.h"
#include "usys_services.h"

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

static bool deserialzie_wimc_response(json_t *json, char **path) {

    const char *type, *result;
    json_t *jResp, *obj;

    jResp = json_object_get(json, JSON_WIMC_RESPONSE);
    if (jResp == NULL) {
        return USYS_FALSE;
    }
    
    obj = json_object_get(jResp, JSON_TYPE);
    if (obj == NULL) {
        log_error("Missing response type");
        return USYS_FALSE;
    }
    type = json_string_value(obj);

    obj = json_object_get(jResp, JSON_VOID_STR);
    if (obj == NULL) {
        log_error("Missing str response.");
        return USYS_FALSE;
    }
    result = json_string_value(obj);

    if (strcmp(type, WIMC_RESP_TYPE_RESULT) == 0) {
        *path = strdup(result);
        return USYS_TRUE;
    } else if (strcmp(type, WIMC_RESP_TYPE_ERROR) == 0) {
        *path = NULL;
        log_error("WIMC responded with an error: %s", result);
        return USYS_FALSE;
    } else if (strcmp(type, WIMC_RESP_TYPE_PROCESSING) == 0) {
        *path = NULL;
        log_error("WIMC is processing the request.");
        return USYS_FALSE;
    }

    return USYS_FALSE;
}

static URequest* wc_create_http_request(char *url,
                                        char *method,
                                        JsonObj *body) {

    URequest *httpReq;
    
    httpReq = (URequest *)usys_calloc(1, sizeof(URequest));
    if (httpReq == NULL) {
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

/*
 * get_capp_path -- location of the capp referred by name:tag
 *
 */
int get_capp_path(Config *config, char *name, char *tag,
                  char **path, int *retCode) {

    int ret = USYS_NOK;
    char url[128] = {0};
    UResponse *httpResp = NULL;
    URequest *httpReq = NULL;
    JsonObj *json = NULL;
    JsonErrObj jErr;

    sprintf(url, "%s:%d/%s?name=%s&tag=%s",
            DEF_WIMC_HOST,
            config->wimcPort,
            API_RES_EP("content/containers"),
            name, tag);

    httpReq = wc_create_http_request(url, "POST", NULL);
    if (!httpReq) {
        return USYS_NOK;
    }
    usys_log_debug("Sending capp path request. URL: %s", url);

    ret = wc_send_http_request(httpReq, &httpResp);
    if (ret != STATUS_OK) {
        usys_log_error("Failed sending rquest to wimc.d");
        *retCode = 0;
        ret = USYS_NOK;
        goto done;
    }

    *retCode = httpResp->status;
    json = ulfius_get_json_body_response(httpResp, &jErr);
    if (json) {
        if (deserialzie_wimc_response(json, path) == USYS_FALSE) {
            usys_log_error("Failed to get path from wimc.d for %s:%s",
                           name, tag);
            ret = STATUS_NOK;
        } else {
            ret = STATUS_OK;
        }
    }

    json_decref(json);

done:
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

bool ping_capp(char *name) {

    bool status = USYS_FALSE;
    int port, httpStatus;
    CURL *curl = NULL;
    CURLcode res;
    char url[MAX_BUFFER] = {0};

    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (curl == NULL) {
        usys_log_error("%s: Unable to initialize curl", name);
        return USYS_FALSE;
    }

    port = usys_find_service_port(name);
    if (port <= 0) {
        usys_log_error("Unable to find service in /etc/services: %s", name);
        return USYS_FALSE;
    }

    sprintf(url, "http://localhost:%d/v1/ping", port);
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");

    res = curl_easy_perform(curl);
    if (res != CURLE_OK) {
        usys_log_error("%s: Error sending ping. URL: %s Error: %s",
                       name, url, curl_easy_strerror(res));
        status = USYS_FALSE;
    } else {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &httpStatus);
        if (httpStatus != HttpStatus_OK) {
            usys_log_error("%s: Recevied wrong status for ping request: %s URL: %s",
                           name, HttpStatusStr(httpStatus), url);
            status = USYS_FALSE;
        } else {
            status = USYS_TRUE;
        }
    }

    curl_easy_cleanup(curl);
    curl_global_cleanup();

    return status;
}
