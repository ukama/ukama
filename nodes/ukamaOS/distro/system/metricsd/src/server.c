/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>

#include "config.h"
#include "microhttpd.h"
#include "prom.h"
#include "promhttp.h"

#include "usys_log.h"

static struct MHD_Daemon *mhdDaemon = NULL;

void metric_server_free_kpi(KPIConfig *kpi) {

    if ((kpi == NULL) || (kpi->registry == NULL)) {
        return;
    }

    switch (kpi->type) {
    case METRICTYPE_COUNTER:
        if (kpi->registry->counter != NULL) {
            prom_counter_destroy(kpi->registry->counter);
        }
        break;

    case METRICTYPE_GAUGE:
        if (kpi->registry->gauge != NULL) {
            prom_gauge_destroy(kpi->registry->gauge);
        }
        break;

    case METRICTYPE_HISTOGRAM:
        if (kpi->registry->histogram != NULL) {
            prom_histogram_destroy(kpi->registry->histogram);
        }
        break;

    default:
        usys_log_error("invalid kpi type %d for %s", kpi->type, kpi->name);
        break;
    }

    free(kpi->registry);
    kpi->registry = NULL;
}

int metric_server_register_kpi(KPIConfig *kpi) {

    usys_log_trace("registering prometheus metric for %s", kpi->fqname);

    kpi->registry = calloc(1, sizeof(PromRegistry));
    if (kpi->registry == NULL) {
        return RETURN_NOTOK;
    }

    switch (kpi->type) {
    case METRICTYPE_COUNTER:
        kpi->registry->counter =
            prom_collector_registry_must_register_metric(
                prom_counter_new(kpi->fqname, kpi->desc,
                                 kpi->numLabels,
                                 (const char **)kpi->labels));
        if (kpi->registry->counter == NULL) {
            return RETURN_NOTOK;
        }
        break;

    case METRICTYPE_GAUGE:
        kpi->registry->gauge =
            prom_collector_registry_must_register_metric(
                prom_gauge_new(kpi->fqname, kpi->desc, kpi->numLabels,
                               (const char **)kpi->labels));
        if (kpi->registry->gauge == NULL) {
            return RETURN_NOTOK;
        }
        break;

    case METRICTYPE_HISTOGRAM:
        kpi->registry->histogram =
            prom_collector_registry_must_register_metric(
                prom_histogram_new(kpi->fqname, kpi->desc, kpi->buckets,
                                   kpi->numLabels,
                                   (const char **)kpi->labels));
        if (kpi->registry->histogram == NULL) {
            return RETURN_NOTOK;
        }
        break;

    default:
        usys_log_error("invalid kpi type %d for %s", kpi->type, kpi->name);
        return RETURN_NOTOK;
    }

    return RETURN_OK;
}

int metric_server_add_kpi_data(KPIConfig *kpi, void *value) {

    int ret = RETURN_OK;

    usys_log_trace("adding kpi %s type %d", kpi->name, kpi->type);

    switch (kpi->type) {
    case METRICTYPE_COUNTER:
        ret = prom_counter_inc(kpi->registry->counter,
                               (const char **)kpi->labels);
        break;

    case METRICTYPE_GAUGE:
        ret = prom_gauge_set(kpi->registry->gauge, *(double *)value,
                             (const char **)kpi->labels);
        break;

    case METRICTYPE_HISTOGRAM:
        ret = prom_histogram_observe(kpi->registry->histogram,
                                     *(double *)value,
                                     (const char **)kpi->labels);
        break;

    default:
        ret = RETURN_NOTOK;
        usys_log_error("invalid kpi type %d for %s", kpi->type, kpi->name);
        break;
    }

    return ret;
}

void metric_server_registry_init(void) {

    prom_collector_registry_default_init();
}

void metric_server_set_active_registry(void) {

    promhttp_set_active_collector_registry(NULL);
}

void metric_server_registry_destroy(void) {

    prom_collector_registry_destroy(PROM_COLLECTOR_REGISTRY_DEFAULT);
}

int metric_server_start(int port) {

    mhdDaemon = promhttp_start_daemon(MHD_USE_SELECT_INTERNALLY, port,
                                      NULL, NULL);
    if (mhdDaemon == NULL) {
        return RETURN_NOTOK;
    }

    return RETURN_OK;
}

void metric_server_stop(void) {

    if (mhdDaemon != NULL) {
        MHD_stop_daemon(mhdDaemon);
        mhdDaemon = NULL;
    }
}
