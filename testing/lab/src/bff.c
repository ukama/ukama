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
#include <unistd.h>

#include "bff.h"
#include "log.h"
#include "util.h"

extern const char *BFF_ADD_NETWORK;
extern const char *BFF_ADD_SITE;
extern const char *BFF_ADD_PACKAGE;
extern const char *BFF_UPDATE_PACKAGE;
extern const char *BFF_GET_PACKAGE;
extern const char *BFF_GET_PACKAGES;
extern const char *BFF_PACKAGE_NAME_AVAILABLE;
extern const char *BFF_PACKAGES_DASHBOARD;
extern const char *BFF_REVENUE_OVERVIEW;
extern const char *BFF_NETWORK_HOME;
extern const char *BFF_NODES_LIST;
extern const char *BFF_INVENTORY_OVERVIEW;
extern const char *BFF_SIM_POOL_OVERVIEW;
extern const char *BFF_SUBSCRIBER_DETAIL;
extern const char *BFF_ADD_SUBSCRIBER;
extern const char *BFF_ALLOCATE_SIM;
extern const char *BFF_GET_DATA_USAGE;
extern const char *BFF_GET_SIMS_USAGE_BY_NETWORK;
extern const char *BFF_GET_SIM_PACKAGES;
extern const char *BFF_ADD_PAYMENT;
extern const char *BFF_GET_PAYMENTS;
extern const char *BFF_GET_KPI_VALUES;
extern const char *BFF_GET_PERFORMANCE_REPORT;
extern const char *BFF_GET_NODE;
extern const char *BFF_GET_RELEASE_CATALOG;
extern const char *BFF_PROMOTE_RELEASE;
extern const char *BFF_UPDATE_SOFTWARE;
extern const char *BFF_GET_APPS;
extern const char *BFF_NETWORK_OVERVIEW;
extern const char *BFF_SITE_VIEW;
extern const char *BFF_GET_NETWORKS;
extern const char *BFF_GET_SITES;
extern const char *BFF_GET_SITES_LEGACY;
extern const char *BFF_GET_NODES_FOR_SITE;
extern const char *BFF_GET_COMPONENTS_BY_USER_ID;
extern const char *BFF_GET_NODES;

typedef struct {
    char *buf;
    size_t len;
} http_buf_t;

static size_t write_cb(void *ptr, size_t size, size_t nmemb, void *data) {
    http_buf_t *b;
    size_t n;
    char *p;

    b = data;
    n = size * nmemb;
    p = realloc(b->buf, b->len + n + 1);

    if (p == NULL) {
        return 0;
    }

    b->buf = p;
    memcpy(b->buf + b->len, ptr, n);
    b->len += n;
    b->buf[b->len] = '\0';

    return n;
}

static int json_get_str(json_t *obj, const char *key, char *out,
                        size_t out_len) {
    json_t *v;
    const char *s;

    v = json_object_get(obj, key);
    if (v == NULL || !json_is_string(v)) {
        return ULAB_ERR;
    }

    s = json_string_value(v);
    if (s == NULL) {
        return ULAB_ERR;
    }

    return ulab_copy(out, out_len, s);
}

static void json_get_optional_str(json_t *obj, const char *key,
                                  char *out, size_t out_len) {
    json_t *v;
    const char *s;

    if (out == NULL || out_len == 0) {
        return;
    }
    out[0] = '\0';
    v = obj ? json_object_get(obj, key) : NULL;
    if (v == NULL || !json_is_string(v)) {
        return;
    }
    s = json_string_value(v);
    if (s != NULL) {
        ulab_copy(out, out_len, s);
    }
}

static int json_get_nested_str(json_t *obj,
                               const char *a,
                               const char *b,
                               char *out,
                               size_t out_len) {
    json_t *x;
    json_t *y;
    const char *v;

    x = json_object_get(obj, a);
    if (x == NULL) {
        return ULAB_ERR;
    }

    y = b == NULL ? x : json_object_get(x, b);
    if (y == NULL || !json_is_string(y)) {
        return ULAB_ERR;
    }

    v = json_string_value(y);
    if (v == NULL || v[0] == '\0') {
        return ULAB_ERR;
    }

    return ulab_copy(out, out_len, v);
}

static json_t *dig(json_t *root, const char *a, const char *b) {
    json_t *x;
    json_t *y;

    x = json_object_get(root, a);
    if (x == NULL) {
        return NULL;
    }

    if (b == NULL) {
        return x;
    }

    y = json_object_get(x, b);
    if (y == NULL) {
        return NULL;
    }

    return y;
}

static const char *env_or_default(const char *name, const char *def) {
    const char *v;

    v = getenv(name);
    if (v != NULL && v[0] != '\0') {
        return v;
    }

    return def;
}

static int build_url(char *out,
                     size_t out_len,
                     const char *base,
                     const char *path,
                     ulab_error_t *err) {
    int n;

    if (out == NULL || base == NULL || path == NULL) {
        snprintf(err->msg, sizeof(err->msg), "invalid URL argument");
        return ULAB_ERR;
    }

    n = snprintf(out, out_len, "%s%s", base, path);
    if (n < 0 || (size_t)n >= out_len) {
        snprintf(err->msg, sizeof(err->msg),
                 "URL too long: base=%s path=%s", base, path);
        return ULAB_ERR;
    }

    return ULAB_OK;
}

static int http_json_request(const char *op,
                             const char *url,
                             const char *method,
                             const char *body,
                             struct curl_slist *extra_hdrs,
                             json_t **out,
                             ulab_error_t *err) {
    CURL *curl;
    CURLcode ret;
    struct curl_slist *hdr;
    http_buf_t resp;
    long code;
    json_t *root;
    json_error_t json_err;

    hdr = NULL;
    resp.buf = NULL;
    resp.len = 0;
    code = 0;
    root = NULL;

    curl = curl_easy_init();
    if (curl == NULL) {
        snprintf(err->msg, sizeof(err->msg), "%s: curl init failed", op);
        return ULAB_ERR;
    }

    hdr = curl_slist_append(hdr, "accept: application/json");

    if (body != NULL) {
        hdr = curl_slist_append(hdr, "Content-Type: application/json");
    }

    while (extra_hdrs != NULL) {
        hdr = curl_slist_append(hdr, extra_hdrs->data);
        extra_hdrs = extra_hdrs->next;
    }

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, hdr);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &resp);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 30L);
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);

    if (method != NULL && ulab_streq(method, "PATCH")) {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "PATCH");
    } else if (method != NULL && ulab_streq(method, "POST")) {
        curl_easy_setopt(curl, CURLOPT_POST, 1L);
    }

    if (body != NULL) {
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body);
    }

    ret = curl_easy_perform(curl);
    if (ret != CURLE_OK) {
        snprintf(err->msg, sizeof(err->msg), "%s: HTTP request failed: %s",
                 op, curl_easy_strerror(ret));
        curl_slist_free_all(hdr);
        curl_easy_cleanup(curl);
        free(resp.buf);
        return ULAB_ERR;
    }

    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    curl_slist_free_all(hdr);
    curl_easy_cleanup(curl);

    if (code < 200 || code >= 300) {
        snprintf(err->msg, sizeof(err->msg), "%s: HTTP %ld: %s",
                 op, code, resp.buf ? resp.buf : "");
        free(resp.buf);
        return ULAB_ERR;
    }

    root = json_loads(resp.buf ? resp.buf : "", 0, &json_err);
    free(resp.buf);

    if (root == NULL) {
        snprintf(err->msg, sizeof(err->msg), "%s: invalid JSON: %s",
                 op, json_err.text);
        return ULAB_ERR;
    }

    *out = root;

    return ULAB_OK;
}

static int json_get_path_str(json_t *root,
                             const char *key,
                             char *out,
                             size_t out_len) {
    json_t *v;
    const char *s;

    v = json_object_get(root, key);
    if (v == NULL || !json_is_string(v)) {
        return ULAB_ERR;
    }

    s = json_string_value(v);
    if (s == NULL || s[0] == '\0') {
        return ULAB_ERR;
    }

    return ulab_copy(out, out_len, s);
}

static void derive_bff_base_url(const char *graphql_url,
                                char *out,
                                size_t out_len) {
    const char *p;
    size_t n;

    p = strstr(graphql_url, "/graphql");
    if (p == NULL) {
        ulab_copy(out, out_len, graphql_url);
        return;
    }

    n = (size_t)(p - graphql_url);
    if (n >= out_len) {
        n = out_len - 1;
    }

    memcpy(out, graphql_url, n);
    out[n] = '\0';
}

static void shell_quote(FILE *f, const char *s) {
    const char *p;

    fputc('\'', f);
    for (p = s; p != NULL && *p != '\0'; p++) {
        if (*p == '\'') {
            fprintf(f, "'\\''");
        } else {
            fputc(*p, f);
        }
    }
    fputc('\'', f);
}

static void bff_dump_curl(bff_client_t *c,
                          const char *op,
                          const char *body) {
    const char *dump;

    dump = getenv("UKAMA_LAB_DUMP_BFF_CURL");
    if (dump == NULL || dump[0] == '\0' || ulab_streq(dump, "0")) {
        return;
    }

    if (c == NULL || c->logf == NULL) {
        return;
    }

    fprintf(c->logf, "--- %s curl ---\n", op);
    fprintf(c->logf, "curl --location ");
    shell_quote(c->logf, c->url);
    fprintf(c->logf, " \\\n");
    fprintf(c->logf, "  -H 'Content-Type: application/json'");

    if (c->authenticated) {
        fprintf(c->logf, " \\\n");
        fprintf(c->logf, "  -H ");
        fprintf(c->logf, "'X-Session-Token: %s'", c->token);
    }

    fprintf(c->logf, " \\\n");
    fprintf(c->logf, "  --data-raw ");
    shell_quote(c->logf, body);
    fprintf(c->logf, "\n");
    fflush(c->logf);
}

static int bff_call(bff_client_t *c, const char *op, const char *query,
                    const char *vars, json_t **out, ulab_error_t *err) {
    CURL *curl;
    CURLcode ret;
    struct curl_slist *hdr;
    http_buf_t resp;
    char qesc[8192];
    char body[16384];
    long code;
    json_t *root;
    json_t *errors;
    json_error_t json_err;

    hdr = NULL;
    resp.buf = NULL;
    resp.len = 0;
    code = 0;
    root = NULL;
    errors = NULL;

    ulab_json_escape(query, qesc, sizeof(qesc));
    snprintf(body, sizeof(body), "{\"query\":\"%s\",\"variables\":%s}",
             qesc, vars ? vars : "{}");

    if (c->logf) {
        fprintf(c->logf, "--- %s request ---\n%s\n", op, body);
        fflush(c->logf);
    }

    curl = curl_easy_init();
    if (curl == NULL) {
        snprintf(err->msg, sizeof(err->msg), "curl init failed");
        return ULAB_ERR;
    }

    hdr = curl_slist_append(hdr, "Content-Type: application/json");

    if (c->authenticated) {
        char token_hdr[8192];

        snprintf(token_hdr, sizeof(token_hdr),
                 "X-Session-Token: %s", c->token);

        hdr = curl_slist_append(hdr, token_hdr);
    }

    curl_easy_setopt(curl, CURLOPT_URL, c->url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, hdr);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &resp);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 30L);

    bff_dump_curl(c, op, body);

    ret = curl_easy_perform(curl);
    if (ret != CURLE_OK) {
        snprintf(err->msg, sizeof(err->msg), "%s: HTTP request failed: %s",
                 op, curl_easy_strerror(ret));
        curl_slist_free_all(hdr);
        curl_easy_cleanup(curl);
        free(resp.buf);
        return ULAB_ERR;
    }

    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    curl_slist_free_all(hdr);
    curl_easy_cleanup(curl);

    if (c->logf) {
        fprintf(c->logf, "--- %s response %ld ---\n%s\n", op, code,
                resp.buf ? resp.buf : "");
        fflush(c->logf);
    }

    if (code < 200 || code >= 300) {
        snprintf(err->msg, sizeof(err->msg), "%s: HTTP %ld", op, code);
        free(resp.buf);
        return ULAB_ERR;
    }

    root = json_loads(resp.buf ? resp.buf : "", 0, &json_err);
    free(resp.buf);

    if (root == NULL) {
        snprintf(err->msg, sizeof(err->msg), "%s: invalid JSON: %s", op,
                 json_err.text);
        return ULAB_ERR;
    }

    errors = json_object_get(root, "errors");
    if (errors != NULL) {
        char *err_json;

        err_json = json_dumps(errors, JSON_COMPACT);
        snprintf(err->msg, sizeof(err->msg), "%s: GraphQL error: %s", op,
                 err_json ? err_json : "unknown");
        free(err_json);
        json_decref(root);
        return ULAB_ERR;
    }

    *out = root;

    return ULAB_OK;
}

