/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "config.h"
#include "server.h"
#include "toml.h"

#include "usys_log.h"

static char gNode[32] = {'\0'};

static char *alloc_str(int size) {

    char *str = NULL;

    str = calloc(1, sizeof(char) * size);
    if (str == NULL) {
        usys_log_error("failed to allocate string of size %d", size);
    }

    return str;
}

static char *read_str_value(toml_table_t *tab, char *key) {

    int len            = 0;
    char *value        = NULL;
    toml_datum_t str   = toml_string_in(tab, key);

    if (!str.ok) {
        usys_log_error("failed to read string value for key %s", key);
        return NULL;
    }

    len = strlen(str.u.s);
    usys_log_trace("string key=%s value=%s length=%d", key, str.u.s, len);

    value = calloc(len + 1, sizeof(char));
    if (value != NULL) {
        memcpy(value, str.u.s, len);
        value[len] = '\0';
    }

    free(str.u.s);

    return value;
}

static char *read_opt_str_value(toml_table_t *tab, char *key) {

    toml_datum_t str = toml_string_in(tab, key);

    if (!str.ok) {
        return NULL;
    }

    return read_str_value(tab, key);
}

static int read_int_value(toml_table_t *tab, char *key) {

    toml_datum_t val = toml_int_in(tab, key);

    if (!val.ok) {
        usys_log_error("failed to read integer value for key %s", key);
        return 0;
    }

    usys_log_trace("integer key=%s value=%d", key, val.u.i);

    return val.u.i;
}

static void free_str_value(char *str) {

    if (str != NULL) {
        free(str);
    }
}

static char **alloc_labels(int count) {

    return calloc(count, sizeof(char *));
}

static void free_labels(char **labels, int count) {

    int idx = 0;

    if (labels == NULL) {
        return;
    }

    for (idx = 0; idx < count; idx++) {
        if (labels[idx] != NULL) {
            free(labels[idx]);
            labels[idx] = NULL;
        }
    }

    free(labels);
}

static KPIConfig *alloc_kpi(int count) {

    return calloc(count, sizeof(KPIConfig));
}

static void free_kpi(KPIConfig *kpi, int count) {

    int idx = 0;

    if (kpi == NULL) {
        return;
    }

    for (idx = 0; idx < count; idx++) {
        free_str_value(kpi[idx].name);
        free_str_value(kpi[idx].fqname);
        free_str_value(kpi[idx].ext);
        free_str_value(kpi[idx].desc);
        free_str_value(kpi[idx].unit);
        free_labels(kpi[idx].labels, kpi[idx].numLabels);
        metric_server_free_kpi(&kpi[idx]);
    }

    free(kpi);
}

static int *alloc_range(int count) {

    return calloc(count, sizeof(int));
}

static void free_range(int *range) {

    if (range != NULL) {
        free(range);
    }
}

static MetricsCatConfig *alloc_stat(int count) {

    return calloc(count, sizeof(MetricsCatConfig));
}

static void free_stat(MetricsCatConfig *stat, int count) {

    int idx = 0;

    if (stat == NULL) {
        return;
    }

    for (idx = 0; idx < count; idx++) {
        int kpiInstanceCount = 1;

        free_str_value(stat[idx].source);
        free_str_value(stat[idx].agent);
        free_str_value(stat[idx].url);

        if (stat[idx].range != NULL) {
            free_range(stat[idx].range);
            stat[idx].range = NULL;
        }

        if (stat[idx].kpi != NULL) {
            kpiInstanceCount = (stat[idx].instances > 0) ?
                               stat[idx].instances : 1;
            free_kpi(stat[idx].kpi,
                     stat[idx].kpiCount * kpiInstanceCount);
        }
    }

    free(stat);
}

void free_stat_cfg(MetricsConfig *statCfg, int count) {

    int idx = 0;

    if (statCfg == NULL) {
        return;
    }

    for (idx = 0; idx < count; idx++) {
        free_str_value(statCfg[idx].name);

        if (statCfg[idx].metricsCategory != NULL) {
            free_stat(statCfg[idx].metricsCategory,
                      statCfg[idx].eachCategoryCount);
        }
    }

    free(statCfg);
}

