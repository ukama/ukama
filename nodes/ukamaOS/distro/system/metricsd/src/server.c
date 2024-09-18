/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include "config.h"
#include "microhttpd.h"
#include "prom.h"
#include "promhttp.h"

#include "usys_log.h"

static struct MHD_Daemon *mhd_daemon;

/* Free the KPI config and registry */
void metric_server_free_kpi(KPIConfig *kpi) {
  if (kpi) {
      if (kpi->registry) {
          switch (kpi->type) {
          case METRICTYPE_COUNTER:
              if (kpi->registry->counter) {
                  prom_counter_destroy(kpi->registry->counter);
              }
              break;

          case METRICTYPE_GAUGE:
              if (kpi->registry->gauge) {
                  prom_gauge_destroy(kpi->registry->gauge);
              }
              break;

          case METRICTYPE_HISTOGRAM:
              if (kpi->registry->histogram) {
                  prom_histogram_destroy(kpi->registry->histogram);
              }
              break;

          default:
              usys_log_error("Invalid KPI type %d for KPI %s.",
                             kpi->type, kpi->name);
              break;
          }

          free(kpi->registry);
          kpi->registry = NULL;
      }
  }
}

/* Initialize metric */
int metric_server_register_kpi(KPIConfig *kpi) {

  usys_log_trace("Initializing prometheus metric  for KPI %s \n",
                 kpi->fqname);

  kpi->registry = calloc(1, sizeof(PromRegistry));
  if (!kpi->registry) {
    return RETURN_NOTOK;
  }
  switch (kpi->type) {

  case METRICTYPE_COUNTER:
    /* Create Counter */
    kpi->registry->counter = prom_collector_registry_must_register_metric(
        prom_gauge_new(kpi->fqname, kpi->desc, kpi->numLabels,
                       (const char **)kpi->labels));
    break;

  case METRICTYPE_GAUGE:
    /* Create gauge */
    kpi->registry->gauge = prom_collector_registry_must_register_metric(
        prom_gauge_new(kpi->fqname, kpi->desc, kpi->numLabels,
                       (const char **)kpi->labels));
    break;

  case METRICTYPE_HISTOGRAM:
    /* Create histogram */
    kpi->registry->histogram = prom_collector_registry_must_register_metric(
        prom_histogram_new(kpi->fqname, kpi->desc, kpi->buckets, kpi->numLabels,
                           (const char **)kpi->labels));
    break;

  default:
    usys_log_error("Invalid KPI type %d for KPI %s.", kpi->type,
                   kpi->name);
  }

  if (!kpi->registry->counter) {
    return RETURN_NOTOK;
  }

  return RETURN_OK;
}

/* Add metric value */
int metric_server_add_kpi_data(KPIConfig *kpi, void *value) {
  int ret = RETURN_OK;

  usys_log_trace("Adding KPI %s  %s and type %d", kpi->name, kpi->desc,
                 kpi->type);
  switch (kpi->type) {

  case METRICTYPE_COUNTER:
    /* Increment counter */
    ret = prom_counter_inc(kpi->registry->counter, (const char **)kpi->labels);
    break;

  case METRICTYPE_GAUGE:
    /* Add gauge value*/
    ret = prom_gauge_set(kpi->registry->gauge, *(double *)value,
                         (const char **)kpi->labels);
    break;

  case METRICTYPE_HISTOGRAM:
    /* Add Histogram Value */
    ret = prom_histogram_observe(kpi->registry->histogram, *(double *)value,
                                 (const char **)kpi->labels);
    break;

  default:
    ret = RETURN_NOTOK;
    usys_log_error("Invalid KPI type %d for KPI %s.", kpi->type,
                   kpi->name);
  }

  return ret;
}

/* Initialize the Default registry */
void metric_server_registry_init() { prom_collector_registry_default_init(); }

/* Set the active registry for the HTTP handler */
void metric_server_set_active_registry() {
  promhttp_set_active_collector_registry(NULL);
}

void metric_server_registry_destroy() {
  prom_collector_registry_destroy(PROM_COLLECTOR_REGISTRY_DEFAULT);
}

int metric_server_start(int port) {

  /* Star HTTP server daemon */
  mhd_daemon =
      promhttp_start_daemon(MHD_USE_SELECT_INTERNALLY, port, NULL, NULL);
  if (mhd_daemon == NULL) {
    return RETURN_NOTOK;
  }

  return RETURN_OK;
}

void metric_server_stop() {
  /* Stop HTTP server daemon */
  MHD_stop_daemon(mhd_daemon);
}
