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

bool wc_lifecycle_check_in(Config *config, bool bootHealthy) {

    char url[256];
    char *body;
    URequest *req;
    UResponse *resp;
    JsonObj *json;
    bool accepted;

    body = NULL;
    req = NULL;
    resp = NULL;
    json = NULL;
    accepted = false;

    if (!config || !config->lifecycleHost || config->lifecyclePort <= 0) {
        return false;
    }

    if (!wc_build_url(url,
                      sizeof(url),
                      config->lifecycleHost,
                      config->lifecyclePort,
                      "/v1/check-in")) {
        return false;
    }

    json = json_pack("{s:s}",
                     "bootResult",
                     bootHealthy ? "ready" : "degraded");
    if (!json) return false;

    body = json_dumps(json, JSON_COMPACT);
    json_decref(json);
    if (!body) return false;

    req = wc_create_request(url, "POST", config->pingTimeoutSec);
    if (!req) {
        free(body);
        return false;
    }

    ulfius_set_string_body_request(req, body);
    u_map_put(req->map_header, "Content-Type", "application/json");

    if (wc_send(req, &resp) && resp &&
        (resp->status == HttpStatus_Accepted ||
         resp->status == HttpStatus_OK)) {
        accepted = true;
    }

    free(body);
    wc_clean(req, resp);
    return accepted;
}

