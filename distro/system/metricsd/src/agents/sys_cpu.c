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

#define CPU_SOURCE "cpu"
#define CORE_SOURCE "cpu_core"

static SysCPUMetrics pCpuStat;

/* Lookup for KPI in stat cfg */
KPIConfig *kpi_lookup(MetricsCatConfig *cpuCfg, char *name) {

  for (int idx = 0; idx < cpuCfg->kpiCount; idx++) {
    KPIConfig *kpi = &(cpuCfg->kpi[idx]);
    if ((kpi) && (kpi->name)) {
      if (strstr(kpi->fqname, name)) {
        return kpi;
      }
    }
  }

  return NULL;
}

/* Check if present in range */
int check_if_present_in_range(int id, int *range, int count) {
  for (int idx = 0; idx < count; idx++) {
    if (id == range[idx]) {
      return 1;
    }
  }
  return 0;
}

/* Read CPU count */
int sys_cpu_read_count(int hidden) {
  DIR *dir;
  struct dirent *drd;
  struct stat buf;
  char line[MAX_PF_NAME];
  int coreNum, coreCount = -1;

  if ((dir = opendir(SYS_DEV_CPU)) == NULL)
    return 0;

  while ((drd = readdir(dir)) != NULL) {

    if (!strncmp(drd->d_name, "cpu", 3) && isdigit(drd->d_name[3])) {
      snprintf(line, sizeof(line), "%s/%s", SYS_DEV_CPU, drd->d_name);
      line[sizeof(line) - 1] = '\0';
      if (stat(line, &buf) < 0)
        continue;
      if (S_ISDIR(buf.st_mode)) {
        if (hidden) {
          sscanf(drd->d_name + 3, "%d", &coreNum);
          if (coreNum > coreCount) {
            coreCount = coreNum;
          }
        } else {
          coreCount++;
        }
      }
    }
  }
  closedir(dir);

  return (coreCount + 1);
}

/* Read CPU stats from the file */
int sys_cpu_read_stat(SysCPUMetrics *cpu, int cpuMax) {
  int ret = RETURN_OK;
  FILE *fp;
  SysCPUMetrics *stCpuIdx;
  SysCPUMetrics icpu;
  char line[8192];
  int cpuId;

  if ((fp = fopen(PROC_CPU_STAT, "r")) == NULL) {
    log_error("Metrics:: Cannot open %s: %s\n", PROC_CPU_STAT);
    return RETURN_NOTOK;
  }
  memset(cpu, 0, sizeof(SysCPUMetrics) * (cpuMax + 1));
  while (fgets(line, sizeof(line), fp) != NULL) {
    if (!strncmp(line, "cpu ", 4)) {
      /* Aggregated CPU */
      sscanf(line + 5, "%llu %llu %llu %llu %llu %llu %llu %llu %llu %llu",
             &cpu->cpuUser, &cpu->cpuNice, &cpu->cpuSys, &cpu->cpuIdle,
             &cpu->cpuIowait, &cpu->cpuHardirq, &cpu->cpuSoftirq,
             &cpu->cpuSteal, &cpu->cpuGuest, &cpu->cpuGuestNice);
    } else if (!strncmp(line, "cpu", 3)) {
      /* CPUx */
      memset(&icpu, 0, sizeof(SysCPUMetrics));
      sscanf(line + 3, "%d %llu %llu %llu %llu %llu %llu %llu %llu %llu %llu",
             &cpuId, &icpu.cpuUser, &icpu.cpuNice, &icpu.cpuSys,
             &icpu.cpuIdle, &icpu.cpuIowait, &icpu.cpuHardirq,
             &icpu.cpuSoftirq, &icpu.cpuSteal, &icpu.cpuGuest,
             &icpu.cpuGuestNice);

      if (cpuId >= cpuMax) {
        log_error("Metrics:: Read cpu id is %d expected %d or less.", cpuId,
                  cpuMax);
        ret = RETURN_NOTOK;
        break;
      }

      stCpuIdx = cpu + cpuId + 1;
      memcpy(stCpuIdx, &icpu, sizeof(SysCPUMetrics));
    }
  }

  fclose(fp);
  return ret;
}