static MetricsConfig *alloc_stat_cfg(int count) {

    return calloc(count, sizeof(MetricsConfig));
}

static void copy_lower(char *dst, size_t size, char *src) {

    size_t i = 0;

    if ((dst == NULL) || (size == 0)) {
        return;
    }

    dst[0] = '\0';

    if (src == NULL) {
        return;
    }

    while ((src[i] != '\0') && (i < (size - 1))) {
        if ((src[i] >= 'A') && (src[i] <= 'Z')) {
            dst[i] = src[i] + 32;
        } else {
            dst[i] = src[i];
        }
        i++;
    }

    dst[i] = '\0';
}

static char *set_fqkname(char *node,
                         char *category,
                         char *source,
                         int range,
                         char *kpi,
                         char *unit,
                         char *ext) {

    int len                                  = 0;
    char *name                               = NULL;
    char fqkn[MAX_KPI_KEY_NAME_LENGTH]       = {'\0'};
    char extKpiName[MAX_KPI_KEY_NAME_LENGTH] = {'\0'};
    char nodeName[MAX_KPI_KEY_NAME_LENGTH]   = {'\0'};
    char catName[MAX_KPI_KEY_NAME_LENGTH]    = {'\0'};
    char srcName[MAX_KPI_KEY_NAME_LENGTH]    = {'\0'};
    char kpiName[MAX_KPI_KEY_NAME_LENGTH]    = {'\0'};
    char unitName[MAX_KPI_KEY_NAME_LENGTH]   = {'\0'};
    char extName[MAX_KPI_KEY_NAME_LENGTH]    = {'\0'};
    char srcWithRange[MAX_KPI_KEY_NAME_LENGTH] = {'\0'};

    copy_lower(nodeName, sizeof(nodeName), node);
    copy_lower(catName, sizeof(catName), category);
    copy_lower(srcName, sizeof(srcName), source);
    copy_lower(kpiName, sizeof(kpiName), kpi);
    copy_lower(unitName, sizeof(unitName), unit);
    copy_lower(extName, sizeof(extName), ext);

    if (unitName[0] != '\0') {
        len = snprintf(extKpiName, sizeof(extKpiName), "%s%s%s",
                       kpiName, TAG_SEP, unitName);
    } else {
        len = snprintf(extKpiName, sizeof(extKpiName), "%s", kpiName);
    }

    if ((len < 0) || (len >= (int)sizeof(extKpiName))) {
        return NULL;
    }

    if (range < 0) {
        len = snprintf(srcWithRange, sizeof(srcWithRange), "%s", srcName);
    } else {
        len = snprintf(srcWithRange, sizeof(srcWithRange), "%s%d",
                       srcName, range);
    }

    if ((len < 0) || (len >= (int)sizeof(srcWithRange))) {
        return NULL;
    }

    if (extName[0] != '\0') {
        len = snprintf(fqkn, sizeof(fqkn), "%s%s%s%s%s%s%s%s%s",
                       nodeName, TAG_SEP,
                       catName, TAG_SEP,
                       srcWithRange, TAG_SEP,
                       extName, TAG_SEP,
                       extKpiName);
    } else {
        len = snprintf(fqkn, sizeof(fqkn), "%s%s%s%s%s%s%s",
                       nodeName, TAG_SEP,
                       catName, TAG_SEP,
                       srcWithRange, TAG_SEP,
                       extKpiName);
    }

    if ((len < 0) || (len >= (int)sizeof(fqkn))) {
        return NULL;
    }

    name = alloc_str(len + 1);
    if (name != NULL) {
        memcpy(name, fqkn, (size_t)len);
        name[len] = '\0';
    }

    return name;
}

