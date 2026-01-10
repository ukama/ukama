/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
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

static long send_http_request(char *url, Request *request, json_t *json,
                              char **retStr) {

    long code = 0;
    CURL *curl = NULL;
    CURLcode res;
    char *json_str = NULL;
    struct curl_slist *headers = NULL;
    struct Response response;
    const char *method = "UNKNOWN";

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

    /* Add headers */
    headers = curl_slist_append(headers, "Accept: application/json");
    headers = curl_slist_append(headers, "Content-Type: application/json");
    headers = curl_slist_append(headers, "charset: utf-8");

    curl_easy_setopt(curl, CURLOPT_URL, url);

    if (request->reqType == (ReqType)REQ_REGISTER) {
        method = "PUT";
        json_str = json_dumps(json, 0);
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json_str);

    } else if (request->reqType == (ReqType)REQ_UNREGISTER) {
        method = "DELETE";
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);

    } else if (request->reqType == (ReqType)REQ_QUERY ||
               request->reqType == (ReqType)REQ_QUERY_SYSTEM) {
        method = "GET";
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);

    } else if (request->reqType == (ReqType)REQ_UPDATE) {
        method = "PATCH";
        json_str = json_dumps(json, 0);
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json_str);
    }

    /* ---- LOG REQUEST ---- */
    log_debug("HTTP Request:");
    log_debug("  Method : %s", method);
    log_debug("  URL    : %s", url);
    if (json_str) {
        log_debug("  Body   : %s", json_str);
    }

    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "initClient/0.1");

    res = curl_easy_perform(curl);

    if (res != CURLE_OK) {
        log_error("HTTP request failed: %s",
                  curl_easy_strerror(res));
    } else {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);

        /* ---- LOG RESPONSE ---- */
        log_debug("HTTP Response:");
        log_debug("  Status : %ld", code);
        log_debug("  Body   : %s",
                  response.buffer ? response.buffer : "(null)");

        *retStr = strdup(response.buffer);
    }

    free(json_str);
    free(response.buffer);
    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    curl_global_cleanup();

    return code;
}

static void create_url(char *url, Config *config, char* org, char *name,
					   ReqType reqType, int global) {

	char *systemName=NULL;

	if (reqType == (ReqType)REQ_REGISTER ||
		reqType == (ReqType)REQ_UPDATE ||
		reqType == (ReqType)REQ_UNREGISTER ||
		reqType == (ReqType)REQ_QUERY) {
		systemName = config->systemName;
	} else if (reqType == (ReqType)REQ_QUERY_SYSTEM && name) {
		systemName = name;
	} else {
		systemName = "";
	}

	/* URL -> host:port/v1/orgs/{org}/systems/{system} */
	if (global) {
		sprintf(url, "http://%s:%s/%s/%s/%s/%s/%s",
					config->globalInitSystemAddr,
					config->globalInitSystemPort,
					config->initSystemAPIVer,
					ORGS_STR, org,
					SYSTEMS_STR, systemName);
	} else {
		sprintf(url, "http://%s:%s/%s/%s/%s/%s/%s",
				config->initSystemAddr,
				config->initSystemPort,
				config->initSystemAPIVer,
				ORGS_STR, org,
				SYSTEMS_STR, systemName);
	}
	log_debug("Request URL: %s", url);
}

static int create_request(Request **request, Config *config) {

	Register *reg=NULL;

	if ((*request)->reqType == (ReqType)REQ_REGISTER ||
		((*request)->reqType == (ReqType)REQ_UPDATE)) {

		reg = (Register *)calloc(1, sizeof(Register));
		if (reg == NULL) return FALSE;

		reg->org        = strdup(config->systemOrg);
		reg->name       = strdup(config->systemName);
		reg->apiGwIp    = strdup(config->systemAddr);
		reg->apiGwPort  = strdup(config->systemPort);
		reg->cert       = strdup(config->systemCert);
        reg->nodeGwIp   = strdup(config->systemNodeGwAddr);
        reg->nodeGwPort = strdup(config->systemNodeGwPort);

		(*request)->reg = reg;
	}

	return TRUE;
}

static void free_request(Request *request) {

	Register *reg=NULL;

	if (request == NULL) return;

	reg = request->reg;

	if (request->reqType == (ReqType) REQ_REGISTER ||
		request->reqType == (ReqType) REQ_UPDATE) {

		if (reg == NULL) return;

		if (reg->org)        free(reg->org);
		if (reg->name)       free(reg->name);
		if (reg->cert)       free(reg->cert);
		if (reg->apiGwIp)    free(reg->apiGwIp);
		if (reg->apiGwPort)  free(reg->apiGwPort);
        if (reg->nodeGwIp)   free(reg->nodeGwIp);
        if (reg->nodeGwPort) free(reg->nodeGwPort);

		free(reg);
	}

	free(request);
}