int bff_login(bff_client_t *c,
              const char *identifier,
              const char *password,
              ulab_error_t *err) {

    char url[ULAB_MAX_PATH];
    char body[8192];
    char flow_id[1024];
    char session_token[4096];
    char token[4096];
    json_t *root;
    struct curl_slist *hdrs;
    char session_hdr[8192];

    root = NULL;
    hdrs = NULL;
    flow_id[0] = '\0';
    session_token[0] = '\0';
    token[0] = '\0';

    if (identifier == NULL || identifier[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "bff auth missing identifier");
        return ULAB_ERR;
    }

    if (password == NULL || password[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "bff auth missing password");
        return ULAB_ERR;
    }

    if (build_url(url, sizeof(url), c->pauth_url,
                  "/.api/self-service/login/api?refresh=false", err)) {
        return ULAB_ERR;
    }

    if (http_json_request("authFlow", url, "GET", NULL, NULL, &root, err)) {
        return ULAB_ERR;
    }

    if (json_get_path_str(root, "id", flow_id, sizeof(flow_id))) {
        snprintf(err->msg, sizeof(err->msg), "authFlow missing id");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);
    root = NULL;

    snprintf(body, sizeof(body),
             "{\"method\":\"password\","
             "\"password\":\"%s\","
             "\"identifier\":\"%s\"}",
             password, identifier);

    {
        char path[2048];
        int n;

        n = snprintf(path, sizeof(path),
                     "/.api/self-service/login?flow=%s", flow_id);
        if (n < 0 || (size_t)n >= sizeof(path)) {
            snprintf(err->msg, sizeof(err->msg),
                     "auth login flow URL path too long");
            return ULAB_ERR;
        }

        if (build_url(url, sizeof(url), c->pauth_url, path, err)) {
            return ULAB_ERR;
        }
    }

    if (http_json_request("authLogin", url, "POST", body, NULL, &root, err)) {
        return ULAB_ERR;
    }

    if (json_get_path_str(root, "session_token",
        session_token, sizeof(session_token))) {
        snprintf(err->msg, sizeof(err->msg), "authLogin missing session_token");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);
    root = NULL;

    snprintf(session_hdr, sizeof(session_hdr),
             "X-Session-Token: %s", session_token);
    hdrs = curl_slist_append(hdrs, session_hdr);

    if (build_url(url, sizeof(url), c->bff_base_url,
                  "/gateway/get-user", err)) {
        curl_slist_free_all(hdrs);
        return ULAB_ERR;
    }

    if (http_json_request("getUser", url, "GET", NULL, hdrs, &root, err)) {
        curl_slist_free_all(hdrs);
        return ULAB_ERR;
    }

    curl_slist_free_all(hdrs);

    if (json_get_path_str(root, "token", token, sizeof(token))) {
        snprintf(err->msg, sizeof(err->msg), "getUser missing token");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);

    ulab_copy(c->session_token, sizeof(c->session_token), session_token);
    ulab_copy(c->token, sizeof(c->token), token);
    c->authenticated = ULAB_TRUE;

    if (c->logf) {
        fprintf(c->logf, "--- bff auth ---\n");
        fprintf(c->logf, "authenticated identifier=%s\n", identifier);
        fflush(c->logf);
    }

    return ULAB_OK;
}

int bff_init(bff_client_t *c, const char *url, const char *run_dir) {

    char path[ULAB_MAX_PATH];
    const char *identifier;
    const char *password;
    ulab_error_t err;

    memset(c, 0, sizeof(*c));
    ulab_copy(c->url, sizeof(c->url), url);

    ulab_copy(c->pauth_url, sizeof(c->pauth_url),
              env_or_default("PAUTH_URL", "https://pauth.udev.ukama.com"));

    derive_bff_base_url(url, c->bff_base_url, sizeof(c->bff_base_url));

    if (getenv("BFF_BASE_URL") != NULL && getenv("BFF_BASE_URL")[0] != '\0') {
        ulab_copy(c->bff_base_url, sizeof(c->bff_base_url),
                  getenv("BFF_BASE_URL"));
    }

    snprintf(path, sizeof(path), "%s/bff.log", run_dir);
    c->logf = fopen(path, "w");

    curl_global_init(CURL_GLOBAL_ALL);

    if (getenv("UKAMA_SESSION_TOKEN") != NULL &&
        getenv("UKAMA_BFF_TOKEN") != NULL &&
        getenv("UKAMA_SESSION_TOKEN")[0] != '\0' &&
        getenv("UKAMA_BFF_TOKEN")[0] != '\0') {
        ulab_copy(c->session_token, sizeof(c->session_token),
                  getenv("UKAMA_SESSION_TOKEN"));
        ulab_copy(c->token, sizeof(c->token), getenv("UKAMA_BFF_TOKEN"));
        c->authenticated = ULAB_TRUE;

        return ULAB_OK;
    }

    identifier = getenv("UKAMA_IDENTIFIER");
    password = getenv("UKAMA_PASSWORD");

    if (identifier != NULL && identifier[0] != '\0' &&
        password != NULL && password[0] != '\0') {
        memset(&err, 0, sizeof(err));
        if (bff_login(c, identifier, password, &err)) {
            if (c->logf) {
                fprintf(c->logf, "--- bff auth failed ---\n%s\n", err.msg);
                fflush(c->logf);
            }
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}

void bff_close(bff_client_t *c) {
    if (c->logf) {
        fclose(c->logf);
    }

    curl_global_cleanup();
}

int bff_add_network(bff_client_t *c, network_t *n, ulab_error_t *err) {
    char vars[4096];
    json_t *root;
    json_t *obj;

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"name\":\"%s\","
             "\"budget\":0,"
             "\"countries\":[\"USA\"],"
             "\"networks\":[\"A3\"]}}",
             n->name);

    if (bff_call(c, "addNetwork", BFF_ADD_NETWORK, vars, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "addNetwork");
    if (obj == NULL || json_get_str(obj, "id", n->bff_id,
        sizeof(n->bff_id))) {
        snprintf(err->msg, sizeof(err->msg), "addNetwork missing id");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);

    return ULAB_OK;
}

static int bff_component_query(bff_client_t *c,
                               const char *category,
                               json_t **out,
                               ulab_error_t *err) {
    char vars[512];

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"category\":\"%s\"}}", category);

    return bff_call(c, "getComponentsByUserId",
                    BFF_GET_COMPONENTS_BY_USER_ID, vars, out, err);
}

static int bff_pick_first_component(bff_client_t *c,
                                    const char *category,
                                    char *out,
                                    size_t out_len,
                                    ulab_error_t *err) {
    json_t *root;
    json_t *obj;
    json_t *arr;
    json_t *it;
    json_t *idv;
    const char *id;

    root = NULL;
    obj = NULL;
    arr = NULL;
    it = NULL;
    idv = NULL;
    id = NULL;

    if (bff_component_query(c, category, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "getComponentsByUserId");
    arr = obj ? json_object_get(obj, "components") : NULL;

    if (arr == NULL || !json_is_array(arr) || json_array_size(arr) == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "no inventory component found for category=%s", category);
        json_decref(root);
        return ULAB_ERR;
    }

    it = json_array_get(arr, 0);
    idv = it ? json_object_get(it, "id") : NULL;
    if (idv == NULL || !json_is_string(idv)) {
        snprintf(err->msg, sizeof(err->msg),
                 "component missing id for category=%s", category);
        json_decref(root);
        return ULAB_ERR;
    }

    id = json_string_value(idv);
    if (id == NULL || id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "empty component id for category=%s", category);
        json_decref(root);
        return ULAB_ERR;
    }

    ulab_copy(out, out_len, id);
    json_decref(root);

    return ULAB_OK;
}

static int bff_pick_access_component_for_site(bff_client_t *c,
                                              const site_t *site,
                                              ulab_error_t *err) {
    json_t *root;
    json_t *obj;
    json_t *arr;
    json_t *it;
    json_t *idv;
    json_t *pnv;
    const char *id;
    const char *part;
    size_t i;

    root = NULL;
    obj = NULL;
    arr = NULL;
    it = NULL;
    idv = NULL;
    pnv = NULL;
    id = NULL;
    part = NULL;

    if (site == NULL || site->tnode_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "site has no tower node id");
        return ULAB_ERR;
    }

    if (bff_component_query(c, "access", &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "getComponentsByUserId");
    arr = obj ? json_object_get(obj, "components") : NULL;

    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getComponentsByUserId missing access components");
        json_decref(root);
        return ULAB_ERR;
    }

    for (i = 0; i < json_array_size(arr); i++) {
        it = json_array_get(arr, i);
        if (it == NULL || !json_is_object(it)) {
            continue;
        }

        pnv = json_object_get(it, "partNumber");
        if (pnv == NULL || !json_is_string(pnv)) {
            continue;
        }

        part = json_string_value(pnv);
        if (part == NULL || !ulab_streq(part, site->tnode_id)) {
            continue;
        }

        idv = json_object_get(it, "id");
        if (idv == NULL || !json_is_string(idv)) {
            continue;
        }

        id = json_string_value(idv);
        if (id == NULL || id[0] == '\0') {
            continue;
        }

        ulab_copy(c->access_id, sizeof(c->access_id), id);
        json_decref(root);
        return ULAB_OK;
    }

    snprintf(err->msg, sizeof(err->msg),
             "no access component found for partNumber=%s",
             site->tnode_id);
    json_decref(root);
    return ULAB_ERR;
}

static int bff_load_site_components(bff_client_t *c,
                                    const site_t *site,
                                    ulab_error_t *err) {

    if (bff_pick_access_component_for_site(c, site, err)) {
        return ULAB_ERR;
    }

    if (bff_pick_first_component(c, "backhaul",
        c->backhaul_id, sizeof(c->backhaul_id), err)) {
        return ULAB_ERR;
    }

    if (bff_pick_first_component(c, "power",
        c->power_id, sizeof(c->power_id), err)) {
        return ULAB_ERR;
    }

    if (bff_pick_first_component(c, "spectrum",
        c->spectrum_id, sizeof(c->spectrum_id), err)) {
        return ULAB_ERR;
    }

    if (bff_pick_first_component(c, "switch",
        c->switch_id, sizeof(c->switch_id), err)) {
        return ULAB_ERR;
    }

    return ULAB_OK;
}

int bff_wait_site_anchor_online(bff_client_t *c,
                                site_t *site,
                                ulab_error_t *err) {
    const char *timeout_env;
    const char *sleep_env;
    unsigned int timeout_sec;
    unsigned int sleep_sec;
    unsigned int elapsed;
    char vars[1024];
    json_t *root;
    json_t *obj;
    json_t *arr;
    json_t *it;
    json_t *idv;
    char state[ULAB_MAX_REF];
    char connectivity[ULAB_MAX_REF];
    const char *id;
    size_t i;
    int found;

    if (site == NULL || site->tnode_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "cannot wait for site without tower node id");
        return ULAB_ERR;
    }

    timeout_env = getenv("ULAB_BFF_NODE_ONLINE_TIMEOUT_SEC");
    sleep_env = getenv("ULAB_BFF_NODE_ONLINE_SLEEP_SEC");
    timeout_sec = timeout_env ? (unsigned int)strtoul(timeout_env, NULL, 10) :
        180u;
    sleep_sec = sleep_env ? (unsigned int)strtoul(sleep_env, NULL, 10) : 5u;
    if (timeout_sec == 0) timeout_sec = 180u;
    if (sleep_sec == 0) sleep_sec = 5u;

    elapsed = 0;

    while (elapsed <= timeout_sec) {
        snprintf(vars, sizeof(vars),
                 "{\"data\":{\"type\":\"%s\"}}",
                 ULAB_NODE_KIND_TOWER);

        root = NULL;
        if (bff_call(c, "getNodes", BFF_GET_NODES, vars, &root, err)) {
            return ULAB_ERR;
        }

        obj = dig(root, "data", "getNodes");
        arr = obj ? json_object_get(obj, "nodes") : NULL;
        found = 0;

        if (arr != NULL && json_is_array(arr)) {
            for (i = 0; i < json_array_size(arr); i++) {
                it = json_array_get(arr, i);
                idv = it ? json_object_get(it, "id") : NULL;
                id = idv && json_is_string(idv) ? json_string_value(idv) :
                    NULL;

                if (id == NULL || !ulab_streq(id, site->tnode_id)) {
                    continue;
                }

                found = 1;
                state[0] = '\0';
                connectivity[0] = '\0';

                json_get_nested_str(it, "status", "state", state,
                                    sizeof(state));
                json_get_nested_str(it, "status", "connectivity",
                                    connectivity, sizeof(connectivity));
                json_get_nested_str(it, "latitude", NULL, site->latitude,
                                    sizeof(site->latitude));
                json_get_nested_str(it, "longitude", NULL, site->longitude,
                                    sizeof(site->longitude));

                if (ulab_streq(connectivity, "Online") &&
                    ulab_streq(state, "Unknown") &&
                    site->latitude[0] != '\0' &&
                    site->longitude[0] != '\0') {
                    json_decref(root);
                    ulab_status("BACKEND", "site %s anchor online %s lat=%s lng=%s",
                                site->ref, site->tnode_id, site->latitude,
                                site->longitude);
                    return ULAB_OK;
                }
            }
        }

        json_decref(root);

        if (found) {
            ulab_status("BACKEND", "waiting site %s anchor online/location",
                        site->ref);
        } else {
            ulab_status("BACKEND", "waiting site %s anchor in registry",
                        site->ref);
        }

        sleep(sleep_sec);
        elapsed += sleep_sec;
    }

    snprintf(err->msg, sizeof(err->msg),
             "site %s tnode %s did not become Online/Unknown with lat/lng",
             site->ref, site->tnode_id);
    return ULAB_ERR;
}

