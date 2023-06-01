/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <curl/curl.h>
#include <curl/easy.h>
#include <jansson.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <errno.h>

#include "initClient.h"
#include "httpStatus.h"
#include "jserdes.h"
#include "config.h"
#include "log.h"

/* Functions related to communicate with init system */

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

/*
 * send_http_request --
 *
 */
static long send_http_request(char *url, Request *request, json_t *json,
							  char **retStr) {

	long code=0;
	CURL *curl=NULL;
	CURLcode res;
	char *json_str=NULL;
	struct curl_slist *headers=NULL;
	struct Response response;

	/* sanity check */
	if (url == NULL) {
		return FALSE;
	}

	curl_global_init(CURL_GLOBAL_ALL);
	curl = curl_easy_init();
	if (curl == NULL) {
		return FALSE;
	}

	response.buffer = malloc(1);
	response.size   = 0;

	/* Add to the header. */
	headers = curl_slist_append(headers, "Accept: application/json");
	headers = curl_slist_append(headers, "Content-Type: application/json");
	headers = curl_slist_append(headers, "charset: utf-8");

	curl_easy_setopt(curl, CURLOPT_URL, url);

	if (request->reqType == (ReqType)REQ_REGISTER) {
		json_str = json_dumps(json, 0);
		curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "PUT");
		curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json_str);
	} else if (request->reqType == (ReqType)REQ_UNREGISTER) {
		curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "DELETE");
	} else if (request->reqType == (ReqType)REQ_QUERY ||
			   request->reqType == (ReqType)REQ_QUERY_SYSTEM) {
		curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");
	} else if (request->reqType == (ReqType)REQ_UPDATE) {
		json_str = json_dumps(json, 0);
		curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "PATCH");
		curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json_str);
	}

	curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

	curl_easy_setopt(curl, CURLOPT_USERAGENT, "initClient/0.1");

	res = curl_easy_perform(curl);

	if (res != CURLE_OK) {
		log_error("Error sending request to init system at url %s: %s", url,
				  curl_easy_strerror(res));
	} else {
		/* get status code */
		curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
		*retStr = strdup(response.buffer);
	}

	free(json_str);
	free(response.buffer);
	curl_slist_free_all(headers);
	curl_easy_cleanup(curl);
	curl_global_cleanup();

	return code;
}

/*
 * create_url --
 *
 */
static void create_url(char *url, Config *config, char *name,
					   ReqType reqType) {

	char *systemName=NULL;

	if (reqType == (ReqType)REQ_REGISTER ||
		reqType == (ReqType)REQ_UNREGISTER ||
		reqType == (ReqType)REQ_QUERY) {
		systemName = config->systemName;
	} else if (reqType == (ReqType)REQ_QUERY_SYSTEM && name) {
		systemName = name;
	} else {
		systemName = "";
	}

	/* URL -> host:port/v1/orgs/{org}/systems/{system} */
	sprintf(url, "http://%s:%s/%s/%s/%s/%s/%s",
			config->initSystemAddr,
			config->initSystemPort,
			config->initSystemAPIVer,
			ORGS_STR, config->systemOrg,
			SYSTEMS_STR, systemName);

	log_debug("Request URL: %s", url);
}

/*
 * create_request --
 *
 */
static int create_request(Request **request, Config *config) {

	Register *reg=NULL;

	if ((*request)->reqType == (ReqType)REQ_REGISTER ||
		((*request)->reqType == (ReqType)REQ_UPDATE)) {

		reg = (Register *)calloc(1, sizeof(Register));
		if (reg == NULL) return FALSE;

		reg->org  = strdup(config->systemOrg);
		reg->name = strdup(config->systemName);
		reg->ip   = strdup(config->systemAddr);
		reg->port = strdup(config->systemPort);
		reg->cert = strdup(config->systemCert);

		(*request)->reg = reg;
	}

	return TRUE;
}

/*
 * free_request --
 *
 */
static void free_request(Request *request) {

	Register *reg=NULL;

	if (request == NULL) return;

	reg = request->reg;

	if (request->reqType == (ReqType) REQ_REGISTER ||
		request->reqType == (ReqType) REQ_UPDATE) {

		if (reg == NULL) return;

		if (reg->org)  free(reg->org);
		if (reg->name) free(reg->name);
		if (reg->cert) free(reg->cert);
		if (reg->ip)   free(reg->ip);
		if (reg->port) free(reg->port);

		free(reg);
	}

	free(request);
}

/*
 * free_query_response --
 *
 */
void free_query_response(QueryResponse *response) {

	if (response == NULL) return;

	if (response->systemName)  free(response->systemName);
	if (response->systemID)    free(response->systemID);
	if (response->certificate) free(response->certificate);
	if (response->ip)          free(response->ip);

	free(response);
}

/*
 * read_cache_uuid --
 *
 */