LifecycleGateState wc_lifecycle_gate(Config *config) {

    char url[256];
    URequest *req;
    UResponse *resp;
    JsonObj *root;
    JsonObj *proceed;
    JsonErrObj error;
    LifecycleGateState state;

    req = NULL;
    resp = NULL;
    root = NULL;
    state = LIFECYCLE_GATE_UNAVAILABLE;

    if (!config || !config->lifecycleHost || config->lifecyclePort <= 0) {
        return state;
    }

    if (!wc_build_url(url,
                      sizeof(url),
                      config->lifecycleHost,
                      config->lifecyclePort,
                      "/v1/gate")) {
        return state;
    }

    req = wc_create_request(url, "GET", config->pingTimeoutSec);
    if (!req) return state;

    if (!wc_send(req, &resp) || !resp) {
        wc_clean(req, resp);
        return state;
    }

    if (resp->status == HttpStatus_Accepted) {
        state = LIFECYCLE_GATE_WAITING;
    } else if (resp->status == HttpStatus_ServiceUnavailable) {
        state = LIFECYCLE_GATE_FAULTY;
    } else if (resp->status == HttpStatus_OK &&
               resp->binary_body &&
               resp->binary_body_length > 0) {
        memset(&error, 0, sizeof(error));
        root = json_loadb((const char *)resp->binary_body,
                          resp->binary_body_length,
                          0,
                          &error);
        proceed = root ? json_object_get(root, "proceed") : NULL;
        if (json_is_boolean(proceed)) {
            state = json_is_true(proceed) ?
                LIFECYCLE_GATE_OPEN : LIFECYCLE_GATE_WAITING;
        }
    }

    json_decref(root);
    wc_clean(req, resp);
    return state;
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
    int ret;

    req = NULL;
    resp = NULL;
    body = NULL;
    start = 0;

    if (pathOut != NULL) {
        *pathOut = NULL;
    }

    if (versionOut != NULL) {
        *versionOut = NULL;
    }

    if (config == NULL || appName == NULL || tag == NULL) {
        return false;
    }

    ret = snprintf(path,
                   sizeof(path),
                   "/v1/apps/%s/%s/status",
                   appName,
                   tag);
    if (ret < 0 || (size_t)ret >= sizeof(path)) {
        usys_log_error("wimc: status path too long %s:%s",
                       appName,
                       tag);
        return false;
    }

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
        if (req == NULL) {
            return false;
        }

        resp = NULL;
        body = NULL;

        if (wc_send(req, &resp)) {
            if (resp != NULL &&
                resp->status == HttpStatus_OK &&
                resp->binary_body != NULL &&
                resp->binary_body_length > 0) {

                body = wc_copy_response_body(resp);
                if (body != NULL) {
                    if (wc_json_status_is(body, "available")) {
                        if (pathOut != NULL) {
                            *pathOut = wc_json_dup_string(body, "path");
                        }

                        if (versionOut != NULL) {
                            *versionOut = wc_json_dup_string(body,
                                                             "actualVersion");
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
            } else if (resp != NULL &&
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

    if (!wc_build_url(url, sizeof(url),
                      "127.0.0.1",
                      probePort,
                      "/v1/ping")) {
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

bool wc_app_ready(Config *config, App *app, AppReadyResponse *result) {

    char url[256];
    URequest *req;
    UResponse *resp;
    JsonObj *root;
    JsonObj *value;
    JsonErrObj error;
    const char *reason;
    const char *requestId;
    int probePort;
    bool valid;

    req = NULL;
    resp = NULL;
    root = NULL;
    valid = false;

    if (!config || !app || !result) return false;

    memset(result, 0, sizeof(*result));
    probePort = wc_get_probe_port(app);
    if (probePort <= 0) return false;

    if (!wc_build_url(url, sizeof(url),
                      "127.0.0.1",
                      probePort,
                      "/v1/ready")) {
        return false;
    }

    req = wc_create_request(url, "GET", config->pingTimeoutSec);
    if (!req) return false;

    if (!wc_send(req, &resp)) {
        wc_clean(req, NULL);
        return false;
    }

    result->status = resp->status;
    if (resp->status != HttpStatus_OK &&
        resp->status != HttpStatus_Accepted &&
        resp->status != HttpStatus_ServiceUnavailable) {
        wc_clean(req, resp);
        return false;
    }

    if (!resp->binary_body || resp->binary_body_length == 0) {
        wc_clean(req, resp);
        return false;
    }

    memset(&error, 0, sizeof(error));
    root = json_loadb((const char *)resp->binary_body,
                      resp->binary_body_length,
                      0,
                      &error);
    if (!root) {
        wc_clean(req, resp);
        return false;
    }

    value = json_object_get(root, "ready");
    if (json_is_boolean(value)) {
        result->ready = json_is_true(value);
        valid = true;
    }

    value = json_object_get(root, "reason");
    reason = json_is_string(value) ? json_string_value(value) : NULL;
    snprintf(result->reason,
             sizeof(result->reason),
             "%s",
             (reason && *reason) ? reason :
             (result->ready ? "ready" : "not ready"));

    value = json_object_get(root, "requestId");
    requestId = json_is_string(value) ? json_string_value(value) : NULL;
    snprintf(result->requestId,
             sizeof(result->requestId),
             "%s",
             requestId ? requestId : "");

    if ((resp->status == HttpStatus_OK && !result->ready) ||
        (resp->status != HttpStatus_OK && result->ready)) {
        valid = false;
    }

    json_decref(root);
    wc_clean(req, resp);
    return valid;
}

bool wc_mesh_status(Config *config,
                    App *app,
                    bool *connected,
                    char *reason,
                    size_t reasonSize) {

    char url[256];
    URequest *req;
    UResponse *resp;
    JsonObj *root;
    JsonObj *value;
    JsonErrObj error;
    const char *text;
    int probePort;
    bool valid;

    req = NULL;
    resp = NULL;
    root = NULL;
    valid = false;

    if (!config || !app || !connected || !reason || reasonSize == 0) {
        return false;
    }

    *connected = false;
    reason[0] = '\0';

    probePort = wc_get_probe_port(app);
    if (probePort <= 0) return false;

    if (!wc_build_url(url, sizeof(url),
                      "127.0.0.1",
                      probePort,
                      "/v1/status")) {
        return false;
    }

    req = wc_create_request(url, "GET", config->pingTimeoutSec);
    if (!req) return false;

    if (!wc_send(req, &resp)) {
        wc_clean(req, NULL);
        return false;
    }

    if (resp->status != HttpStatus_OK ||
        !resp->binary_body ||
        resp->binary_body_length == 0) {
        wc_clean(req, resp);
        return false;
    }

    memset(&error, 0, sizeof(error));
    root = json_loadb((const char *)resp->binary_body,
                      resp->binary_body_length,
                      0,
                      &error);
    if (!root) {
        wc_clean(req, resp);
        return false;
    }

    value = json_object_get(root, "connected");
    if (json_is_boolean(value)) {
        *connected = json_is_true(value);
        valid = true;
    }

    value = json_object_get(root, "reason");
    text = json_is_string(value) ? json_string_value(value) : NULL;
    snprintf(reason,
             reasonSize,
             "%s",
             (text && *text) ? text :
             (*connected ? "connected" : "disconnected"));

    json_decref(root);
    wc_clean(req, resp);
    return valid;
}

static const char *wc_strip_v_prefix(const char *s) {

    if (s != NULL && s[0] == 'v' && s[1] != '\0') {
        return s + 1;
    }

    return s;
}

static bool wc_versions_equal(const char *a, const char *b) {

    if (a == NULL || b == NULL) {
        return false;
    }

    if (strcmp(a, b) == 0) {
        return true;
    }

    a = wc_strip_v_prefix(a);
    b = wc_strip_v_prefix(b);

    return strcmp(a, b) == 0;
}

bool wc_app_version_matches(Config *config,
                            App *app,
                            const char *tag) {

    char url[512];
    URequest *req;
    UResponse *resp;
    char *copy;
    char *p;
    char *end;
    int probePort;
    bool ok;

    req = NULL;
    resp = NULL;
    copy = NULL;
    ok = false;

    if (!config || !app || !tag) {
        return false;
    }

    probePort = wc_get_probe_port(app);
    if (probePort <= 0) return false;

    if (!wc_build_url(url, sizeof(url),
                      "127.0.0.1",
                      probePort,
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

            if (wc_versions_equal(p, tag)) {
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
                      char **pathOut,
                      char **versionOut) {

    char url[512];
    char path[256];
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

    if (pathOut != NULL) {
        *pathOut = NULL;
    }

    if (versionOut != NULL) {
        *versionOut = NULL;
    }

    if (config == NULL || appName == NULL || tag == NULL) {
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
        usys_log_error("wimc: request path too long %s:%s",
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

        if (hub[0] == '[') {
            json_error_t jerr;

            memset(&jerr, 0, sizeof(jerr));

            jhub = json_loads(hub, 0, &jerr);
            if (jhub == NULL || !json_is_array(jhub)) {
                if (jhub != NULL) {
                    json_decref(jhub);
                    jhub = NULL;
                }

                usys_log_error("wimc: invalid hub array payload");
                goto done;
            }
        } else {
            jhub = json_string(hub);
            if (jhub == NULL) {
                usys_log_error("wimc: failed creating hub json string");
                goto done;
            }
        }

        if (json_object_set_new(jreq, "hub", jhub) != 0) {
            json_decref(jhub);
            jhub = NULL;
            usys_log_error("wimc: failed setting hub in json request");
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

    if (availablePath == NULL || *availablePath == '\0') {
        usys_log_error("wimc: package available but path missing %s:%s",
                       appName,
                       tag);
        goto done;
    }

    if (!wc_file_exists_non_empty(availablePath)) {
        usys_log_error("wimc: package missing or empty %s", availablePath);
        goto done;
    }

    if (pathOut != NULL) {
        *pathOut = availablePath;
        availablePath = NULL;
    }

    if (versionOut != NULL) {
        *versionOut = actualVersion;
        actualVersion = NULL;
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
