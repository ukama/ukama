/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <jansson.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <time.h>
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
extern const char *BFF_GET_NETWORK;
extern const char *BFF_GET_SITE;
extern const char *BFF_GET_SITES;
extern const char *BFF_CONSOLE_SITES_LIST;
extern const char *BFF_CONSOLE_SITE_NODE_COUNTS;
extern const char *BFF_CONSOLE_SITE_DETAIL;
extern const char *BFF_GET_NODES;
extern const char *BFF_CONSOLE_NODES_LIST;
extern const char *BFF_CONSOLE_NODE_POOL;
extern const char *BFF_CONSOLE_NODE_DETAIL;
extern const char *BFF_GET_SUBSCRIBERS_BY_NETWORK;
extern const char *BFF_GET_SUBSCRIBER;
extern const char *BFF_GET_SIMS_BY_NETWORK;
extern const char *BFF_SIM_POOL_OVERVIEW;
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
extern const char *BFF_GET_SOFTWARES;
extern const char *BFF_GET_NODE_OPERATION_STATUS;
extern const char *BFF_GET_SITE_OPERATION_STATUS;
extern const char *BFF_GET_KPI_TIMESERIES;
extern const char *BFF_GET_NETWORKS;
extern const char *BFF_GET_NODES_FOR_SITE;
extern const char *BFF_GET_COMPONENTS_BY_USER_ID;

typedef struct {
    char *buf;
    size_t len;
} http_buf_t;

static int payment_status_settled(const char *status);

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

static int console_section_ok(json_t *section,
                              const char *operation,
                              const char *section_name,
                              ulab_error_t *err) {
    json_t *error;
    json_t *code_value;
    json_t *message_value;
    const char *code;
    const char *message;

    if (section == NULL || !json_is_object(section)) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s missing %s section", operation, section_name);
        return ULAB_ERR;
    }

    error = json_object_get(section, "error");
    if (error == NULL || json_is_null(error)) {
        return ULAB_OK;
    }
    if (!json_is_object(error)) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s %s section returned invalid error",
                 operation, section_name);
        return ULAB_ERR;
    }

    code_value = json_object_get(error, "code");
    message_value = json_object_get(error, "message");
    code = code_value && json_is_string(code_value) ?
        json_string_value(code_value) : "UNKNOWN";
    message = message_value && json_is_string(message_value) ?
        json_string_value(message_value) : "section failed";
    snprintf(err->msg, sizeof(err->msg),
             "%s %s section failed: %.48s: %.160s",
             operation, section_name, code, message);
    return ULAB_ERR;
}

static json_t *console_find_id(json_t *arr, const char *id) {
    size_t i;

    if (arr == NULL || !json_is_array(arr) || id == NULL) {
        return NULL;
    }
    for (i = 0; i < json_array_size(arr); i++) {
        json_t *item;
        json_t *value;
        const char *actual;

        item = json_array_get(arr, i);
        value = item ? json_object_get(item, "id") : NULL;
        actual = value && json_is_string(value) ?
            json_string_value(value) : NULL;
        if (actual != NULL && ulab_streq(actual, id)) {
            return item;
        }
    }
    return NULL;
}

static int console_nodes_array(bff_client_t *c,
                               const network_t *network,
                               const char *view,
                               json_t **root,
                               json_t **arr,
                               ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    const char *operation;
    const char *query;
    json_t *nodes_view;
    json_t *section;

    if (root == NULL || arr == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "console node query has invalid output storage");
        return ULAB_ERR;
    }
    *root = NULL;
    *arr = NULL;

    if (ulab_streq(view, "node_pool")) {
        operation = "NodePool";
        query = BFF_CONSOLE_NODE_POOL;
        snprintf(vars, sizeof(vars), "{}");
    } else if (ulab_streq(view, "nodes_list")) {
        if (network == NULL || network->bff_id[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "NodesList requires a configured network");
            return ULAB_ERR;
        }
        ulab_json_escape(network->bff_id, network_esc,
                         sizeof(network_esc));
        snprintf(vars, sizeof(vars),
                 "{\"networkId\":\"%s\"}", network_esc);
        operation = "NodesList";
        query = BFF_CONSOLE_NODES_LIST;
    } else {
        snprintf(err->msg, sizeof(err->msg),
                 "unsupported console node view: %s",
                 view ? view : "");
        return ULAB_ERR;
    }

    if (bff_call(c, operation, query, vars, root, err)) {
        return ULAB_ERR;
    }
    nodes_view = dig(*root, "data", "nodesView");
    section = nodes_view ? json_object_get(nodes_view, "nodes") : NULL;
    if (console_section_ok(section, operation, "nodes", err)) {
        json_decref(*root);
        *root = NULL;
        return ULAB_ERR;
    }
    *arr = json_object_get(section, "nodes");
    if (*arr == NULL || json_is_null(*arr)) {
        *arr = NULL;
        return ULAB_OK;
    }
    if (!json_is_array(*arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s nodes section missing nodes list", operation);
        json_decref(*root);
        *root = NULL;
        *arr = NULL;
        return ULAB_ERR;
    }
    return ULAB_OK;
}

static int console_sites_array(bff_client_t *c,
                               const network_t *network,
                               json_t **root,
                               json_t **arr,
                               ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *sites_view;
    json_t *section;

    if (network == NULL || network->bff_id[0] == '\0' ||
        root == NULL || arr == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "SitesList requires a configured network");
        return ULAB_ERR;
    }
    *root = NULL;
    *arr = NULL;

    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
             network_esc);
    if (bff_call(c, "SitesList", BFF_CONSOLE_SITES_LIST,
                 vars, root, err)) {
        return ULAB_ERR;
    }
    sites_view = dig(*root, "data", "sitesView");
    section = sites_view ? json_object_get(sites_view, "sites") : NULL;
    if (console_section_ok(section, "SitesList", "sites", err)) {
        json_decref(*root);
        *root = NULL;
        return ULAB_ERR;
    }
    *arr = json_object_get(section, "sites");
    if (*arr == NULL || json_is_null(*arr)) {
        *arr = NULL;
        return ULAB_OK;
    }
    if (!json_is_array(*arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "SitesList sites section missing sites list");
        json_decref(*root);
        *root = NULL;
        *arr = NULL;
        return ULAB_ERR;
    }
    return ULAB_OK;
}

static int console_site_detail_view(bff_client_t *c,
                                    const char *site_id,
                                    json_t **root,
                                    json_t **site_view,
                                    ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char site_esc[ULAB_MAX_ID * 2];

    if (site_id == NULL || site_id[0] == '\0' ||
        root == NULL || site_view == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "NetworkSiteDetail requires a configured site");
        return ULAB_ERR;
    }
    *root = NULL;
    *site_view = NULL;

    ulab_json_escape(site_id, site_esc, sizeof(site_esc));
    snprintf(vars, sizeof(vars), "{\"siteId\":\"%s\"}", site_esc);
    if (bff_call(c, "NetworkSiteDetail", BFF_CONSOLE_SITE_DETAIL,
                 vars, root, err)) {
        return ULAB_ERR;
    }
    *site_view = dig(*root, "data", "siteView");
    if (*site_view == NULL || !json_is_object(*site_view)) {
        snprintf(err->msg, sizeof(err->msg),
                 "NetworkSiteDetail missing siteView");
        json_decref(*root);
        *root = NULL;
        *site_view = NULL;
        return ULAB_ERR;
    }
    return ULAB_OK;
}

static int console_site_detail_node_item(bff_client_t *c,
                                         const char *site_id,
                                         const char *node_id,
                                         json_t **root,
                                         json_t **item,
                                         ulab_error_t *err) {
    json_t *site_view;
    json_t *section;
    json_t *arr;

    if (console_site_detail_view(c, site_id, root, &site_view, err)) {
        return ULAB_ERR;
    }
    section = json_object_get(site_view, "nodes");
    if (console_section_ok(section, "NetworkSiteDetail", "nodes", err)) {
        json_decref(*root);
        *root = NULL;
        return ULAB_ERR;
    }
    arr = json_object_get(section, "nodes");
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "NetworkSiteDetail nodes section missing nodes list");
        json_decref(*root);
        *root = NULL;
        return ULAB_ERR;
    }
    *item = console_find_id(arr, node_id);
    if (*item == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "node id=%s is absent from NetworkSiteDetail", node_id);
        json_decref(*root);
        *root = NULL;
        return ULAB_ERR;
    }
    return ULAB_OK;
}

static int console_node_detail(bff_client_t *c,
                               const node_t *node,
                               json_t **root,
                               json_t **node_view,
                               ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char node_esc[ULAB_MAX_ID * 2];

    if (node == NULL || node->bff_id[0] == '\0' ||
        root == NULL || node_view == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "NodeDetail requires node and output storage");
        return ULAB_ERR;
    }
    *root = NULL;
    *node_view = NULL;
    ulab_json_escape(node->bff_id, node_esc, sizeof(node_esc));
    snprintf(vars, sizeof(vars), "{\"nodeId\":\"%s\"}", node_esc);
    if (bff_call(c, "NodeDetail", BFF_CONSOLE_NODE_DETAIL,
                 vars, root, err)) {
        return ULAB_ERR;
    }
    *node_view = dig(*root, "data", "nodeView");
    if (*node_view == NULL || !json_is_object(*node_view)) {
        snprintf(err->msg, sizeof(err->msg),
                 "NodeDetail missing nodeView");
        json_decref(*root);
        *root = NULL;
        return ULAB_ERR;
    }
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
    if (obj == NULL ||
        json_get_str(obj, "id", s->bff_id, sizeof(s->bff_id))) {
        snprintf(err->msg, sizeof(err->msg), "addSite missing id");
        json_decref(root);
        return ULAB_ERR;
    }
    json_get_optional_str(obj, "createdAt", s->created_at,
                          sizeof(s->created_at));
    json_get_optional_str(obj, "installDate", s->install_date,
                          sizeof(s->install_date));
    json_get_optional_str(obj, "location", s->location,
                          sizeof(s->location));

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


