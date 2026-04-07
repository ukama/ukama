/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <ctype.h>
#include <stdlib.h>
#include <string.h>

#include "usys_log.h"

#include "agents.h"
#include "collector.h"
#include "file.h"
#include "metrics.h"

/* check if file exists */
int fp_check_for_kpi_source(char *source) {
    int ret = RETURN_OK;

    if (!file_exist(source)) {
        ret = RETURN_NOTOK;
    }

    return ret;
}

/*
 * check if string s2 is a prefix followed by a space in string s1.
 * on failure return 0, otherwise return 1.
 */
int fp_is_prefix(char *s1, char *s2) {
    int ret     = 1;
    size_t n1   = strlen(s1);
    size_t n2   = strlen(s2);
    unsigned int idx = 0;

    if (n1 > n2) {
        for (idx = 0; idx < n2; idx++) {
            if (tolower(s1[idx]) != tolower(s2[idx])) {
                ret = 0;
                break;
            }
        }

        if ((ret == 1) && (s1[n2] != ' ')) {
            ret = 0;
        }
    } else {
        ret = 0;
    }

    return ret;
}

/* parse kpi file line */
KPIData *fp_parse_kpi(KPIConfig *kpi, int count, char *kpiData) {
    KPIData *kdata      = NULL;
    unsigned int idx    = 0;

    for (idx = 0; idx < (unsigned int)count; idx++) {
        if (fp_is_prefix(kpiData, kpi[idx].name)) {
            int dataOffset = 0;

            usys_log_trace("match found for kpi %s", kpi[idx].name);

            kdata = calloc(1, sizeof(KPIData));
            if (kdata == NULL) {
                break;
            }

            kdata->kpi = &kpi[idx];
            dataOffset = strlen(kpi[idx].name);
            kdata->value = atof(&kpiData[dataOffset]);
            break;
        }
    }

    return kdata;
}

/* read kpi file */
int fp_read_kpi_from_file(MetricsCatConfig *stat, metricAddFunc addFunc) {
    int ret      = RETURN_NOTOK;
    FILE *fp     = NULL;
    char *line   = NULL;
    size_t len   = 0;
    ssize_t read = 0;

    if (fp_check_for_kpi_source(stat->url) != RETURN_OK) {
        usys_log_error("file %s does not exist", stat->url);
        return RETURN_NOTOK;
    }

    fp = fopen(stat->url, "r");
    if (fp == NULL) {
        usys_log_error("file %s does not exist", stat->url);
        return RETURN_NOTOK;
    }

    while ((read = getline(&line, &len, fp)) != -1) {
        KPIData *kdata = NULL;

        (void)read;
        kdata = fp_parse_kpi(stat->kpi, stat->kpiCount, line);
        if (kdata == NULL) {
            continue;
        }

        ret = addFunc(kdata->kpi, &kdata->value);
        if (ret != RETURN_OK) {
            usys_log_error("failed to add kpi for %s", kdata->kpi->fqname);
            ret = RETURN_NOTOK;
        } else {
            usys_log_trace("added kpi %s value %lf", kdata->kpi->fqname,
                           kdata->value);
            ret = RETURN_OK;
        }

        free(kdata);
    }

    if (fp != NULL) {
        fclose(fp);
    }

    if (line != NULL) {
        free(line);
    }

    return ret;
}
