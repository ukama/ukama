/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include "metrics.h"
#include "sys_stat.h"

#include "usys_log.h"

/* Read system uptime */
int sys_read_uptime(SysGenMetrics *sysGen) {
  int ret = RETURN_OK;
  FILE *fp = NULL;
  char line[128];
  double uptime = 0;

  if ((fp = fopen(PROC_UPTIME, "r")) == NULL) {
    return RETURN_NOTOK;
  } else if (fgets(line, sizeof(line), fp)) {
    sscanf(line, "%lf", &uptime);
    sysGen->uptime = uptime;
  }

  if (fp != NULL) {
    fclose(fp);
  }
  return ret;
}

/* Collect and add  stats */
int sys_generic_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                           SysGenMetrics *genStat,
                                           metricAddFunc addFunc) {
  int ret = RETURN_OK;
  double val = 0;

  /* Start Collecting KPI */
  for (int idx = 0; idx < (cfgStat->kpiCount); idx++) {
    KPIConfig *kpi = &(cfgStat->kpi[idx]);
    if ((kpi) && (kpi->fqname)) {

      if (strstr(kpi->fqname, "uptime")) {
        val = genStat->uptime;
      } else {
        continue;
      }
      /* Add KPI to server*/
      addFunc(kpi, &val);
    }
  }
  return ret;
}

/* Collect generic stats */
int sys_gen_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {
  int ret = RETURN_OK;

  SysGenMetrics *genStat = calloc(1, sizeof(SysGenMetrics));
  if (!genStat) {
    usys_log_error("Failed to allocate memory for generic stat collection.");
    return RETURN_NOTOK;
  }

  if (sys_read_uptime(genStat) != RETURN_OK) {
    usys_log_error("Failed to collect generic stats.");
    free(genStat);
    return RETURN_NOTOK;
  } else if (sys_generic_push_stat_to_metric_server(cfgStat, genStat,
                                                    addFunc) != RETURN_OK) {
    usys_log_error("Failed to add generic stats to metric server.");
    ret = RETURN_NOTOK;
  }

  free(genStat);
  return ret;
}
