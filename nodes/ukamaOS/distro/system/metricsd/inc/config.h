/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef SRC_PARSER_CONFIG_H_
#define SRC_PARSER_CONFIG_H_

#include "metrics.h"

#define MAX_KPI_KEY_NAME_LENGTH         128
#define TAG_SEP                            "_"
#define TAG_VERSION                     "version"
#define TAG_SCRAPING_TIME_PERIOD        "scraping_time_period"
#define TAG_NODE                        "ukamnode"
#define TAG_GENERIC                     "generic"
#define TAG_SYSTEM                      "system"
#define TAG_CAT_SOC                     "SoC"
#define TAG_SOURCE_LIST                 "sources"
#define TAG_SOURCE                      "source"
#define TAG_AGENT                       "agent"
#define TAG_URL                         "url"
#define TAG_RANGE                       "range"
#define TAG_NAME                        "name"
#define TAG_EXT                         "ext"
#define TAG_TABLE                       "stats"
#define TAG_KPI                         "kpi"
#define TAG_IOR_LTE_STACK               "IoR-lte-stack"
#define TAG_IOR_CORE                    "Ior-core"
#define TAG_SENSORS                     "sensors"
#define TAG_CPU                         "cpu"
#define TAG_MEMORY                      "memory"
#define TAG_NETWORK                     "network"
#define TAG_DESC                        "desc"
#define TAG_UNIT                        "unit"
#define TAG_NUM_LABELS                  "numLabels"
#define TAG_METRIC_TYPE                 "type"
#define TAG_LABELS                      "labels"

int toml_parse_config(char *cfg, char **version, int *scraping_time_period,
                      MetricsConfig **pstat_cfg, int *source_count);
void free_stat_cfg(MetricsConfig *stat_cfg, int count);

#endif /* SRC_PARSER_CONFIG_H_ */
