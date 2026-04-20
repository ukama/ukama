/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "usys_log.h"

#include "metrics.h"
#include "sys_stat.h"

#define MEM_SOURCE_DDR  "ddr"
#define MEM_SOURCE_SWAP "swap"

/* translate kB to MB */
static double translate_mem_value(unsigned long long mem) {
    return ((double)mem) / 1024;
}

/* read mem stats */
int sys_read_memStat(SysMemMetrics *statMem) {
    int ret      = RETURN_OK;
    FILE *fp     = NULL;
    char line[128];

    if ((fp = fopen(PROC_MEM_STAT, "r")) == NULL) {
        usys_log_error("cannot open %s: %s", PROC_MEM_STAT,
                       strerror(errno));
        return RETURN_NOTOK;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        if (!strncmp(line, "MemTotal:", 9)) {
            sscanf(line + 9, "%llu", &statMem->ddr.memTotal);
        } else if (!strncmp(line, "MemFree:", 8)) {
            sscanf(line + 8, "%llu", &statMem->ddr.memFree);
        } else if (!strncmp(line, "MemAvailable:", 13)) {
            sscanf(line + 13, "%llu", &statMem->ddr.memAvail);
        } else if (!strncmp(line, "Buffers:", 8)) {
            sscanf(line + 8, "%llu", &statMem->ddr.memBuffer);
        } else if (!strncmp(line, "Cached:", 7)) {
            sscanf(line + 7, "%llu", &statMem->ddr.memCached);
        } else if (!strncmp(line, "SwapTotal:", 10)) {
            sscanf(line + 10, "%llu", &statMem->swap.total);
        } else if (!strncmp(line, "SwapFree:", 9)) {
            sscanf(line + 9, "%llu", &statMem->swap.free);
        }
    }

    statMem->ddr.memUsed = statMem->ddr.memTotal - statMem->ddr.memFree -
                           statMem->ddr.memBuffer - statMem->ddr.memCached;
    statMem->swap.used = statMem->swap.total - statMem->swap.free;

    fclose(fp);

    return ret;
}

/* collect and add ddr mem stats */
int sys_mem_ddr_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                           SysMemDDRMetrics *ddrStat,
                                           metricAddFunc addFunc) {
    int ret                     = RETURN_OK;
    int idx                     = 0;
    unsigned long long val      = 0;

    for (idx = 0; idx < cfgStat->kpiCount; idx++) {
        KPIConfig *kpi = &(cfgStat->kpi[idx]);

        if ((kpi != NULL) && (kpi->fqname != NULL)) {
            if (strstr(kpi->fqname, "total")) {
                val = ddrStat->memTotal;
            } else if (strstr(kpi->fqname, "used")) {
                val = ddrStat->memUsed;
            } else if (strstr(kpi->fqname, "free")) {
                val = ddrStat->memFree;
            } else {
                continue;
            }

            {
                double castVal = translate_mem_value(val);
                addFunc(kpi, &castVal);
            }
        }
    }

    return ret;
}

/* collect and add swap mem stats */
int sys_mem_swap_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                            SysMemSwapMetrics *swapStat,
                                            metricAddFunc addFunc) {
    int ret                     = RETURN_OK;
    int idx                     = 0;
    unsigned long long val      = 0;

    for (idx = 0; idx < cfgStat->kpiCount; idx++) {
        KPIConfig *kpi = &(cfgStat->kpi[idx]);

        if ((kpi != NULL) && (kpi->fqname != NULL)) {
            if (strstr(kpi->fqname, "total")) {
                val = swapStat->total;
            } else if (strstr(kpi->fqname, "used")) {
                val = swapStat->used;
            } else if (strstr(kpi->fqname, "free")) {
                val = swapStat->free;
            } else {
                continue;
            }

            {
                double castVal = translate_mem_value(val);
                addFunc(kpi, &castVal);
            }
        }
    }

    return ret;
}

/* collect memory stats */
int sys_mem_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {
    int ret                   = RETURN_OK;
    SysMemMetrics *memStat    = NULL;

    memStat = calloc(1, sizeof(SysMemMetrics));
    if (memStat == NULL) {
        usys_log_error("failed to allocate memory for mem stat");
        return RETURN_NOTOK;
    }

    if (sys_read_memStat(memStat) != RETURN_OK) {
        usys_log_error("failed to read memory stats");
        free(memStat);
        return RETURN_NOTOK;
    }

    if (!strcmp(cfgStat->source, MEM_SOURCE_DDR)) {
        if (sys_mem_ddr_push_stat_to_metric_server(cfgStat, &memStat->ddr,
                                                   addFunc) != RETURN_OK) {
            usys_log_error("failed to add %s memory stats to metric server",
                           MEM_SOURCE_DDR);
            ret = RETURN_NOTOK;
        }
    } else if (!strcmp(cfgStat->source, MEM_SOURCE_SWAP)) {
        if (sys_mem_swap_push_stat_to_metric_server(cfgStat, &memStat->swap,
                                                    addFunc) != RETURN_OK) {
            usys_log_error("failed to add memory stats %s to metric server",
                           MEM_SOURCE_SWAP);
            ret = RETURN_NOTOK;
        }
    } else {
        usys_log_error("failed to find memory stats source %s",
                       cfgStat->source);
        ret = RETURN_NOTOK;
    }

    free(memStat);

    return ret;
}
