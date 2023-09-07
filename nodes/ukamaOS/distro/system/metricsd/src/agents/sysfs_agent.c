/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "file.h"
#include "log.h"
#include "metrics.h"

#include <stdio.h>
#include <stdlib.h>

#define PS_MAX_LENGTH_NUMBER 20
#define PS_DEF_OFFSET 0

/* Raw read from sysfs file.*/
int ps_read_block(char *name, void *buff, uint16_t size) {
  int read_bytes = 0;
  int fd = file_open(name, O_RDONLY);
  if (fd < 0) {
    read_bytes = -1;
  } else {
    lseek(fd, PS_DEF_OFFSET, SEEK_SET);
    read_bytes = read(fd, buff, size);
    file_close(fd);
  }
  log_trace("Metrics:: Read %d bytes from %s file from offset 0x%x.",
            read_bytes, name, PS_DEF_OFFSET);
  return read_bytes;
}

/* Check if sysfs file exist */
int sysfs_check_for_kpi_source(char *source) {
  int ret = RETURN_OK;
  if (!file_exist(source)) {
    ret = RETURN_NOTOK;
  }
  return ret;
}

/* Read KPI data for the sysfs  file */
int sysfs_read_kpi_data(char *source, double *nval) {

  int ret = RETURN_OK;
  FILE *fp;
  char line[32];

  if ((fp = fopen(source, "r")) == NULL) {
    log_error("Metrics:: Cannot open %s: %s\n", source);
    return RETURN_NOTOK;
  }

  if (fgets(line, sizeof(line), fp) != NULL) {
    sscanf(line, "%lf", nval);
  }

  if (fp) {
    fclose(fp);
  }

  return ret;
}

int sysfs_push_kpi_metric_server(KPIConfig *kpi, char *source,
                                 metricAddFunc addFunc) {
  int ret = RETURN_NOTOK;

  double val = 0;

  /* Check for source */
  if (sysfs_check_for_kpi_source(source) != RETURN_OK) {
    log_error("Metrics:: Source %s missing for KPI %s", source, kpi->name);
    return ret;
  }

  /* Read KPI data */
  if (sysfs_read_kpi_data(source, &val) != RETURN_OK) {
    log_error("Metrics:: Failed to read KPI %s from file %s ", kpi->name,
              source);
    return ret;
  }

  /* Push Metrics */
  addFunc(kpi, &val);

  return RETURN_OK;
}

/* Collect KPI data from sysfs files */
int sysfs_collect_kpi(MetricsCatConfig *stat, metricAddFunc addFunc) {

  int ret = RETURN_NOTOK;
  for (int idx = 0; idx < stat->kpiCount; idx++) {
    int length =
        sizeof(char) * ((strlen(stat->url)) + (strlen(stat->kpi[idx].ext)));
    char *source = calloc(1, length);
    if (source) {
      strcpy(source, stat->url);
      strcat(source, stat->kpi[idx].ext);
      if (sysfs_push_kpi_metric_server(&(stat->kpi[idx]), source, addFunc) !=
          RETURN_OK) {
        /* failed to push KPI but anyways continue to next KPI */
        log_error("Metrics:: Failed to push data for kpi %s from source %s",
                  stat->kpi[idx].name, source);
      }
      free(source);
      ret = RETURN_OK;
    }
  }
  return ret;
}