int bff_update_package_name(bff_client_t *c,
                            package_t *p,
                            const char *name,
                            ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char id_esc[ULAB_MAX_ID * 2];
    char name_esc[ULAB_MAX_NAME * 2];
    json_t *root;
    json_t *obj;
    char actual_name[ULAB_MAX_NAME];

    if (p == NULL || p->bff_id[0] == '\0' || name == NULL ||
        name[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "updatePackage requires package id and name");
        return ULAB_ERR;
    }
    ulab_json_escape(p->bff_id, id_esc, sizeof(id_esc));
    ulab_json_escape(name, name_esc, sizeof(name_esc));
    snprintf(vars, sizeof(vars),
             "{\"packageId\":\"%s\",\"data\":{\"name\":\"%s\","
             "\"active\":%s}}",
             id_esc, name_esc, p->active ? "true" : "false");
    root = NULL;
    if (bff_call(c, "updatePackage", BFF_UPDATE_PACKAGE, vars,
                 &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "updatePackage");
    if (obj == NULL ||
        json_get_str(obj, "name", actual_name, sizeof(actual_name)) ||
        !ulab_streq(actual_name, name)) {
        snprintf(err->msg, sizeof(err->msg),
                 "updatePackage returned unexpected name");
        json_decref(root);
        return ULAB_ERR;
    }
    if (ulab_copy(p->name, sizeof(p->name), name)) {
        snprintf(err->msg, sizeof(err->msg),
                 "updated package name is too long");
        json_decref(root);
        return ULAB_ERR;
    }
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

int bff_package_name_available(bff_client_t *c,
                               const char *name,
                               int *available,
                               ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char name_esc[ULAB_MAX_NAME * 2];
    json_t *root;
    json_t *obj;
    json_t *value;

    if (name == NULL || name[0] == '\0' || available == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "package name availability requires name");
        return ULAB_ERR;
    }
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

int bff_invalid_package_name_available(bff_client_t *c,
                                       const package_t *pkg,
                                       const char *variant,
                                       int *available,
                                       ulab_error_t *err) {
    char name[ULAB_MAX_NAME];

    if (pkg == NULL || variant == NULL || variant[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "package name availability requires package and variant");
        return ULAB_ERR;
    }
    invalid_package_name(pkg, variant, name, sizeof(name));
    return bff_package_name_available(c, name, available, err);
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

static int text_equals_ci(const char *left, const char *right) {
    return left != NULL && right != NULL && strcasecmp(left, right) == 0;
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
    json_t *obj;
    json_t *sims;
    bff_package_t actual;
    double revenue;
    char unit[ULAB_MAX_REF];
    int revenue_found;
    size_t i;

    if (pkg == NULL || pkg->bff_id[0] == '\0' || network == NULL ||
        network->bff_id[0] == '\0' || metrics == NULL || found == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "package metrics require package, network and output");
        return ULAB_ERR;
    }

    memset(metrics, 0, sizeof(*metrics));
    *found = 0;
    if (bff_get_package(c, pkg, &actual, err)) {
        return ULAB_ERR;
    }
    ulab_copy(metrics->package_id, sizeof(metrics->package_id),
              actual.uuid);
    *found = 1;

    revenue = 0;
    unit[0] = '\0';
    revenue_found = 0;
    if (bff_get_performance_report_cell(
            c, "package_performance", "daily", network->bff_id,
            pkg->bff_id, "revenue", &revenue, unit, sizeof(unit),
            &revenue_found, err)) {
        return ULAB_ERR;
    }
    metrics->revenue = revenue_found ? revenue : 0;

    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
             network_esc);
    root = NULL;
    if (bff_call(c, "getSimsByNetwork", BFF_GET_SIMS_BY_NETWORK,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getSimsByNetwork");
    sims = obj ? json_object_get(obj, "sims") : NULL;
    if (sims == NULL || !json_is_array(sims)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSimsByNetwork missing sims list");
        json_decref(root);
        return ULAB_ERR;
    }
    for (i = 0; i < json_array_size(sims); i++) {
        json_t *sim;
        json_t *assignment;
        json_t *active;
        char package_id[ULAB_MAX_ID];

        sim = json_array_get(sims, i);
        assignment = sim ? json_object_get(sim, "package") : NULL;
        if (assignment == NULL || json_is_null(assignment) ||
            json_get_str(assignment, "package_id", package_id,
                         sizeof(package_id)) ||
            !ulab_streq(package_id, pkg->bff_id)) {
            continue;
        }
        active = json_object_get(assignment, "is_currently_in_use");
        if (active != NULL && json_is_true(active)) {
            metrics->attach_count++;
        }
    }
    metrics->has_attach_count = 1;
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

static int section_has_error(json_t *section) {
    json_t *error;

    error = section ? json_object_get(section, "error") : NULL;
    return error != NULL && json_is_object(error);
}

static int payment_status_pending(const char *status) {
    return status != NULL &&
        (text_equals_ci(status, "pending") ||
         text_equals_ci(status, "processing") ||
         text_equals_ci(status, "initiated"));
}

static void payment_month_keys(char current[8], char previous[8]) {
    time_t now;
    struct tm current_tm;
    struct tm previous_tm;

    now = time(NULL);
    gmtime_r(&now, &current_tm);
    previous_tm = current_tm;
    if (previous_tm.tm_mon == 0) {
        previous_tm.tm_mon = 11;
        previous_tm.tm_year--;
    } else {
        previous_tm.tm_mon--;
    }
    strftime(current, 8, "%Y-%m", &current_tm);
    strftime(previous, 8, "%Y-%m", &previous_tm);
}

int bff_get_revenue_summary(bff_client_t *c,
                            const network_t *network,
                            bff_revenue_summary_t *summary,
                            ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    char current_month[8];
    char previous_month[8];
    json_t *root;
    json_t *obj;
    json_t *packages;
    size_t i;

    if (network == NULL || network->bff_id[0] == '\0' || summary == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "revenue summary requires network and output");
        return ULAB_ERR;
    }
    memset(summary, 0, sizeof(*summary));
    payment_month_keys(current_month, previous_month);

    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
             network_esc);
    root = NULL;
    if (bff_call(c, "getPackages", BFF_GET_PACKAGES, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getPackages");
    packages = obj ? json_object_get(obj, "packages") : NULL;
    if (packages == NULL || !json_is_array(packages)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPackages missing packages list");
        json_decref(root);
        return ULAB_ERR;
    }

    for (i = 0; i < json_array_size(packages); i++) {
        json_t *package_obj;
        char package_id[ULAB_MAX_ID];
        char package_esc[ULAB_MAX_ID * 2];
        json_t *payments_root;
        json_t *payments_obj;
        json_t *payments;
        size_t j;

        package_obj = json_array_get(packages, i);
        if (json_get_str(package_obj, "uuid", package_id,
                         sizeof(package_id))) {
            snprintf(err->msg, sizeof(err->msg),
                     "getPackages returned a package without uuid");
            json_decref(root);
            return ULAB_ERR;
        }
        ulab_json_escape(package_id, package_esc, sizeof(package_esc));
        snprintf(vars, sizeof(vars),
                 "{\"data\":{\"type\":\"package\","
                 "\"itemId\":\"%s\"}}", package_esc);
        payments_root = NULL;
        if (bff_call(c, "getPayments", BFF_GET_PAYMENTS, vars,
                     &payments_root, err)) {
            json_decref(root);
            return ULAB_ERR;
        }
        payments_obj = dig(payments_root, "data", "getPayments");
        payments = payments_obj ?
            json_object_get(payments_obj, "payments") : NULL;
        if (payments == NULL || !json_is_array(payments)) {
            snprintf(err->msg, sizeof(err->msg),
                     "getPayments missing payments list");
            json_decref(payments_root);
            json_decref(root);
            return ULAB_ERR;
        }
        for (j = 0; j < json_array_size(payments); j++) {
            json_t *payment;
            char status[ULAB_MAX_REF];
            char amount[ULAB_MAX_REF];
            char paid_at[ULAB_MAX_REF];
            double value;

            payment = json_array_get(payments, j);
            if (json_get_str(payment, "status", status, sizeof(status)) ||
                json_get_str(payment, "amount", amount, sizeof(amount))) {
                snprintf(err->msg, sizeof(err->msg),
                         "getPayments returned an invalid payment");
                json_decref(payments_root);
                json_decref(root);
                return ULAB_ERR;
            }
            paid_at[0] = '\0';
            json_get_optional_str(payment, "paidAt", paid_at,
                                  sizeof(paid_at));
            value = strtod(amount, NULL);
            if (payment_status_settled(status)) {
                summary->total_paid += value;
                if (strncmp(paid_at, current_month, 7) == 0) {
                    summary->month_paid += value;
                } else if (strncmp(paid_at, previous_month, 7) == 0) {
                    summary->previous_month_paid += value;
                }
            } else if (payment_status_pending(status)) {
                summary->total_pending += value;
            }
        }
        json_decref(payments_root);
    }
    json_decref(root);

    if (summary->previous_month_paid != 0) {
        summary->month_over_month_percent =
            ((summary->month_paid - summary->previous_month_paid) /
             summary->previous_month_paid) * 100.0;
    }
    return ULAB_OK;
}

int bff_get_package_kpis(bff_client_t *c,
                         const network_t *network,
                         bff_package_kpis_t *kpis,
                         ulab_error_t *err) {
    bff_kpi_value_t value;
    int found;

    if (network == NULL || network->bff_id[0] == '\0' || kpis == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "package KPIs require network and output");
        return ULAB_ERR;
    }
    memset(kpis, 0, sizeof(*kpis));

    found = 0;
    if (bff_get_kpi_value(c, "MRR", "daily", NULL,
                          network->bff_id, "network_id",
                          network->bff_id, &value, &found, err)) {
        return ULAB_ERR;
    }
    if (found) {
        kpis->mrr = value.value;
        kpis->has_mrr = 1;
    }

    found = 0;
    if (bff_get_kpi_value(c, "ARPU", "daily", NULL,
                          network->bff_id, "network_id",
                          network->bff_id, &value, &found, err)) {
        return ULAB_ERR;
    }
    if (found) {
        kpis->arpu = value.value;
        kpis->has_arpu = 1;
    }
    return ULAB_OK;
}

int bff_get_network_summary(bff_client_t *c,
                            const network_t *network,
                            bff_network_summary_t *summary,
                            ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *items;
    size_t i;

    if (network == NULL || network->bff_id[0] == '\0' || summary == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "network summary requires network and output");
        return ULAB_ERR;
    }
    memset(summary, 0, sizeof(*summary));
    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"networkId\":\"%s\"}}", network_esc);
    root = NULL;
    if (bff_call(c, "getSites", BFF_GET_SITES, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getSites");
    items = obj ? json_object_get(obj, "sites") : NULL;
    if (items == NULL || !json_is_array(items)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSites missing sites list");
        json_decref(root);
        return ULAB_ERR;
    }
    summary->sites_total = (uint32_t)json_array_size(items);
    json_decref(root);

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"networkId\":\"%s\"}}", network_esc);
    root = NULL;
    if (bff_call(c, "getNodes", BFF_GET_NODES, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getNodes");
    items = obj ? json_object_get(obj, "nodes") : NULL;
    if (items == NULL || !json_is_array(items)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getNodes missing nodes list");
        json_decref(root);
        return ULAB_ERR;
    }
    summary->nodes_total = (uint32_t)json_array_size(items);
    for (i = 0; i < json_array_size(items); i++) {
        json_t *node;
        json_t *status;
        char connectivity[ULAB_MAX_REF];

        node = json_array_get(items, i);
        status = node ? json_object_get(node, "status") : NULL;
        if (status != NULL &&
            json_get_str(status, "connectivity", connectivity,
                         sizeof(connectivity)) == ULAB_OK &&
            text_equals_ci(connectivity, "online")) {
            summary->nodes_online++;
        } else {
            summary->nodes_offline++;
        }
    }
    json_decref(root);

    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
             network_esc);
    root = NULL;
    if (bff_call(c, "getSubscribersByNetwork",
                 BFF_GET_SUBSCRIBERS_BY_NETWORK, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getSubscribersByNetwork");
    items = obj ? json_object_get(obj, "subscribers") : NULL;
    if (items == NULL || !json_is_array(items)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSubscribersByNetwork missing subscribers list");
        json_decref(root);
        return ULAB_ERR;
    }
    summary->subscribers_total = (uint32_t)json_array_size(items);
    for (i = 0; i < json_array_size(items); i++) {
        json_t *subscriber;
        json_t *sims;
        size_t j;
        int active;

        subscriber = json_array_get(items, i);
        sims = subscriber ? json_object_get(subscriber, "sim") : NULL;
        active = 0;
        if (sims != NULL && json_is_array(sims)) {
            for (j = 0; j < json_array_size(sims); j++) {
                json_t *sim;
                char status[ULAB_MAX_REF];

                sim = json_array_get(sims, j);
                if (json_get_str(sim, "status", status,
                                 sizeof(status)) == ULAB_OK &&
                    text_equals_ci(status, "active")) {
                    active = 1;
                    break;
                }
            }
        }
        if (active) {
            summary->subscribers_active++;
        } else {
            summary->subscribers_inactive++;
        }
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_get_nodes_count(bff_client_t *c,
                        const network_t *network,
                        uint32_t *count,
                        ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *nodes;

    if (network == NULL || network->bff_id[0] == '\0' || count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "node count requires network and output");
        return ULAB_ERR;
    }
    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars),
             "{\"data\":{\"networkId\":\"%s\"}}", network_esc);
    root = NULL;
    if (bff_call(c, "getNodes", BFF_GET_NODES, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getNodes");
    nodes = obj ? json_object_get(obj, "nodes") : NULL;
    if (nodes == NULL || !json_is_array(nodes)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getNodes missing nodes list");
        json_decref(root);
        return ULAB_ERR;
    }
    *count = (uint32_t)json_array_size(nodes);
    json_decref(root);
    return ULAB_OK;
}

int bff_get_component_inventory_summary(
    bff_client_t *c,
    bff_inventory_summary_t *summary,
    ulab_error_t *err) {
    json_t *root;
    json_t *obj;
    json_t *components;
    size_t i;

    if (summary == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "component inventory requires output storage");
        return ULAB_ERR;
    }
    memset(summary, 0, sizeof(*summary));
    root = NULL;
    if (bff_component_query(c, "all", &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getComponentsByUserId");
    components = obj ? json_object_get(obj, "components") : NULL;
    if (components == NULL || !json_is_array(components)) {
        snprintf(err->msg, sizeof(err->msg),
                 "GetComponentsByUserId returned no components list");
        json_decref(root);
        return ULAB_ERR;
    }
    summary->component_total = (uint32_t)json_array_size(components);
    for (i = 0; i < json_array_size(components); i++) {
        json_t *component;
        json_t *category_value;
        const char *category;

        component = json_array_get(components, i);
        category_value = component ?
            json_object_get(component, "category") : NULL;
        category = category_value && json_is_string(category_value) ?
            json_string_value(category_value) : NULL;
        if (category != NULL &&
            (ulab_streq(category, "access") ||
             ulab_streq(category, "backhaul") ||
             ulab_streq(category, "power") ||
             ulab_streq(category, "switch") ||
             ulab_streq(category, "spectrum"))) {
            summary->component_category_total++;
        }
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_get_sim_pool_summary(bff_client_t *c,
                             const char *sim_type,
                             bff_inventory_summary_t *summary,
                             ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char sim_type_esc[ULAB_MAX_REF * 2];
    json_t *root;
    json_t *pool;
    json_t *pool_stats;

    if (sim_type == NULL || sim_type[0] == '\0' || summary == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "SIM pool summary requires SIM type and output storage");
        return ULAB_ERR;
    }
    memset(summary, 0, sizeof(*summary));

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
        json_u32_field(pool_stats, "total", &summary->sim_total) ||
        json_u32_field(pool_stats, "available",
                       &summary->sim_available) ||
        json_u32_field(pool_stats, "consumed",
                       &summary->sim_consumed) ||
        json_u32_field(pool_stats, "failed", &summary->sim_failed)) {
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
        (text_equals_ci(status, "success") ||
         text_equals_ci(status, "paid") ||
         text_equals_ci(status, "completed") ||
         text_equals_ci(status, "processed"));
}

int bff_get_subscriber_payment_summary(
    bff_client_t *c,
    const subscriber_t *subscriber,
    bff_subscriber_billing_t *billing,
    ulab_error_t *err) {
    json_t *root;
    json_t *obj;
    json_t *payments;
    size_t i;

    if (subscriber == NULL || subscriber->bff_id[0] == '\0' ||
        subscriber->email[0] == '\0' || billing == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "subscriber payment summary requires subscriber and output");
        return ULAB_ERR;
    }
    root = NULL;
    if (bff_call(c, "getPayments", BFF_GET_PAYMENTS,
                 "{\"data\":{\"type\":\"package\"}}",
                 &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getPayments");
    payments = obj ? json_object_get(obj, "payments") : NULL;
    if (payments == NULL || !json_is_array(payments)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getPayments missing payments list");
        json_decref(root);
        return ULAB_ERR;
    }

    memset(billing, 0, sizeof(*billing));
    for (i = 0; i < json_array_size(payments); i++) {
        json_t *payment;
        char email[ULAB_MAX_NAME];
        char status[ULAB_MAX_REF];
        char amount[ULAB_MAX_REF];

        payment = json_array_get(payments, i);
        email[0] = '\0';
        json_get_optional_str(payment, "payerEmail", email,
                              sizeof(email));
        if (!ulab_streq(email, subscriber->email)) {
            continue;
        }
        if (json_get_str(payment, "status", status, sizeof(status)) ||
            json_get_str(payment, "amount", amount, sizeof(amount))) {
            snprintf(err->msg, sizeof(err->msg),
                     "getPayments returned an invalid payment");
            json_decref(root);
            return ULAB_ERR;
        }
        billing->payment_count++;
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

/* GraphQL KpiValuesInput field for a KPI scope dimension, or NULL when the
 * dimension is not a filter the gateway accepts. Sending an unknown filter is
 * an InvalidArgument, so an unmapped scope_key stays match-only. */
static const char *kpi_scope_filter_field(const char *scope_key) {
    if (scope_key == NULL || scope_key[0] == '\0') {
        return NULL;
    }
    if (ulab_streq(scope_key, "site_id")) {
        return "siteId";
    }
    if (ulab_streq(scope_key, "package_id")) {
        return "packageId";
    }
    if (ulab_streq(scope_key, "sim_package_id")) {
        return "simPackageId";
    }
    if (ulab_streq(scope_key, "iccid")) {
        return "iccid";
    }
    return NULL;
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
    char scope_esc[ULAB_MAX_ID * 2];
    const char *scope_field;
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
    /* Send the scope dimension as a FILTER, not just as a match predicate.
     * The gateway folds the answer to the filtered grain, so without this a
     * site/package/sim read comes back folded to the network and its own
     * dimension is absent from the response scope. */
    scope_field = kpi_scope_filter_field(scope_key);
    if (scope_field != NULL && scope_value != NULL && scope_value[0] != '\0') {
        ulab_json_escape(scope_value, scope_esc, sizeof(scope_esc));
        snprintf(optional + strlen(optional),
                 sizeof(optional) - strlen(optional),
                 ",\"%s\":\"%s\"", scope_field, scope_esc);
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


static int append_json_string_field(char *buffer,
                                    size_t buffer_len,
                                    const char *field,
                                    const char *value,
                                    ulab_error_t *err) {
    char escaped[ULAB_MAX_ID * 2];
    size_t used;
    size_t remaining;
    int written;

    if (value == NULL || value[0] == '\0') {
        return ULAB_OK;
    }
    used = strlen(buffer);
    if (used >= buffer_len) {
        snprintf(err->msg, sizeof(err->msg),
                 "JSON variables buffer is full");
        return ULAB_ERR;
    }
    remaining = buffer_len - used;
    ulab_json_escape(value, escaped, sizeof(escaped));
    written = snprintf(buffer + used, remaining,
                       ",\"%s\":\"%s\"", field, escaped);
    if (written < 0 || (size_t)written >= remaining) {
        snprintf(err->msg, sizeof(err->msg),
                 "JSON variables buffer is too long");
        return ULAB_ERR;
    }
    return ULAB_OK;
}

int bff_get_kpi_timeseries(bff_client_t *c,
                           const char *key,
                           const char *span,
                           const char *op,
                           const char *from,
                           const char *to,
                           const char *network_id,
                           const char *site_id,
                           bff_kpi_value_t values[],
                           size_t max_values,
                           size_t *value_count,
                           ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char key_esc[ULAB_MAX_REF * 2];
    char span_esc[ULAB_MAX_REF * 2];
    char optional[ULAB_MAX_QUERY / 2];
    json_t *root;
    json_t *obj;
    json_t *arr;
    size_t count;
    size_t i;
    int n;

    if (key == NULL || key[0] == '\0' || values == NULL ||
        max_values == 0 || value_count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getKpiTimeSeries requires key and output storage");
        return ULAB_ERR;
    }

    ulab_json_escape(key, key_esc, sizeof(key_esc));
    ulab_json_escape(span && span[0] ? span : "daily",
                     span_esc, sizeof(span_esc));
    optional[0] = '\0';
    if (append_json_string_field(optional, sizeof(optional), "op", op, err) ||
        append_json_string_field(optional, sizeof(optional), "from", from,
                                 err) ||
        append_json_string_field(optional, sizeof(optional), "to", to,
                                 err) ||
        append_json_string_field(optional, sizeof(optional), "networkId",
                                 network_id, err) ||
        append_json_string_field(optional, sizeof(optional), "siteId",
                                 site_id, err)) {
        return ULAB_ERR;
    }

    n = snprintf(vars, sizeof(vars),
                 "{\"data\":{\"keys\":[\"%s\"],"
                 "\"span\":\"%s\"%s}}",
                 key_esc, span_esc, optional);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getKpiTimeSeries variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "getKpiTimeSeries", BFF_GET_KPI_TIMESERIES,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getKpiTimeSeries");
    arr = obj ? json_object_get(obj, "values") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getKpiTimeSeries missing values list");
        json_decref(root);
        return ULAB_ERR;
    }

    count = json_array_size(arr);
    if (count > max_values) {
        snprintf(err->msg, sizeof(err->msg),
                 "getKpiTimeSeries returned %zu values, max=%zu",
                 count, max_values);
        json_decref(root);
        return ULAB_ERR;
    }
    for (i = 0; i < count; i++) {
        if (parse_kpi_value(json_array_get(arr, i), &values[i], err)) {
            json_decref(root);
            return ULAB_ERR;
        }
    }
    *value_count = count;
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
    if (obj == NULL ||
        json_get_str(obj, "id", ue->bff_id, sizeof(ue->bff_id)) ||
        json_get_str(obj, "iccid", ue->iccid, sizeof(ue->iccid))) {
        snprintf(err->msg, sizeof(err->msg),
                 "allocateSim missing id or iccid");
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
        act = it ? json_object_get(it, "is_currently_in_use") : NULL;
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

int bff_get_node_status_for_view(bff_client_t *c,
                                 const network_t *network,
                                 const node_t *node,
                                 const char *view,
                                 bff_node_status_t *status,
                                 ulab_error_t *err) {
    json_t *root;
    json_t *arr;
    json_t *item;
    json_t *node_status;

    if (view == NULL || view[0] == '\0') {
        return bff_get_node_status(c, node, status, err);
    }
    if (node == NULL || node->bff_id[0] == '\0' || status == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "node status view requires node and output storage");
        return ULAB_ERR;
    }

    root = NULL;
    arr = NULL;
    if (console_nodes_array(c, network, view, &root, &arr, err)) {
        return ULAB_ERR;
    }
    item = console_find_id(arr, node->bff_id);
    if (item == NULL) {
        ulab_copy(status->id, sizeof(status->id), node->bff_id);
        json_decref(root);
        return ULAB_OK;
    }
    node_status = json_object_get(item, "status");
    if (node_status == NULL || !json_is_object(node_status) ||
        json_get_str(item, "id", status->id, sizeof(status->id)) ||
        json_get_str(node_status, "connectivity", status->connectivity,
                     sizeof(status->connectivity)) ||
        json_get_str(node_status, "state", status->state,
                     sizeof(status->state))) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s node id=%s is missing status fields",
                 view, node->bff_id);
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


static int json_get_required_bool(json_t *obj, const char *key, int *out) {
    json_t *value;

    value = obj ? json_object_get(obj, key) : NULL;
    if (value == NULL || !json_is_boolean(value) || out == NULL) {
        return ULAB_ERR;
    }
    *out = json_is_true(value);
    return ULAB_OK;
}

static int parse_operation(json_t *obj, bff_operation_t *operation,
                           ulab_error_t *err) {
    if (obj == NULL || !json_is_object(obj) || operation == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "operation status has invalid operation");
        return ULAB_ERR;
    }

    memset(operation, 0, sizeof(*operation));
    if (json_get_str(obj, "id", operation->id,
                     sizeof(operation->id)) ||
        json_get_str(obj, "type", operation->type,
                     sizeof(operation->type)) ||
        json_get_str(obj, "status", operation->status,
                     sizeof(operation->status))) {
        snprintf(err->msg, sizeof(err->msg),
                 "operation status is missing id/type/status");
        return ULAB_ERR;
    }
    json_get_optional_str(obj, "requestedBy", operation->requested_by,
                          sizeof(operation->requested_by));
    json_get_optional_str(obj, "startedAt", operation->started_at,
                          sizeof(operation->started_at));
    json_get_optional_str(obj, "leaseExpiresAt",
                          operation->lease_expires_at,
                          sizeof(operation->lease_expires_at));
    return ULAB_OK;
}

static int parse_node_operation_status(
    json_t *obj,
    bff_node_operation_status_t *status,
    ulab_error_t *err) {
    json_t *operation;

    if (obj == NULL || !json_is_object(obj) || status == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "node operation status has invalid payload");
        return ULAB_ERR;
    }

    memset(status, 0, sizeof(*status));
    if (json_get_str(obj, "nodeId", status->node_id,
                     sizeof(status->node_id)) ||
        json_get_required_bool(obj, "busy", &status->busy)) {
        snprintf(err->msg, sizeof(err->msg),
                 "node operation status is missing nodeId/busy");
        return ULAB_ERR;
    }
    json_get_optional_str(obj, "type", status->node_type,
                          sizeof(status->node_type));

    operation = json_object_get(obj, "operation");
    if (operation != NULL && !json_is_null(operation)) {
        if (parse_operation(operation, &status->operation, err)) {
            return ULAB_ERR;
        }
        status->has_operation = 1;
    }
    return ULAB_OK;
}

static int parse_action_availability(json_t *obj,
                                     bff_action_availability_t *action,
                                     ulab_error_t *err) {
    if (obj == NULL || !json_is_object(obj) || action == NULL ||
        json_get_required_bool(obj, "available", &action->available)) {
        snprintf(err->msg, sizeof(err->msg),
                 "site operation status has invalid action availability");
        return ULAB_ERR;
    }
    json_get_optional_str(obj, "reason", action->reason,
                          sizeof(action->reason));
    return ULAB_OK;
}

static int console_software_list(bff_client_t *c,
                                 const node_t *node,
                                 bff_software_t software[],
                                 size_t max_software,
                                 size_t *software_count,
                                 ulab_error_t *err) {
    json_t *root;
    json_t *node_view;
    json_t *section;
    json_t *softwares;
    json_t *arr;
    size_t i;

    root = NULL;
    node_view = NULL;
    if (console_node_detail(c, node, &root, &node_view, err)) {
        return ULAB_ERR;
    }
    section = json_object_get(node_view, "software");
    if (console_section_ok(section, "NodeDetail", "software", err)) {
        json_decref(root);
        return ULAB_ERR;
    }
    softwares = json_object_get(section, "softwares");
    arr = softwares && json_is_object(softwares) ?
        json_object_get(softwares, "software") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "NodeDetail software section missing software list");
        json_decref(root);
        return ULAB_ERR;
    }

    *software_count = json_array_size(arr);
    if (*software_count > max_software) {
        snprintf(err->msg, sizeof(err->msg),
                 "NodeDetail returned %zu software rows; maximum is %zu",
                 *software_count, max_software);
        json_decref(root);
        return ULAB_ERR;
    }
    for (i = 0; i < *software_count; i++) {
        json_t *item;
        bff_software_t *row;

        item = json_array_get(arr, i);
        row = &software[i];
        memset(row, 0, sizeof(*row));
        if (json_get_str(item, "releaseDate", row->release_date,
                         sizeof(row->release_date)) ||
            json_get_str(item, "status", row->status,
                         sizeof(row->status)) ||
            json_get_str(item, "currentVersion", row->current_version,
                         sizeof(row->current_version)) ||
            json_get_str(item, "name", row->name, sizeof(row->name))) {
            snprintf(err->msg, sizeof(err->msg),
                     "NodeDetail returned incomplete software row");
            json_decref(root);
            return ULAB_ERR;
        }
        json_get_optional_str(item, "desiredVersion",
                              row->desired_version,
                              sizeof(row->desired_version));
        ulab_copy(row->node_id, sizeof(row->node_id), node->bff_id);
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_get_software_list(bff_client_t *c,
                          const node_t *node,
                          const char *view,
                          bff_software_t software[],
                          size_t max_software,
                          size_t *software_count,
                          ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char node_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *arr;
    size_t i;
    int n;

    if (node == NULL || node->bff_id[0] == '\0' ||
        software == NULL || max_software == 0 || software_count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSoftwares list requires node and output storage");
        return ULAB_ERR;
    }
    if (view != NULL && view[0] != '\0') {
        if (!ulab_streq(view, "node_detail")) {
            snprintf(err->msg, sizeof(err->msg),
                     "unsupported software view: %s", view);
            return ULAB_ERR;
        }
        return console_software_list(c, node, software, max_software,
                                     software_count, err);
    }
    ulab_json_escape(node->bff_id, node_esc, sizeof(node_esc));
    n = snprintf(vars, sizeof(vars),
                 "{\"data\":{\"name\":\"\",\"nodeId\":\"%s\","
                 "\"status\":\"unknown\"}}", node_esc);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSoftwares variables too long");
        return ULAB_ERR;
    }
    root = NULL;
    if (bff_call(c, "getSoftwares", BFF_GET_SOFTWARES, vars,
                 &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getSoftwares");
    arr = obj ? json_object_get(obj, "software") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSoftwares missing software list");
        json_decref(root);
        return ULAB_ERR;
    }
    *software_count = json_array_size(arr);
    if (*software_count > max_software) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSoftwares returned %zu rows; maximum is %zu",
                 *software_count, max_software);
        json_decref(root);
        return ULAB_ERR;
    }
    for (i = 0; i < *software_count; i++) {
        json_t *item;
        bff_software_t *row;

        item = json_array_get(arr, i);
        row = &software[i];
        memset(row, 0, sizeof(*row));
        if (json_get_str(item, "id", row->id, sizeof(row->id)) ||
            json_get_str(item, "releaseDate", row->release_date,
                         sizeof(row->release_date)) ||
            json_get_str(item, "nodeId", row->node_id,
                         sizeof(row->node_id)) ||
            json_get_str(item, "status", row->status,
                         sizeof(row->status)) ||
            json_get_str(item, "currentVersion", row->current_version,
                         sizeof(row->current_version)) ||
            json_get_str(item, "desiredVersion", row->desired_version,
                         sizeof(row->desired_version)) ||
            json_get_str(item, "name", row->name, sizeof(row->name))) {
            snprintf(err->msg, sizeof(err->msg),
                     "getSoftwares returned incomplete software row");
            json_decref(root);
            return ULAB_ERR;
        }
        json_get_optional_str(item, "createdAt", row->created_at,
                              sizeof(row->created_at));
        json_get_optional_str(item, "updatedAt", row->updated_at,
                              sizeof(row->updated_at));
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_get_software(bff_client_t *c,
                     const node_t *node,
                     const char *app,
                     const char *view,
                     bff_software_t *software,
                     int *found,
                     ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char app_esc[ULAB_MAX_NAME * 2];
    char node_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *arr;
    size_t i;
    int n;

    if (node == NULL || node->bff_id[0] == '\0' || app == NULL ||
        app[0] == '\0' || software == NULL || found == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSoftwares requires node, app and output storage");
        return ULAB_ERR;
    }

    memset(software, 0, sizeof(*software));
    *found = 0;
    if (view != NULL && view[0] != '\0') {
        bff_software_t rows[ULAB_MAX_LIST];
        size_t count;

        if (!ulab_streq(view, "node_detail")) {
            snprintf(err->msg, sizeof(err->msg),
                     "unsupported software view: %s", view);
            return ULAB_ERR;
        }
        memset(rows, 0, sizeof(rows));
        count = 0;
        if (console_software_list(c, node, rows, ULAB_MAX_LIST,
                                  &count, err)) {
            return ULAB_ERR;
        }
        for (i = 0; i < count; i++) {
            if (ulab_streq(rows[i].name, app)) {
                *software = rows[i];
                *found = 1;
                break;
            }
        }
        return ULAB_OK;
    }
    ulab_json_escape(app, app_esc, sizeof(app_esc));
    ulab_json_escape(node->bff_id, node_esc, sizeof(node_esc));
    n = snprintf(vars, sizeof(vars),
                 "{\"data\":{\"name\":\"%s\","
                 "\"nodeId\":\"%s\",\"status\":\"unknown\"}}",
                 app_esc, node_esc);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSoftwares variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "getSoftwares", BFF_GET_SOFTWARES, vars,
                 &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getSoftwares");
    arr = obj ? json_object_get(obj, "software") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSoftwares missing software list");
        json_decref(root);
        return ULAB_ERR;
    }

    for (i = 0; i < json_array_size(arr); i++) {
        json_t *item;
        json_t *name;
        const char *value;

        item = json_array_get(arr, i);
        name = item ? json_object_get(item, "name") : NULL;
        value = name && json_is_string(name) ?
            json_string_value(name) : NULL;
        if (value == NULL || !ulab_streq(value, app)) {
            continue;
        }
        if (json_get_str(item, "id", software->id,
                         sizeof(software->id)) ||
            json_get_str(item, "releaseDate", software->release_date,
                         sizeof(software->release_date)) ||
            json_get_str(item, "nodeId", software->node_id,
                         sizeof(software->node_id)) ||
            json_get_str(item, "status", software->status,
                         sizeof(software->status)) ||
            json_get_str(item, "currentVersion",
                         software->current_version,
                         sizeof(software->current_version)) ||
            json_get_str(item, "desiredVersion",
                         software->desired_version,
                         sizeof(software->desired_version)) ||
            json_get_str(item, "name", software->name,
                         sizeof(software->name))) {
            snprintf(err->msg, sizeof(err->msg),
                     "getSoftwares returned incomplete software row");
            json_decref(root);
            return ULAB_ERR;
        }
        json_get_optional_str(item, "createdAt", software->created_at,
                              sizeof(software->created_at));
        json_get_optional_str(item, "updatedAt", software->updated_at,
                              sizeof(software->updated_at));
        *found = 1;
        break;
    }

    json_decref(root);
    return ULAB_OK;
}

int bff_get_software_count(bff_client_t *c,
                           const node_t *node,
                           size_t *count,
                           ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char node_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *arr;
    int n;

    if (node == NULL || node->bff_id[0] == '\0' || count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSoftwares count requires node and output storage");
        return ULAB_ERR;
    }

    ulab_json_escape(node->bff_id, node_esc, sizeof(node_esc));
    n = snprintf(vars, sizeof(vars),
                 "{\"data\":{\"name\":\"\",\"nodeId\":\"%s\","
                 "\"status\":\"unknown\"}}", node_esc);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSoftwares count variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "getSoftwares", BFF_GET_SOFTWARES, vars,
                 &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getSoftwares");
    arr = obj ? json_object_get(obj, "software") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSoftwares missing software list");
        json_decref(root);
        return ULAB_ERR;
    }
    *count = json_array_size(arr);
    json_decref(root);
    return ULAB_OK;
}

int bff_get_node_operation_status(
    bff_client_t *c,
    const node_t *node,
    bff_node_operation_status_t *status,
    ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char node_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    int n;

    if (node == NULL || node->bff_id[0] == '\0' || status == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getNodeOperationStatus requires node and output");
        return ULAB_ERR;
    }
    ulab_json_escape(node->bff_id, node_esc, sizeof(node_esc));
    n = snprintf(vars, sizeof(vars), "{\"nodeId\":\"%s\"}",
                 node_esc);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getNodeOperationStatus variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "getNodeOperationStatus",
                 BFF_GET_NODE_OPERATION_STATUS, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getNodeOperationStatus");
    if (parse_node_operation_status(obj, status, err)) {
        json_decref(root);
        return ULAB_ERR;
    }
    json_decref(root);
    return ULAB_OK;
}

int bff_get_site_operation_status(
    bff_client_t *c,
    const site_t *site,
    bff_site_operation_status_t *status,
    ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char site_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *obj;
    json_t *nodes;
    json_t *actions;
    size_t count;
    size_t i;
    int n;

    if (site == NULL || site->bff_id[0] == '\0' || status == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSiteOperationStatus requires site and output");
        return ULAB_ERR;
    }
    ulab_json_escape(site->bff_id, site_esc, sizeof(site_esc));
    n = snprintf(vars, sizeof(vars), "{\"siteId\":\"%s\"}",
                 site_esc);
    if (n < 0 || (size_t)n >= sizeof(vars)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSiteOperationStatus variables too long");
        return ULAB_ERR;
    }

    root = NULL;
    if (bff_call(c, "getSiteOperationStatus",
                 BFF_GET_SITE_OPERATION_STATUS, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getSiteOperationStatus");
    memset(status, 0, sizeof(*status));
    if (obj == NULL ||
        json_get_str(obj, "siteId", status->site_id,
                     sizeof(status->site_id)) ||
        json_get_required_bool(obj, "busy", &status->busy) ||
        json_get_required_bool(obj, "degraded", &status->degraded)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSiteOperationStatus missing site fields");
        json_decref(root);
        return ULAB_ERR;
    }

    nodes = json_object_get(obj, "nodes");
    actions = json_object_get(obj, "actions");
    if (nodes == NULL || !json_is_array(nodes) || actions == NULL ||
        !json_is_object(actions)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSiteOperationStatus missing nodes/actions");
        json_decref(root);
        return ULAB_ERR;
    }
    count = json_array_size(nodes);
    if (count > ULAB_MAX_LIST) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSiteOperationStatus returned too many nodes");
        json_decref(root);
        return ULAB_ERR;
    }
    for (i = 0; i < count; i++) {
        if (parse_node_operation_status(json_array_get(nodes, i),
                                        &status->nodes[i], err)) {
            json_decref(root);
            return ULAB_ERR;
        }
    }
    status->node_count = count;
    if (parse_action_availability(json_object_get(actions, "restartSite"),
                                  &status->restart_site, err) ||
        parse_action_availability(json_object_get(actions, "rf"),
                                  &status->rf, err) ||
        parse_action_availability(json_object_get(actions, "service"),
                                  &status->service, err)) {
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);
    return ULAB_OK;
}

int bff_console_network_loads(bff_client_t *c,
                              const network_t *net,
                              ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *network;
    bff_network_summary_t summary;
    char id[ULAB_MAX_ID];

    if (net == NULL || net->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "console network load requires network id");
        return ULAB_ERR;
    }
    ulab_json_escape(net->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
             network_esc);
    root = NULL;
    if (bff_call(c, "getNetwork", BFF_GET_NETWORK, vars, &root, err)) {
        return ULAB_ERR;
    }
    network = dig(root, "data", "getNetwork");
    if (network == NULL ||
        json_get_str(network, "id", id, sizeof(id)) ||
        !ulab_streq(id, net->bff_id)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getNetwork returned the wrong network");
        json_decref(root);
        return ULAB_ERR;
    }
    json_decref(root);

    return bff_get_network_summary(c, net, &summary, err);
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
    char network_esc[ULAB_MAX_ID * 2];

    if (root != NULL) {
        *root = NULL;
    }
    ulab_json_escape(network_id ? network_id : "", network_esc,
                     sizeof(network_esc));
    snprintf(vars, sizeof(vars),
             "{\"data\":{\"networkId\":\"%s\"}}", network_esc);
    return bff_call(c, "getSites", BFF_GET_SITES, vars, root, err);
}

static int backend_get_nodes_for_site(bff_client_t *c, const char *site_id,
                                      json_t **root, ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];

    snprintf(vars, sizeof(vars), "{\"siteId\":\"%s\"}",
             site_id ? site_id : "");
    return bff_call(c, "getNodesForSite", BFF_GET_NODES_FOR_SITE, vars,
                    root, err);
}


static int direct_network_list_call(bff_client_t *c,
                                    const char *target,
                                    const network_t *network,
                                    json_t **root,
                                    const char **envelope,
                                    const char **list_key,
                                    ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];

    if (root == NULL || envelope == NULL || list_key == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "list count has invalid output storage");
        return ULAB_ERR;
    }
    *root = NULL;
    if (ulab_streq(target, "networks")) {
        *envelope = "getNetworks";
        *list_key = "networks";
        return backend_get_networks(c, root, err);
    }
    if ((ulab_streq(target, "packages") ||
         ulab_streq(target, "plans")) && network == NULL) {
        *envelope = "getPackages";
        *list_key = "packages";
        return bff_call(c, "getPackages", BFF_GET_PACKAGES,
                        "{}", root, err);
    }
    if (network == NULL || network->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "list count target=%s requires a network", target);
        return ULAB_ERR;
    }

    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    if (ulab_streq(target, "sites")) {
        *envelope = "getSites";
        *list_key = "sites";
        return backend_get_sites(c, network->bff_id, root, err);
    }
    if (ulab_streq(target, "nodes")) {
        snprintf(vars, sizeof(vars),
                 "{\"data\":{\"networkId\":\"%s\"}}",
                 network_esc);
        *envelope = "getNodes";
        *list_key = "nodes";
        return bff_call(c, "getNodes", BFF_GET_NODES, vars, root, err);
    }
    if (ulab_streq(target, "subscribers") ||
        ulab_streq(target, "customers")) {
        snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
                 network_esc);
        *envelope = "getSubscribersByNetwork";
        *list_key = "subscribers";
        return bff_call(c, "getSubscribersByNetwork",
                        BFF_GET_SUBSCRIBERS_BY_NETWORK, vars, root, err);
    }
    if (ulab_streq(target, "packages") ||
        ulab_streq(target, "plans")) {
        snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
                 network_esc);
        *envelope = "getPackages";
        *list_key = "packages";
        return bff_call(c, "getPackages", BFF_GET_PACKAGES,
                        vars, root, err);
    }
    if (ulab_streq(target, "sims")) {
        snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
                 network_esc);
        *envelope = "getSimsByNetwork";
        *list_key = "sims";
        return bff_call(c, "getSimsByNetwork", BFF_GET_SIMS_BY_NETWORK,
                        vars, root, err);
    }

    snprintf(err->msg, sizeof(err->msg),
             "unsupported direct list target: %s", target);
    return ULAB_ERR;
}

int bff_get_list_count(bff_client_t *c,
                       const char *target,
                       const network_t *network,
                       size_t *count,
                       ulab_error_t *err) {
    json_t *root;
    json_t *obj;
    json_t *arr;
    const char *envelope;
    const char *list_key;

    if (target == NULL || target[0] == '\0' || count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "list count requires target and output storage");
        return ULAB_ERR;
    }

    root = NULL;
    envelope = NULL;
    list_key = NULL;
    if (direct_network_list_call(c, target, network, &root,
                                 &envelope, &list_key, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", envelope);
    arr = obj ? json_object_get(obj, list_key) : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s missing %s list", envelope, list_key);
        json_decref(root);
        return ULAB_ERR;
    }
    *count = json_array_size(arr);
    json_decref(root);
    return ULAB_OK;
}

int bff_get_site_list_count(bff_client_t *c,
                            const network_t *network,
                            const char *view,
                            size_t *count,
                            ulab_error_t *err) {
    json_t *root;
    json_t *arr;

    if (network == NULL || network->bff_id[0] == '\0' || count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "site list count requires network");
        return ULAB_ERR;
    }
    if (view == NULL || view[0] == '\0') {
        return bff_get_list_count(c, "sites", network, count, err);
    }
    if (!ulab_streq(view, "sites_list")) {
        snprintf(err->msg, sizeof(err->msg),
                 "unsupported console site list view: %s", view);
        return ULAB_ERR;
    }

    root = NULL;
    arr = NULL;
    if (console_sites_array(c, network, &root, &arr, err)) {
        return ULAB_ERR;
    }
    *count = arr != NULL ? json_array_size(arr) : 0;
    json_decref(root);
    return ULAB_OK;
}


static const char *node_kind_filter(const char *node_type) {
    if (node_type == NULL || node_type[0] == '\0') return "";
    if (ulab_streq(node_type, ULAB_NODE_TOWER)) {
        return ULAB_NODE_KIND_TOWER;
    }
    if (ulab_streq(node_type, ULAB_NODE_AMPLIFIER)) {
        return ULAB_NODE_KIND_AMPLIFIER;
    }
    if (ulab_streq(node_type, ULAB_NODE_CONTROLLER)) {
        return ULAB_NODE_KIND_CONTROLLER;
    }
    return node_type;
}

int bff_get_node_list_count(bff_client_t *c,
                            const network_t *network,
                            const char *node_type,
                            const char *view,
                            size_t *count,
                            ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    char type_esc[ULAB_MAX_REF * 2];
    const char *kind;
    json_t *root;
    json_t *obj;
    json_t *arr;
    size_t i;

    if (network == NULL || network->bff_id[0] == '\0' || count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "node list count requires network");
        return ULAB_ERR;
    }
    kind = node_kind_filter(node_type);
    if (view != NULL && view[0] != '\0') {
        root = NULL;
        arr = NULL;
        if (console_nodes_array(c, network, view, &root, &arr, err)) {
            return ULAB_ERR;
        }
        *count = 0;
        for (i = 0; arr != NULL && i < json_array_size(arr); i++) {
            json_t *item;
            char actual_type[ULAB_MAX_REF];

            item = json_array_get(arr, i);
            actual_type[0] = '\0';
            json_get_optional_str(item, "type", actual_type,
                                  sizeof(actual_type));
            if (kind[0] == '\0' || ulab_streq(actual_type, kind)) {
                (*count)++;
            }
        }
        json_decref(root);
        return ULAB_OK;
    }

    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    ulab_json_escape(kind, type_esc, sizeof(type_esc));
    if (kind[0] != '\0') {
        snprintf(vars, sizeof(vars),
                 "{\"data\":{\"networkId\":\"%s\","
                 "\"type\":\"%s\"}}", network_esc, type_esc);
    } else {
        snprintf(vars, sizeof(vars),
                 "{\"data\":{\"networkId\":\"%s\"}}",
                 network_esc);
    }
    root = NULL;
    if (bff_call(c, "getNodes", BFF_GET_NODES, vars, &root, err)) {
        return ULAB_ERR;
    }
    obj = dig(root, "data", "getNodes");
    arr = obj ? json_object_get(obj, "nodes") : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getNodes missing nodes list");
        json_decref(root);
        return ULAB_ERR;
    }
    *count = json_array_size(arr);
    json_decref(root);
    return ULAB_OK;
}

