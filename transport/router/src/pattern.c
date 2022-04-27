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
#include <regex.h>

#include "pattern.h"
#include "router.h"
#include "log.h"

/*
 * match_key_value --
 *
 */
static int match_all_key_value(Pattern *pattern, char *key, char *value) {

  Pattern *ptr=NULL;
  int ret;
  regex_t re;

  for(ptr=pattern; ptr; ptr=ptr->next) {

    if (strcasecmp(ptr->key, key) == 0) {

      /* Special case for asterik-only. If value is '*', it matches
       * with anything but empty strings */
      if (strcmp(value, ASTERIK_ONLY) == 0 &&
	  strlen(value) == 1 &&
	  strlen(ptr->value)) {
	return TRUE;
      }

      if ((ret=regcomp(&re, value, REG_EXTENDED | REG_NOSUB)) != 0) {
	return FALSE;
      }

      if (regexec(&re, ptr->value, 0, NULL, 0) == 0) {
	regfree(&re);
	return TRUE;
      }

      regfree(&re);
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

  Patterns *patterns, *ptmp;
  Pattern *pattern, *tmp;
  Forward *forward;

  if (service == NULL) return;

  if (service->name) free(service->name);

  patterns = service->patterns;
  forward  = service->forward;

  while(patterns) {

    pattern = patterns->pattern;
    while (pattern) {

      if (pattern->key)   free(pattern->key);
      if (pattern->value) free(pattern->value);

      tmp = pattern->next;
      free(pattern);
      pattern = tmp;
    }

    if (patterns->path) free(patterns->path);

    ptmp = patterns->next;
    free(patterns);
    patterns = ptmp;
  }

  if (forward) {
    if (forward->ip)   free(forward->ip);
    free(forward);
  }

  free(service);
  service=NULL;
}

/*
 * find_matching_service --
 *
 */
int find_matching_service(Router *router, Pattern *requestPattern,
			  Forward **forward, char **ep) {

  Service *services=NULL;
  Patterns *patterns=NULL;
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
    for (patterns=services->patterns; patterns; patterns=patterns->next) {

      count = pattern_count(patterns->pattern);

      if (count != requestCount) continue;

      if (service_pattern_match(patterns->pattern, requestPattern)) {
	*forward = (Forward *)calloc(1, sizeof(Forward));
	if (*forward == NULL) {
	  log_error("Error allocating memory of size: %lu", sizeof(Forward));
	  return FALSE;
	}

	(*forward)->ip   = strdup(services->forward->ip);
	(*forward)->port = services->forward->port;
	(*ep) = strdup(patterns->path);

	return TRUE;
      }
    }
  }

  return FALSE;
}
