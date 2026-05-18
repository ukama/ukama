/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

#include "callback.h"
#include "db.h"
#include "http_status.h"
#include "hub.h"
#include "jserdes.h"
#include "package_cache.h"
#include "tasks.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_types.h"

#include "version.h"

extern void cleanup_wimc_request(WimcReq *request);
extern void create_wimc_request(WimcReq **request, char *name, char *tag,
                                char *indexURL, char *storeURL, char *method,
                                int interval);
extern bool communicate_with_agent(WimcReq *request, char *agentMethod,
                                   Config *config);
extern int process_agent_update_request(WTasks **tasks,
                                        AgentReq *req,
                                        sqlite3 *db);

extern bool deserialize_agent_request_update(Update **update, json_t *json);

typedef struct {

    char *items[WIMC_MAX_HUBS];
    int count;
} HubList;

static bool is_absolute_url(const char *url) {

    if (url == NULL) {
        return false;
    }

    return strncmp(url, "http://", 7) == 0 ||
           strncmp(url, "https://", 8) == 0;
}

static void free_hub_list(HubList *list) {

    int i;

    if (list == NULL) {
        return;
    }

    for (i = 0; i < list->count; i++) {
        free(list->items[i]);
        list->items[i] = NULL;
    }

    list->count = 0;
}

static bool hub_list_add(HubList *list, const char *hub) {

    if (list == NULL || hub == NULL || *hub == '\0') {
        return false;
    }

    if (list->count >= WIMC_MAX_HUBS) {
        return false;
    }

    if (!is_absolute_url(hub)) {
        return false;
    }

    list->items[list->count] = strdup(hub);
    if (list->items[list->count] == NULL) {
        return false;
    }

    list->count++;

    return true;
}

static bool parse_hub_overrides(const URequest *request, HubList *list) {

    json_t *json;
    json_error_t jerr;
    json_t *jhub;
    json_t *item;
    const char *hub;
    size_t i;
    bool ret;

    json = NULL;
    jhub = NULL;
    item = NULL;
    hub = NULL;
    ret = true;

    if (list == NULL) {
        return false;
    }

    memset(list, 0, sizeof(HubList));

    if (request == NULL ||
        request->binary_body == NULL ||
        request->binary_body_length == 0) {
        return true;
    }

    memset(&jerr, 0, sizeof(jerr));

    json = ulfius_get_json_body_request(request, &jerr);
    if (json == NULL) {
        return true;
    }

    jhub = json_object_get(json, "hub");
    if (jhub == NULL) {
        json_decref(json);
        return true;
    }

    if (json_is_string(jhub)) {
        hub = json_string_value(jhub);

        if (hub != NULL && *hub != '\0') {
            ret = hub_list_add(list, hub);
        }

        json_decref(json);
        return ret;
    }

    if (json_is_array(jhub)) {
        if (json_array_size(jhub) == 0 ||
            json_array_size(jhub) > WIMC_MAX_HUBS) {
            json_decref(json);
            return false;
        }

        json_array_foreach(jhub, i, item) {
            hub = json_is_string(item) ? json_string_value(item) : NULL;

            if (hub == NULL || *hub == '\0') {
                ret = false;
                break;
            }

            if (!hub_list_add(list, hub)) {
                ret = false;
                break;
            }
        }

        json_decref(json);
        return ret;
    }

    json_decref(json);

    return false;
}

