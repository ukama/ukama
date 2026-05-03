/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <jansson.h>
#include <sqlite3.h>
#include <stdlib.h>
#include <string.h>
#include <ulfius.h>

#include "http_status.h"
#include "hub.h"
#include "jserdes.h"
#include "log.h"
#include "wimc.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_types.h"

typedef struct {
    char *buffer;
    size_t size;
    size_t maxSize;
} Response;

static int process_response_from_hub(Artifact ***artifacts, void *resp) {

    Response *response;
    ArtifactFormat *format;
    json_t *json;
    int count;
    int i;
    int j;

    response = (Response *)resp;
    format = NULL;
    json = NULL;
    count = 0;

    if (response == NULL || response->buffer == NULL) {
        return 0;
    }

    json = json_loads(response->buffer, JSON_DECODE_ANY, NULL);
    if (json == NULL) {
        usys_log_error("Failure loading hub response JSON");
        return 0;
    }

    if (!deserialize_hub_response(artifacts, &count, json)) {
        usys_log_error("Failed to deserialize hub response");
        json_decref(json);
        return 0;
    }

    if (count == 0) {
        usys_log_debug("No matching capp available");
        json_decref(json);
        return 0;
    }

    usys_log_debug("Received %d artifacts from hub", count);

    for (i = 0; i < count; i++) {
        usys_log_debug("Name:%s Version:%s Formats:%d",
                       (*artifacts)[i]->name,
                       (*artifacts)[i]->version,
                       (*artifacts)[i]->formatsCount);

        for (j = 0; j < (*artifacts)[i]->formatsCount; j++) {
            format = (*artifacts)[i]->formats[j];
            usys_log_debug("Format:%d Type:%s Url:%s Size:%d",
                           j, format->type, format->url, format->size);
        }
    }

    json_decref(json);
    return count;
}

static bool copy_artifact(Artifact *src, Artifact *dest) {

    int i;

    if (src == NULL || dest == NULL) {
        return USYS_FALSE;
    }

    memset(dest, 0, sizeof(*dest));

    dest->name = src->name ? strdup(src->name) : NULL;
    dest->version = src->version ? strdup(src->version) : NULL;
    dest->formatsCount = src->formatsCount;

    if (dest->name == NULL || dest->version == NULL ||
        dest->formatsCount <= 0) {
        goto failure;
    }

    dest->formats = (ArtifactFormat **)calloc(src->formatsCount,
                                              sizeof(ArtifactFormat *));
    if (dest->formats == NULL) {
        goto failure;
    }

    for (i = 0; i < src->formatsCount; i++) {
        dest->formats[i] = (ArtifactFormat *)calloc(1,
                                                    sizeof(ArtifactFormat));
        if (dest->formats[i] == NULL) {
            goto failure;
        }

        dest->formats[i]->type = src->formats[i]->type ?
                                 strdup(src->formats[i]->type) : NULL;
        dest->formats[i]->url = src->formats[i]->url ?
                                strdup(src->formats[i]->url) : NULL;
        dest->formats[i]->createdAt = src->formats[i]->createdAt ?
                                      strdup(src->formats[i]->createdAt) :
                                      NULL;
        dest->formats[i]->size = src->formats[i]->size;

        if (src->formats[i]->extraInfo != NULL) {
            dest->formats[i]->extraInfo = strdup(src->formats[i]->extraInfo);
        }

        if (dest->formats[i]->type == NULL ||
            dest->formats[i]->url == NULL ||
            dest->formats[i]->createdAt == NULL) {
            goto failure;
        }
    }

    return USYS_TRUE;

failure:
    free_artifact(dest);
    memset(dest, 0, sizeof(*dest));
    return USYS_FALSE;
}

void free_artifact(Artifact *artifact) {

    int i;

    if (artifact == NULL) {
        return;
    }

    usys_free(artifact->name);
    usys_free(artifact->version);

    for (i = 0; i < artifact->formatsCount; i++) {
        if (artifact->formats[i] == NULL) {
            continue;
        }
        usys_free(artifact->formats[i]->type);
        usys_free(artifact->formats[i]->url);
        usys_free(artifact->formats[i]->extraInfo);
        usys_free(artifact->formats[i]->createdAt);
        usys_free(artifact->formats[i]);
    }

    usys_free(artifact->formats);
    memset(artifact, 0, sizeof(*artifact));
}

