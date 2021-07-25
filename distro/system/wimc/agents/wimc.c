/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Functions related to wimc. */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <jansson.h>
#include <curl/curl.h>

#include "wimc.h"
#include "err.h"
#include "common/utils.h"
#include "agent/jserdes.h"

#define AGENT_EP "container/"
#define WIMC_EP  "admin/agent/"

struct Response {
  char *buffer;
  size_t size;
};

/* Function def. */
static char *create_cb_url(char *port);
static char *create_wimc_url(char *url);
static void cleanup_agent_request(AgentReq *request);
static AgentReq *create_agent_request(ReqType type, int method, char *cbURL,
				      uuid_t *uuid, TStats *stats);
static size_t response_callback(void *contents, size_t size, size_t nmemb,
				void *userp);
static void process_response_from_wimc(ReqType reqType, long statusCode,
				       void *resp, uuid_t *uuid);
static long send_request_to_wimc(ReqType reqType, char *wimcURL, json_t *json,
				 uuid_t *uuid);

/*
 * create_cb_url --
 */
static char *create_cb_url(char *port) {

  char *cbURL=NULL;
  
  if (port==NULL) {
    return NULL;
  }

  cbURL = (char *)malloc(WIMC_MAX_URL_LEN);
  if (cbURL) {
    sprintf(cbURL, "http://localhost:%s/%s", port, AGENT_EP);
  }

  return cbURL;
}

static char *create_wimc_url(char *url) {

  char *wimcURL=NULL;

  if (!url) {
    return wimcURL;
  }

  wimcURL = (char *)malloc(WIMC_MAX_URL_LEN);
  if (wimcURL) {
    sprintf(wimcURL, "%s/%s", url, WIMC_EP);
  }
	    
  return wimcURL;
}

/* 
 * cleanup_agent_request --
 *
 */
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

/*
 * get_task_status --
 *
 */

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

/*
 * create_agent_request --
 *
 */
static AgentReq *create_agent_request(ReqType type, int method, char *cbURL,
				      uuid_t *uuid, TStats *stats) {

  AgentReq *request=NULL;
  Register *reg=NULL;
  UnRegister *unreg=NULL;
  Update *update=NULL;
  
  request = (AgentReq *)calloc(1, sizeof(AgentReq));
  if (request==NULL) {
    goto done;
  }

  if (type == (ReqType)REQ_REG) {
    
    reg = (Register *)malloc(sizeof(Register));
    if (!reg) {
      goto done;
    }
    
    request->type = REQ_REG;
    
    reg->method = strdup(convert_method_to_str(method));
    reg->url = strdup(cbURL);
    
    /* Sanity check. */
    if (!strlen(reg->method) || reg->url==NULL) {
      goto done;
    }
    
    request->reg = reg;
  } else if (type == (ReqType)REQ_UNREG) {

    unreg = (UnRegister *)malloc(sizeof(UnRegister));
    if (!unreg) {
      goto done;
    }	 

    uuid_copy(unreg->uuid, *uuid);

    request->type = REQ_UNREG;
    request->unReg = unreg;
  } else if (type == (ReqType)REQ_UPDATE) {

    update = (Update *)malloc(sizeof(Update));
    if (!update) {
      goto done;
    }

    uuid_copy(update->uuid, *uuid);
    update->totalKB = stats->total_requests / 1024; /* in kilobytes */
    update->transferKB = stats->total_bytes / 1024;
    update->transferState = get_task_status(stats->status);

    if (update->transferState == WSTATUS_DONE ||  /* content path */
	update->transferState == WSTATUS_ERROR) { /* error str */
      update->voidStr = strdup(stats->statusStr);
    }

    request->type = REQ_UPDATE;
    request->update = update;
  }

  return request;
  
 done:
 if (reg) {
   free(reg->url);
   free(reg->method);
   free(reg);
 }

 if (unreg) {
   free(unreg);
 }
  
 if (request) {
   free(request);
 }


 
 return NULL;
}

/*
 * response_callback --
 *
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

/*
 * process_response_from_wimc --
 *
 */
