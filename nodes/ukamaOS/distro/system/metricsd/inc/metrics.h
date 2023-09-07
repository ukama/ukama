/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_METRICS_HEADER_H_
#define INC_METRICS_HEADER_H_

#include "prom.h"

#include <string.h>
#include <unistd.h>
#include <stdbool.h>

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