int sys_cpu_get_kpi_value(KPIConfig *kpi, SysCPUMetrics *cpuStat,
                          unsigned long long int *val) {
  int ret = RETURN_OK;

  /* Choose KPI for the CPU */
  if (strstr(kpi->name, "frequency")) {
    *val = cpuStat->freq;
  } else if (strstr(kpi->name, "user_usage")) {
    *val = cpuStat->cpuUser;
  } else if (strstr(kpi->name, "nice")) {
    *val = cpuStat->cpuNice;
  } else if (strstr(kpi->name, "system_usage")) {
    *val = cpuStat->cpuSys;
  } else if (strstr(kpi->name, "idle_time")) {
    *val = cpuStat->cpuIdle;
  } else if (strstr(kpi->name, "io_wait_time")) {
    *val = cpuStat->cpuIowait;
  } else if (strstr(kpi->name, "hardirq")) {
    *val = cpuStat->cpuHardirq;
  } else if (strstr(kpi->name, "softirq")) {
    *val = cpuStat->cpuSoftirq;
  } else if (strstr(kpi->name, "steal")) {
    *val = cpuStat->cpuSteal;
  } else {
    log_error("Metrics:: KPI %s not available under CPU.", kpi->name);
    ret = RETURN_NOTOK;
  }

  return ret;
}



/* Push cpu usage */
int sys_read_and_push_cpu_usage(MetricsCatConfig *cpuCfg,
                                SysCPUMetrics *cpuStat,
                                metricAddFunc addFunc) {
  int ret = RETURN_OK;

  unsigned long long  idle = cpuStat->cpuIowait + cpuStat->cpuIdle;

  unsigned long long  nonidle = (cpuStat->cpuUser + cpuStat->cpuNice +
          cpuStat->cpuSys + cpuStat->cpuHardirq +
          cpuStat->cpuSoftirq + cpuStat->cpuSteal +
          cpuStat->cpuGuest + cpuStat->cpuGuestNice);

  unsigned long long total = idle +nonidle;

  unsigned long long  acc_cpu = (cpuStat->cpuUser + cpuStat->cpuNice +
          cpuStat->cpuSys + cpuStat->cpuIdle +
          cpuStat->cpuIowait + cpuStat->cpuHardirq +
          cpuStat->cpuSoftirq + cpuStat->cpuSteal +
          cpuStat->cpuGuest + cpuStat->cpuGuestNice);

  double usage = 100 - ((double)(cpuStat->cpuIdle * 100) / acc_cpu);

  log_trace("cpuUser: %llu cpuNice: %llu cpuSys: %llu cpuIdle: %llu cpuIowait: %llu cpuHardirq: %llu cpuSoftirq: %llu cpuSteal: %llu cpuGuest: %llu cpuGuestNice: %llu \n acc_cpu: %llu Usage: %lf",
		     cpuStat->cpuUser, cpuStat->cpuNice, cpuStat->cpuSys,
			 cpuStat->cpuIdle, cpuStat->cpuIowait, cpuStat->cpuHardirq,
			 cpuStat->cpuSoftirq, cpuStat->cpuSteal, cpuStat->cpuGuest,
			 cpuStat->cpuGuestNice, acc_cpu, usage);

  KPIConfig *kpi = kpi_lookup(cpuCfg, "cpu_usage");
  if (kpi) {
    addFunc(kpi, &usage);
  }
  log_trace("Metrics: Pushed cpu usage %lf to metrics server.", usage);

  /* Calculate real time cpu usage */

  unsigned long long  old_idle =  pCpuStat.cpuIdle + pCpuStat.cpuIowait;

  unsigned long long  old_nonidle =  (pCpuStat.cpuUser + pCpuStat.cpuNice +
          pCpuStat.cpuSys + pCpuStat.cpuHardirq +
          pCpuStat.cpuSoftirq + pCpuStat.cpuSteal +
          pCpuStat.cpuGuest + pCpuStat.cpuGuestNice);

  unsigned long long old_total = old_idle + old_nonidle;

  unsigned long long  old_acc_cpu = (pCpuStat.cpuUser + pCpuStat.cpuNice +
          pCpuStat.cpuSys + pCpuStat.cpuIdle +
          pCpuStat.cpuIowait + pCpuStat.cpuHardirq +
          pCpuStat.cpuSoftirq + pCpuStat.cpuSteal +
          pCpuStat.cpuGuest + pCpuStat.cpuGuestNice);

  double realUsage = 100 - ((double)( (cpuStat->cpuIdle - pCpuStat.cpuIdle)  * 100) / (acc_cpu - old_acc_cpu));

  log_trace("Real usage is %lf previous cpuStatIdle %llu previous acc_cpu %llu", realUsage, pCpuStat.cpuIdle, old_acc_cpu);

  unsigned long long dtotal = total - old_total;

  unsigned long long didle = idle - old_idle;

  double cpu_per = ((double)(dtotal - didle)/dtotal)*100;

  memcpy(&pCpuStat,cpuStat, sizeof(SysCPUMetrics));

  kpi = kpi_lookup(cpuCfg, "cpu_realtime_usage");
  if (kpi) {
	  addFunc(kpi, &cpu_per);
  }
  log_trace("Metrics: Pushed realtime cpu usage %lf to metrics server.", cpu_per);

  return ret;
}

