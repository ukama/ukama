/* Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "log.h"
#include "metrics.h"
#include "sys_stat.h"

#include <sys/statvfs.h>

/* Translate bytes to MB */
static double translate_mem_value(unsigned long long mem) {
  return ((double)(mem)) / (1024 * 1024);
}

/* Verify if the path exist.*/
int verify_path(const char *path) {
  struct stat fstat;
  int ret = RETURN_NOTOK;

  if (stat(path, &fstat) == -1) {
    return ret;
  }

  /* Check path is a directory if it isn't,
   * then it can't be a mountpoint. */
  if (!(fstat.st_mode & S_IFDIR)) {
    log_error("Metrics:: %s is not a directory.\n", path);
    return ret;
  }

  return RETURN_OK;
}

/* Read storage stats */
long sys_read_storage_stats(const char *path, SysStorageMetrics *sysSt) {
  int ret = RETURN_NOTOK;
  struct statvfs stat;

  /* If path is not provided. */
  if (!path) {
    log_error("Metrics: No storage device path specified.");
    return ret;
  }

  /* Verify the input path */
  if (verify_path(path) != RETURN_OK) {
    return ret;
  }

  /* Read the stat for path */
  if (statvfs(path, &stat) != 0) {
    return ret;
  }

  /*  Verify using df -B4K */
  sysSt->blksize = stat.f_bsize;
  sysSt->total = stat.f_bsize * stat.f_blocks;
  sysSt->free = stat.f_bsize * stat.f_bavail;
  sysSt->pfree = stat.f_bsize * stat.f_bfree;
  sysSt->used = sysSt->total - sysSt->pfree;

  return RETURN_OK;
}

/* Collect and add storage stats */
int sys_storage_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                           SysStorageMetrics *storageStat,
                                           metricAddFunc addFunc) {
  int ret = RETURN_OK;
  unsigned long long val = 0;

  /* Start Collecting KPI */
  for (int idx = 0; idx < (cfgStat->kpiCount); idx++) {
    KPIConfig *kpi = &(cfgStat->kpi[idx]);
    if ((kpi) && (kpi->fqname)) {

      if (strstr(kpi->fqname, "total")) {
        val = storageStat->total;
      } else if (strstr(kpi->fqname, "used")) {
        val = storageStat->used;
      } else if (strstr(kpi->fqname, "free")) {
        val = storageStat->free;
      } else {
        continue;
      }

      /* Add KPI to server*/
      double castVal = translate_mem_value(val);
      addFunc(kpi, &castVal);
    }
  }
  return ret;
}

/* Collect storage stats */
int sys_storage_collect_stat(MetricsCatConfig *cfgStat,
                             metricAddFunc addFunc) {
  int ret = RETURN_OK;

  SysStorageMetrics *storageStat = calloc(1, sizeof(SysStorageMetrics));
  if (!storageStat) {
    log_error(
        "Metrics:: Failed to allocate memory for storage stat collection.");
    return RETURN_NOTOK;
  }

  if (sys_read_storage_stats(cfgStat->url, storageStat) != RETURN_OK) {
    log_error("Metrics:: Failed to collect storage stats.");
    free(storageStat);
    return RETURN_NOTOK;
  } else if (sys_storage_push_stat_to_metric_server(cfgStat, storageStat,
                                                    addFunc) != RETURN_OK) {
    log_error("Metrics:: Failed to add storage stats to metric server.");
    ret = RETURN_NOTOK;
  }

  free(storageStat);
  return ret;
}
