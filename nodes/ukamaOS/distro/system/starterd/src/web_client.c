/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "web_client.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "usys_log.h"
#include "usys_mem.h"

#include <ulfius.h>

static URequest* wc_create_request(const char *url,
                                   const char *method,
                                   int timeoutSec) {

    URequest *req;

    req = (URequest *)usys_calloc(1, sizeof(URequest));
    if (!req) return NULL;

    if (ulfius_init_request(req) != U_OK) {
        usys_free(req);
        return NULL;
    }

    ulfius_set_request_properties(req,
                                  U_OPT_HTTP_VERB, method,
                                  U_OPT_HTTP_URL, url,
                                  U_OPT_TIMEOUT, timeoutSec,
                                  U_OPT_NONE);

    return req;
}

static bool wc_send(URequest *req, UResponse **respOut) {

    UResponse *resp;

    resp = (UResponse *)usys_calloc(1, sizeof(UResponse));
    if (!resp) return false;

    if (ulfius_init_response(resp) != U_OK) {
        usys_free(resp);
        return false;
    }

    if (ulfius_send_http_request(req, resp) != U_OK) {
        ulfius_clean_response(resp);
        usys_free(resp);
        return false;
    }

    *respOut = resp;
    return true;
}

static void wc_clean(URequest *req, UResponse *resp) {

    if (req) {
        ulfius_clean_request(req);
        usys_free(req);
    }
    if (resp) {
        ulfius_clean_response(resp);
        usys_free(resp);
    }
}

static bool wc_build_url(char *buf,
                         size_t buflen,
                         const char *addr,
                         int port,
                         const char *path) {

    int n;

    if (!buf || buflen == 0 || !addr || !path) return false;

    n = snprintf(buf, buflen, "http://%s:%d%s", addr, port, path);
    return (n > 0 && (size_t)n < buflen);
}

bool wc_app_ping(Config *config, App *app) {

    char url[256];
    URequest *req;
    UResponse *resp;
    bool ok;

    req = NULL;
    resp = NULL;
    ok = false;

    if (!config || !app) return false;
    if (app->port <= 0) return false;

    if (!wc_build_url(url, sizeof(url), "127.0.0.1", app->port, "/v1/ping")) {
        return false;
    }

    req = wc_create_request(url, "GET", config->pingTimeoutSec);
    if (!req) return false;

    if (!wc_send(req, &resp)) {
        wc_clean(req, NULL);
        return false;
    }

    if (resp->status == 200) ok = true;

    wc_clean(req, resp);
    return ok;
}

bool wc_app_version_matches(Config *config,
                            App *app,
                            const char *tag) {

    char url[256];
    URequest *req;
    UResponse *resp;
    bool ok;
    const char *body;

    req = NULL;
    resp = NULL;
    ok = false;
    body = NULL;

    if (!config || !app || !tag) return false;
    if (app->port <= 0) return false;

    if (!wc_build_url(url, sizeof(url), "127.0.0.1", app->port, "/v1/version")) {
        return false;
    }

    req = wc_create_request(url, "GET", config->commitTimeoutSec);
    if (!req) return false;

    if (!wc_send(req, &resp)) {
        wc_clean(req, NULL);
        return false;
    }

    if (resp->status == 200) {
        body = resp->binary_body ? (const char *)resp->binary_body : NULL;
        if (body && strstr(body, tag) != NULL) ok = true;
    }

    wc_clean(req, resp);
    return ok;
}

bool wc_fetch_package(Config *config,
                      const char *appName,
                      const char *tag,
                      const char *hub,
                      const char *dstPath) {

    char url[512];
    char path[256];
    URequest *req;
    UResponse *resp;
    FILE *f;
    bool ok;
    JsonObj *jreq;
    char *body;

    req = NULL;
    resp = NULL;
    f = NULL;
    ok = false;
    jreq = NULL;
    body = NULL;

    if (!config || !appName || !tag || !dstPath) return false;

    snprintf(path,
             sizeof(path),
             config->wimcPathTemplate ?
                 config->wimcPathTemplate :
                 "/v1/apps/%s/%s/pkg",
             appName,
             tag);

    if (!wc_build_url(url, sizeof(url), config->wimcHost, config->wimcPort, path)) {
        usys_log_error("wimc: url build failed");
        return false;
    }

    if (hub && *hub) {
        req = wc_create_request(url, "POST", 60);
    } else {
        req = wc_create_request(url, "GET", 60);
    }

    if (!req) return false;

    if (hub && *hub) {
        jreq = json_object();
        if (!jreq) {
            wc_clean(req, NULL);
            return false;
        }

        json_object_set_new(jreq, "hub", json_string(hub));
        body = json_dumps(jreq, JSON_COMPACT);
        json_decref(jreq);

        if (!body) {
            wc_clean(req, NULL);
            return false;
        }

        ulfius_set_string_body_request(req, body);
        u_map_put(req->map_header, "Content-Type", "application/json");
    }

    if (!wc_send(req, &resp)) {
        wc_clean(req, NULL);
        free(body);
        usys_log_error("wimc: request failed %s", url);
        return false;
    }

    free(body);

    if (resp->status != 200 || !resp->binary_body || resp->binary_body_length == 0) {
        usys_log_error("wimc: bad response http=%d", resp->status);
        wc_clean(req, resp);
        return false;
    }

    f = fopen(dstPath, "wb");
    if (!f) {
        usys_log_error("wimc: cannot open %s", dstPath);
        wc_clean(req, resp);
        return false;
    }

    if (fwrite(resp->binary_body, 1, resp->binary_body_length, f) != resp->binary_body_length) {
        usys_log_error("wimc: write failed %s", dstPath);
        fclose(f);
        wc_clean(req, resp);
        return false;
    }

    fclose(f);
    ok = true;

    wc_clean(req, resp);
    return ok;
}
