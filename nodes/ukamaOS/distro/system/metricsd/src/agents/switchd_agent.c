/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <stdbool.h>

#include <curl/curl.h>
#include <jansson.h>

#include "agents.h"
#include "server.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_api.h"
#include "usys_file.h"
#include "usys_services.h"

#ifndef RETURN_OK
#define RETURN_OK 0
#endif

#ifndef RETURN_NOTOK
#define RETURN_NOTOK -1
#endif

#define SWITCHD_METRIC_NAME_LEN 128
#define SWITCHD_FQNAME_LEN      256
#define SWITCHD_DESC_LEN        256

typedef struct {
    char    *data;
    size_t   len;
} HttpBuf;

typedef struct {
    char    name[SWITCHD_METRIC_NAME_LEN];
    double  value;
} SwitchMetric;

typedef struct {
    SwitchMetric *items;
    size_t        count;
    size_t        cap;
} SwitchMetricList;

typedef struct DynamicKpi {
    KPIConfig          kpi;
    char               name[SWITCHD_FQNAME_LEN];
    char               desc[SWITCHD_DESC_LEN];
    struct DynamicKpi *next;
} DynamicKpi;

static DynamicKpi *gDynamicKpis = NULL;

static int str_ieq(const char *a, const char *b) {

    if (!a || !b) return 0;

    while (*a && *b) {
        char ca = *a;
        char cb = *b;

        if (ca >= 'A' && ca <= 'Z') ca = (char)(ca - 'A' + 'a');
        if (cb >= 'A' && cb <= 'Z') cb = (char)(cb - 'A' + 'a');

        if (ca != cb) return 0;

        a++;
        b++;
    }

    return (*a == '\0' && *b == '\0');
}

static int str_istartswith(const char *s, const char *prefix) {

    size_t i = 0;

    if (!s || !prefix) return 0;

    while (prefix[i] != '\0') {
        char a = s[i];
        char b = prefix[i];

        if (a == '\0') return 0;

        if (a >= 'A' && a <= 'Z') a = (char)(a - 'A' + 'a');
        if (b >= 'A' && b <= 'Z') b = (char)(b - 'A' + 'a');

        if (a != b) return 0;

        i++;
    }

    return 1;
}

static int str_iendswith(const char *s, const char *suffix) {

    size_t slen;
    size_t tlen;

    if (!s || !suffix) return 0;

    slen = strlen(s);
    tlen = strlen(suffix);

    if (tlen == 0 || slen < tlen) return 0;

    return str_ieq(s + slen - tlen, suffix);
}

static size_t curl_write_cb(void *contents,
                            size_t size,
                            size_t nmemb,
                            void *userp) {

    size_t total;
    HttpBuf *buf;
    char *newMem;

    total = size * nmemb;
    buf = (HttpBuf *)userp;

    newMem = realloc(buf->data, buf->len + total + 1);
    if (!newMem) return 0;

    buf->data = newMem;
    memcpy(buf->data + buf->len, contents, total);
    buf->len += total;
    buf->data[buf->len] = '\0';

    return total;
}

static int http_get_json(const char *url, json_t **outRoot) {

    int ret = RETURN_NOTOK;
    CURL *curl = NULL;
    CURLcode cres;
    HttpBuf buf = {0};
    json_error_t jerr;

    if (!url || !outRoot) return RETURN_NOTOK;

    curl = curl_easy_init();
    if (!curl) {
        usys_log_error("switchd_agent: curl_easy_init failed");
        return RETURN_NOTOK;
    }

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, curl_write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &buf);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, 1500L);
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, 700L);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);

    cres = curl_easy_perform(curl);
    if (cres != CURLE_OK) {
        usys_log_error("switchd_agent: curl failed url=%s err=%s",
                       url, curl_easy_strerror(cres));
        goto done;
    }

    *outRoot = json_loads(buf.data ? buf.data : "{}", 0, &jerr);
    if (!*outRoot) {
        usys_log_error("switchd_agent: json parse failed url=%s line=%d: %s",
                       url, jerr.line, jerr.text);
        goto done;
    }

    ret = RETURN_OK;

done:
    if (curl) curl_easy_cleanup(curl);
    if (buf.data) free(buf.data);

    return ret;
}

static double j_get_num(json_t *root, const char *key) {

    json_t *v;

    if (!root || !key) return 0;

    v = json_object_get(root, key);
    if (!v) return 0;

    if (json_is_real(v)) return json_real_value(v);
    if (json_is_integer(v)) return (double)json_integer_value(v);
    if (json_is_boolean(v)) return json_is_true(v) ? 1.0 : 0.0;

    return 0;
}

static const char *j_get_str(json_t *root, const char *key) {

    json_t *v;

    if (!root || !key) return NULL;

    v = json_object_get(root, key);
    if (!v || !json_is_string(v)) return NULL;

    return json_string_value(v);
}

