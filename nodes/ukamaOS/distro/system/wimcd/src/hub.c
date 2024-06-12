/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <sqlite3.h>
#include <jansson.h>
#include <ulfius.h>
#include <curl/curl.h>
#include <string.h>

#include "wimc.h"
#include "log.h"
#include "hub.h"
#include "jserdes.h"
#include "http_status.h"

#include "usys_types.h"
#include "usys_log.h"
#include "usys_mem.h"

struct Response {
    char *buffer;
    size_t size;
};

static int process_response_from_hub(Artifact ***artifacts, void *resp) {

    struct Response *response=NULL;
    ArtifactFormat  *format=NULL;
    json_t          *json=NULL;
    int count=0;

    response = (struct Response *)resp;

    json = json_loads(response->buffer, JSON_DECODE_ANY, NULL);
    if (!json) {
        usys_log_error("failure loading str into JSON object. Str: %s",
                       response->buffer);
        return 0;
    }

    if (!deserialize_hub_response(artifacts, &count, json)) {
        usys_log_error("Deser failed for hub response: %s", response->buffer);
        json_decref(json);
        return 0;
    }

    if (count==0) {
        usys_log_debug("No matching capp available");
        json_decref(json);
        return 0;
    }

    usys_log_debug("Received Artifacts from the hub. %d:", count);

    for (int i=0; i<count; i++) {
        usys_log_debug("\n\t Name: %s \n\t Version: %s \n\t Formats: %d",
                       (*artifacts)[i]->name,
                       (*artifacts)[i]->version,
                       (*artifacts)[i]->formatsCount);

        for (int j=0; j<(*artifacts)[i]->formatsCount; j++) {
            format = (*artifacts)[i]->formats[j];
            usys_log_debug("\n\t %d:\n \t\t type: %s \n\t\t "
                           "url: %s \n\t\t createdAt: %s \n\t\t size: %d",
                           j, format->type,
                           format->url,
                           format->createdAt,
                           format->size);
            if (format->extraInfo) {
                log_debug("\t\t extra: %s", format->extraInfo);
            }
        }
    }

    json_decref(json);

    return count;
}

static bool copy_artifact(Artifact *src, Artifact *dest) {

    if (src == NULL || dest == NULL) return USYS_FALSE;

    dest->name         = strdup(src->name);
    dest->version      = strdup(src->version);
    dest->formatsCount = src->formatsCount;

    dest->formats = (ArtifactFormat **)calloc(src->formatsCount,
                                              sizeof(ArtifactFormat *));
    if (dest->formats == NULL) {
        goto failure;
    }

    for (int i=0; i<src->formatsCount; i++) {

        dest->formats[i] = (ArtifactFormat *)calloc(1, sizeof(ArtifactFormat));

        dest->formats[i]->type      = strdup(src->formats[i]->type);
        dest->formats[i]->url       = strdup(src->formats[i]->url);
        dest->formats[i]->createdAt = strdup(src->formats[i]->createdAt);
        dest->formats[i]->size      = src->formats[i]->size;

        if (src->formats[i]->extraInfo) {
            dest->formats[i]->extraInfo = strdup(src->formats[i]->extraInfo);
        }
    }

    return USYS_TRUE;

failure:
    usys_free(dest->name);
    usys_free(dest->version);

    return USYS_FALSE;
}

void free_artifact(Artifact *artifact) {

    if (artifact == NULL) return;

    usys_free(artifact->name);
    usys_free(artifact->version);

    for (int i=0; i<artifact->formatsCount; i++) {
        usys_free(artifact->formats[i]->type);
        usys_free(artifact->formats[i]->url);
        usys_free(artifact->formats[i]->extraInfo);
        usys_free(artifact->formats[i]->createdAt);
        usys_free(artifact->formats[i]);
    }

    usys_free(artifact->formats);
}

static size_t response_callback(void *contents, size_t size, size_t nmemb,
                                void *userp) {

    size_t realsize = size * nmemb;
    struct Response *response = (struct Response *)userp;

    response->buffer = realloc(response->buffer, response->size + realsize + 1);

    if(response->buffer == NULL) {
        usys_log_error("Not enough memory to realloc of size: %s",
                       response->size + realsize + 1);
        return 0;
    }

    memcpy(&(response->buffer[response->size]), contents, realsize);
    response->size += realsize;
    response->buffer[response->size] = 0; /* Null terminate. */

    return realsize;
}

/*
 * get_artifacts_info_from_hub --
 *
 * -1: curl error and curlCode is set.
 *  0: error processing response.
 *  1: Success and curlCode is CURLE_OK
 */
bool get_artifact_info_from_hub(Artifact *artifact,
                                Config *config,
                                char *name, char *tag,
                                int *status) {

    int i, count=0;
    bool ret=USYS_FALSE;
    char hubEP[WIMC_MAX_URL_LEN] = {0};

    CURL *curl=NULL;
    CURLcode res;
    struct Response response;

    Artifact **artifacts=NULL;

    if (name == NULL || tag == NULL) return USYS_FALSE;

    /* create HUB EP: http://localhost:18300/v1/hub/apps/:name */
    sprintf(hubEP, "%s/%s/%s", config->hubURL, WIMC_EP_HUB_APPS, name);

    curl = curl_easy_init();
    if (curl == NULL) {
        *status = HttpStatus_InternalServerError;
        return USYS_FALSE;
    }

    response.buffer = (char *)malloc(1);
    response.size   = 0;

    curl_easy_setopt(curl, CURLOPT_URL, hubEP);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "wimc/0.1");

    res = curl_easy_perform(curl);

    if (res != CURLE_OK) {
        usys_log_error("Error sending request to hub for %s:%s: %s",
                       curl_easy_strerror(res), name, tag);
        *status = HttpStatus_InternalServerError;
        goto done;
    }

    /* get status code. */
    count = process_response_from_hub(&artifacts, &response);
    if (count == 0) { /* No matching capp found by 'name' */
        usys_log_debug("No matching capp returned from the hub "
                       "Requested: %s:%s", name, tag);
        *status = HttpStatus_NotFound;
        goto done;
    }

    /* Find matching capp (with right tag/version) */
    for (i=0; i<count; i++) {
        if (strcmp(artifacts[i]->version, tag) == 0) {
            copy_artifact(artifacts[i], artifact);
            break;
        }
    }

    *status = HttpStatus_OK;
    ret = USYS_TRUE;

done:
    for (i=0; i<count; i++) {
        free_artifact(artifacts[i]);
        usys_free(artifacts[i]);
    }
    usys_free(artifacts);
    usys_free(response.buffer);
    curl_easy_cleanup(curl);

    return ret;
}
