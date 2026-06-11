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
#include "log.h"
#include "util.h"

extern const char *BFF_ADD_NETWORK;
extern const char *BFF_ADD_SITE;
extern const char *BFF_ADD_NODE;
extern const char *BFF_ADD_NODE_TO_SITE;
extern const char *BFF_ADD_PACKAGE;
extern const char *BFF_ADD_SUBSCRIBER;
extern const char *BFF_ALLOCATE_SIM;
extern const char *BFF_GET_DATA_USAGE;
extern const char *BFF_GET_SIM_PACKAGES;
extern const char *BFF_GET_NODE_STATE;
extern const char *BFF_NETWORK_OVERVIEW;
extern const char *BFF_SITE_VIEW;
extern const char *BFF_GET_NETWORKS;
extern const char *BFF_GET_SITES;
extern const char *BFF_GET_NODES_FOR_SITE;

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

    curl_easy_setopt(curl, CURLOPT_URL, c->url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, hdr);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &resp);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 30L);

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

int bff_init(bff_client_t *c, const char *url, const char *run_dir) {
    char path[ULAB_MAX_PATH];

    memset(c, 0, sizeof(*c));
    ulab_copy(c->url, sizeof(c->url), url);

    snprintf(path, sizeof(path), "%s/bff.log", run_dir);
    c->logf = fopen(path, "w");

    curl_global_init(CURL_GLOBAL_ALL);

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
             "{\"data\":{\"name\":\"%s\",\"budget\":0,"
             "\"countries\":[\"USA\"],\"networks\":[\"A3\"]}}",
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

int bff_add_site(bff_client_t *c, site_t *s, const network_t *n,
                 ulab_error_t *err) {
    char vars[4096];
    json_t *root;
    json_t *obj;

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"name\":\"%s\",\"network_id\":\"%s\","
             "\"latitude\":\"37.7749\",\"longitude\":\"-122.4194\","
             "\"location\":\"Lab\",\"access_id\":\"\","
             "\"backhaul_id\":\"\",\"install_date\":\"\","
             "\"power_id\":\"\",\"spectrum_id\":\"\","
             "\"switch_id\":\"\"}}", s->name, n->bff_id);

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

int bff_add_node(bff_client_t *c, node_t *n, ulab_error_t *err) {
    char vars[4096];
    json_t *root;
    json_t *obj;

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"id\":\"%s\",\"name\":\"%s\","
             "\"type\":\"%s\"}}", n->id, n->name, n->type);

    if (bff_call(c, "addNode", BFF_ADD_NODE, vars, &root, err)) {
        return ULAB_ERR;
    }

    obj = dig(root, "data", "addNode");
    if (obj == NULL || json_get_str(obj, "id", n->bff_id,
        sizeof(n->bff_id))) {
        snprintf(err->msg, sizeof(err->msg), "addNode missing id");
        json_decref(root);
        return ULAB_ERR;
    }

    json_decref(root);

    return ULAB_OK;
}

int bff_add_node_to_site(bff_client_t *c, const node_t *n,
                         const site_t *s, const network_t *net,
                         ulab_error_t *err) {
    char vars[4096];
    json_t *root;

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"networkId\":\"%s\",\"nodeId\":\"%s\","
             "\"siteId\":\"%s\"}}", net->bff_id, n->bff_id, s->bff_id);

    if (bff_call(c, "addNodeToSite", BFF_ADD_NODE_TO_SITE, vars, &root,
        err)) {
        return ULAB_ERR;
    }

    json_decref(root);

    return ULAB_OK;
}

int bff_add_package(bff_client_t *c, package_t *p, ulab_error_t *err) {
    char vars[4096];
    json_t *root;
    json_t *obj;

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"name\":\"%s\",\"amount\":%.2f,"
             "\"dataUnit\":\"MB\",\"dataVolume\":%llu,"
             "\"duration\":%u,\"currency\":\"USD\","
             "\"country\":\"USA\"}}", p->name, p->amount,
             (unsigned long long)p->data_mb, p->duration_hours);

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
             "{\"data\":{\"email\":\"%s\",\"first_name\":\"Lab\","
             "\"last_name\":\"User\",\"network_id\":\"%s\","
             "\"phone\":\"%s\"}}", sub->email, net->bff_id, sub->phone);

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
                     ulab_error_t *err) {
    char vars[4096];
    json_t *root;
    json_t *obj;

    snprintf(vars, sizeof(vars),
             "{\"data\":{\"iccid\":\"%s\",\"network_id\":\"%s\","
             "\"sim_type\":\"test\",\"package_id\":\"%s\","
             "\"subscriber_id\":\"%s\",\"traffic_policy\":1}}",
             ue->iccid, net->bff_id, pkg->bff_id, sub->bff_id);

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

int bff_get_sim_usage(bff_client_t *c, const ue_t *ue, uint64_t *used_mb,
                      ulab_error_t *err) {
    char vars[4096];
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
    char vars[4096];
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
    char vars[4096];
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
    char vars[4096];
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
    char vars[4096];
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