int bff_get_site_node_count(bff_client_t *c,
                            const site_t *site,
                            const char *node_type,
                            const char *view,
                            size_t *count,
                            ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char site_esc[ULAB_MAX_ID * 2];
    const char *kind;
    json_t *root;
    json_t *obj;
    json_t *arr;
    size_t i;

    if (site == NULL || site->bff_id[0] == '\0' || count == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "site node count requires configured site");
        return ULAB_ERR;
    }
    kind = node_kind_filter(node_type);
    root = NULL;
    arr = NULL;

    if (view != NULL && view[0] != '\0') {
        json_t *site_view;
        json_t *section;

        if (!ulab_streq(view, "site_detail")) {
            snprintf(err->msg, sizeof(err->msg),
                     "unsupported console site node view: %s", view);
            return ULAB_ERR;
        }
        if (console_site_detail_view(c, site->bff_id, &root,
                                     &site_view, err)) {
            return ULAB_ERR;
        }
        section = json_object_get(site_view, "nodes");
        if (console_section_ok(section, "NetworkSiteDetail",
                               "nodes", err)) {
            json_decref(root);
            return ULAB_ERR;
        }
        arr = json_object_get(section, "nodes");
        if (arr == NULL || json_is_null(arr)) {
            *count = 0;
            json_decref(root);
            return ULAB_OK;
        }
        if (!json_is_array(arr)) {
            snprintf(err->msg, sizeof(err->msg),
                     "NetworkSiteDetail nodes section missing nodes list");
            json_decref(root);
            return ULAB_ERR;
        }
    } else {
        ulab_json_escape(site->bff_id, site_esc, sizeof(site_esc));
        snprintf(vars, sizeof(vars), "{\"siteId\":\"%s\"}", site_esc);
        if (bff_call(c, "getNodesForSite", BFF_GET_NODES_FOR_SITE,
                     vars, &root, err)) {
            return ULAB_ERR;
        }
        obj = dig(root, "data", "getNodesForSite");
        arr = obj ? json_object_get(obj, "nodes") : NULL;
        if (arr == NULL || !json_is_array(arr)) {
            snprintf(err->msg, sizeof(err->msg),
                     "getNodesForSite missing nodes list");
            json_decref(root);
            return ULAB_ERR;
        }
    }

    *count = 0;
    for (i = 0; i < json_array_size(arr); i++) {
        json_t *item;
        char actual_type[ULAB_MAX_REF];

        item = json_array_get(arr, i);
        actual_type[0] = '\0';
        json_get_optional_str(item, "type", actual_type,
                              sizeof(actual_type));
        if (kind[0] == '\0' || ulab_streq(actual_type, kind)) {
            (*count)++;
        }
    }
    json_decref(root);
    return ULAB_OK;
}


