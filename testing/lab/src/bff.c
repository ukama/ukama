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
extern const char *BFF_ADD_SUBSCRIBER;
extern const char *BFF_ALLOCATE_SIM;
extern const char *BFF_GET_DATA_USAGE;
extern const char *BFF_GET_SIM_PACKAGES;
extern const char *BFF_GET_NODE_STATE;
extern const char *BFF_NETWORK_OVERVIEW;
extern const char *BFF_SITE_VIEW;
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
                    ulab_status("BFF", "site %s anchor online %s lat=%s lng=%s",
                                site->ref, site->tnode_id, site->latitude,
                                site->longitude);
                    return ULAB_OK;
                }
            }
        }

        json_decref(root);

        if (found) {
            ulab_status("BFF", "waiting site %s anchor online/location",
                        site->ref);
        } else {
            ulab_status("BFF", "waiting site %s anchor in registry",
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

int bff_add_package(bff_client_t *c, package_t *p, ulab_error_t *err) {
    char vars[4096];
    json_t *root;
    json_t *obj;
    uint32_t duration;

    duration = p->duration_days;
    if (duration == 0) {
        duration = 1;
    }

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"name\":\"%s\",\"amount\":%.2f,"
             "\"dataUnit\":\"MB\",\"dataVolume\":%llu,"
             "\"duration\":%u,\"currency\":\"USD\","
             "\"country\":\"USA\"}}", p->name, p->amount,
             (unsigned long long)p->data_mb, duration);

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
                      uint64_t *used_mb,
                      ulab_error_t *err) {

    char vars[ULAB_MAX_QUERY];
    json_t *root;
    json_t *obj;
    json_t *u;

    snprintf(vars, sizeof(vars), "{\"simId\":\"%s\"}", ue->bff_id);

    if (bff_call(c, "getDataUsage", BFF_GET_DATA_USAGE, vars, &root,
        err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "getDataUsage");
    u = obj ? json_object_get(obj, "usage") : NULL;

    if (u == NULL || !json_is_integer(u)) {
        snprintf(err->msg, sizeof(err->msg), "getDataUsage missing usage");
        json_decref(root);
        return ULAB_ERR;
    }

    *used_mb = (uint64_t)json_integer_value(u);

    json_decref(root);

    return ULAB_OK;
}

int bff_get_packages_for_sim(bff_client_t *c, const ue_t *ue,
                             const char *package_id, int *active,
                             ulab_error_t *err) {

    char vars[ULAB_MAX_QUERY];
    json_t *root;
    json_t *obj;
    json_t *arr;
    json_t *it;
    json_t *pid;
    json_t *act;
    size_t i;

    snprintf(vars, sizeof(vars), "{\"data\":{\"simId\":\"%s\"}}",
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

    *active = 0;

    for (i = 0; i < json_array_size(arr); i++) {
        it = json_array_get(arr, i);
        pid = it ? json_object_get(it, "package_id") : NULL;
        act = it ? json_object_get(it, "is_active") : NULL;

        if (pid != NULL && json_is_string(pid) &&
            ulab_streq(json_string_value(pid), package_id)) {
            if (act != NULL) {
                *active = json_is_true(act);
            }
        }
    }

    json_decref(root);

    return ULAB_OK;
}

int bff_get_node_state(bff_client_t *c, const node_t *node,
                       bff_node_state_t *state, ulab_error_t *err) {

    char vars[ULAB_MAX_QUERY];
    json_t *root;
    json_t *obj;

    snprintf(vars, sizeof(vars), "{\"nodeId\":\"%s\"}", node->bff_id);

    if (bff_call(c, "getNodeState", BFF_GET_NODE_STATE, vars, &root,
        err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "getNodeState");
    if (obj == NULL) {
        snprintf(err->msg, sizeof(err->msg), "getNodeState missing data");
        json_decref(root);
        return ULAB_ERR;
    }

    json_get_str(obj, "state", state->state, sizeof(state->state));
    json_get_str(obj, "connectivity", state->connectivity,
                 sizeof(state->connectivity));

    json_decref(root);

    return ULAB_OK;
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

int bff_query_count(bff_client_t *c, const char *target, const world_t *w,
                    size_t *count, ulab_error_t *err) {

    (void)c;
    (void)err;

    if (ulab_streq(target, "networks")) {
        *count = w->network_count;
    } else if (ulab_streq(target, "sites")) {
        *count = w->site_count;
    } else if (ulab_streq(target, "nodes")) {
        *count = w->node_count;
    } else if (ulab_streq(target, "packages")) {
        *count = w->package_count;
    } else if (ulab_streq(target, "subscribers")) {
        *count = w->subscriber_count;
    } else if (ulab_streq(target, "sims") ||
               ulab_streq(target, "ues")) {
        *count = w->ue_count;
    } else {
        return ULAB_ERR;
    }

    return ULAB_OK;
}