/* Push cpu core count */
int sys_read_and_push_cpu_count(MetricsCatConfig *cpuCfg, int core,
                                metricAddFunc addFunc) {
  int ret = RETURN_OK;
  KPIConfig *kpi = kpi_lookup(cpuCfg, "cpu_cores");
  if (kpi) {
    double castVal = (double)core;
    addFunc(kpi, &castVal);
  }
  return ret;
}

/* Push CPU core stats to the Metric server */
int sys_cpu_core_push_stat_to_metric_server(MetricsCatConfig *cpuCfg,
                                            SysCPUMetrics *cstat,
                                            int coreCount,
                                            metricAddFunc addFunc) {

  /* Push Metrics */
  unsigned long long int val = 0;

  /* Start with core 0 */
  for (int cid = 0; cid < coreCount; cid++) {
    char cpuid[32] = {'\0'};
    SysCPUMetrics *cpuStat = NULL;

    /* Only for cpu cores in range  */
    if (check_if_present_in_range(cid, cpuCfg->range, cpuCfg->instances)) {
      sprintf(cpuid, "_%s%d_", CORE_SOURCE, cid);
      cpuStat = &cstat[cid + 1];
    } else {
      continue;
    }

    if (!cpuStat) {
      return RETURN_NOTOK;
    }

    /* Start Collecting KPI */
    for (int idx = 0; idx < ((cpuCfg->kpiCount) * (cpuCfg->instances));
         idx++) {

      KPIConfig *kpi = &(cpuCfg->kpi[idx]);
      if ((kpi) && (kpi->name)) {
        /* Check if KPI is for core<cid> */
        if (!strstr(kpi->fqname, cpuid)) {
          continue;
        }

        if (sys_cpu_get_kpi_value(kpi, cpuStat, &val) != RETURN_OK) {
          continue;
        }

        /* Add KPI to server*/
        double cast_val = (double)val;
        addFunc(kpi, &cast_val);
      }
    }
  }
  log_trace("Metrics:: CPU core KPI pushed to server.");
  return RETURN_OK;
}

