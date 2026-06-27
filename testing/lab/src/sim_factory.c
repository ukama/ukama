/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#include "sim_factory.h"
#include "log.h"
#include "util.h"

#define SIMFACTORY_HTTP_TIMEOUT_SEC 30L
#define SIMFACTORY_WAIT_ATTEMPTS    60
#define SIMFACTORY_WAIT_SLEEP_SEC   2
#define ASR_WAIT_ATTEMPTS           60
#define ASR_WAIT_SLEEP_SEC          2

typedef struct {
    char *buf;
    size_t len;
} sf_http_buf_t;

static size_t sf_write_cb(void *ptr, size_t size, size_t nmemb, void *data) {
    sf_http_buf_t *b;
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

static int sf_join_url(char *out,
                       size_t out_len,
                       const char *base,
                       const char *path,
                       ulab_error_t *err) {
    size_t n;
    int rc;

    if (base == NULL || base[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "missing URL base");
        return ULAB_ERR;
    }

    n = strlen(base);
    while (n > 0 && base[n - 1] == '/') {
        n--;
    }

    rc = snprintf(out, out_len, "%.*s%s", (int)n, base, path);
    if (rc < 0 || (size_t)rc >= out_len) {
        snprintf(err->msg, sizeof(err->msg), "URL too long");
        return ULAB_ERR;
    }

    return ULAB_OK;
}

static int sf_http_request(const char *op,
                           const char *method,
                           const char *url,
                           const char *body,
                           char **out,
                           long *status,
                           ulab_error_t *err) {
    CURL *curl;
    CURLcode ret;
    struct curl_slist *hdr;
    sf_http_buf_t resp;
    long code;

    hdr = NULL;
    resp.buf = NULL;
    resp.len = 0;
    code = 0;

    if (out != NULL) {
        *out = NULL;
    }
    if (status != NULL) {
        *status = 0;
    }

    curl = curl_easy_init();
    if (curl == NULL) {
        snprintf(err->msg, sizeof(err->msg), "%s: curl init failed", op);
        return ULAB_ERR;
    }

    hdr = curl_slist_append(hdr, "Content-Type: application/json");
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, hdr);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, sf_write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &resp);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, SIMFACTORY_HTTP_TIMEOUT_SEC);

    if (ulab_streq(method, "POST")) {
        curl_easy_setopt(curl, CURLOPT_POST, 1L);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body != NULL ? body : "{}");
    } else if (!ulab_streq(method, "GET")) {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);
        if (body != NULL) {
            curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body);
        }
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

    if (status != NULL) {
        *status = code;
    }

    if (out != NULL) {
        if (resp.buf == NULL) {
            resp.buf = calloc(1, 1);
            if (resp.buf == NULL) {
                snprintf(err->msg, sizeof(err->msg),
                         "%s: out of memory reading response", op);
                return ULAB_ERR;
            }
        }
        *out = resp.buf;
    } else {
        free(resp.buf);
    }

    return ULAB_OK;
}

static int sf_url_escape(const char *in,
                         char *out,
                         size_t out_len,
                         ulab_error_t *err) {
    CURL *curl;
    char *esc;

    curl = curl_easy_init();
    if (curl == NULL) {
        snprintf(err->msg, sizeof(err->msg), "curl init failed");
        return ULAB_ERR;
    }

    esc = curl_easy_escape(curl, in != NULL ? in : "", 0);
    if (esc == NULL) {
        curl_easy_cleanup(curl);
        snprintf(err->msg, sizeof(err->msg), "URL escape failed");
        return ULAB_ERR;
    }

    if (ulab_copy(out, out_len, esc)) {
        curl_free(esc);
        curl_easy_cleanup(curl);
        snprintf(err->msg, sizeof(err->msg), "escaped URL value too long");
        return ULAB_ERR;
    }

    curl_free(esc);
    curl_easy_cleanup(curl);

    return ULAB_OK;
}

static int sf_warehouse_sim_exists(const runner_opts_t *opts,
                                   const char *iccid,
                                   int *exists,
                                   ulab_error_t *err) {
    char path[ULAB_MAX_PATH];
    char url[ULAB_MAX_PATH];
    char *body;
    long code;
    int rc;

    body = NULL;
    *exists = 0;

    rc = snprintf(path, sizeof(path), "/v1/sims/%s", iccid);
    if (rc < 0 || (size_t)rc >= sizeof(path)) {
        snprintf(err->msg, sizeof(err->msg), "warehouse sim path too long");
        return ULAB_ERR;
    }

    if (sf_join_url(url, sizeof(url), opts->warehouse_url, path, err)) {
        return ULAB_ERR;
    }

    if (sf_http_request("warehouseGetSim", "GET", url, NULL, &body,
                        &code, err)) {
        return ULAB_ERR;
    }
    free(body);

    if (code == 200) {
        *exists = 1;
        return ULAB_OK;
    }

    if (code == 404) {
        *exists = 0;
        return ULAB_OK;
    }

    snprintf(err->msg, sizeof(err->msg),
             "warehouseGetSim iccid=%s HTTP %ld", iccid, code);
    return ULAB_ERR;
}