int bff_get_console_site_node_counts(bff_client_t *c,
                                     const network_t *network,
                                     const site_t *site,
                                     bff_site_node_counts_t *counts,
                                     ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char network_esc[ULAB_MAX_ID * 2];
    json_t *root;
    json_t *sites_view;
    json_t *section;
    json_t *arr;
    size_t i;

    if (network == NULL || network->bff_id[0] == '\0' ||
        site == NULL || site->bff_id[0] == '\0' || counts == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "site node counts require configured network and site");
        return ULAB_ERR;
    }
    memset(counts, 0, sizeof(*counts));
    ulab_json_escape(network->bff_id, network_esc, sizeof(network_esc));
    snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
             network_esc);
    root = NULL;
    if (bff_call(c, "SitesList", BFF_CONSOLE_SITE_NODE_COUNTS,
                 vars, &root, err)) {
        return ULAB_ERR;
    }
    sites_view = dig(root, "data", "sitesView");
    section = sites_view ? json_object_get(sites_view, "nodeCounts") : NULL;
    if (console_section_ok(section, "SitesList", "nodeCounts", err)) {
        json_decref(root);
        return ULAB_ERR;
    }
    arr = json_object_get(section, "counts");
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "SitesList nodeCounts section missing counts list");
        json_decref(root);
        return ULAB_ERR;
    }
    for (i = 0; i < json_array_size(arr); i++) {
        json_t *item;
        const char *site_id;
        long long total;
        long long online;
        long long offline;

        item = json_array_get(arr, i);
        site_id = json_string_value(json_object_get(item, "siteId"));
        if (site_id == NULL || !ulab_streq(site_id, site->bff_id)) {
            continue;
        }
        if (!json_is_integer(json_object_get(item, "total")) ||
            !json_is_integer(json_object_get(item, "online")) ||
            !json_is_integer(json_object_get(item, "offline"))) {
            snprintf(err->msg, sizeof(err->msg),
                     "SitesList nodeCounts contains invalid values");
            json_decref(root);
            return ULAB_ERR;
        }
        total = json_integer_value(json_object_get(item, "total"));
        online = json_integer_value(json_object_get(item, "online"));
        offline = json_integer_value(json_object_get(item, "offline"));
        if (total < 0 || online < 0 || offline < 0) {
            snprintf(err->msg, sizeof(err->msg),
                     "SitesList nodeCounts contains negative values");
            json_decref(root);
            return ULAB_ERR;
        }
        counts->total = (size_t)total;
        counts->online = (size_t)online;
        counts->offline = (size_t)offline;
        json_decref(root);
        return ULAB_OK;
    }
    snprintf(err->msg, sizeof(err->msg),
             "site id=%s is absent from SitesList nodeCounts",
             site->bff_id);
    json_decref(root);
    return ULAB_ERR;
}


