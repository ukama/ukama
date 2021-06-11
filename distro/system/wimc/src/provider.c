/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Functions to interact with the service provider in the cloud.
 *
 */

#include <sqlite3.h>
#include <jansson.h>

#include "wimc.h"
#include "log.h"
#include "provider.h"
#include "ulfius.h"

static void log_request(req_t *req) {

}

static void log_response(resp_t *resp) {

}

/*
 * Create EP from the passed arguments and cmd option.
 *
 */

static char *create_url_ep(char *url, char *name, char *tag, int cmd, int type) {

  int ret;
  char *ep = NULL;

  /* Sanity check. */
  if (name == NULL || strlen(name) >= WIMC_MAX_NAME_LEN ||
      (cmd == WIMC_CMD_ALL_TAGS && tag != NULL)) {
    return NULL;
  }

  ep = (char *)calloc(WIMC_MAX_EP_LEN, 1);
  
  if (ep == NULL) {
    return NULL;
  }
  
  /* EP is of following format: url/type/name/tag/cmd */
 
  switch(type) {
    
  case WIMC_TYPE_CONTAINER:
    sprintf(ep, "%s/%s", url, WIMC_TYPE_CONTAINER_STR);
    break;
    
  case WIMC_TYPE_DATA:
    sprintf(ep, "%s/%s", url, WIMC_TYPE_DATA_STR);
    break;
    
  default:
    /* invalid type. */
    goto failure;
  }

  /* Append name and tag */
  if (tag) {
    sprintf(ep, "%s/%s/%s", ep, name, tag);
  } else {
    /* for containers, if tag is missing the default is 'latest' */
    if (type == WIMC_TYPE_CONTAINER) {
      sprintf(ep, "%s/%s/latest", ep, name);
    }
  }

  switch(cmd) {

  case WIMC_CMD_TRANSFER:
    sprintf(ep, "%s/%s", ep, WIMC_CMD_TRANSFER_STR);
    break;
    
  case WIMC_CMD_INFO:
    sprintf(ep, "%s/%s", ep, WIMC_CMD_INFO_STR);
    break;
    
  case WIMC_CMD_INSPECT:
    sprintf(ep, "%s/%s", ep, WIMC_CMD_INSPECT_STR);
    break;

  case WIMC_CMD_ALL_TAGS:
    sprintf(ep, "%s/%s", ep, WIMC_CMD_ALL_TAGS_STR);
    break;

  default:
    if (type == WIMC_TYPE_CONTAINER){
      sprintf(ep, "%s/%s", ep, WIMC_CMD_TRANSFER_STR);
    }
  }
      
  return ep;
  
 failure:
  free(ep);
  return NULL;
}

/*
 * process_response -- 
 *
 */
static void *process_response(WimcCfg *cfg, resp_t *resp, int cmd) {

  /* Depending on the cmd line, we expect different data from server. 
   * Data is always returned as JSON. 
   */
  
  json_t *jresp = NULL;
  AgentCB **agent = NULL;

  int count, i;
  
  /* transfer - JSON with URLs for the agent(s).
   * info     - return the latest tag for the container.
   * inspect  - return JSON information regarding the container. (OCI/chunk)
   * all-tags - return all tags for the container.
   */
  
  jresp = ulfius_get_json_body_response(resp, NULL);

  if (jresp != NULL) {
    
    log_debug("JSON string: %s \n", json_dumps(jresp, JSON_ENCODE_ANY));

    /* If response is HTTP_RESPONSE_OK, de-seralize JSON response. */
    if (cmd == WIMC_CMD_TRANSFER) {
      
      deserialize_transfer_response(jresp, &agent[0], &count);

      for (i=0; i<count; i++) {
	log_debug("%s: type: %s URL: %s", i, agent[i]->method, agent[i]->url);
      }
      
    } else if (cmd == WIMC_CMD_INFO) {
      
    } else if (cmd == WIMC_CMD_INSPECT) {
      
    } else if (cmd == WIMC_CMD_ALL_TAGS) {

    }
  }

}
/*
 * send_get_request -- send 'cmd' command to the provider.
 *
 */

static int send_request(WimcCfg *cfg, char *name, char *tag, int cmd,
			int type) {
  
  int ret;
  req_t *req=NULL;
  resp_t *resp=NULL;
  char *ep=NULL;

  ep = create_url_ep(cfg->cloud, name, tag, cmd, type);

  if (ep == NULL) {
    log_error("Error creating EP for URL. URL: %s, name: %s, tag: %s, cmd:%d",
	      cfg->cloud, name, tag, cmd);
    return FALSE;
  }
  
  /* create request call. */
  req  = create_http_request(cfg->cloud, ep, WIMC_METHOD_TYPE_GET);
  if (req) {
    log_request(req);
  }
  
  /* Send request call and receive response */
  resp = send_http_request(req); 
  
  if (resp) {
    log_response(resp);
    process_response(cfg, resp, cmd);
  }

  return TRUE;
}

/*
 * fetch_container_from_service_provider -- 
 *
 */ 

 int fetch_content_from_service_provider(WimcCfg *cfg, char *name, char *tag,
					 int type){

  /* Logic is as follows:
   * 1. Issue GET command to the cloud-based service provider for name:tag
   * 2. Provider will either:
   *    a. reject the request with 404 or
   *    b. accept and return remote_cb URL. 
   * 3. Remote_cb URL is then passed to the Agent, along with status_CB.
   * 4. status_cb keep the db updated for the content.
   */

   if (type == WIMC_TYPE_CONTAINER) {
     send_request(cfg, name, tag, WIMC_CMD_TRANSFER, type); /* Issue GET */
   }
}
