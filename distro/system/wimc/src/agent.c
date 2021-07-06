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
#include "wimc.h"

/*
 * register_agent -- register new agent
 */

int register_agent(Agent **agents, char *method, char *url, int *id) {

  int i;
  Agent *ptr = *agents;

  for (i=0; i < MAX_AGENTS; i++) {
  
    if (ptr[i].id) { /* have valid agent id. */
      if (strcmp(method, ptr[i].method)==0 &&
	  strcmp(url, ptr[i].url)==0) {
	*id = ptr[i].id;
	/* An existing entry. */
	log_debug("Found similar agent at id: %d, method %s and url: %s",
		  ptr[i].id, ptr[i].method, ptr[i].url);
	return WIMC_ERROR_EXIST;
      }
    } else {
      ptr[i].id     = i+1;  /* XXX, need better way to create ID. */
      ptr[i].method = strndup(method, strlen(method));
      ptr[i].url    = strndup(url, strlen(url));
      ptr[i].state  = WIMC_AGENT_STATE_REGISTER;

      /* Return the ID. */
      *id = i+1;

      return WIMC_OK;
    }
  }

  /* Max. reached */
  log_debug("Max. allowable number of agents reached. Ignoring");
  return WIMC_ERROR_MAX_AGENTS;
}

/*
 * process_agent_request --
 *
 */

int process_agent_request(Agent **agents, AgentReq *req, int *id){

  int ret=WIMC_OK;
  Register *reg;
  
  if (req->type == (ReqType)REQ_REG) {

    reg = req->reg;
    
    /* validate the URL. */
    ret = validate_url(reg->url);
    if (ret != WIMC_OK) {
      log_debug("Agent process failed, unreachable URL: %s: %s", reg->url,
		error_to_str(ret));
      goto done;
    }
    
    ret = register_agent(agents, reg->method, reg->url, id);
    if (ret != WIMC_OK) {
      goto done;
    }

    log_debug("Agent successfully registered. Id: %d Method: %s URL: %s",
	      *id, reg->method, reg->url);
  } else if (req->type == (ReqType)REQ_UNREG) {
    
  } else if (req->type == (ReqType)REQ_UPDATE) {

  } else {
    log_debug("Invalid Agent request command: %d", req->type);
    ret = WIMC_ERROR_BAD_METHOD;
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
#if 0
Agent *find_matching_agent(char *method, Agent *agents) {

  Agent *curr;

  for (curr=agents; curr!=NULL; curr=curr->next) {
    if (curr->id) {
      if (strcmp(curr->method, method)==0 &&
	  curr->state != WIMC_AGENT_STATE_UNREGISTER) {
	return curr;
      }
    }
  }

  return NULL;
}
#endif
/*
 * send_agent_request -- send Agent the request to do something. Currently,
 *                       they fetch data from provider using the specified 
 *                       method and send us status update on provided CB URL.
 */
int send_request_to_agent(char *method, Agent *agents) {

  Agent *dest=NULL;
#if 0
  dest = find_matching_agent(method, agents);

  if (!dest) {
    log_debug("No matching agent found for method: %s", method);
    /* currently we ignore this request. In future, we might want to 
     * cache the request and retry again.
     */
    return FALSE;
  }
#endif
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