int bff_add_site(bff_client_t *c,
                 site_t *s,
                 const network_t *n,
                 ulab_error_t *err) {
    char vars[4096];
    json_t *root;
    json_t *obj;

    if (s == NULL || s->tnode_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "site has no tower node id");
        return ULAB_ERR;
    }

    if (s->latitude[0] == '\0' || s->longitude[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "site %s missing anchor lat/lng", s->ref);
        return ULAB_ERR;
    }

    if (bff_load_site_components(c, s, err)) {
        return ULAB_ERR;
    }

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"name\":\"%s\",\"network_id\":\"%s\","
             "\"backhaul_id\":\"%s\","
             "\"power_id\":\"%s\","
             "\"access_id\":\"%s\","
             "\"spectrum_id\":\"%s\","
             "\"switch_id\":\"%s\","
             "\"latitude\":\"%s\","
             "\"longitude\":\"%s\","
             "\"install_date\":\"2026-01-01T00:00:00Z\","
             "\"location\":\"Lab\"}}",
             s->name,
             n->bff_id,
             c->backhaul_id,
             c->power_id,
             c->access_id,
             c->spectrum_id,
             c->switch_id,
             s->latitude,
             s->longitude);

    if (bff_call(c, "addSite", BFF_ADD_SITE, vars, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "addSite");
    if (obj == NULL || json_get_str(obj, "id", s->bff_id,
        sizeof(s->bff_id))) {
        snprintf(err->msg, sizeof(err->msg), "addSite missing id");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);

    return ULAB_OK;
}

int bff_add_package(bff_client_t *c,
                    package_t *p,
                    const network_t *net,
                    ulab_error_t *err) {
    char vars[4096];
    json_t *root;
    json_t *obj;
    uint32_t duration;
    uint64_t duration64;

    if (net != NULL && net->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "addPackage has unresolved network id");
        return ULAB_ERR;
    }

    if (p->duration_minutes > 0) {
        duration64 = p->duration_minutes;
    } else {
        duration64 = (uint64_t)p->duration_days * 1440u;
    }
    if (duration64 == 0 || duration64 > UINT32_MAX) {
        snprintf(err->msg, sizeof(err->msg),
                 "addPackage invalid duration");
        return ULAB_ERR;
    }
    duration = (uint32_t)duration64;

    if (net != NULL) {
        snprintf(vars, sizeof(vars),
                 "{\"data\":{\"name\":\"%s\",\"amount\":%.2f,"
                 "\"dataUnit\":\"MB\",\"dataVolume\":%llu,"
                 "\"duration\":%u,\"currency\":\"%s\","
                 "\"country\":\"%s\",\"networkId\":\"%s\"}}",
                 p->name, p->amount,
                 (unsigned long long)p->data_mb, duration,
                 p->currency[0] ? p->currency : "USD",
                 p->country[0] ? p->country : "USA", net->bff_id);
    } else {
        snprintf(vars, sizeof(vars),
                 "{\"data\":{\"name\":\"%s\",\"amount\":%.2f,"
                 "\"dataUnit\":\"MB\",\"dataVolume\":%llu,"
                 "\"duration\":%u,\"currency\":\"%s\","
                 "\"country\":\"%s\"}}",
                 p->name, p->amount,
                 (unsigned long long)p->data_mb, duration,
                 p->currency[0] ? p->currency : "USD",
                 p->country[0] ? p->country : "USA");
    }

    if (bff_call(c, "addPackage", BFF_ADD_PACKAGE, vars, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "addPackage");
    if (obj == NULL || json_get_str(obj, "uuid", p->bff_id,
        sizeof(p->bff_id))) {
        snprintf(err->msg, sizeof(err->msg), "addPackage missing uuid");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);

    return ULAB_OK;
}

int bff_set_package_active(bff_client_t *c,
                           package_t *p,
                           int active,
                           ulab_error_t *err) {
    char vars[4096];
    json_t *root;
    json_t *obj;
    json_t *value;

    if (p == NULL || p->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "updatePackage missing package id");
        return ULAB_ERR;
    }

    snprintf(vars, sizeof(vars),
             "{\"packageId\":\"%s\",\"data\":{\"name\":\"%s\","
             "\"active\":%s}}",
             p->bff_id, p->name, active ? "true" : "false");

    if (bff_call(c, "updatePackage", BFF_UPDATE_PACKAGE, vars,
                 &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "updatePackage");
    value = obj ? json_object_get(obj, "active") : NULL;
    if (value == NULL || !json_is_boolean(value) ||
        json_is_true(value) != (active != 0)) {
        snprintf(err->msg, sizeof(err->msg),
                 "updatePackage returned unexpected active state");
        json_decref(root);
        return ULAB_ERR;
    }

    p->active = active != 0;
    json_decref(root);
    return ULAB_OK;
}

static double package_json_number(json_t *obj, const char *key) {
    json_t *value;

    value = obj ? json_object_get(obj, key) : NULL;
    return value != NULL && json_is_number(value) ?
        json_number_value(value) : 0.0;
}

static void invalid_package_name(const package_t *pkg,
                                 const char *variant,
                                 char *out,
                                 size_t out_len) {
    snprintf(out, out_len, "%.180s invalid %.48s", pkg->name, variant);
}