static json_t *j_get_arr(json_t *root, const char *key) {

    json_t *v;

    if (!root || !key) return NULL;

    v = json_object_get(root, key);
    if (!v || !json_is_array(v)) return NULL;

    return v;
}

static void switchd_metrics_free(SwitchMetricList *list) {

    if (!list) return;

    if (list->items) free(list->items);

    list->items = NULL;
    list->count = 0;
    list->cap = 0;
}

static int switchd_metrics_add(SwitchMetricList *list,
                               const char *name,
                               double value) {

    SwitchMetric *items;
    size_t nextCap;

    if (!list || !name || !*name) return RETURN_NOTOK;

    if (list->count == list->cap) {
        nextCap = (list->cap == 0) ? 32 : list->cap * 2;

        items = realloc(list->items, nextCap * sizeof(SwitchMetric));
        if (!items) return RETURN_NOTOK;

        list->items = items;
        list->cap = nextCap;
    }

    snprintf(list->items[list->count].name,
             sizeof(list->items[list->count].name),
             "%s",
             name);

    list->items[list->count].value = value;
    list->count++;

    return RETURN_OK;
}

static int switchd_read_metrics(const char *url, SwitchMetricList *out) {

    int ret = RETURN_NOTOK;
    json_t *root = NULL;
    json_t *metrics = NULL;
    size_t i;

    if (!out) return RETURN_NOTOK;

    memset(out, 0, sizeof(*out));

    if (http_get_json(url, &root) != RETURN_OK) {
        return RETURN_NOTOK;
    }

    metrics = j_get_arr(root, "metrics");
    if (!metrics) {
        usys_log_error("switchd_agent: missing metrics array in %s", url);
        goto done;
    }

    for (i = 0; i < json_array_size(metrics); i++) {
        json_t *metric;
        const char *name;
        double value;

        metric = json_array_get(metrics, i);
        if (!metric || !json_is_object(metric)) continue;

        name = j_get_str(metric, "name");
        if (!name || !*name) continue;

        value = j_get_num(metric, "value");

        if (switchd_metrics_add(out, name, value) != RETURN_OK) {
            usys_log_error("switchd_agent: failed to store metric %s", name);
            goto done;
        }
    }

    ret = RETURN_OK;

done:
    if (ret != RETURN_OK) {
        switchd_metrics_free(out);
    }

    if (root) {
        json_decref(root);
    }

    return ret;
}

static int build_prefix_from_config(MetricsCatConfig *cfgStat,
                                    SwitchMetricList *metrics,
                                    char *prefix,
                                    size_t prefixSize) {

    int idx;
    size_t i;

    if (!cfgStat || !metrics || !prefix || prefixSize == 0) {
        return RETURN_NOTOK;
    }

    prefix[0] = '\0';

    for (idx = 0; idx < cfgStat->kpiCount; idx++) {
        KPIConfig *kpi = &(cfgStat->kpi[idx]);

        if (!kpi || !kpi->fqname) continue;

        for (i = 0; i < metrics->count; i++) {
            const char *metricName = metrics->items[i].name;
            size_t fqLen;
            size_t metricLen;
            size_t prefixLen;

            if (str_istartswith(metricName, "port_")) {
                continue;
            }

            if (!str_iendswith(kpi->fqname, metricName)) {
                continue;
            }

            fqLen = strlen(kpi->fqname);
            metricLen = strlen(metricName);

            if (fqLen <= metricLen) continue;

            prefixLen = fqLen - metricLen;
            if (prefixLen >= prefixSize) {
                return RETURN_NOTOK;
            }

            memcpy(prefix, kpi->fqname, prefixLen);
            prefix[prefixLen] = '\0';

            return RETURN_OK;
        }
    }

    return RETURN_NOTOK;
}

static DynamicKpi *dynamic_kpi_find(const char *fqname) {

    DynamicKpi *cur = gDynamicKpis;

    while (cur != NULL) {
        if (strcmp(cur->name, fqname) == 0) {
            return cur;
        }

        cur = cur->next;
    }

    return NULL;
}

static KPIConfig *dynamic_kpi_get_or_create(const char *fqname) {

    DynamicKpi *dyn;

    if (!fqname || !*fqname) return NULL;

    dyn = dynamic_kpi_find(fqname);
    if (dyn != NULL) {
        return &dyn->kpi;
    }

    dyn = calloc(1, sizeof(DynamicKpi));
    if (dyn == NULL) {
        usys_log_error("switchd_agent: oom creating dynamic kpi %s",
                       fqname);
        return NULL;
    }

    snprintf(dyn->name, sizeof(dyn->name), "%s", fqname);
    snprintf(dyn->desc, sizeof(dyn->desc), "Switch metric %s", fqname);

    dyn->kpi.name      = dyn->name;
    dyn->kpi.fqname    = dyn->name;
    dyn->kpi.desc      = dyn->desc;
    dyn->kpi.type      = METRICTYPE_GAUGE;
    dyn->kpi.numLabels = 0;
    dyn->kpi.labels    = NULL;
    dyn->kpi.unit      = NULL;
    dyn->kpi.ext       = NULL;

    if (metric_server_register_kpi(&dyn->kpi) != RETURN_OK) {
        usys_log_error("switchd_agent: failed to register dynamic kpi %s",
                       fqname);
        free(dyn);
        return NULL;
    }

    dyn->next = gDynamicKpis;
    gDynamicKpis = dyn;

    usys_log_info("switchd_agent: registered dynamic kpi %s", fqname);

    return &dyn->kpi;
}