static bool get_artifacts_info_from_any_hub(Artifact *artifact,
                                            Config *config,
                                            HubList *hubList,
                                            const char *defaultHubURL,
                                            char *name,
                                            char *tag,
                                            char *selectedHub,
                                            size_t selectedHubSize,
                                            int *status) {

    int i;
    int httpStatus;
    const char *hubURL;

    if (artifact == NULL ||
        config == NULL ||
        name == NULL ||
        tag == NULL ||
        selectedHub == NULL ||
        selectedHubSize == 0) {
        if (status != NULL) {
            *status = HttpStatus_InternalServerError;
        }
        return false;
    }

    selectedHub[0] = '\0';

    if (status != NULL) {
        *status = HttpStatus_InternalServerError;
    }

    if (hubList != NULL && hubList->count > 0) {
        for (i = 0; i < hubList->count; i++) {
            hubURL = hubList->items[i];

            if (hubURL == NULL || *hubURL == '\0') {
                continue;
            }

            httpStatus = 0;

            usys_log_debug("Trying hub %s for %s:%s",
                           hubURL, name, tag);

            if (get_artifacts_info_from_hub(artifact, config, hubURL,
                                            name, tag, &httpStatus)) {
                snprintf(selectedHub, selectedHubSize, "%s", hubURL);

                if (status != NULL) {
                    *status = HttpStatus_OK;
                }

                usys_log_debug("Selected hub %s for %s:%s",
                               selectedHub, name, tag);

                return true;
            }

            usys_log_error("Hub failed %s for %s:%s http=%d",
                           hubURL, name, tag, httpStatus);

            if (status != NULL && httpStatus != 0) {
                *status = httpStatus;
            }
        }

        return false;
    }

    if (defaultHubURL == NULL || *defaultHubURL == '\0') {
        if (status != NULL) {
            *status = HttpStatus_BadRequest;
        }
        return false;
    }

    httpStatus = 0;

    if (get_artifacts_info_from_hub(artifact, config, defaultHubURL,
                                    name, tag, &httpStatus)) {
        snprintf(selectedHub, selectedHubSize, "%s", defaultHubURL);

        if (status != NULL) {
            *status = HttpStatus_OK;
        }

        return true;
    }

    if (status != NULL) {
        *status = httpStatus;
    }

    return false;
}

static void free_agent_request_update(AgentReq *req) {

    if (req == NULL || req->update == NULL) {
        return;
    }

    usys_free(req->update->voidStr);
    usys_free(req->update);
}

static void create_hub_url_for_agent(char *hubURL, char *srcURL,
                                     char *destURL) {

    size_t hubLen;
    bool hubSlash;
    bool srcSlash;

    if (hubURL == NULL || srcURL == NULL || destURL == NULL) {
        return;
    }

    if (is_absolute_url(srcURL)) {
        snprintf(destURL, WIMC_MAX_URL_LEN, "%s", srcURL);
        return;
    }

    hubLen = strlen(hubURL);
    if (hubLen == 0 || srcURL[0] == '\0') {
        return;
    }

    hubSlash = hubURL[hubLen - 1] == '/';
    srcSlash = srcURL[0] == '/';

    if (hubSlash && srcSlash) {
        snprintf(destURL, WIMC_MAX_URL_LEN, "%s%s", hubURL, srcURL + 1);
    } else if (!hubSlash && !srcSlash) {
        snprintf(destURL, WIMC_MAX_URL_LEN, "%s/%s", hubURL, srcURL);
    } else {
        snprintf(destURL, WIMC_MAX_URL_LEN, "%s%s", hubURL, srcURL);
    }
}

static void create_hub_urls_for_agent(char *hubURL,
                                      char *srcURL, char *destURL,
                                      char *srcExtraURL, char *destExtraURL) {

    if (hubURL == NULL || srcURL == NULL || destURL == NULL ||
        srcExtraURL == NULL || destExtraURL == NULL) {
        return;
    }

    create_hub_url_for_agent(hubURL, srcURL, destURL);
    create_hub_url_for_agent(hubURL, srcExtraURL, destExtraURL);
}

static char *trim_priority_token(char *value) {

    char *end;

    if (value == NULL) {
        return NULL;
    }

    while (*value == ' ' || *value == '\t' || *value == '\n' ||
           *value == '\r') {
        value++;
    }

    end = value + strlen(value);
    while (end > value && (*(end - 1) == ' ' || *(end - 1) == '\t' ||
           *(end - 1) == '\n' || *(end - 1) == '\r')) {
        end--;
    }
    *end = '\0';

    return value;
}

static ArtifactFormat *find_artifact_format(Artifact *artifact,
                                            const char *method) {

    int i;

    if (artifact == NULL || method == NULL || *method == '\0') {
        return NULL;
    }

    for (i = 0; i < artifact->formatsCount; i++) {
        if (artifact->formats[i] == NULL ||
            artifact->formats[i]->type == NULL) {
            continue;
        }

        if (strcmp(artifact->formats[i]->type, method) == 0) {
            return artifact->formats[i];
        }
    }

    return NULL;
}

