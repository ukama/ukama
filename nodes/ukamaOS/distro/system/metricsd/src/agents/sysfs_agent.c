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
#include <unistd.h>

#include "file.h"
#include "metrics.h"

#include "usys_log.h"

#define PS_MAX_LENGTH_NUMBER 20
#define PS_DEF_OFFSET        0

/* raw read from sysfs file */
int ps_read_block(char *name, void *buff, uint16_t size) {
    int readBytes = 0;
    int fd        = 0;

    fd = file_open(name, O_RDONLY);
    if (fd < 0) {
        readBytes = -1;
    } else {
        lseek(fd, PS_DEF_OFFSET, SEEK_SET);
        readBytes = read(fd, buff, size);
        file_close(fd);
    }

    usys_log_trace("read %d bytes from %s at offset 0x%x", readBytes,
                   name, PS_DEF_OFFSET);

    return readBytes;
}

/* check if sysfs file exists */
int sysfs_check_for_kpi_source(char *source) {
    int ret = RETURN_OK;

    if (!file_exist(source)) {
        ret = RETURN_NOTOK;
    }

    return ret;
}

/* read kpi data from the sysfs file */
int sysfs_read_kpi_data(char *source, double *nval) {
    int ret      = RETURN_OK;
    FILE *fp     = NULL;
    char line[32];

    if ((fp = fopen(source, "r")) == NULL) {
        usys_log_error("cannot open %s: %s", source, strerror(errno));
        return RETURN_NOTOK;
    }

    if (fgets(line, sizeof(line), fp) != NULL) {
        sscanf(line, "%lf", nval);
    }

    if (fp != NULL) {
        fclose(fp);
    }

    return ret;
}

int sysfs_push_kpi_metric_server(KPIConfig *kpi, char *source,
                                 metricAddFunc addFunc) {
    int ret    = RETURN_NOTOK;
    double val = 0;

    if (sysfs_check_for_kpi_source(source) != RETURN_OK) {
        usys_log_error("source %s missing for kpi %s", source, kpi->name);
        return ret;
    }

    if (sysfs_read_kpi_data(source, &val) != RETURN_OK) {
        usys_log_error("failed to read kpi %s from file %s", kpi->name,
                       source);
        return ret;
    }

    addFunc(kpi, &val);

    return RETURN_OK;
}

/* collect kpi data from sysfs files */
int sysfs_collect_kpi(MetricsCatConfig *stat, metricAddFunc addFunc) {
    int ret = RETURN_NOTOK;
    int idx = 0;

    for (idx = 0; idx < stat->kpiCount; idx++) {
        int length   = 0;
        char *source = NULL;

        length = strlen(stat->url) + strlen(stat->kpi[idx].ext) + 1;
        source = calloc(1, length);
        if (source == NULL) {
            continue;
        }

        strcpy(source, stat->url);
        strcat(source, stat->kpi[idx].ext);

        if (sysfs_push_kpi_metric_server(&(stat->kpi[idx]), source,
                                         addFunc) != RETURN_OK) {
            usys_log_error("failed to push data for kpi %s from source %s",
                           stat->kpi[idx].name, source);
        }

        free(source);
        ret = RETURN_OK;
    }

    return ret;
}
