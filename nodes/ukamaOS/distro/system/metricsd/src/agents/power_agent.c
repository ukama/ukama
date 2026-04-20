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

#include <curl/curl.h>
#include <jansson.h>

#include "agents.h"

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

#ifndef SERVICE_POWER
#define SERVICE_POWER "power"
#endif

typedef struct {
    char    *data;
    size_t  len;
} HttpBuf;

typedef struct {
    int     ok;
    int     haveLm75;
    int     haveLm25066;
    int     haveAds1015;

    double  totalWatts;
    double  energyWh;

    double  boardTempC;

    double  inVolts;
    double  outVolts;
    double  inAmps;
    double  inWatts;
    double  hsTempC;
    double  statusWord;
    double  diagnosticWord;
    double  assumedDirect;

    double  adcVin;
    double  adcVpa;
    double  adcAux;
} PowerMetrics;

static int str_icontains(const char *hay, const char *needle) {

    size_t i = 0;
    size_t nlen = 0;
    size_t hlen = 0;

    if (!hay || !needle) return 0;

    hlen = strlen(hay);
    nlen = strlen(needle);

    if (nlen == 0 || hlen < nlen) return 0;

    for (i = 0; i + nlen <= hlen; i++) {
        size_t j = 0;

        for (j = 0; j < nlen; j++) {
            char c1 = hay[i + j];
            char c2 = needle[j];

            if (c1 >= 'A' && c1 <= 'Z') c1 = (char)(c1 - 'A' + 'a');
            if (c2 >= 'A' && c2 <= 'Z') c2 = (char)(c2 - 'A' + 'a');

            if (c1 != c2) break;
        }

        if (j == nlen) return 1;
    }

    return 0;
}

