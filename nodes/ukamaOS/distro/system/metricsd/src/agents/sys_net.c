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

/* Iface count */
int sys_net_dev_if_count(int *ifCount) {
  int ret = RETURN_OK;
  FILE *fp;
  char line[128];

  if ((fp = fopen(PROC_NET_DEV, "r")) == NULL) {
    return RETURN_NOTOK; /* No network device file */
  }
  while (fgets(line, sizeof(line), fp) != NULL) {
    if (strchr(line, ':')) {
      (*ifCount)++;
    }
  }

  fclose(fp);

  return ret;
}

/* Stats for net devices */
int sys_net_dev_read_stat(SysNetDevMetrics *netDev, int maxDev,
                          char *ifName) {
  int ret = RETURN_NOTOK;
  FILE *fp;
  char line[256];
  if ((fp = fopen(PROC_NET_DEV, "r")) == NULL)
    return RETURN_NOTOK;

  while (fgets(line, sizeof(line), fp) != NULL) {
    if (!strncmp(line, ifName, strlen(ifName))) {
      sscanf(line + strlen(ifName) + 2,
             "%llu %llu %llu %llu %llu %llu %llu %llu %llu %llu "
             "%llu %llu %llu %llu %llu %llu",
             &netDev->rxBytes, &netDev->rxPackets, &netDev->rxErrors,
             &netDev->rxDropped, &netDev->rxFifoErrors,
             &netDev->rxOverruns, &netDev->rxCompressed,
             &netDev->multicast, &netDev->txBytes, &netDev->txPackets,
             &netDev->txErrors, &netDev->txDropped,
             &netDev->txFifoErrors, &netDev->collisions,
             &netDev->txCarrierErrors, &netDev->txCompressed);
      ret = RETURN_OK;
      break;
    }
  }

  fclose(fp);
  return ret;
}

/* Push Net stats to the Metric server */
int sys_net_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                       SysNetDevMetrics *nstat,
                                       metricAddFunc addFunc) {
  int ret = RETURN_OK;
  unsigned long long val = 0;

  /* Start Collecting KPI */
  for (int idx = 0; idx < (cfgStat->kpiCount); idx++) {
    KPIConfig *kpi = &(cfgStat->kpi[idx]);
    if ((kpi) && (kpi->fqname)) {

      if (strstr(kpi->fqname, "rxBytes")) {
        val = nstat->rxBytes;
      } else if (strstr(kpi->fqname, "rx_errors")) {
        val = nstat->rxErrors;
      } else if (strstr(kpi->fqname, "rx_dropped")) {
        val = nstat->rxDropped;
      } else if (strstr(kpi->fqname, "rx_overruns")) {
        val = nstat->rxOverruns;
      } else if (strstr(kpi->fqname, "rx_packets")) {
        val = nstat->rxPackets;
      } else if (strstr(kpi->fqname, "tx_bytes")) {
        val = nstat->txBytes;
      } else if (strstr(kpi->fqname, "tx_errors")) {
        val = nstat->txErrors;
      } else if (strstr(kpi->fqname, "tx_dropped")) {
        val = nstat->txDropped;
      } else if (strstr(kpi->fqname, "tx_carrier_errors")) {
        val = nstat->txCarrierErrors;
      } else if (strstr(kpi->fqname, "tx_packets")) {
        val = nstat->txPackets;
      } else if (strstr(kpi->fqname, "link")) {
        val = nstat->linkstatus;
      } else if (strstr(kpi->fqname, "speed")) {
        val = nstat->linkspeed;
      } else if (strstr(kpi->fqname, "latency")) {
        val = nstat->latency;
      } else {
        continue;
      }

      /* Add KPI to server*/
      double cast_val = (double)(val);
      addFunc(kpi, &cast_val);
    }
  }

  return ret;
}

/* Collect Network stats */
int sys_net_collect_stat(MetricsCatConfig *nstat, metricAddFunc addFunc) {
  int ret = RETURN_OK;
  int if_count = 1;

  SysNetDevMetrics *netStat = calloc(if_count, sizeof(SysNetDevMetrics));
  if (!netStat) {
    log_error("Metrics:: Failed to allocate memory for %d network stat.",
              if_count);
    return RETURN_NOTOK;
  }

  ret = sys_net_dev_read_stat(netStat, if_count, nstat->source);
  if (ret != RETURN_OK) {
    log_error("Metrics:: Failed to read network stats.");
    free(netStat);
    return ret;
  }

  ret = sys_net_push_stat_to_metric_server(nstat, netStat, addFunc);

  free(netStat);

  return ret;
}
