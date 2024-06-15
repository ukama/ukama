/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/*
 * Agent related functions.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <curl/curl.h>
#include <curl/easy.h>

#include "agent.h"
#include "wimc.h"
#include "tasks.h"
#include "jserdes.h"
#include "common/utils.h"
#include "http_status.h"

#include "usys_types.h"
#include "usys_mem.h"
#include "usys_log.h"

struct Response {
    char *buffer;
    size_t size;
};

bool register_agent(Agent **agents,
                    char *agentID,
                    char *agentMethod,
                    char *agentURL) {

    int i;
    Agent *ptr = NULL;
    uuid_t uuid;
    
    uuid_parse(agentID, uuid);
    ptr = *agents;

    for (i=0; i < MAX_AGENTS; i++) {
        if (uuid_is_null(ptr[i].uuid) == 0) { /* have valid agent id. */
            if (strcasecmp(agentMethod, ptr[i].method) == 0 &&
                strcasecmp(agentURL, ptr[i].url) == 0) {
                /* An existing entry. */
                log_debug("Found similar agent. id: %s, method %s and url: %s",
                          agentID, agentMethod, agentURL);
                return USYS_FALSE;
            }
        } else {

            ptr[i].method = strdup(agentMethod);
            ptr[i].url    = strdup(agentURL);
            ptr[i].state  = WIMC_AGENT_STATE_REGISTER;

            uuid_copy(ptr[i].uuid, uuid);
            
            return USYS_TRUE;
        }
    }

    usys_log_debug("Max. allowable number of agents reached. Ignoring");
    return USYS_FALSE;
}

bool delete_agent(Agent **agents,
                  char *agentID) {

    Agent *ptr = NULL;
    uuid_t uuid;

    uuid_parse(agentID, uuid);
    ptr = *agents;

    for (int i=0; i < MAX_AGENTS; i++) {
        if (uuid_is_null(ptr[i].uuid) == 0) { /* have valid agent id. */
            if (uuid_compare(uuid, ptr[i].uuid) == 0) { /* match */

                usys_log_debug("Agent removed id: %s, method %s and url: %s",
                               agentID, ptr[i].method, ptr[i].url);

                uuid_clear(ptr[i].uuid);
                usys_free(ptr[i].method);
                usys_free(ptr[i].url);
                ptr[i].state = 0;

                return USYS_TRUE;
            }
        }
    }

    usys_log_debug("Agent with UUID not found: %s", agentID);
    return USYS_FALSE;
}

int process_agent_update_request(WTasks **tasks,
                                 AgentReq *req,
                                 sqlite3 *db) {

  int ret=WIMC_OK;
  Update *update;
  char idStr1[36+1] = {0};
  char idStr2[36+1] = {0};
  WTasks *task=NULL;
  
  if (*tasks == NULL) return HttpStatus_InternalServerError;

  update = req->update;
  task = *tasks;

  uuid_unparse(update->uuid, &idStr1[0]);
  usys_log_debug("Looking up task with ID: %s", idStr1);

  /* Find matching task in our list. */
  while (task != NULL) {
      uuid_unparse(task->uuid, &idStr2[0]);
      if (uuid_compare(task->uuid, update->uuid) == 0) {
          usys_log_debug("Found. Ask: %s Match: %s", idStr1, idStr2);
          break;
      }
      task = task->next;
  }

  if (task == NULL) {
      usys_log_error("No record found for ID: %s", idStr1);
      return HttpStatus_BadRequest;
  }

  /* update the task entry. */
  task->update->totalKB = req->update->totalKB;
  task->update->transferKB = req->update->transferKB;

  /* Update the status */
  task->update->transferState = req->update->transferState;
  task->state = req->update->transferState;
  if (req->update->voidStr) {
      task->update->voidStr = strdup(req->update->voidStr);
  }

  if (task->state == DONE) {
      task->localPath = strdup(req->update->voidStr);
      update_local_db(db, task->content->name, task->content->tag,
                      task->localPath);
  }

  return HttpStatus_OK;
}

Agent *find_matching_agent(Agent *agents, char *method) {

    int i;
    Agent *ptr = agents;

    /* Sanity check. */
    if (agents == NULL)
        return NULL;

    for (i=0; i < MAX_AGENTS; i++) {
        if (uuid_is_null(ptr[i].uuid)==0) { /* have valid agent id. */
            if (strcmp(method, ptr[i].method)==0) {
                return ptr;
            }
        }
    }

    /* NULL if no match found. */
    return NULL;
}

