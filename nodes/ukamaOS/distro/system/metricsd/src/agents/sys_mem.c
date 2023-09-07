/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "log.h"
#include "metrics.h"
#include "sys_stat.h"

#define MEM_SOURCE_DDR "ddr"
#define MEM_SOURCE_SWAP "swap"

/* Translate kB to MB */
static double translate_mem_value(unsigned long long mem) {
  return ((double)(mem)) / (1024);
}

/* Read mem stats */
int sys_read_memStat(SysMemMetrics *statMem) {
  int ret = RETURN_OK;
  FILE *fp;
  char line[128];

  if ((fp = fopen(PROC_MEM_STAT, "r")) == NULL) {
    log_error("Metrics:: Cannot open %s: %s\n", PROC_MEM_STAT);
    return RETURN_NOTOK;
  }
  
  while (fgets(line, sizeof(line), fp) != NULL) {

    /* Read the total amount of memory in kB */
    if (!strncmp(line, "MemTotal:", 9)) {
      sscanf(line + 9, "%llu", &statMem->ddr.memTotal);
    } else if (!strncmp(line, "MemFree:", 8)) {
      /* Read the amount of free memory in kB */
      sscanf(line + 8, "%llu", &statMem->ddr.memFree);
    } else if (!strncmp(line, "MemAvailable:", 13)) {
      /* Read the amount of available memory in kB */
      sscanf(line + 13, "%llu", &statMem->ddr.memAvail);
    } else if (!strncmp(line, "Buffers:", 8)) {
      /* Read the amount of buffered memory in kB */
      sscanf(line + 8, "%llu", &statMem->ddr.memBuffer);
    } else if (!strncmp(line, "Cached:", 7)) {
      /* Read the amount of cached memory in kB */
      sscanf(line + 7, "%llu", &statMem->ddr.memCached);
    } else if (!strncmp(line, "SwapTotal:", 10)) {
      /* Read the total amount of swap memory in kB */
      sscanf(line + 10, "%llu", &statMem->swap.total);
    } else if (!strncmp(line, "SwapFree:", 9)) {
      /* Read the amount of free swap memory in kB */
      sscanf(line + 9, "%llu", &statMem->swap.free);
    }
  }

  /* calculate mem used */
  statMem->ddr.memUsed = statMem->ddr.memTotal - statMem->ddr.memFree -
                           statMem->ddr.memBuffer - statMem->ddr.memCached;

  /* calculate swap mem used */
  statMem->swap.used = statMem->swap.total - statMem->swap.free;

  fclose(fp);

  return ret;
}

/* Collect and add ddr mem stats */
int sys_mem_ddr_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                           SysMemDDRMetrics *ddr_stat,
                                           metricAddFunc addFunc) {
  int ret = RETURN_OK;
  unsigned long long val = 0;

  /* Start Collecting KPI */
  for (int idx = 0; idx < (cfgStat->kpiCount); idx++) {
    KPIConfig *kpi = &(cfgStat->kpi[idx]);
    if ((kpi) && (kpi->fqname)) {

      if (strstr(kpi->fqname, "total")) {
        val = ddr_stat->memTotal;
      } else if (strstr(kpi->fqname, "used")) {
        val = ddr_stat->memUsed;
      } else if (strstr(kpi->fqname, "free")) {
        val = ddr_stat->memFree;
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

/* Collect and add swap mem stats */
int sys_mem_swap_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                            SysMemSwapMetrics *swapStat,
                                            metricAddFunc addFunc) {
  int ret = RETURN_OK;
  unsigned long long val = 0;

  /* Start Collecting KPI */
  for (int idx = 0; idx < (cfgStat->kpiCount); idx++) {
    KPIConfig *kpi = &(cfgStat->kpi[idx]);
    if ((kpi) && (kpi->fqname)) {

      if (strstr(kpi->fqname, "total")) {
        val = swapStat->total;
      } else if (strstr(kpi->fqname, "used")) {
        val = swapStat->used;
      } else if (strstr(kpi->fqname, "free")) {
        val = swapStat->free;
      } else {
        continue;
      }

      /* Add KPI to server*/
      double cast_val = translate_mem_value(val);
      addFunc(kpi, &cast_val);
    }
  }
  return ret;
}

/* Collect memory  stats */
int sys_mem_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {
  int ret = RETURN_OK;

  SysMemMetrics *memStat = calloc(1, sizeof(SysMemMetrics));
  if (!memStat) {
    log_error("Metrics:: Failed to allocate memory for mem stat.");
    return RETURN_NOTOK;
  }

  if (sys_read_memStat(memStat) != RETURN_OK) {
    log_error("Metrics:: Failed to read memory stats.");
    free(memStat);
    return RETURN_NOTOK;
  }

  if (!strcmp(cfgStat->source, MEM_SOURCE_DDR)) {
    if (sys_mem_ddr_push_stat_to_metric_server(cfgStat, &memStat->ddr,
                                               addFunc) != RETURN_OK) {
      log_error("Metrics:: Failed to add %s memory stats to metric server.",
                MEM_SOURCE_DDR);
      ret = RETURN_NOTOK;
    }
  } else if (!strcmp(cfgStat->source, MEM_SOURCE_SWAP)) {
    if (sys_mem_swap_push_stat_to_metric_server(cfgStat, &memStat->swap,
                                                addFunc) != RETURN_OK) {
      log_error("Metrics:: Failed to add memory stats %s to metric server.",
                MEM_SOURCE_SWAP);
      ret = RETURN_NOTOK;
    } else {
      log_error("Metrics:: Failed to find memory stats source  %s.",
                cfgStat->source);
      ret = RETURN_NOTOK;
    }
  }

  free(memStat);

  return ret;
}
