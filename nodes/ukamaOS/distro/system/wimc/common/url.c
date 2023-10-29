/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/* Utility functions related to URL. */

#include <curl/curl.h>
#include <curl/easy.h>

#include "wimc.h"
#include "err.h"

#include "usys_log.h"
#include "usys_types.h"

bool validate_url(char *url) {
  
    CURL *curl;
    CURLcode response;

    if (url == NULL) return USYS_FALSE;

    curl = curl_easy_init();
    if (curl) {
    
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_NOBODY, 1);

        response = curl_easy_perform(curl);

        curl_easy_cleanup(curl);
    }

    if (response != CURLE_OK) {
        return USYS_FALSE;
    }

    return USYS_TRUE;
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