int bff_get_package(bff_client_t *c,
                    const package_t *pkg,
                    bff_package_t *actual,
                    ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char id_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *value;

    if (pkg == NULL || pkg->bff_id[0] == '\0' || actual == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPackage requires package and output storage");
        return ULAB_ERR;
    }
    ulab_json_escape(pkg->bff_id, id_esc, sizeof(id_esc));
    snprintf(vars, sizeof(vars), "{\"packageId\":\"%s\"}", id_esc);

    root = NULL;
    if (bff_call(c, "getPackage", BFF_GET_PACKAGE, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getPackage");
    memset(actual, 0, sizeof(*actual));
    if (obj == NULL ||
        json_get_str(obj, "uuid", actual->uuid, sizeof(actual->uuid)) ||
        json_get_str(obj, "name", actual->name, sizeof(actual->name)) ||
        json_get_str(obj, "dataUnit", actual->data_unit,
                     sizeof(actual->data_unit)) ||
        json_get_str(obj, "currency", actual->currency,
                     sizeof(actual->currency)) ||
        json_get_str(obj, "country", actual->country,
                     sizeof(actual->country))) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPackage missing required fields");
        json_decref(root);
        return ULAB_ERR;
    }
    actual->data_volume = (uint64_t)package_json_number(obj, "dataVolume");
    actual->duration_minutes =
        (uint32_t)package_json_number(obj, "duration");
    actual->amount = package_json_number(obj, "amount");
    value = json_object_get(obj, "active");
    actual->active = value != NULL && json_is_true(value);
    json_get_optional_str(obj, "networkId", actual->network_id,
                          sizeof(actual->network_id));
    json_decref(root);
    return ULAB_OK;
}

int bff_package_visible_for_network(bff_client_t *c,
                                    const package_t *pkg,
                                    const network_t *network,
                                    int *visible,
                                    ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *arr;
    size_t i;

    if (pkg == NULL || pkg->bff_id[0] == '\0' || network == NULL ||
        network->bff_id[0] == '\0' || visible == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPackages visibility requires package and network");
        return ULAB_ERR;
    }
    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
             network_esc);
    root = NULL;
    if (bff_call(c, "getPackages", BFF_GET_PACKAGES, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getPackages");
    arr = obj ? json_object_get(obj, "packages") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPackages missing packages list");
        json_decref(root);
        return ULAB_ERR;
    }
    *visible = 0;
    for (i = 0; i < json_array_size(arr); i++) {
        char id[ULAB_MAX_ID];

        if (json_get_str(json_array_get(arr, i), "uuid", id,
                         sizeof(id)) == ULAB_OK &&
            ulab_streq(id, pkg->bff_id)) {
            *visible = 1;
            break;
        }
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_invalid_package_name_available(bff_client_t *c,
                                       const package_t *pkg,
                                       const char *variant,
                                       int *available,
                                       ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char name[ULAB_MAX_NAME];
    char name_esc[ULAB_MAX_NAME * 2];
    json_t *root;
    json_t *obj;
    json_t *value;

    if (pkg == NULL || variant == NULL || variant[0] == '\0' ||
        available == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "package name availability requires package and variant");
        return ULAB_ERR;
    }
    invalid_package_name(pkg, variant, name, sizeof(name));
    ulab_json_escape(name, name_esc, sizeof(name_esc));
    snprintf(vars, sizeof(vars), "{\"name\":\"%s\"}", name_esc);
    root = NULL;
    if (bff_call(c, "isPackageNameAvailable", BFF_PACKAGE_NAME_AVAILABLE,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "isPackageNameAvailable");
    value = obj ? json_object_get(obj, "isAvailable") : NULL;
    if (value == NULL || !json_is_boolean(value)) {
        snprintf(err->msg, sizeof(err->msg),
                 "isPackageNameAvailable missing boolean result");
        json_decref(root);
        return ULAB_ERR;
    }
    *available = json_is_true(value);
    json_decref(root);
    return ULAB_OK;
}

int bff_add_invalid_package(bff_client_t *c,
                            const package_t *pkg,
                            const network_t *network,
                            const char *variant,
                            char *created_id,
                            size_t created_id_len,
                            ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char name[ULAB_MAX_NAME];
    char name_esc[ULAB_MAX_NAME * 2];
    char network_esc[ULAB_MAX_ID * 2];
    const char *currency;
    long long data_volume;
    long long duration;
    double amount;
    json_t *root;
    json_t *obj;

    if (pkg == NULL || network == NULL || network->bff_id[0] == '\0' ||
        variant == NULL || variant[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "invalid package request requires package, network and "
                 "variant");
        return ULAB_ERR;
    }
    data_volume = (long long)pkg->data_mb;
    duration = pkg->duration_minutes > 0 ?
        (long long)pkg->duration_minutes :
        (long long)pkg->duration_days * 1440LL;
    amount = pkg->amount;
    currency = pkg->currency[0] ? pkg->currency : "USD";
    if (ulab_streq(variant, "allowance")) data_volume = -1;
    else if (ulab_streq(variant, "duration")) duration = 0;
    else if (ulab_streq(variant, "price")) amount = -1.0;
    else if (ulab_streq(variant, "currency")) currency = "";
    else {
        snprintf(err->msg, sizeof(err->msg),
                 "unknown invalid package variant %.64s", variant);
        return ULAB_ERR;
    }
    invalid_package_name(pkg, variant, name, sizeof(name));
    ulab_json_escape(name, name_esc, sizeof(name_esc));
    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars),
             "{\"data\":{\"name\":\"%s\",\"amount\":%.2f,"
             "\"dataUnit\":\"MB\",\"dataVolume\":%lld,"
             "\"duration\":%lld,\"currency\":\"%s\","
             "\"country\":\"%s\",\"networkId\":\"%s\"}}",
             name_esc, amount, data_volume, duration, currency,
             pkg->country[0] ? pkg->country : "USA", network_esc);
    root = NULL;
    if (bff_call(c, "addPackage", BFF_ADD_PACKAGE, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "addPackage");
    if (obj == NULL || created_id == NULL || created_id_len == 0 ||
        json_get_str(obj, "uuid", created_id, created_id_len)) {
        snprintf(err->msg, sizeof(err->msg),
                 "unexpected accepted invalid package missing uuid");
        json_decref(root);
        return ULAB_ERR;
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_get_package_metrics(bff_client_t *c,
                            const package_t *pkg,
                            const network_t *network,
                            bff_package_metrics_t *metrics,
                            int *found,
                            ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *view;
    json_t *plans_section;
    json_t *plans;
    size_t i;

    if (pkg == NULL || pkg->bff_id[0] == '\0' || network == NULL ||
        metrics == NULL || found == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "PackagesDashboard requires package, network and output");
        return ULAB_ERR;
    }
    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
             network_esc);
    root = NULL;
    if (bff_call(c, "PackagesDashboard", BFF_PACKAGES_DASHBOARD,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    view = dig(root, "data", "commerceView");
    plans_section = view ? json_object_get(view, "plans") : NULL;
    plans = plans_section ? json_object_get(plans_section, "plans") : NULL;
    if (plans == NULL || !json_is_array(plans)) {
        snprintf(err->msg, sizeof(err->msg),
                 "PackagesDashboard missing plans list");
        json_decref(root);
        return ULAB_ERR;
    }
    *found = 0;
    memset(metrics, 0, sizeof(*metrics));
    for (i = 0; i < json_array_size(plans); i++) {
        json_t *plan;
        json_t *attach;
        char id[ULAB_MAX_ID];

        plan = json_array_get(plans, i);
        if (json_get_str(plan, "packageId", id, sizeof(id)) ||
            !ulab_streq(id, pkg->bff_id)) {
            continue;
        }
        ulab_copy(metrics->package_id, sizeof(metrics->package_id), id);
        metrics->revenue = package_json_number(plan, "revenue");
        attach = json_object_get(plan, "attachCount");
        if (attach != NULL && json_is_number(attach)) {
            metrics->attach_count = (uint32_t)json_number_value(attach);
            metrics->has_attach_count = 1;
        }
        *found = 1;
        break;
    }
    json_decref(root);
    return ULAB_OK;
}

static int json_u32_field(json_t *obj, const char *key, uint32_t *out) {
    json_t *value;
    double number;

    value = obj ? json_object_get(obj, key) : NULL;
    if (value == NULL || !json_is_number(value) || out == NULL) {
        return ULAB_ERR;
    }
    number = json_number_value(value);
    if (number < 0 || number > 4294967295.0) {
        return ULAB_ERR;
    }
    *out = (uint32_t)number;
    return ULAB_OK;
}

static int json_double_field(json_t *obj, const char *key, double *out) {
    json_t *value;

    value = obj ? json_object_get(obj, key) : NULL;
    if (value == NULL || !json_is_number(value) || out == NULL) {
        return ULAB_ERR;
    }
    *out = json_number_value(value);
    return ULAB_OK;
}

static int section_has_error(json_t *section) {
    json_t *error;

    error = section ? json_object_get(section, "error") : NULL;
    return error != NULL && json_is_object(error);
}

int bff_get_revenue_summary(bff_client_t *c,
                            const network_t *network,
                            bff_revenue_summary_t *summary,
                            ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *view;
    json_t *revenue;

    if (network == NULL || network->bff_id[0] == '\0' || summary == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "RevenueOverview requires network and output");
        return ULAB_ERR;
    }
    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}", network_esc);
    root = NULL;
    if (bff_call(c, "RevenueOverview", BFF_REVENUE_OVERVIEW,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    view = dig(root, "data", "commerceView");
    revenue = view ? json_object_get(view, "revenue") : NULL;
    memset(summary, 0, sizeof(*summary));
    if (revenue == NULL || section_has_error(revenue) ||
        json_double_field(revenue, "totalPaid", &summary->total_paid) ||
        json_double_field(revenue, "totalPending",
                          &summary->total_pending) ||
        json_double_field(revenue, "monthPaid", &summary->month_paid) ||
        json_double_field(revenue, "prevMonthPaid",
                          &summary->previous_month_paid)) {
        snprintf(err->msg, sizeof(err->msg),
                 "RevenueOverview returned an incomplete revenue section");
        json_decref(root);
        return ULAB_ERR;
    }
    if (json_double_field(revenue, "momPct",
                          &summary->month_over_month_percent)) {
        summary->month_over_month_percent = 0;
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_get_package_dashboard(bff_client_t *c,
                              const network_t *network,
                              bff_package_dashboard_t *dashboard,
                              ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *view;
    json_t *plans;

    if (network == NULL || network->bff_id[0] == '\0' || dashboard == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "PackagesDashboard requires network and output");
        return ULAB_ERR;
    }
    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}", network_esc);
    root = NULL;
    if (bff_call(c, "PackagesDashboard", BFF_PACKAGES_DASHBOARD,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    view = dig(root, "data", "commerceView");
    plans = view ? json_object_get(view, "plans") : NULL;
    memset(dashboard, 0, sizeof(*dashboard));
    if (plans == NULL || section_has_error(plans)) {
        snprintf(err->msg, sizeof(err->msg),
                 "PackagesDashboard returned a plans error");
        json_decref(root);
        return ULAB_ERR;
    }
    dashboard->has_mrr =
        json_double_field(plans, "mrr", &dashboard->mrr) == ULAB_OK;
    dashboard->has_arpu =
        json_double_field(plans, "arpu", &dashboard->arpu) == ULAB_OK;
    json_decref(root);
    return ULAB_OK;
}

int bff_get_network_overview(bff_client_t *c,
                             const network_t *network,
                             bff_network_overview_t *overview,
                             ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *view;
    json_t *subscribers;
    json_t *sites;
    json_t *site_list;
    json_t *nodes;

    if (network == NULL || network->bff_id[0] == '\0' || overview == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "NetworkHome requires network and output");
        return ULAB_ERR;
    }
    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}", network_esc);
    root = NULL;
    if (bff_call(c, "NetworkHome", BFF_NETWORK_HOME, vars, &root, err)) {
        return ULAB_ERR;
    }
    view = dig(root, "data", "networkOverview");
    subscribers = view ? json_object_get(view, "subscriberStats") : NULL;
    sites = view ? json_object_get(view, "siteStats") : NULL;
    site_list = sites ? json_object_get(sites, "sites") : NULL;
    nodes = view ? json_object_get(view, "nodeStats") : NULL;
    memset(overview, 0, sizeof(*overview));
    if (subscribers == NULL || sites == NULL || nodes == NULL ||
        section_has_error(subscribers) || section_has_error(sites) ||
        section_has_error(nodes) || site_list == NULL ||
        !json_is_array(site_list) ||
        json_u32_field(subscribers, "total",
                       &overview->subscribers_total) ||
        json_u32_field(subscribers, "active",
                       &overview->subscribers_active) ||
        json_u32_field(subscribers, "inactive",
                       &overview->subscribers_inactive) ||
        json_u32_field(nodes, "total", &overview->nodes_total) ||
        json_u32_field(nodes, "online", &overview->nodes_online) ||
        json_u32_field(nodes, "offline", &overview->nodes_offline)) {
        snprintf(err->msg, sizeof(err->msg),
                 "NetworkHome returned incomplete summary sections");
        json_decref(root);
        return ULAB_ERR;
    }
    overview->sites_total = (uint32_t)json_array_size(site_list);
    json_decref(root);
    return ULAB_OK;
}

int bff_get_nodes_view_count(bff_client_t *c,
                             const network_t *network,
                             uint32_t *count,
                             ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *view;
    json_t *section;
    json_t *nodes;

    if (network == NULL || network->bff_id[0] == '\0' || count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "NodesList requires network and output");
        return ULAB_ERR;
    }
    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}", network_esc);
    root = NULL;
    if (bff_call(c, "NodesList", BFF_NODES_LIST, vars, &root, err)) {
        return ULAB_ERR;
    }
    view = dig(root, "data", "nodesView");
    section = view ? json_object_get(view, "nodes") : NULL;
    nodes = section ? json_object_get(section, "nodes") : NULL;
    if (section == NULL || section_has_error(section) || nodes == NULL ||
        !json_is_array(nodes)) {
        snprintf(err->msg, sizeof(err->msg),
                 "NodesList returned an incomplete nodes section");
        json_decref(root);
        return ULAB_ERR;
    }
    *count = (uint32_t)json_array_size(nodes);
    json_decref(root);
    return ULAB_OK;
}

int bff_get_inventory_summary(bff_client_t *c,
                              const char *sim_type,
                              bff_inventory_summary_t *summary,
                              ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char sim_type_esc[ULAB_MAX_REF * 2];
    json_t *root;
    json_t *view;
    json_t *components;
    json_t *categories;
    json_t *stock;
    json_t *pool;
    json_t *pool_stats;
    size_t i;

    if (sim_type == NULL || sim_type[0] == '\0' || summary == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "inventory reconciliation requires SIM type and output");
        return ULAB_ERR;
    }
    memset(summary, 0, sizeof(*summary));
    root = NULL;
    if (bff_call(c, "InventoryOverview", BFF_INVENTORY_OVERVIEW,
                 "{}", &root, err)) {
        return ULAB_ERR;
    }
    view = dig(root, "data", "inventoryView");
    components = view ? json_object_get(view, "components") : NULL;
    categories = components ? json_object_get(components, "byCategory") :
        NULL;
    stock = view ? json_object_get(view, "simStock") : NULL;
    if (components == NULL || stock == NULL ||
        section_has_error(components) || section_has_error(stock) ||
        categories == NULL || !json_is_array(categories) ||
        json_u32_field(components, "total", &summary->component_total) ||
        json_u32_field(stock, "total", &summary->sim_total) ||
        json_u32_field(stock, "available", &summary->sim_available) ||
        json_u32_field(stock, "consumed", &summary->sim_consumed)) {
        snprintf(err->msg, sizeof(err->msg),
                 "InventoryOverview returned incomplete inventory sections");
        json_decref(root);
        return ULAB_ERR;
    }
    for (i = 0; i < json_array_size(categories); i++) {
        uint32_t count;

        count = 0;
        if (json_u32_field(json_array_get(categories, i), "count", &count)) {
            snprintf(err->msg, sizeof(err->msg),
                     "InventoryOverview returned an invalid category count");
            json_decref(root);
            return ULAB_ERR;
        }
        summary->component_category_total += count;
    }
    json_decref(root);

    ulab_json_escape(sim_type, sim_type_esc, sizeof(sim_type_esc));
    snprintf(vars, sizeof(vars),
             "{\"simType\":\"%s\",\"limit\":1}", sim_type_esc);
    root = NULL;
    if (bff_call(c, "SimPoolOverview", BFF_SIM_POOL_OVERVIEW,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    pool = dig(root, "data", "simPoolView");
    pool_stats = pool ? json_object_get(pool, "stats") : NULL;
    if (pool_stats == NULL || section_has_error(pool_stats) ||
        json_u32_field(pool_stats, "total", &summary->sim_pool_total) ||
        json_u32_field(pool_stats, "available",
                       &summary->sim_pool_available) ||
        json_u32_field(pool_stats, "consumed",
                       &summary->sim_pool_consumed)) {
        snprintf(err->msg, sizeof(err->msg),
                 "SimPoolOverview returned an incomplete stats section");
        json_decref(root);
        return ULAB_ERR;
    }
    json_decref(root);
    return ULAB_OK;
}

static int payment_status_settled(const char *status) {
    return status != NULL &&
        (ulab_streq(status, "success") || ulab_streq(status, "paid") ||
         ulab_streq(status, "completed") ||
         ulab_streq(status, "processed"));
}

int bff_get_subscriber_billing(bff_client_t *c,
                               const subscriber_t *subscriber,
                               bff_subscriber_billing_t *billing,
                               ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char subscriber_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *view;
    json_t *section;
    json_t *payments;
    size_t i;

    if (subscriber == NULL || subscriber->bff_id[0] == '\0' ||
        billing == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "SubscriberDetail requires subscriber and output");
        return ULAB_ERR;
    }
    ulab_json_escape(subscriber->bff_id, subscriber_esc,
                     sizeof(subscriber_esc));
    snprintf(vars, sizeof(vars), "{\"subscriberId\":\"%s\"}",
             subscriber_esc);
    root = NULL;
    if (bff_call(c, "SubscriberDetail", BFF_SUBSCRIBER_DETAIL,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    view = dig(root, "data", "subscriberView");
    section = view ? json_object_get(view, "billing") : NULL;
    payments = section ? json_object_get(section, "payments") : NULL;
    if (section == NULL || section_has_error(section) || payments == NULL ||
        !json_is_array(payments)) {
        snprintf(err->msg, sizeof(err->msg),
                 "SubscriberDetail returned an incomplete billing section");
        json_decref(root);
        return ULAB_ERR;
    }
    memset(billing, 0, sizeof(*billing));
    billing->payment_count = (uint32_t)json_array_size(payments);
    for (i = 0; i < json_array_size(payments); i++) {
        json_t *payment;
        json_t *status_value;
        json_t *amount_value;
        const char *status;
        const char *amount;

        payment = json_array_get(payments, i);
        status_value = payment ? json_object_get(payment, "status") : NULL;
        amount_value = payment ? json_object_get(payment, "amount") : NULL;
        status = status_value && json_is_string(status_value) ?
            json_string_value(status_value) : NULL;
        amount = amount_value && json_is_string(amount_value) ?
            json_string_value(amount_value) : NULL;
        if (status == NULL || amount == NULL) {
            snprintf(err->msg, sizeof(err->msg),
                     "SubscriberDetail returned an invalid payment");
            json_decref(root);
            return ULAB_ERR;
        }
        if (payment_status_settled(status)) {
            billing->settled_count++;
            billing->settled_amount += strtod(amount, NULL);
        }
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_sim_is_unallocated(bff_client_t *c,
                           const ue_t *ue,
                           const char *sim_type,
                           int *unallocated,
                           ulab_error_t *err) {
    char (*iccids)[ULAB_MAX_ID];
    char (*ids)[ULAB_MAX_ID];
    size_t count;
    size_t i;
    size_t max_sims;
    int rc;

    if (ue == NULL || ue->iccid[0] == '\0' || sim_type == NULL ||
        sim_type[0] == '\0' || unallocated == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "SIM pool check requires UE and sim type");
        return ULAB_ERR;
    }
    max_sims = 4096;
    iccids = calloc(max_sims, sizeof(*iccids));
    ids = calloc(max_sims, sizeof(*ids));
    if (iccids == NULL || ids == NULL) {
        free(iccids);
        free(ids);
        snprintf(err->msg, sizeof(err->msg),
                 "out of memory checking SIM pool");
        return ULAB_ERR;
    }
    count = 0;
    rc = bff_get_sims_from_pool(c, sim_type, iccids, ids, max_sims,
                                &count, err);
    if (rc != ULAB_OK) {
        free(iccids);
        free(ids);
        return rc;
    }
    *unallocated = 0;
    for (i = 0; i < count; i++) {
        if (ulab_streq(ue->iccid, iccids[i])) {
            *unallocated = 1;
            break;
        }
    }
    free(iccids);
    free(ids);
    return ULAB_OK;
}

static int parse_payment(json_t *obj, bff_payment_t *payment,
                         ulab_error_t *err) {
    if (obj == NULL || payment == NULL) {
        snprintf(err->msg, sizeof(err->msg), "payment response is empty");
        return ULAB_ERR;
    }

    memset(payment, 0, sizeof(*payment));
    if (json_get_str(obj, "id", payment->id, sizeof(payment->id)) ||
        json_get_str(obj, "itemId", payment->item_id,
                     sizeof(payment->item_id)) ||
        json_get_str(obj, "itemType", payment->item_type,
                     sizeof(payment->item_type)) ||
        json_get_str(obj, "amount", payment->amount,
                     sizeof(payment->amount)) ||
        json_get_str(obj, "currency", payment->currency,
                     sizeof(payment->currency)) ||
        json_get_str(obj, "paymentMethod", payment->payment_method,
                     sizeof(payment->payment_method)) ||
        json_get_str(obj, "status", payment->status,
                     sizeof(payment->status))) {
        snprintf(err->msg, sizeof(err->msg),
                 "payment response is missing required fields");
        return ULAB_ERR;
    }

    json_get_optional_str(obj, "paidAt", payment->paid_at,
                          sizeof(payment->paid_at));
    json_get_optional_str(obj, "payerEmail", payment->payer_email,
                          sizeof(payment->payer_email));
    json_get_optional_str(obj, "payerPhone", payment->payer_phone,
                          sizeof(payment->payer_phone));
    json_get_optional_str(obj, "metadata", payment->metadata,
                          sizeof(payment->metadata));
    return ULAB_OK;
}

int bff_record_cash_package_sale(bff_client_t *c,
                                 ue_t *ue,
                                 const package_t *pkg,
                                 const subscriber_t *subscriber,
                                 double amount,
                                 const char *currency,
                                 bff_payment_t *payment,
                                 ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char item_esc[ULAB_MAX_ID * 2];
    char sim_esc[ULAB_MAX_ID * 2];
    char currency_esc[ULAB_MAX_REF * 2];
    char email_esc[ULAB_MAX_NAME * 2];
    char phone_esc[ULAB_MAX_REF * 2];
    json_t *root;
    json_t *obj;
    int n;

    if (ue == NULL || ue->bff_id[0] == '\0' ||
        pkg == NULL || pkg->bff_id[0] == '\0' || subscriber == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "cash package sale requires SIM, package and subscriber");
        return ULAB_ERR;
    }

    if (amount <= 0) {
        amount = pkg->amount;
    }
    ulab_json_escape(pkg->bff_id, item_esc, sizeof(item_esc));
    ulab_json_escape(ue->bff_id, sim_esc, sizeof(sim_esc));
    ulab_json_escape(currency && currency[0] ? currency :
                     (pkg->currency[0] ? pkg->currency : "USD"),
                     currency_esc, sizeof(currency_esc));
    ulab_json_escape(subscriber->email, email_esc, sizeof(email_esc));
    ulab_json_escape(subscriber->phone, phone_esc, sizeof(phone_esc));

    n = snprintf(vars, sizeof(vars),
                 "{\"data\":{\"itemId\":\"%s\","
                 "\"itemType\":\"package\",\"paymentMethod\":\"cash\","
                 "\"amount\":\"%.2f\",\"currency\":\"%s\","
                 "\"sim\":\"%s\",\"payerEmail\":\"%s\","
                 "\"payerPhone\":\"%s\"}}",
                 item_esc, amount, currency_esc, sim_esc,
                 email_esc, phone_esc);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "cash package sale variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "addPayment", BFF_ADD_PAYMENT, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "addPayment");
    if (parse_payment(obj, payment, err)) {
        json_decref(root);
        return ULAB_ERR;
    }

    ulab_copy(ue->last_payment_id, sizeof(ue->last_payment_id),
              payment->id);
    ulab_copy(ue->last_payment_status, sizeof(ue->last_payment_status),
              payment->status);
    json_decref(root);
    return ULAB_OK;
}

int bff_get_package_payments(bff_client_t *c,
                             const package_t *pkg,
                             bff_payment_t payments[],
                             size_t max_payments,
                             size_t *payment_count,
                             ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char item_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *arr;
    size_t i;
    size_t count;

    if (pkg == NULL || pkg->bff_id[0] == '\0' ||
        payments == NULL || payment_count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPayments requires package and output storage");
        return ULAB_ERR;
    }

    ulab_json_escape(pkg->bff_id, item_esc, sizeof(item_esc));
    snprintf(vars, sizeof(vars),
             "{\"data\":{\"type\":\"package\",\"itemId\":\"%s\"}}",
             item_esc);

    root = NULL;
    if (bff_call(c, "getPayments", BFF_GET_PAYMENTS, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getPayments");
    arr = obj ? json_object_get(obj, "payments") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPayments missing payments list");
        json_decref(root);
        return ULAB_ERR;
    }

    count = json_array_size(arr);
    if (count > max_payments) {
        count = max_payments;
    }
    for (i = 0; i < count; i++) {
        if (parse_payment(json_array_get(arr, i), &payments[i], err)) {
            json_decref(root);
            return ULAB_ERR;
        }
    }
    *payment_count = count;
    json_decref(root);
    return ULAB_OK;
}

static double json_optional_number(json_t *obj, const char *key) {
    json_t *v;

    v = obj ? json_object_get(obj, key) : NULL;
    return v != NULL && json_is_number(v) ? json_number_value(v) : 0.0;
}

static int kpi_scope_matches(const bff_kpi_value_t *value,
                             const char *scope_key,
                             const char *scope_value) {
    size_t i;

    if (scope_key == NULL || scope_key[0] == '\0') {
        return 1;
    }
    for (i = 0; i < value->scope_count; i++) {
        if (ulab_streq(value->scope[i].key, scope_key) &&
            (scope_value == NULL || scope_value[0] == '\0' ||
             ulab_streq(value->scope[i].value, scope_value))) {
            return 1;
        }
    }
    return 0;
}

static int parse_kpi_value(json_t *obj, bff_kpi_value_t *value,
                           ulab_error_t *err) {
    json_t *scope;
    json_t *trend;
    json_t *flag;
    size_t i;
    size_t count;

    if (obj == NULL || value == NULL) {
        snprintf(err->msg, sizeof(err->msg), "KPI response is empty");
        return ULAB_ERR;
    }
    memset(value, 0, sizeof(*value));
    if (json_get_str(obj, "kpi", value->kpi, sizeof(value->kpi))) {
        snprintf(err->msg, sizeof(err->msg), "KPI response missing key");
        return ULAB_ERR;
    }
    value->value = json_optional_number(obj, "value");
    json_get_optional_str(obj, "span", value->span, sizeof(value->span));
    json_get_optional_str(obj, "op", value->op, sizeof(value->op));
    json_get_optional_str(obj, "unit", value->unit, sizeof(value->unit));
    json_get_optional_str(obj, "symbol", value->symbol,
                          sizeof(value->symbol));
    json_get_optional_str(obj, "from", value->from, sizeof(value->from));
    json_get_optional_str(obj, "to", value->to, sizeof(value->to));
    json_get_optional_str(obj, "computedAt", value->computed_at,
                          sizeof(value->computed_at));
    flag = json_object_get(obj, "isPartial");
    value->is_partial = flag != NULL && json_is_true(flag);

    scope = json_object_get(obj, "scope");
    if (scope != NULL && json_is_array(scope)) {
        count = json_array_size(scope);
        if (count > ULAB_MAX_BFF_KPI_SCOPES) {
            count = ULAB_MAX_BFF_KPI_SCOPES;
        }
        for (i = 0; i < count; i++) {
            json_t *entry;

            entry = json_array_get(scope, i);
            if (json_get_str(entry, "key", value->scope[i].key,
                             sizeof(value->scope[i].key)) ||
                json_get_str(entry, "value", value->scope[i].value,
                             sizeof(value->scope[i].value))) {
                snprintf(err->msg, sizeof(err->msg),
                         "KPI response has invalid scope");
                return ULAB_ERR;
            }
        }
        value->scope_count = count;
    }

    trend = json_object_get(obj, "trend");
    if (trend != NULL && json_is_object(trend)) {
        json_get_optional_str(trend, "direction", value->trend.direction,
                              sizeof(value->trend.direction));
        value->trend.change_pct = json_optional_number(trend, "changePct");
        value->trend.change_abs = json_optional_number(trend, "changeAbs");
        value->trend.previous_value =
            json_optional_number(trend, "prevValue");
        flag = json_object_get(trend, "hasPrevious");
        value->trend.has_previous = flag != NULL && json_is_true(flag);
    }
    return ULAB_OK;
}

int bff_get_kpi_value(bff_client_t *c,
                      const char *key,
                      const char *span,
                      const char *op,
                      const char *network_id,
                      const char *scope_key,
                      const char *scope_value,
                      bff_kpi_value_t *value,
                      int *found,
                      ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char key_esc[ULAB_MAX_REF * 2];
    char span_esc[ULAB_MAX_REF * 2];
    char op_esc[ULAB_MAX_REF * 2];
    char network_esc[ULAB_MAX_ID * 2];
    char optional[ULAB_MAX_QUERY / 2];
    json_t *root;
    json_t *obj;
    json_t *arr;
    size_t i;
    int n;

    if (key == NULL || key[0] == '\0' || value == NULL || found == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getKpiValues requires key and output storage");
        return ULAB_ERR;
    }
    ulab_json_escape(key, key_esc, sizeof(key_esc));
    ulab_json_escape(span && span[0] ? span : "daily",
                     span_esc, sizeof(span_esc));
    optional[0] = '\0';
    if (op != NULL && op[0] != '\0') {
        ulab_json_escape(op, op_esc, sizeof(op_esc));
        snprintf(optional + strlen(optional),
                 sizeof(optional) - strlen(optional),
                 ",\"op\":\"%s\"", op_esc);
    }
    if (network_id != NULL && network_id[0] != '\0') {
        ulab_json_escape(network_id, network_esc, sizeof(network_esc));
        snprintf(optional + strlen(optional),
                 sizeof(optional) - strlen(optional),
                 ",\"networkId\":\"%s\"", network_esc);
    }
    n = snprintf(vars, sizeof(vars),
                 "{\"data\":{\"keys\":[\"%s\"],\"span\":\"%s\"%s}}",
                 key_esc, span_esc, optional);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg), "KPI variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "getKpiValues", BFF_GET_KPI_VALUES,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getKpiValues");
    arr = obj ? json_object_get(obj, "values") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getKpiValues missing values list");
        json_decref(root);
        return ULAB_ERR;
    }

    *found = 0;
    for (i = 0; i < json_array_size(arr); i++) {
        bff_kpi_value_t candidate;

        if (parse_kpi_value(json_array_get(arr, i), &candidate, err)) {
            json_decref(root);
            return ULAB_ERR;
        }
        if (ulab_streq(candidate.kpi, key) &&
            kpi_scope_matches(&candidate, scope_key, scope_value)) {
            *value = candidate;
            *found = 1;
            break;
        }
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_get_performance_report_cell(bff_client_t *c,
                                    const char *report,
                                    const char *span,
                                    const char *network_id,
                                    const char *entity_id,
                                    const char *column,
                                    double *value,
                                    char *unit,
                                    size_t unit_len,
                                    int *found,
                                    ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char report_esc[ULAB_MAX_REF * 2];
    char span_esc[ULAB_MAX_REF * 2];
    char network_esc[ULAB_MAX_ID * 2];
    char optional[ULAB_MAX_ID * 2 + 32];
    json_t *root;
    json_t *obj;
    json_t *rows;
    size_t i;
    int n;

    if (report == NULL || report[0] == '\0' || entity_id == NULL ||
        entity_id[0] == '\0' || column == NULL || column[0] == '\0' ||
        value == NULL || found == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "performance report requires report, entity and column");
        return ULAB_ERR;
    }
    ulab_json_escape(report, report_esc, sizeof(report_esc));
    ulab_json_escape(span && span[0] ? span : "daily",
                     span_esc, sizeof(span_esc));
    optional[0] = '\0';
    if (network_id != NULL && network_id[0] != '\0') {
        ulab_json_escape(network_id, network_esc, sizeof(network_esc));
        snprintf(optional, sizeof(optional),
                 ",\"networkId\":\"%s\"", network_esc);
    }
    n = snprintf(vars, sizeof(vars),
                 "{\"data\":{\"report\":\"%s\",\"span\":\"%s\"%s}}",
                 report_esc, span_esc, optional);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "performance report variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "getPerformanceReport", BFF_GET_PERFORMANCE_REPORT,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getPerformanceReport");
    rows = obj ? json_object_get(obj, "rows") : NULL;
    if (rows == NULL || !json_is_array(rows)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPerformanceReport missing rows");
        json_decref(root);
        return ULAB_ERR;
    }

    *found = 0;
    for (i = 0; i < json_array_size(rows); i++) {
        json_t *row;
        json_t *cells;
        char row_id[ULAB_MAX_ID];
        size_t j;

        row = json_array_get(rows, i);
        if (json_get_str(row, "entityId", row_id, sizeof(row_id)) ||
            !ulab_streq(row_id, entity_id)) {
            continue;
        }
        cells = json_object_get(row, "cells");
        if (cells == NULL || !json_is_array(cells)) {
            continue;
        }
        for (j = 0; j < json_array_size(cells); j++) {
            json_t *cell;
            char cell_column[ULAB_MAX_REF];

            cell = json_array_get(cells, j);
            if (json_get_str(cell, "column", cell_column,
                             sizeof(cell_column)) ||
                !ulab_streq(cell_column, column)) {
                continue;
            }
            *value = json_optional_number(cell, "value");
            if (unit != NULL && unit_len > 0) {
                json_get_optional_str(cell, "unit", unit, unit_len);
            }
            *found = 1;
            break;
        }
        if (*found) {
            break;
        }
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_get_performance_report_row(bff_client_t *c,
                                   const char *report,
                                   const char *span,
                                   const char *network_id,
                                   const char *entity_id,
                                   bff_performance_row_t *row,
                                   int *found,
                                   ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char report_esc[ULAB_MAX_REF * 2];
    char span_esc[ULAB_MAX_REF * 2];
    char network_esc[ULAB_MAX_ID * 2];
    char optional[ULAB_MAX_ID * 2 + 32];
    json_t *root;
    json_t *obj;
    json_t *rows;
    size_t i;
    int n;

    if (report == NULL || report[0] == '\0' || entity_id == NULL ||
        entity_id[0] == '\0' || row == NULL || found == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "performance report row requires report and entity");
        return ULAB_ERR;
    }
    ulab_json_escape(report, report_esc, sizeof(report_esc));
    ulab_json_escape(span && span[0] ? span : "daily",
                     span_esc, sizeof(span_esc));
    optional[0] = '\0';
    if (network_id != NULL && network_id[0] != '\0') {
        ulab_json_escape(network_id, network_esc, sizeof(network_esc));
        snprintf(optional, sizeof(optional),
                 ",\"networkId\":\"%s\"", network_esc);
    }
    n = snprintf(vars, sizeof(vars),
                 "{\"data\":{\"report\":\"%s\",\"span\":\"%s\"%s}}",
                 report_esc, span_esc, optional);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "performance report variables too long");
        return ULAB_ERR;
    }
    root = NULL;
    if (bff_call(c, "getPerformanceReport", BFF_GET_PERFORMANCE_REPORT,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getPerformanceReport");
    rows = obj ? json_object_get(obj, "rows") : NULL;
    if (rows == NULL || !json_is_array(rows)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPerformanceReport missing rows");
        json_decref(root);
        return ULAB_ERR;
    }
    memset(row, 0, sizeof(*row));
    row->row_count = (uint32_t)json_array_size(rows);
    *found = 0;
    for (i = 0; i < json_array_size(rows); i++) {
        json_t *candidate;
        json_t *attributes;
        char id[ULAB_MAX_ID];
        size_t j;

        candidate = json_array_get(rows, i);
        if (json_get_str(candidate, "entityId", id, sizeof(id)) ||
            !ulab_streq(id, entity_id)) {
            continue;
        }
        row->row_index = (uint32_t)i;
        json_get_optional_str(candidate, "status", row->status,
                              sizeof(row->status));
        attributes = json_object_get(candidate, "attributes");
        if (attributes == NULL || !json_is_array(attributes)) {
            snprintf(err->msg, sizeof(err->msg),
                     "performance report row has no attributes");
            json_decref(root);
            return ULAB_ERR;
        }
        for (j = 0; j < json_array_size(attributes); j++) {
            json_t *attribute;
            char key[ULAB_MAX_REF];

            attribute = json_array_get(attributes, j);
            if (json_get_str(attribute, "key", key, sizeof(key))) {
                continue;
            }
            if (ulab_streq(key, "name")) row->has_name = 1;
            else if (ulab_streq(key, "price")) row->has_price = 1;
            else if (ulab_streq(key, "validity")) row->has_validity = 1;
            else if (ulab_streq(key, "active")) {
                row->has_active = 1;
                json_get_optional_str(attribute, "value", row->active,
                                      sizeof(row->active));
            }
        }
        *found = 1;
        break;
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_add_subscriber(bff_client_t *c, subscriber_t *sub,
                       const network_t *net, ulab_error_t *err) {
    char vars[4096];
    json_t *root;
    json_t *obj;

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"email\":\"%s\","
             "\"name\":\"%s\","
             "\"network_id\":\"%s\","
             "\"phone\":\"%s\"}}",
             sub->email, sub->name, net->bff_id, sub->phone);

    if (bff_call(c, "addSubscriber", BFF_ADD_SUBSCRIBER, vars, &root,
        err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "addSubscriber");
    if (obj == NULL || json_get_str(obj, "uuid", sub->bff_id,
        sizeof(sub->bff_id))) {
        snprintf(err->msg, sizeof(err->msg), "addSubscriber missing uuid");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);

    return ULAB_OK;
}

int bff_allocate_sim(bff_client_t *c, ue_t *ue, const subscriber_t *sub,
                     const network_t *net, const package_t *pkg,
                     const char *sim_type,
                     ulab_error_t *err) {

    char vars[ULAB_MAX_QUERY];
    json_t *root;
    json_t *obj;

    snprintf(vars, sizeof(vars),
         "{\"data\":{\"iccid\":\"%s\",\"network_id\":\"%s\","
         "\"sim_type\":\"%s\",\"package_id\":\"%s\","
         "\"subscriber_id\":\"%s\",\"traffic_policy\":1}}",
         ue->iccid, net->bff_id, sim_type,
         pkg->bff_id, sub->bff_id);

    if (bff_call(c, "allocateSim", BFF_ALLOCATE_SIM, vars, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "allocateSim");
    if (obj == NULL || json_get_str(obj, "id", ue->bff_id,
        sizeof(ue->bff_id))) {
        snprintf(err->msg, sizeof(err->msg), "allocateSim missing id");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);

    return ULAB_OK;
}

int bff_get_sim_usage(bff_client_t *c,
                      const ue_t *ue,
                      const network_t *network,
                      uint64_t *used_mb,
                      ulab_error_t *err) {

    char vars[ULAB_MAX_QUERY];
    json_t *root;
    json_t *arr;
    const char *usage_str;
    size_t i;

    if (ue == NULL || ue->bff_id[0] == '\0' || network == NULL ||
        network->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "getSimsUsageByNetwork missing SIM or network id");
        return ULAB_ERR;
    }

    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
             network->bff_id);

    if (bff_call(c, "getSimsUsageByNetwork",
                 BFF_GET_SIMS_USAGE_BY_NETWORK, vars, &root, err)) {
        return ULAB_ERR;
    }
    arr = dig(root, "data", "getSimsUsageByNetwork");
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSimsUsageByNetwork missing usage list");
        json_decref(root);
        return ULAB_ERR;
    }
    for (i = 0; i < json_array_size(arr); i++) {
        json_t *entry;
        json_t *u;
        char sim_id[ULAB_MAX_ID];

        entry = json_array_get(arr, i);
        if (json_get_str(entry, "simId", sim_id, sizeof(sim_id)) ||
            !ulab_streq(sim_id, ue->bff_id)) {
            continue;
        }
        u = json_object_get(entry, "usage");
        if (u != NULL && json_is_integer(u)) {
            *used_mb = (uint64_t)json_integer_value(u);
            json_decref(root);
            return ULAB_OK;
        }
        if (u != NULL && json_is_string(u)) {
            usage_str = json_string_value(u);
            if (usage_str != NULL &&
                ulab_parse_u64(usage_str, used_mb) == ULAB_OK) {
                json_decref(root);
                return ULAB_OK;
            }
        }
    }

    snprintf(err->msg, sizeof(err->msg),
             "getSimsUsageByNetwork missing selected SIM usage");
    json_decref(root);

    return ULAB_ERR;
}

int bff_get_packages_for_sim(bff_client_t *c, const ue_t *ue,
                             const char *package_id, int *active,
                             ulab_error_t *err) {

    bff_sim_package_t packages[ULAB_MAX_BFF_SIM_PACKAGES];
    size_t count;
    size_t i;

    if (bff_get_sim_packages(c, ue, packages,
                             ULAB_MAX_BFF_SIM_PACKAGES,
                             &count, err)) {
        return ULAB_ERR;
    }

    *active = 0;
    for (i = 0; i < count; i++) {
        if ((package_id == NULL || package_id[0] == '\0' ||
             ulab_streq(packages[i].package_id, package_id)) &&
            packages[i].active) {
            *active = 1;
            if (package_id != NULL && package_id[0] != '\0') {
                break;
            }
        }
    }
    return ULAB_OK;
}

int bff_get_sim_packages(bff_client_t *c,
                         const ue_t *ue,
                         bff_sim_package_t packages[],
                         size_t max_packages,
                         size_t *package_count,
                         ulab_error_t *err) {

    char vars[ULAB_MAX_QUERY];
    json_t *root;
    json_t *obj;
    json_t *arr;
    json_t *it;
    json_t *act;
    size_t i;
    size_t count;

    if (ue == NULL || ue->bff_id[0] == '\0' || packages == NULL ||
        package_count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPackagesForSim requires SIM and output storage");
        return ULAB_ERR;
    }

    snprintf(vars, sizeof(vars), "{\"data\":{\"sim_id\":\"%s\"}}",
             ue->bff_id);

    if (bff_call(c, "getPackagesForSim", BFF_GET_SIM_PACKAGES, vars,
        &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "getPackagesForSim");
    arr = obj ? json_object_get(obj, "packages") : NULL;

    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg), "getPackagesForSim missing list");
        json_decref(root);
        return ULAB_ERR;
    }

    count = json_array_size(arr);
    if (count > max_packages) {
        count = max_packages;
    }
    for (i = 0; i < count; i++) {
        it = json_array_get(arr, i);
        memset(&packages[i], 0, sizeof(packages[i]));
        if (json_get_str(it, "id", packages[i].id,
                         sizeof(packages[i].id)) ||
            json_get_str(it, "package_id", packages[i].package_id,
                         sizeof(packages[i].package_id)) ||
            json_get_str(it, "start_date", packages[i].start_date,
                         sizeof(packages[i].start_date)) ||
            json_get_str(it, "end_date", packages[i].end_date,
                         sizeof(packages[i].end_date))) {
            snprintf(err->msg, sizeof(err->msg),
                     "getPackagesForSim returned incomplete assignment");
            json_decref(root);
            return ULAB_ERR;
        }
        act = it ? json_object_get(it, "is_active") : NULL;
        packages[i].active = act != NULL && json_is_true(act);
    }

    *package_count = count;
    json_decref(root);

    return ULAB_OK;
}

int bff_get_node_status(bff_client_t *c, const node_t *node,
                        bff_node_status_t *status, ulab_error_t *err) {

    char vars[ULAB_MAX_QUERY];
    json_t *root;
    json_t *obj;
    json_t *node_status;

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"id\":\"%s\"}}",
             node->bff_id);

    if (bff_call(c, "getNode", BFF_GET_NODE, vars, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "getNode");
    if (obj == NULL) {
        snprintf(err->msg, sizeof(err->msg), "getNode missing data");
        json_decref(root);
        return ULAB_ERR;
    }

    node_status = json_object_get(obj, "status");
    if (node_status == NULL || !json_is_object(node_status)) {
        snprintf(err->msg, sizeof(err->msg), "getNode missing status");
        json_decref(root);
        return ULAB_ERR;
    }

    if (json_get_str(obj, "id", status->id, sizeof(status->id)) ||
        json_get_str(node_status, "connectivity", status->connectivity,
                     sizeof(status->connectivity)) ||
        json_get_str(node_status, "state", status->state,
                     sizeof(status->state))) {
        snprintf(err->msg, sizeof(err->msg),
                 "getNode missing id/status fields");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);

    return ULAB_OK;
}

int bff_get_release(bff_client_t *c,
                    const char *name,
                    const char *type,
                    const char *version,
                    bff_release_t *release,
                    int *found,
                    ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char name_esc[ULAB_MAX_NAME * 2];
    char type_esc[ULAB_MAX_REF * 2];
    json_t *root;
    json_t *obj;
    json_t *arr;
    json_t *it;
    json_t *value;
    const char *release_version;
    size_t i;
    int n;

    if (name == NULL || name[0] == '\0' ||
        type == NULL || type[0] == '\0' ||
        version == NULL || version[0] == '\0' ||
        release == NULL || found == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getReleaseCatalog requires name, type and version");
        return ULAB_ERR;
    }

    memset(release, 0, sizeof(*release));
    *found = 0;
    ulab_json_escape(name, name_esc, sizeof(name_esc));
    ulab_json_escape(type, type_esc, sizeof(type_esc));

    n = snprintf(vars, sizeof(vars),
                 "{\"name\":\"%s\",\"type\":\"%s\"}",
                 name_esc, type_esc);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getReleaseCatalog variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "getReleaseCatalog", BFF_GET_RELEASE_CATALOG,
                 vars, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "getReleaseCatalog");
    arr = obj ? json_object_get(obj, "releases") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getReleaseCatalog missing releases list");
        json_decref(root);
        return ULAB_ERR;
    }

    for (i = 0; i < json_array_size(arr); i++) {
        it = json_array_get(arr, i);
        value = it ? json_object_get(it, "version") : NULL;
        release_version = value && json_is_string(value) ?
            json_string_value(value) : NULL;
        if (release_version == NULL ||
            !ulab_streq(release_version, version)) {
            continue;
        }

        if (json_get_str(it, "name", release->name,
                         sizeof(release->name)) ||
            json_get_str(it, "type", release->type,
                         sizeof(release->type)) ||
            json_get_str(it, "version", release->version,
                         sizeof(release->version))) {
            snprintf(err->msg, sizeof(err->msg),
                     "getReleaseCatalog returned incomplete release");
            json_decref(root);
            return ULAB_ERR;
        }

        value = json_object_get(it, "available");
        release->available = value != NULL && json_is_true(value);
        value = json_object_get(it, "chunked");
        release->chunked = value != NULL && json_is_true(value);
        value = json_object_get(it, "desired");
        release->desired = value != NULL && json_is_true(value);
        value = json_object_get(it, "uploadedAt");
        if (value != NULL && json_is_string(value)) {
            ulab_copy(release->uploaded_at,
                      sizeof(release->uploaded_at),
                      json_string_value(value));
        }

        *found = 1;
        json_decref(root);
        return ULAB_OK;
    }

    json_decref(root);
    return ULAB_OK;
}