static int read_metric_type(toml_table_t *tabKpi) {

    int metricType = METRICTYPE_GAUGE;
    char *type     = NULL;

    type = read_str_value(tabKpi, TAG_METRIC_TYPE);
    if (type != NULL) {
        if (strcmp(type, "METRICTYPE_COUNTER") == 0) {
            metricType = METRICTYPE_COUNTER;
        } else if (strcmp(type, "METRICTYPE_GAUGE") == 0) {
            metricType = METRICTYPE_GAUGE;
        } else if (strcmp(type, "METRICTYPE_HISTOGRAM") == 0) {
            metricType = METRICTYPE_HISTOGRAM;
        }

        free_str_value(type);
    } else {
        usys_log_error("failed to read metric type, defaulting to gauge");
    }

    return metricType;
}

static int *parse_range_array(int *inst, toml_table_t *tab, char *key) {

    int idx                = 0;
    int *range             = NULL;
    toml_array_t *arrRange = toml_array_in(tab, key);

    if (arrRange == NULL) {
        usys_log_error("missing %s", key);
        *inst = 0;
        return NULL;
    }

    *inst = toml_array_nelem(arrRange);
    usys_log_trace("%d range values available", *inst);

    if (*inst <= 0) {
        return NULL;
    }

    range = alloc_range(*inst);
    if (range == NULL) {
        usys_log_error("failed to allocate %d range elements", *inst);
        return NULL;
    }

    for (idx = 0; idx < *inst; idx++) {
        toml_datum_t data = toml_int_at(arrRange, idx);

        if (!data.ok) {
            free(range);
            return NULL;
        }

        range[idx] = (int)data.u.i;
        usys_log_trace("range value at [%d] is %d", idx, range[idx]);
    }

    return range;
}

static int *parse_range(int *inst, toml_table_t *tab, char *key) {

    toml_array_t *arrRange = toml_array_in(tab, key);

    if (arrRange == NULL) {
        return NULL;
    }

    return parse_range_array(inst, tab, key);
}

static char **read_labels(int *count, toml_table_t *tabKpi) {

    int idx                = 0;
    int labelsCount        = 0;
    int len                = 0;
    char **labels          = NULL;
    toml_array_t *arrKpi   = toml_array_in(tabKpi, TAG_LABELS);

    if (arrKpi == NULL) {
        usys_log_error("missing %s", TAG_LABELS);
        return NULL;
    }

    *count = toml_array_nelem(arrKpi);
    usys_log_trace("%d labels available", *count);

    if (*count <= 0) {
        return NULL;
    }

    labels = alloc_labels(*count);
    if (labels == NULL) {
        usys_log_error("failed to allocate %d labels", *count);
        return NULL;
    }

    for (idx = 0; idx < *count; idx++) {
        toml_datum_t dlabel = toml_string_at(arrKpi, idx);

        if (!dlabel.ok) {
            *count = labelsCount;
            break;
        }

        len = strlen(dlabel.u.s);
        labels[idx] = calloc(len + 1, sizeof(char));
        if (labels[idx] != NULL) {
            memcpy(labels[idx], dlabel.u.s, len);
            labels[idx][len] = '\0';
        }

        labelsCount++;
        free(dlabel.u.s);
    }

    return labels;
}

static int toml_parse_kpi_table(char *category,
                                char *source,
                                int range,
                                KPIConfig *kpi,
                                toml_table_t *tabKpi) {

    int ret = RETURN_NOTOK;

    kpi->name   = read_str_value(tabKpi, TAG_NAME);
    kpi->ext    = read_opt_str_value(tabKpi, TAG_EXT);
    kpi->desc   = read_str_value(tabKpi, TAG_DESC);
    kpi->unit   = read_opt_str_value(tabKpi, TAG_UNIT);
    kpi->type   = read_metric_type(tabKpi);
    kpi->labels = read_labels(&kpi->numLabels, tabKpi);

    kpi->fqname = set_fqkname(gNode, category, source, range, kpi->name,
                              kpi->unit, kpi->ext);

    ret = metric_server_register_kpi(kpi);

    return ret;
}