static int metric_matches_config_kpi(const SwitchMetric *metric,
                                     KPIConfig *kpi) {

    if (!metric || !kpi || !kpi->fqname) return 0;

    if (str_ieq(metric->name, kpi->fqname)) {
        return 1;
    }

    if (str_iendswith(kpi->fqname, metric->name)) {
        return 1;
    }

    return 0;
}

static int switchd_push_configured_kpis(MetricsCatConfig *cfgStat,
                                        SwitchMetricList *metrics,
                                        metricAddFunc addFunc) {

    int idx;
    size_t i;

    if (!cfgStat || !metrics || !addFunc) return RETURN_NOTOK;

    for (idx = 0; idx < cfgStat->kpiCount; idx++) {
        KPIConfig *kpi = &(cfgStat->kpi[idx]);

        if (!kpi || !kpi->fqname) continue;

        for (i = 0; i < metrics->count; i++) {
            double val;

            if (!metric_matches_config_kpi(&metrics->items[i], kpi)) {
                continue;
            }

            val = metrics->items[i].value;
            addFunc(kpi, &val);
            break;
        }
    }

    return RETURN_OK;
}

static int switchd_push_dynamic_port_kpis(MetricsCatConfig *cfgStat,
                                          SwitchMetricList *metrics,
                                          metricAddFunc addFunc) {

    char prefix[SWITCHD_FQNAME_LEN];
    size_t i;

    if (!cfgStat || !metrics || !addFunc) return RETURN_NOTOK;

    if (build_prefix_from_config(cfgStat,
                                 metrics,
                                 prefix,
                                 sizeof(prefix)) != RETURN_OK) {
        usys_log_error("switchd_agent: failed to derive switch kpi prefix");
        return RETURN_OK;
    }

    for (i = 0; i < metrics->count; i++) {
        char fqname[SWITCHD_FQNAME_LEN];
        KPIConfig *kpi;
        double val;
        int ret;

        if (!str_istartswith(metrics->items[i].name, "port_")) {
            continue;
        }

        ret = snprintf(fqname,
                       sizeof(fqname),
                       "%s%s",
                       prefix,
                       metrics->items[i].name);
        if (ret < 0 || ret >= (int)sizeof(fqname)) {
            usys_log_error("switchd_agent: dynamic fqname too long for %s",
                           metrics->items[i].name);
            continue;
        }

        kpi = dynamic_kpi_get_or_create(fqname);
        if (kpi == NULL) {
            continue;
        }

        val = metrics->items[i].value;
        addFunc(kpi, &val);
    }

    return RETURN_OK;
}

static int switchd_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                              SwitchMetricList *metrics,
                                              metricAddFunc addFunc) {

    if (!cfgStat || !metrics || !addFunc) return RETURN_NOTOK;

    switchd_push_configured_kpis(cfgStat, metrics, addFunc);
    switchd_push_dynamic_port_kpis(cfgStat, metrics, addFunc);

    return RETURN_OK;
}

int switchd_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {

    int port;
    char urlBuf[256];
    const char *path;
    SwitchMetricList metrics;

    memset(&metrics, 0, sizeof(metrics));

    port = usys_find_service_port(SERVICE_SWITCH);
    if (port <= 0) {
        usys_log_error("switchd_agent: could not resolve service port for %s",
                       SERVICE_SWITCH);
        return RETURN_NOTOK;
    }

    path = (cfgStat && cfgStat->url && cfgStat->url[0]) ?
           cfgStat->url : "/v1/metrics";

    snprintf(urlBuf, sizeof(urlBuf), "http://127.0.0.1:%d%s", port, path);

    if (switchd_read_metrics(urlBuf, &metrics) != RETURN_OK) {
        usys_log_error("switchd_agent: failed to read %s", urlBuf);
        return RETURN_NOTOK;
    }

    if (switchd_push_stat_to_metric_server(cfgStat,
                                           &metrics,
                                           addFunc) != RETURN_OK) {
        usys_log_error("switchd_agent: failed to push metrics to server");
        switchd_metrics_free(&metrics);
        return RETURN_NOTOK;
    }

    switchd_metrics_free(&metrics);

    return RETURN_OK;
}