int bff_promote_release(bff_client_t *c,
                        const char *name,
                        const char *type,
                        const char *version,
                        ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char name_esc[ULAB_MAX_NAME * 2];
    char type_esc[ULAB_MAX_REF * 2];
    char version_esc[ULAB_MAX_REF * 2];
    char desired_version[ULAB_MAX_REF];
    char returned_name[ULAB_MAX_NAME];
    json_t *root;
    json_t *obj;
    int n;

    if (name == NULL || name[0] == '\0' ||
        type == NULL || type[0] == '\0' ||
        version == NULL || version[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "promoteRelease requires name, type and version");
        return ULAB_ERR;
    }

    ulab_json_escape(name, name_esc, sizeof(name_esc));
    ulab_json_escape(type, type_esc, sizeof(type_esc));
    ulab_json_escape(version, version_esc, sizeof(version_esc));

    n = snprintf(vars, sizeof(vars),
                 "{\"name\":\"%s\",\"type\":\"%s\","
                 "\"version\":\"%s\"}",
                 name_esc, type_esc, version_esc);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "promoteRelease variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "promoteRelease", BFF_PROMOTE_RELEASE,
                 vars, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "promoteRelease");
    if (obj == NULL ||
        json_get_str(obj, "name", returned_name,
                     sizeof(returned_name)) ||
        json_get_str(obj, "desiredVersion", desired_version,
                     sizeof(desired_version))) {
        snprintf(err->msg, sizeof(err->msg),
                 "promoteRelease missing name or desiredVersion");
        json_decref(root);
        return ULAB_ERR;
    }

    if (!ulab_streq(returned_name, name) ||
        !ulab_streq(desired_version, version)) {
        snprintf(err->msg, sizeof(err->msg),
                 "promoteRelease returned name=%s desired=%s",
                 returned_name, desired_version);
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);
    return ULAB_OK;
}