static int sf_warehouse_add_sim(const runner_opts_t *opts,
                                const char *batch_id,
                                const char *iccid,
                                const char *imsi,
                                ulab_error_t *err) {
    char url[ULAB_MAX_PATH];
    char body[ULAB_MAX_QUERY];
    char iccid_esc[ULAB_MAX_ID * 2];
    char imsi_esc[ULAB_MAX_ID * 2];
    char batch_esc[ULAB_MAX_ID * 2];
    char form_esc[ULAB_MAX_REF * 2];
    char profile_esc[ULAB_MAX_REF * 2];
    char vendor_esc[ULAB_MAX_REF * 2];
    char org_esc[ULAB_MAX_REF * 2];
    char *resp;
    long code;
    int rc;

    resp = NULL;

    ulab_json_escape(iccid, iccid_esc, sizeof(iccid_esc));
    ulab_json_escape(imsi, imsi_esc, sizeof(imsi_esc));
    ulab_json_escape(batch_id, batch_esc, sizeof(batch_esc));
    ulab_json_escape(opts->sim_form_factor, form_esc, sizeof(form_esc));
    ulab_json_escape(opts->sim_profile, profile_esc, sizeof(profile_esc));
    ulab_json_escape(opts->sim_vendor, vendor_esc, sizeof(vendor_esc));
    ulab_json_escape(opts->sim_org, org_esc, sizeof(org_esc));

    rc = snprintf(body, sizeof(body),
                  "{\"batch_id\":\"%s\","
                  "\"form_factor\":\"%s\","
                  "\"iccid\":\"%s\","
                  "\"imsi\":\"%s\","
                  "\"profile\":\"%s\","
                  "\"vendor\":\"%s\","
                  "\"org_name\":\"%s\"}",
                  batch_esc, form_esc, iccid_esc, imsi_esc,
                  profile_esc, vendor_esc, org_esc);
    if (rc < 0 || (size_t)rc >= sizeof(body)) {
        snprintf(err->msg, sizeof(err->msg), "warehouse add body too long");
        return ULAB_ERR;
    }

    if (sf_join_url(url, sizeof(url), opts->warehouse_url, "/v1/sims", err)) {
        return ULAB_ERR;
    }

    if (sf_http_request("warehouseAddSim", "POST", url, body, &resp,
                        &code, err)) {
        return ULAB_ERR;
    }

    if (code < 200 || code >= 300) {
        snprintf(err->msg, sizeof(err->msg),
                 "warehouseAddSim iccid=%s HTTP %ld: %s", iccid, code,
                 resp != NULL ? resp : "");
        free(resp);
        return ULAB_ERR;
    }

    free(resp);
    return ULAB_OK;
}

static int sf_factory_batch_count(const runner_opts_t *opts,
                                  const char *batch_id,
                                  size_t expected_count,
                                  size_t *actual_count,
                                  ulab_error_t *err) {
    char batch_q[ULAB_MAX_ID * 3];
    char org_q[ULAB_MAX_REF * 3];
    char type_q[ULAB_MAX_REF * 3];
    char path[ULAB_MAX_PATH];
    char url[ULAB_MAX_PATH];
    char *body;
    long code;
    const char *p;
    int rc;

    body = NULL;
    *actual_count = 0;

    if (sf_url_escape(batch_id, batch_q, sizeof(batch_q), err) ||
        sf_url_escape(opts->sim_org, org_q, sizeof(org_q), err) ||
        sf_url_escape(opts->sim_type, type_q, sizeof(type_q), err)) {
        return ULAB_ERR;
    }

    rc = snprintf(path, sizeof(path),
                  "/v1/sims?imsi=&batch_id=%s&org_name=%s&sim_type=%s&count=%zu&sort=true",
                  batch_q, org_q, type_q, expected_count);
    if (rc < 0 || (size_t)rc >= sizeof(path)) {
        snprintf(err->msg, sizeof(err->msg), "factory list URL too long");
        return ULAB_ERR;
    }

    if (sf_join_url(url, sizeof(url), opts->factory_url, path, err)) {
        return ULAB_ERR;
    }

    if (sf_http_request("factoryListSims", "GET", url, NULL, &body,
                        &code, err)) {
        return ULAB_ERR;
    }

    if (code < 200 || code >= 300) {
        snprintf(err->msg, sizeof(err->msg),
                 "factoryListSims batch=%s HTTP %ld: %s", batch_id, code,
                 body != NULL ? body : "");
        free(body);
        return ULAB_ERR;
    }

    p = body;
    while (p != NULL && *p != '\0') {
        p = strstr(p, "\"iccid\"");
        if (p == NULL) {
            break;
        }
        (*actual_count)++;
        p += strlen("\"iccid\"");
    }

    free(body);

    return ULAB_OK;
}