static int toml_parse_stat_table(char *category,
                                 MetricsCatConfig *stat,
                                 toml_table_t *tabStat) {

    int ret                 = RETURN_OK;
    int clmn                = 0;
    int rkpiCount           = 0;
    int maxInstances        = 1;
    int inst                = 0;
    int idx                 = 0;
    toml_table_t *tabSource = NULL;
    toml_array_t *arrKpi    = NULL;

    clmn = toml_table_nkval(tabStat);
    usys_log_trace("number of columns in stat table %d", clmn);

    stat->source = read_str_value(tabStat, TAG_SOURCE);
    stat->agent  = read_str_value(tabStat, TAG_AGENT);
    stat->url    = read_str_value(tabStat, TAG_URL);
    stat->range  = parse_range(&(stat->instances), tabStat, TAG_RANGE);

    tabSource = toml_table_in(tabStat, stat->source);
    if (tabSource == NULL) {
        usys_log_error("missing table %s", stat->source);
        return RETURN_NOTOK;
    }

    arrKpi = toml_array_in(tabSource, TAG_KPI);
    if (arrKpi == NULL) {
        usys_log_error("missing %s for source", TAG_KPI);
        return RETURN_NOTOK;
    }

    stat->kpiCount = toml_array_nelem(arrKpi);
    usys_log_trace("%d kpi available", stat->kpiCount);

    rkpiCount = stat->kpiCount;
    if ((stat->instances > 0) && (stat->range != NULL)) {
        rkpiCount = stat->kpiCount * stat->instances;
    }

    stat->kpi = alloc_kpi(rkpiCount);
    if (stat->kpi == NULL) {
        usys_log_error("failed to allocate %d kpi", rkpiCount);
        return RETURN_NOTOK;
    }

    maxInstances = (stat->instances > 0) ? stat->instances : 1;

    usys_log_trace("range has %d instances, max instances %d, category %s, "
                   "source %s",
                   stat->instances, maxInstances, category, stat->source);

    for (inst = 0; inst < maxInstances; inst++) {
        for (idx = 0; idx < stat->kpiCount; idx++) {
            int rangeValue          = -1;
            toml_table_t *tabKpi    = toml_table_at(arrKpi, idx);

            if (tabKpi == NULL) {
                return RETURN_NOTOK;
            }

            if (stat->range != NULL) {
                rangeValue = stat->range[inst];
            }

            ret = toml_parse_kpi_table(category, stat->source, rangeValue,
                                       &(stat->kpi[idx +
                                         (inst * stat->kpiCount)]),
                                       tabKpi);
            if (ret != RETURN_OK) {
                return RETURN_NOTOK;
            }
        }
    }

    return ret;
}

static int toml_parse_stats_cat(char *category,
                                MetricsCatConfig *stat,
                                toml_table_t *tabStats) {

    const char *key = toml_key_in(tabStats, 0);

    usys_log_trace("in stat category table is %s", key);

    if (toml_parse_stat_table(category, stat, tabStats) != RETURN_OK) {
        return RETURN_NOTOK;
    }

    return RETURN_OK;
}

static int toml_parse_stat_cat_array(MetricsConfig *statCfg,
                                     toml_array_t *statCatArr) {

    int idx        = 0;
    int tableCount = 0;
    MetricsCatConfig *stats = NULL;

    tableCount = toml_array_nelem(statCatArr);
    usys_log_trace("stats category contains %d tables", tableCount);

    statCfg->eachCategoryCount = tableCount;
    statCfg->metricsCategory   = alloc_stat(tableCount);
    if (statCfg->metricsCategory == NULL) {
        usys_log_error("failed to allocate %d stat tables", tableCount);
        return RETURN_NOTOK;
    }

    stats = statCfg->metricsCategory;

    for (idx = 0; idx < tableCount; idx++) {
        toml_table_t *tabStat = toml_table_at(statCatArr, idx);

        if (tabStat == NULL) {
            usys_log_error("failed to read stats table at %d for %s",
                           idx, statCfg->name);
            return RETURN_NOTOK;
        }

        if (toml_parse_stats_cat(statCfg->name, &stats[idx], tabStat) !=
            RETURN_OK) {
            usys_log_error("failed to parse stats table at %d for %s",
                           idx, statCfg->name);
            return RETURN_NOTOK;
        }
    }

    return RETURN_OK;
}