/* Push CPU stats to the Metric server */
int sys_cpu_push_stat_to_metric_server(MetricsCatConfig *cpuCfg,
                                       SysCPUMetrics *cstat, int coreCount,
                                       metricAddFunc addFunc) {

  /* Push Metrics */
  unsigned long long int val = 0;

  SysCPUMetrics *cpuStat = NULL;
  /* CPU stats are always stored at index 0. */
  cpuStat = &cstat[0];
  if (!cpuStat) {
    return RETURN_NOTOK;
  }

  /* CPU count */
  if (sys_read_and_push_cpu_count(cpuCfg, coreCount, addFunc) != RETURN_OK) {
    log_error("Metrics:: Failed to collect cpu_core count kpi info.");
  }

  /* CPU usage */
  if (sys_read_and_push_cpu_usage(cpuCfg, cpuStat, addFunc) != RETURN_OK) {
    log_error("Metrics:: Failed to collect cpu_usage kpi info.");
  }

  /* Start Collecting other KPI */
  for (int idx = 0; idx < cpuCfg->kpiCount; idx++) {
    KPIConfig *kpi = &(cpuCfg->kpi[idx]);

    if (sys_cpu_get_kpi_value(kpi, cpuStat, &val) != RETURN_OK) {
      continue;
    }

    /* Add KPI to server*/
    double castVal = (double)val;
    addFunc(kpi, &castVal);
  }

  log_trace("Metrics:: CPU KPI pushed to server.");
  return RETURN_OK;
}

/* Read frequency */
int sys_cpu_read_freq(SysCPUMetrics *cpuStat, int cpuMax) {
  int ret = RETURN_OK;
  FILE *fp;
  SysCPUMetrics *tmpCpu;
  char line[1024];
  unsigned int coreId = 0;
  double decfreq;
  int idx = 0;

  if ((fp = fopen(PROC_CPU_INFO, "r")) == NULL)
    return 0;

  while (fgets(line, sizeof(line), fp) != NULL) {

    if (!strncmp(line, "processor\t", 10)) {
      sscanf(strchr(line, ':') + 1, "%u", &coreId);
      if (coreId > (cpuMax - 1)) {
        return RETURN_NOTOK;
        break;
      }
    }

    else if (!strncmp(line, "cpu MHz\t", 8)) {
      char *pos = strchr(line, ':');
      sscanf(pos + 1, "%lf", &decfreq);

      /* Store frequency. Note In CPU stat first entry is for CPU, second for
       * the core 0 and so on. */
      tmpCpu = cpuStat + coreId + 1;
      tmpCpu->freq = decfreq;
      idx++;
      if (idx > (cpuMax - 1)) {
        break;
      }
    }
  }

  return ret;
}

/* Collect CPU stats */
int sys_cpu_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {
  int ret = RETURN_OK;
  int coreCount = 0;

  coreCount = sys_cpu_read_count(false);
  if (coreCount <= 0) {
    log_error("Metrics:: Failed to read processor count.", coreCount);
    return ret;
  }

  SysCPUMetrics *cpuStat = calloc(coreCount + 1, sizeof(SysCPUMetrics));
  if (!cpuStat) {
    log_error("Metrics:: Failed to allocate memory for %d cpu stat.",
              coreCount);
    return RETURN_NOTOK;
  }

  ret = sys_cpu_read_stat(cpuStat, coreCount);
  if (ret != RETURN_OK) {
    log_error("Metrics:: Failed to read processor stats.");
    free(cpuStat);
    return ret;
  }

  ret = sys_cpu_read_freq(cpuStat, coreCount);
  if (ret != RETURN_OK) {
    log_error("Metrics:: Failed to read processor frequency.");
    free(cpuStat);
    return ret;
  }

  if (!strcmp(cfgStat->source, CPU_SOURCE)) {
    ret = sys_cpu_push_stat_to_metric_server(cfgStat, cpuStat, coreCount,
                                             addFunc);
  } else if (!strcmp(cfgStat->source, CORE_SOURCE)) {
    ret = sys_cpu_core_push_stat_to_metric_server(cfgStat, cpuStat,
                                                  coreCount, addFunc);
  } else {
    log_error("Metrics: %s source for the CPU metrics not implemented.",
              cfgStat->source);
    ret = RETURN_NOTOK;
  }

  free(cpuStat);

  return ret;
}