static int sf_factory_wait_batch(const runner_opts_t *opts,
                                 const char *batch_id,
                                 size_t expected_count,
                                 ulab_error_t *err) {
    size_t got;
    int i;

    got = 0;
    ulab_status("FACTORY", "wait batch %s count=%zu",
                batch_id, expected_count);

    for (i = 0; i < SIMFACTORY_WAIT_ATTEMPTS; i++) {
        if (sf_factory_batch_count(opts, batch_id, expected_count, &got,
                                   err) == ULAB_OK && got >= expected_count) {
            return ULAB_OK;
        }

        sleep(SIMFACTORY_WAIT_SLEEP_SEC);
    }

    snprintf(err->msg, sizeof(err->msg),
             "factory batch %s not ready: need=%zu got=%zu", batch_id,
             expected_count, got);
    return ULAB_ERR;
}

static int sf_factory_export_csv(const runner_opts_t *opts,
                                 const char *batch_id,
                                 size_t count,
                                 const char *csv_path,
                                 ulab_error_t *err) {
    char batch_q[ULAB_MAX_ID * 3];
    char org_q[ULAB_MAX_REF * 3];
    char type_q[ULAB_MAX_REF * 3];
    char path[ULAB_MAX_PATH];
    char url[ULAB_MAX_PATH];
    char *body;
    FILE *f;
    long code;
    int rc;

    body = NULL;

    if (sf_url_escape(batch_id, batch_q, sizeof(batch_q), err) ||
        sf_url_escape(opts->sim_org, org_q, sizeof(org_q), err) ||
        sf_url_escape(opts->sim_type, type_q, sizeof(type_q), err)) {
        return ULAB_ERR;
    }

    rc = snprintf(path, sizeof(path),
                  "/v1/sims/csv?imsi=&batch_id=%s&org_name=%s&sim_type=%s&count=%zu&sort=true",
                  batch_q, org_q, type_q, count);
    if (rc < 0 || (size_t)rc >= sizeof(path)) {
        snprintf(err->msg, sizeof(err->msg), "factory csv URL too long");
        return ULAB_ERR;
    }

    if (sf_join_url(url, sizeof(url), opts->factory_url, path, err)) {
        return ULAB_ERR;
    }

    if (sf_http_request("factoryExportSimsCsv", "GET", url, NULL, &body,
                        &code, err)) {
        return ULAB_ERR;
    }

    if (code < 200 || code >= 300) {
        snprintf(err->msg, sizeof(err->msg),
                 "factoryExportSimsCsv batch=%s HTTP %ld: %s", batch_id,
                 code, body != NULL ? body : "");
        free(body);
        return ULAB_ERR;
    }

    f = fopen(csv_path, "w");
    if (f == NULL) {
        snprintf(err->msg, sizeof(err->msg), "failed to write %s", csv_path);
        free(body);
        return ULAB_ERR;
    }

    if (fputs(body != NULL ? body : "", f) == EOF) {
        snprintf(err->msg, sizeof(err->msg), "failed to write %s", csv_path);
        fclose(f);
        free(body);
        return ULAB_ERR;
    }

    fclose(f);
    free(body);

    return ULAB_OK;
}

static int sf_luhn_digit(const char *body) {
    int sum;
    int alt;
    int i;
    int n;

    sum = 0;
    alt = 1;
    for (i = (int)strlen(body) - 1; i >= 0; i--) {
        n = body[i] - '0';
        if (alt) {
            n *= 2;
            if (n > 9) {
                n -= 9;
            }
        }
        sum += n;
        alt = !alt;
    }

    return (10 - (sum % 10)) % 10;
}

