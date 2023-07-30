/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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

static void free_agent_request(AgentReq *req) {
#if 0
    if (req->type == REQ_REG) {
        free(req->reg->method);
        free(req->reg->url);
        free(req->reg);
    } else if (req->type == REQ_UPDATE) {
        if (req->update->voidStr) free(req->update->voidStr);
        free(req->update);
    }

    free(req);
#endif
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
        sprintf(destURL, "%s/%s", hubURL, srcURL);
    } else {
        strncpy(destURL, srcURL, strlen(srcURL));
    }

    if (strstr(srcExtraURL, "https://") == NULL ||
        strstr(srcExtraURL, "http://") == NULL) {
        sprintf(destExtraURL, "%s/%s", hubURL, srcExtraURL);
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
    char url[WIMC_MAX_URL_LEN]={0};
    char cbURL[WIMC_MAX_URL_LEN] = {0};
    char extraURL[WIMC_MAX_URL_LEN]={0};

    Artifact artifact;
    Config *config;

    char *errStr=NULL;
    char *respBody=NULL, *name=NULL, *tag=NULL;

    WRespType respType=WRESP_ERROR;
    Agent *agent=NULL;
    ArtifactFormat *artifactFormat=NULL;
  
    config = (Config *)data;
    uuid_clear(uuid);

    name = (char *)u_map_get(request->map_url, "name");
    tag  = (char *)u_map_get(request->map_url, "tag");

    if (name == NULL || tag == NULL) {
        usys_log_error("capp name:tag not found");
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        goto cleanup;
    }

#if 0    
    if (db_read_path(config->db, name, tag, &path[0])) {
        usys_log_debug("Path found in db. name:%s tag:%s path:%s",
                       name, tag, path[0]);
        ulfius_set_string_body_response(response, HttpStatus_OK, path);
        goto cleanup;
    }
#endif

    /* Check with hub */
    if (!get_artifact_info_from_hub(&artifact,
                                    config,
                                    name, tag,
                                    &httpStatus)) {
        if (httpStatus == HttpStatus_InternalServerError) {
            usys_log_error("Unable to connect with hub at: %s",
                           config->hubURL);
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

        goto cleanup;
    } else {
        usys_log_debug("capp %s:%s is available at hub: %s",
                       name, tag, config->hubURL);
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
                       name, tag);
        ulfius_set_string_body_response(response,
                                HttpStatus_ServiceUnavailable,
                                HttpStatusStr(HttpStatus_ServiceUnavailable));
        goto cleanup;
    }

    /* create callback URL */
    sprintf(cbURL, "http://localhost:%s/%s",
            config->servicePort,
            WIMC_EP_AGENT_UPDATE);

    /* create URLs for agent to fetch the artifacts from */
    if (strcmp(agent->method, WIMC_METHOD_CHUNK_STR) == 0) {
        create_hub_urls_for_agent(config->hubURL,
                                  artifactFormat->url,
                                  &url[0],
                                  artifactFormat->extraInfo,
                                  &extraURL[0]);
    }

    /* create request */
    request = create_wimc_request(name, tag, config->hubURL, cbURL,
                                  url, extraURL, artifactFormat->type,
                                  DEFAULT_INTERVAL);

    /* Send the request to agent */
    resCode = communicate_with_agent(request, agent->url, config, &uuid);
    if (resCode == HttpStatus_OK) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_Accepted,
                                        HttpStatusStr(HttpStatus_Accepted));
        goto cleanup;
    } else {
        usys_log_debug("No matching agent/resource found for capp %s:%s",
                       name, tag);
    
        ulfius_set_string_body_response(response,
                               HttpStatus_InternalServerError,
                               HttpStatusStr(HttpStatus_InternalServerError));
    }

cleanup:
    usys_free(respBody);
    usys_free(errStr);

    return U_CALLBACK_CONTINUE;
}

#if 0
int callback_put_agent_update(const struct _u_request *request,
			      struct _u_response *response,
			      void *user_data) {

    int ret=WIMC_OK, retCode;
    char *resBody;
    json_t *jreq=NULL;
    json_error_t jerr;
    AgentReq *req=NULL;

    Config *config = NULL;


    config = (Config *)user_data;

    req = (AgentReq *)calloc(sizeof(AgentReq), 1);

    jreq = ulfius_get_json_body_request(request, &jerr);
    if (!jreq) {
        log_error("json error: %s", jerr.text);
    } else {
        deserialize_agent_request(&req, jreq);
    }

    ret = process_agent_update_request(cfg->tasks, req, &uuid, cfg->db);
    
    if (ret == WIMC_OK) {
        retCode = 200;
    } else if (ret == WIMC_ERROR_BAD_ID){
        retCode = 404;
    } else {
        retCode = 400;
    }

    resBody = msprintf("%s\n", error_to_str(ret));
    ulfius_set_string_body_response(response, retCode, resBody);
    o_free(resBody);
    free_agent_request(req);
    json_decref(jreq);
    
    return U_CALLBACK_CONTINUE;
}

#endif

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

    int ret=WIMC_OK, retCode;
    uuid_t uuid;

    char *agentID  = NULL;
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

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_Unauthorized,
                                    HttpStatusStr(HttpStatus_Unauthorized));

    return U_CALLBACK_CONTINUE;
}