static size_t response_callback(void *contents, size_t size, size_t nmemb,
                                void *userp) {

    size_t realSize;
    char *newBuffer;
    Response *response;

    realSize = size * nmemb;
    response = (Response *)userp;

    if (response == NULL) {
        return 0;
    }

    if (response->size + realSize + 1 > response->maxSize) {
        usys_log_error("Hub response too large");
        return 0;
    }

    newBuffer = realloc(response->buffer,
                        response->size + realSize + 1);
    if (newBuffer == NULL) {
        usys_log_error("Unable to allocate hub response buffer");
        return 0;
    }

    response->buffer = newBuffer;
    memcpy(&(response->buffer[response->size]), contents, realSize);
    response->size += realSize;
    response->buffer[response->size] = '\0';

    return realSize;
}

static int map_http_status(long code) {

    if (code == HttpStatus_OK) {
        return HttpStatus_OK;
    }

    if (code == HttpStatus_NotFound) {
        return HttpStatus_NotFound;
    }

    if (code == HttpStatus_BadRequest) {
        return HttpStatus_BadRequest;
    }

    if (code >= 500) {
        return HttpStatus_ServiceUnavailable;
    }

    return HttpStatus_InternalServerError;
}

bool get_artifacts_info_from_hub(Artifact *artifact,
                                 Config *config,
                                 const char *hubURL,
                                 char *name,
                                 char *tag,
                                 int *status) {

    int i;
    int count;
    bool ret;
    char hubEP[WIMC_MAX_URL_LEN];
    CURL *curl;
    CURLcode res;
    Response response;
    Artifact **artifacts;
    long httpCode;

    (void)config;

    i = 0;
    count = 0;
    ret = USYS_FALSE;
    curl = NULL;
    artifacts = NULL;
    httpCode = 0;
    memset(hubEP, 0, sizeof(hubEP));
    memset(&response, 0, sizeof(response));

    if (status != NULL) {
        *status = HttpStatus_InternalServerError;
    }

    if (artifact == NULL || name == NULL || tag == NULL ||
        hubURL == NULL || *hubURL == '\0') {
        return USYS_FALSE;
    }

    if (snprintf(hubEP, sizeof(hubEP), "%s/%s/%s", hubURL,
                 WIMC_EP_HUB_APPS, name) >= (int)sizeof(hubEP)) {
        if (status != NULL) {
            *status = HttpStatus_BadRequest;
        }
        return USYS_FALSE;
    }

    curl = curl_easy_init();
    if (curl == NULL) {
        return USYS_FALSE;
    }

    response.buffer = (char *)malloc(1);
    response.size = 0;
    response.maxSize = WIMC_MAX_HTTP_RESPONSE_BYTES;
    if (response.buffer == NULL) {
        goto done;
    }
    response.buffer[0] = '\0';

    curl_easy_setopt(curl, CURLOPT_URL, hubEP);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "wimc/0.1");
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT,
                     WIMC_HTTP_CONNECT_TIMEOUT_SEC);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, WIMC_HTTP_TIMEOUT_SEC);

    res = curl_easy_perform(curl);
    if (res != CURLE_OK) {
        usys_log_error("Hub request failed for %s:%s: %s",
                       name, tag, curl_easy_strerror(res));
        if (status != NULL) {
            *status = HttpStatus_ServiceUnavailable;
        }
        goto done;
    }

    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &httpCode);
    if (httpCode != HttpStatus_OK) {
        usys_log_error("Hub returned HTTP %ld for %s:%s", httpCode,
                       name, tag);
        if (status != NULL) {
            *status = map_http_status(httpCode);
        }
        goto done;
    }

    count = process_response_from_hub(&artifacts, &response);
    if (count == 0) {
        usys_log_debug("No matching capp returned from hub: %s:%s",
                       name, tag);
        if (status != NULL) {
            *status = HttpStatus_NotFound;
        }
        goto done;
    }

    for (i = 0; i < count; i++) {
        if (strcmp(artifacts[i]->version, tag) == 0) {
            if (!copy_artifact(artifacts[i], artifact)) {
                if (status != NULL) {
                    *status = HttpStatus_InternalServerError;
                }
                goto done;
            }
            if (status != NULL) {
                *status = HttpStatus_OK;
            }
            ret = USYS_TRUE;
            goto done;
        }
    }

    if (status != NULL) {
        *status = HttpStatus_NotFound;
    }

 done:
    for (i = 0; i < count; i++) {
        free_artifact(artifacts[i]);
        usys_free(artifacts[i]);
    }

    usys_free(artifacts);
    usys_free(response.buffer);

    if (curl != NULL) {
        curl_easy_cleanup(curl);
    }

    return ret;
}