static ArtifactFormat *select_artifact_format(Artifact *artifact) {

    char priority[WIMC_MAX_ARGS_LEN];
    char *savePtr;
    char *token;
    char *method;
    const char *env;
    ArtifactFormat *format;
    int i;

    if (artifact == NULL) {
        return NULL;
    }

    env = getenv(WIMC_METHOD_PRIORITY_ENV);
    if (env != NULL && *env != '\0') {
        snprintf(priority, sizeof(priority), "%s", env);
        savePtr = NULL;
        token = strtok_r(priority, ",", &savePtr);

        while (token != NULL) {
            method = trim_priority_token(token);
            format = find_artifact_format(artifact, method);

            if (format != NULL && get_agent_port_by_method(format->type)) {
                usys_log_debug("Selected WIMC method from priority: %s",
                               format->type);
                return format;
            }

            token = strtok_r(NULL, ",", &savePtr);
        }
    }

    for (i = 0; i < artifact->formatsCount; i++) {
        if (artifact->formats[i] == NULL ||
            artifact->formats[i]->type == NULL) {
            continue;
        }

        if (get_agent_port_by_method(artifact->formats[i]->type)) {
            usys_log_debug("Selected WIMC method from hub order: %s",
                           artifact->formats[i]->type);
            return artifact->formats[i];
        }
    }

    return NULL;
}

static json_t *status_json(char *name, char *tag, char *status,
                           char *path, char *actualVersion, char *error) {

    json_t *json;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "name", json_string(name ? name : ""));
    json_object_set_new(json, "tag", json_string(tag ? tag : ""));
    json_object_set_new(json, "status", json_string(status ? status : ""));
    json_object_set_new(json, "message", json_string(status ? status : ""));

    if (path != NULL) {
        json_object_set_new(json, "path", json_string(path));
    }

    if (actualVersion != NULL) {
        json_object_set_new(json, "actualVersion",
                            json_string(actualVersion));
    }

    if (error != NULL) {
        json_object_set_new(json, "error", json_string(error));
    }

    return json;
}

static int validate_cached_package(Config *config, char *name, char *tag,
                                   char **pathOut, char **versionOut) {

    char path[WIMC_MAX_PATH_LEN];
    PackageInfo info;

    if (pathOut != NULL) {
        *pathOut = NULL;
    }
    if (versionOut != NULL) {
        *versionOut = NULL;
    }

    if (pkg_path_for_tag(name, tag, path, sizeof(path)) != 0) {
        return 0;
    }

    if (!pkg_validate_tar(name, tag, path, &info)) {
        return 0;
    }

    pthread_mutex_lock(&config->dbMutex);
    db_update_package_status(config->db, name, tag, path,
                             WIMC_STATUS_AVAILABLE,
                             info.actualVersion, NULL);
    db_update_package_status(config->db, name, info.actualVersion, path,
                             WIMC_STATUS_AVAILABLE,
                             info.actualVersion, NULL);
    pthread_mutex_unlock(&config->dbMutex);

    if (pathOut != NULL) {
        *pathOut = strdup(path);
    }
    if (versionOut != NULL) {
        *versionOut = strdup(info.actualVersion);
    }

    return 1;
}

