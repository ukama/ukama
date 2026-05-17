/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <curl/easy.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "agent.h"
#include "common/utils.h"
#include "db.h"
#include "http_status.h"
#include "jserdes.h"
#include "package_cache.h"
#include "tasks.h"
#include "wimc.h"

#include "usys_file.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_services.h"
#include "usys_types.h"

typedef struct {
    char *buffer;
    size_t size;
    size_t maxSize;
} Response;

int get_agent_port_by_method(char *method) {

    char buffer[128];

    if (method == NULL || *method == '\0') {
        return 0;
    }

    if (snprintf(buffer, sizeof(buffer), "wimc-agent-%s", method) >=
        (int)sizeof(buffer)) {
        return 0;
    }

    return usys_find_service_port(buffer);
}

int process_agent_update_request(WTasks **tasks,
                                 AgentReq *req,
                                 sqlite3 *db) {

    Update *update;
    char idStr[36 + 1];
    WTasks *task;
    PackageInfo info;
    int finalState;
    int httpStatus;

    if (tasks == NULL || req == NULL || req->update == NULL || db == NULL) {
        return HttpStatus_InternalServerError;
    }

    if (*tasks == NULL) {
        return HttpStatus_InternalServerError;
    }

    update = req->update;
    uuid_unparse(update->uuid, idStr);
    usys_log_debug("Looking up task with ID: %s", idStr);

    task = find_task_by_uuid(*tasks, update->uuid);
    if (task == NULL) {
        usys_log_error("No record found for ID: %s", idStr);
        return HttpStatus_BadRequest;
    }

    task->update->totalKB = update->totalKB;
    task->update->transferKB = update->transferKB;
    task->update->transferState = update->transferState;
    task->state = update->transferState;

    free(task->update->voidStr);
    task->update->voidStr = update->voidStr ? strdup(update->voidStr) : NULL;

    finalState = 0;
    httpStatus = HttpStatus_OK;

    if (task->state == DONE) {
        free(task->localPath);
        task->localPath = update->voidStr ? strdup(update->voidStr) : NULL;
        finalState = 1;

        if (task->localPath == NULL) {
            db_update_package_status(db, task->content->name,
                                     task->content->tag, NULL,
                                     WIMC_STATUS_FAILED, NULL,
                                     "agent did not return package path");
            httpStatus = HttpStatus_InternalServerError;
            goto done;
        }

        if (pkg_validate_tar(task->content->name, task->content->tag,
                             task->localPath, &info)) {
            db_update_package_status(db, task->content->name,
                                     task->content->tag,
                                     task->localPath,
                                     WIMC_STATUS_AVAILABLE,
                                     info.actualVersion, NULL);
            db_update_package_status(db, task->content->name,
                                     info.actualVersion,
                                     task->localPath,
                                     WIMC_STATUS_AVAILABLE,
                                     info.actualVersion, NULL);
        } else {
            db_update_package_status(db, task->content->name,
                                     task->content->tag,
                                     task->localPath,
                                     WIMC_STATUS_CORRUPT,
                                     info.actualVersion[0] ?
                                     info.actualVersion : NULL,
                                     info.error[0] ? info.error :
                                     "invalid package");
        }
    } else if (task->state == ERR) {
        finalState = 1;
        db_update_package_status(db, task->content->name,
                                 task->content->tag, NULL,
                                 WIMC_STATUS_FAILED, NULL,
                                 update->voidStr ? update->voidStr :
                                 "agent error");
    } else {
        db_update_package_status(db, task->content->name,
                                 task->content->tag, NULL,
                                 WIMC_STATUS_DOWNLOAD, NULL, NULL);
    }

 done:
    if (finalState) {
        delete_from_tasks(tasks, task);
    }

    return httpStatus;
}

void cleanup_wimc_request(WimcReq *request) {

    WFetch *fetch;
    WContent *content;

    if (request == NULL) {
        return;
    }

    fetch = request->fetch;
    if (fetch != NULL) {
        content = fetch->content;
        if (content != NULL) {
            usys_free(content->name);
            usys_free(content->tag);
            usys_free(content->method);
            usys_free(content->indexURL);
            usys_free(content->storeURL);
            usys_free(content);
        }
        usys_free(fetch->cbURL);
        usys_free(fetch);
    }

    usys_free(request);
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
        usys_log_error("Agent response too large");
        return 0;
    }

    newBuffer = realloc(response->buffer,
                        response->size + realSize + 1);
    if (newBuffer == NULL) {
        usys_log_error("Unable to allocate agent response buffer");
        return 0;
    }

    response->buffer = newBuffer;
    memcpy(&(response->buffer[response->size]), contents, realSize);
    response->size += realSize;
    response->buffer[response->size] = '\0';

    return realSize;
}

