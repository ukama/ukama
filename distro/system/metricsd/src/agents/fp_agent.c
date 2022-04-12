/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "agents.h"
#include "collector.h"
#include "file.h"
#include "log.h"
#include "metrics.h"

#include <ctype.h>
#include <string.h>

/* Check if file exist */
int fp_check_for_kpi_source(char *source) {
  int ret = RETURN_OK;
  if (!file_exist(source)) {
    ret = RETURN_NOTOK;
  }
  return ret;
}

/* Check if string s2 is prefix + space in string s1
 * on failure it return 0 else a 1.
 */
int fp_is_prefix(char *s1, char *s2) {
  int ret = 1;
  size_t n1 = strlen(s1);
  size_t n2 = strlen(s2);
  if (n1 > n2) {
    for (unsigned int i = 0; i < n2; i++) {
      if (tolower(s1[i]) != tolower(s2[i])) {
        ret = 0;
        break;
      }
    }
    /* check for space */
    if ((ret == 1) && (s1[n2] != ' ')) {
      ret = 0;
    }
  } else {
    ret = 0;
  }
  return ret;
}

/* Parse KPI file */
KPIData *fp_parse_kpi(KPIConfig *kpi, int count, char *kpi_data) {
  KPIData *kdata = NULL;
  for (unsigned int id = 0; id < count; id++) {

    /* Check which  KPI metric */
    if (fp_is_prefix(kpi_data, kpi[id].name)) {
      log_trace("Metrics:: Match found with KPI %s\n", kpi[id].name);
      kdata = calloc(1, sizeof(KPIData));
      if (kdata) {
        kdata->kpi = &kpi[id];
        int dataoffset = strlen(kpi[id].name);
        kdata->value = atof(&kpi_data[dataoffset]);
      }
      break;
    }
  }
  return kdata;
}

/* Read KPI file */
int fp_read_kpi_from_file(MetricsCatConfig *stat, metricAddFunc addFunc) {
  int ret = RETURN_NOTOK;
  FILE *fp;
  char *line = NULL;
  size_t len = 0;
  size_t read;

  /* url is the path of the file */
  if (fp_check_for_kpi_source(stat->url) != RETURN_OK) {
    log_error("Metrics:: Error:: File %s doesn't exist.\n", stat->url);
    return RETURN_NOTOK;
  }

  /* Open KPI file */
  fp = fopen(stat->url, "r");
  if (fp == NULL) {
    log_error("Metrics:: Error:: File %s doesn't exist.\n", stat->url);
    return RETURN_NOTOK;
  }

  /* Read KPI entries */
  while ((read = getline(&line, &len, fp)) != -1) {

    log_trace(" Metrics::  Retrieved line of length %zu: Data %s\n", read,
              line);

    /* Parse KPI data */
    KPIData *kdata = fp_parse_kpi(stat->kpi, stat->kpiCount, line);
    if (kdata) {

      /* Add metric  data for prometheus to scrape */
      ret = addFunc(kdata->kpi, &kdata->value);
      if (ret) {
        log_error(" Metrics:: Failed to add KPI for %s.", kdata->kpi->fqname);
        ret = RETURN_NOTOK;
      } else {
        log_trace(" Metrics:: Added KPI For %s Value %lf.", kdata->kpi->fqname,
                  kdata->value);
        ret = RETURN_OK;
      }
      /* clean */
      free(kdata);
    }
  }
  if (fp) {
    fclose(fp);
  }

  if (line) {
    free(line);
  }

  return ret;
}