int web_service_cb_get_app_status(const URequest *request,
                                  UResponse *response,
                                  void *data) {

    char *name;
    char *tag;
    char *status = NULL;
    char *path = NULL;
    char *actualVersion = NULL;
    char *error = NULL;
    Config *config;
    json_t *jResponse;
    PackageInfo info;
    int httpStatus;

    config = (Config *)data;
    name = (char *)u_map_get(request->map_url, "name");
    tag  = (char *)u_map_get(request->map_url, "tag");

    if (!pkg_is_valid_identifier(name) || !pkg_is_valid_identifier(tag)) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&config->dbMutex);
    if (!db_read_package(config->db, name, tag, &status, &path,
                         &actualVersion, &error)) {
        pthread_mutex_unlock(&config->dbMutex);
        if (validate_cached_package(config, name, tag,
                                    &path, &actualVersion)) {
            status = strdup(WIMC_STATUS_AVAILABLE);
            httpStatus = HttpStatus_OK;
        } else {
            status = strdup(WIMC_STATUS_MISSING);
            httpStatus = HttpStatus_NotFound;
        }
    } else {
        if (status != NULL && strcmp(status, WIMC_STATUS_AVAILABLE) == 0) {
            if (path == NULL || !pkg_validate_tar(name, tag, path, &info)) {
                db_update_package_status(config->db, name, tag, path,
                                         WIMC_STATUS_CORRUPT,
                                         actualVersion,
                                         "package missing or invalid");
                free(status);
                status = strdup(WIMC_STATUS_CORRUPT);
                free(error);
                error = strdup("package missing or invalid");
            } else if (actualVersion == NULL) {
                actualVersion = strdup(info.actualVersion);
            }
        }
        pthread_mutex_unlock(&config->dbMutex);
        httpStatus = HttpStatus_OK;
    }

    jResponse = status_json(name, tag, status, path, actualVersion, error);
    ulfius_set_json_body_response(response, httpStatus, jResponse);
    json_decref(jResponse);

    free(status);
    free(path);
    free(actualVersion);
    free(error);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_app(const URequest *request,
                            UResponse *response,
                            void *data) {

    int httpStatus;
    int responseStatus;
    char indexURL[WIMC_MAX_URL_LEN];
    char storeURL[WIMC_MAX_URL_LEN];
    char selectedHub[WIMC_MAX_URL_LEN];
    char *name;
    char *tag;
    char *status;
    char *path;
    char *actualVersion;
    char *error;
    const char *hubURL;
    Artifact artifact;
    Config *config;
    WimcReq *wimcRequest;
    ArtifactFormat *artifactFormat;
    json_t *jResponse;
    HubList hubList;
    bool responseSet;
    bool artifactLoaded;

    httpStatus = 0;
    responseStatus = HttpStatus_InternalServerError;
    status = NULL;
    path = NULL;
    actualVersion = NULL;
    error = NULL;
    hubURL = NULL;
    wimcRequest = NULL;
    artifactFormat = NULL;
    jResponse = NULL;
    responseSet = false;
    artifactLoaded = false;

    memset(indexURL, 0, sizeof(indexURL));
    memset(storeURL, 0, sizeof(storeURL));
    memset(selectedHub, 0, sizeof(selectedHub));
    memset(&artifact, 0, sizeof(artifact));
    memset(&hubList, 0, sizeof(hubList));

    config = (Config *)data;
    name = (char *)u_map_get(request->map_url, "name");
    tag  = (char *)u_map_get(request->map_url, "tag");

    if (!pkg_is_valid_identifier(name) || !pkg_is_valid_identifier(tag)) {
        responseStatus = HttpStatus_BadRequest;
        goto done;
    }

    if (validate_cached_package(config, name, tag, &path, &actualVersion)) {
        jResponse = status_json(name, tag, WIMC_STATUS_AVAILABLE, path,
                                actualVersion, NULL);
        ulfius_set_json_body_response(response, HttpStatus_OK, jResponse);
        json_decref(jResponse);
        jResponse = NULL;
        responseSet = true;
        goto done;
    }

    pthread_mutex_lock(&config->dbMutex);
    if (db_read_package(config->db, name, tag, &status, &path,
                        &actualVersion, &error)) {
        if (strcmp(status, WIMC_STATUS_DOWNLOAD) == 0 ||
            strcmp(status, WIMC_STATUS_DOWNLOADING) == 0 ||
            strcmp(status, WIMC_STATUS_QUEUED) == 0) {
            jResponse = status_json(name, tag, status, path,
                                    actualVersion, error);
            ulfius_set_json_body_response(response, HttpStatus_Conflict,
                                          jResponse);
            json_decref(jResponse);
            jResponse = NULL;
            responseSet = true;
            pthread_mutex_unlock(&config->dbMutex);
            goto done;
        }
    }

    db_update_package_status(config->db, name, tag, NULL,
                             WIMC_STATUS_QUEUED, NULL, NULL);
    pthread_mutex_unlock(&config->dbMutex);

    free(status);
    free(path);
    free(actualVersion);
    free(error);

    status = NULL;
    path = NULL;
    actualVersion = NULL;
    error = NULL;

    if (!parse_hub_overrides(request, &hubList)) {
        pthread_mutex_lock(&config->dbMutex);
        db_update_package_status(config->db, name, tag, NULL,
                                 WIMC_STATUS_FAILED, NULL,
                                 "invalid hub list");
        pthread_mutex_unlock(&config->dbMutex);

        responseStatus = HttpStatus_BadRequest;
        goto done;
    }

    if (!get_artifacts_info_from_any_hub(&artifact, config, &hubList,
                                         config->hubURL, name, tag,
                                         selectedHub, sizeof(selectedHub),
                                         &httpStatus)) {
        pthread_mutex_lock(&config->dbMutex);
        db_update_package_status(config->db, name, tag, NULL,
                                 WIMC_STATUS_FAILED, NULL,
                                 "hub lookup failed");
        pthread_mutex_unlock(&config->dbMutex);

        responseStatus = httpStatus ? httpStatus :
                         HttpStatus_InternalServerError;
        goto done;
    }

    artifactLoaded = true;
    hubURL = selectedHub;

    artifactFormat = select_artifact_format(&artifact);
    if (artifactFormat == NULL) {
        pthread_mutex_lock(&config->dbMutex);
        db_update_package_status(config->db, name, tag, NULL,
                                 WIMC_STATUS_FAILED, NULL,
                                 "no matching agent");
        pthread_mutex_unlock(&config->dbMutex);

        responseStatus = HttpStatus_ServiceUnavailable;
        goto done;
    }

    if (strcmp(artifactFormat->type, WIMC_METHOD_CHUNK_STR) == 0) {
        create_hub_urls_for_agent((char *)hubURL, artifactFormat->url,
                                  indexURL, artifactFormat->extraInfo,
                                  storeURL);
    } else if (strcmp(artifactFormat->type, WIMC_METHOD_TARGZ_STR) == 0) {
        create_hub_url_for_agent((char *)hubURL, artifactFormat->url,
                                 indexURL);
        storeURL[0] = '\0';
    } else {
        pthread_mutex_lock(&config->dbMutex);
        db_update_package_status(config->db, name, tag, NULL,
                                 WIMC_STATUS_FAILED, NULL,
                                 "unsupported package method");
        pthread_mutex_unlock(&config->dbMutex);

        responseStatus = HttpStatus_ServiceUnavailable;
        goto done;
    }

    create_wimc_request(&wimcRequest, name, tag, indexURL, storeURL,
                        artifactFormat->type, DEFAULT_INTERVAL);
    if (wimcRequest == NULL) {
        responseStatus = HttpStatus_InternalServerError;
        goto done;
    }

    if (artifactFormat->size > 0) {
        if ((long)artifactFormat->size > WIMC_MAX_PACKAGE_BYTES) {
            pthread_mutex_lock(&config->dbMutex);
            db_update_package_status(config->db, name, tag, NULL,
                                     WIMC_STATUS_FAILED, NULL,
                                     "package too large");
            pthread_mutex_unlock(&config->dbMutex);

            responseStatus = HttpStatus_BadRequest;
            goto done;
        }

        wimcRequest->fetch->content->expectedSizeBytes =
            (long)artifactFormat->size;
    }

    if (communicate_with_agent(wimcRequest, artifactFormat->type, config)) {
        pthread_mutex_lock(&config->dbMutex);
        db_update_package_status(config->db, name, tag, NULL,
                                 WIMC_STATUS_DOWNLOAD, NULL, NULL);
        pthread_mutex_unlock(&config->dbMutex);

        jResponse = status_json(name, tag, WIMC_STATUS_DOWNLOAD,
                                NULL, NULL, NULL);
        ulfius_set_json_body_response(response, HttpStatus_Accepted,
                                      jResponse);
        json_decref(jResponse);
        jResponse = NULL;
        responseSet = true;
        goto done;
    }

    pthread_mutex_lock(&config->dbMutex);
    db_update_package_status(config->db, name, tag, NULL,
                             WIMC_STATUS_FAILED, NULL,
                             "agent request failed");
    pthread_mutex_unlock(&config->dbMutex);

    responseStatus = HttpStatus_ServiceUnavailable;

done:
    if (wimcRequest != NULL) {
        cleanup_wimc_request(wimcRequest);
        wimcRequest = NULL;
    }

    if (artifactLoaded) {
        free_artifact(&artifact);
    }

    free_hub_list(&hubList);

    free(status);
    free(path);
    free(actualVersion);
    free(error);

    if (!responseSet) {
        ulfius_set_string_body_response(response,
                                        responseStatus,
                                        HttpStatusStr(responseStatus));
    }

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_put_app_stats_update(const struct _u_request *request,
                                        struct _u_response *response,
                                        void *data) {

    int retCode;
    json_t *json;
    json_error_t jerr;
    AgentReq *agentRequest;
    Config *config;

    retCode = 0;
    agentRequest = (AgentReq *)calloc(sizeof(AgentReq), 1);
    config = (Config *)data;
    json = ulfius_get_json_body_request(request, &jerr);

    if (json == NULL) {
        usys_free(agentRequest);
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    if (agentRequest == NULL || config == NULL) {
        json_decref(json);
        usys_free(agentRequest);
        ulfius_set_string_body_response(response,
                              HttpStatus_InternalServerError,
                              HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    if (!deserialize_agent_request_update(&agentRequest->update, json)) {
        json_decref(json);
        usys_free(agentRequest);
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&config->taskMutex);
    pthread_mutex_lock(&config->dbMutex);
    retCode = process_agent_update_request(config->tasks,
                                           agentRequest,
                                           config->db);
    pthread_mutex_unlock(&config->dbMutex);
    pthread_mutex_unlock(&config->taskMutex);

    ulfius_set_string_body_response(response, retCode,
                                    HttpStatusStr(retCode));

    free_agent_request_update(agentRequest);
    usys_free(agentRequest);
    json_decref(json);

    return U_CALLBACK_CONTINUE;
}

static json_t *agent_manager_status_json(AgentManager *mgr,
                                         json_t **summaryOut) {

    json_t *agents;
    json_t *summary;
    json_t *agentJson;
    ManagedAgent *agent;
    int configured;
    int running;
    int failed;
    int i;

    agents = json_array();
    summary = json_object();

    if (summaryOut != NULL) {
        *summaryOut = summary;
    }

    if (agents == NULL || summary == NULL) {
        if (agents != NULL) {
            json_decref(agents);
        }
        if (summary != NULL) {
            json_decref(summary);
        }
        return json_array();
    }

    configured = 0;
    running = 0;
    failed = 0;

    if (mgr == NULL) {
        json_object_set_new(summary, "configured", json_integer(0));
        json_object_set_new(summary, "running", json_integer(0));
        json_object_set_new(summary, "failed", json_integer(0));
        return agents;
    }

    pthread_mutex_lock(&mgr->mutex);

    configured = mgr->count;

    for (i = 0; i < mgr->count; i++) {
        agent = &mgr->agents[i];

        if (agent->running) {
            running++;
        } else {
            failed++;
        }

        agentJson = json_pack("{s:s, s:s, s:b, s:i, s:i, s:i, s:s}",
                              "method", agent->method,
                              "service", agent->service,
                              "running", agent->running ? 1 : 0,
                              "pid", (int)agent->pid,
                              "port", agent->port,
                              "restartCount", agent->restartCount,
                              "execPath", agent->execPath);

        if (agentJson != NULL) {
            json_array_append_new(agents, agentJson);
        }
    }

    pthread_mutex_unlock(&mgr->mutex);

    json_object_set_new(summary, "configured", json_integer(configured));
    json_object_set_new(summary, "running", json_integer(running));
    json_object_set_new(summary, "failed", json_integer(failed));

    return agents;
}

int web_service_cb_get_status(const URequest *request,
                              UResponse *response,
                              void *epConfig) {

    Config *config;
    json_t *json;
    json_t *agents;
    json_t *summary;

    (void)request;

    config = (Config *)epConfig;
    summary = NULL;

    json = json_object();
    if (json == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(
                                        HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    agents = agent_manager_status_json(config ? config->agentManager : NULL,
                                       &summary);

    json_object_set_new(json, "service", json_string(SERVICE_NAME));
    json_object_set_new(json, "status", json_string("ok"));
    json_object_set_new(json, "port",
                        json_integer(config ? config->servicePort : 0));
    json_object_set_new(json, "pkgDir",
                        json_string(DEFAULT_APPS_PKGS_PATH));

    if (summary != NULL) {
        json_object_set_new(json, "agentsSummary", summary);
    }

    if (agents != NULL) {
        json_object_set_new(json, "agents", agents);
    }

    ulfius_set_json_body_response(response, HttpStatus_OK, json);
    json_decref(json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_metrics(const URequest *request,
                               UResponse *response,
                               void *epConfig) {

    Config *config;
    json_t *json;

    config = (Config *)epConfig;
    if (config == NULL || config->db == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(
                                        HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&config->dbMutex);
    json = json_pack("{s:i, s:i, s:i, s:i}",
                     "packages_available",
                     db_count_status(config->db, WIMC_STATUS_AVAILABLE),
                     "packages_downloading",
                     db_count_status(config->db, WIMC_STATUS_DOWNLOAD),
                     "packages_failed",
                     db_count_status(config->db, WIMC_STATUS_FAILED),
                     "packages_corrupt",
                     db_count_status(config->db, WIMC_STATUS_CORRUPT));
    pthread_mutex_unlock(&config->dbMutex);

    ulfius_set_json_body_response(response, HttpStatus_OK, json);
    json_decref(json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig) {

    ulfius_set_empty_body_response(response, HttpStatus_OK);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_MethodNotAllowed,
                                    HttpStatusStr(HttpStatus_MethodNotAllowed));
    return U_CALLBACK_CONTINUE;
}