/*
 * cleanup_wimc_request --
 *
 */
void cleanup_wimc_request(WimcReq *request) {

    if (request == NULL) return;

    if (request->fetch) {
        if (request->fetch->content) {

            WFetch   *fetch   = request->fetch;
            WContent *content = fetch->content;

            usys_free(content->name);
            usys_free(content->tag);
            usys_free(content->method);
            usys_free(content->indexURL);
            usys_free(content->storeURL);
            usys_free(content);
        }

        usys_free(request->fetch);
    }

    usys_free(request);
}

/*
 * response_callback --
 */
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

void create_wimc_request(WimcReq **request,
                         char *name, char *tag,
                         char *indexURL,
                         char *storeURL,
                         char *method,
                         int interval) {

  WFetch   *fetch=NULL;
  WContent *content=NULL;

  *request = (WimcReq *) calloc(1, sizeof(WimcReq));
  fetch    = (WFetch *)  calloc(1, sizeof(WFetch));
  content  = (WContent *)calloc(1, sizeof(WContent));
  
  if (*request == NULL || fetch == NULL || content == NULL) {
      usys_free(*request);
      usys_free(fetch);
      usys_free(content);

      return NULL;
  }

  (*request)->type = WREQ_FETCH;
  uuid_generate(fetch->uuid);
  fetch->interval = interval;

  content->name        = strdup(name);
  content->tag         = strdup(tag);
  content->method      = strdup(method);
  content->indexURL    = strdup(indexURL);
  content->storeURL    = strdup(storeURL);
  
  fetch->content    = content;
  (*request)->fetch = fetch;
}

static bool send_request_to_agent(char *agentURL,
                                  json_t *json,
                                  int *statusCode) {

    bool  ret  = USYS_FALSE;
    CURL *curl = NULL;
    char *jStr = NULL;

    struct curl_slist *headers=NULL;
    struct Response response;

    *statusCode = 0;
  
    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (curl == NULL) {
        usys_log_error("Error initializing curl");
        return USYS_FALSE;
    }

    response.buffer = malloc(1);
    response.size   = 0;
    jStr = json_dumps(json, 0);

    /* Add to the header. */
    headers = curl_slist_append(headers, "Accept: application/json");
    headers = curl_slist_append(headers, "Content-Type: application/json");
    headers = curl_slist_append(headers, "charset: utf-8");

    curl_easy_setopt(curl, CURLOPT_URL, agentURL);

    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "POST");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER,    headers);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS,    jStr);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA,     (void *)&response);

    curl_easy_setopt(curl, CURLOPT_USERAGENT, "wimc/0.1");

    if (curl_easy_perform(curl) != CURLE_OK) {
        log_error("Error sending request to Agent: %s");
    } else {
        *statusCode = HttpStatus_OK;
        ret = USYS_TRUE;
    }

    usys_free(jStr);
    usys_free(response.buffer);

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    curl_global_cleanup();

    return ret;
}

bool communicate_with_agent(WimcReq *request,
                            char *url,
                            Config *config,
                            uuid_t *uuid) {

    long code=0;
    json_t *json=NULL;
    int agentRetCode=0;

    if (!serialize_wimc_request(request, &json)) {
        usys_log_error("Error serializing wimc request to agent");
        return USYS_FALSE;
    }

    add_to_tasks(config->tasks, request);
    uuid_copy(*uuid, request->fetch->uuid);

    if (send_request_to_agent(url, json, &agentRetCode)) {
        if (agentRetCode == HttpStatus_OK) {
            usys_log_debug("Agent iniated to fetch capp");
        } else {
            usys_log_error("Agent reutrned an error: %d", agentRetCode);
            json_decref(json);

            return USYS_FALSE;
        }
    } else {
        usys_log_error("Error communicating with Agent");
        json_decref(json);

        return USYS_FALSE;
    }

    json_decref(json);
    return USYS_TRUE;
}

void clear_agents(Agent *agent) {

  int i;

  for (i=0; i<MAX_AGENTS; i++){
    if (uuid_is_null(agent[i].uuid)==0) {
      free(agent[i].method);
      free(agent[i].url);
    }
  }
}
