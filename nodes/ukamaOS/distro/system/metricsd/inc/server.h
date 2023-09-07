/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef METRICSERVER_INC_METRIC_H_
#define METRICSERVER_INC_METRIC_H_

#include "config.h"
#include "metrics.h"

int metric_server_add_kpi_data(KPIConfig *kpi, void *value);
int metric_server_register_kpi(KPIConfig *kpi);
int metric_server_start(int port);

void metric_server_registry_init();
void metric_server_registry_destroy();
void metric_server_set_active_registry();
void metric_server_stop();

#endif /* METRICSERVER_INC_METRIC_H_ */