typedef struct {
    char id[ULAB_MAX_ID];
    char name[ULAB_MAX_NAME];
    char network_id[ULAB_MAX_ID];
    char site_id[ULAB_MAX_ID];
    char type[ULAB_MAX_REF];
    char latitude[ULAB_MAX_REF];
    char longitude[ULAB_MAX_REF];
    char location[ULAB_MAX_NAME];
    char created_at[ULAB_MAX_REF];
    char install_date[ULAB_MAX_REF];
    char email[ULAB_MAX_NAME];
    char phone[ULAB_MAX_REF];
    char connectivity[ULAB_MAX_REF];
    char state[ULAB_MAX_REF];
    char sim_id[ULAB_MAX_ID];
    char sim_status[ULAB_MAX_REF];
    char package_id[ULAB_MAX_ID];
    int  active;
    int  has_active;
    int  deactivated;
    int  has_deactivated;
} bff_entity_snapshot_t;

static json_t *json_array_find_id(json_t *arr, const char *id) {
    size_t i;

    if (arr == NULL || !json_is_array(arr) || id == NULL) {
        return NULL;
    }
    for (i = 0; i < json_array_size(arr); i++) {
        json_t *item;
        json_t *value;
        const char *actual;

        item = json_array_get(arr, i);
        value = item ? json_object_get(item, "id") : NULL;
        if (value == NULL || !json_is_string(value)) {
            value = item ? json_object_get(item, "uuid") : NULL;
        }
        actual = value && json_is_string(value) ?
            json_string_value(value) : NULL;
        if (actual != NULL && ulab_streq(actual, id)) {
            return item;
        }
    }
    return NULL;
}

