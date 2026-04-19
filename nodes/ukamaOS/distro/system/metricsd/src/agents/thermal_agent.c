/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "config.h"
#include "metrics.h"

#include "usys_log.h"

static int thermal_read_value(char *path, double *val) {

    FILE *fp          = NULL;
    char line[64]     = {'\0'};
    double raw        = 0.0;

    if ((path == NULL) || (val == NULL)) {
        return RETURN_NOTOK;
    }

    fp = fopen(path, "r");
    if (fp == NULL) {
        usys_log_error("thermal_agent: cannot open %s: %s",
                       path, strerror(errno));
        return RETURN_NOTOK;
    }

    if (fgets(line, sizeof(line), fp) == NULL) {
        fclose(fp);
        usys_log_error("thermal_agent: failed to read %s", path);
        return RETURN_NOTOK;
    }

    fclose(fp);

    if (sscanf(line, "%lf", &raw) != 1) {
        usys_log_error("thermal_agent: invalid temperature in %s", path);
        return RETURN_NOTOK;
    }

    /*
     * Linux thermal sysfs commonly reports millidegree Celsius.
     * Example: 42375 => 42.375 C
     */
    if (raw > 1000.0) {
        raw = raw / 1000.0;
    }

    *val = raw;

    return RETURN_OK;
}

int thermal_collect_stat(MetricsCatConfig *stat, metricAddFunc addFunc) {

    int idx = 0;

    if ((stat == NULL) || (addFunc == NULL) || (stat->url == NULL)) {
        return RETURN_NOTOK;
    }

    for (idx = 0; idx < stat->kpiCount; idx++) {
        int pathLen   = 0;
        double val    = 0.0;
        char *path    = NULL;
        KPIConfig *kpi = &(stat->kpi[idx]);

        if ((kpi == NULL) || (kpi->ext == NULL)) {
            continue;
        }

        pathLen = strlen(stat->url) + strlen(kpi->ext) + 1;
        path = calloc(pathLen, sizeof(char));
        if (path == NULL) {
            continue;
        }

        snprintf(path, pathLen, "%s%s", stat->url, kpi->ext);

        if (thermal_read_value(path, &val) != RETURN_OK) {
            usys_log_error("thermal_agent: failed for source %s path %s",
                           stat->source, path);
            free(path);
            continue;
        }

        addFunc(kpi, &val);
        free(path);
    }

    return RETURN_OK;
}