int bff_update_software(bff_client_t *c,
                        const node_t *node,
                        const char *app,
                        const char *tag,
                        ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char app_esc[ULAB_MAX_NAME * 2];
    char tag_esc[ULAB_MAX_REF * 2];
    char node_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *message;
    int n;

    if (node == NULL || node->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "updateSoftware missing node id");
        return ULAB_ERR;
    }

    if (app == NULL || app[0] == '\0' || tag == NULL || tag[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "updateSoftware requires app and tag");
        return ULAB_ERR;
    }

    ulab_json_escape(app, app_esc, sizeof(app_esc));
    ulab_json_escape(tag, tag_esc, sizeof(tag_esc));
    ulab_json_escape(node->bff_id, node_esc, sizeof(node_esc));

    n = snprintf(vars, sizeof(vars),
                 "{\"data\":{\"name\":\"%s\","
                 "\"nodeId\":\"%s\",\"tag\":\"%s\"}}",
                 app_esc, node_esc, tag_esc);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "updateSoftware variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "updateSoftware", BFF_UPDATE_SOFTWARE, vars,
                 &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "updateSoftware");
    message = obj ? json_object_get(obj, "message") : NULL;
    if (message == NULL || !json_is_string(message)) {
        snprintf(err->msg, sizeof(err->msg),
                 "updateSoftware missing message");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);
    return ULAB_OK;
}

