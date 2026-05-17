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
#include <sys/stat.h>
#include <time.h>
#include <unistd.h>
#include <jansson.h>

#include "web_client.h"
#include "http_status.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_file.h"
#include "usys_services.h"

static bool wc_file_exists_non_empty(const char *path) {

    struct stat st;

    if (!path) {
        return false;
    }

    if (stat(path, &st) != 0) {
        return false;
    }

    if (!S_ISREG(st.st_mode)) {
        return false;
    }

    return st.st_size > 0;
}

static bool wc_json_status_is(const char *body, const char *expected) {

    json_t *root;
    json_t *status;
    json_error_t error;
    const char *value;
    bool ok;

    root   = NULL;
    status = NULL;
    value  = NULL;
    ok     = false;

    if (!body || !expected) {
        return false;
    }

    root = json_loads(body, 0, &error);
    if (!root) {
        return false;
    }

    status = json_object_get(root, "status");
    if (json_is_string(status)) {
        value = json_string_value(status);
        if (value && strcmp(value, expected) == 0) {
            ok = true;
        }
    }

    json_decref(root);
    return ok;
}

static char *wc_json_dup_string(const char *body, const char *key) {

    json_t *root;
    json_t *value;
    json_error_t error;
    const char *s;
    char *dup;

    root = NULL;
    value = NULL;
    s = NULL;
    dup = NULL;

    if (body == NULL || key == NULL) {
        return NULL;
    }

    root = json_loads(body, 0, &error);
    if (root == NULL) {
        return NULL;
    }

    value = json_object_get(root, key);
    if (json_is_string(value)) {
        s = json_string_value(value);
        if (s != NULL && *s != '\0') {
            dup = strdup(s);
        }
    }

    json_decref(root);
    return dup;
}

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
                                  const char *tag,
                                  char **pathOut,
                                  char **versionOut) {

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
        body = NULL;

        if (wc_send(req, &resp)) {

            if (resp &&
                resp->status == HttpStatus_OK &&
                resp->binary_body &&
                resp->binary_body_length > 0) {

                body = wc_copy_response_body(resp);
                if (body) {

                    if (wc_json_status_is(body, "available")) {
                        if (pathOut != NULL) {
                            *pathOut = wc_json_dup_string(body, "path");
                        }

                        if (versionOut != NULL) {
                            *versionOut = wc_json_dup_string(body, "actualVersion");
                        }

                        free(body);
                        wc_clean(req, resp);
                        return true;
                    }

                    if (wc_json_status_is(body, "failed") ||
                        wc_json_status_is(body, "corrupt") ||
                        wc_json_status_is(body, "missing")) {
                        usys_log_error("wimc: fetch failed %s:%s",
                                       appName,
                                       tag);
                        free(body);
                        wc_clean(req, resp);
                        return false;
                    }

                    free(body);
                    body = NULL;
                }

            } else if (resp &&
                       (resp->status == HttpStatus_NotFound ||
                        resp->status == HttpStatus_InternalServerError)) {

                wc_clean(req, resp);
                return false;
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
    const char *pkgsDir;
    const char *pathTag;
    URequest *req;
    UResponse *resp;
    JsonObj *jreq;
    JsonObj *jhub;
    char *body;
    char *availablePath;
    char *actualVersion;
    bool ok;
    int ret;

    req = NULL;
    resp = NULL;
    jreq = NULL;
    jhub = NULL;
    body = NULL;
    availablePath = NULL;
    actualVersion = NULL;
    ok = false;

    if (config == NULL || appName == NULL || tag == NULL || dstPath == NULL) {
        return false;
    }

    ret = snprintf(path,
                   sizeof(path),
                   config->wimcPathTemplate ?
                       config->wimcPathTemplate :
                       "/v1/apps/%s/%s",
                   appName,
                   tag);
    if (ret < 0 || (size_t)ret >= sizeof(path)) {
        usys_log_error("wimc: request path too long for %s:%s",
                       appName,
                       tag);
        return false;
    }

    if (!wc_build_url(url,
                      sizeof(url),
                      config->wimcHost,
                      config->wimcPort,
                      path)) {
        usys_log_error("wimc: url build failed");
        return false;
    }

    req = wc_create_request(url, "POST", 30);
    if (req == NULL) {
        usys_log_error("wimc: failed creating request %s", url);
        goto done;
    }

    if (hub != NULL && *hub != '\0') {
        jreq = json_object();
        if (jreq == NULL) {
            usys_log_error("wimc: failed creating json request");
            goto done;
        }

        jhub = json_string(hub);
        if (jhub == NULL) {
            usys_log_error("wimc: failed creating hub json string");
            goto done;
        }

        if (json_object_set_new(jreq, "hub", jhub) != 0) {
            usys_log_error("wimc: failed setting hub in json request");
            json_decref(jhub);
            jhub = NULL;
            goto done;
        }
        jhub = NULL;

        body = json_dumps(jreq, JSON_COMPACT);
        if (body == NULL) {
            usys_log_error("wimc: failed dumping json request");
            goto done;
        }

        ulfius_set_string_body_request(req, body);
        u_map_put(req->map_header, "Content-Type", "application/json");
    }

    if (!wc_send(req, &resp)) {
        usys_log_error("wimc: request failed %s", url);
        goto done;
    }

    if (resp == NULL) {
        usys_log_error("wimc: empty response from %s", url);
        goto done;
    }

    if (resp->status != HttpStatus_OK &&
        resp->status != HttpStatus_Accepted &&
        resp->status != HttpStatus_NotModified &&
        resp->status != HttpStatus_Conflict) {

        usys_log_error("wimc: unexpected response http=%d", resp->status);
        goto done;
    }

    wc_clean(req, resp);
    req = NULL;
    resp = NULL;

    if (!wc_wait_for_available(config,
                               appName,
                               tag,
                               &availablePath,
                               &actualVersion)) {
        usys_log_error("wimc: package not available %s:%s",
                       appName,
                       tag);
        goto done;
    }

    pkgsDir = config->pkgsDir ? config->pkgsDir : "/ukama/apps/pkgs";

    if (availablePath != NULL && *availablePath != '\0') {
        ret = snprintf(srcPath, sizeof(srcPath), "%s", availablePath);
    } else {
        /*
         * If wimc/agent validated the package VERSION and returned the real
         * version, prefer it for fallback path construction. This handles:
         *
         *   requested tag: 1.0.1-abcdefgh
         *   actual tag:    v1.0.1-abcdefgh
         *
         * The normal contract should still be: use availablePath from wimc.
         */
        pathTag = actualVersion && *actualVersion ? actualVersion : tag;

        ret = snprintf(srcPath,
                       sizeof(srcPath),
                       "%s/%s_%s.tar.gz",
                       pkgsDir,
                       appName,
                       pathTag);
    }

    if (ret < 0 || (size_t)ret >= sizeof(srcPath)) {
        usys_log_error("wimc: source package path too long %s:%s",
                       appName,
                       tag);
        goto done;
    }

    /*
     * In the normal starter/wimc contract both daemons share the same
     * package cache. Do not copy the package onto itself.
     */
    if (strcmp(srcPath, dstPath) == 0) {
        if (!wc_file_exists_non_empty(dstPath)) {
            usys_log_error("wimc: package missing or empty %s", dstPath);
            goto done;
        }

        ok = true;
        goto done;
    }

    if (!wc_copy_file(srcPath, dstPath)) {
        usys_log_error("wimc: failed copying package %s -> %s",
                       srcPath,
                       dstPath);
        goto done;
    }

    if (!wc_file_exists_non_empty(dstPath)) {
        usys_log_error("wimc: copied package is missing or empty %s",
                       dstPath);
        goto done;
    }

    ok = true;

done:
    if (jhub != NULL) {
        json_decref(jhub);
    }

    if (jreq != NULL) {
        json_decref(jreq);
    }

    free(body);
    free(availablePath);
    free(actualVersion);

    if (req != NULL || resp != NULL) {
        wc_clean(req, resp);
    }

    return ok;
}
