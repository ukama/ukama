/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
#include "log.h"

/*
 * agent_error_to_str -- return string representation of the error code.
 *
 */

const char *agent_error_to_str(int error) {

  switch (error) {

  case WIMC_AGENT_OK:
    return WIMC_AGENT_OK_STR;
    
  case WIMC_AGENT_ERROR_EXIST:
    return  WIMC_AGENT_ERROR_EXIST_STR;

  case WIMC_AGENT_ERROR_BAD_METHOD:
    return  WIMC_AGENT_ERROR_BAD_METHOD_STR;

  case WIMC_AGENT_ERROR_BAD_URL:
    return  WIMC_AGENT_ERROR_BAD_URL_STR;

  case WIMC_AGENT_ERROR_MEMORY:
    return  WIMC_AGENT_ERROR_MEMORY_STR;

  default:
    return "";
  }

  return "";
}

/*
 * validate_agent_url -- validate agent URL by doing CURL request. 
 */
static int validate_url(char *url) {

  CURL *curl;
  CURLcode response;

  curl = curl_easy_init();

  if (curl) {
    
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_NOBODY, 1);

    response = curl_easy_perform(curl);

    curl_easy_cleanup(curl);
  }

  if (response != CURLE_OK) {
    return WIMC_AGENT_ERROR_BAD_URL;
  }

  return WIMC_AGENT_OK;
}

/*
 * register_agent -- register new agent
 */

int register_agent(Agent *agents, char *method, char *url) {

  Agent *ptr;

  for (ptr = agents; ptr != NULL; ptr=ptr->next) {

    if (strcmp(method, ptr->method)==0 &&
	strcmp(url, ptr->url)==0) {
      /* An existing entry. */
      return WIMC_AGENT_ERROR_EXIST;
    }
    
    if (ptr->method == NULL && ptr->url == NULL) {
      ptr->method = method;
      ptr->url    = url;

      ptr->work = (AgentWork *)calloc(sizeof(AgentWork), 1);
      if (ptr->work == NULL){
	log_error("Error allocating the memory: %d", sizeof(AgentWork));
	return WIMC_AGENT_ERROR_MEMORY;
      }
	  
      ptr->work->state = WIMC_AGENT_STATE_REGISTER;
    }
  }

  return WIMC_AGENT_OK;
}

/*
 * process_agent_request --
 *
 */

int process_agent_request(Agent *agents, AgentReq *req){

  int ret=WIMC_AGENT_OK;
  Register *reg;
  
  if (req->type == (ReqType)REQ_REG) {

    reg = req->reg;
    
    /* validate the URL. */
    ret = validate_url(reg->url);
    if (ret != WIMC_AGENT_OK) {
      goto done;
    }
    
    ret = register_agent(agents, reg->method, reg->url);
    if (ret != WIMC_AGENT_OK) {
      goto done;
    }

    log_debug("Agent successfully registered. Method: %s URL: %s",
	      reg->method, reg->url);
  } else if (req->type == (ReqType)REQ_UNREG) {
    
  } else if (req->type == (ReqType)REQ_UPDATE) {

  } else {
    log_debug("Invalid Agent request command: %d", req->type);
    ret = WIMC_AGENT_ERROR_BAD_METHOD;
    goto done;
  }
  
 done:
    return ret;
}

/*
 * find_matching_agent -- return the Agent which matches the given
 *                        method. If there are multiple URL in the list
 *                        currently, we always send the first match.
 */

Agent *find_matching_agent(char *method, Agent *agents) {

  Agent *curr;

  for (curr=agents; curr!=NULL; curr=curr->next) {
    if (strcmp(curr->method, method)==0 &&
	curr->work->state != WIMC_AGENT_STATE_UNREGISTER) {
      return curr;
    }
  }

  return NULL;
}

/*
 * send_agent_request -- send Agent the request to do something. Currently,
 *                       they fetch data from provider using the specified 
 *                       method and send us status update on provided CB URL.
 */
int send_request_to_agent(char *method, Agent *agents) {

  Agent *dest=NULL;

  dest = find_matching_agent(method, agents);

  if (!dest) {
    log_debug("No matching agent found for method: %s", method);
    /* currently we ignore this request. In future, we might want to 
     * cache the request and retry again.
     */
    return FALSE;
  }

  /* req: WIMC.d --> Agent
   * req -> { id: "some_id", 
   *          cmd: "fetch",
   *          content: {name: "name", 
   *                     tag: "tag", 
   *                     provider_url: "http://www/www/www"},
   *          callback_url: "http://www.xyz.ccc/cc/cc/", 
   *          update_interval: 10}}
   *        }
   *
   * cmd: fetch, update, cancel
   *
   * updates: Agent --> WIMC.d
   * updates -> { event: "update", 
   *              update: {
   *                      id: "same_id", 
   *                      total_kbytes: "1234"
   *                      transfer_kbytes:  "34"
   *                      state: "fetch"
   *                      void: "some_string_"
   *			}
   *              }	 
   *
   * state: fetch, unpack, done, error
   * void: error -> error string
   *            done  -> path where data is stored.
   *            ""    
   */

}
