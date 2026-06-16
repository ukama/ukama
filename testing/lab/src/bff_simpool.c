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

#define SIMPOOL_HTTP_TIMEOUT_SEC 30L

typedef struct {
    char *buf;
    size_t len;
} sim_http_buf_t;

static size_t sim_write_cb(void *ptr, size_t size, size_t nmemb,
                           void *data) {
    sim_http_buf_t *b;
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

static json_t *sim_dig(json_t *root, const char *a, const char *b) {
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

static int sim_json_get_str(json_t *obj,
                            const char *key,
                            char *out,
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

static char *sim_read_file_all(const char *path,
                               size_t *len,
                               ulab_error_t *err) {
    FILE *f;
    long n;
    char *buf;
    size_t got;

    f = fopen(path, "rb");
    if (f == NULL) {
        snprintf(err->msg, sizeof(err->msg), "failed to open %s", path);
        return NULL;
    }

    if (fseek(f, 0, SEEK_END) != 0) {
        snprintf(err->msg, sizeof(err->msg), "failed to seek %s", path);
        fclose(f);
        return NULL;
    }

    n = ftell(f);
    if (n < 0) {
        snprintf(err->msg, sizeof(err->msg), "failed to size %s", path);
        fclose(f);
        return NULL;
    }

    if (fseek(f, 0, SEEK_SET) != 0) {
        snprintf(err->msg, sizeof(err->msg), "failed to rewind %s", path);
        fclose(f);
        return NULL;
    }

    buf = calloc(1, (size_t)n + 1);
    if (buf == NULL) {
        snprintf(err->msg, sizeof(err->msg), "out of memory reading %s",
                 path);
        fclose(f);
        return NULL;
    }

    got = fread(buf, 1, (size_t)n, f);
    fclose(f);

    if (got != (size_t)n) {
        snprintf(err->msg, sizeof(err->msg), "failed to read %s", path);
        free(buf);
        return NULL;
    }

    buf[got] = '\0';
    *len = got;

    return buf;
}

static char *sim_base64_encode(const unsigned char *src,
                               size_t len,
                               ulab_error_t *err) {
    static const char tbl[] =
        "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    char *out;
    size_t out_len;
    size_t i;
    size_t j;
    unsigned int v;
    unsigned int a;
    unsigned int b;
    unsigned int c;
    int have_b;
    int have_c;

    out_len = ((len + 2) / 3) * 4;
    out = calloc(1, out_len + 1);
    if (out == NULL) {
        snprintf(err->msg, sizeof(err->msg), "out of memory base64 encode");
        return NULL;
    }

    i = 0;
    j = 0;
    while (i < len) {
        a = src[i++];
        have_b = i < len;
        b = have_b ? src[i++] : 0;
        have_c = i < len;
        c = have_c ? src[i++] : 0;

        v = (a << 16) | (b << 8) | c;
        out[j++] = tbl[(v >> 18) & 0x3f];
        out[j++] = tbl[(v >> 12) & 0x3f];
        out[j++] = have_b ? tbl[(v >> 6) & 0x3f] : '=';
        out[j++] = have_c ? tbl[v & 0x3f] : '=';
    }

    out[out_len] = '\0';

    return out;
}

static char *sim_make_graphql_body(const char *query,
                                   ulab_error_t *err) {
    char *qesc;
    char *body;
    size_t q_len;
    size_t qesc_len;
    size_t body_len;
    int n;

    q_len = strlen(query);
    qesc_len = (q_len * 2) + 1;
    body_len = qesc_len + 64;

    qesc = calloc(1, qesc_len);
    body = calloc(1, body_len);
    if (qesc == NULL || body == NULL) {
        snprintf(err->msg, sizeof(err->msg), "out of memory bff body");
        free(qesc);
        free(body);
        return NULL;
    }

    ulab_json_escape(query, qesc, qesc_len);
    n = snprintf(body, body_len, "{\"query\":\"%s\",\"variables\":{}}",
                 qesc);
    free(qesc);

    if (n < 0 || (size_t)n >= body_len) {
        snprintf(err->msg, sizeof(err->msg), "bff request body too long");
        free(body);
        return NULL;
    }

    return body;
}

static int sim_graphql_call(bff_client_t *c,
                            const char *op,
                            const char *query,
                            json_t **out,
                            ulab_error_t *err) {
    CURL *curl;
    CURLcode ret;
    struct curl_slist *hdr;
    sim_http_buf_t resp;
    char token_hdr[8192];
    char *body;
    long code;
    json_t *root;
    json_t *errors;
    json_error_t json_err;

    hdr = NULL;
    body = NULL;
    root = NULL;
    errors = NULL;
    code = 0;
    resp.buf = NULL;
    resp.len = 0;

    body = sim_make_graphql_body(query, err);
    if (body == NULL) {
        return ULAB_ERR;
    }

    if (c->logf) {
        fprintf(c->logf, "--- %s request ---\n%s\n", op, body);
        fflush(c->logf);
    }

    curl = curl_easy_init();
    if (curl == NULL) {
        snprintf(err->msg, sizeof(err->msg), "%s: curl init failed", op);
        free(body);
        return ULAB_ERR;
    }

    hdr = curl_slist_append(hdr, "Content-Type: application/json");

    if (c->authenticated) {
        snprintf(token_hdr, sizeof(token_hdr),
                 "X-Session-Token: %s", c->token);
        hdr = curl_slist_append(hdr, token_hdr);
    }

    curl_easy_setopt(curl, CURLOPT_URL, c->url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, hdr);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, sim_write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &resp);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, SIMPOOL_HTTP_TIMEOUT_SEC);

    ret = curl_easy_perform(curl);
    if (ret != CURLE_OK) {
        snprintf(err->msg, sizeof(err->msg), "%s: HTTP request failed: %s",
                 op, curl_easy_strerror(ret));
        curl_slist_free_all(hdr);
        curl_easy_cleanup(curl);
        free(resp.buf);
        free(body);
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

    free(body);

    if (code < 200 || code >= 300) {
        snprintf(err->msg, sizeof(err->msg), "%s: HTTP %ld", op, code);
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

static int sim_copy_pool_iccid(json_t *item,
                               char *out,
                               size_t out_len) {
    json_t *v;
    const char *s;

    if (item == NULL) {
        return ULAB_ERR;
    }

    if (json_is_string(item)) {
        s = json_string_value(item);
    } else if (json_is_object(item)) {
        v = json_object_get(item, "iccid");
        if (v == NULL || !json_is_string(v)) {
            return ULAB_ERR;
        }
        s = json_string_value(v);
    } else {
        return ULAB_ERR;
    }

    if (s == NULL || s[0] == '\0') {
        return ULAB_ERR;
    }

    return ulab_copy(out, out_len, s);
}

int bff_upload_sims_from_csv(bff_client_t *c,
                             const char *csv_path,
                             const char *sim_type,
                             ulab_error_t *err) {
    char *csv;
    char *b64;
    char *b64_esc;
    char *type_esc;
    char *query;
    size_t csv_len;
    size_t b64_len;
    size_t query_len;
    json_t *root;
    int n;

    csv = NULL;
    b64 = NULL;
    b64_esc = NULL;
    type_esc = NULL;
    query = NULL;
    root = NULL;

    if (csv_path == NULL || csv_path[0] == '\0') {
        return ULAB_OK;
    }

    if (sim_type == NULL || sim_type[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "missing SIM type");
        return ULAB_ERR;
    }

    csv = sim_read_file_all(csv_path, &csv_len, err);
    if (csv == NULL) {
        return ULAB_ERR;
    }

    b64 = sim_base64_encode((const unsigned char *)csv, csv_len, err);
    free(csv);
    if (b64 == NULL) {
        return ULAB_ERR;
    }

    b64_len = strlen(b64);
    b64_esc = calloc(1, (b64_len * 2) + 1);
    type_esc = calloc(1, (strlen(sim_type) * 2) + 1);
    if (b64_esc == NULL || type_esc == NULL) {
        snprintf(err->msg, sizeof(err->msg), "out of memory upload sims");
        free(b64);
        free(b64_esc);
        free(type_esc);
        return ULAB_ERR;
    }

    ulab_json_escape(b64, b64_esc, (b64_len * 2) + 1);
    ulab_json_escape(sim_type, type_esc, (strlen(sim_type) * 2) + 1);
    free(b64);

    query_len = strlen(b64_esc) + strlen(type_esc) + 256;
    query = calloc(1, query_len);
    if (query == NULL) {
        snprintf(err->msg, sizeof(err->msg), "out of memory upload query");
        free(b64_esc);
        free(type_esc);
        return ULAB_ERR;
    }

    n = snprintf(query, query_len,
                 "mutation { uploadsSims(data:{data:\"%s\","
                 "simType:\"%s\"}) { iccid } }",
                 b64_esc, type_esc);
    free(b64_esc);
    free(type_esc);

    if (n < 0 || (size_t)n >= query_len) {
        snprintf(err->msg, sizeof(err->msg), "upload sims query too long");
        free(query);
        return ULAB_ERR;
    }

    if (sim_graphql_call(c, "uploadsSims", query, &root, err)) {
        free(query);
        return ULAB_ERR;
    }

    free(query);
    json_decref(root);

    return ULAB_OK;
}

int bff_get_sims_from_pool(bff_client_t *c,
                           const char *sim_type,
                           char iccids[][ULAB_MAX_ID],
                           size_t max_iccids,
                           size_t *iccid_count,
                           ulab_error_t *err) {
    char type_esc[ULAB_MAX_REF * 2];
    char query[ULAB_MAX_QUERY];
    json_t *root;
    json_t *obj;
    json_t *arr;
    json_t *it;
    size_t i;
    size_t n;

    root = NULL;
    obj = NULL;
    arr = NULL;
    n = 0;

    if (iccid_count == NULL) {
        snprintf(err->msg, sizeof(err->msg), "invalid SIM pool count arg");
        return ULAB_ERR;
    }

    *iccid_count = 0;

    if (sim_type == NULL || sim_type[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "missing SIM type");
        return ULAB_ERR;
    }

    ulab_json_escape(sim_type, type_esc, sizeof(type_esc));
    snprintf(query, sizeof(query),
             "query { getSimsFromPool(data:{type:\"%s\","
             "status:UNASSIGNED}) { sims { iccid } } }",
             type_esc);

    if (sim_graphql_call(c, "getSimsFromPool", query, &root, err)) {
        return ULAB_ERR;
    }

    obj = sim_dig(root, "data", "getSimsFromPool");
    if (obj != NULL && json_is_object(obj)) {
        arr = json_object_get(obj, "sims");
    } else if (obj != NULL && json_is_array(obj)) {
        arr = obj;
    }

    if (arr == NULL || !json_is_array(arr)) {
        snprintf(err->msg, sizeof(err->msg),
                 "getSimsFromPool missing sims list");
        json_decref(root);
        return ULAB_ERR;
    }

    for (i = 0; i < json_array_size(arr) && n < max_iccids; i++) {
        it = json_array_get(arr, i);
        if (sim_copy_pool_iccid(it, iccids[n], ULAB_MAX_ID) == ULAB_OK) {
            n++;
        }
    }

    *iccid_count = n;
    json_decref(root);

    return ULAB_OK;
}

int bff_allocate_sim_from_pool(bff_client_t *c,
                               ue_t *ue,
                               const subscriber_t *sub,
                               const network_t *net,
                               const package_t *pkg,
                               const char *sim_type,
                               ulab_error_t *err) {
    char iccid_esc[ULAB_MAX_ID * 2];
    char net_esc[ULAB_MAX_ID * 2];
    char type_esc[ULAB_MAX_REF * 2];
    char pkg_esc[ULAB_MAX_ID * 2];
    char sub_esc[ULAB_MAX_ID * 2];
    char query[ULAB_MAX_QUERY * 2];
    json_t *root;
    json_t *obj;

    root = NULL;

    if (ue == NULL || ue->iccid[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "allocateSim missing ICCID");
        return ULAB_ERR;
    }

    if (sub == NULL || sub->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "allocateSim missing subscriber id");
        return ULAB_ERR;
    }

    if (net == NULL || net->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "allocateSim missing network id");
        return ULAB_ERR;
    }

    if (pkg == NULL || pkg->bff_id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "allocateSim missing package id");
        return ULAB_ERR;
    }

    if (sim_type == NULL || sim_type[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "allocateSim missing SIM type");
        return ULAB_ERR;
    }

    ulab_json_escape(ue->iccid, iccid_esc, sizeof(iccid_esc));
    ulab_json_escape(net->bff_id, net_esc, sizeof(net_esc));
    ulab_json_escape(sim_type, type_esc, sizeof(type_esc));
    ulab_json_escape(pkg->bff_id, pkg_esc, sizeof(pkg_esc));
    ulab_json_escape(sub->bff_id, sub_esc, sizeof(sub_esc));

    snprintf(query, sizeof(query),
             "mutation { allocateSim(data:{iccid:\"%s\","
             "network_id:\"%s\",sim_type:\"%s\","
             "package_id:\"%s\",subscriber_id:\"%s\","
             "traffic_policy:0}) { id subscriber_id network_id iccid imsi "
             "status package { packageId isActive startDate endDate } } }",
             iccid_esc, net_esc, type_esc, pkg_esc, sub_esc);

    if (sim_graphql_call(c, "allocateSim", query, &root, err)) {
        return ULAB_ERR;
    }

    obj = sim_dig(root, "data", "allocateSim");
    if (obj == NULL ||
        sim_json_get_str(obj, "id", ue->bff_id, sizeof(ue->bff_id))) {
        snprintf(err->msg, sizeof(err->msg), "allocateSim missing id");
        json_decref(root);
        return ULAB_ERR;
    }

    sim_json_get_str(obj, "iccid", ue->iccid, sizeof(ue->iccid));
    sim_json_get_str(obj, "imsi", ue->imsi, sizeof(ue->imsi));

    json_decref(root);

    return ULAB_OK;
}