static int sf_make_iccid(char *out,
                         size_t out_len,
                         uint64_t serial,
                         ulab_error_t *err) {
    char body[32];
    int luhn;
    int rc;

    rc = snprintf(body, sizeof(body), "891030%012" PRIu64,
                  (uint64_t)(serial % 1000000000000ULL));
    if (rc < 0 || (size_t)rc >= sizeof(body)) {
        snprintf(err->msg, sizeof(err->msg), "ICCID body too long");
        return ULAB_ERR;
    }

    luhn = sf_luhn_digit(body);
    rc = snprintf(out, out_len, "%s%d", body, luhn);
    if (rc < 0 || (size_t)rc >= out_len) {
        snprintf(err->msg, sizeof(err->msg), "ICCID too long");
        return ULAB_ERR;
    }

    return ULAB_OK;
}

static int sf_make_imsi(char *out,
                        size_t out_len,
                        uint64_t serial,
                        ulab_error_t *err) {
    int rc;

    rc = snprintf(out, out_len, "001010%09" PRIu64,
                  (uint64_t)(serial % 1000000000ULL));
    if (rc < 0 || (size_t)rc >= out_len) {
        snprintf(err->msg, sizeof(err->msg), "IMSI too long");
        return ULAB_ERR;
    }

    return ULAB_OK;
}

static uint64_t sf_run_serial_seed(const world_t *world) {
    uint32_t h1;
    uint32_t h2;

    h1 = ulab_hash32(world->run_id, world->seed);
    h2 = ulab_hash32(world->run_id, h1 ^ 0xa5a5u);

    return ((((uint64_t)h1) << 32) | h2) % 1000000000000ULL;
}

static int sf_csv_field(char **p, char *out, size_t out_len) {
    char *s;
    char *e;
    int rc;

    if (p == NULL || *p == NULL || out == NULL || out_len == 0) {
        return ULAB_ERR;
    }

    s = *p;
    e = strchr(s, ',');
    if (e != NULL) {
        *e = '\0';
        *p = e + 1;
    } else {
        *p = s + strlen(s);
    }

    s = ulab_trim(s);
    rc = ulab_copy(out, out_len, s);

    return rc;
}

static int sf_assign_ues_from_csv(world_t *world,
                                  const char *csv_path,
                                  ulab_error_t *err) {
    FILE *f;
    char line[ULAB_MAX_LINE * 4];
    size_t idx;

    f = fopen(csv_path, "r");
    if (f == NULL) {
        snprintf(err->msg, sizeof(err->msg), "failed to read %s", csv_path);
        return ULAB_ERR;
    }

    if (fgets(line, sizeof(line), f) == NULL) {
        snprintf(err->msg, sizeof(err->msg), "factory CSV is empty: %s",
                 csv_path);
        fclose(f);
        return ULAB_ERR;
    }

    idx = 0;
    while (idx < world->ue_count && fgets(line, sizeof(line), f) != NULL) {
        char iccid[ULAB_MAX_ID];
        char imsi[ULAB_MAX_ID];
        char *p;

        line[strcspn(line, "\r\n")] = '\0';
        if (line[0] == '\0') {
            continue;
        }

        p = line;
        if (sf_csv_field(&p, iccid, sizeof(iccid)) ||
            sf_csv_field(&p, imsi, sizeof(imsi))) {
            snprintf(err->msg, sizeof(err->msg),
                     "invalid factory CSV row %zu in %s", idx + 2,
                     csv_path);
            fclose(f);
            return ULAB_ERR;
        }

        if (iccid[0] == '\0' || imsi[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "factory CSV row %zu missing ICCID/IMSI", idx + 2);
            fclose(f);
            return ULAB_ERR;
        }

        ulab_copy(world->ues[idx].iccid, sizeof(world->ues[idx].iccid), iccid);
        ulab_copy(world->ues[idx].imsi, sizeof(world->ues[idx].imsi), imsi);
        ulab_status("FACTORY", "ue %s iccid=%s imsi=%s",
                    world->ues[idx].ref, world->ues[idx].iccid,
                    world->ues[idx].imsi);
        idx++;
    }

    fclose(f);

    if (idx != world->ue_count) {
        snprintf(err->msg, sizeof(err->msg),
                 "factory CSV has %zu usable sims, need=%zu", idx,
                 world->ue_count);
        return ULAB_ERR;
    }

    return ULAB_OK;
}