int bff_get_node_app(bff_client_t *c,
                     const node_t *node,
                     const char *app,
                     bff_app_state_t *state,
                     ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char app_esc[ULAB_MAX_NAME * 2];
    char node_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *arr;
    json_t *it;
    json_t *name;
    const char *value;
    size_t i;
    int n;

    if (node == NULL || node->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "getApps missing node id");
        return ULAB_ERR;
    }

    if (app == NULL || app[0] == '\0' || state == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getApps requires app and output state");
        return ULAB_ERR;
    }

    memset(state, 0, sizeof(*state));
    ulab_json_escape(app, app_esc, sizeof(app_esc));
    ulab_json_escape(node->bff_id, node_esc, sizeof(node_esc));

    n = snprintf(vars, sizeof(vars),
                 "{\"data\":{\"nodeId\":\"%s\","
                 "\"appName\":\"%s\"}}",
                 node_esc, app_esc);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg), "getApps variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "getApps", BFF_GET_APPS, vars, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "getApps");
    arr = obj ? json_object_get(obj, "apps") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg), "getApps missing apps list");
        json_decref(root);
        return ULAB_ERR;
    }

    for (i = 0; i < json_array_size(arr); i++) {
        it = json_array_get(arr, i);
        name = it ? json_object_get(it, "name") : NULL;
        value = name && json_is_string(name) ? json_string_value(name) : NULL;
        if (value == NULL || !ulab_streq(value, app)) {
            continue;
        }

        if (json_get_str(it, "name", state->name, sizeof(state->name)) ||
            json_get_str(it, "version", state->version,
                         sizeof(state->version)) ||
            json_get_str(it, "tag", state->tag, sizeof(state->tag)) ||
            json_get_str(it, "status", state->status,
                         sizeof(state->status))) {
            snprintf(err->msg, sizeof(err->msg),
                     "getApps returned incomplete app state for %s", app);
            json_decref(root);
            return ULAB_ERR;
        }

        json_decref(root);
        return ULAB_OK;
    }

    snprintf(err->msg, sizeof(err->msg),
             "getApps app not found: node=%s app=%s", node->bff_id, app);
    json_decref(root);
    return ULAB_ERR;
}

int bff_network_overview_loads(bff_client_t *c, const network_t *net,
                               ulab_error_t *err) {

    char vars[ULAB_MAX_QUERY];
    json_t *root;

    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}", net->bff_id);

    if (bff_call(c, "networkOverview", BFF_NETWORK_OVERVIEW, vars,
        &root, err)) {
        return ULAB_ERR;
    }

    json_decref(root);

    return ULAB_OK;
}

int bff_site_view_loads(bff_client_t *c, const site_t *site,
                        ulab_error_t *err) {

    char vars[ULAB_MAX_QUERY];
    json_t *root;

    snprintf(vars, sizeof(vars), "{\"siteId\":\"%s\"}", site->bff_id);

    if (bff_call(c, "siteView", BFF_SITE_VIEW, vars, &root, err)) {
        return ULAB_ERR;
    }

    json_decref(root);

    return ULAB_OK;
}


static int json_array_has_id(json_t *arr, const char *id) {
    size_t i;

    if (arr == NULL || !json_is_array(arr) || id == NULL || id[0] == '\0') {
        return 0;
    }

    for (i = 0; i < json_array_size(arr); i++) {
        json_t *it = json_array_get(arr, i);
        json_t *v;
        const char *s;

        if (it == NULL || !json_is_object(it)) {
            continue;
        }

        v = json_object_get(it, "id");
        if (v == NULL || !json_is_string(v)) {
            v = json_object_get(it, "uuid");
        }
        if (v == NULL || !json_is_string(v)) {
            continue;
        }

        s = json_string_value(v);
        if (s != NULL && ulab_streq(s, id)) {
            return 1;
        }
    }

    return 0;
}

static int backend_get_networks(bff_client_t *c, json_t **root,
                                ulab_error_t *err) {
    return bff_call(c, "getNetworks", BFF_GET_NETWORKS, "{}", root, err);
}

static int backend_get_sites(bff_client_t *c, const char *network_id,
                             json_t **root, ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char legacy_vars[ULAB_MAX_QUERY];
    ulab_error_t first_err;

    if (root != NULL) {
        *root = NULL;
    }

    snprintf(vars, sizeof(vars), "{\"data\":{\"networkId\":\"%s\"}}",
             network_id ? network_id : "");

    memset(&first_err, 0, sizeof(first_err));
    if (bff_call(c, "getSites", BFF_GET_SITES, vars, root, &first_err) == ULAB_OK) {
        return ULAB_OK;
    }

    /*
     * Keep compatibility with older BFF deployments that still exposed
     * getSites(networkId: String!). Current BFF expects
     * getSites(data: SitesInputDto!).
     */
    if (strstr(first_err.msg, "Unknown argument \"data\"") != NULL ||
        strstr(first_err.msg, "SitesInputDto") != NULL ||
        strstr(first_err.msg, "Field \"getSites\" argument \"networkId\"") != NULL) {
        if (root != NULL && *root != NULL) {
            json_decref(*root);
            *root = NULL;
        }

        snprintf(legacy_vars, sizeof(legacy_vars), "{\"networkId\":\"%s\"}",
                 network_id ? network_id : "");
        return bff_call(c, "getSites", BFF_GET_SITES_LEGACY, legacy_vars,
                        root, err);
    }

    if (err != NULL) {
        *err = first_err;
    }

    return ULAB_ERR;
}

static int backend_get_nodes_for_site(bff_client_t *c, const char *site_id,
                                      json_t **root, ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];

    snprintf(vars, sizeof(vars), "{\"siteId\":\"%s\"}",
             site_id ? site_id : "");
    return bff_call(c, "getNodesForSite", BFF_GET_NODES_FOR_SITE, vars,
                    root, err);
}

