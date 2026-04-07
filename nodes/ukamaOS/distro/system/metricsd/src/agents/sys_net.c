/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "usys_log.h"

#include "metrics.h"
#include "sys_stat.h"

/* iface count */
int sys_net_dev_if_count(int *ifCount) {
    int ret      = RETURN_OK;
    FILE *fp     = NULL;
    char line[128];

    if ((fp = fopen(PROC_NET_DEV, "r")) == NULL) {
        return RETURN_NOTOK;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        if (strchr(line, ':') != NULL) {
            (*ifCount)++;
        }
    }

    fclose(fp);

    return ret;
}

/* stats for net devices */
int sys_net_dev_read_stat(SysNetDevMetrics *netDev, int maxDev,
                          char *ifName) {
    int ret      = RETURN_NOTOK;
    FILE *fp     = NULL;
    char line[256];

    (void)maxDev;

    if ((fp = fopen(PROC_NET_DEV, "r")) == NULL) {
        return RETURN_NOTOK;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        if (!strncmp(line, ifName, strlen(ifName))) {
            sscanf(line + strlen(ifName) + 2,
                   "%llu %llu %llu %llu %llu %llu %llu %llu %llu %llu "
                   "%llu %llu %llu %llu %llu %llu",
                   &netDev->rxBytes, &netDev->rxPackets,
                   &netDev->rxErrors, &netDev->rxDropped,
                   &netDev->rxFifoErrors, &netDev->rxOverruns,
                   &netDev->rxCompressed, &netDev->multicast,
                   &netDev->txBytes, &netDev->txPackets,
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

/* push net stats to the metric server */
int sys_net_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                       SysNetDevMetrics *nstat,
                                       metricAddFunc addFunc) {
    int ret                = RETURN_OK;
    int idx                = 0;
    unsigned long long val = 0;

    for (idx = 0; idx < cfgStat->kpiCount; idx++) {
        KPIConfig *kpi = &(cfgStat->kpi[idx]);

        if ((kpi != NULL) && (kpi->fqname != NULL)) {
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

            {
                double castVal = (double)val;
                addFunc(kpi, &castVal);
            }
        }
    }

    return ret;
}

/* collect network stats */
int sys_net_collect_stat(MetricsCatConfig *nstat, metricAddFunc addFunc) {
    int ret                       = RETURN_OK;
    int ifCount                   = 1;
    SysNetDevMetrics *netStat     = NULL;

    netStat = calloc(ifCount, sizeof(SysNetDevMetrics));
    if (netStat == NULL) {
        usys_log_error("failed to allocate memory for %d network stat",
                       ifCount);
        return RETURN_NOTOK;
    }

    ret = sys_net_dev_read_stat(netStat, ifCount, nstat->source);
    if (ret != RETURN_OK) {
        usys_log_error("failed to read network stats");
        free(netStat);
        return ret;
    }

    ret = sys_net_push_stat_to_metric_server(nstat, netStat, addFunc);

    free(netStat);

    return ret;
}
