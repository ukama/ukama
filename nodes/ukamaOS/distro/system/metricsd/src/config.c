/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "config.h"
#include "log.h"
#include "server.h"
#include "toml.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char g_node[32] = {'\0'};

/* Allocate string */
char *alloc_str(int size) {
  char *str = calloc(1, (sizeof(char) * size));
  if (!str) {
    log_error("Metrics:: Memory allocation failed for string of size %d.",
              size);
  }
  return str;
}

/* Allocate memory and read a string from config */
char *read_str_value(toml_table_t *tab, char *key) {
  char *value = NULL;

  toml_datum_t str = toml_string_in(tab, key);
  if (!str.ok) {
    log_error("Metrics:: Failed to read string value for key %s", key);
    return value;
  } else {
    int len = strlen(str.u.s);
    log_trace("Metrics:: Type: String Key: %s Value: %s Length: %d", key,
              str.u.s, len);
    value = calloc(len + 1, sizeof(char));
    if (value) {
      memcpy(value, str.u.s, len);
    }
    free(str.u.s);
  }
  return value;
}

/* Read int value */
int read_int_value(toml_table_t *tab, char *key) {
  toml_datum_t val = toml_int_in(tab, key);
  if (!val.ok) {
    log_error("Metrics:: Failed to read integer value for key %s", key);
    return 0;
  }
  log_trace("Metrics:: Type: Integer Key: %s Value: %d", key, val.u.i);
  return val.u.i;
}

/* Free string value */
void free_str_value(char *str) {
  if (str) {
    free(str);
  }
  str = NULL;
}

/* Allocate labels */
char **alloc_lables(int count) {
  char **labels = (char **)calloc(count, sizeof(char *));
  if (!labels) {
    labels = NULL;
  }
  return labels;
}

/* Free labels */
void free_labels(char **labels, int count) {
  for (int idx = 0; idx < count; idx++) {
    if (labels[idx]) {
      free(labels[idx]);
      labels[idx] = NULL;
    }
  }
  free(labels);
  labels = NULL;
}

/* Allocate KPI */
KPIConfig *alloc_kpi(int count) {
  KPIConfig *kpi = calloc(count, sizeof(KPIConfig));
  if (!kpi) {
    return NULL;
  }
  return kpi;
}

/* Free KPI */
void free_kpi(KPIConfig *kpi, int count) {
  if (kpi) {
    for (int idx = 0; idx < count; idx++) {
      free_str_value(kpi[idx].name);
      free_str_value(kpi[idx].fqname);
      free_str_value(kpi[idx].ext);
      free_str_value(kpi[idx].desc);
      free_str_value(kpi[idx].unit);
      free_labels(kpi[idx].labels, kpi[idx].numLabels);
    }
    free(kpi);
    kpi = NULL;
  }
}

/* Allocate range */
int *alloc_range(int count) {
  int *range = calloc(count, sizeof(int));
  if (!range) {
    return NULL;
  }
  return range;
}

/* Free range */
void free_range(int *range) {
  if (range) {
    free(range);
    range = NULL;
  }
}

/* Allocate stat */
MetricsCatConfig *alloc_stat(int count) {
  MetricsCatConfig *stat = calloc(count, sizeof(MetricsCatConfig));
  if (!stat) {
    return NULL;
  }
  return stat;
}

/* Free stat */
void free_stat(MetricsCatConfig *stat, int count) {
  if (stat) {
    for (int idx = 0; idx < count; idx++) {
      free_str_value(stat[idx].source);
      free_str_value(stat[idx].agent);
      free_str_value(stat[idx].url);
      free_range(stat[idx].range);
      int kcount = (stat->instances) ? (stat->instances) : 1;
      free_kpi(stat[idx].kpi, stat[idx].kpiCount * kcount);
    }
    free(stat);
    stat = NULL;
  }
}

/* Allocate stat table */
MetricsConfig *alloc_stat_cfg(int count) {
  MetricsConfig *stat_cfg = calloc(count, sizeof(MetricsConfig));
  if (!stat_cfg) {
    return NULL;
  }
  return stat_cfg;
}

