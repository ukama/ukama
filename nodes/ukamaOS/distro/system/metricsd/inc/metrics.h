/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_METRICS_HEADER_H_
#define INC_METRICS_HEADER_H_

#include "prom.h"

#include <string.h>
#include <unistd.h>
#include <stdbool.h>

#include "usys_types.h"
#include "usys_services.h"

#define SERVICE_NAME SERVICE_METRICS

#define RETURN_OK                       0
#define RETURN_NOTOK                    (-1)

/* PROM Metric Type */
typedef union {
    prom_counter_t *counter;
    prom_gauge_t *gauge;
    prom_histogram_t *histogram;
} PromRegistry;

typedef prom_histogram_buckets_t histogram_buckets;

/* METRICS TYPE */
#define METRICTYPE_COUNTER              0
#define METRICTYPE_GAUGE                1
#define METRICTYPE_HISTOGRAM            2

typedef struct {
    char *name;
    char *fqname;
    char *ext;
    int type;
    char *desc;
    char *unit;
    int numLabels;
    char **labels;
    PromRegistry *registry;
    histogram_buckets *buckets;
} KPIConfig;

typedef struct {
    KPIConfig *kpi;
    double value;
} KPIData;

typedef struct {
    char *source;
    char *agent;
    char *url;
    int instances;
    int *range;
    int kpiCount;
    KPIConfig *kpi;
} MetricsCatConfig;

typedef struct {
    char *name;
    MetricsCatConfig *metricsCategory;
    int eachCategoryCount;
} MetricsConfig;

typedef struct {
    char source[32];
} Sources;

typedef int (*metricAddFunc)(KPIConfig *kpi, void *value);

#endif /* INC_METRICS_HEADER_H_ */
