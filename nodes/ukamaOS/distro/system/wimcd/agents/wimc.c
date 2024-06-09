/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/* Functions related to wimc. */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <jansson.h>
#include <curl/curl.h>

#include "wimc.h"
#include "err.h"
#include "http_status.h"
#include "common/utils.h"
#include "agent/jserdes.h"

#include "usys_types.h"
#include "usys_log.h"

#define AGENT_CB_EP "app"
#define WIMC_EP     "v1/agents"

struct Response {
    char *buffer;
    size_t size;
};

static void cleanup_agent_request(AgentReq *request) {

  if (request->reg) {
    Register *reg = request->reg;
    
    if (reg->method)
      free(reg->method);

    if (reg->url)
      free(reg->url);
    
    free(reg);
  }

  if (request->unReg) {
    free(request->unReg);
  }

  if (request->update) {
    if (request->update->voidStr)
      free(request->update->voidStr);
    free(request->update);
  }
  
  free(request);
}

static int get_task_status(TaskStatus state) {

  if (state == (TaskStatus)WSTATUS_PEND) {
    return REQUEST;
  } else if (state == (TaskStatus)WSTATUS_START ||
	     state == (TaskStatus)WSTATUS_RUNNING) {
    return FETCH;
  } else if (state == (TaskStatus)WSTATUS_DONE) {
    return DONE;
  } else if (state == (TaskStatus)WSTATUS_ERROR) {
    return ERR;
  }
}

#if 0
static AgentReq *create_agent_request(ReqType type,
                                      int method,
                                      char *cbURL,
                                      uuid_t *uuid,
                                      TStats *stats) {

    AgentReq   *request=NULL;
    Register   *reg=NULL;
    UnRegister *unreg=NULL;
    Update     *update=NULL;
  
    request = (AgentReq *)calloc(1, sizeof(AgentReq));
    if (request == NULL) {
        usys_log_error("Unable to allocate memory: %ld", sizeof(AgentReq));
        return NULL;
    }

    if (type == (ReqType)REQ_REG) {
    
        reg = (Register *)malloc(sizeof(Register));
        if (reg == NULL) {
            usys_log_error("Unable to allocate memory: %ld", sizeof(Register));
            goto done;
        }

        request->type = REQ_REG;
        reg->method   = strdup(convert_method_to_str(method));
        reg->url      = strdup(cbURL);

        request->reg = reg;
    } else if (type == (ReqType)REQ_UNREG) {

        unreg = (UnRegister *)malloc(sizeof(UnRegister));
        if (unreg == NULL) {
           usys_log_error("Unable to allocate memory: %ld", sizeof(UnRegister));
           goto done;
        }

        request->type = REQ_UNREG;
        uuid_copy(unreg->uuid, *uuid);

        request->unReg = unreg;
  } else if (type == (ReqType)REQ_UPDATE) {

        update = (Update *)calloc(1, sizeof(Update));
        if (update == NULL) {
            usys_log_error("Unable to allocate memory: %ld", sizeof(Update));
            goto done;
        }

        request->type = REQ_UPDATE;
       
        uuid_copy(update->uuid, *uuid);
        update->totalKB       = stats->total_bytes / 1024; /* in kilobytes */
        update->transferKB    = stats->total_bytes / 1024;
        update->transferState = get_task_status(stats->status);

        if (stats->stop == TRUE) {
            update->voidStr = strdup(stats->statusStr);
        } else {
            update->voidStr = strdup("");
        }

        request->update = update;
    }

    return request;
  
done:
    if (reg) {
        usys_free(reg->url);
        usys_free(reg->method);
        usys_free(reg);
    }

    usys_free(unreq);
    usys_free(request);
 
    return NULL;
}
#endif

static size_t response_callback(void *contents, size_t size, size_t nmemb,
				void *userp) {

  size_t realsize = size * nmemb;
  struct Response *response = (struct Response *)userp;

  response->buffer = realloc(response->buffer, response->size + realsize + 1);
  
  if(response->buffer == NULL) {
    log_error("Not enough memory to realloc of size: %s",
	      response->size + realsize + 1);
    return 0;
  }

  memcpy(&(response->buffer[response->size]), contents, realsize);
  response->size += realsize;
  response->buffer[response->size] = 0; /* Null terminate. */
  
  return realsize;
}

