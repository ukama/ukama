/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <curl/curl.h>

#include "db.h"
#include "callback.h"
#include "tasks.h"
#include "jserdes.h"
#include "hub.h"
#include "http_status.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_types.h"

#include "version.h"

/* agent.c */
extern void cleanup_wimc_request(WimcReq *request);
extern void create_wimc_request(WimcReq **request, char *name, char *tag,
                                char *indexURL, char *storeURL, char *method,
                                int interval);
extern bool communicate_with_agent(WimcReq *request, char *agentMethod,
                                   Config *config);
extern int process_agent_update_request(WTasks **tasks,
                                        AgentReq *req,
                                        sqlite3 *db);

/* jserdes.c */
bool deserialize_agent_request_update(Update **update, json_t *json);

static bool is_absolute_url(const char *url) {

    if (url == NULL) {
        return false;
    }

    if (strncmp(url, "http://", 7) == 0 ||
        strncmp(url, "https://", 8) == 0) {
        return true;
    }

    return false;
}

static bool is_valid_identifier(const char *value) {

    size_t i;

    if (value == NULL || *value == '\0') {
        return false;
    }

    if (strlen(value) >= WIMC_MAX_NAME_LEN) {
        return false;
    }

    for (i = 0; value[i] != '\0'; i++) {
        if ((value[i] >= 'a' && value[i] <= 'z') ||
            (value[i] >= 'A' && value[i] <= 'Z') ||
            (value[i] >= '0' && value[i] <= '9') ||
            value[i] == '-' || value[i] == '_' || value[i] == '.') {
            continue;
        }

        return false;
    }

    return true;
}

static void free_agent_request_update(AgentReq *req) {

    usys_free(req->update->voidStr);
    usys_free(req->update);
}

static void create_hub_urls_for_agent(char *hubURL,
                                      char *srcURL, char *destURL,
                                      char *srcExtraURL, char *destExtraURL) {

    if (hubURL == NULL || srcURL == NULL ||
        destURL == NULL || srcExtraURL == NULL ||
        destExtraURL == NULL)
        return;

    if (!is_absolute_url(srcURL)) {
        snprintf(destURL, WIMC_MAX_URL_LEN, "%s%s", hubURL, srcURL);
    } else {
        snprintf(destURL, WIMC_MAX_URL_LEN, "%s", srcURL);
    }

    if (!is_absolute_url(srcExtraURL)) {
        snprintf(destExtraURL, WIMC_MAX_URL_LEN, "%s%s", hubURL, srcExtraURL);
    } else {
        snprintf(destExtraURL, WIMC_MAX_URL_LEN, "%s", srcExtraURL);
    }
}

static int file_exists_and_non_empty(char *name, char *tag) {


    char *fileName = NULL;
    FILE *file     = NULL;
    long filesize  = 0;

    fileName = (char *)malloc((strlen(DEFAULT_APPS_PKGS_PATH) +
                               strlen(name) + strlen(tag) + 16)*sizeof(char));

    snprintf(fileName, strlen(DEFAULT_APPS_PKGS_PATH) + strlen(name) + strlen(tag) + 16, "%s/%s_%s.tar.gz",
            DEFAULT_APPS_PKGS_PATH,
            name, tag);

    file = fopen(fileName, "r");
    if (file == NULL) {
        free(fileName);
        return 0;
    }

    fseek(file, 0, SEEK_END);
    filesize = ftell(file);
    fclose(file);

    free(fileName);

    if (filesize > 0) {
        return 1;
    }

    return 0;
}