static void snapshot_optional_bool(json_t *obj, const char *key,
                                   int *value, int *present) {
    json_t *field;

    *value = 0;
    *present = 0;
    field = obj ? json_object_get(obj, key) : NULL;
    if (field != NULL && json_is_boolean(field)) {
        *value = json_is_true(field);
        *present = 1;
    }
}

static int parse_entity_snapshot(const char *entity,
                                 json_t *obj,
                                 bff_entity_snapshot_t *snapshot,
                                 ulab_error_t *err) {
    json_t *site;
    json_t *status;
    json_t *sim;
    json_t *package;

    if (entity == NULL || obj == NULL || !json_is_object(obj) ||
        snapshot == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "entity snapshot has invalid payload");
        return ULAB_ERR;
    }
    memset(snapshot, 0, sizeof(*snapshot));

    if (json_get_str(obj,
                     ulab_streq(entity, "subscriber") ||
                     ulab_streq(entity, "customer") ||
                     ulab_streq(entity, "package") ||
                     ulab_streq(entity, "plan") ? "uuid" : "id",
                     snapshot->id, sizeof(snapshot->id)) ||
        json_get_str(obj, "name", snapshot->name,
                     sizeof(snapshot->name))) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s snapshot is missing id/name", entity);
        return ULAB_ERR;
    }
    json_get_optional_str(obj, "networkId", snapshot->network_id,
                          sizeof(snapshot->network_id));
    json_get_optional_str(obj, "type", snapshot->type,
                          sizeof(snapshot->type));
    json_get_optional_str(obj, "latitude", snapshot->latitude,
                          sizeof(snapshot->latitude));
    json_get_optional_str(obj, "longitude", snapshot->longitude,
                          sizeof(snapshot->longitude));
    json_get_optional_str(obj, "location", snapshot->location,
                          sizeof(snapshot->location));
    json_get_optional_str(obj, "createdAt", snapshot->created_at,
                          sizeof(snapshot->created_at));
    json_get_optional_str(obj, "installDate", snapshot->install_date,
                          sizeof(snapshot->install_date));
    json_get_optional_str(obj, "email", snapshot->email,
                          sizeof(snapshot->email));
    json_get_optional_str(obj, "phone", snapshot->phone,
                          sizeof(snapshot->phone));
    snapshot_optional_bool(obj, "active", &snapshot->active,
                           &snapshot->has_active);
    snapshot_optional_bool(obj, "isDeactivated", &snapshot->deactivated,
                           &snapshot->has_deactivated);

    site = json_object_get(obj, "site");
    if (site != NULL && json_is_object(site)) {
        json_get_optional_str(site, "siteId", snapshot->site_id,
                              sizeof(snapshot->site_id));
        if (snapshot->network_id[0] == '\0') {
            json_get_optional_str(site, "networkId",
                                  snapshot->network_id,
                                  sizeof(snapshot->network_id));
        }
    }
    status = json_object_get(obj, "status");
    if (status != NULL && json_is_object(status)) {
        json_get_optional_str(status, "connectivity",
                              snapshot->connectivity,
                              sizeof(snapshot->connectivity));
        json_get_optional_str(status, "state", snapshot->state,
                              sizeof(snapshot->state));
    }
    sim = json_object_get(obj, "sim");
    if (sim != NULL && json_is_object(sim)) {
        json_get_optional_str(sim, "id", snapshot->sim_id,
                              sizeof(snapshot->sim_id));
        json_get_optional_str(sim, "status", snapshot->sim_status,
                              sizeof(snapshot->sim_status));
        package = json_object_get(sim, "package");
        if (package != NULL && json_is_object(package)) {
            json_get_optional_str(package, "package_id",
                                  snapshot->package_id,
                                  sizeof(snapshot->package_id));
            snapshot_optional_bool(package, "is_currently_in_use",
                                   &snapshot->active,
                                   &snapshot->has_active);
        }
    }
    return ULAB_OK;
}