static void process_response_from_wimc(ReqType reqType, long statusCode,
				       void *resp, uuid_t *uuid) {

  struct Response *response;

  response = (struct Response *)resp;
  
  if (reqType == (ReqType)REQ_REG) {
    if (statusCode == 200) { /* Success, response has ID. */
      uuid_parse(response->buffer, *uuid);
      log_debug("Registration successful. Status code: 200 Recevied ID: %s",
		response->buffer);
    } else if (statusCode == 400) { /* Failure. Report. */
      log_debug("Registration unsuccessful. Status Code: 400 Response: %s",
		response->buffer);
    }
  } else if (reqType == (ReqType)REQ_UNREG) {
    
  }


  return;
}

/*
 * send_request_to_wimc -- 
 *
 */
static long send_request_to_wimc(ReqType reqType, char *wimcURL,
				 json_t *json, uuid_t *uuid) {

  long code=0;
  CURL *curl=NULL;
  CURLcode res;
  char *json_str;
  struct curl_slist *headers=NULL;
  struct Response response;
  
  curl_global_init(CURL_GLOBAL_ALL);
  curl = curl_easy_init();
  if (curl == NULL) {
    return code;
  }

  response.buffer = malloc(1);
  response.size   = 0;
  json_str = json_dumps(json, 0);
  
  /* Add to the header. */
  headers = curl_slist_append(headers, "Accept: application/json");
  headers = curl_slist_append(headers, "Content-Type: application/json");
  headers = curl_slist_append(headers, "charset: utf-8");

  curl_easy_setopt(curl, CURLOPT_URL, wimcURL);

  if (reqType == REQ_UPDATE) {
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "PUT");
  } else {
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "POST");
  }

  curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
  curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json_str);
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

  curl_easy_setopt(curl, CURLOPT_USERAGENT, "agent/0.1");

  res = curl_easy_perform(curl);

  if (res != CURLE_OK) {
    log_error("Error sending request to WIMC at URL %s: %s", wimcURL,
	      curl_easy_strerror(res));
  } else {
    /* get status code. */
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    process_response_from_wimc(reqType, code, &response, uuid);
  }

  free(json_str);
  free(response.buffer);
  curl_slist_free_all(headers);
  curl_easy_cleanup(curl);
  curl_global_cleanup();

  return code;
}

/*
 * communicate_with_wimc -- Function agent uses to communicate with the wimc.d
 *
 */

long communicate_with_wimc(ReqType reqType, char *url, char *port,
			   int method, uuid_t *uuid, void *data) {

  int ret;
  long code=0;
  char *cbURL=NULL, *wimcURL=NULL;
  AgentReq *request=NULL;
  json_t *json=NULL;
  TStats *stats=NULL;

  /* Sanity check. Method can be NULL; only for REQ_REG */
  if (reqType == (ReqType)REQ_UPDATE) {
    if (!url && !data) {
      return code;
    }

    wimcURL = strdup(url);
    stats = (TStats *)data;
  } else if (reqType == (ReqType)REQ_REG) {
    if (!url && !port) {
      return code;
    }

    cbURL   = create_cb_url(port);
    wimcURL = create_wimc_url(url);

    if (!cbURL || !wimcURL) {
      goto done;
    }
  }

  request = create_agent_request(reqType, method, cbURL, uuid, stats);
  if (!request) {
    goto done;
  }

  ret = serialize_agent_request(request, &json);
  if (!ret) {
    goto done;
  }

  code = send_request_to_wimc(reqType, wimcURL, json, uuid);
  if (code == 200) {
    log_debug("WIMC.d %s: success. URL: %s Return code: %d", 
	      convert_type_to_str(reqType), wimcURL, code);
  } else {
    log_error("WIMC.d %s: failed. URL: %s Return code: %d",
	      convert_type_to_str(reqType), wimcURL, code);
  }

 done:

  json_decref(json);
  cleanup_agent_request(request);
  if (cbURL) {
    free(cbURL);
  }

  if (wimcURL) {
    free(wimcURL);
  }

  return code;
}