/* Free stat table*/
void free_stat_cfg(MetricsConfig *stat_cfg, int count) {
  if (stat_cfg) {
    for (int idx = 0; idx < count; idx++) {
      free_str_value(stat_cfg[idx].name);
      free_stat(stat_cfg[idx].metricsCategory, stat_cfg[idx].eachCategoryCount);
    }
    free(stat_cfg);
    stat_cfg = NULL;
  }
}

/* convert to lower_case */
void lower_string(char **string) {
  int c = 0;
  if (!(*string)) {
      return;
  }
  char *s = *string;
  while (s[c] != '\0') {
    if (s[c] >= 'A' && s[c] <= 'Z') {
      s[c] = s[c] + 32;
    }
    c++;
  }
}

/* return fully qualified kpi name */
char *set_fqkname(char *node, char *category, char *source, int range,
                  char *kpi, char* unit) {
  char *name = NULL;
  char fqkn[MAX_KPI_KEY_NAME_LENGTH] = {'\0'};

  /* Make sure every thing is lower case */
  lower_string(&node);
  lower_string(&category);
  lower_string(&source);
  lower_string(&kpi);
  lower_string(&unit);
  int len = 0;

  /* Extended KPI name with units */
  char extKpiName[MAX_KPI_KEY_NAME_LENGTH] = {'\0'};
  if (unit) {
      snprintf(extKpiName, MAX_KPI_KEY_NAME_LENGTH, "%s%s%s", kpi, TAG_SEP, unit);
  } else {
      strcpy(extKpiName, kpi);
  }

  if (range < 0) {
    len = snprintf(fqkn, MAX_KPI_KEY_NAME_LENGTH, "%s%s%s%s%s%s%s", g_node,
                   TAG_SEP, category, TAG_SEP, source, TAG_SEP, extKpiName);
  } else {
    len = snprintf(fqkn, MAX_KPI_KEY_NAME_LENGTH, "%s%s%s%s%s%d%s%s", g_node,
                   TAG_SEP, category, TAG_SEP, source, range, TAG_SEP, extKpiName);
  }

  name = alloc_str(len + 1);
  if (name) {
    memcpy(name, fqkn, strlen(fqkn));
  }

  return name;
}

/* Parsing Metric Type */
int read_metric_type(toml_table_t *tab_kpi) {
  int metrictype = METRICTYPE_GAUGE;
  char *type = read_str_value(tab_kpi, TAG_METRIC_TYPE);
  if (type) {
    if (strcmp(type, "METRICTYPE_COUNTER") == 0) {
      metrictype = METRICTYPE_COUNTER;
    } else if (strcmp(type, "METRICTYPE_GAUGE") == 0) {
      metrictype = METRICTYPE_GAUGE;
    } else if (strcmp(type, "METRICTYPE_HISTOGRAM") == 0) {
      metrictype = METRICTYPE_HISTOGRAM;
    }
  } else {
    /* Set default GUAGE TYPE */
    log_error("Metrics:: Err reading metric type.");
  }
  return metrictype;
}

/* Parse range array */
int *parse_range_array(int *inst, toml_table_t *tab, char *key) {

  int *range = NULL;
  /* Range array */
  toml_array_t *arr_range = toml_array_in(tab, key);
  if (!arr_range) {
    log_error("Metrics:: Error missing %s", key);
  }

  *inst = toml_array_nelem(arr_range);
  log_trace("Metrics:: %d range values available.", *inst);

  /* Allocate memory for kpi */
  range = alloc_range(*inst);
  if (!range) {
    log_error(
        "Metrics:: Error: failed to allocate memory for %d range elements.",
        *inst);
    return range;
  }

  /* Parse each value in range */
  for (uint8_t idx = 0; idx < *inst; idx++) {
    toml_datum_t data = toml_int_at(arr_range, idx);
    if (!data.ok) {
      free(range);
      range = NULL;

    } else {
      range[idx] = data.u.i;
      log_trace("Metrics:: Range value at [%d] is %d", idx, range[idx]);
    }
  }

  return range;
}