static int snapshot_common_equal(const char *entity,
                                 const bff_entity_snapshot_t *left,
                                 const bff_entity_snapshot_t *right) {
    if (!ulab_streq(left->id, right->id) ||
        !ulab_streq(left->name, right->name)) {
        return 0;
    }
    if (ulab_streq(entity, "network")) {
        return 1;
    }
    if (!ulab_streq(left->network_id, right->network_id)) {
        return 0;
    }
    if (ulab_streq(entity, "site")) {
        return ulab_streq(left->latitude, right->latitude) &&
            ulab_streq(left->longitude, right->longitude) &&
            ulab_streq(left->location, right->location) &&
            left->created_at[0] != '\0' &&
            right->created_at[0] != '\0' &&
            ulab_streq(left->install_date, right->install_date) &&
            left->has_deactivated == right->has_deactivated &&
            (!left->has_deactivated ||
             left->deactivated == right->deactivated);
    }
    if (ulab_streq(entity, "node")) {
        return ulab_streq(left->type, right->type) &&
            ulab_streq(left->site_id, right->site_id) &&
            ulab_streq(left->connectivity, right->connectivity) &&
            ulab_streq(left->state, right->state);
    }
    if (ulab_streq(entity, "subscriber") ||
        ulab_streq(entity, "customer")) {
        return ulab_streq(left->email, right->email) &&
            ulab_streq(left->phone, right->phone) &&
            ulab_streq(left->sim_id, right->sim_id) &&
            ulab_streq(left->sim_status, right->sim_status) &&
            ulab_streq(left->package_id, right->package_id) &&
            left->has_active == right->has_active &&
            (!left->has_active || left->active == right->active);
    }
    if (ulab_streq(entity, "package") || ulab_streq(entity, "plan")) {
        return left->has_active == right->has_active &&
            (!left->has_active || left->active == right->active);
    }
    return 0;
}

static int entity_context(const char *entity,
                          const char *ref,
                          const world_t *world,
                          const char **id,
                          const network_t **network,
                          ulab_error_t *err) {
    world_t *mutable_world;

    mutable_world = (world_t *)world;
    *id = NULL;
    *network = NULL;
    if (ulab_streq(entity, "network")) {
        network_t *item;

        item = world_network_by_ref(mutable_world, ref);
        if (item != NULL) {
            *id = item->bff_id;
            *network = item;
        }
    } else if (ulab_streq(entity, "site")) {
        site_t *item;

        item = world_site_by_ref(mutable_world, ref);
        if (item != NULL) {
            *id = item->bff_id;
            *network = world_network_by_ref(mutable_world,
                                             item->network_ref);
        }
    } else if (ulab_streq(entity, "node")) {
        node_t *item;

        item = world_node_by_ref(mutable_world, ref);
        if (item != NULL) {
            *id = item->bff_id;
            *network = world_network_by_ref(mutable_world,
                                             item->network_ref);
        }
    } else if (ulab_streq(entity, "subscriber") ||
               ulab_streq(entity, "customer")) {
        subscriber_t *item;

        item = world_subscriber_by_ref(mutable_world, ref);
        if (item != NULL) {
            *id = item->bff_id;
            *network = world_network_by_ref(mutable_world,
                                             item->network_ref);
        }
    } else if (ulab_streq(entity, "package") ||
               ulab_streq(entity, "plan")) {
        package_t *item;

        item = world_package_by_ref(mutable_world, ref);
        if (item == NULL) {
            item = world_package_by_base_ref(mutable_world, ref);
        }
        if (item != NULL) {
            *id = item->bff_id;
            *network = world_network_by_ref(mutable_world,
                                             item->network_ref);
        }
    } else {
        snprintf(err->msg, sizeof(err->msg),
                 "unsupported entity check: %s", entity);
        return ULAB_ERR;
    }

    if (*id == NULL || (*id)[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "cannot resolve %s ref=%s", entity, ref);
        return ULAB_ERR;
    }
    return ULAB_OK;
}

static int query_entity_detail(bff_client_t *c,
                               const char *entity,
                               const char *id,
                               const char *view,
                               json_t **root,
                               json_t **item,
                               ulab_error_t *err) {
    char vars[ULAB_MAX_QUERY];
    char id_esc[ULAB_MAX_ID * 2];
    const char *operation;
    const char *query;

    if (ulab_streq(entity, "node") &&
        (ulab_streq(view, "nodes_list") ||
         ulab_streq(view, "node_detail") ||
         ulab_streq(view, "site_detail"))) {
        json_t *node_view;
        json_t *section;

        ulab_json_escape(id, id_esc, sizeof(id_esc));
        snprintf(vars, sizeof(vars), "{\"nodeId\":\"%s\"}", id_esc);
        *root = NULL;
        if (bff_call(c, "NodeDetail", BFF_CONSOLE_NODE_DETAIL,
                     vars, root, err)) {
            return ULAB_ERR;
        }
        node_view = dig(*root, "data", "nodeView");
        section = node_view ? json_object_get(node_view, "node") : NULL;
        if (console_section_ok(section, "NodeDetail", "node", err)) {
            json_decref(*root);
            *root = NULL;
            return ULAB_ERR;
        }
        *item = json_object_get(section, "node");
        if (*item == NULL || !json_is_object(*item)) {
            snprintf(err->msg, sizeof(err->msg),
                     "NodeDetail node section missing node");
            json_decref(*root);
            *root = NULL;
            return ULAB_ERR;
        }
        return ULAB_OK;
    }

    if (ulab_streq(entity, "site") &&
        (ulab_streq(view, "sites_list") ||
         ulab_streq(view, "site_detail"))) {
        json_t *site_view;
        json_t *section;

        if (console_site_detail_view(c, id, root, &site_view, err)) {
            return ULAB_ERR;
        }
        section = json_object_get(site_view, "site");
        if (console_section_ok(section, "NetworkSiteDetail",
                               "site", err)) {
            json_decref(*root);
            *root = NULL;
            return ULAB_ERR;
        }
        *item = json_object_get(section, "site");
        if (*item == NULL || !json_is_object(*item)) {
            snprintf(err->msg, sizeof(err->msg),
                     "NetworkSiteDetail site section missing site");
            json_decref(*root);
            *root = NULL;
            return ULAB_ERR;
        }
        return ULAB_OK;
    }

    ulab_json_escape(id, id_esc, sizeof(id_esc));
    operation = NULL;
    query = NULL;
    if (ulab_streq(entity, "network")) {
        snprintf(vars, sizeof(vars), "{\"networkId\":\"%s\"}",
                 id_esc);
        operation = "getNetwork";
        query = BFF_GET_NETWORK;
    } else if (ulab_streq(entity, "site")) {
        snprintf(vars, sizeof(vars), "{\"siteId\":\"%s\"}", id_esc);
        operation = "getSite";
        query = BFF_GET_SITE;
    } else if (ulab_streq(entity, "node")) {
        snprintf(vars, sizeof(vars), "{\"data\":{\"id\":\"%s\"}}",
                 id_esc);
        operation = "getNode";
        query = BFF_GET_NODE;
    } else if (ulab_streq(entity, "subscriber") ||
               ulab_streq(entity, "customer")) {
        snprintf(vars, sizeof(vars),
                 "{\"subscriberId\":\"%s\"}", id_esc);
        operation = "getSubscriber";
        query = BFF_GET_SUBSCRIBER;
    } else if (ulab_streq(entity, "package") ||
               ulab_streq(entity, "plan")) {
        snprintf(vars, sizeof(vars), "{\"packageId\":\"%s\"}",
                 id_esc);
        operation = "getPackage";
        query = BFF_GET_PACKAGE;
    }
    if (operation == NULL || query == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "unsupported entity detail: %s", entity);
        return ULAB_ERR;
    }

    *root = NULL;
    if (bff_call(c, operation, query, vars, root, err)) {
        return ULAB_ERR;
    }
    *item = dig(*root, "data", operation);
    if (*item == NULL || !json_is_object(*item)) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s missing entity detail", operation);
        json_decref(*root);
        *root = NULL;
        return ULAB_ERR;
    }
    return ULAB_OK;
}

static int query_entity_list_item(bff_client_t *c,
                                  const char *entity,
                                  const char *id,
                                  const network_t *network,
                                  const char *view,
                                  json_t **root,
                                  json_t **item,
                                  ulab_error_t *err) {
    const char *target;
    const char *envelope;
    const char *list_key;
    json_t *obj;
    json_t *arr;

    if (ulab_streq(entity, "node") &&
        ulab_streq(view, "nodes_list")) {
        if (console_nodes_array(c, network, view, root, &arr, err)) {
            return ULAB_ERR;
        }
        *item = console_find_id(arr, id);
        if (*item == NULL) {
            snprintf(err->msg, sizeof(err->msg),
                     "node id=%s is absent from NodesList", id);
            json_decref(*root);
            *root = NULL;
            return ULAB_ERR;
        }
        return ULAB_OK;
    }

    if (ulab_streq(entity, "site") &&
        (ulab_streq(view, "sites_list") ||
         ulab_streq(view, "site_detail"))) {
        if (console_sites_array(c, network, root, &arr, err)) {
            return ULAB_ERR;
        }
        *item = console_find_id(arr, id);
        if (*item == NULL) {
            snprintf(err->msg, sizeof(err->msg),
                     "site id=%s is absent from SitesList", id);
            json_decref(*root);
            *root = NULL;
            return ULAB_ERR;
        }
        return ULAB_OK;
    }

    if (ulab_streq(entity, "network")) target = "networks";
    else if (ulab_streq(entity, "site")) target = "sites";
    else if (ulab_streq(entity, "node")) target = "nodes";
    else if (ulab_streq(entity, "subscriber") ||
             ulab_streq(entity, "customer")) target = "subscribers";
    else if (ulab_streq(entity, "package") ||
             ulab_streq(entity, "plan")) target = "packages";
    else {
        snprintf(err->msg, sizeof(err->msg),
                 "unsupported entity list: %s", entity);
        return ULAB_ERR;
    }

    envelope = NULL;
    list_key = NULL;
    if (direct_network_list_call(c, target, network, root,
                                 &envelope, &list_key, err)) {
        return ULAB_ERR;
    }
    obj = dig(*root, "data", envelope);
    arr = obj ? json_object_get(obj, list_key) : NULL;
    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s missing %s list", envelope, list_key);
        json_decref(*root);
        *root = NULL;
        return ULAB_ERR;
    }
    *item = json_array_find_id(arr, id);
    if (*item == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s id=%s is absent from %s", entity, id, envelope);
        json_decref(*root);
        *root = NULL;
        return ULAB_ERR;
    }
    return ULAB_OK;
}