void free_system_registration(SystemRegistrationId* sysReg) {
	if (sysReg == NULL) return;

	if (sysReg->globalUUID) {
		free(sysReg->globalUUID);
	}

	if (sysReg->localUUID) {
		free(sysReg->localUUID);
	}

	free(sysReg);
}

void free_query_response(QueryResponse *response) {

	if (response == NULL) return;

	if (response->systemName)  free(response->systemName);
	if (response->systemID)    free(response->systemID);
	if (response->certificate) free(response->certificate);
	if (response->apiGwIp)     free(response->apiGwIp);
    if (response->nodeGwIp)    free(response->nodeGwIp);

	free(response);
}

int parse_cache_uuid(char *fileName, SystemRegistrationId **sysReg) {

	FILE *fp = NULL;
    struct stat sb;
    char *str = NULL;

	/* Check to see if the cache file exist. */
	if (stat(fileName, &sb) == -1) {
		log_debug("Cache file does not exist: %s Error: %s",
				fileName, strerror(errno));
		return FALSE;
	}

	/* Try to open it */
	fp = fopen(fileName, "r");
	if (fp == NULL) {
		log_error("Error opening cache file: %s Error: %s",
				fileName, strerror(errno));
		return FALSE;
	}

	fseek(fp, 0, SEEK_END);
	long fsize = ftell(fp);
	fseek(fp, 0, SEEK_SET);  /* same as rewind(f); */

	str = malloc(fsize + 1);
	/* Try to read the uuid */
	if (fread(str, 1, fsize, fp) == 0) {
		log_error("Error reading from the cache file: %s Error :%s",
				fileName, strerror(errno));
		return FALSE;
	}
	fclose(fp);

	if (!deserialize_uuids_from_file(str, sysReg)) {
		log_error("Error parsing the cache file: %s Error :%s",
				fileName, strerror(errno));
		return FALSE;
	}

	return TRUE;
}

static int read_cache_uuid(char *fileName, char** uuid, int global) {

	SystemRegistrationId *sysReg = NULL;

	if (parse_cache_uuid(fileName, &sysReg)) {
		if (global && sysReg->globalUUID) {

			*uuid = strdup(sysReg->globalUUID);
            free_system_registration(sysReg);
			return REG_STATUS_HAVE_UUID;
		} else if (sysReg->localUUID){

			*uuid = strdup(sysReg->localUUID);
            free_system_registration(sysReg);
			return REG_STATUS_HAVE_UUID;
		}
	}

	free_system_registration(sysReg);
	return REG_STATUS_NO_UUID;
}

int send_request_to_init(ReqType reqType, Config *config, char* org,
						 char *systemName, char **response, int global ) {

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
	create_url(&url[0], config, org, systemName, reqType, global);

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

int existing_registration(Config *config, char **cacheUUID, char **systemUUID,
		 int global) {

	int status=REG_STATUS_NONE;
	char *str=NULL;
	QueryResponse *queryResponse=NULL;
	if (send_request_to_init(REQ_QUERY, config, config->systemOrg, NULL, &str, global)) {
		if (deserialize_response(REQ_QUERY, &queryResponse, str) != TRUE) {
			log_error("Error deserialize query response. Str: %s", str);
			return -1;
		}
	} else {
		status = REG_STATUS_NO_MATCH;
		goto return_function;
	}

	status = read_cache_uuid(config->tempFile, cacheUUID, global);

	/* match? */
	if (strcmp(config->systemName, queryResponse->systemName) == 0 &&
		strcmp(config->systemAddr, queryResponse->apiGwIp) == 0 &&
		strcmp(config->systemCert, queryResponse->certificate) == 0 &&
		atoi(config->systemPort) == queryResponse->apiGwPort) {

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
	log_info("Returning status 0x%X for %s registration",
             status, (queryResponse->systemID)?queryResponse->systemID:"null");
 return_function:
	if (str)  free(str);
	if (*cacheUUID) free (*cacheUUID);
	free_query_response(queryResponse);
	return status;
}

int get_system_info(Config *config, char *org,
                    char *systemName, char **systemInfo,
                    int global) {

	int status=QUERY_OK;
	char *str=NULL;
	QueryResponse *queryResponse=NULL;

	if (send_request_to_init(REQ_QUERY_SYSTEM, config, org, systemName, &str, global)) {
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