void create_wimc_request(WimcReq **request,
                         char *name,
                         char *tag,
                         char *indexURL,
                         char *storeURL,
                         char *method,
                         int interval) {

    WFetch *fetch;
    WContent *content;

    if (request == NULL) {
        return;
    }

    *request = NULL;
    fetch = NULL;
    content = NULL;

    *request = (WimcReq *)calloc(1, sizeof(WimcReq));
    fetch    = (WFetch *)calloc(1, sizeof(WFetch));
    content  = (WContent *)calloc(1, sizeof(WContent));

    if (*request == NULL || fetch == NULL || content == NULL) {
        usys_free(*request);
        usys_free(fetch);
        usys_free(content);
        *request = NULL;
        return;
    }

    (*request)->type = WREQ_FETCH;
    uuid_generate(fetch->uuid);
    fetch->interval = interval;

    content->name     = name ? strdup(name) : NULL;
    content->tag      = tag ? strdup(tag) : NULL;
    content->method   = method ? strdup(method) : NULL;
    content->indexURL = indexURL ? strdup(indexURL) : NULL;
    content->storeURL = storeURL ? strdup(storeURL) : strdup("");
    content->expectedSizeBytes = 0;

    if (content->name     == NULL ||
        content->tag      == NULL ||
        content->method   == NULL ||
        content->indexURL == NULL ||
        content->storeURL == NULL) {
        cleanup_wimc_request(*request);
        *request = NULL;
        return;
    }

    fetch->content    = content;
    (*request)->fetch = fetch;
}

static bool send_request_to_agent(char *name, char *tag,
                                  char *agentMethod,
                                  const json_t *json,
                                  long *statusCode) {

    bool ret;
    CURL *curl;
    char *jStr;
    char agentURL[WIMC_MAX_URL_LEN];
    CURLcode res;
    struct curl_slist *headers;
    Response response;
    int agentPort;

    ret = USYS_FALSE;
    curl = NULL;
    jStr = NULL;
    headers = NULL;
    agentPort = 0;
    memset(agentURL, 0, sizeof(agentURL));
    memset(&response, 0, sizeof(response));

    if (statusCode != NULL) {
        *statusCode = 0;
    }

    agentPort = get_agent_port_by_method(agentMethod);
    if (agentPort <= 0) {
        usys_log_error("No agent port for method: %s", agentMethod);
        return USYS_FALSE;
    }

    if (snprintf(agentURL, sizeof(agentURL),
                 "http://localhost:%d/v1/apps/%s/%s", agentPort,
                 name, tag) >= (int)sizeof(agentURL)) {
        usys_log_error("Agent URL too long");
        return USYS_FALSE;
    }

    curl = curl_easy_init();
    if (curl == NULL) {
        usys_log_error("Error initializing curl");
        return USYS_FALSE;
    }

    response.buffer = malloc(1);
    response.size = 0;
    response.maxSize = WIMC_MAX_HTTP_RESPONSE_BYTES;
    if (response.buffer == NULL) {
        goto done;
    }
    response.buffer[0] = '\0';

    jStr = json_dumps(json, 0);
    if (jStr == NULL) {
        goto done;
    }

    usys_log_debug("Sending request to Agent: %s", jStr);

    headers = curl_slist_append(headers, "Accept: application/json");
    headers = curl_slist_append(headers, "Content-Type: application/json");
    headers = curl_slist_append(headers, "charset: utf-8");

    curl_easy_setopt(curl, CURLOPT_URL, agentURL);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "POST");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, jStr);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "wimc/0.1");
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT,
                     WIMC_HTTP_CONNECT_TIMEOUT_SEC);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, WIMC_HTTP_TIMEOUT_SEC);

    res = curl_easy_perform(curl);
    if (res != CURLE_OK) {
        usys_log_error("Error sending request to Agent: %s",
                       curl_easy_strerror(res));
        goto done;
    }

    if (statusCode != NULL) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, statusCode);
    }
    ret = USYS_TRUE;

 done:
    usys_free(jStr);
    usys_free(response.buffer);
    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);

    return ret;
}

bool communicate_with_agent(WimcReq *request,
                            char *agentMethod,
                            Config *config) {

    json_t *json;
    long agentRetCode;
    WTasks *task;

    json = NULL;
    agentRetCode = 0;
    task = NULL;

    if (request == NULL || request->fetch == NULL || config == NULL) {
        return USYS_FALSE;
    }

    if (!serialize_wimc_request(request, &json)) {
        usys_log_error("Error serializing wimc request to agent");
        return USYS_FALSE;
    }

    /*
     * Add task before contacting the agent.
     *
     * The agent can accept the request and immediately call back with FETCH.
     * If the task is inserted only after send_request_to_agent() returns,
     * the callback can arrive first and WIMC will reject the update because
     * it cannot find the UUID.
     */
    pthread_mutex_lock(&config->taskMutex);
    add_to_tasks(config->tasks, request);
    pthread_mutex_unlock(&config->taskMutex);

    if (send_request_to_agent(request->fetch->content->name,
                              request->fetch->content->tag,
                              agentMethod, json, &agentRetCode)) {
        if (agentRetCode == HttpStatus_OK) {
            usys_log_debug("Agent initiated to fetch capp");
            json_decref(json);
            return USYS_TRUE;
        }

        usys_log_error("Agent returned an error: %ld", agentRetCode);
    } else {
        usys_log_error("Error communicating with Agent");
    }

    pthread_mutex_lock(&config->taskMutex);
    task = find_task_by_uuid(*config->tasks, request->fetch->uuid);
    if (task != NULL) {
        delete_from_tasks(config->tasks, task);
    }
    pthread_mutex_unlock(&config->taskMutex);

    json_decref(json);
    return USYS_FALSE;
}

void clear_agents(Agent *agent) {

    int i;

    if (agent == NULL) {
        return;
    }

    for (i = 0; i < MAX_AGENTS; i++) {
        if (uuid_is_null(agent[i].uuid) == 0) {
            free(agent[i].method);
            free(agent[i].url);
        }
    }
}
