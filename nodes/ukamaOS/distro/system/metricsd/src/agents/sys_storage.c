/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <sys/stat.h>
#include <sys/statvfs.h>

#include "usys_log.h"

#include "metrics.h"
#include "sys_stat.h"

/* translate bytes to MB */
static double translate_mem_value(unsigned long long mem) {
    return ((double)mem) / (1024 * 1024);
}

/* verify if the path exists */
int verify_path(const char *path) {
    int ret         = RETURN_NOTOK;
    struct stat sb;

    if (stat(path, &sb) == -1) {
        return ret;
    }

    if (!(sb.st_mode & S_IFDIR)) {
        usys_log_error("%s is not a directory", path);
        return ret;
    }

    return RETURN_OK;
}

/* read storage stats */
long sys_read_storage_stats(const char *path, SysStorageMetrics *sysSt) {
    int ret               = RETURN_NOTOK;
    struct statvfs statFs;

    if (path == NULL) {
        usys_log_error("no storage device path specified");
        return ret;
    }

    if (verify_path(path) != RETURN_OK) {
        return ret;
    }

    if (statvfs(path, &statFs) != 0) {
        return ret;
    }

    sysSt->blksize = statFs.f_bsize;
    sysSt->total   = statFs.f_bsize * statFs.f_blocks;
    sysSt->free    = statFs.f_bsize * statFs.f_bavail;
    sysSt->pfree   = statFs.f_bsize * statFs.f_bfree;
    sysSt->used    = sysSt->total - sysSt->pfree;

    return RETURN_OK;
}

/* collect and add storage stats */
int sys_storage_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                           SysStorageMetrics *storageStat,
                                           metricAddFunc addFunc) {
    int ret                = RETURN_OK;
    int idx                = 0;
    unsigned long long val = 0;

    for (idx = 0; idx < cfgStat->kpiCount; idx++) {
        KPIConfig *kpi = &(cfgStat->kpi[idx]);

        if ((kpi != NULL) && (kpi->fqname != NULL)) {
            if (strstr(kpi->fqname, "total")) {
                val = storageStat->total;
            } else if (strstr(kpi->fqname, "used")) {
                val = storageStat->used;
            } else if (strstr(kpi->fqname, "free")) {
                val = storageStat->free;
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

/* collect storage stats */
int sys_storage_collect_stat(MetricsCatConfig *cfgStat,
                             metricAddFunc addFunc) {
    int ret                           = RETURN_OK;
    SysStorageMetrics *storageStat    = NULL;

    storageStat = calloc(1, sizeof(SysStorageMetrics));
    if (storageStat == NULL) {
        usys_log_error("failed to allocate memory for storage stat");
        return RETURN_NOTOK;
    }

    if (sys_read_storage_stats(cfgStat->url, storageStat) != RETURN_OK) {
        usys_log_error("failed to collect storage stats");
        free(storageStat);
        return RETURN_NOTOK;
    } else if (sys_storage_push_stat_to_metric_server(cfgStat, storageStat,
                                                      addFunc) != RETURN_OK) {
        usys_log_error("failed to add storage stats to metric server");
        ret = RETURN_NOTOK;
    }

    free(storageStat);

    return ret;
}
