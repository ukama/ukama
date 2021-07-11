/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Utility functions related to URL. */

#include <curl/curl.h>
#include <curl/easy.h>

#include "wimc.h"
#include "err.h"

/*
 * validate_agent_url -- validate agent URL by doing CURL request. 
 */
int validate_url(char *url) {
  
#if 0
  CURL *curl;
  CURLcode response;

  if (url==NULL) {
    return WIMC_ERROR_BAD_URL;
  }

  curl = curl_easy_init();

  if (curl) {
    
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_NOBODY, 1);

    response = curl_easy_perform(curl);

    curl_easy_cleanup(curl);
  }

  if (response != CURLE_OK) {
    return WIMC_ERROR_BAD_URL;
  }
#endif
  return WIMC_OK;
}

/*
 * valid_url_format -- A valid URL is of format http://host:port/
 */

int valid_url_format(char *url) {

  if (url == NULL) {
    return FALSE;
  }

  /* XXX */

  return TRUE;
}

