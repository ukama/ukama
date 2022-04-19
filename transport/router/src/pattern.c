/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Methods related to service router pattern matching
 *
 */

#include <curl/curl.h>
#include <curl/easy.h>
#include <string.h>
#include <strings.h>

#include "router.h"
#include "log.h"

/*
 * match_key_value --
 *
 */
static int match_all_key_value(Pattern *pattern, char *key, char *value) {

  Pattern *ptr=NULL;

  for(ptr=pattern; ptr; ptr=ptr->next) {
    if (strcasecmp(ptr->key, key) == 0 &&
	strcasecmp(ptr->value, value) == 0) {
      return TRUE;
    }
  }

  return FALSE;
}

/*
 * service_pattern_match --
 *
 */
static int service_pattern_match(Pattern *sPattern, Pattern *rPattern) {

  Pattern *sPtr=NULL;

  if (sPattern == NULL || rPattern == NULL) return FALSE;

  for (sPtr=sPattern; sPtr; sPtr=sPtr->next) {
    if (!match_all_key_value(rPattern, sPtr->key, sPtr->value)) {
      return FALSE;
    }
  }

  return TRUE;
}

/*
 * pattern_count --
 *
 */
static int pattern_count(Pattern *pattern) {

  Pattern *ptr=NULL;
  int count=0;

  ptr = pattern;

  while(ptr) {
    count++;
    ptr=ptr->next;
  }

  return count;
}

/*
 * free_service --
 *
 */
void free_service(Service *service) {

  Pattern *ptr=NULL, *tmp=NULL;
  Forward *fPtr=NULL;

  if (service == NULL) return;

  ptr  = service->pattern;
  fPtr = service->forward;

  if (fPtr) {
    if (fPtr->ip)   free(fPtr->ip);
    if (fPtr->port) free(fPtr->port);
    free(fPtr);
  }

  while (ptr) {
    if (ptr->key)   free(ptr->key);
    if (ptr->value) free(ptr->value);
    tmp = ptr->next;
    free(ptr);
    ptr = tmp;
  }

  free(service);
}

/*
 * find_matching_service --
 *
 */
int find_matching_service(Router *router, Pattern *requestPattern,
			  Forward **forward) {

  Service *services=NULL;
  int requestCount=0, count=0;

  /* two basic matching rules:
   *
   * 1. # of k-v pairs must match
   * 2. order doesn't matter
   */

  if (router == NULL || requestPattern == NULL) return FALSE;

  requestCount = pattern_count(requestPattern);
  if (!requestCount) {
    log_info("Requested pattern count: %s", requestCount);
    return FALSE;
  }

  for (services=router->services; services; services=services->next) {

    count = pattern_count(services->pattern);

    if (count != requestCount) continue;

    if (service_pattern_match(services->pattern, requestPattern)) {
      *forward = (Forward *)calloc(1, sizeof(Forward));
      if (*forward == NULL) {
	log_error("Error allocating memory of size: %lu", sizeof(Forward));
	return FALSE;
      }

      (*forward)->ip = strdup(services->forward->ip);
      (*forward)->port = strdup(services->forward->port);

      return TRUE;
    }
  }

  return FALSE;
}
