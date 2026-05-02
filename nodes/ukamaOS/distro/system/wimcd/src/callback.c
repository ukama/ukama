/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <stdbool.h>
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

bool deserialize_agent_request_update(Update **update, json_t *json);

static bool is_absolute_url(const char *url) {

    if (url == NULL) {
        return false;
    }

    return strncmp(url, "http://", 7) == 0 ||
           strncmp(url, "https://", 8) == 0;
}

static char *parse_hub_override(const URequest *request) {

    json_t *json;
    json_error_t jerr;
    json_t *jhub;
    const char *hub;
    char *ret;

    json = NULL;
    jhub = NULL;
    hub = NULL;
    ret = NULL;

    if (!request || !request->binary_body || request->binary_body_length == 0) {
        return NULL;
    }

    json = ulfius_get_json_body_request(request, &jerr);
    if (!json) {
        return NULL;
    }

    jhub = json_object_get(json, "hub");
    hub = json_is_string(jhub) ? json_string_value(jhub) : NULL;

    if (hub && *hub && is_absolute_url(hub)) {
        ret = strdup(hub);
    }

    json_decref(json);
    return ret;
}

static void free_agent_request_update(AgentReq *req) {

    if (req == NULL || req->update == NULL) {
        return;
    }

    usys_free(req->update->voidStr);
    usys_free(req->update);
}

static void create_hub_urls_for_agent(char *hubURL,
                                      char *srcURL, char *destURL,
                                      char *srcExtraURL, char *destExtraURL) {

    if (hubURL == NULL || srcURL == NULL || destURL == NULL ||
        srcExtraURL == NULL || destExtraURL == NULL) {
        return;
    }

    if (!is_absolute_url(srcURL)) {
        snprintf(destURL, WIMC_MAX_URL_LEN, "%s%s", hubURL, srcURL);
    } else {
        snprintf(destURL, WIMC_MAX_URL_LEN, "%s", srcURL);
    }

    if (!is_absolute_url(srcExtraURL)) {
        snprintf(destExtraURL, WIMC_MAX_URL_LEN, "%s%s", hubURL,
                 srcExtraURL);
    } else {
        snprintf(destExtraURL, WIMC_MAX_URL_LEN, "%s", srcExtraURL);
    }
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

    int httpStatus = 0;
    char indexURL[WIMC_MAX_URL_LEN] = {0};
    char storeURL[WIMC_MAX_URL_LEN] = {0};
    char *name;
    char *tag;
    char *status = NULL;
    char *path = NULL;
    char *actualVersion = NULL;
    char *error = NULL;
    char *hubOverride = NULL;
    const char *hubURL;
    Artifact artifact;
    Config *config;
    WimcReq *wimcRequest = NULL;
    ArtifactFormat *artifactFormat = NULL;
    json_t *jResponse;
    int i;

    memset(&artifact, 0, sizeof(artifact));
    config = (Config *)data;
    name = (char *)u_map_get(request->map_url, "name");
    tag  = (char *)u_map_get(request->map_url, "tag");

    if (!pkg_is_valid_identifier(name) || !pkg_is_valid_identifier(tag)) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    if (validate_cached_package(config, name, tag, &path, &actualVersion)) {
        jResponse = status_json(name, tag, WIMC_STATUS_AVAILABLE, path,
                                actualVersion, NULL);
        ulfius_set_json_body_response(response, HttpStatus_NotModified,
                                      jResponse);
        json_decref(jResponse);
        free(path);
        free(actualVersion);
        return U_CALLBACK_CONTINUE;
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
            free(status);
            free(path);
            free(actualVersion);
            free(error);
            pthread_mutex_unlock(&config->dbMutex);
            return U_CALLBACK_CONTINUE;
        }
    }

    db_update_package_status(config->db, name, tag, NULL,
                             WIMC_STATUS_QUEUED, NULL, NULL);
    pthread_mutex_unlock(&config->dbMutex);

    free(status);
    free(path);
    free(actualVersion);
    free(error);

    hubOverride = parse_hub_override(request);
    hubURL = (hubOverride && *hubOverride) ? hubOverride : config->hubURL;

    if (!get_artifacts_info_from_hub(&artifact, config, hubURL,
                                     name, tag, &httpStatus)) {
        pthread_mutex_lock(&config->dbMutex);
        db_update_package_status(config->db, name, tag, NULL,
                                 WIMC_STATUS_FAILED, NULL,
                                 "hub lookup failed");
        pthread_mutex_unlock(&config->dbMutex);

        ulfius_set_string_body_response(response,
                                        httpStatus ? httpStatus :
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(httpStatus ? httpStatus :
                                        HttpStatus_InternalServerError));
        free(hubOverride);
        return U_CALLBACK_CONTINUE;
    }

    for (i = 0; i < artifact.formatsCount; i++) {
        if (get_agent_port_by_method(artifact.formats[i]->type)) {
            artifactFormat = artifact.formats[i];
            break;
        }
    }

    if (artifactFormat == NULL) {
        pthread_mutex_lock(&config->dbMutex);
        db_update_package_status(config->db, name, tag, NULL,
                                 WIMC_STATUS_FAILED, NULL,
                                 "no matching agent");
        pthread_mutex_unlock(&config->dbMutex);

        free_artifact(&artifact);
        free(hubOverride);
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        HttpStatusStr(
                                        HttpStatus_ServiceUnavailable));
        return U_CALLBACK_CONTINUE;
    }

    if (strcmp(artifactFormat->type, WIMC_METHOD_CHUNK_STR) == 0) {
        create_hub_urls_for_agent((char *)hubURL, artifactFormat->url,
                                  indexURL, artifactFormat->extraInfo,
                                  storeURL);
    }

    create_wimc_request(&wimcRequest, name, tag, indexURL, storeURL,
                        artifactFormat->type, DEFAULT_INTERVAL);
    if (wimcRequest == NULL) {
        free_artifact(&artifact);
        free(hubOverride);
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(
                                        HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
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
    } else {
        pthread_mutex_lock(&config->dbMutex);
        db_update_package_status(config->db, name, tag, NULL,
                                 WIMC_STATUS_FAILED, NULL,
                                 "agent request failed");
        pthread_mutex_unlock(&config->dbMutex);
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        HttpStatusStr(
                                        HttpStatus_ServiceUnavailable));
    }

    cleanup_wimc_request(wimcRequest);
    free_artifact(&artifact);
    free(hubOverride);

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

int web_service_cb_get_status(const URequest *request,
                              UResponse *response,
                              void *epConfig) {

    Config *config;
    json_t *json;

    config = (Config *)epConfig;
    json = json_pack("{s:s, s:s, s:i, s:s}",
                     "service", SERVICE_NAME,
                     "status", "ok",
                     "port", config ? config->servicePort : 0,
                     "pkgDir", DEFAULT_APPS_PKGS_PATH);
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

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

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