int toml_parse_config(char *cfg,
                      char **version,
                      int *scrapingTimePeriod,
                      MetricsConfig **pstatCfg,
                      int *catCount) {

    int idx                   = 0;
    int categoryCounts        = 0;
    FILE *fp                  = NULL;
    char errbuf[200]          = {'\0'};
    char *node                = NULL;
    toml_table_t *conf        = NULL;
    toml_array_t *arrStatCfg  = NULL;
    MetricsConfig *statCfg    = NULL;

    fp = fopen(cfg, "r");
    if (fp == NULL) {
        usys_log_error("cannot open %s", cfg);
        return RETURN_NOTOK;
    }

    conf = toml_parse_file(fp, errbuf, sizeof(errbuf));
    fclose(fp);

    if (conf == NULL) {
        usys_log_error("cannot parse %s", errbuf);
        return RETURN_NOTOK;
    }

    *version = read_str_value(conf, TAG_VERSION);
    if (*version != NULL) {
        usys_log_trace("config version %s", *version);
    }

    node = read_str_value(conf, TAG_NODE);
    if (node != NULL) {
        usys_log_trace("node is %s", node);
        snprintf(gNode, sizeof(gNode), "%s", node);
        free_str_value(node);
    }

    *scrapingTimePeriod = read_int_value(conf, TAG_SCRAPING_TIME_PERIOD);
    usys_log_trace("scraping time period is %d", *scrapingTimePeriod);

    arrStatCfg = toml_array_in(conf, TAG_TABLE);
    if (arrStatCfg == NULL) {
        usys_log_error("missing %s", TAG_TABLE);
        toml_free(conf);
        return RETURN_NOTOK;
    }

    categoryCounts = toml_array_nelem(arrStatCfg);
    usys_log_trace("stat config contains %d categories", categoryCounts);
    *catCount = categoryCounts;

    *pstatCfg = alloc_stat_cfg(categoryCounts);
    if (*pstatCfg == NULL) {
        usys_log_error("failed to allocate %d stat categories",
                       categoryCounts);
        toml_free(conf);
        return RETURN_NOTOK;
    }

    statCfg = *pstatCfg;

    for (idx = 0; idx < categoryCounts; idx++) {
        const char *key              = NULL;
        toml_table_t *statTab        = toml_table_at(arrStatCfg, idx);
        toml_array_t *statCatArr     = NULL;

        if (statTab == NULL) {
            usys_log_error("failed to read stats category at %d", idx);
            toml_free(conf);
            return RETURN_NOTOK;
        }

        key = toml_key_in(statTab, 0);
        usys_log_trace("%d in stat category array is %s", idx, key);

        statCfg[idx].name = alloc_str(strlen(key) + 1);
        if (statCfg[idx].name != NULL) {
            snprintf(statCfg[idx].name, strlen(key) + 1, "%s", key);
        }

        statCatArr = toml_array_in(statTab, key);
        if (statCatArr == NULL) {
            usys_log_error("failed to read stats category %s at %d",
                           key, idx);
            toml_free(conf);
            return RETURN_NOTOK;
        }

        if (toml_parse_stat_cat_array(&statCfg[idx], statCatArr) !=
            RETURN_OK) {
            usys_log_error("failed to parse stats category %s at %d",
                           key, idx);
            toml_free(conf);
            return RETURN_NOTOK;
        }
    }

    usys_log_trace("completed parsing for %s", cfg);

    toml_free(conf);

    return RETURN_OK;
}
