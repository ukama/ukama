/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <curl/curl.h>

#include "callback.h"
#include "tasks.h"
#include "jserdes.h"
#include "hub.h"
#include "http_status.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_types.h"

#include "version.h"

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

    if (strstr(srcURL, "https://") == NULL ||
        strstr(srcURL, "http://") == NULL) {
        sprintf(destURL, "%s%s", hubURL, srcURL);
    } else {
        strncpy(destURL, srcURL, strlen(srcURL));
    }

    if (strstr(srcExtraURL, "https://") == NULL ||
        strstr(srcExtraURL, "http://") == NULL) {
        sprintf(destExtraURL, "%s%s", hubURL, srcExtraURL);
    } else {
        strncpy(destExtraURL, srcExtraURL, strlen(srcExtraURL));
    }
}

int web_service_cb_get_capp(const URequest *request,
                            UResponse *response,
                            void *data) {

    int ret=TRUE, resCode=200, i=0;
    int httpStatus=0;
    uuid_t uuid;

    char idStr[36+1]={0};
    char path[WIMC_MAX_PATH_LEN]={0};
    char cbURL[WIMC_MAX_URL_LEN] = {0};

    char indexURL[WIMC_MAX_URL_LEN] = {0};
    char storeURL[WIMC_MAX_URL_LEN] = {0};

    Artifact artifact;
    Config *config;

    char *cappName=NULL, *cappTag=NULL;

    WimcReq   *wimcRequest  = NULL;
    WRespType respType=WRESP_ERROR;
    
    Agent *agent=NULL;
    ArtifactFormat *artifactFormat=NULL;
  
    config = (Config *)data;
    uuid_clear(uuid);

    cappName = (char *)u_map_get(request->map_url, "name");
    cappTag  = (char *)u_map_get(request->map_url, "tag");

    if (cappName == NULL || cappTag == NULL) {
        usys_log_error("capp name:tag not found");
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    if (db_read_path(config->db, cappName, cappTag, &path[0])) {
        usys_log_debug("Path found in db. name:%s tag:%s path:%s",
                       cappName, cappTag, path[0]);
        ulfius_set_string_body_response(response, HttpStatus_OK, path);
        return U_CALLBACK_CONTINUE;
    }

    /* Check with hub */
    if (!get_artifact_info_from_hub(&artifact,
                                    config,
                                    cappName, cappTag,
                                    &httpStatus)) {
        if (httpStatus == HttpStatus_InternalServerError) {
            usys_log_error("Unable to connect with hub at: %s",
                           config->hubURL);
            ulfius_set_string_body_response(response,
                                HttpStatus_InternalServerError,
                                HttpStatusStr(HttpStatus_InternalServerError));
        } else if (httpStatus == HttpStatus_NotFound) {
            usys_log_error("No matching capp %s:%s found by hub: %s",
                           cappName, cappTag, config->hubURL);
            ulfius_set_string_body_response(response,
                                            HttpStatus_NotFound,
                                            HttpStatusStr(HttpStatus_NotFound));
        }

        return U_CALLBACK_CONTINUE;
    } else {
        usys_log_debug("capp %s:%s is available at hub: %s",
                       cappName, cappTag, config->hubURL);
    }

    /* Find matching agent. */
    for (i=0; i < artifact.formatsCount; i++) {
        agent = find_matching_agent(*config->agents,
                                    artifact.formats[i]->type);
        if (agent) {
            artifactFormat = artifact.formats[i];
            uuid_unparse(agent->uuid, &idStr[0]);
            usys_log_debug("Matching agent: %s method: %s URL: %s",
                           idStr, agent->method, agent->url);
            break;
        }
    }

    if (agent == NULL) {
        usys_log_error("No matching agent found for capp %s:%s",
                       cappName, cappTag);
        free_artifact(&artifact);
        ulfius_set_string_body_response(response,
                                HttpStatus_ServiceUnavailable,
                                HttpStatusStr(HttpStatus_ServiceUnavailable));
        return U_CALLBACK_CONTINUE;
    }

    /* create URLs for agent to fetch the artifacts from */
    if (strcmp(agent->method, WIMC_METHOD_CHUNK_STR) == 0) {
        create_hub_urls_for_agent(config->hubURL,
                                  artifactFormat->url,
                                  &indexURL[0],
                                  artifactFormat->extraInfo,
                                  &storeURL[0]);
    }

    /* create request */
    create_wimc_request(&wimcRequest,
                        cappName, cappTag,
                        indexURL,
                        storeURL,
                        artifactFormat->type,
                        DEFAULT_INTERVAL);

    /* Send the request to agent */
    if (communicate_with_agent(wimcRequest, agent->url, config, &uuid)) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_Accepted,
                                        HttpStatusStr(HttpStatus_Accepted));
    } else {
        usys_log_error("Error sending capp fetch request to agent %s:%s",
                       cappName, cappTag);
        ulfius_set_string_body_response(response,
                               HttpStatus_ServiceUnavailable,
                               HttpStatusStr(HttpStatus_ServiceUnavailable));
    }

    cleanup_wimc_request(wimcRequest);
    free_artifact(&artifact);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_agent_update(const struct _u_request *request,
                                     struct _u_response *response,
                                     void *data) {

    int retCode=0;
    char *agentID = NULL;
    json_t *json  = NULL;
    json_error_t jerr;
    AgentReq *agentRequest=NULL;

    Config *config = NULL;

    agentID      = (char *)u_map_get(request->map_url, "id");
    agentRequest = (AgentReq *)calloc(sizeof(AgentReq), 1);
    config       = (Config *)data;
    json         = ulfius_get_json_body_request(request, &jerr);

    if (agentID == NULL || json == NULL ) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }
    
    if (agentRequest == NULL || config == NULL ) {
        ulfius_set_string_body_response(response,
                              HttpStatus_InternalServerError,
                              HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    if (!deserialize_agent_request_update(&agentRequest->update, json)) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    retCode = process_agent_update_request(config->tasks,
                                           agentRequest,
                                           config->db);

    ulfius_set_string_body_response(response,
                                    retCode,
                                    HttpStatusStr(retCode));
    
    free_agent_request_update(agentRequest->update);
    usys_free(agentRequest);
    json_decref(json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_agent(const URequest *request,
                              UResponse *response,
                              void *data) {

    int ret=WIMC_OK, retCode;
    uuid_t uuid;
    json_t *json=NULL, *jMethod=NULL, *jURL=NULL;
    json_error_t jerr;

    char *agentID     = NULL;
    char *agentURL    = NULL;
    char *agentMethod = NULL;

    Config *config = NULL;

    config = (Config *)data;
    
    agentID = (char *)u_map_get(request->map_url, "id");
    if (agentID == NULL) {
        usys_log_error("agent id not found");
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    /* convert id to uuid */
    if (uuid_parse(agentID, uuid) == -1) {
        usys_log_error("agent id not found");
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    json = ulfius_get_json_body_request(request, &jerr);
    if (!json) {
        usys_log_error("JSON error for the agent register request: %s",
                       jerr.text);
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    jMethod = json_object_get(json, JSON_METHOD);
    jURL    = json_object_get(json, JSON_URL);
    if (jURL == NULL || jMethod == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        json_decref(json);
        return U_CALLBACK_CONTINUE;
    } else {
        agentURL    = json_string_value(jURL);
        agentMethod = json_string_value(jMethod);

        if (agentURL == NULL || agentMethod == NULL) {
            ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
            json_decref(json);
            return U_CALLBACK_CONTINUE;
        }
    }

    if (!register_agent(config->agents,agentID, agentMethod, agentURL)) {
        ulfius_set_string_body_response(response,
                                HttpStatus_Conflict,
                                HttpStatusStr(HttpStatus_Conflict));
        
        json_decref(json);
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));
    json_decref(json);
    
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_delete_agent(const URequest *request,
                                UResponse *response,
                                void *data) {

    int retCode;
    uuid_t uuid;

    char *agentID  = NULL;
    Config *config = NULL;

    config  = (Config *)data;
    agentID = (char *)u_map_get(request->map_url, "id");
    if (agentID == NULL) {
        usys_log_error("agent id not found");
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    /* convert id to uuid */
    if (uuid_parse(agentID, uuid) == -1) {
        usys_log_error("agent id not found");
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    /* Find matching agent and delete it */
    if (!delete_agent(config->agents, agentID)) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
    } else {
        ulfius_set_string_body_response(response, HttpStatus_OK,
                                        HttpStatusStr(HttpStatus_OK));
    }
    
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