/* Parse int value in range */
int *parse_range_int(toml_table_t *tab, char *key) {

  int *range = NULL;

  /* Allocate memory for range */
  range = alloc_range(1);
  if (!range) {
    log_error("Metrics:: Error: failed to allocate memory for range value.");
    return range;
  }

  /* Parse each value in range */

  toml_datum_t data = toml_int_in(tab, key);
  if (!data.ok) {
    free(range);
    range = NULL;

  } else {
    int idx = 0;
    range[idx] = data.u.i;
    log_trace("Metrics:: Range value at [%d] is %d", idx, range[idx]);
  }

  return range;
}

/* Parse range */
int *parse_range(int *inst, toml_table_t *tab, char *key) {
  /* Range array */
  toml_array_t *arr_range = toml_array_in(tab, key);
  if (!arr_range) {
    /* Disable for now
     * Need to see if need is there or not
     */
    //*inst = 0;
    // return parse_range_int(tab, key);
    return NULL;
  } else {
    return parse_range_array(inst, tab, key);
  }
}

/* Parsing Metric Type */
char **read_labels(int *count, toml_table_t *tab_kpi) {
  char **labels = NULL;
  int labels_count = 0;

  /* KPI Table array */
  toml_array_t *arr_kpi = toml_array_in(tab_kpi, TAG_LABELS);
  if (!arr_kpi) {
    log_error("Metrics:: Error missing %s", TAG_LABELS);
    return NULL;
  }

  *count = toml_array_nelem(arr_kpi);
  log_trace("Metrics:: %d labels available.", *count);

  if (*count <= 0) {
    return NULL;
  }

  /* Allocate memory for kpi */
  labels = alloc_lables(*count);
  if (!(labels)) {
    log_error("Metrics:: Error: failed to allocate memory for %d labels.",
              *count);
    return NULL;
  }

  /* Parse each KPI */
  for (uint8_t idx = 0; idx < *count; idx++) {
    toml_datum_t dlabel = toml_string_at(arr_kpi, idx);
    if (!dlabel.ok) {
      /* If we read less than actual count */
      *count = labels_count;
      break;
    } else {
      int len = strlen(dlabel.u.s);
      labels[idx] = calloc(len + 1, sizeof(char));
      memcpy(labels[idx], dlabel.u.s, len + 1);
      labels_count++;
      free(dlabel.u.s);
    }
  }

  return labels;
}

/* Parsing KPI */
static int toml_parse_kpi_table(char *category, char *source, int range,
                                KPIConfig *kpi, toml_table_t *tab_kpi) {
  int ret = RETURN_NOTOK;

  /* Name */
  kpi->name = read_str_value(tab_kpi, TAG_NAME);

  /* EXT */
  kpi->ext = read_str_value(tab_kpi, TAG_EXT);

  /* Description */
  kpi->desc = read_str_value(tab_kpi, TAG_DESC);

  /* Unit */
  kpi->unit = read_str_value(tab_kpi, TAG_UNIT);

  /* Type */
  kpi->type = read_metric_type(tab_kpi);

  /* Labels */
  kpi->labels = read_labels(&kpi->numLabels, tab_kpi);

  /* Get FQKN */
  kpi->fqname = set_fqkname(g_node, category, source, range, kpi->name, kpi->unit);

  /* register to prometheus metrics */
  ret = metric_server_register_kpi(kpi);

  return ret;
}

