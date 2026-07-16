/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <jansson.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "bff.h"
#include "util.h"

#define BFF_CONTROL_TIMEOUT_SEC 30L

static const char *BFF_RESTART_NODE =
"mutation RestartNode($data: RestartNodeInputDto!) {"
" restartNode(data: $data) { success } }";

static const char *BFF_TOGGLE_SERVICE =
"mutation ToggleService($data: ToggleSiteStatusInputDto!) {"
" toggleService(data: $data) { success } }";

static const char *BFF_TOGGLE_RADIO =
"mutation ToggleRFStatus($data: ToggleSiteStatusInputDto!) {"
" toggleRFStatus(data: $data) { success } }";

typedef struct {
    char   *buf;
    size_t  len;
} control_buf_t;

static size_t control_write_cb(void *ptr, size_t size, size_t nmemb,
                               void *data) {
    control_buf_t *resp;
    size_t len;
    char *next;

    resp = data;
    len = size * nmemb;
    next = realloc(resp->buf, resp->len + len + 1);
    if (next == NULL) {
        return 0;
    }

    resp->buf = next;
    memcpy(resp->buf + resp->len, ptr, len);
    resp->len += len;
    resp->buf[resp->len] = '\0';

    return len;
}

static int control_response_success(json_t *root, const char *field,
                                    ulab_error_t *err) {
    json_t *data;
    json_t *result;
    json_t *success;

    data = json_object_get(root, "data");
    result = data ? json_object_get(data, field) : NULL;
    success = result ? json_object_get(result, "success") : NULL;

    if (success == NULL || !json_is_boolean(success)) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s missing boolean success", field);
        return ULAB_ERR;
    }

    if (!json_is_true(success)) {
        snprintf(err->msg, sizeof(err->msg), "%s returned success=false",
                 field);
        return ULAB_ERR;
    }

    return ULAB_OK;
}

static int control_call(bff_client_t *c, const char *op,
                        const char *query, const char *vars,
                        const char *field, ulab_error_t *err) {
    CURL *curl;
    CURLcode curl_rc;
    struct curl_slist *headers;
    control_buf_t resp;
    char escaped_query[4096];
    char body[8192];
    char token_header[8192];
    long http_code;
    json_error_t json_err;
    json_t *root;
    json_t *errors;
    int rc;

    headers = NULL;
    resp.buf = NULL;
    resp.len = 0;
    http_code = 0;
    root = NULL;
    rc = ULAB_ERR;

    if (c == NULL || c->url[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "%s missing BFF URL", op);
        return ULAB_ERR;
    }

    ulab_json_escape(query, escaped_query, sizeof(escaped_query));
    if (snprintf(body, sizeof(body),
                 "{\"query\":\"%s\",\"variables\":%s}",
                 escaped_query, vars) >= (int)sizeof(body)) {
        snprintf(err->msg, sizeof(err->msg), "%s request too long", op);
        return ULAB_ERR;
    }

    if (c->logf != NULL) {
        fprintf(c->logf, "--- %s request ---\n%s\n", op, body);
        fflush(c->logf);
    }

    curl = curl_easy_init();
    if (curl == NULL) {
        snprintf(err->msg, sizeof(err->msg), "%s curl init failed", op);
        return ULAB_ERR;
    }

    headers = curl_slist_append(headers, "Content-Type: application/json");
    if (c->authenticated) {
        if (snprintf(token_header, sizeof(token_header),
                     "X-Session-Token: %s", c->token) >=
            (int)sizeof(token_header)) {
            snprintf(err->msg, sizeof(err->msg),
                     "%s session token too long", op);
            goto done;
        }
        headers = curl_slist_append(headers, token_header);
    }

    curl_easy_setopt(curl, CURLOPT_URL, c->url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, control_write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &resp);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, BFF_CONTROL_TIMEOUT_SEC);

    curl_rc = curl_easy_perform(curl);
    if (curl_rc != CURLE_OK) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s HTTP request failed: %s", op,
                 curl_easy_strerror(curl_rc));
        goto done;
    }

    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_code);

    if (c->logf != NULL) {
        fprintf(c->logf, "--- %s response %ld ---\n%s\n", op,
                http_code, resp.buf ? resp.buf : "");
        fflush(c->logf);
    }

    if (http_code < 200 || http_code >= 300) {
        snprintf(err->msg, sizeof(err->msg), "%s HTTP %ld: %.512s", op,
                 http_code, resp.buf ? resp.buf : "");
        goto done;
    }

    root = json_loads(resp.buf ? resp.buf : "", 0, &json_err);
    if (root == NULL) {
        snprintf(err->msg, sizeof(err->msg), "%s invalid JSON: %s", op,
                 json_err.text);
        goto done;
    }

    errors = json_object_get(root, "errors");
    if (errors != NULL) {
        char *errors_json;

        errors_json = json_dumps(errors, JSON_COMPACT);
        snprintf(err->msg, sizeof(err->msg),
                 "%s GraphQL error: %.512s", op,
                 errors_json ? errors_json : "unknown");
        free(errors_json);
        goto done;
    }

    rc = control_response_success(root, field, err);

done:
    if (root != NULL) {
        json_decref(root);
    }
    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    free(resp.buf);

    return rc;
}

int bff_restart_node(bff_client_t *c, const node_t *node,
                     ulab_error_t *err) {
    char vars[1024];

    if (node == NULL || node->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "restartNode missing node id");
        return ULAB_ERR;
    }

    if (snprintf(vars, sizeof(vars),
                 "{\"data\":{\"nodeId\":\"%s\"}}",
                 node->bff_id) >= (int)sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "restartNode variables too long");
        return ULAB_ERR;
    }

    return control_call(c, "restartNode", BFF_RESTART_NODE, vars,
                        "restartNode", err);
}

int bff_toggle_site_service(bff_client_t *c, const site_t *site,
                            int enabled, ulab_error_t *err) {
    char vars[1024];

    if (site == NULL || site->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "toggleService missing site id");
        return ULAB_ERR;
    }

    if (snprintf(vars, sizeof(vars),
                 "{\"data\":{\"siteId\":\"%s\","
                 "\"status\":%s}}",
                 site->bff_id, enabled ? "true" : "false") >=
        (int)sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "toggleService variables too long");
        return ULAB_ERR;
    }

    return control_call(c, "toggleService", BFF_TOGGLE_SERVICE, vars,
                        "toggleService", err);
}

int bff_toggle_site_radio(bff_client_t *c, const site_t *site,
                          int enabled, ulab_error_t *err) {
    char vars[1024];

    if (site == NULL || site->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "toggleRFStatus missing site id");
        return ULAB_ERR;
    }

    if (snprintf(vars, sizeof(vars),
                 "{\"data\":{\"siteId\":\"%s\","
                 "\"status\":%s}}",
                 site->bff_id, enabled ? "true" : "false") >=
        (int)sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "toggleRFStatus variables too long");
        return ULAB_ERR;
    }

    return control_call(c, "toggleRFStatus", BFF_TOGGLE_RADIO, vars,
                        "toggleRFStatus", err);
}