static long send_request_to_wimc(int reqType,
                                 char *wimcURL,
                                 json_t *json) {

    long code=0;
    char *jsonStr = NULL;
    CURL *curl = NULL;
    CURLcode res;
    
    struct curl_slist *headers=NULL;
    struct Response response;
  
    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (curl == NULL) return 0;

    response.buffer = malloc(1);
    response.size   = 0;
    if (json) jsonStr = json_dumps(json, 0);
  
    /* Add to the header. */
    headers = curl_slist_append(headers, "Accept: application/json");
    headers = curl_slist_append(headers, "Content-Type: application/json");
    headers = curl_slist_append(headers, "charset: utf-8");

    curl_easy_setopt(curl, CURLOPT_URL, wimcURL);

    if (reqType == REQUEST_UPDATE) {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "PUT");
    } else if (reqType == REQUEST_UNREGISTER) {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "DELETE");
    } else if (reqType == REQUEST_REGISTER) {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "POST");
    }

    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    if (jsonStr) curl_easy_setopt(curl, CURLOPT_POSTFIELDS, jsonStr);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

    curl_easy_setopt(curl, CURLOPT_USERAGENT, "agent/0.1");

    res = curl_easy_perform(curl);
    if (res != CURLE_OK) {
        usys_log_error("Error sending request to WIMC at URL %s: %s",
                       wimcURL,
                       curl_easy_strerror(res));
    } else {
        /* get status code. */
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    }

    usys_free(jsonStr);
    usys_free(response.buffer);

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    curl_global_cleanup();
    
    return code;
}

long communicate_with_wimc(int reqType,
                           char *wimcURL,
                           char *cbURL,
                           void *data) {

    int ret;
    long code=0;

    json_t   *json=NULL;
    TStats   *stats=NULL;
    char     *method = NULL;

    if (reqType == REQUEST_UPDATE) {
        stats = (TStats *)data;
    } else if (reqType == REQUEST_REGISTER) {
        method = (char *)data;

        if (!serialize_agent_register_request(method, cbURL, &json)) {
            usys_log_error("Unable to serialize register request");
            return USYS_FALSE;
        }
    }

    code = send_request_to_wimc(reqType, wimcURL, json);
    if (code == HttpStatus_OK) {
        usys_log_debug("Communication with wimc: %s code: %d", 
                       wimcURL, code);
    } else {
        usys_log_error("Communication with WIMC %s: failed. Code: %d",
                       wimcURL, code);
    }

    if (json) json_decref(json);
    
    return code;
}

int register_agent_with_wimc(char *url,
                             char *agentMethod,
                             uuid_t uuid) {

    char idStr[36+1]               = {0};
    char cbURL[WIMC_MAX_URL_LEN]   = {0};
    char wimcURL[WIMC_MAX_URL_LEN] = {0};

    uuid_unparse(uuid, &idStr[0]);

    /* setup cbURL and wimcURL EP */
    sprintf(cbURL, "http://localhost:%d/v1/%s",
            usys_find_service_port(SERVICE_WIMC_AGENT),
            AGENT_CB_EP);
    sprintf(wimcURL, "%s/%s/%s", url, WIMC_EP, idStr);

    if (communicate_with_wimc(REQUEST_REGISTER,
                              wimcURL,
                              cbURL,
                              agentMethod) != HttpStatus_OK) {
        usys_log_error("Error registering agent with WIMC");
        return USYS_FALSE;
    }

    usys_log_debug("Agent registerd");
    return USYS_TRUE;
}

int unregister_agent_with_wimc(char *url,
                               uuid_t uuid) {

    char idStr[36+1]               = {0};
    char wimcURL[WIMC_MAX_URL_LEN] = {0};

    uuid_unparse(uuid, &idStr[0]);
    sprintf(wimcURL, "%s/%s/%s", url, WIMC_EP, idStr);

    if (communicate_with_wimc(REQUEST_UNREGISTER,
                              wimcURL,
                              NULL,
                              NULL) != HttpStatus_OK) {
        usys_log_error("Error un-registering agent with WIMC");
        return USYS_FALSE;
    }

    usys_log_debug("Agent de-registerd");
    return USYS_TRUE;
}