static int backend_network_contains(bff_client_t *c, const char *id,
                                    int *found, ulab_error_t *err) {
    json_t *root;
    json_t *obj;
    json_t *arr;

    root = NULL;
    if (backend_get_networks(c, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "getNetworks");
    arr = obj ? json_object_get(obj, "networks") : NULL;
    *found = json_array_has_id(arr, id);
    json_decref(root);

    return ULAB_OK;
}

static int backend_site_contains(bff_client_t *c, const world_t *w,
                                 const char *id, int *found,
                                 ulab_error_t *err) {
    size_t i;

    *found = 0;
    for (i = 0; i < w->network_count; i++) {
        json_t *root;
        json_t *obj;
        json_t *arr;

        if (w->networks[i].bff_id[0] == '\0') {
            continue;
        }

        root = NULL;
        if (backend_get_sites(c, w->networks[i].bff_id, &root, err)) {
            return ULAB_ERR;
        }

        obj = dig(root, "data", "getSites");
        arr = obj ? json_object_get(obj, "sites") : NULL;
        if (json_array_has_id(arr, id)) {
            *found = 1;
            json_decref(root);
            return ULAB_OK;
        }
        json_decref(root);
    }

    return ULAB_OK;
}

static int backend_node_contains(bff_client_t *c, const world_t *w,
                                 const char *id, int *found,
                                 ulab_error_t *err) {
    size_t i;

    *found = 0;
    for (i = 0; i < w->site_count; i++) {
        json_t *root;
        json_t *obj;
        json_t *arr;

        if (w->sites[i].bff_id[0] == '\0') {
            continue;
        }

        root = NULL;
        if (backend_get_nodes_for_site(c, w->sites[i].bff_id, &root, err)) {
            return ULAB_ERR;
        }

        obj = dig(root, "data", "getNodesForSite");
        arr = obj ? json_object_get(obj, "nodes") : NULL;
        if (json_array_has_id(arr, id)) {
            *found = 1;
            json_decref(root);
            return ULAB_OK;
        }
        json_decref(root);
    }

    return ULAB_OK;
}

static int backend_sim_contains(bff_client_t *c, const world_t *w,
                                const char *id, int *found,
                                ulab_error_t *err) {
    size_t i;

    *found = 0;

    for (i = 0; i < w->ue_count; i++) {
        int active;

        if (!ulab_streq(w->ues[i].bff_id, id)) {
            continue;
        }

        active = 0;
        if (bff_get_packages_for_sim(c, &w->ues[i], NULL, &active, err)) {
            return ULAB_ERR;
        }
        *found = 1;
        return ULAB_OK;
    }

    return ULAB_OK;
}

int bff_backend_contains(bff_client_t *c, const char *view, const char *id,
                         const world_t *w, int *found, ulab_error_t *err) {
    if (view == NULL || id == NULL || w == NULL || found == NULL) {
        snprintf(err->msg, sizeof(err->msg), "backend_contains bad args");
        return ULAB_ERR;
    }

    if (ulab_streq(view, "networks")) {
        return backend_network_contains(c, id, found, err);
    }
    if (ulab_streq(view, "sites")) {
        return backend_site_contains(c, w, id, found, err);
    }
    if (ulab_streq(view, "nodes")) {
        return backend_node_contains(c, w, id, found, err);
    }
    if (ulab_streq(view, "sims") || ulab_streq(view, "ues")) {
        return backend_sim_contains(c, w, id, found, err);
    }

    snprintf(err->msg, sizeof(err->msg),
             "unsupported backend list view: %s", view);
    return ULAB_ERR;
}

int bff_backend_count(bff_client_t *c, const char *target, const world_t *w,
                      size_t *count, ulab_error_t *err) {
    size_t i;
    size_t n;

    if (target == NULL || w == NULL || count == NULL) {
        snprintf(err->msg, sizeof(err->msg), "backend_count bad args");
        return ULAB_ERR;
    }

    n = 0;

    if (ulab_streq(target, "networks")) {
        for (i = 0; i < w->network_count; i++) {
            int found = 0;
            if (bff_backend_contains(c, "networks", w->networks[i].bff_id,
                                     w, &found, err)) {
                return ULAB_ERR;
            }
            if (found) n++;
        }
    } else if (ulab_streq(target, "sites")) {
        for (i = 0; i < w->site_count; i++) {
            int found = 0;
            if (bff_backend_contains(c, "sites", w->sites[i].bff_id,
                                     w, &found, err)) {
                return ULAB_ERR;
            }
            if (found) n++;
        }
    } else if (ulab_streq(target, "nodes")) {
        for (i = 0; i < w->node_count; i++) {
            int found = 0;
            if (bff_backend_contains(c, "nodes", w->nodes[i].bff_id,
                                     w, &found, err)) {
                return ULAB_ERR;
            }
            if (found) n++;
        }
    } else if (ulab_streq(target, "sims") || ulab_streq(target, "ues")) {
        for (i = 0; i < w->ue_count; i++) {
            int found = 0;
            if (bff_backend_contains(c, "sims", w->ues[i].bff_id,
                                     w, &found, err)) {
                return ULAB_ERR;
            }
            if (found) n++;
        }
    } else {
        snprintf(err->msg, sizeof(err->msg),
                 "unsupported backend_count target: %s", target);
        return ULAB_ERR;
    }

    *count = n;
    return ULAB_OK;
}

static int bff_cleanup_call(bff_client_t *c,
                            const char *op,
                            const char *query) {
    json_t *root;
    ulab_error_t err;

    root = NULL;
    memset(&err, 0, sizeof(err));

    if (query == NULL || query[0] == '\0') {
        return ULAB_OK;
    }

    if (bff_call(c, op, query, "{}", &root, &err)) {
        if ((ulab_streq(op, "setInactivePackageForSim") ||
             ulab_streq(op, "removePackageForSim")) &&
            strstr(err.msg, "package record not found") != NULL) {
            if (c != NULL && c->logf != NULL) {
                fprintf(c->logf,
                        "cleanup ignore: %s: package record not found\n",
                        op);
                fflush(c->logf);
            }
            return ULAB_OK;
        }

        if (ulab_streq(op, "toggleSimStatus") &&
            strstr(err.msg, "inactive is invalid for deactivation") != NULL) {
            if (c != NULL && c->logf != NULL) {
                fprintf(c->logf,
                        "cleanup ignore: %s: sim already inactive\n",
                        op);
                fflush(c->logf);
            }
            return ULAB_OK;
        }

        if (c != NULL && c->logf != NULL) {
            fprintf(c->logf, "cleanup warning: %s: %s\n", op, err.msg);
            fflush(c->logf);
        }
        return ULAB_ERR;
    }

    json_decref(root);
    return ULAB_OK;
}

static const char *cleanup_package_for_ue(const world_t *w,
                                          const ue_t *ue) {
    size_t i;

    if (ue->sim_package_id[0] != '\0') {
        return ue->sim_package_id;
    }

    for (i = 0; i < w->package_count; i++) {
        if (ulab_streq(w->packages[i].ref, ue->package_ref) &&
            w->packages[i].bff_id[0] != '\0') {
            return w->packages[i].bff_id;
        }
    }

    return NULL;
}


static int cleanup_sim_packages_from_bff(bff_client_t *c,
                                         const ue_t *ue,
                                         int *failures) {
    char vars[ULAB_MAX_QUERY];
    char query[ULAB_MAX_QUERY];
    char package_ids[32][ULAB_MAX_ID];
    json_t *root;
    json_t *obj;
    json_t *arr;
    json_t *it;
    json_t *pid;
    ulab_error_t qerr;
    size_t i;
    size_t count;
    int n;

    root = NULL;
    count = 0;
    memset(&qerr, 0, sizeof(qerr));

    if (ue == NULL || ue->bff_id[0] == '\0') {
        return ULAB_OK;
    }

    snprintf(vars, sizeof(vars), "{\"data\":{\"sim_id\":\"%s\"}}",
             ue->bff_id);

    if (bff_call(c, "getPackagesForSim", BFF_GET_SIM_PACKAGES, vars,
        &root, &qerr)) {
        if (c != NULL && c->logf != NULL) {
            fprintf(c->logf,
                    "cleanup warning: getPackagesForSim failed for sim %s: %s\n",
                    ue->bff_id, qerr.msg);
            fflush(c->logf);
        }
        return ULAB_ERR;
    }

    obj = dig(root, "data", "getPackagesForSim");
    arr = obj ? json_object_get(obj, "packages") : NULL;
    if (arr != NULL && json_is_array(arr)) {
        for (i = 0; i < json_array_size(arr) && count < 32; i++) {
            it = json_array_get(arr, i);
            pid = it ? json_object_get(it, "package_id") : NULL;
            if (pid != NULL && json_is_string(pid) &&
                json_string_value(pid) != NULL &&
                json_string_value(pid)[0] != '\0') {
                ulab_copy(package_ids[count], sizeof(package_ids[count]),
                          json_string_value(pid));
                count++;
            }
        }
    }

    json_decref(root);

    for (i = 0; i < count; i++) {

        n = snprintf(query, sizeof(query),
                     "mutation { removePackageForSim(data: {"
                     "packageId: \"%s\", simId: \"%s\"}) { packageId } }",
                     package_ids[i], ue->bff_id);
        if (n >= 0 && (size_t)n < sizeof(query) &&
            bff_cleanup_call(c, "removePackageForSim", query)) {
            (*failures)++;
        }
    }

    return ULAB_OK;
}

static int cleanup_should_delete_pool_sim(void) {
    const char *v;

    v = getenv("ULAB_DELETE_SIM_POOL_ON_CLEANUP");
    return v != NULL && v[0] != '\0' && !ulab_streq(v, "0");
}

int bff_cleanup_world(bff_client_t *c,
                      const world_t *w,
                      ulab_error_t *err) {
    char query[ULAB_MAX_QUERY];
    const char *pkg_id;
    size_t i;
    int failures;
    int n;

    failures = 0;

    if (c == NULL || c->url[0] == '\0' || w == NULL) {
        return ULAB_OK;
    }

    for (i = 0; i < w->ue_count; i++) {
        const ue_t *ue = &w->ues[i];

        if (ue->bff_id[0] == '\0') {
            continue;
        }

        if (cleanup_sim_packages_from_bff(c, ue, &failures)) {
            pkg_id = cleanup_package_for_ue(w, ue);
            if (pkg_id != NULL && pkg_id[0] != '\0') {
                n = snprintf(query, sizeof(query),
                             "mutation { setInactivePackageForSim(data: {"
                             "packageId: \"%s\", simId: \"%s\"}) { packageId } }",
                             pkg_id, ue->bff_id);
                if (n >= 0 && (size_t)n < sizeof(query) &&
                    bff_cleanup_call(c, "setInactivePackageForSim", query)) {
                    failures++;
                }

                n = snprintf(query, sizeof(query),
                             "mutation { removePackageForSim(data: {"
                             "packageId: \"%s\", simId: \"%s\"}) { packageId } }",
                             pkg_id, ue->bff_id);
                if (n >= 0 && (size_t)n < sizeof(query) &&
                    bff_cleanup_call(c, "removePackageForSim", query)) {
                    failures++;
                }
            }
        }

        n = snprintf(query, sizeof(query),
                     "mutation { toggleSimStatus(data: {sim_id: \"%s\", "
                     "status: \"inactive\"}) { success } }",
                     ue->bff_id);
        if (n >= 0 && (size_t)n < sizeof(query) &&
            bff_cleanup_call(c, "toggleSimStatus", query)) {
            failures++;
        }

        n = snprintf(query, sizeof(query),
                     "mutation { deleteSim(data: {simId: \"%s\"}) { simId } }",
                     ue->bff_id);
        if (n >= 0 && (size_t)n < sizeof(query) &&
            bff_cleanup_call(c, "deleteSim", query)) {
            failures++;
        }

        if (ue->pool_sim_id[0] != '\0' && cleanup_should_delete_pool_sim()) {
            n = snprintf(query, sizeof(query),
                         "mutation { deleteSimFromPool(simId: \"%s\") { success } }",
                         ue->pool_sim_id);
            if (n >= 0 && (size_t)n < sizeof(query) &&
                bff_cleanup_call(c, "deleteSimFromPool", query)) {
                failures++;
            }
        } else if (ue->pool_sim_id[0] != '\0' &&
                   c != NULL && c->logf != NULL) {
            fprintf(c->logf,
                    "cleanup keep SIM pool record id=%s iccid=%s\n",
                    ue->pool_sim_id, ue->iccid);
            fflush(c->logf);
        }

    }

    for (i = 0; i < w->package_count; i++) {
        if (w->packages[i].bff_id[0] == '\0') {
            continue;
        }

        n = snprintf(query, sizeof(query),
                     "mutation { deletePackage(packageId: \"%s\") { uuid } }",
                     w->packages[i].bff_id);
        if (n >= 0 && (size_t)n < sizeof(query) &&
            bff_cleanup_call(c, "deletePackage", query)) {
            failures++;
        }
    }

    for (i = 0; i < w->node_count; i++) {
        if (w->nodes[i].bff_id[0] == '\0') {
            continue;
        }

        n = snprintf(
            query,
            sizeof(query),
            "mutation { releaseNodeFromSite(data: {id: \"%s\"}) "
            "{ success } }",
            w->nodes[i].bff_id);

        if (n >= 0 && (size_t)n < sizeof(query) &&
            bff_cleanup_call(c, "releaseNodeFromSite", query)) {
            failures++;
        }
    }

    for (i = 0; i < w->node_count; i++) {
        if (w->nodes[i].bff_id[0] == '\0') {
            continue;
        }

        n = snprintf(
            query,
            sizeof(query),
            "mutation { deleteNode(data: {id: \"%s\"}) { id } }",
            w->nodes[i].bff_id);

        if (n >= 0 && (size_t)n < sizeof(query) &&
            bff_cleanup_call(c, "deleteNode", query)) {
            failures++;
        }
    }

    for (i = 0; i < w->site_count; i++) {
        if (w->sites[i].bff_id[0] == '\0') {
            continue;
        }

        n = snprintf(
            query,
            sizeof(query),
            "mutation { deleteSite(id: \"%s\") { success } }",
            w->sites[i].bff_id);

        if (n >= 0 && (size_t)n < sizeof(query) &&
            bff_cleanup_call(c, "deleteSite", query)) {
            failures++;
        }
    }

    for (i = 0; i < w->network_count; i++) {
        if (w->networks[i].bff_id[0] == '\0') {
            continue;
        }

        n = snprintf(
            query,
            sizeof(query),
            "mutation { deleteNetwork(id: \"%s\") { success } }",
            w->networks[i].bff_id);

        if (n >= 0 && (size_t)n < sizeof(query) &&
            bff_cleanup_call(c, "deleteNetwork", query)) {
            failures++;
        }
    }

    if (failures > 0 && err != NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "backend cleanup had %d failed step(s)", failures);
        return ULAB_ERR;
    }

    return ULAB_OK;
}