int sim_factory_prepare_world(const runner_opts_t *opts,
                              world_t *world,
                              const char *run_dir,
                              char *csv_path,
                              size_t csv_path_len,
                              ulab_error_t *err) {
    char batch_id[ULAB_MAX_ID];
    uint64_t seed;
    size_t i;
    int rc;

    if (world->ue_count == 0) {
        csv_path[0] = '\0';
        return ULAB_OK;
    }

    if (opts->warehouse_url[0] == '\0' || opts->factory_url[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "SIM provisioning requires --warehouse-url and --factory-url");
        return ULAB_ERR;
    }

    rc = snprintf(batch_id, sizeof(batch_id), "%s-%s",
                  opts->sim_batch_prefix, world->run_id);
    if (rc < 0 || (size_t)rc >= sizeof(batch_id)) {
        snprintf(err->msg, sizeof(err->msg), "SIM batch id too long");
        return ULAB_ERR;
    }

    rc = snprintf(csv_path, csv_path_len, "%s/factory-sims.csv", run_dir);
    if (rc < 0 || (size_t)rc >= csv_path_len) {
        snprintf(err->msg, sizeof(err->msg), "factory CSV path too long");
        return ULAB_ERR;
    }

    ulab_status("FACTORY", "batch %s count=%zu", batch_id,
                world->ue_count);

    seed = sf_run_serial_seed(world);
    for (i = 0; i < world->ue_count; i++) {
        char iccid[ULAB_MAX_ID];
        char imsi[ULAB_MAX_ID];
        uint64_t serial;
        int exists;
        int attempt;
        int added;

        added = 0;
        for (attempt = 0; attempt < 50; attempt++) {
            serial = seed + (uint64_t)i + ((uint64_t)attempt * 1000000ULL);
            if (sf_make_iccid(iccid, sizeof(iccid), serial, err) ||
                sf_make_imsi(imsi, sizeof(imsi), serial, err)) {
                return ULAB_ERR;
            }

            if (sf_warehouse_sim_exists(opts, iccid, &exists, err)) {
                return ULAB_ERR;
            }
            if (exists) {
                continue;
            }

            ulab_status("FACTORY", "add %s iccid=%s imsi=%s",
                        world->ues[i].ref, iccid, imsi);
            if (sf_warehouse_add_sim(opts, batch_id, iccid, imsi, err) == ULAB_OK) {
                added = 1;
                break;
            }

            if (strstr(err->msg, "409") == NULL &&
                strstr(err->msg, "duplicate") == NULL) {
                return ULAB_ERR;
            }
            err->msg[0] = '\0';
        }

        if (!added) {
            snprintf(err->msg, sizeof(err->msg),
                     "failed to generate unused SIM for %s", world->ues[i].ref);
            return ULAB_ERR;
        }
    }

    if (sf_factory_wait_batch(opts, batch_id, world->ue_count, err)) {
        return ULAB_ERR;
    }

    ulab_status("FACTORY", "export %s", csv_path);
    if (sf_factory_export_csv(opts, batch_id, world->ue_count, csv_path, err)) {
        return ULAB_ERR;
    }

    return sf_assign_ues_from_csv(world, csv_path, err);
}

int sim_factory_wait_asr(const runner_opts_t *opts,
                         const ue_t *ue,
                         ulab_error_t *err) {
    char path[ULAB_MAX_PATH];
    char url[ULAB_MAX_PATH];
    int i;
    int rc;

    if (opts->asr_url[0] == '\0') {
        return ULAB_OK;
    }

    rc = snprintf(path, sizeof(path), "/v1/asr/%s", ue->iccid);
    if (rc < 0 || (size_t)rc >= sizeof(path)) {
        snprintf(err->msg, sizeof(err->msg), "ASR path too long");
        return ULAB_ERR;
    }

    if (sf_join_url(url, sizeof(url), opts->asr_url, path, err)) {
        return ULAB_ERR;
    }

    ulab_status("ASR", "wait iccid=%s", ue->iccid);
    for (i = 0; i < ASR_WAIT_ATTEMPTS; i++) {
        char *body;
        long code;

        body = NULL;
        err->msg[0] = '\0';
        if (sf_http_request("asrGet", "GET", url, NULL, &body, &code, err) == ULAB_OK) {
            if (code >= 200 && code < 300) {
                free(body);
                ulab_status("ASR", "ready iccid=%s", ue->iccid);
                return ULAB_OK;
            }
        }
        free(body);
        sleep(ASR_WAIT_SLEEP_SEC);
    }

    snprintf(err->msg, sizeof(err->msg),
             "ASR missing for ICCID=%s after allocation; SIM may not be registered by ukama-agent",
             ue->iccid);
    return ULAB_ERR;
}