/* Parsing stat table */
static int toml_parse_stat_table(char *category, MetricsCatConfig *stat,
                                 toml_table_t *tab_stat) {
  int ret = RETURN_OK;

  int clmn = toml_table_nkval(tab_stat);
  log_trace("Metrics:: Number of columns in Stat table %d", clmn);

  /* Source */
  stat->source = read_str_value(tab_stat, TAG_SOURCE);

  /* Agent */
  stat->agent = read_str_value(tab_stat, TAG_AGENT);

  /* URL */
  stat->url = read_str_value(tab_stat, TAG_URL);

  /* Range */
  stat->range = parse_range(&(stat->instances), tab_stat, TAG_RANGE);

  /* Table source */
  toml_table_t *tab_source = toml_table_in(tab_stat, stat->source);
  if (!tab_source) {
    log_error("Metrics:: Error missing table %s", stat->source);
    return RETURN_NOTOK;
  }

  /* KPI Table array */
  toml_array_t *arr_kpi = toml_array_in(tab_source, TAG_KPI);
  if (!arr_kpi) {
    log_error("Metrics:: Error missing %s for source.", TAG_KPI, stat->source);
    return RETURN_NOTOK;
  }

  stat->kpiCount = toml_array_nelem(arr_kpi);
  log_trace("Metrics:: %d KPI available.", stat->kpiCount);

  int rkpi_count = stat->kpiCount;

  /* Handle range:  Range means KPI has to be repeated for each value in range +
   * master kpi example like for cpu range is 1,2,3,4 kpi would be available for
   * cpu, cpu0, cpu1, cpu2, cpu3
   *   */
  if (stat->instances) {
    if (stat->range) {
      rkpi_count = stat->kpiCount * (stat->instances);
    }
  }

  /* Allocate memory for kpi */
  stat->kpi = alloc_kpi(rkpi_count);
  if (!(stat->kpi)) {
    log_error("Metrics:: Error: failed to allocate memory for %d KPI.",
              rkpi_count);
    return RETURN_NOTOK;
  }

  int max_instances = (stat->instances) ? (stat->instances) : 1;

  log_trace("Metrics:: Range has %d instances max instances is set to %d for "
            "category %s source %s.",
            stat->instances, max_instances, category, stat->source);

  /* For each value in range */
  for (uint8_t inst = 0; inst < max_instances; inst++) {

    /* Parse each KPI */
    for (uint8_t idx = 0; idx < stat->kpiCount; idx++) {

      toml_table_t *tab_kpi = toml_table_at(arr_kpi, idx);
      if (!tab_kpi) {
        return RETURN_NOTOK;
      }

      if (stat->range) {

        /* For Range values */
        int rval = stat->range[inst];
        if (RETURN_OK !=
            toml_parse_kpi_table(category, stat->source, rval,
                                 &(stat->kpi[idx + (inst * stat->kpiCount)]),
                                 tab_kpi)) {
          return RETURN_NOTOK;
        }

      } else {

        /* In case if no range is available */
        int rval = -1;
        if (RETURN_OK !=
            toml_parse_kpi_table(category, stat->source, rval,
                                 &(stat->kpi[idx + (inst * stat->kpiCount)]),
                                 tab_kpi)) {
          return RETURN_NOTOK;
        }
      }
    }
  }

  return ret;
}

/* Parsing stat table */
static int toml_parse_stats_cat(char *category, MetricsCatConfig *stat,
                                toml_table_t *tab_stats) {

  /* Get table key */
  const char *key = toml_key_in(tab_stats, 0);
  log_trace("Metrics::Key in stat category table is %s.", key);

  /* parse table */
  if (RETURN_OK != toml_parse_stat_table(category, stat, tab_stats)) {
    return RETURN_NOTOK;
  }

  return RETURN_OK;
}

/* Parsing stat category array */
static int toml_parse_stat_cat_array(MetricsConfig *stat_cfg,
                                     toml_array_t *stat_cat_arr) {

  /* Device under each category count */
  int table_count = toml_array_nelem(stat_cat_arr);
  log_trace("Metrics:: Stats category contains %d tables.", table_count);
  stat_cfg->eachCategoryCount = table_count;

  /* Allocate memory */
  stat_cfg->metricsCategory = alloc_stat(table_count);
  if (!(stat_cfg->metricsCategory)) {
    log_error("Metrics:: Error: failed to allocate memory for %d tables to "
              "stat category.",
              table_count);
    return RETURN_NOTOK;
  }

  MetricsCatConfig *stats = stat_cfg->metricsCategory;

  for (uint8_t idx = 0; idx < table_count; idx++) {

    toml_table_t *tab_stat = toml_table_at(stat_cat_arr, idx);
    if (!tab_stat) {
      log_error(
          "Metrics:: Error: reading Stats table at %d for category for %s.",
          idx, stat_cfg->name);
      return RETURN_NOTOK;
    }

    if (RETURN_OK !=
        toml_parse_stats_cat(stat_cfg->name, &stats[idx], tab_stat)) {
      log_error(
          "Metrics:: Error: Parsing stats table for idx %d category %s failed.",
          idx, stat_cfg->name);
      return RETURN_NOTOK;
    }
  }

  return RETURN_OK;
}