static int read_cache_uuid(char *fileName, char **uuid) {

	FILE *fp;
	struct stat sb;
	char buffer[MAX_BUFFER_SIZE] = {0};

	/* Check to see if the cache file exist. */
	if (stat(fileName, &sb) == -1) {
		log_debug("Cache file does not exist: %s Error: %s",
				  fileName, strerror(errno));
		return REG_STATUS_NO_UUID;
	}

	/* Try to open it */
	fp = fopen(fileName, "r");
	if (fp == NULL) {
		log_error("Error opening cache file: %s Error: %s",
				  fileName, strerror(errno));
		return REG_STATUS_NO_UUID;
	}

	/* Try to read the uuid */
	if (fread(buffer, 1, MAX_UUID_LEN, fp) == 0) {
		log_error("Error reading from the cache file: %s Error :%s",
				  fileName, strerror(errno));
		return REG_STATUS_NO_UUID;
	}

	*uuid = strndup(buffer, MAX_UUID_LEN);
	return REG_STATUS_HAVE_UUID;
}

/*
 * send_request_to_init --
 *
 * create_request
 * serialize
 * send to init
 *
 */
int send_request_to_init(ReqType reqType, Config *config,
						 char *systemName, char **response) {

	Request *request=NULL;
	json_t *json=NULL;
	char url[MAX_URL_LEN]={0};
	long respCode;
	int ret=FALSE;

	if (config == NULL) return FALSE;

	request = (Request *)calloc(1, sizeof(Request));
	if (request == NULL) {
		log_error("Error allocating memory of size: %d", sizeof(Request));
		return FALSE;
	}

	request->reqType = reqType;

	/* Step-1 create request */
	if (!create_request(&request, config)) {
		free(request);
		return FALSE;
	}

	/* Step-2 serialize the request */
	if (!serialize_request(request, &json)) {
		log_error("Unable to serialize the request for init");
		json_decref(json);
		free(request);
		return FALSE;
	}

	/* Step-3 create URL for init system */
	create_url(&url[0], config, systemName, reqType);

	/* Step-3 send over the wire */
	respCode = send_http_request(&url[0], request, json, response);

	switch(respCode) {
	case HttpStatus_OK:
		if (reqType == (ReqType)REQ_UNREGISTER) {
			log_debug("Successful unregister");
			ret = TRUE;
		} else if (reqType == (ReqType)REQ_QUERY ||
				   reqType == (ReqType)REQ_QUERY_SYSTEM) {
			log_debug("Query successful");
			ret = TRUE;
		} else if (reqType == (ReqType)REQ_UPDATE) {
			log_debug("Update successful");
			ret = TRUE;
		}
		break;
	case HttpStatus_Created:
		if (reqType == (ReqType)REQ_REGISTER) {
			log_debug("Successful register");
			ret = TRUE;
		}
		break;
	case HttpStatus_BadRequest:
		if (reqType == (ReqType)REQ_QUERY_SYSTEM) {
			log_debug("Invalid system name: %s Response code: %s", systemName,
					  HttpStatusStr(respCode));
			ret = FALSE;
		}
		break;
	default:
		log_error("Error sending request to init: %s", HttpStatusStr(respCode));
		ret=FALSE;
	}

	free_request(request);
	json_decref(json);

	return ret;
}

/*
 * existing_registration --
 *
 */
int existing_registration(Config *config, char **cacheUUID,
						  char **systemUUID) {

	int status=REG_STATUS_NONE;
	char *str=NULL;
	QueryResponse *queryResponse=NULL;

	if (send_request_to_init(REQ_QUERY, config, NULL, &str)) {
		if (deserialize_response(REQ_QUERY, &queryResponse, str) != TRUE) {
			log_error("Error deserialize query response. Str: %s", str);
			return -1;
		}
	} else {
		status = REG_STATUS_NO_MATCH;
		goto return_function;
	}

	status = read_cache_uuid(config->tempFile, cacheUUID);

	/* match? */
	if (strcmp(config->systemName, queryResponse->systemName) == 0 &&
		strcmp(config->systemAddr, queryResponse->ip) == 0 &&
		strcmp(config->systemCert, queryResponse->certificate) == 0 &&
		atoi(config->systemPort) == queryResponse->port) {

		if (status == REG_STATUS_HAVE_UUID) {
			if (strcmp(*cacheUUID, queryResponse->systemID) == 0){
				status |= REG_STATUS_MATCH;
			} else {
				status |= REG_STATUS_NO_MATCH;
			}
		} else {
			status |= REG_STATUS_MATCH;
		}
	} else {
		status |= REG_STATUS_NO_MATCH;
	}

	if (queryResponse->systemID) {
		*systemUUID = strdup(queryResponse->systemID);
	}

 return_function:
	if (str)  free(str);
	free_query_response(queryResponse);

	return status;
}

/*
 * get_system_info -- get info about 'system' from the init.
 *
 */
int get_system_info(Config *config, char *systemName, char **systemInfo) {

	int status=QUERY_OK;
	char *str=NULL;
	QueryResponse *queryResponse=NULL;

	if (send_request_to_init(REQ_QUERY_SYSTEM, config, systemName, &str)) {
		if (deserialize_response(REQ_QUERY, &queryResponse, str) != TRUE) {
			free(str);
			log_error("Error deserialize query response. Str: %s", str);
			return -1;
		}

		*systemInfo = strdup(str);
	} else {
		status = QUERY_ERROR;
	}

	if (str)  free(str);
	free_query_response(queryResponse);

	return status;
}