int bff_entity_list_detail_reconciles(bff_client_t *c,
                                      const char *entity,
                                      const char *ref,
                                      const world_t *world,
                                      const char *view,
                                      int *matched,
                                      char *detail,
                                      size_t detail_len,
                                      ulab_error_t *err) {
    const char *id;
    const network_t *network;
    json_t *list_root;
    json_t *detail_root;
    json_t *list_item;
    json_t *detail_item;
    bff_entity_snapshot_t list_snapshot;
    bff_entity_snapshot_t detail_snapshot;

    if (entity == NULL || ref == NULL || world == NULL ||
        matched == NULL || detail == NULL || detail_len == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "entity reconciliation has invalid arguments");
        return ULAB_ERR;
    }
    if (entity_context(entity, ref, world, &id, &network, err)) {
        return ULAB_ERR;
    }
    list_root = NULL;
    detail_root = NULL;
    if (ulab_streq(entity, "node") && ulab_streq(view, "site_detail")) {
        world_t *mutable_world;
        node_t *expected_node;
        site_t *expected_site;

        mutable_world = (world_t *)world;
        expected_node = world_node_by_ref(mutable_world, ref);
        expected_site = expected_node ?
            world_site_by_ref(mutable_world, expected_node->site_ref) : NULL;
        if (expected_node == NULL || expected_site == NULL ||
            expected_site->bff_id[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "site_detail node reconciliation has unresolved site");
            return ULAB_ERR;
        }
        if (console_site_detail_node_item(c, expected_site->bff_id, id,
                                          &list_root, &list_item, err)) {
            return ULAB_ERR;
        }
    } else if (query_entity_list_item(c, entity, id, network, view,
                                      &list_root, &list_item, err)) {
        return ULAB_ERR;
    }
    if (parse_entity_snapshot(entity, list_item, &list_snapshot, err)) {
        json_decref(list_root);
        return ULAB_ERR;
    }
    if (query_entity_detail(c, entity, id, view, &detail_root,
                            &detail_item, err)) {
        json_decref(list_root);
        return ULAB_ERR;
    }
    if (parse_entity_snapshot(entity, detail_item, &detail_snapshot, err)) {
        json_decref(list_root);
        json_decref(detail_root);
        return ULAB_ERR;
    }

    *matched = snapshot_common_equal(entity, &list_snapshot,
                                     &detail_snapshot);
    snprintf(detail, detail_len,
             "entity=%s ref=%s id=%s list_name=%.96s detail_name=%.96s "
             "network=%.64s/%.64s state=%.32s/%.32s connectivity=%.32s/%.32s",
             entity, ref, id, list_snapshot.name, detail_snapshot.name,
             list_snapshot.network_id, detail_snapshot.network_id,
             list_snapshot.state, detail_snapshot.state,
             list_snapshot.connectivity, detail_snapshot.connectivity);
    json_decref(list_root);
    json_decref(detail_root);
    return ULAB_OK;
}

static const char *world_node_kind(const node_t *node) {
    if (node == NULL) {
        return "";
    }
    if (ulab_streq(node->type, ULAB_NODE_TOWER)) {
        return ULAB_NODE_KIND_TOWER;
    }
    if (ulab_streq(node->type, ULAB_NODE_AMPLIFIER)) {
        return ULAB_NODE_KIND_AMPLIFIER;
    }
    if (ulab_streq(node->type, ULAB_NODE_CONTROLLER)) {
        return ULAB_NODE_KIND_CONTROLLER;
    }
    return node->type;
}

int bff_entity_fields_match_world(bff_client_t *c,
                                  const char *entity,
                                  const char *ref,
                                  const world_t *world,
                                  const char *view,
                                  int *matched,
                                  char *detail,
                                  size_t detail_len,
                                  ulab_error_t *err) {
    const char *id;
    const network_t *network;
    json_t *root;
    json_t *item;
    bff_entity_snapshot_t actual;
    world_t *mutable_world;

    if (entity == NULL || ref == NULL || world == NULL ||
        matched == NULL || detail == NULL || detail_len == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "entity field check has invalid arguments");
        return ULAB_ERR;
    }
    if (entity_context(entity, ref, world, &id, &network, err)) {
        return ULAB_ERR;
    }
    root = NULL;
    if (ulab_streq(entity, "node") && ulab_streq(view, "site_detail")) {
        node_t *expected_node;
        site_t *expected_site;

        mutable_world = (world_t *)world;
        expected_node = world_node_by_ref(mutable_world, ref);
        expected_site = expected_node ?
            world_site_by_ref(mutable_world, expected_node->site_ref) : NULL;
        if (expected_node == NULL || expected_site == NULL ||
            expected_site->bff_id[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "site_detail node field check has unresolved site");
            return ULAB_ERR;
        }
        if (console_site_detail_node_item(c, expected_site->bff_id, id,
                                          &root, &item, err)) {
            return ULAB_ERR;
        }
    } else if (ulab_streq(entity, "node") &&
               ulab_streq(view, "nodes_list")) {
        if (query_entity_list_item(c, entity, id, network, view,
                                   &root, &item, err)) {
            return ULAB_ERR;
        }
    } else {
        if (query_entity_detail(c, entity, id, view,
                                &root, &item, err)) {
            return ULAB_ERR;
        }
    }
    if (parse_entity_snapshot(entity, item, &actual, err)) {
        json_decref(root);
        return ULAB_ERR;
    }

    mutable_world = (world_t *)world;
    *matched = 0;
    if (ulab_streq(entity, "network")) {
        network_t *expected;

        expected = world_network_by_ref(mutable_world, ref);
        *matched = expected != NULL &&
            ulab_streq(actual.id, expected->bff_id) &&
            ulab_streq(actual.name, expected->name);
    } else if (ulab_streq(entity, "site")) {
        site_t *expected;

        expected = world_site_by_ref(mutable_world, ref);
        *matched = expected != NULL && network != NULL &&
            ulab_streq(actual.id, expected->bff_id) &&
            ulab_streq(actual.name, expected->name) &&
            ulab_streq(actual.network_id, network->bff_id) &&
            ulab_streq(actual.latitude, expected->latitude) &&
            ulab_streq(actual.longitude, expected->longitude) &&
            ulab_streq(actual.location, expected->location) &&
            expected->created_at[0] != '\0' &&
            actual.created_at[0] != '\0' &&
            expected->install_date[0] != '\0' &&
            ulab_streq(actual.install_date, expected->install_date);
    } else if (ulab_streq(entity, "node")) {
        node_t *expected;
        site_t *site;

        expected = world_node_by_ref(mutable_world, ref);
        site = expected ? world_site_by_ref(mutable_world,
                                             expected->site_ref) : NULL;
        *matched = expected != NULL && network != NULL && site != NULL &&
            ulab_streq(actual.id, expected->bff_id) &&
            ulab_streq(actual.name, expected->name) &&
            ulab_streq(actual.type, world_node_kind(expected)) &&
            ulab_streq(actual.site_id, site->bff_id) &&
            ulab_streq(actual.network_id, network->bff_id);
    } else if (ulab_streq(entity, "subscriber") ||
               ulab_streq(entity, "customer")) {
        subscriber_t *expected;

        expected = world_subscriber_by_ref(mutable_world, ref);
        *matched = expected != NULL && network != NULL &&
            ulab_streq(actual.id, expected->bff_id) &&
            ulab_streq(actual.name, expected->name) &&
            ulab_streq(actual.email, expected->email) &&
            ulab_streq(actual.phone, expected->phone) &&
            ulab_streq(actual.network_id, network->bff_id);
    } else if (ulab_streq(entity, "package") ||
               ulab_streq(entity, "plan")) {
        package_t *expected;
        bff_package_t package_actual;
        int duration;

        expected = world_package_by_ref(mutable_world, ref);
        if (expected == NULL) {
            expected = world_package_by_base_ref(mutable_world, ref);
        }
        memset(&package_actual, 0, sizeof(package_actual));
        if (expected != NULL && bff_get_package(c, expected,
                                                &package_actual, err)) {
            json_decref(root);
            return ULAB_ERR;
        }
        duration = expected && expected->duration_minutes > 0 ?
            (int)expected->duration_minutes :
            (expected ? (int)(expected->duration_days * 1440u) : 0);
        *matched = expected != NULL &&
            ulab_streq(package_actual.uuid, expected->bff_id) &&
            ulab_streq(package_actual.name, expected->name) &&
            package_actual.data_volume == expected->data_mb &&
            package_actual.duration_minutes == (uint32_t)duration &&
            fabs(package_actual.amount - expected->amount) <= 0.000001 &&
            ulab_streq(package_actual.currency, expected->currency) &&
            ulab_streq(package_actual.country, expected->country) &&
            package_actual.active == expected->active &&
            ((expected->network_ref[0] == '\0' &&
              package_actual.network_id[0] == '\0') ||
             (expected->network_ref[0] != '\0' && network != NULL &&
              ulab_streq(package_actual.network_id, network->bff_id)));
    }

    snprintf(detail, detail_len,
             "entity=%s ref=%s id=%s name=%.96s network=%.64s "
             "site=%.64s type=%.24s",
             entity, ref, actual.id, actual.name, actual.network_id,
             actual.site_id, actual.type);
    json_decref(root);
    return ULAB_OK;
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
        if ((ulab_streq(op, "unsetPackageInUseForSim") ||
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

        if (ulab_streq(op, "toggleSimServiceStatus") &&
            strstr(err.msg, "is invalid for turning off") != NULL) {
            if (c != NULL && c->logf != NULL) {
                fprintf(c->logf,
                        "cleanup ignore: %s: sim already off\n",
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
    char package_record_ids[32][ULAB_MAX_ID];
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
            pid = it ? json_object_get(it, "id") : NULL;
            if (pid == NULL || !json_is_string(pid) ||
                json_string_value(pid) == NULL ||
                json_string_value(pid)[0] == '\0') {
                pid = it ? json_object_get(it, "package_id") : NULL;
            }
            if (pid != NULL && json_is_string(pid) &&
                json_string_value(pid) != NULL &&
                json_string_value(pid)[0] != '\0') {
                ulab_copy(package_record_ids[count],
                          sizeof(package_record_ids[count]),
                          json_string_value(pid));
                count++;
            }
        }
    }

    json_decref(root);

    for (i = 0; i < count; i++) {

        n = snprintf(query, sizeof(query),
                     "mutation { unsetPackageInUseForSim(data: {"
                     "packageId: \"%s\", simId: \"%s\"}) { packageId } }",
                     package_record_ids[i], ue->bff_id);
        if (n >= 0 && (size_t)n < sizeof(query) &&
            bff_cleanup_call(c, "unsetPackageInUseForSim", query)) {
            (*failures)++;
        }

        n = snprintf(query, sizeof(query),
                     "mutation { removePackageForSim(data: {"
                     "packageId: \"%s\", simId: \"%s\"}) { packageId } }",
                     package_record_ids[i], ue->bff_id);
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
                             "mutation { unsetPackageInUseForSim(data: {"
                             "packageId: \"%s\", simId: \"%s\"}) { packageId } }",
                             pkg_id, ue->bff_id);
                if (n >= 0 && (size_t)n < sizeof(query) &&
                    bff_cleanup_call(c, "unsetPackageInUseForSim", query)) {
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
                     "mutation { toggleSimServiceStatus(data: {sim_id: \"%s\", "
                     "status: \"service_off\"}) { success } }",
                     ue->bff_id);
        if (n >= 0 && (size_t)n < sizeof(query) &&
            bff_cleanup_call(c, "toggleSimServiceStatus", query)) {
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