int toml_parse_config(char *cfg, char **version, int *scraping_time_period,
                      int *server_port, MetricsConfig **pstat_cfg,
                      int *cat_count) {
  FILE *fp;
  char errbuf[200];
  char *node;

  // 1. Read and parse toml file
  fp = fopen(cfg, "r");
  if (!fp) {
    log_error("Metrics:: cannot open %s.", cfg);
  }

  toml_table_t *conf = toml_parse_file(fp, errbuf, sizeof(errbuf));
  fclose(fp);

  if (!conf) {
    log_error("Metrics:: cannot parse  %s ", errbuf);
  }

  /* Read version */
  *version = read_str_value(conf, TAG_VERSION);
  if (*version) {
    log_trace("Metrics:: config version %s.", *version);
  }

  /* Node type */
  node = read_str_value(conf, TAG_NODE);
  if (node) {
    log_trace("Metrics:: Node read is %s.", node);
    memcpy(g_node, node, strlen(node));
  }

  *scraping_time_period = read_int_value(conf, TAG_SCRAPING_TIME_PERIOD);
  log_trace("Metrics:: Scraping time period read is %d.",
            *scraping_time_period);

  *server_port = read_int_value(conf, TAG_SERVER_PORT);
  log_trace("Metrics:: Metric server port %d.", *server_port);

#if 0
    /* Parse Source list */
    toml_array_t* arr_src_list = toml_array_in(conf, TAG_SOURCE_LIST);
    if (!arr_src_list) {
    	log_error("missing %s", TAG_SOURCE_LIST);
    }
    if( RETURN_OK != toml_parse_source_list(arr_src_list)) {
    	return RETURN_NOTOK;
    }
#endif
  /* Traverse ARRAY of table of stats category. */
  toml_array_t *arr_stat_cfg = toml_array_in(conf, TAG_TABLE);
  if (!arr_stat_cfg) {
    log_error("missing %s", TAG_TABLE);
  }

  /* Stats config category  count */
  int category_counts = toml_array_nelem(arr_stat_cfg);
  log_trace("Metrics:: Stat config contains %d category.", category_counts);
  *cat_count = category_counts;

  /* Allocate memory */
  *pstat_cfg = alloc_stat_cfg(category_counts);
  if (!(*pstat_cfg)) {
    log_error(
        "Metrics:: Error: failed to allocate memory for %d category of stats.",
        category_counts);
    return RETURN_NOTOK;
  }
  MetricsConfig *stat_cfg = *pstat_cfg;

  for (uint8_t idx = 0; idx < category_counts; idx++) {

    toml_table_t *stat_tab = toml_table_at(arr_stat_cfg, idx);
    if (!stat_tab) {
      log_error("Metrics:: Err reading Stats category table from array at %d.",
                idx);
      return RETURN_NOTOK;
    }

    const char *key = toml_key_in(stat_tab, 0);
    log_trace("Metrics:: %d in stat category array is %s.", idx, key);

    /* Copy category name of stats */
    stat_cfg[idx].name = alloc_str(strlen(key) + 1);
    if (stat_cfg[idx].name) {
      memcpy(stat_cfg[idx].name, key, strlen(key));
    }

    toml_array_t *stat_cat_arr = toml_array_in(stat_tab, key);
    if (!stat_cat_arr) {
      log_error("Metrics:: Err reading Stats category %s array at %d.", key,
                idx);
      return RETURN_NOTOK;
    }

    if (RETURN_OK != toml_parse_stat_cat_array(&stat_cfg[idx], stat_cat_arr)) {
      log_error(
          "Metrics:: Error: Parsing stats category array %s for idx %d failed.",
          key, idx);
      return RETURN_NOTOK;
    }
  }
  log_trace("Metrics:: Parsing completed for %s.", cfg);
  toml_free(conf);
  return RETURN_OK;
}