int web_service_cb_get_app_status(const URequest *request,
                                  UResponse *response,
                                  void *data) {

    char   *name=NULL, *tag=NULL, *status=NULL;
    Config *config=NULL;
    json_t *jResponse = NULL;

    config = (Config *)data;

    name = (char *)u_map_get(request->map_url, "name");
    tag  = (char *)u_map_get(request->map_url, "tag");

    if (name == NULL || tag == NULL || !is_valid_identifier(name) || !is_valid_identifier(tag)) {
        usys_log_error("app name:tag not found in the request.");
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&config->dbMutex);
    if (db_read_status(config->db, name, tag, &status)) {
        if (strcmp(status, "download") == 0) {
            jResponse = json_pack("{s:s}",
                                  "message", "download");
            ulfius_set_json_body_response(response,
                                          HttpStatus_OK,
                                          jResponse);
            json_decref(jResponse);
        } else if (strcmp(status, "available") == 0) {
            if (file_exists_and_non_empty(name, tag)) {
                usys_log_debug("app found in the default location");
                jResponse = json_pack("{s:s}",
                                      "message", "available");
                ulfius_set_json_body_response(response,
                                              HttpStatus_OK,
                                              jResponse);
                json_decref(jResponse);
            } else {
                usys_log_error("app is corrupted at default location");
                jResponse = json_pack("{s:s}",
                                      "message", "App corrupted at default location.");
                ulfius_set_json_body_response(response,
                                            HttpStatus_InternalServerError,
                                            jResponse);
                json_decref(jResponse);
            }
        } else if (strcmp(status, "failed") == 0) {
            jResponse = json_pack("{s:s}",
                                  "message", "failed");
            ulfius_set_json_body_response(response,
                                          HttpStatus_OK,
                                          jResponse);
            json_decref(jResponse);
        } else {
            usys_log_error("Unknown status found for app '%s:%s'.", name, tag);
            jResponse = json_pack("{s:s}",
                                  "message", "Unknown app status.");
            ulfius_set_json_body_response(response,
                                          HttpStatus_InternalServerError,
                                          jResponse);
            json_decref(jResponse);
        }
        free(status);
    } else {
        jResponse = json_pack("{s:s}",
                              "message", HttpStatusStr(HttpStatus_NotFound));
        ulfius_set_json_body_response(response,
                                      HttpStatus_NotFound,
                                      jResponse);
        json_decref(jResponse);
    }
    pthread_mutex_unlock(&config->dbMutex);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_app(const URequest *request,
                            UResponse *response,
                            void *data) {

    int httpStatus=0;

    char indexURL[WIMC_MAX_URL_LEN] = {0};
    char storeURL[WIMC_MAX_URL_LEN] = {0};

    Artifact artifact;
    Config *config;

    char *name=NULL, *tag=NULL, *status=NULL;

    WimcReq   *wimcRequest  = NULL;
    ArtifactFormat *artifactFormat=NULL;

    config = (Config *)data;

    name = (char *)u_map_get(request->map_url, "name");
    tag  = (char *)u_map_get(request->map_url, "tag");

    if (name == NULL || tag == NULL || !is_valid_identifier(name) || !is_valid_identifier(tag)) {
        usys_log_error("capp name:tag not found");
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    /*
     * if app is downloading -> 409 (conflict)
     * if app is 'available' but not found in pkg -> start downloading - 202
     * if app is 'available and also found in pkg -> 304
     */
    pthread_mutex_lock(&config->dbMutex);
    if (db_read_status(config->db, name, tag, &status)) {
        if (strcmp(status, "download") == 0) {
            usys_log_debug("capp found in db. name:%s tag:%s status:%s",
                           name, tag, status);
            ulfius_set_string_body_response(response,
                                            HttpStatus_Conflict,
                                            HttpStatusStr(HttpStatus_Conflict));
            free(status);
            pthread_mutex_unlock(&config->dbMutex);
            return U_CALLBACK_CONTINUE;
        } else if (strcmp(status, "available") == 0) {
            if (file_exists_and_non_empty(name, tag)) {
                usys_log_debug("capp found in the default location");
                ulfius_set_string_body_response(response,
                                                HttpStatus_NotModified,
                                                HttpStatusStr(HttpStatus_NotModified));
                free(status);
                pthread_mutex_unlock(&config->dbMutex);
                return U_CALLBACK_CONTINUE;
            }
        }
        free(status);
    }
    pthread_mutex_unlock(&config->dbMutex);

    /* Check with hub */
    if (!get_artifacts_info_from_hub(&artifact, config, name, tag, &httpStatus)) {
        if (httpStatus == HttpStatus_InternalServerError) {
            usys_log_error("Unable to connect with hub at: %s", config->hubURL);
            ulfius_set_string_body_response(response,
                                HttpStatus_InternalServerError,
                                HttpStatusStr(HttpStatus_InternalServerError));
        } else if (httpStatus == HttpStatus_NotFound) {
            usys_log_error("No matching capp %s:%s found by hub: %s",
                           name, tag, config->hubURL);
            ulfius_set_string_body_response(response,
                                            HttpStatus_NotFound,
                                            HttpStatusStr(HttpStatus_NotFound));
        }

        return U_CALLBACK_CONTINUE;
    } else {
        usys_log_debug("capp %s:%s is available at hub", name, tag);
    }

    /* Find matching agent */
    for (int i=0; i < artifact.formatsCount; i++) {
        if (get_agent_port_by_method(artifact.formats[i]->type)) {
            artifactFormat = artifact.formats[i];
            usys_log_debug("Matching agent for method: %s",
                           artifact.formats[i]->type);
            break;
        }
    }

    if (artifactFormat == NULL) {
        usys_log_error("No matching agent found for app %s:%s", name, tag);
        free_artifact(&artifact);
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    /* create URLs for agent to fetch the artifacts from */
    if (strcmp(artifactFormat->type, WIMC_METHOD_CHUNK_STR) == 0) {
        create_hub_urls_for_agent(config->hubURL,
                                  artifactFormat->url,
                                  &indexURL[0],
                                  artifactFormat->extraInfo,
                                  &storeURL[0]);
    }

    /* create request */
    create_wimc_request(&wimcRequest,
                        name, tag,
                        indexURL,
                        storeURL,
                        artifactFormat->type,
                        DEFAULT_INTERVAL);
    if (wimcRequest == NULL) {
        free_artifact(&artifact);
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    /* Send the request to agent */
    if (communicate_with_agent(wimcRequest, artifactFormat->type, config)) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_Accepted,
                                        HttpStatusStr(HttpStatus_Accepted));
        pthread_mutex_lock(&config->dbMutex);
        db_update_status(config->db, name, tag, "download");
        pthread_mutex_unlock(&config->dbMutex);
    } else {
        usys_log_error("Error sending capp fetch request to agent %s:%s", name, tag);
        ulfius_set_string_body_response(response,
                               HttpStatus_ServiceUnavailable,
                               HttpStatusStr(HttpStatus_ServiceUnavailable));
    }

    cleanup_wimc_request(wimcRequest);
    free_artifact(&artifact);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_put_app_stats_update(const struct _u_request *request,
                                        struct _u_response *response,
                                        void *data) {

    int retCode=0;
    json_t *json  = NULL;
    json_error_t jerr;
    AgentReq *agentRequest=NULL;

    Config *config = NULL;

    agentRequest = (AgentReq *)calloc(sizeof(AgentReq), 1);
    config       = (Config *)data;
    json         = ulfius_get_json_body_request(request, &jerr);

    if (json == NULL ) {
        usys_free(agentRequest);
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    if (agentRequest == NULL || config == NULL ) {
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

    ulfius_set_string_body_response(response,
                                    retCode,
                                    HttpStatusStr(retCode));
    
    free_agent_request_update(agentRequest);
    usys_free(agentRequest);
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
