/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ulfius.h>

#include "web_client.h"
#include "http_status.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_file.h"
#include "usys_services.h"

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

static char *wc_copy_response_body(UResponse *resp) {

    char *copy;
    size_t len;

    copy = NULL;
    len  = 0;

    if (!resp || !resp->binary_body || resp->binary_body_length <= 0) {
        return NULL;
    }

    len = (size_t)resp->binary_body_length;

    copy = (char *)calloc(1, len + 1);
    if (!copy) {
        return NULL;
    }

    memcpy(copy, resp->binary_body, len);
    copy[len] = '\0';

    return copy;
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

static int wc_get_probe_port(App *app) {

    int port;

    port = -1;

    if (!app) return -1;

    /* some apps expose /v1/ping and /v1/version on their admin port,
     * not on the main app port.
     */
    if (strcmp(app->name, "metrics") == 0) {
        port = usys_find_service_port(SERVICE_METRICS_ADMIN);
    }
    else if (strcmp(app->name, "notify") == 0) {
        port = usys_find_service_port(SERVICE_NOTIFY_ADMIN);
    }
    else if (strcmp(app->name, "rlog") == 0) {
        port = usys_find_service_port(SERVICE_RLOG_ADMIN);
    }
    else {
        port = app->port;
    }

    if (port <= 0) {
        port = app->port;
    }

    return port;
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
                                  const char *tag) {

    char url[512];
    char path[256];
    URequest *req;
    UResponse *resp;
    char *body;
    time_t start;

    req   = NULL;
    resp  = NULL;
    body  = NULL;
    start = 0;

    if (!config || !appName || !tag) {
        return false;
    }

    snprintf(path,
             sizeof(path),
             "/v1/apps/%s/%s/status",
             appName,
             tag);

    if (!wc_build_url(url,
                      sizeof(url),
                      config->wimcHost,
                      config->wimcPort,
                      path)) {
        usys_log_error("wimc: status url build failed");
        return false;
    }

    start = time(NULL);

    while (true) {

        req = wc_create_request(url, "GET", config->pingTimeoutSec);
        if (!req) {
            return false;
        }

        resp = NULL;

        if (wc_send(req, &resp)) {

            if (resp &&
                resp->status == HttpStatus_OK &&
                resp->binary_body &&
                resp->binary_body_length > 0) {

                body = wc_copy_response_body(resp);
                if (body) {

                    if (strstr(body, "\"available\"") != NULL) {
                        free(body);
                        wc_clean(req, resp);
                        return true;
                    }

                    if (strstr(body, "\"failed\"") != NULL) {
                        free(body);
                        wc_clean(req, resp);
                        return false;
                    }

                    free(body);
                    body = NULL;
                }
            }
        }

        wc_clean(req, resp);

        if ((int)(time(NULL) - start) >= config->commitTimeoutSec) {
            break;
        }

        usleep(200 * 1000);
    }

    usys_log_error("wimc: timed out waiting for %s:%s", appName, tag);
    return false;
}

bool wc_app_ping(Config *config, App *app) {

    char url[256];
    URequest *req;
    UResponse *resp;
    bool ok;
    int probePort;

    req       = NULL;
    resp      = NULL;
    ok        = false;
    probePort = -1;

    if (!config || !app) return false;

    probePort = wc_get_probe_port(app);
    if (probePort <= 0) return false;

    if (!wc_build_url(url, sizeof(url), "127.0.0.1", probePort, "/v1/ping")) {
        return false;
    }

    req = wc_create_request(url, "GET", config->pingTimeoutSec);
    if (!req) return false;

    if (!wc_send(req, &resp)) {
        wc_clean(req, NULL);
        return false;
    }

    if (resp->status == HttpStatus_OK) {
        ok = true;
    }

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
    char *copy;
    char *p;
    char *end;
    int probePort;

    req       = NULL;
    resp      = NULL;
    ok        = false;
    copy      = NULL;
    probePort = -1;

    if (!config || !app || !tag) return false;

    probePort = wc_get_probe_port(app);
    if (probePort <= 0) return false;

    if (!wc_build_url(url, sizeof(url), "127.0.0.1", probePort,
                      "/v1/version")) {
        return false;
    }

    req = wc_create_request(url, "GET", config->commitTimeoutSec);
    if (!req) return false;

    if (!wc_send(req, &resp)) {
        wc_clean(req, NULL);
        return false;
    }

    if (resp->status == HttpStatus_OK) {
        copy = wc_copy_response_body(resp);
        if (copy) {
            p = copy;

            while (*p == ' ' || *p == '\t' ||
                   *p == '\r' || *p == '\n') {
                p++;
            }

            end = p + strlen(p);
            while (end > p &&
                   (end[-1] == ' ' || end[-1] == '\t' ||
                    end[-1] == '\r' || end[-1] == '\n')) {
                end--;
            }
            *end = '\0';

            if (strcmp(p, tag) == 0) {
                ok = true;
            }
        }
    }

    free(copy);
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

    if (resp->status != HttpStatus_Accepted &&
        resp->status != HttpStatus_NotModified &&
        resp->status != HttpStatus_Conflict) {
        usys_log_error("wimc: unexpected response http=%d", resp->status);
        wc_clean(req, resp);
        return false;
    }

    wc_clean(req, resp);

    if (!wc_wait_for_available(config, appName, tag)) {
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
