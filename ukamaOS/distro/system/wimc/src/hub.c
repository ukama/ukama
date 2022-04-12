/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Interact with ukama's hub.
 *
 */

#include <sqlite3.h>
#include <jansson.h>
#include <ulfius.h>
#include <curl/curl.h>
#include <string.h>

#include "wimc.h"
#include "log.h"
#include "hub.h"
#include "jserdes.h"

struct Response {
  char *buffer;
  size_t size;
};

static int copy_artifact(Artifact *src, Artifact *dest);
static void create_hub_url(WimcCfg *cfg, char *url, char *name);
static int process_response_from_hub(Artifact ***artifacts, void *resp);
static size_t response_callback(void *contents, size_t size, size_t nmemb,
				void *userp);

/*
 * process_response_from_hub --
 *
 */
static int process_response_from_hub(Artifact ***artifacts, void *resp) {

  struct Response *response=NULL;
  ArtifactFormat *format=NULL;
  json_t *json=NULL;
  int count=0, i=0, j=0, ret=FALSE;

  response = (struct Response *)resp;

  json = json_loads(response->buffer, JSON_DECODE_ANY, NULL);

  if (!json) {
    log_error("Can not load str into JSON object. Str: %s", response->buffer);
    goto done;
  }

  ret = deserialize_hub_response(artifacts, &count, json);

  if (ret==FALSE) {
    log_error("Deserialization failed for response: %s", response->buffer);
    goto done;
  }

  if (count==0) {
    log_debug("No matching capp available");
    goto done;
  }

  log_debug("Received Artifacts from the hub. %d:", count);

  for (i=0; i<count; i++) {
    log_debug("\n\t Name: %s \n\t Version: %s", (*artifacts)[i]->name,
	      (*artifacts)[i]->version);
    log_debug("\t Formats: %d", (*artifacts)[i]->formatsCount);

    for (j=0; j<(*artifacts)[i]->formatsCount; j++) {
      format = (*artifacts)[i]->formats[j];

      log_debug("\n\t %d:\n \t\t type: %s \n\t\t url: %s \n\t\t createdAt: %s",
		j, format->type, format->url, format->createdAt);
      log_debug("\t\t size: %d", format->size);
      if (format->extraInfo) {
	log_debug("\t\t extra: %s", format->extraInfo);
      }
    }
  }

 done:
  json_decref(json);
  return count;
}

/*
 * create_hub_url --
 *
 */
static void create_hub_url(WimcCfg *cfg, char *name, char *url) {

  if (!cfg || !name || !url) return;

  sprintf(url, "%s/%s/%s", cfg->hubURL, WIMC_EP_HUB_CAPPS, name);

  return;
}

/*
 * copy_artifact --
 *
 */
static int copy_artifact(Artifact *src, Artifact *dest) {

  int i;

  if (src == NULL || dest == NULL) return FALSE;

  dest->name         = strdup(src->name);
  dest->version      = strdup(src->version);
  dest->formatsCount = src->formatsCount;

  dest->formats = (ArtifactFormat **)calloc(src->formatsCount,
					    sizeof(ArtifactFormat *));
  if (dest->formats == NULL) {
    goto failure;
  }

  for (i=0; i<src->formatsCount; i++) {

    dest->formats[i] = (ArtifactFormat *)calloc(1, sizeof(ArtifactFormat));

    dest->formats[i]->type      = strdup(src->formats[i]->type);
    dest->formats[i]->url       = strdup(src->formats[i]->url);
    dest->formats[i]->createdAt = strdup(src->formats[i]->createdAt);
    dest->formats[i]->size      = src->formats[i]->size;

    if (src->formats[i]->extraInfo) {
      dest->formats[i]->extraInfo = strdup(src->formats[i]->extraInfo);
    }
  }

  return TRUE;

 failure:
  if (dest->name)    free(dest->name);
  if (dest->version) free(dest->version);

  return FALSE;
}

/*
 * free_artifact --
 *
 */
void free_artifact(Artifact *artifact) {

  int i;

  if (artifact == NULL) return;

  if (artifact->name)    free(artifact->name);
  if (artifact->version) free(artifact->version);

  for (i=0; i<artifact->formatsCount; i++) {
    if(artifact->formats[i]->type)      free(artifact->formats[i]->type);
    if(artifact->formats[i]->url)       free(artifact->formats[i]->url);
    if(artifact->formats[i]->extraInfo) free(artifact->formats[i]->extraInfo);
    if(artifact->formats[i]->createdAt) free(artifact->formats[i]->createdAt);

    free(artifact->formats[i]);
  }

  if (artifact->formats) free(artifact->formats);

  return;
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

/*
 * get_artifacts_info_from_hub --
 *
 * -1: curl error and curlCode is set.
 *  0: error processing response.
 *  1: Success and curlCode is CURLE_OK
 */
int get_artifacts_info_from_hub(Artifact *artifact, WimcCfg *cfg,
				char *name, char *tag,
				CURLcode *curlCode) {

  int i, ret=TRUE, count=0;
  char hubEP[WIMC_MAX_URL_LEN] = {0};
  CURL *curl=NULL;
  CURLcode res;
  struct Response response;
  Artifact **artifacts=NULL;

  /* Sanity check. */
  if (!name || !tag) {
    ret = FALSE;
    goto done;
  }

  create_hub_url(cfg, name, &hubEP[0]);

  curl = curl_easy_init();
  if (curl == NULL) {
    ret = -1;
    return ret;
  }

  response.buffer = (char *)malloc(1);
  response.size   = 0;

  curl_easy_setopt(curl, CURLOPT_URL, hubEP);

  curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

  curl_easy_setopt(curl, CURLOPT_USERAGENT, "wimc/0.1");

  *curlCode = curl_easy_perform(curl);

  if (*curlCode != CURLE_OK) {
    ret = -1;
    log_error("Error sending request to hub: %s",
	      curl_easy_strerror(*curlCode));
    goto done;
  }

  /* get status code. */
  count = process_response_from_hub(&artifacts, &response);

  if (count == 0) { /* No matching capp found by 'name' */
    ret = FALSE;
    log_debug("No matching capp returned from the hub. Requested: %s tag: %s",
	      name, tag);
    goto done;
  }

  /* Validate the name */
  if (strcmp(artifacts[0]->name, name) != 0) {
    log_error("Got wrong capp. Requested: %s Got %s", name,
	      artifacts[0]->name);
    goto done;
  }

  /* Find matching capp */
  for (i=0; i<count; i++) {
    if (strcmp(artifacts[i]->version, tag)==0) {
      copy_artifact(artifacts[i], artifact);
      break;
    }
  }

 done:
  for (i=0; i<count; i++) {
    free_artifact(artifacts[i]);
    free(artifacts[i]);
  }
  free(artifacts);
  free(response.buffer);
  curl_easy_cleanup(curl);

  return ret;
}

