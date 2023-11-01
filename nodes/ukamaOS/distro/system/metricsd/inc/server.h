/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