static size_t curl_write_cb(void *contents, size_t size, size_t nmemb,
                            void *userp) {

    size_t total = size * nmemb;
    HttpBuf *buf = (HttpBuf *)userp;
    char *newMem = NULL;

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
    CURLcode cres = CURLE_OK;
    HttpBuf buf = {0};
    json_error_t jerr;

    if (!url || !outRoot) return RETURN_NOTOK;

    curl = curl_easy_init();
    if (!curl) {
        usys_log_error("power_agent: curl_easy_init failed");
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
        usys_log_error("power_agent: curl failed url=%s err=%s",
                       url, curl_easy_strerror(cres));
        goto done;
    }

    *outRoot = json_loads(buf.data ? buf.data : "{}", 0, &jerr);
    if (!*outRoot) {
        usys_log_error("power_agent: json parse failed url=%s line=%d: %s",
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

    json_t *v = NULL;

    if (!root || !key) return 0;

    v = json_object_get(root, key);
    if (!v) return 0;

    if (json_is_real(v)) return json_real_value(v);
    if (json_is_integer(v)) return (double)json_integer_value(v);

    return 0;
}

static int j_get_bool(json_t *root, const char *key, int *out) {

    json_t *v = NULL;

    if (!root || !key || !out) return RETURN_NOTOK;

    v = json_object_get(root, key);
    if (!v || !json_is_boolean(v)) return RETURN_NOTOK;

    *out = json_boolean_value(v) ? 1 : 0;
    return RETURN_OK;
}

static json_t *j_get_obj(json_t *root, const char *key) {

    json_t *v = NULL;

    if (!root || !key) return NULL;

    v = json_object_get(root, key);
    if (!v || !json_is_object(v)) return NULL;

    return v;
}

static int power_read_metrics(const char *url, PowerMetrics *out) {

    int ret = RETURN_NOTOK;
    json_t *root = NULL;
    json_t *lm75 = NULL;
    json_t *lm25066 = NULL;
    json_t *ads1015 = NULL;

    if (!out) return RETURN_NOTOK;
    memset(out, 0, sizeof(*out));

    if (http_get_json(url, &root) != RETURN_OK) {
        return RETURN_NOTOK;
    }

    j_get_bool(root, "ok", &out->ok);
    out->totalWatts = j_get_num(root, "totalWatts");
    out->energyWh   = j_get_num(root, "energyWh");

    lm75 = j_get_obj(root, "lm75");
    if (lm75) {
        out->haveLm75 = 1;
        out->boardTempC = j_get_num(lm75, "boardTempC");
    }

    lm25066 = j_get_obj(root, "lm25066");
    if (lm25066) {
        out->haveLm25066 = 1;
        out->inVolts = j_get_num(lm25066, "inVolts");
        out->outVolts = j_get_num(lm25066, "outVolts");
        out->inAmps = j_get_num(lm25066, "inAmps");
        out->inWatts = j_get_num(lm25066, "inWatts");
        out->hsTempC = j_get_num(lm25066, "hsTempC");
        out->statusWord = j_get_num(lm25066, "statusWord");
        out->diagnosticWord = j_get_num(lm25066, "diagnosticWord");
        out->assumedDirect = j_get_num(lm25066, "assumedDirect");
    }

    ads1015 = j_get_obj(root, "ads1015");
    if (ads1015) {
        out->haveAds1015 = 1;
        out->adcVin = j_get_num(ads1015, "adcVin");
        out->adcVpa = j_get_num(ads1015, "adcVpa");
        out->adcAux = j_get_num(ads1015, "adcAux");
    }

    ret = RETURN_OK;

    if (root) json_decref(root);
    return ret;
}

static int power_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                            PowerMetrics *m,
                                            metricAddFunc addFunc) {

    int ret = RETURN_OK;

    for (int idx = 0; idx < cfgStat->kpiCount; idx++) {
        KPIConfig *kpi = &(cfgStat->kpi[idx]);
        double val = 0;

        if (!kpi || !kpi->fqname) continue;

        if (str_icontains(kpi->fqname, "power_ok")) {
            val = m->ok ? 1.0 : 0.0;

        } else if (str_icontains(kpi->fqname, "power_total_watts")) {
            val = m->totalWatts;

        } else if (str_icontains(kpi->fqname, "power_energy_wh")) {
            val = m->energyWh;

        } else if (str_icontains(kpi->fqname, "power_board_temp_c")) {
            if (!m->haveLm75) continue;
            val = m->boardTempC;

        } else if (str_icontains(kpi->fqname, "power_lm25066_in_volts")) {
            if (!m->haveLm25066) continue;
            val = m->inVolts;

        } else if (str_icontains(kpi->fqname, "power_lm25066_out_volts")) {
            if (!m->haveLm25066) continue;
            val = m->outVolts;

        } else if (str_icontains(kpi->fqname, "power_lm25066_in_amps")) {
            if (!m->haveLm25066) continue;
            val = m->inAmps;

        } else if (str_icontains(kpi->fqname, "power_lm25066_in_watts")) {
            if (!m->haveLm25066) continue;
            val = m->inWatts;

        } else if (str_icontains(kpi->fqname, "power_lm25066_temp_c")) {
            if (!m->haveLm25066) continue;
            val = m->hsTempC;

        } else if (str_icontains(kpi->fqname, "power_lm25066_status_word")) {
            if (!m->haveLm25066) continue;
            val = m->statusWord;

        } else if (str_icontains(kpi->fqname, "power_lm25066_diagnostic_word")) {
            if (!m->haveLm25066) continue;
            val = m->diagnosticWord;

        } else if (str_icontains(kpi->fqname, "power_lm25066_assumed_direct")) {
            if (!m->haveLm25066) continue;
            val = m->assumedDirect;

        } else if (str_icontains(kpi->fqname, "power_ads1015_adc_vin")) {
            if (!m->haveAds1015) continue;
            val = m->adcVin;

        } else if (str_icontains(kpi->fqname, "power_ads1015_adc_vpa")) {
            if (!m->haveAds1015) continue;
            val = m->adcVpa;

        } else if (str_icontains(kpi->fqname, "power_ads1015_adc_aux")) {
            if (!m->haveAds1015) continue;
            val = m->adcAux;

        } else {
            continue;
        }

        addFunc(kpi, &val);
    }

    return ret;
}

int power_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {

    int ret = RETURN_OK;
    int port = 0;
    char urlBuf[256] = {0};
    const char *path = NULL;

    PowerMetrics *m = calloc(1, sizeof(PowerMetrics));
    if (!m) {
        usys_log_error("power_agent: oom allocating metrics");
        return RETURN_NOTOK;
    }

    port = usys_find_service_port(SERVICE_POWER);
    if (port <= 0) {
        usys_log_error("power_agent: could not resolve service port for '%s'",
                       SERVICE_POWER);
        free(m);
        return RETURN_NOTOK;
    }

    path = (cfgStat && cfgStat->url && cfgStat->url[0]) ?
           cfgStat->url : "/v1/metrics";

    snprintf(urlBuf, sizeof(urlBuf), "http://127.0.0.1:%d%s", port, path);

    if (power_read_metrics(urlBuf, m) != RETURN_OK) {
        usys_log_error("power_agent: failed to read %s", urlBuf);
        free(m);
        return RETURN_NOTOK;
    }

    if (power_push_stat_to_metric_server(cfgStat, m, addFunc) != RETURN_OK) {
        usys_log_error("power_agent: failed to push metrics to server");
        ret = RETURN_NOTOK;
    }

    free(m);
    return ret;
}
