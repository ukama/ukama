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

static bool wc_copy_file(const char *srcPath, const char *dstPath) {

    FILE *src;
    FILE *dst;
    char buf[8192];
    size_t n;

    src = NULL;
    dst = NULL;

    if (!srcPath || !dstPath) return false;

    src = fopen(srcPath, "rb");
    if (!src) return false;

    dst = fopen(dstPath, "wb");
    if (!dst) {
        fclose(src);
        return false;
    }

    while ((n = fread(buf, 1, sizeof(buf), src)) > 0) {
        if (fwrite(buf, 1, n, dst) != n) {
            fclose(src);
            fclose(dst);
            return false;
        }
    }

    if (ferror(src)) {
        fclose(src);
        fclose(dst);
        return false;
    }

    fclose(src);
    fclose(dst);
    return true;
}

static bool wc_wait_for_available(Config *config,
                                  const char *appName,
                                  const char *tag,
                                  int timeoutSec) {

    char url[512];
    char path[256];
    URequest *req;
    UResponse *resp;
    time_t start;
    const char *body;

    req = NULL;
    resp = NULL;
    body = NULL;

    snprintf(path, sizeof(path), "/v1/apps/%s/%s/status", appName, tag);
    if (!wc_build_url(url, sizeof(url), config->wimcHost, config->wimcPort, path)) {
        usys_log_error("wimc: status url build failed");
        return false;
    }

    start = time(NULL);
    while ((int)(time(NULL) - start) < timeoutSec) {

        req = wc_create_request(url, "GET", 10);
        if (!req) {
            return false;
        }

        if (!wc_send(req, &resp)) {
            wc_clean(req, NULL);
            usleep(500 * 1000);
            continue;
        }

        if (resp->status == 200 && resp->binary_body) {
            body = (const char *)resp->binary_body;

            if (strstr(body, "\"available\"") != NULL) {
                wc_clean(req, resp);
                return true;
            }

            if (strstr(body, "\"failed\"") != NULL) {
                wc_clean(req, resp);
                return false;
            }
        }

        wc_clean(req, resp);
        req = NULL;
        resp = NULL;

        usleep(500 * 1000);
    }

    return false;
}

bool wc_fetch_package(Config *config,
                      const char *appName,
                      const char *tag,
                      const char *hub,
                      const char *dstPath) {

    char url[512];
    char path[256];
    char srcPath[512];
    URequest *req;
    UResponse *resp;
    JsonObj *jreq;
    char *body;
    bool ok;

    req = NULL;
    resp = NULL;
    jreq = NULL;
    body = NULL;
    ok = false;

    if (!config || !appName || !tag || !dstPath) return false;

    snprintf(path,
             sizeof(path),
             config->wimcPathTemplate ?
                 config->wimcPathTemplate :
                 "/v1/apps/%s/%s",
             appName,
             tag);

    if (!wc_build_url(url, sizeof(url), config->wimcHost, config->wimcPort, path)) {
        usys_log_error("wimc: url build failed");
        return false;
    }

    req = wc_create_request(url, "POST", 30);
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
        free(body);
        wc_clean(req, NULL);
        usys_log_error("wimc: request failed %s", url);
        return false;
    }

    free(body);

    if (resp->status != 202 && resp->status != 304 && resp->status != 409) {
        usys_log_error("wimc: unexpected response http=%d", resp->status);
        wc_clean(req, resp);
        return false;
    }

    wc_clean(req, resp);

    if (!wc_wait_for_available(config, appName, tag, 120)) {
        usys_log_error("wimc: package not available %s:%s", appName, tag);
        return false;
    }

    snprintf(srcPath, sizeof(srcPath), "/ukama/apps/pkgs/%s_%s.tar.gz", appName, tag);

    if (!wc_copy_file(srcPath, dstPath)) {
        usys_log_error("wimc: failed copying package %s -> %s", srcPath, dstPath);
        return false;
    }

    ok = true;
    return ok;
}
