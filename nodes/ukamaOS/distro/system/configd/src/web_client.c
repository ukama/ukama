/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "httpStatus.h"
#include "web_client.h"
#include "configd.h"
#include "jserdes.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static bool wc_send_http_request(URequest* httpReq, UResponse** httpResp) {

	*httpResp = (UResponse *)usys_calloc(1, sizeof(UResponse));
	if (*httpResp == NULL) {
		usys_log_error("Error allocating memory of size: %lu for http response",
                       sizeof(UResponse));
		return USYS_FALSE;
	}

	if (ulfius_init_response(*httpResp)) {
		usys_log_error("Error initializing new http response.");
		return USYS_FALSE;
	}

	if (ulfius_send_http_request(httpReq, *httpResp) != STATUS_OK) {
		usys_log_error( "Web client failed to send %s web request to %s",
                        httpReq->http_verb, httpReq->http_url);

        ulfius_clean_response(*httpResp);
        usys_free(*httpResp);
        *httpResp = NULL;

        return USYS_FALSE;
	}

	return USYS_TRUE;
}

URequest* wc_create_http_request(char *url,
                                 const char *method,
                                 JsonObj *body) {

	/* Preparing Request */
	URequest* httpReq = (URequest *)usys_calloc(1, sizeof(URequest));
	if (!httpReq) {
		usys_log_error("Error allocating memory of size: %lu for http Request",
				sizeof(URequest));
		return NULL;
	}

	if (ulfius_init_request(httpReq)) {
		usys_log_error("Error initializing new http request.");
		return NULL;
	}

	ulfius_set_request_properties(httpReq,
			U_OPT_HTTP_VERB, method,
			U_OPT_HTTP_URL, url,
			U_OPT_TIMEOUT, 20,
			U_OPT_NONE);

	if (body) {
		if (STATUS_OK != ulfius_set_json_body_request(httpReq, body)) {
			ulfius_clean_request(httpReq);
			usys_free(httpReq);
			httpReq = NULL;
		}
	}

	return httpReq;
}

int wc_send_node_info_request(char *url, char *method, char **nodeID) {

	int ret = STATUS_NOK;
	JsonObj *json = NULL;
	JsonErrObj jErr;
	UResponse *httpResp = NULL;
	URequest *httpReq = NULL;

	httpReq = wc_create_http_request(url, method, NULL);
	if (!httpReq) {
		return ret;
	}

	if (wc_send_http_request(httpReq, &httpResp) == USYS_FALSE) {
		usys_log_error("Failed to send http request.");
		goto cleanup;
	}

	if (httpResp->status == HttpStatus_OK) {
		json = ulfius_get_json_body_response(httpResp, &jErr);
		if (json) {
			ret = json_deserialize_node_id(nodeID, json);
			if (!ret) {
				usys_log_error("Failed to parse NodeInfo response from noded.");
				return STATUS_NOK;
			}
			ret = STATUS_OK;
		}
	} else {
		ret = STATUS_NOK;
	}

	json_decref(json);

cleanup:
	if (httpReq) {
		ulfius_clean_request(httpReq);
		usys_free(httpReq);
	}
	if (httpResp) {
		ulfius_clean_response(httpResp);
		usys_free(httpResp);
	}

	return ret;
}

int wc_read_node_info(Config* config) {

	int ret = STATUS_NOK;
	char url[128]={0};

	sprintf(url,"http://%s:%d/%s",
            config->nodedHost,
            config->nodedPort,
			config->nodedEP);

	ret = wc_send_node_info_request(url, "GET", &config->nodeId);
	if (ret) {
		usys_log_error("Failed to parse NodeInfo response from noded.");
		return ret;
	}

	return ret;
}

int get_nodeid_from_noded(Config *config) {

	if (wc_read_node_info(config)) {
		usys_log_error("Error reading NodeID from noded.d");
		return STATUS_NOK;
	}

	usys_log_info("notify.d: Node ID: %s", config->nodeId);

	return STATUS_OK;
}

static bool wc_send_app_request(Config *config,
                                char *app,
                                const char *action,
                                const char *httpMethod,
                                int expectedStatus) {

    UResponse *httpResp = NULL;
    URequest  *httpReq  = NULL;
    char url[MAX_URL] = {0};
    bool result;

    if (strcasecmp(action, "exec") == 0) {
        snprintf(url, sizeof(url), "http://%s:%d/v1/%s/%s/%s/latest",
                 config->starterHost,
                 config->starterPort,
                 action,
                 DEF_SPACE_NAME,
                 app);
    } else {
        snprintf(url, sizeof(url), "http://%s:%d/v1/%s/%s/%s",
                 config->starterHost,
                 config->starterPort,
                 action,
                 DEF_SPACE_NAME,
                 app);
    }

    httpReq = wc_create_http_request(url, httpMethod, NULL);
    if (!httpReq) {
        return USYS_FALSE;
    }

    if (wc_send_http_request(httpReq, &httpResp) == USYS_FALSE) {
        usys_log_error("Failed to send http request.");
        ulfius_clean_request(httpReq);
        usys_free(httpReq);
        return USYS_FALSE;
    }

    result = (httpResp->status == expectedStatus) ? USYS_TRUE : USYS_FALSE;

    ulfius_clean_request(httpReq);
    ulfius_clean_response(httpResp);
    usys_free(httpReq);
    usys_free(httpResp);

    return result;
}

bool wc_send_app_restart_request(Config *config, char *app) {

    bool result;

    result = wc_send_app_request(config,
                                 app,
                                 "terminate",
                                 "POST",
                                 HttpStatus_Accepted);

    if (result) {
        return  wc_send_app_request(config,
                               app,
                               "exec",
                               "POST",
                               HttpStatus_Accepted);
    }

    return USYS_FALSE;
}

bool wc_is_app_valid(Config *config, char *app) {

    bool result;

    result = wc_send_app_request(config,
                                 app,
                                 "status",
                                 "GET",
                                 HttpStatus_OK);

    if (result) {
        usys_log_debug("App found by starter.d. Is valid: %s", app);
    } else {
        usys_log_error("App not found by starter.d: %s", app);
    }

    return result;
}
